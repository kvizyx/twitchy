package helix

import (
	"context"
	"fmt"
)

type SearchCategoriesRequest struct {
	Query string  `query:"query"`
	First *int    `query:"first,omitempty"`
	After *string `query:"after,omitempty"`
}

type SearchCategory struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"box_art_url"`
}

type SearchCategoriesData []SearchCategory

type SearchChannelsRequest struct {
	Query    string  `query:"query"`
	LiveOnly *bool   `query:"live_only,omitempty"`
	First    *int    `query:"first,omitempty"`
	After    *string `query:"after,omitempty"`
}

type SearchChannel struct {
	BroadcasterLanguage string   `json:"broadcaster_language"`
	BroadcasterLogin    string   `json:"broadcaster_login"`
	DisplayName         string   `json:"display_name"`
	GameID              string   `json:"game_id"`
	GameName            string   `json:"game_name"`
	ID                  string   `json:"id"`
	IsLive              bool     `json:"is_live"`
	TagIDs              []string `json:"tag_ids"`
	Tags                []string `json:"tags"`
	ThumbnailURL        string   `json:"thumbnail_url"`
	Title               string   `json:"title"`
	StartedAt           string   `json:"started_at"`
}

type SearchChannelsData []SearchChannel

func (s *SearchService) SearchCategories(ctx context.Context, req SearchCategoriesRequest) (*Response[SearchCategoriesData], error) {
	return executeEndpoint[SearchCategoriesData](s.client, ctx, "search-categories", req)
}

func (s *SearchService) SearchCategoriesPager(req SearchCategoriesRequest, opts ...PagerOption) (*Pager[SearchCategoriesData], error) {
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
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[SearchCategoriesData], error) {
		if cursor == "" {
			return s.SearchCategories(ctx, baseRequest)
		}
		requestValue, err := plan.withCursor(baseRequest, cursor)
		if err != nil {
			return nil, err
		}
		request, ok := requestValue.(SearchCategoriesRequest)
		if !ok {
			return nil, fmt.Errorf("helix: search categories pagination request has type %T", requestValue)
		}
		return s.SearchCategories(ctx, request)
	}, initialCursor, opts...)
}

func (s *SearchService) SearchChannels(ctx context.Context, req SearchChannelsRequest) (*Response[SearchChannelsData], error) {
	return executeEndpoint[SearchChannelsData](s.client, ctx, "search-channels", req)
}

func (s *SearchService) SearchChannelsPager(req SearchChannelsRequest, opts ...PagerOption) (*Pager[SearchChannelsData], error) {
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
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[SearchChannelsData], error) {
		if cursor == "" {
			return s.SearchChannels(ctx, baseRequest)
		}
		requestValue, err := plan.withCursor(baseRequest, cursor)
		if err != nil {
			return nil, err
		}
		request, ok := requestValue.(SearchChannelsRequest)
		if !ok {
			return nil, fmt.Errorf("helix: search channels pagination request has type %T", requestValue)
		}
		return s.SearchChannels(ctx, request)
	}, initialCursor, opts...)
}
