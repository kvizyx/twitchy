package helix

import (
	"context"
	"fmt"
)

type GetTopGamesRequest struct {
	First  *int    `query:"first,omitempty"`
	After  *string `query:"after,omitempty"`
	Before *string `query:"before,omitempty"`
}

type GetTopGamesData []Game

func (s *GamesService) GetTopGames(ctx context.Context, req GetTopGamesRequest) (*Response[GetTopGamesData], error) {
	return executeEndpoint[GetTopGamesData](s.client, ctx, "get-top-games", req)
}

func (s *GamesService) GetTopGamesPager(req GetTopGamesRequest, opts ...PagerOption) (*Pager[GetTopGamesData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	initialCursor := ""
	baseRequest := req
	if req.Before != nil {
		initialCursor = *req.Before
		baseRequest.Before = nil
	} else if req.After != nil {
		initialCursor = *req.After
		baseRequest.After = nil
	}
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetTopGamesData], error) {
		if cursor == "" {
			return s.GetTopGames(ctx, baseRequest)
		}
		requestValue, err := plan.withCursor(baseRequest, cursor)
		if err != nil {
			return nil, err
		}
		request, ok := requestValue.(GetTopGamesRequest)
		if !ok {
			return nil, fmt.Errorf("helix: games pagination request has type %T", requestValue)
		}
		return s.GetTopGames(ctx, request)
	}, initialCursor, opts...)
}
