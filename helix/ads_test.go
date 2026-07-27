package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestAdsStartCommercial_sendsJSONBody(t *testing.T) {
	// Given a commercial request and an authorized user token.
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {verticalSliceRateReset}},
		Body:       `{"data":[{"length":60,"message":"","retry_after":480}]}`,
	})
	client := newTask14Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "141981764", Scopes: []helix.AuthorizationScope{helix.ScopeChannelEditCommercial}})

	// When the commercial is started.
	broadcasterID := "141981764"
	length := 60
	result, err := client.Ads.StartCommercial(context.Background(), helix.StartCommercialRequest{BroadcasterID: &broadcasterID, Length: &length})

	// Then the exact body and typed response are returned.
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 || result.Data[0].Length != 60 || result.Data[0].RetryAfter != 480 {
		t.Fatalf("commercial data = %#v", result.Data)
	}
	requests := transport.Requests()
	if len(requests) != 1 || requests[0].Method != http.MethodPost || requests[0].Path != "/helix/channels/commercial" || string(requests[0].Body) != `{"broadcaster_id":"141981764","length":60}` {
		t.Fatalf("commercial request = %#v", requests)
	}
}

func TestAdsScheduleAndSnooze_decodeRFC3339AndCooldown(t *testing.T) {
	// Given schedule and snooze responses followed by an endpoint cooldown.
	transport := testkit.NewRecordingRoundTripper(
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[{"snooze_count":2,"snooze_refresh_at":"2026-07-27T12:00:00Z","next_ad_at":"2026-07-27T12:05:00Z","duration":60,"last_ad_at":"2026-07-27T11:00:00Z","preroll_free_time":90}]}`},
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[{"snooze_count":1,"snooze_refresh_at":"2026-07-27T12:01:00Z","next_ad_at":"2026-07-27T12:10:00Z"}]}`},
		testkit.RoundTripResponse{StatusCode: http.StatusTooManyRequests, Body: `{"error":"Too Many Requests","status":429,"message":"channel has no snoozes left"}`},
	)
	client := newTask14Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp, Scopes: []helix.AuthorizationScope{helix.ScopeChannelReadAds, helix.ScopeChannelManageAds}})

	// When the schedule, snooze, and a second snooze are requested.
	schedule, err := client.Ads.GetAdSchedule(context.Background(), helix.GetAdScheduleRequest{BroadcasterID: "141981764"})
	if err != nil {
		t.Fatal(err)
	}
	snooze, err := client.Ads.SnoozeNextAd(context.Background(), helix.SnoozeNextAdRequest{BroadcasterID: "141981764"})
	if err != nil {
		t.Fatal(err)
	}
	_, cooldownErr := client.Ads.SnoozeNextAd(context.Background(), helix.SnoozeNextAdRequest{BroadcasterID: "141981764"})

	// Then timestamps are parsed, and the endpoint 429 is distinct from bucket exhaustion.
	if schedule.Data[0].SnoozeRefreshAt.Time.Format("2006-01-02T15:04:05Z07:00") != "2026-07-27T12:00:00Z" || schedule.Data[0].SnoozeCount != 2 {
		t.Fatalf("schedule data = %#v", schedule.Data)
	}
	if snooze.Data[0].SnoozeCount != 1 || snooze.Data[0].NextAdAt.IsZero() {
		t.Fatalf("snooze data = %#v", snooze.Data)
	}
	var typedCooldown *helix.CooldownError
	if !errors.As(cooldownErr, &typedCooldown) {
		t.Fatalf("cooldown error = %T %v", cooldownErr, cooldownErr)
	}
	if len(transport.Requests()) != 3 {
		t.Fatalf("request count = %d, want 3", len(transport.Requests()))
	}
}

func TestAdsStartCommercial_cooldownDoesNotWaitForBucket(t *testing.T) {
	// Given a 429 response and a policy that would wait for bucket exhaustion.
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"0"}, "Ratelimit-Reset": {verticalSliceRateReset}},
		Body:       `{"error":"Too Many Requests","status":429,"message":"commercial cooldown"}`,
	})
	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithRateLimitPolicy(helix.RateLimitPolicy{Wait: true}),
		helix.WithStaticToken(helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "141981764", Scopes: []helix.AuthorizationScope{helix.ScopeChannelEditCommercial}}),
	)
	if err != nil {
		t.Fatal(err)
	}

	// When a commercial is started during its endpoint cooldown.
	broadcasterID := "141981764"
	length := 60
	_, callErr := client.Ads.StartCommercial(context.Background(), helix.StartCommercialRequest{BroadcasterID: &broadcasterID, Length: &length})

	// Then cooldown is typed and the non-waitable endpoint makes one request.
	var cooldownErr *helix.CooldownError
	if !errors.As(callErr, &cooldownErr) || len(transport.Requests()) != 1 {
		t.Fatalf("error = %T %v, requests = %d", callErr, callErr, len(transport.Requests()))
	}
}

func TestAds_wrongTokenClassRejectedBeforeNetwork(t *testing.T) {
	// Given an extension token that is not allowed by the Ads manifest rows.
	transport := testkit.NewRecordingRoundTripper()
	client := newTask14Client(t, transport, helix.Credential{AccessToken: "extension-token", ClientID: "client-id", TokenClass: helix.TokenClassExtension})

	// When an Ads operation is invoked.
	_, err := client.Ads.GetAdSchedule(context.Background(), helix.GetAdScheduleRequest{BroadcasterID: "141981764"})

	// Then auth fails before the network.
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) || len(transport.Requests()) != 0 {
		t.Fatalf("error = %T %v, requests = %d", err, err, len(transport.Requests()))
	}
}

func newTask14Client(t *testing.T, transport *testkit.RecordingRoundTripper, credential helix.Credential) *helix.Client {
	t.Helper()
	client, err := helix.New(helix.WithBaseURL("https://api.twitch.test/helix"), helix.WithHTTPClient(&http.Client{Transport: transport}), helix.WithStaticToken(credential))
	if err != nil {
		t.Fatal(err)
	}
	return client
}
