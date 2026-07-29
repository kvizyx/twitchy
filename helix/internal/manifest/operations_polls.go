package manifest

import "net/http"

var pollsOperations = operationsForGroup("Polls",
	defineOperation("get-polls", Operation{
		Name:           "Get Polls",
		Method:         http.MethodGet,
		Path:           "/helix/polls",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:read:polls", "channel:manage:polls"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "broadcaster_id", Type: "String", Required: true},
					{Name: "id", Type: "String"},
					{Name: "first", Type: "String"},
					{Name: "after", Type: "String"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 404}},
		Pagination: PaginationSpec{Shape: "cursor", CursorParameter: "after"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:       "Client.Polls",
			ServiceType:    "PollsService",
			Method:         "GetPolls",
			Signature:      "func (s *PollsService) GetPolls(ctx context.Context, req GetPollsRequest) (*Response[GetPollsData], error)",
			PagerSignature: pagerSignature("func (s *PollsService) GetPollsPager(req GetPollsRequest, opts ...PagerOption) (*Pager[GetPollsData], error)"),
			RequestType:    "GetPollsRequest",
			DataType:       "GetPollsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-polls",
	}),
	defineOperation("create-poll", Operation{
		Name:           "Create Poll",
		Method:         http.MethodPost,
		Path:           "/helix/polls",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:manage:polls"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"body": {
					{Name: "broadcaster_id", Type: "String", Required: true},
					{Name: "title", Type: "String", Required: true},
					{Name: "choices", Type: "Object[]", Required: true},
					{Name: "title", Type: "String", Required: true},
					{Name: "duration", Type: "Integer", Required: true},
					{Name: "channel_points_voting_enabled", Type: "Boolean"},
					{Name: "channel_points_per_vote", Type: "Integer"},
				},
			},
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Polls",
			ServiceType: "PollsService",
			Method:      "CreatePoll",
			Signature:   "func (s *PollsService) CreatePoll(ctx context.Context, req CreatePollRequest) (*Response[CreatePollData], error)",
			RequestType: "CreatePollRequest",
			DataType:    "CreatePollData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#create-poll",
	}),
	defineOperation("end-poll", Operation{
		Name:           "End Poll",
		Method:         http.MethodPatch,
		Path:           "/helix/polls",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:manage:polls"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"body": {
					{Name: "broadcaster_id", Type: "String", Required: true},
					{Name: "id", Type: "String", Required: true},
					{Name: "status", Type: "String", Required: true},
				},
			},
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Polls",
			ServiceType: "PollsService",
			Method:      "EndPoll",
			Signature:   "func (s *PollsService) EndPoll(ctx context.Context, req EndPollRequest) (*Response[EndPollData], error)",
			RequestType: "EndPollRequest",
			DataType:    "EndPollData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#end-poll",
	}),
)
