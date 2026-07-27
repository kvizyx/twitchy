package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

const task19RateReset = "4102444800"

func task19Body(t *testing.T, name string) string {
	t.Helper()
	body, err := testkit.LoadText("testdata/task19", name)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func task19Client(t *testing.T, transport *testkit.RecordingRoundTripper, credential helix.Credential) *helix.Client {
	t.Helper()
	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(credential),
	)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func task19Response(status int, body string) testkit.RoundTripResponse {
	return testkit.RoundTripResponse{
		StatusCode: status,
		Header: http.Header{
			"Ratelimit-Limit":     {"8000"},
			"Ratelimit-Remaining": {"7999"},
			"Ratelimit-Reset":     {task19RateReset},
		},
		Body: body,
	}
}

func task19MetaContract(t *testing.T, anchor string, fixture testkit.ContractFixture, transport *testkit.RecordingRoundTripper, meta helix.ResponseMeta, callErr error) {
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

func task19Headers(scheme, token string) http.Header {
	return http.Header{
		"Authorization": {scheme + " " + token},
		"Client-Id":     {"client-id"},
	}
}

func task19Fixture(query map[string][]string, body string, headers http.Header, response testkit.ContractResponse) testkit.ContractFixture {
	fixture := testkit.ContractFixture{
		Request:  testkit.ContractRequest{Query: query, Headers: headers, Body: body},
		Response: response,
		Want:     testkit.ContractExpectation{Attempts: 1, RateLimitValid: true},
	}
	fixture.Want.RateLimit.Limit = 8000
	fixture.Want.RateLimit.Remaining = 7999
	return fixture
}
