package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

const (
	guestStarRateReset = "4102444800"
)

func guestStarCredential(scopes ...helix.AuthorizationScope) helix.Credential {
	return helix.Credential{AccessToken: "guest-star-token", ClientID: "guest-star-client", TokenClass: helix.TokenClassUser, UserID: "5678", Scopes: scopes}
}

func guestStarClient(t *testing.T, transport *testkit.RecordingRoundTripper, credential helix.Credential) (*helix.Client, error) {
	t.Helper()
	return helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(credential),
	)
}

func guestStarBody(t *testing.T, name string) string {
	t.Helper()
	body, err := testkit.LoadText("testdata/task21", name)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func guestStarFixture(query map[string][]string, body string, status int, responseBody string) testkit.ContractFixture {
	return testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query: query,
			Body:  body,
			Headers: http.Header{
				"Authorization": {"Bearer guest-star-token"},
				"Client-Id":     {"guest-star-client"},
			},
		},
		Response: testkit.ContractResponse{
			Status: status,
			Headers: http.Header{
				"Ratelimit-Limit":     {"8000"},
				"Ratelimit-Remaining": {"7999"},
				"Ratelimit-Reset":     {guestStarRateReset},
			},
			Body:    responseBody,
			Success: status >= http.StatusOK && status < http.StatusMultipleChoices,
		},
		Want: testkit.ContractExpectation{Attempts: 1, RateLimitValid: true},
	}
}

func runGuestStarContract(t *testing.T, fixture testkit.ContractFixture, anchor string, transport *testkit.RecordingRoundTripper, meta helix.ResponseMeta, callErr error) error {
	t.Helper()
	fixture.Want.RateLimit.Limit = 8000
	fixture.Want.RateLimit.Remaining = 7999
	return testkit.RunManifestContract(context.Background(), manifestOperation(t, anchor), fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return meta, callErr
	})
}

func guestStarMeta[T any](response *helix.Response[T], err error) (helix.ResponseMeta, error) {
	if err != nil {
		return helix.ResponseMeta{}, err
	}
	return response.Meta, nil
}

func guestStarGroupLayoutPointer(value helix.GuestStarGroupLayout) *helix.GuestStarGroupLayout {
	return &value
}

func guestStarString(value string) *string {
	return &value
}

func requireGuestStarAuthError(t *testing.T, callErr error) {
	t.Helper()
	var authErr *helix.AuthError
	if !errors.As(callErr, &authErr) {
		t.Fatalf("error = %T %v, want AuthError", callErr, callErr)
	}
}
