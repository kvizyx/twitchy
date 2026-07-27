package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

const channelRateReset = "4102444800"

func channelClient(t *testing.T, transport *testkit.RecordingRoundTripper) *helix.Client {
	t.Helper()
	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(helix.Credential{
			AccessToken: "user-token",
			ClientID:    "client-id",
			TokenClass:  helix.TokenClassUser,
			UserID:      "123456",
			Scopes: []helix.AuthorizationScope{
				helix.ScopeChannelManageBroadcast,
				helix.ScopeChannelReadEditors,
				helix.ScopeUserReadFollows,
				helix.ScopeModeratorReadFollowers,
				helix.ScopeChannelManageRedemptions,
				helix.ScopeChannelReadRedemptions,
				helix.ScopeChannelReadCharity,
			},
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func channelResponse(body string) testkit.RoundTripResponse {
	return testkit.RoundTripResponse{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Ratelimit-Limit":     {"8000"},
			"Ratelimit-Remaining": {"7999"},
			"Ratelimit-Reset":     {channelRateReset},
		},
		Body: body,
	}
}

func channelContract(t *testing.T, anchor string, fixture testkit.ContractFixture, transport *testkit.RecordingRoundTripper, meta helix.ResponseMeta, callErr error) {
	t.Helper()
	if callErr != nil {
		t.Fatal(callErr)
	}
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, anchor), fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return meta, nil
	}); err != nil {
		t.Fatal(err)
	}
}

func channelHeaders() http.Header {
	return http.Header{
		"Authorization": {"Bearer user-token"},
		"Client-Id":     {"client-id"},
	}
}

func channelRateHeaders() http.Header {
	headers := channelHeaders()
	headers.Set("Ratelimit-Limit", "8000")
	headers.Set("Ratelimit-Remaining", "7999")
	headers.Set("Ratelimit-Reset", channelRateReset)
	return headers
}

func channelFixture(query map[string][]string, body string, response testkit.ContractResponse) testkit.ContractFixture {
	fixture := testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query:   query,
			Headers: channelHeaders(),
			Body:    body,
		},
		Response: response,
		Want:     testkit.ContractExpectation{Attempts: 1, RateLimitValid: true},
	}
	fixture.Want.RateLimit.Limit = 8000
	fixture.Want.RateLimit.Remaining = 7999
	return fixture
}

func channelSuccess(body string) testkit.ContractResponse {
	return testkit.ContractResponse{Status: http.StatusOK, Headers: channelRateHeaders(), Body: body, Success: true}
}

func channelNoContent() testkit.ContractResponse {
	return testkit.ContractResponse{Status: http.StatusNoContent, Headers: channelRateHeaders(), Success: true}
}

func stringPointer(value string) *string { return &value }
func boolPointer(value bool) *bool       { return &value }
func intPointer(value int) *int          { return &value }
func int64Pointer(value int64) *int64    { return &value }
