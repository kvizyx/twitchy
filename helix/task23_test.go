package helix_test

import (
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func task23Client(t *testing.T, transport *testkit.RecordingRoundTripper) *helix.Client {
	return task23ClientWithUser(t, transport, "123456")
}

func task23ClientWithUser(t *testing.T, transport *testkit.RecordingRoundTripper, userID string) *helix.Client {
	return task23ClientWithScopes(t, transport, userID,
		helix.ScopeChannelReadPolls,
		helix.ScopeChannelManagePolls,
		helix.ScopeChannelReadPredictions,
		helix.ScopeChannelManagePredictions,
		helix.ScopeChannelManageRaids,
	)
}

func task23ClientWithScopes(t *testing.T, transport *testkit.RecordingRoundTripper, userID string, scopes ...helix.AuthorizationScope) *helix.Client {
	t.Helper()
	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: userID, Scopes: scopes}),
	)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func task23Response(body string) testkit.RoundTripResponse {
	return testkit.RoundTripResponse{StatusCode: http.StatusOK, Header: task23RateHeaders(), Body: body}
}

func task23RateHeaders() http.Header {
	return http.Header{"Authorization": {"Bearer user-token"}, "Client-Id": {"client-id"}, "Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {"4102444800"}}
}
