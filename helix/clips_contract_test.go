package helix_test

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func TestConsumerCompile_acceptsClipSelectors(t *testing.T) {
	directory := t.TempDir()
	moduleRoot := filepath.Dir(mustWorkingDirectory(t))
	goMod := "module fixture\n\ngo 1.24\n\nrequire github.com/kvizyx/twitchy v0.0.0\n\nreplace github.com/kvizyx/twitchy => " + moduleRoot + "\n"
	source := `package fixture

import (
    "context"
    "github.com/kvizyx/twitchy/helix"
)

func valid(client *helix.Client) {
    _, _ = client.Clips.CreateClip(context.Background(), helix.CreateClipRequest{})
    _, _ = client.Clips.GetClips(context.Background(), helix.GetClipsRequest{})
    _, _ = client.Clips.GetClipsPager(helix.GetClipsRequest{})
    _, _ = client.Experimental.Clips.CreateClipFromVOD(context.Background(), helix.CreateClipFromVODRequest{})
    _, _ = client.Experimental.Clips.GetClipsDownload(context.Background(), helix.GetClipsDownloadRequest{})
}
`
	writeCompileFixture(t, directory, goMod, source)
	if output, err := runCompileFixture(directory); err != nil {
		t.Fatalf("valid clip selectors failed to compile: %v\n%s", err, output)
	}
}

func TestConsumerCompile_rejectsStableAccessToClipNEWMethods(t *testing.T) {
	directory := t.TempDir()
	moduleRoot := filepath.Dir(mustWorkingDirectory(t))
	goMod := "module fixture\n\ngo 1.24\n\nrequire github.com/kvizyx/twitchy v0.0.0\n\nreplace github.com/kvizyx/twitchy => " + moduleRoot + "\n"
	source := `package fixture

import "github.com/kvizyx/twitchy/helix"

func invalid(client *helix.Client) {
    _, _ = client.Clips.CreateClipFromVOD(nil, helix.CreateClipFromVODRequest{})
    _, _ = client.Clips.GetClipsDownload(nil, helix.GetClipsDownloadRequest{})
}
`
	writeCompileFixture(t, directory, goMod, source)
	output, err := runCompileFixture(directory)
	if err == nil {
		t.Fatal("stable clip access to NEW methods unexpectedly compiled")
	}
	if !bytes.Contains(output, []byte("CreateClipFromVOD")) && !bytes.Contains(output, []byte("GetClipsDownload")) && !strings.Contains(string(output), "undefined") {
		t.Fatalf("unexpected compile failure: %s", output)
	}
}

func newTask17Client(transport http.RoundTripper, credential helix.Credential) (*helix.Client, error) {
	return helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(credential),
	)
}

func stringPtr(value string) *string    { return &value }
func float64Ptr(value float64) *float64 { return &value }
func intPtr(value int) *int             { return &value }
func boolPtr(value bool) *bool          { return &value }

func timestamp(value string) helix.Timestamp {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		panic(err)
	}
	return helix.Timestamp{Time: parsed}
}

func mustWorkingDirectory(t *testing.T) string {
	t.Helper()
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return workingDirectory
}

func writeCompileFixture(t *testing.T, directory, goMod, source string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(directory, "go.mod"), []byte(goMod), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "fixture.go"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
}

func runCompileFixture(directory string) ([]byte, error) {
	command := exec.Command("go", "build", ".")
	command.Dir = directory
	command.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOPROXY=off", "GOSUMDB=off", "GOWORK=off")
	return command.CombinedOutput()
}
