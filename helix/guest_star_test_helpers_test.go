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

var guestStarBodies = map[string]string{
	"invites.json":  `{"data":[{"user_id":"9012","invited_at":"2026-07-27T12:00:00Z","status":"READY","is_video_enabled":true,"is_audio_enabled":false,"is_video_available":true,"is_audio_available":true}]}`,
	"session.json":  `{"data":[{"id":"session-1","guests":[{"slot_id":"1","user_id":"9012","user_display_name":"Guest","user_login":"guest","is_live":true,"volume":75,"assigned_at":"2026-07-27T12:00:00Z","audio_settings":{"is_available":true,"is_host_enabled":true,"is_guest_enabled":false},"video_settings":{"is_available":true,"is_host_enabled":true,"is_guest_enabled":true}}]}]}`,
	"settings.json": `{"data":[{"is_moderator_send_live_enabled":true,"slot_count":4,"is_browser_source_audio_enabled":false,"group_layout":"SCREENSHARE_LAYOUT","browser_source_token":"browser-token"}]}`,
}

func guestStarBody(t *testing.T, name string) string {
	t.Helper()
	body, ok := guestStarBodies[name]
	if !ok {
		t.Fatalf("unknown guest star fixture %q", name)
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
