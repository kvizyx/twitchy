package manifest

import "net/http"

var adsOperations = operationsForGroup("Ads",
	defineOperation("start-commercial", Operation{
		Name:           "Start Commercial",
		Method:         http.MethodPost,
		Path:           "/helix/channels/commercial",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"channel:edit:commercial"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"body": {
					{Name: "broadcaster_id", Type: "String"},
					{Name: "length", Type: "Integer"},
				},
			},
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 404, 429}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Implementation: ImplementationSpec{
			Selector:    "Client.Ads",
			ServiceType: "AdsService",
			Method:      "StartCommercial",
			Signature:   "func (s *AdsService) StartCommercial(ctx context.Context, req StartCommercialRequest) (*Response[StartCommercialData], error)",
			RequestType: "StartCommercialRequest",
			DataType:    "StartCommercialData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#start-commercial",
	}),
	defineOperation("get-ad-schedule", Operation{
		Name:           "Get Ad Schedule",
		Method:         http.MethodGet,
		Path:           "/helix/channels/ads",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"channel:read:ads"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "broadcaster_id", Type: "String", Required: true},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 500}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Ads",
			ServiceType: "AdsService",
			Method:      "GetAdSchedule",
			Signature:   "func (s *AdsService) GetAdSchedule(ctx context.Context, req GetAdScheduleRequest) (*Response[GetAdScheduleData], error)",
			RequestType: "GetAdScheduleRequest",
			DataType:    "GetAdScheduleData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-ad-schedule",
	}),
	defineOperation("snooze-next-ad", Operation{
		Name:           "Snooze Next Ad",
		Method:         http.MethodPost,
		Path:           "/helix/channels/ads/schedule/snooze",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"channel:manage:ads"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "broadcaster_id", Type: "String", Required: true},
				},
			},
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 429, 500}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Implementation: ImplementationSpec{
			Selector:    "Client.Ads",
			ServiceType: "AdsService",
			Method:      "SnoozeNextAd",
			Signature:   "func (s *AdsService) SnoozeNextAd(ctx context.Context, req SnoozeNextAdRequest) (*Response[SnoozeNextAdData], error)",
			RequestType: "SnoozeNextAdRequest",
			DataType:    "SnoozeNextAdData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#snooze-next-ad",
	}),
)
