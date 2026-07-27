package helix_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestWhispersSendWhisper(t *testing.T) {
	// Given a scoped sender and a message that must not escape error text.
	transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusNoContent, ""))
	client := task26Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "sender", Scopes: []helix.AuthorizationScope{helix.ScopeUserManageWhispers}})
	message := "secret whisper message"

	// When a whisper is sent.
	result, err := client.Whispers.SendWhisper(context.Background(), helix.SendWhisperRequest{FromUserID: "sender", ToUserID: "recipient", Message: message})

	// Then query/body encoding is exact and the empty success is accepted.
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || len(transport.Requests()) != 1 || transport.Requests()[0].Path != "/helix/whispers?from_user_id=sender&to_user_id=recipient" || string(transport.Requests()[0].Body) != `{"message":"secret whisper message"}` {
		t.Fatalf("whisper request = %#v", transport.Requests())
	}
}

func TestWhispersSendWhisperRedactsMessageFromTypedErrors(t *testing.T) {
	// Given an upstream error body that repeats the secret message.
	message := "secret whisper message"
	transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusBadRequest, `{"error":"Bad Request","status":400,"message":"secret whisper message"}`))
	client := task26Client(t, transport, helix.Credential{AccessToken: "secret-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "sender", Scopes: []helix.AuthorizationScope{helix.ScopeUserManageWhispers}})

	// When the whisper fails.
	_, err := client.Whispers.SendWhisper(context.Background(), helix.SendWhisperRequest{FromUserID: "sender", ToUserID: "recipient", Message: message})

	// Then the typed error is returned without exposing token or message content.
	var apiErr *helix.APIError
	if !errors.As(err, &apiErr) || strings.Contains(err.Error(), message) || strings.Contains(err.Error(), "secret-token") {
		t.Fatalf("redaction error = %T %v", err, err)
	}
}

func TestWhispersSendWhisperDoesNotReplay(t *testing.T) {
	// Given a whisper mutation and a transient upstream failure.
	transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusServiceUnavailable, `{"error":"Unavailable","status":503,"message":"try again"}`))
	client := task26Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "sender", Scopes: []helix.AuthorizationScope{helix.ScopeUserManageWhispers}})

	// When the whisper is sent.
	_, err := client.Whispers.SendWhisper(context.Background(), helix.SendWhisperRequest{FromUserID: "sender", ToUserID: "recipient", Message: "hello"})

	// Then the non-replayable mutation makes exactly one request.
	var apiErr *helix.APIError
	if !errors.As(err, &apiErr) || len(transport.Requests()) != 1 {
		t.Fatalf("error = %T %v, requests = %#v", err, err, transport.Requests())
	}
}
