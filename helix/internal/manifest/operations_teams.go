package manifest

import "net/http"

var teamsOperations = operationsForGroup("Teams",
	defineOperation("get-channel-teams", Operation{
		Name:           "Get Channel Teams",
		Method:         http.MethodGet,
		Path:           "/helix/teams/channel",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "broadcaster_id", Type: "String", Required: true},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 404}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Teams",
			ServiceType: "TeamsService",
			Method:      "GetChannelTeams",
			Signature:   "func (s *TeamsService) GetChannelTeams(ctx context.Context, req GetChannelTeamsRequest) (*Response[GetChannelTeamsData], error)",
			RequestType: "GetChannelTeamsRequest",
			DataType:    "GetChannelTeamsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-channel-teams",
	}),
	defineOperation("get-teams", Operation{
		Name:           "Get Teams",
		Method:         http.MethodGet,
		Path:           "/helix/teams",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "name", Type: "String", Required: true},
					{Name: "id", Type: "String", Required: true},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 404}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Teams",
			ServiceType: "TeamsService",
			Method:      "GetTeams",
			Signature:   "func (s *TeamsService) GetTeams(ctx context.Context, req GetTeamsRequest) (*Response[GetTeamsData], error)",
			RequestType: "GetTeamsRequest",
			DataType:    "GetTeamsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-teams",
	}),
)
