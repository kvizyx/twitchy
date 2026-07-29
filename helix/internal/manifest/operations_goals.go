package manifest

import "net/http"

var goalsOperations = operationsForGroup("Goals",
	defineOperation("get-creator-goals", Operation{
		Name:           "Get Creator Goals",
		Method:         http.MethodGet,
		Path:           "/helix/goals",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:read:goals"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "broadcaster_id", Type: "String", Required: true},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Goals",
			ServiceType: "GoalsService",
			Method:      "GetCreatorGoals",
			Signature:   "func (s *GoalsService) GetCreatorGoals(ctx context.Context, req GetCreatorGoalsRequest) (*Response[GetCreatorGoalsData], error)",
			RequestType: "GetCreatorGoalsRequest",
			DataType:    "GetCreatorGoalsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-creator-goals",
	}),
)
