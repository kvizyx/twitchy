package oauth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func TestRefreshingTokenSource_closeWinsAndCancelsRefresh(t *testing.T) {
	clock := newTestClock(time.Unix(8_000, 0))
	refreshStarted := make(chan struct{})
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/oauth2/validate" {
			_, _ = io.WriteString(writer, `{"client_id":"client","login":"streamer","scopes":[],"user_id":"42","expires_in":3600}`)
			return
		}
		close(refreshStarted)
		<-release
		_, _ = io.WriteString(writer, `{"access_token":"late","refresh_token":"late-refresh","expires_in":3600,"token_type":"bearer"}`)
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL), WithClock(clock))
	if err != nil {
		t.Fatal(err)
	}
	source, err := NewRefreshingTokenSource(client, TokenPair{AccessToken: "old", RefreshToken: "old-refresh", ExpiresIn: time.Hour}, nil)
	if err != nil {
		t.Fatal(err)
	}
	clock.Advance(time.Hour)
	result := make(chan error, 1)
	go func() {
		_, refreshErr := source.Token(context.Background())
		result <- refreshErr
	}()
	<-refreshStarted
	if err = source.Close(); err != nil {
		t.Fatal(err)
	}
	if refreshErr := <-result; !errors.Is(refreshErr, helix.ErrSessionClosed) {
		t.Fatalf("in-flight refresh error = %v", refreshErr)
	}
	close(release)
	if _, err = source.Token(context.Background()); !errors.Is(err, helix.ErrSessionClosed) {
		t.Fatalf("closed Token() error = %v", err)
	}
	if err = source.RetryCommit(context.Background()); !errors.Is(err, helix.ErrSessionClosed) {
		t.Fatalf("closed RetryCommit() error = %v", err)
	}
}

func TestRefreshingTokenSource_validatesRefreshSkewBounds(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for _, skew := range []time.Duration{-time.Nanosecond, 15*time.Minute + time.Nanosecond} {
		if _, err = NewRefreshingTokenSource(client, TokenPair{AccessToken: "access", ExpiresIn: time.Hour}, nil, WithRefreshSkew(skew)); !errors.Is(err, ErrInvalidOption) {
			t.Fatalf("skew %s error = %v", skew, err)
		}
	}
	if _, err = NewRefreshingTokenSource(client, TokenPair{AccessToken: "access", ExpiresIn: time.Hour}, nil, WithRefreshSkew(0)); err != nil {
		t.Fatalf("zero skew error = %v", err)
	}
}

func TestManagedSession_rejectsShortIntervalWithProductionClock(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}
	source, err := NewRefreshingTokenSource(client, TokenPair{AccessToken: "access", ExpiresIn: time.Hour}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer source.Close()
	if _, err = NewManagedSession(context.Background(), source, WithValidationInterval(time.Minute)); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("short interval error = %v", err)
	}
}
