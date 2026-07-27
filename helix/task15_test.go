package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

const task15RateReset = "4102444800"

func task15Client(t *testing.T, transport *testkit.RecordingRoundTripper) *helix.Client {
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

func task15Response(body string) testkit.RoundTripResponse {
	return testkit.RoundTripResponse{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Ratelimit-Limit":     {"8000"},
			"Ratelimit-Remaining": {"7999"},
			"Ratelimit-Reset":     {task15RateReset},
		},
		Body: body,
	}
}

func task15Contract(t *testing.T, anchor string, fixture testkit.ContractFixture, transport *testkit.RecordingRoundTripper, meta helix.ResponseMeta, callErr error) {
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

func task15Headers() http.Header {
	return http.Header{
		"Authorization": {"Bearer user-token"},
		"Client-Id":     {"client-id"},
	}
}

func task15RateHeaders() http.Header {
	headers := task15Headers()
	headers.Set("Ratelimit-Limit", "8000")
	headers.Set("Ratelimit-Remaining", "7999")
	headers.Set("Ratelimit-Reset", task15RateReset)
	return headers
}

func task15Fixture(query map[string][]string, body string, response testkit.ContractResponse) testkit.ContractFixture {
	fixture := testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query:   query,
			Headers: task15Headers(),
			Body:    body,
		},
		Response: response,
		Want:     testkit.ContractExpectation{Attempts: 1, RateLimitValid: true},
	}
	fixture.Want.RateLimit.Limit = 8000
	fixture.Want.RateLimit.Remaining = 7999
	return fixture
}

func task15Success(body string) testkit.ContractResponse {
	return testkit.ContractResponse{Status: http.StatusOK, Headers: task15RateHeaders(), Body: body, Success: true}
}

func task15NoContent() testkit.ContractResponse {
	return testkit.ContractResponse{Status: http.StatusNoContent, Headers: task15RateHeaders(), Success: true}
}

func stringPointer(value string) *string { return &value }
func boolPointer(value bool) *bool       { return &value }
func intPointer(value int) *int          { return &value }
func int64Pointer(value int64) *int64    { return &value }
