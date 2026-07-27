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

func TestEventSubManagement_createTransportAuthAndEnvelope(t *testing.T) {
	tests := []struct {
		name      string
		token     helix.TokenClass
		transport helix.EventSubTransport
		body      string
	}{
		{name: "webhook app", token: helix.TokenClassApp, transport: helix.EventSubTransport{Method: "webhook", Callback: "https://callback.test", Secret: "secret-value"}, body: `{"type":"user.update","version":"1","condition":{"user_id":"1234"},"transport":{"method":"webhook","callback":"https://callback.test","secret":"secret-value"}}`},
		{name: "conduit app", token: helix.TokenClassApp, transport: helix.EventSubTransport{Method: "conduit", ConduitID: "conduit-1"}, body: `{"type":"user.update","version":"1","condition":{"user_id":"1234"},"transport":{"method":"conduit","conduit_id":"conduit-1"}}`},
		{name: "websocket user", token: helix.TokenClassUser, transport: helix.EventSubTransport{Method: "websocket", SessionID: "session-1"}, body: `{"type":"user.update","version":"1","condition":{"user_id":"1234"},"transport":{"method":"websocket","session_id":"session-1"}}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{
				StatusCode: http.StatusAccepted,
				Header:     http.Header{"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {verticalSliceRateReset}},
				Body:       `{"data":[{"id":"subscription-1","status":"enabled","type":"user.update","version":"1","condition":{"user_id":"1234"},"created_at":"2020-11-10T14:32:18.730260295Z","transport":{"method":"webhook","callback":"https://callback.test"},"cost":1}],"total":1,"total_cost":1,"max_total_cost":10000}`,
			})
			credential := helix.Credential{AccessToken: "token", ClientID: "client-id", TokenClass: test.token, Scopes: []helix.AuthorizationScope{helix.ScopeChannelReadSubscriptions}}
			client, err := helix.New(helix.WithBaseURL("https://api.twitch.test/helix"), helix.WithHTTPClient(&http.Client{Transport: transport}), helix.WithStaticToken(credential))
			if err != nil {
				t.Fatal(err)
			}
			result, err := client.EventSub.CreateEventSubSubscription(context.Background(), helix.CreateEventSubSubscriptionRequest{Type: "user.update", Version: "1", Condition: helix.EventSubCondition{"user_id": "1234"}, Transport: test.transport})
			if err != nil {
				t.Fatal(err)
			}
			if result.Data.TotalCost != 1 || result.Data.MaxTotalCost != 10000 || result.Data.Subscriptions[0].Version != "1" {
				t.Fatalf("result = %#v", result.Data)
			}
			if got := transport.Requests()[0].Header.Get("Authorization"); got != "Bearer token" {
				t.Fatalf("authorization = %q", got)
			}
			if got := string(transport.Requests()[0].Body); got != test.body {
				t.Fatalf("body = %s, want %s", got, test.body)
			}
		})
	}
}

func TestEventSubManagement_deleteAndGetWire(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(
		testkit.RoundTripResponse{StatusCode: http.StatusNoContent},
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[{"id":"subscription-1","status":"enabled","type":"user.update","version":"1","condition":{"user_id":"1234"},"created_at":"2020-11-10T14:32:18.730260295Z","transport":{"method":"webhook","callback":"https://callback.test"},"cost":1}],"total":1,"total_cost":1,"max_total_cost":10000}`},
	)
	client, err := helix.New(helix.WithBaseURL("https://api.twitch.test/helix"), helix.WithHTTPClient(&http.Client{Transport: transport}), helix.WithStaticToken(helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp}))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.EventSub.DeleteEventSubSubscription(context.Background(), helix.DeleteEventSubSubscriptionRequest{ID: "subscription-1", TransportMethod: "webhook"}); err != nil {
		t.Fatal(err)
	}
	result, err := client.EventSub.GetEventSubSubscriptions(context.Background(), helix.GetEventSubSubscriptionsRequest{Status: "enabled", After: "cursor-1", TransportMethod: "webhook"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Data.Total != 1 || result.Data.Subscriptions[0].ID != "subscription-1" {
		t.Fatalf("subscriptions = %#v", result.Data)
	}
	requests := transport.Requests()
	if requests[0].Path != "/helix/eventsub/subscriptions?id=subscription-1" || requests[1].Path != "/helix/eventsub/subscriptions?after=cursor-1&status=enabled" {
		t.Fatalf("paths = %q, %q", requests[0].Path, requests[1].Path)
	}
}

func TestEventSubManagement_rejectsTransportTokenPairBeforeNetwork(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper()
	client, err := helix.New(helix.WithBaseURL("https://api.twitch.test/helix"), helix.WithHTTPClient(&http.Client{Transport: transport}), helix.WithStaticToken(helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp}))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.EventSub.CreateEventSubSubscription(context.Background(), helix.CreateEventSubSubscriptionRequest{Type: "user.update", Version: "1", Condition: helix.EventSubCondition{"user_id": "1234"}, Transport: helix.EventSubTransport{Method: "websocket", SessionID: "session-1"}})
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("error = %T %v, want AuthError", err, err)
	}
	if len(transport.Requests()) != 0 {
		t.Fatal("invalid token pairing reached network")
	}
}

func TestEventSubManagement_getRejectsFirstBeforeNetwork(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper()
	client, err := helix.New(helix.WithBaseURL("https://api.twitch.test/helix"), helix.WithHTTPClient(&http.Client{Transport: transport}), helix.WithStaticToken(helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp}))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.EventSub.GetEventSubSubscriptions(context.Background(), helix.GetEventSubSubscriptionsRequest{First: 1, TransportMethod: "webhook"})
	if err == nil || !strings.Contains(err.Error(), "first") {
		t.Fatalf("error = %v, want first rejection", err)
	}
	if len(transport.Requests()) != 0 {
		t.Fatal("unsupported first parameter reached network")
	}
}

func TestEventSubManagement_redactsCallbackSecretInErrors(t *testing.T) {
	secret := "secret-value"
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: http.StatusBadRequest, Body: `{"error":"Bad Request","status":400,"message":"callback secret ` + secret + ` rejected"}`})
	client, err := helix.New(helix.WithBaseURL("https://api.twitch.test/helix"), helix.WithHTTPClient(&http.Client{Transport: transport}), helix.WithStaticToken(helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp}))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.EventSub.CreateEventSubSubscription(context.Background(), helix.CreateEventSubSubscriptionRequest{Type: "user.update", Version: "1", Condition: helix.EventSubCondition{"user_id": "1234"}, Transport: helix.EventSubTransport{Method: "webhook", Callback: "https://callback.test", Secret: secret}})
	if err == nil || strings.Contains(err.Error(), secret) || !strings.Contains(err.Error(), "[redacted]") {
		t.Fatalf("error = %v, want redacted callback secret", err)
	}
}
