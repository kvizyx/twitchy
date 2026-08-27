package oauth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func registryServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/oauth2/validate":
			_, _ = io.WriteString(writer, `{"client_id":"client","login":"streamer","scopes":["user:read:chat"],"user_id":"validated","expires_in":7200}`)
		default:
			writer.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func newRegistry(t *testing.T) *Registry {
	t.Helper()
	server := registryServer(t)
	client, err := New(WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}
	registry, err := NewRegistry(client)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = registry.Close() })
	return registry
}

func noopHook(context.Context, TokenPair) error { return nil }

func registryPair(accessToken string) TokenPair {
	return TokenPair{AccessToken: accessToken, RefreshToken: "refresh-" + accessToken, ExpiresIn: time.Hour}
}

func TestRegistryAddUser_registersManagedCredential(t *testing.T) {
	registry := newRegistry(t)
	ctx := context.Background()
	if err := registry.AddUser(ctx, "111", registryPair("one"), noopHook, "chat"); err != nil {
		t.Fatal(err)
	}
	source, err := registry.SourceForUser("111")
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := source.Token(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got := snapshot.AccessToken(); got != "one" {
		t.Fatalf("AccessToken = %q, want %q", got, "one")
	}
}

func TestRegistryAddUser_rejectsDuplicateAndInvalidInput(t *testing.T) {
	registry := newRegistry(t)
	ctx := context.Background()
	if err := registry.AddUser(ctx, "111", registryPair("one"), noopHook); err != nil {
		t.Fatal(err)
	}
	if err := registry.AddUser(ctx, "111", registryPair("two"), noopHook); !errors.Is(err, ErrUserExists) {
		t.Fatalf("duplicate AddUser error = %v, want ErrUserExists", err)
	}
	if err := registry.AddUser(ctx, "", registryPair("three"), noopHook); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("empty userID error = %v, want ErrInvalidOption", err)
	}
	if err := registry.AddUser(ctx, "222", registryPair("three"), nil); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("nil hook error = %v, want ErrInvalidOption", err)
	}
}

func TestRegistryRemoveUser_closesSession(t *testing.T) {
	registry := newRegistry(t)
	ctx := context.Background()
	if err := registry.AddUser(ctx, "111", registryPair("one"), noopHook); err != nil {
		t.Fatal(err)
	}
	if err := registry.RemoveUser("111"); err != nil {
		t.Fatal(err)
	}
	if _, err := registry.SourceForUser("111"); !errors.Is(err, helix.ErrUserNotFound) {
		t.Fatalf("SourceForUser error = %v, want ErrUserNotFound", err)
	}
	if err := registry.RemoveUser("111"); !errors.Is(err, helix.ErrUserNotFound) {
		t.Fatalf("second RemoveUser error = %v, want ErrUserNotFound", err)
	}
}

