package manifest

import "net/http"

var entitlementsOperations = operationsForGroup("Entitlements",
	defineOperation("get-drops-entitlements", Operation{
		Name:           "Get Drops Entitlements",
		Method:         http.MethodGet,
		Path:           "/helix/entitlements/drops",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "id", Type: "String"},
					{Name: "user_id", Type: "String"},
					{Name: "game_id", Type: "String"},
					{Name: "fulfillment_status", Type: "String"},
					{Name: "after", Type: "String"},
					{Name: "first", Type: "Integer"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 403, 500}},
		Pagination: PaginationSpec{Shape: "cursor", CursorParameter: "after"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:       "Client.Entitlements",
			ServiceType:    "EntitlementsService",
			Method:         "GetDropsEntitlements",
			Signature:      "func (s *EntitlementsService) GetDropsEntitlements(ctx context.Context, req GetDropsEntitlementsRequest) (*Response[GetDropsEntitlementsData], error)",
			PagerSignature: pagerSignature("func (s *EntitlementsService) GetDropsEntitlementsPager(req GetDropsEntitlementsRequest, opts ...PagerOption) (*Pager[GetDropsEntitlementsData], error)"),
			RequestType:    "GetDropsEntitlementsRequest",
			DataType:       "GetDropsEntitlementsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-drops-entitlements",
	}),
	defineOperation("update-drops-entitlements", Operation{
		Name:           "Update Drops Entitlements",
		Method:         http.MethodPatch,
		Path:           "/helix/entitlements/drops",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"body": {
					{Name: "entitlement_ids", Type: "String[]"},
					{Name: "fulfillment_status", Type: "String"},
				},
			},
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 500}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Entitlements",
			ServiceType: "EntitlementsService",
			Method:      "UpdateDropsEntitlements",
			Signature:   "func (s *EntitlementsService) UpdateDropsEntitlements(ctx context.Context, req UpdateDropsEntitlementsRequest) (*Response[UpdateDropsEntitlementsData], error)",
			RequestType: "UpdateDropsEntitlementsRequest",
			DataType:    "UpdateDropsEntitlementsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#update-drops-entitlements",
	}),
)
