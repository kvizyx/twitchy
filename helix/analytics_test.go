package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestAnalytics_returnsCSVURLsWithoutFetchingThem(t *testing.T) {
	// Given an analytics response whose URL points at a forbidden external host.
	transport := testkit.NewFailingRoundTripper(testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{
		StatusCode: http.StatusOK,
		Body:       `{"data":[{"extension_id":"ext-1","URL":"https://csv.invalid/report.csv","type":"overview_v2","date_range":{"started_at":"2026-07-01T00:00:00Z","ended_at":"2026-07-02T00:00:00Z"}}],"pagination":{"cursor":"next"}}`,
	}), "api.twitch.test")
	client, err := helix.New(helix.WithBaseURL("https://api.twitch.test/helix"), helix.WithHTTPClient(&http.Client{Transport: transport}), helix.WithStaticToken(helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, Scopes: []helix.AuthorizationScope{helix.ScopeAnalyticsReadExtensions}}))
	if err != nil {
		t.Fatal(err)
	}

	// When analytics is requested.
	result, err := client.Analytics.GetExtensionAnalytics(context.Background(), helix.GetExtensionAnalyticsRequest{})

	// Then only the Helix response is decoded and the URL is returned untouched.
	if err != nil {
		t.Fatal(err)
	}
	if result.Data[0].URL != "https://csv.invalid/report.csv" || result.Pagination.Cursor() != "next" {
		t.Fatalf("analytics result = %#v", result)
	}
}

func TestAnalyticsPager_forwardsAfterCursor(t *testing.T) {
	// Given two pages from the extension analytics endpoint.
	transport := testkit.NewRecordingRoundTripper(
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[{"extension_id":"ext-1","URL":"https://csv.invalid/one.csv","type":"overview_v2","date_range":{"started_at":"2026-07-01T00:00:00Z","ended_at":"2026-07-02T00:00:00Z"}}],"pagination":{"cursor":"next"}}`},
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[{"extension_id":"ext-2","URL":"https://csv.invalid/two.csv","type":"overview_v2","date_range":{"started_at":"2026-07-03T00:00:00Z","ended_at":"2026-07-04T00:00:00Z"}}]}`},
	)
	client := newTask14Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, Scopes: []helix.AuthorizationScope{helix.ScopeAnalyticsReadExtensions}})

	// When both pages are consumed.
	pager, err := client.Analytics.GetExtensionAnalyticsPager(helix.GetExtensionAnalyticsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || pager.Page().Data[0].ExtensionID != "ext-1" || !pager.Next(context.Background()) || pager.Page().Data[0].ExtensionID != "ext-2" {
		t.Fatalf("pager did not return both pages")
	}

	// Then the second request forwards the first response cursor.
	requests := transport.Requests()
	if len(requests) != 2 || requests[1].Path != "/helix/analytics/extensions?after=next" {
		t.Fatalf("analytics pager requests = %#v", requests)
	}
}
