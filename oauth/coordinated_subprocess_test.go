package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type subprocessPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type subprocessResult struct {
	PID         int    `json:"pid"`
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
}

const subprocessChildEnv = "TWITCHY_SUBPROCESS_CHILD"

func TestCoordinatedRefreshSubprocess(t *testing.T) {
	if os.Getenv(subprocessChildEnv) == "1" {
		os.Exit(runCoordinatedSubprocessChild())
	}
	address := os.Getenv("TWITCHY_TEST_REDIS_ADDR")
	if address == "" {
		t.Skip("TWITCHY_TEST_REDIS_ADDR is not set")
	}

	var refreshes atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/oauth2/validate":
			_, _ = io.WriteString(writer, `{"client_id":"client","login":"streamer",`+
				`"scopes":[],"user_id":"42","expires_in":1}`)
		case "/oauth2/token":
			refreshes.Add(1)
			_, _ = io.WriteString(writer, coordinatedRotationResponse)
		default:
			http.Error(writer, "unexpected path", http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	directory := t.TempDir()
	storePath := filepath.Join(directory, "store.json")
	initial := subprocessPair{
		AccessToken:  "subprocess-initial-access",
		RefreshToken: "subprocess-initial-refresh",
		ExpiresAt:    time.Now().Add(time.Minute).UnixNano(),
	}
	if err := writeSubprocessPair(storePath, initial); err != nil {
		t.Fatal(err)
	}
	commands := startSubprocessChildren(t, directory, storePath, server.URL, address)
	openSubprocessStartGate(t, directory)
	results := waitSubprocessChildren(t, directory, commands)
	assertSubprocessOutcome(t, storePath, results, refreshes.Load())
}

func startSubprocessChildren(
	t *testing.T,
	directory string,
	storePath string,
	serverURL string,
	address string,
) []*exec.Cmd {
	t.Helper()
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	leasePrefix := fmt.Sprintf("oauth:subprocess:%d:", time.Now().UnixNano())
	commands := make([]*exec.Cmd, 2)
	for index, role := range []string{"first", "second"} {
		command := exec.Command(self, "-test.run", "^TestCoordinatedRefreshSubprocess$", "-test.count=1")
		command.Env = append(os.Environ(),
			subprocessChildEnv+"=1",
			"TWITCHY_SUBPROCESS_SERVER="+serverURL,
			"TWITCHY_SUBPROCESS_REDIS="+address,
			"TWITCHY_SUBPROCESS_STORE="+storePath,
			"TWITCHY_SUBPROCESS_RESULT="+filepath.Join(directory, role+".json"),
			"TWITCHY_SUBPROCESS_PREFIX="+leasePrefix,
			"TWITCHY_SUBPROCESS_READY="+filepath.Join(directory, role+".ready"),
			"TWITCHY_SUBPROCESS_GO="+filepath.Join(directory, "go"),
		)
		if err := command.Start(); err != nil {
			t.Fatalf("start %s child: %v", role, err)
		}
		commands[index] = command
	}
	return commands
}

func openSubprocessStartGate(t *testing.T, directory string) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for _, role := range []string{"first", "second"} {
		readyPath := filepath.Join(directory, role+".ready")
		for {
			if _, err := os.Stat(readyPath); err == nil {
				break
			}
			if time.Now().After(deadline) {
				t.Fatalf("%s child did not register in time", role)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	if err := os.WriteFile(filepath.Join(directory, "go"), []byte("go"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func waitSubprocessChildren(t *testing.T, directory string, commands []*exec.Cmd) []subprocessResult {
	t.Helper()
	results := make([]subprocessResult, 2)
	for index, role := range []string{"first", "second"} {
		if err := commands[index].Wait(); err != nil {
			t.Fatalf("%s child failed: %v", role, err)
		}
		payload, err := os.ReadFile(filepath.Join(directory, role+".json"))
		if err != nil {
			t.Fatalf("read %s result: %v", role, err)
		}
		if err := json.Unmarshal(payload, &results[index]); err != nil {
			t.Fatalf("decode %s result: %v", role, err)
		}
	}
	return results
}

func assertSubprocessOutcome(t *testing.T, storePath string, results []subprocessResult, refreshes int32) {
	t.Helper()
	if results[0].PID == results[1].PID || results[0].PID == os.Getpid() {
		t.Fatalf("child PIDs = %d/%d, parent = %d", results[0].PID, results[1].PID, os.Getpid())
	}
	for _, result := range results {
		if result.Error != "" {
			t.Fatalf("child %d error: %s", result.PID, result.Error)
		}
		if result.AccessToken != "rotated-access" {
			t.Fatalf("child %d access token = %q, want rotated pair", result.PID, result.AccessToken)
		}
	}
	if refreshes != 1 {
		t.Fatalf("remote refreshes = %d, want 1", refreshes)
	}
	pair, err := readSubprocessPair(storePath)
	if err != nil {
		t.Fatal(err)
	}
	if pair.RefreshToken != "rotated-refresh" {
		t.Fatalf("durable refresh token = %q, want one committed rotation", pair.RefreshToken)
	}
	t.Logf("DISTINCT_PIDS=true REMOTE_REFRESHES=%d DURABLE_COMMITS=1", refreshes)
}

func runCoordinatedSubprocessChild() int {
	ctx := context.Background()
	result := subprocessResult{PID: os.Getpid()}
	defer func() {
		_ = writeSubprocessResult(os.Getenv("TWITCHY_SUBPROCESS_RESULT"), result)
	}()

	client := redis.NewClient(&redis.Options{Addr: os.Getenv("TWITCHY_SUBPROCESS_REDIS")})
	defer client.Close()
	prefix := os.Getenv("TWITCHY_SUBPROCESS_PREFIX")
	coordinator, err := NewRedisRefreshCoordinator(
		client,
		func(userID string) string { return prefix + userID },
		WithRefreshLeaseTTL(300*time.Millisecond),
		WithRefreshLeaseRenewal(100*time.Millisecond),
	)
	if err != nil {
		result.Error = err.Error()
		return 1
	}
	oauthClient, err := New(WithBaseURL(os.Getenv("TWITCHY_SUBPROCESS_SERVER")))
	if err != nil {
		result.Error = err.Error()
		return 1
	}
	registry, err := NewCoordinatedRegistry(oauthClient, coordinator)
	if err != nil {
		result.Error = err.Error()
		return 1
	}
	defer registry.Close()

	storePath := os.Getenv("TWITCHY_SUBPROCESS_STORE")
	loader := func(_ context.Context, _ string) (TokenPair, error) {
		pair, err := readSubprocessPair(storePath)
		if err != nil {
			return TokenPair{}, err
		}
		return TokenPair{
			AccessToken:  pair.AccessToken,
			RefreshToken: pair.RefreshToken,
			ExpiresIn:    time.Until(time.Unix(0, pair.ExpiresAt)),
			TokenType:    "bearer",
		}, nil
	}
	hook := func(_ context.Context, pair TokenPair) error {
		return writeSubprocessPair(storePath, subprocessPair{
			AccessToken:  pair.AccessToken,
			RefreshToken: pair.RefreshToken,
			ExpiresAt:    time.Now().Add(pair.ExpiresIn).UnixNano(),
		})
	}
	if err := registry.AddCoordinatedUser(ctx, "42", loader, hook, "chat"); err != nil {
		result.Error = err.Error()
		return 1
	}
	source, err := registry.SourceForUser("42")
	if err != nil {
		result.Error = err.Error()
		return 1
	}
	if err := os.WriteFile(os.Getenv("TWITCHY_SUBPROCESS_READY"), []byte("ready"), 0o600); err != nil {
		result.Error = err.Error()
		return 1
	}
	goPath := os.Getenv("TWITCHY_SUBPROCESS_GO")
	deadline := time.Now().Add(30 * time.Second)
	for {
		if _, err := os.Stat(goPath); err == nil {
			break
		}
		if time.Now().After(deadline) {
			result.Error = "start gate was not opened"
			return 1
		}
		time.Sleep(10 * time.Millisecond)
	}
	snapshot, err := source.Token(ctx)
	if err != nil {
		result.Error = err.Error()
		return 1
	}
	result.AccessToken = snapshot.AccessToken()
	return 0
}

func writeSubprocessPair(path string, pair subprocessPair) error {
	payload, err := json.Marshal(pair)
	if err != nil {
		return err
	}
	temporary := path + ".tmp"
	if err := os.WriteFile(temporary, payload, 0o600); err != nil {
		return err
	}
	return os.Rename(temporary, path)
}

func readSubprocessPair(path string) (subprocessPair, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return subprocessPair{}, err
	}
	var pair subprocessPair
	return pair, json.Unmarshal(payload, &pair)
}

func writeSubprocessResult(path string, result subprocessResult) error {
	payload, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return os.WriteFile(path, payload, 0o600)
}
