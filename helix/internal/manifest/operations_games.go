package manifest

import "net/http"

var gamesOperations = operationsForGroup("Games",
	defineOperation("get-top-games", Operation{
		Name:           "Get Top Games",
		Method:         http.MethodGet,
		Path:           "/helix/games/top",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "first", Type: "Integer"},
					{Name: "after", Type: "String"},
					{Name: "before", Type: "String"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401}},
		Pagination: PaginationSpec{Shape: "cursor", CursorParameter: "after"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:       "Client.Games",
			ServiceType:    "GamesService",
			Method:         "GetTopGames",
			Signature:      "func (s *GamesService) GetTopGames(ctx context.Context, req GetTopGamesRequest) (*Response[GetTopGamesData], error)",
			PagerSignature: pagerSignature("func (s *GamesService) GetTopGamesPager(req GetTopGamesRequest, opts ...PagerOption) (*Pager[GetTopGamesData], error)"),
			RequestType:    "GetTopGamesRequest",
			DataType:       "GetTopGamesData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-top-games",
	}),
	defineOperation("get-games", Operation{
		Name:           "Get Games",
		Method:         http.MethodGet,
		Path:           "/helix/games",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "id", Type: "String", Required: true},
					{Name: "name", Type: "String", Required: true},
					{Name: "igdb_id", Type: "String", Required: true},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Games",
			ServiceType: "GamesService",
			Method:      "GetGames",
			Signature:   "func (s *GamesService) GetGames(ctx context.Context, req GetGamesRequest) (*Response[GetGamesData], error)",
			RequestType: "GetGamesRequest",
			DataType:    "GetGamesData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-games",
	}),
)
