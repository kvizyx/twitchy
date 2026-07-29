package manifest

import "net/http"

var tagsOperations = operationsForGroup("Tags",
	defineOperation("get-all-stream-tags", Operation{
		Name:           "Get All Stream Tags",
		Method:         http.MethodGet,
		Path:           "/helix/tags/streams",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "tag_id", Type: "String"},
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
			Selector:       "Client.Tags",
			ServiceType:    "TagsService",
			Method:         "GetAllStreamTags",
			Signature:      "func (s *TagsService) GetAllStreamTags(ctx context.Context, req GetAllStreamTagsRequest) (*Response[GetAllStreamTagsData], error)",
			PagerSignature: pagerSignature("func (s *TagsService) GetAllStreamTagsPager(req GetAllStreamTagsRequest, opts ...PagerOption) (*Pager[GetAllStreamTagsData], error)"),
			RequestType:    "GetAllStreamTagsRequest",
			DataType:       "GetAllStreamTagsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-all-stream-tags",
	}),
	defineOperation("get-stream-tags", Operation{
		Name:           "Get Stream Tags",
		Method:         http.MethodGet,
		Path:           "/helix/streams/tags",
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
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Tags",
			ServiceType: "TagsService",
			Method:      "GetStreamTags",
			Signature:   "func (s *TagsService) GetStreamTags(ctx context.Context, req GetStreamTagsRequest) (*Response[GetStreamTagsData], error)",
			RequestType: "GetStreamTagsRequest",
			DataType:    "GetStreamTagsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-stream-tags",
	}),
)
