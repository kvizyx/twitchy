package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestEntitlementsGetDropsEntitlements(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(extensionResponse(http.StatusOK, extensionBody(t, "get_drops_entitlements.json")))
	client := extensionClient(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})

	result, callErr := client.Entitlements.GetDropsEntitlements(context.Background(), helix.GetDropsEntitlementsRequest{
		IDs:               []string{"entitlement-1", "entitlement-2"},
		UserID:            "user-1",
		GameID:            "game-1",
		FulfillmentStatus: helix.EntitlementFulfillmentStatus("FUTURE_STATUS"),
		After:             stringPointer("cursor-1"),
		First:             intPointer(100),
	})

	fixture := extensionFixture(
		map[string][]string{
			"after":              {"cursor-1"},
			"first":              {"100"},
			"fulfillment_status": {"FUTURE_STATUS"},
			"game_id":            {"game-1"},
			"id":                 {"entitlement-1", "entitlement-2"},
			"user_id":            {"user-1"},
		},
		"",
		extensionHeaders("Bearer", "app-token"),
		testkit.ContractResponse{Status: http.StatusOK, Headers: http.Header{
			"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {extensionRateReset},
		}, Body: extensionBody(t, "get_drops_entitlements.json"), Success: true},
	)
	extensionMetaContract(t, "get-drops-entitlements", fixture, transport, result.Meta, callErr)
	if len(result.Data) != 1 || result.Data[0].FulfillmentStatus != helix.EntitlementFulfillmentStatus("FUTURE_STATUS") {
		t.Fatalf("entitlements = %#v", result.Data)
	}
}

func TestEntitlementsUpdateDropsEntitlements(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(extensionResponse(http.StatusOK, extensionBody(t, "update_drops_entitlements.json")))
	client := extensionClient(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser})

	result, callErr := client.Entitlements.UpdateDropsEntitlements(context.Background(), helix.UpdateDropsEntitlementsRequest{
		EntitlementIDs:    []string{"entitlement-1", "entitlement-2"},
		FulfillmentStatus: helix.EntitlementFulfillmentStatusFulfilled,
	})
	fixture := extensionFixture(nil, `{"entitlement_ids":["entitlement-1","entitlement-2"],"fulfillment_status":"FULFILLED"}`, extensionHeaders("Bearer", "user-token"), testkit.ContractResponse{
		Status: http.StatusOK, Headers: http.Header{"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {extensionRateReset}}, Body: extensionBody(t, "update_drops_entitlements.json"), Success: true,
	})
	extensionMetaContract(t, "update-drops-entitlements", fixture, transport, result.Meta, callErr)
	if len(result.Data) != 2 || result.Data[1].Status != helix.EntitlementUpdateStatus("FUTURE_UPDATE") {
		t.Fatalf("updates = %#v", result.Data)
	}
}

func TestEntitlementsPagerUsesBearerCursor(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(
		extensionResponse(http.StatusOK, `{"data":[{"id":"entitlement-1","benefit_id":"benefit-1","timestamp":"2025-01-01T00:00:00Z","user_id":"user-1","game_id":"game-1","fulfillment_status":"CLAIMED","last_updated":"2025-01-01T00:00:00Z"}],"pagination":{"cursor":"next-cursor"}}`),
		extensionResponse(http.StatusOK, `{"data":[],"pagination":{}}`),
	)
	client := extensionClient(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	pager, err := client.Entitlements.GetDropsEntitlementsPager(helix.GetDropsEntitlementsRequest{First: intPointer(1)})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) {
		t.Fatalf("first pager page = %v", pager.Err())
	}
	if !pager.Next(context.Background()) || pager.Err() != nil {
		t.Fatalf("pager state = page=%v err=%v", pager.Page(), pager.Err())
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[1].Path != "/helix/entitlements/drops?after=next-cursor&first=1" {
		t.Fatalf("pager requests = %#v", requests)
	}
}
