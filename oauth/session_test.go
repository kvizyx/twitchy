package oauth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func TestManagedSession_validatesImmediatelyAndOnExactIntervals(t *testing.T) {
	clock := newTestClock(time.Unix(4_000, 0))
	var validations atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		validations.Add(1)
		_, _ = io.WriteString(writer, `{"client_id":"client","login":"streamer","scopes":[],"user_id":"42","expires_in":7200}`)
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL), WithClock(clock))
	if err != nil {
		t.Fatal(err)
	}
	source, err := NewRefreshingTokenSource(client, TokenPair{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: time.Hour}, nil)
	if err != nil {
		t.Fatal(err)
	}
	session, err := NewManagedSession(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}
	if got := validations.Load(); got != 1 {
		t.Fatalf("initial validations = %d, want 1", got)
	}
	clock.Advance(time.Hour - time.Nanosecond)
	if got := validations.Load(); got != 1 {
		t.Fatalf("early validations = %d, want 1", got)
	}
	clock.Advance(time.Nanosecond)
	select {
	case <-waitForAtomic(&validations, 2):
	case <-time.After(time.Second):
		t.Fatal("hourly validation did not run")
	}
	if _, err = source.Token(context.Background()); err != nil {
		t.Fatalf("active Token() error = %v", err)
	}
	if err = session.Close(); err != nil {
		t.Fatal(err)
	}
	if err = session.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestManagedSession_periodicInvalidationStopsServingAndScheduling(t *testing.T) {
	clock := newTestClock(time.Unix(5_000, 0))
	var validations atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if validations.Add(1) == 1 {
			_, _ = io.WriteString(writer, `{"client_id":"client","login":"streamer","scopes":[],"user_id":"42","expires_in":7200}`)
			return
		}
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(writer, `{"error":"invalid_token","error_description":"expired"}`)
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL), WithClock(clock))
	if err != nil {
		t.Fatal(err)
	}
	source, err := NewRefreshingTokenSource(client, TokenPair{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: time.Hour}, nil)
	if err != nil {
		t.Fatal(err)
	}
	session, err := NewManagedSession(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	clock.Advance(time.Hour)
	select {
	case <-waitForAtomic(&validations, 2):
	case <-time.After(time.Second):
		t.Fatal("periodic validation did not run")
	}
	if _, err = source.Token(context.Background()); !errors.Is(err, helix.ErrInvalidSession) {
		t.Fatalf("invalidated Token() error = %v", err)
	}
	clock.Advance(time.Hour)
	time.Sleep(10 * time.Millisecond)
	if got := validations.Load(); got != 2 {
		t.Fatalf("validations after terminal state = %d, want 2", got)
	}
}

func TestManagedSession_transientValidationFailureHooksAndRetriesNextTick(t *testing.T) {
	clock := newTestClock(time.Unix(6_000, 0))
	var validations atomic.Int32
	var hookCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		count := validations.Add(1)
		if count == 2 {
			writer.WriteHeader(http.StatusTooManyRequests)
			_, _ = io.WriteString(writer, `{"error":"slow_down"}`)
			return
		}
		_, _ = io.WriteString(writer, `{"client_id":"client","login":"streamer","scopes":[],"user_id":"42","expires_in":7200}`)
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL), WithClock(clock))
	if err != nil {
		t.Fatal(err)
	}
	source, err := NewRefreshingTokenSource(client, TokenPair{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: time.Hour}, nil)
	if err != nil {
		t.Fatal(err)
	}
	session, err := NewManagedSession(context.Background(), source, WithValidationErrorHook(func(error) { hookCalls.Add(1) }))
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	clock.Advance(time.Hour)
	select {
	case <-waitForAtomic(&validations, 2):
	case <-time.After(time.Second):
		t.Fatal("transient validation did not run")
	}
	if got := hookCalls.Load(); got != 1 {
		t.Fatalf("validation hook calls = %d, want 1", got)
	}
	if _, err = source.Token(context.Background()); err != nil {
		t.Fatalf("transient failure Token() error = %v", err)
	}
	clock.Advance(time.Hour)
	select {
	case <-waitForAtomic(&validations, 3):
	case <-time.After(time.Second):
		t.Fatal("validation retry did not run")
	}
}

func TestManagedSession_initialInvalidValidationReturnsNoSession(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   error
	}{
		{name: "unauthorized", status: http.StatusUnauthorized, body: `{"error":"invalid_token"}`, want: helix.ErrInvalidSession},
		{name: "other client error", status: http.StatusBadRequest, body: `{"error":"invalid_request"}`},
		{name: "malformed success", status: http.StatusOK, body: `{"client_id":"client"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clock := newTestClock(time.Unix(7_000, 0))
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(test.status)
				_, _ = io.WriteString(writer, test.body)
			}))
			defer server.Close()
			client, err := New(WithBaseURL(server.URL), WithClock(clock))
			if err != nil {
				t.Fatal(err)
			}
			source, err := NewRefreshingTokenSource(client, TokenPair{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: time.Hour}, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer source.Close()
			_, err = NewManagedSession(context.Background(), source)
			if err == nil || (test.want != nil && !errors.Is(err, test.want)) {
				t.Fatalf("NewManagedSession() error = %v", err)
			}
		})
	}
}

func waitForAtomic(value *atomic.Int32, want int32) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		for value.Load() < want {
			time.Sleep(time.Millisecond)
		}
		close(done)
	}()
	return done
}
