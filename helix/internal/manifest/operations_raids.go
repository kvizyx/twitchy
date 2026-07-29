package manifest

import "net/http"

var raidsOperations = operationsForGroup("Raids",
	defineOperation("start-a-raid", Operation{
		Name:           "Start a raid",
		Method:         http.MethodPost,
		Path:           "/helix/raids",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:manage:raids"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "from_broadcaster_id", Type: "String", Required: true},
					{Name: "to_broadcaster_id", Type: "String", Required: true},
				},
			},
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 404, 409, 429}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Implementation: ImplementationSpec{
			Selector:    "Client.Raids",
			ServiceType: "RaidsService",
			Method:      "StartRaid",
			Signature:   "func (s *RaidsService) StartRaid(ctx context.Context, req StartRaidRequest) (*Response[StartRaidData], error)",
			RequestType: "StartRaidRequest",
			DataType:    "StartRaidData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#start-a-raid",
	}),
	defineOperation("cancel-a-raid", Operation{
		Name:           "Cancel a raid",
		Method:         http.MethodDelete,
		Path:           "/helix/raids",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:manage:raids"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "broadcaster_id", Type: "String", Required: true},
				},
			},
		},
		Response:   ResponseSpec{Format: "unknown", Status: []int{204, 400, 401, 404, 429}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Implementation: ImplementationSpec{
			Selector:    "Client.Raids",
			ServiceType: "RaidsService",
			Method:      "CancelRaid",
			Signature:   "func (s *RaidsService) CancelRaid(ctx context.Context, req CancelRaidRequest) (*Response[CancelRaidData], error)",
			RequestType: "CancelRaidRequest",
			DataType:    "CancelRaidData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#cancel-a-raid",
	}),
)
