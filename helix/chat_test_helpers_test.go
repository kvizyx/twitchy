package helix_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

type chatOperationCase struct {
	anchor     string
	fixture    testkit.ContractFixture
	credential helix.Credential
	call       func(*helix.Client) (helix.ResponseMeta, error)
}

func chatSuccessFixture(query url.Values, status int, body string) testkit.ContractFixture {
	return testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query: query,
			Headers: http.Header{
				"Authorization": {"Bearer chat-token"},
				"Client-Id":     {"chat-client"},
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
		Want: testkit.ContractExpectation{Attempts: 1, RateLimitValid: true},
	}
}

func chatRequestBodyFixture(fixture testkit.ContractFixture, body string) testkit.ContractFixture {
	fixture.Request.Body = body
	return fixture
}

func chatRateHeaders() http.Header {
	return http.Header{
		"Ratelimit-Limit":     {"8000"},
		"Ratelimit-Remaining": {"7999"},
		"Ratelimit-Reset":     {verticalSliceRateReset},
	}
}

func chatClient(t *testing.T, testCase chatOperationCase) (*helix.Client, *testkit.RecordingRoundTripper) {
	t.Helper()
	testCase.fixture.Want.RateLimit.Limit = 8000
	testCase.fixture.Want.RateLimit.Remaining = 7999
	return chatClientForCredential(t, testCase.credential, testkit.ResponseFromFixture(testCase.fixture.Response))
}

func chatClientForCredential(t *testing.T, credential helix.Credential, responses ...testkit.RoundTripResponse) (*helix.Client, *testkit.RecordingRoundTripper) {
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

func runChatOperation(t *testing.T, testCase chatOperationCase) {
	t.Helper()
	testCase.fixture.Want.RateLimit.Limit = 8000
	testCase.fixture.Want.RateLimit.Remaining = 7999
	client, transport := chatClient(t, testCase)
	meta, callErr := testCase.call(client)
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, testCase.anchor), testCase.fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return meta, callErr
	}); err != nil {
		t.Fatal(err)
	}
}

func chatMeta[T any](response *helix.Response[T], err error) (helix.ResponseMeta, error) {
	if err != nil {
		return helix.ResponseMeta{}, err
	}
	return response.Meta, nil
}

func chatCredential(tokenClass helix.TokenClass, userID string, scopes ...helix.AuthorizationScope) helix.Credential {
	return helix.Credential{
		AccessToken: "chat-token",
		ClientID:    "chat-client",
		TokenClass:  tokenClass,
		UserID:      userID,
		Scopes:      scopes,
	}
}

func chatBool(value bool) *bool {
	return &value
}

func chatInt(value int) *int {
	return &value
}

func chatString(value string) *string {
	return &value
}
