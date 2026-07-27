package helix

import "context"

type GetAllStreamTagsRequest struct {
	TagIDs []string `query:"tag_id,omitempty"`
	First  *int     `query:"first,omitempty"`
	After  *string  `query:"after,omitempty"`
}

type StreamTag struct {
	TagID                    string            `json:"tag_id"`
	IsAuto                   bool              `json:"is_auto"`
	LocalizationNames        map[string]string `json:"localization_names"`
	LocalizationDescriptions map[string]string `json:"localization_descriptions"`
}

type GetAllStreamTagsData []StreamTag

type GetStreamTagsRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type GetStreamTagsData []StreamTag

func (s *TagsService) GetAllStreamTags(ctx context.Context, req GetAllStreamTagsRequest) (*Response[GetAllStreamTagsData], error) {
	return executeEndpoint[GetAllStreamTagsData](s.client, ctx, "get-all-stream-tags", req)
}

func (s *TagsService) GetAllStreamTagsPager(req GetAllStreamTagsRequest, opts ...PagerOption) (*Pager[GetAllStreamTagsData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	initialCursor := ""
	baseRequest := req
	if req.After != nil {
		initialCursor = *req.After
		baseRequest.After = nil
	}
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetAllStreamTagsData], error) {
		if cursor == "" {
			return s.GetAllStreamTags(ctx, baseRequest)
		}
		requestValue, err := plan.withCursor(baseRequest, cursor)
		if err != nil {
			return nil, err
		}
		request, ok := requestValue.(GetAllStreamTagsRequest)
		if !ok {
			return nil, &paginationRequestError{reason: "tags pagination request has unexpected type"}
		}
		return s.GetAllStreamTags(ctx, request)
	}, initialCursor, opts...)
}

func (s *TagsService) GetStreamTags(ctx context.Context, req GetStreamTagsRequest) (*Response[GetStreamTagsData], error) {
	return executeEndpoint[GetStreamTagsData](s.client, ctx, "get-stream-tags", req)
}
