package manifest

import "net/http"

var videosOperations = operationsForGroup("Videos",
	defineOperation("get-videos", Operation{
		Name:           "Get Videos",
		Method:         http.MethodGet,
		Path:           "/helix/videos",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "id", Type: "String", Required: true},
					{Name: "user_id", Type: "String", Required: true},
					{Name: "game_id", Type: "String", Required: true},
					{Name: "language", Type: "String"},
					{Name: "period", Type: "String"},
					{Name: "sort", Type: "String"},
					{Name: "type", Type: "String"},
					{Name: "first", Type: "String"},
					{Name: "after", Type: "String"},
					{Name: "before", Type: "String"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 404}},
		Pagination: PaginationSpec{Shape: "cursor", CursorParameter: "after"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:       "Client.Videos",
			ServiceType:    "VideosService",
			Method:         "GetVideos",
			Signature:      "func (s *VideosService) GetVideos(ctx context.Context, req GetVideosRequest) (*Response[GetVideosData], error)",
			PagerSignature: pagerSignature("func (s *VideosService) GetVideosPager(req GetVideosRequest, opts ...PagerOption) (*Pager[GetVideosData], error)"),
			RequestType:    "GetVideosRequest",
			DataType:       "GetVideosData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-videos",
	}),
	defineOperation("delete-videos", Operation{
		Name:           "Delete Videos",
		Method:         http.MethodDelete,
		Path:           "/helix/videos",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:manage:videos"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "id", Type: "String", Required: true},
				},
			},
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Videos",
			ServiceType: "VideosService",
			Method:      "DeleteVideos",
			Signature:   "func (s *VideosService) DeleteVideos(ctx context.Context, req DeleteVideosRequest) (*Response[DeleteVideosData], error)",
			RequestType: "DeleteVideosRequest",
			DataType:    "DeleteVideosData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#delete-videos",
	}),
)
