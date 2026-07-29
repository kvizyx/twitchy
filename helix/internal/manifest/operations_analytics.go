package manifest

import "net/http"

var analyticsOperations = operationsForGroup("Analytics",
	defineOperation("get-extension-analytics", Operation{
		Name:           "Get Extension Analytics",
		Method:         http.MethodGet,
		Path:           "/helix/analytics/extensions",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"analytics:read:extensions"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "extension_id", Type: "String"},
					{Name: "type", Type: "String"},
					{Name: "started_at", Type: "String"},
					{Name: "ended_at", Type: "String"},
					{Name: "first", Type: "Integer"},
					{Name: "after", Type: "String"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 404}},
		Pagination: PaginationSpec{Shape: "cursor", CursorParameter: "after"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:       "Client.Analytics",
			ServiceType:    "AnalyticsService",
			Method:         "GetExtensionAnalytics",
			Signature:      "func (s *AnalyticsService) GetExtensionAnalytics(ctx context.Context, req GetExtensionAnalyticsRequest) (*Response[GetExtensionAnalyticsData], error)",
			PagerSignature: pagerSignature("func (s *AnalyticsService) GetExtensionAnalyticsPager(req GetExtensionAnalyticsRequest, opts ...PagerOption) (*Pager[GetExtensionAnalyticsData], error)"),
			RequestType:    "GetExtensionAnalyticsRequest",
			DataType:       "GetExtensionAnalyticsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-extension-analytics",
	}),
	defineOperation("get-game-analytics", Operation{
		Name:           "Get Game Analytics",
		Method:         http.MethodGet,
		Path:           "/helix/analytics/games",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"analytics:read:games"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "game_id", Type: "String"},
					{Name: "type", Type: "String"},
					{Name: "started_at", Type: "String"},
					{Name: "ended_at", Type: "String"},
					{Name: "first", Type: "Integer"},
					{Name: "after", Type: "String"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 404}},
		Pagination: PaginationSpec{Shape: "cursor", CursorParameter: "after"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:       "Client.Analytics",
			ServiceType:    "AnalyticsService",
			Method:         "GetGameAnalytics",
			Signature:      "func (s *AnalyticsService) GetGameAnalytics(ctx context.Context, req GetGameAnalyticsRequest) (*Response[GetGameAnalyticsData], error)",
			PagerSignature: pagerSignature("func (s *AnalyticsService) GetGameAnalyticsPager(req GetGameAnalyticsRequest, opts ...PagerOption) (*Pager[GetGameAnalyticsData], error)"),
			RequestType:    "GetGameAnalyticsRequest",
			DataType:       "GetGameAnalyticsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-game-analytics",
	}),
)
