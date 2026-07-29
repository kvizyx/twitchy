package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

const extensionRateReset = "4102444800"

var extensionBodies = map[string]string{
	"extension_bits_products.json":        `{"data":[{"sku":"sku-1","cost":{"amount":100,"type":"bits"},"in_development":false,"display_name":"Product","expiration":"2026-01-01T00:00:00Z","is_broadcast":true}]}`,
	"extension_configuration.json":        `{"data":[{"segment":"broadcaster","broadcaster_id":"broadcaster-1","content":"hello config","version":"1"},{"segment":"global","content":"global config","version":"2"}]}`,
	"extension_live_channels_page_1.json": `{"data":[{"broadcaster_id":"broadcaster-1","broadcaster_name":"Broadcaster One","game_name":"Game","game_id":"game-1","title":"Title"}],"pagination":"next-cursor"}`,
	"extension_live_channels_page_2.json": `{"data":[{"broadcaster_id":"broadcaster-2","broadcaster_name":"Broadcaster Two","game_name":"Game","game_id":"game-1","title":"Title"}],"pagination":""}`,
	"extension_secrets.json":              `{"data":[{"format_version":1,"secrets":[{"content":"secret-value","active_at":"2025-01-01T00:00:00Z","expires_at":"2026-01-01T00:00:00Z"}]}]}`,
	"extensions.json":                     `{"data":[{"author_name":"Twitch Developers","bits_enabled":true,"can_install":true,"configuration_location":"hosted","description":"Extension description","eula_tos_url":"https://example.test/eula","has_chat_support":true,"icon_url":"https://example.test/icon.png","icon_urls":{"24x24":"https://example.test/icon-small.png"},"id":"extension-1","name":"Extension","privacy_policy_url":"https://example.test/privacy","request_identity_link":true,"screenshot_urls":["https://example.test/screenshot.png"],"state":"Released","subscriptions_support_level":"optional","summary":"Summary","support_email":"support@example.test","version":"1","viewer_summary":"Viewer summary","views":{"mobile":{"viewer_url":"https://example.test/mobile"},"panel":{"viewer_url":"https://example.test/panel","height":300,"can_link_external_content":false},"video_overlay":{"viewer_url":"https://example.test/overlay","can_link_external_content":false},"component":{"viewer_url":"https://example.test/component","aspect_ratio_x":16,"aspect_ratio_y":9,"autoscale":true,"scale_pixels":640,"target_height":100,"can_link_external_content":false},"config":{"viewer_url":"https://example.test/config","can_link_external_content":false}},"allowlisted_config_urls":["https://example.test/config"],"allowlisted_panel_urls":["https://example.test/panel"]}]}`,
	"get_drops_entitlements.json":         `{"data":[{"id":"entitlement-1","benefit_id":"benefit-1","timestamp":"2025-01-01T00:00:00Z","user_id":"user-1","game_id":"game-1","fulfillment_status":"FUTURE_STATUS","last_updated":"2025-01-02T00:00:00Z"}],"pagination":{"cursor":"next-cursor"}}`,
	"update_drops_entitlements.json":      `{"data":[{"status":"SUCCESS","ids":["entitlement-1"]},{"status":"FUTURE_UPDATE","ids":["entitlement-2"]}]}`,
}

func extensionBody(t *testing.T, name string) string {
	t.Helper()
	body, ok := extensionBodies[name]
	if !ok {
		t.Fatalf("unknown extension fixture %q", name)
	}
	return body
}

func extensionClient(t *testing.T, transport *testkit.RecordingRoundTripper, credential helix.Credential) *helix.Client {
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

func extensionResponse(status int, body string) testkit.RoundTripResponse {
	return testkit.RoundTripResponse{
		StatusCode: status,
		Header: http.Header{
			"Ratelimit-Limit":     {"8000"},
			"Ratelimit-Remaining": {"7999"},
			"Ratelimit-Reset":     {extensionRateReset},
		},
		Body: body,
	}
}

func extensionMetaContract(t *testing.T, anchor string, fixture testkit.ContractFixture, transport *testkit.RecordingRoundTripper, meta helix.ResponseMeta, callErr error) {
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

func extensionHeaders(scheme, token string) http.Header {
	return http.Header{
		"Authorization": {scheme + " " + token},
		"Client-Id":     {"client-id"},
	}
}

func extensionFixture(query map[string][]string, body string, headers http.Header, response testkit.ContractResponse) testkit.ContractFixture {
	fixture := testkit.ContractFixture{
		Request:  testkit.ContractRequest{Query: query, Headers: headers, Body: body},
		Response: response,
		Want:     testkit.ContractExpectation{Attempts: 1, RateLimitValid: true},
	}
	fixture.Want.RateLimit.Limit = 8000
	fixture.Want.RateLimit.Remaining = 7999
	return fixture
}