func TestRegistrySourceForIntent_coversDeterministically(t *testing.T) {
	registry := newRegistry(t)
	ctx := context.Background()
	if err := registry.AddUser(ctx, "222", registryPair("two"), noopHook, "chat", "eventsub"); err != nil {
		t.Fatal(err)
	}
	if err := registry.AddUser(ctx, "111", registryPair("one"), noopHook, "chat"); err != nil {
		t.Fatal(err)
	}

	source, err := registry.SourceForIntent("chat")
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := source.Token(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got := snapshot.AccessToken(); got != "one" {
		t.Fatalf("chat intent resolved to %q, want %q (sorted user IDs)", got, "one")
	}

	source, err = registry.SourceForIntent("chat", "eventsub")
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err = source.Token(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got := snapshot.AccessToken(); got != "two" {
		t.Fatalf("chat+eventsub intent resolved to %q, want %q", got, "two")
	}

	if _, err := registry.SourceForIntent("whispers"); !errors.Is(err, helix.ErrIntentNotCovered) {
		t.Fatalf("SourceForIntent error = %v, want ErrIntentNotCovered", err)
	}
}

func TestRegistrySourceForIntent_skipsTerminatedSessions(t *testing.T) {
	registry := newRegistry(t)
	ctx := context.Background()
	if err := registry.AddUser(ctx, "111", registryPair("one"), noopHook, "chat"); err != nil {
		t.Fatal(err)
	}
	if err := registry.AddUser(ctx, "222", registryPair("two"), noopHook, "chat"); err != nil {
		t.Fatal(err)
	}

	registry.mu.Lock()
	registry.users["111"].source.terminal = helix.ErrInvalidSession
	registry.mu.Unlock()

	source, err := registry.SourceForIntent("chat")
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := source.Token(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got := snapshot.AccessToken(); got != "two" {
		t.Fatalf("intent resolved to %q, want %q (terminated session skipped)", got, "two")
	}

	if _, err := registry.SourceForUser("111"); !errors.Is(err, helix.ErrInvalidSession) {
		t.Fatalf("SourceForUser error = %v, want ErrInvalidSession", err)
	}
}

func TestRegistryClose_isIdempotentAndRejectsFurtherUse(t *testing.T) {
	registry := newRegistry(t)
	ctx := context.Background()
	if err := registry.AddUser(ctx, "111", registryPair("one"), noopHook, "chat"); err != nil {
		t.Fatal(err)
	}
	if err := registry.Close(); err != nil {
		t.Fatal(err)
	}
	if err := registry.Close(); err != nil {
		t.Fatal(err)
	}
	if err := registry.AddUser(ctx, "222", registryPair("two"), noopHook); !errors.Is(err, ErrRegistryClosed) {
		t.Fatalf("AddUser after close error = %v, want ErrRegistryClosed", err)
	}
	if _, err := registry.SourceForUser("111"); !errors.Is(err, helix.ErrUserNotFound) {
		t.Fatalf("SourceForUser after close error = %v, want ErrUserNotFound", err)
	}
}

func TestRegistry_concurrentAccess(t *testing.T) {
	registry := newRegistry(t)
	ctx := context.Background()
	var wg sync.WaitGroup
	for index := 0; index < 8; index++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			userID := fmt.Sprintf("user-%d", index)
			if err := registry.AddUser(ctx, userID, registryPair(userID), noopHook, "chat"); err != nil {
				t.Error(err)
				return
			}
			if _, err := registry.SourceForUser(userID); err != nil {
				t.Error(err)
			}
			if _, err := registry.SourceForIntent("chat"); err != nil {
				t.Error(err)
			}
		}(index)
	}
	wg.Wait()
}

func TestNewRegistry_clientSecretOptionValidation(t *testing.T) {
	server := registryServer(t)
	client, err := New(WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewRegistry(client, WithClientSecret("")); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("empty client secret error = %v, want ErrInvalidOption", err)
	}
	if _, err := NewRegistry(client, nil); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("nil option error = %v, want ErrInvalidOption", err)
	}
	if _, err := NewCoordinatedRegistry(client, newMemoryRefreshCoordinator(), WithClientSecret("")); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("coordinated empty client secret error = %v, want ErrInvalidOption", err)
	}
}

func TestRegistryAddUser_forwardsClientSecretToRefresh(t *testing.T) {
	clock := newTestClock(time.Unix(80_000, 0))
	secrets := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		if err := request.ParseForm(); err != nil {
			http.Error(writer, "invalid form", http.StatusBadRequest)
			return
		}
		switch request.URL.Path {
		case "/oauth2/validate":
			_, _ = io.WriteString(writer, `{"client_id":"client","login":"streamer","scopes":[],"user_id":"111","expires_in":3600}`)
		case "/oauth2/token":
			secrets <- request.Form.Get("client_secret")
			_, _ = io.WriteString(writer, `{"access_token":"rotated","refresh_token":"rotated-refresh","expires_in":3600,"token_type":"bearer"}`)
		default:
			writer.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	client, err := New(WithBaseURL(server.URL), WithHTTPClient(server.Client()), WithClock(clock))
	if err != nil {
		t.Fatal(err)
	}
	registry, err := NewRegistry(client, WithClientSecret("registry-secret"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = registry.Close() })
	ctx := context.Background()
	if err := registry.AddUser(ctx, "111", registryPair("one"), noopHook, "chat"); err != nil {
		t.Fatal(err)
	}
	source, err := registry.SourceForUser("111")
	if err != nil {
		t.Fatal(err)
	}
	clock.Advance(time.Hour)

	if _, err := source.Token(ctx); err != nil {
		t.Fatal(err)
	}
	select {
	case secret := <-secrets:
		if secret != "registry-secret" {
			t.Fatalf("refresh client_secret = %q, want %q", secret, "registry-secret")
		}
	case <-time.After(time.Second):
		t.Fatal("refresh did not run")
	}
}
