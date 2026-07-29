package manifest

import "net/http"

var hypeTrainOperations = operationsForGroup("Hype Train",
	defineOperation("get-hype-train-status", Operation{
		Name:           "Get Hype Train Status",
		Method:         http.MethodGet,
		Path:           "/helix/hypetrain/status",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:read:hype_train"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "broadcaster_id", Type: "String", Required: true},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 500}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.HypeTrain",
			ServiceType: "HypeTrainService",
			Method:      "GetHypeTrainStatus",
			Signature:   "func (s *HypeTrainService) GetHypeTrainStatus(ctx context.Context, req GetHypeTrainStatusRequest) (*Response[GetHypeTrainStatusData], error)",
			RequestType: "GetHypeTrainStatusRequest",
			DataType:    "GetHypeTrainStatusData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-hype-train-status",
	}),
)
