package helix_test

import (
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func interactiveClient(t *testing.T, transport *testkit.RecordingRoundTripper) *helix.Client {
	return interactiveClientWithUser(t, transport, "123456")
}

func interactiveClientWithUser(t *testing.T, transport *testkit.RecordingRoundTripper, userID string) *helix.Client {
	return interactiveClientWithScopes(t, transport, userID,
		helix.ScopeChannelReadPolls,
		helix.ScopeChannelManagePolls,
		helix.ScopeChannelReadPredictions,
		helix.ScopeChannelManagePredictions,
		helix.ScopeChannelManageRaids,
	)
}

func interactiveClientWithScopes(t *testing.T, transport *testkit.RecordingRoundTripper, userID string, scopes ...helix.AuthorizationScope) *helix.Client {
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

func interactiveResponse(body string) testkit.RoundTripResponse {
	return testkit.RoundTripResponse{StatusCode: http.StatusOK, Header: interactiveRateHeaders(), Body: body}
}

func interactiveRateHeaders() http.Header {
	return http.Header{"Authorization": {"Bearer user-token"}, "Client-Id": {"client-id"}, "Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {"4102444800"}}
}
