package manifest

import "net/http"

var searchOperations = operationsForGroup("Search",
	defineOperation("search-categories", Operation{
		Name:           "Search Categories",
		Method:         http.MethodGet,
		Path:           "/helix/search/categories",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "query", Type: "String", Required: true},
					{Name: "first", Type: "Integer"},
					{Name: "after", Type: "String"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401}},
		Pagination: PaginationSpec{Shape: "cursor", CursorParameter: "after"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:       "Client.Search",
			ServiceType:    "SearchService",
			Method:         "SearchCategories",
			Signature:      "func (s *SearchService) SearchCategories(ctx context.Context, req SearchCategoriesRequest) (*Response[SearchCategoriesData], error)",
			PagerSignature: pagerSignature("func (s *SearchService) SearchCategoriesPager(req SearchCategoriesRequest, opts ...PagerOption) (*Pager[SearchCategoriesData], error)"),
			RequestType:    "SearchCategoriesRequest",
			DataType:       "SearchCategoriesData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#search-categories",
	}),
	defineOperation("search-channels", Operation{
		Name:           "Search Channels",
		Method:         http.MethodGet,
		Path:           "/helix/search/channels",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "query", Type: "String", Required: true},
					{Name: "live_only", Type: "Boolean"},
					{Name: "first", Type: "Integer"},
					{Name: "after", Type: "String"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401}},
		Pagination: PaginationSpec{Shape: "cursor", CursorParameter: "after"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:       "Client.Search",
			ServiceType:    "SearchService",
			Method:         "SearchChannels",
			Signature:      "func (s *SearchService) SearchChannels(ctx context.Context, req SearchChannelsRequest) (*Response[SearchChannelsData], error)",
			PagerSignature: pagerSignature("func (s *SearchService) SearchChannelsPager(req SearchChannelsRequest, opts ...PagerOption) (*Pager[SearchChannelsData], error)"),
			RequestType:    "SearchChannelsRequest",
			DataType:       "SearchChannelsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#search-channels",
	}),
)
