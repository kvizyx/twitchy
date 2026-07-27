package oauth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func TestRefreshingTokenSource_coalescesExpiredRequests(t *testing.T) {
	clock := newTestClock(time.Unix(1_000, 0))
	started := make(chan struct{})
	release := make(chan struct{})
	var refreshes atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/oauth2/validate":
			_, _ = io.WriteString(writer, `{"client_id":"client","login":"streamer","scopes":[],"user_id":"42","expires_in":3600}`)
		case "/oauth2/token":
			refreshes.Add(1)
			close(started)
			<-release
			_, _ = io.WriteString(writer, `{"access_token":"rotated","refresh_token":"rotated-refresh","expires_in":3600,"token_type":"bearer"}`)
		default:
			http.Error(writer, "unexpected path", http.StatusNotFound)
		}
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL), WithClock(clock))
	if err != nil {
		t.Fatal(err)
	}
	pair := TokenPair{AccessToken: "old", RefreshToken: "old-refresh", ExpiresIn: time.Hour, TokenType: "bearer"}
	source, err := NewRefreshingTokenSource(client, pair, func(context.Context, TokenPair) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	defer source.Close()
	clock.Advance(time.Hour)

	type tokenResult struct {
		snapshot helix.CredentialSnapshot
		err      error
	}
	results := make(chan tokenResult, 100)
	for range 100 {
		go func() {
			snapshot, tokenErr := source.Token(context.Background())
			results <- tokenResult{snapshot: snapshot, err: tokenErr}
		}()
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("refresh did not start")
	}
	if got := refreshes.Load(); got != 1 {
		t.Fatalf("refresh calls = %d, want 1", got)
	}
	close(release)
	for range 100 {
		result := <-results
		if result.err != nil {
			t.Fatalf("Token() error = %v", result.err)
		}
		if result.snapshot.AccessToken() != "rotated" || result.snapshot.Generation() != 1 {
			t.Fatalf("snapshot = %q generation %d", result.snapshot.AccessToken(), result.snapshot.Generation())
		}
	}
}

func TestRefreshingTokenSource_hookRejectionRetainsPairForRetryOnly(t *testing.T) {
	clock := newTestClock(time.Unix(2_000, 0))
	var refreshToken string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		body := make([]byte, request.ContentLength)
		_, _ = request.Body.Read(body)
		values, parseErr := url.ParseQuery(string(body))
		if parseErr != nil {
			t.Fatal(parseErr)
		}
		refreshToken = values.Get("refresh_token")
		if request.URL.Path == "/oauth2/validate" {
			_, _ = io.WriteString(writer, `{"client_id":"client","login":"streamer","scopes":[],"user_id":"42","expires_in":3600}`)
			return
		}
		_, _ = io.WriteString(writer, `{"access_token":"rotated","refresh_token":"rotated-refresh","expires_in":3600,"token_type":"bearer"}`)
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL), WithClock(clock))
	if err != nil {
		t.Fatal(err)
	}
	var hooks atomic.Int32
	commitErr := errors.New("store unavailable")
	source, err := NewRefreshingTokenSource(client, TokenPair{AccessToken: "old", RefreshToken: "old-refresh", ExpiresIn: time.Hour}, func(context.Context, TokenPair) error {
		if hooks.Add(1) == 1 {
			return commitErr
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	defer source.Close()
	clock.Advance(time.Hour)
	_, err = source.Token(context.Background())
	if !errors.Is(err, helix.ErrCredentialCommit) || !errors.Is(err, commitErr) {
		t.Fatalf("initial commit error = %v", err)
	}
	if _, err = source.Token(context.Background()); !errors.Is(err, helix.ErrCredentialCommit) {
		t.Fatalf("terminal Token() error = %v", err)
	}
	if refreshToken != "old-refresh" {
		t.Fatalf("refresh token = %q", refreshToken)
	}
	if err = source.RetryCommit(context.Background()); err != nil {
		t.Fatalf("RetryCommit() error = %v", err)
	}
	snapshot, err := source.Token(context.Background())
	if err != nil || snapshot.AccessToken() != "rotated" || snapshot.Generation() != 1 {
		t.Fatalf("published snapshot = %q/%d, error = %v", snapshot.AccessToken(), snapshot.Generation(), err)
	}
}

func TestRefreshingTokenSource_canceledWaiterDoesNotCancelSharedRefresh(t *testing.T) {
	clock := newTestClock(time.Unix(3_000, 0))
	started := make(chan struct{})
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/oauth2/validate" {
			_, _ = io.WriteString(writer, `{"client_id":"client","login":"streamer","scopes":[],"user_id":"42","expires_in":3600}`)
			return
		}
		close(started)
		<-release
		_, _ = io.WriteString(writer, `{"access_token":"rotated","refresh_token":"rotated-refresh","expires_in":3600,"token_type":"bearer"}`)
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL), WithClock(clock))
	if err != nil {
		t.Fatal(err)
	}
	source, err := NewRefreshingTokenSource(client, TokenPair{AccessToken: "old", RefreshToken: "old-refresh", ExpiresIn: time.Hour}, func(context.Context, TokenPair) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	defer source.Close()
	clock.Advance(time.Hour)
	leaderResult := make(chan error, 1)
	go func() {
		_, leaderErr := source.Token(context.Background())
		leaderResult <- leaderErr
	}()
	<-started
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err = source.Token(canceled); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled waiter error = %v", err)
	}
	close(release)
	if err = <-leaderResult; err != nil {
		t.Fatalf("shared refresh error = %v", err)
	}
}

func TestRefreshingTokenSource_doesNotRefreshAppToken(t *testing.T) {
	var refreshes atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) { refreshes.Add(1) }))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	source, err := NewRefreshingTokenSource(client, TokenPair{AccessToken: "app", ExpiresIn: time.Nanosecond}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer source.Close()
	time.Sleep(2 * time.Millisecond)
	snapshot, err := source.Token(context.Background())
	if err != nil || snapshot.AccessToken() != "app" || refreshes.Load() != 0 {
		t.Fatalf("app token = %q, refreshes = %d, error = %v", snapshot.AccessToken(), refreshes.Load(), err)
	}
	_, err = source.Refresh(context.Background(), snapshot, helix.RefreshReasonUnauthorized)
	if !errors.Is(err, helix.ErrNotRefreshable) {
		t.Fatalf("app Refresh() error = %v", err)
	}
}

var _ helix.RefreshableTokenSource = (*RefreshingTokenSource)(nil)
