package helix_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

type moderationOperationCase struct {
	anchor     string
	fixture    testkit.ContractFixture
	credential helix.Credential
	call       func(*helix.Client) (helix.ResponseMeta, error)
}

func moderationSuccessFixture(query url.Values, status int, body string) testkit.ContractFixture {
	return testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query: query,
			Headers: http.Header{
				"Authorization": {"Bearer moderation-token"},
				"Client-Id":     {"moderation-client"},
			},
		},
		Response: testkit.ContractResponse{
			Status: status,
			Headers: http.Header{
				"Ratelimit-Limit":     {"8000"},
				"Ratelimit-Remaining": {"7999"},
				"Ratelimit-Reset":     {verticalSliceRateReset},
			},
			Body:    body,
			Success: true,
		},
		Want: testkit.ContractExpectation{
			Attempts:       1,
			RateLimitValid: true,
			RateLimit: struct {
				Limit     int `json:"limit"`
				Remaining int `json:"remaining"`
			}{Limit: 8000, Remaining: 7999},
		},
	}
}

func moderationBodyFixture(fixture testkit.ContractFixture, body string) testkit.ContractFixture {
	fixture.Request.Body = body
	return fixture
}

func moderationClient(t *testing.T, credential helix.Credential, responses ...testkit.RoundTripResponse) (*helix.Client, *testkit.RecordingRoundTripper) {
	t.Helper()
	transport := testkit.NewRecordingRoundTripper(responses...)
	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(credential),
	)
	if err != nil {
		t.Fatal(err)
	}
	return client, transport
}

func moderationCredential(tokenClass helix.TokenClass, userID string, scopes ...helix.AuthorizationScope) helix.Credential {
	return helix.Credential{
		AccessToken: "moderation-token",
		ClientID:    "moderation-client",
		TokenClass:  tokenClass,
		UserID:      userID,
		Scopes:      scopes,
	}
}

func runModerationOperation(t *testing.T, testCase moderationOperationCase) {
	t.Helper()
	client, transport := moderationClient(t, testCase.credential, testkit.ResponseFromFixture(testCase.fixture.Response))
	meta, callErr := testCase.call(client)
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, testCase.anchor), testCase.fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return meta, callErr
	}); err != nil {
		t.Fatal(err)
	}
}

func moderationMeta[T any](response *helix.Response[T], err error) (helix.ResponseMeta, error) {
	if err != nil {
		return helix.ResponseMeta{}, err
	}
	return response.Meta, nil
}

func moderationString(value string) *string { return &value }
func moderationInt(value int) *int          { return &value }
