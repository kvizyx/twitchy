package helix

import "context"

type GetExtensionAnalyticsRequest struct {
	ExtensionID string  `query:"extension_id,omitempty"`
	Type        string  `query:"type,omitempty"`
	StartedAt   *string `query:"started_at,omitempty"`
	EndedAt     *string `query:"ended_at,omitempty"`
	First       *int    `query:"first,omitempty"`
	After       *string `query:"after,omitempty"`
}

type AnalyticsDateRange struct {
	StartedAt Timestamp `json:"started_at"`
	EndedAt   Timestamp `json:"ended_at"`
}

type ExtensionAnalytics struct {
	ExtensionID string             `json:"extension_id"`
	URL         string             `json:"URL"`
	Type        string             `json:"type"`
	DateRange   AnalyticsDateRange `json:"date_range"`
}

type GetExtensionAnalyticsData []ExtensionAnalytics

type GetGameAnalyticsRequest struct {
	GameID    string  `query:"game_id,omitempty"`
	Type      string  `query:"type,omitempty"`
	StartedAt *string `query:"started_at,omitempty"`
	EndedAt   *string `query:"ended_at,omitempty"`
	First     *int    `query:"first,omitempty"`
	After     *string `query:"after,omitempty"`
}

type GameAnalytics struct {
	GameID    string             `json:"game_id"`
	URL       string             `json:"URL"`
	Type      string             `json:"type"`
	DateRange AnalyticsDateRange `json:"date_range"`
}

type GetGameAnalyticsData []GameAnalytics

func (s *AnalyticsService) GetExtensionAnalytics(ctx context.Context, req GetExtensionAnalyticsRequest) (*Response[GetExtensionAnalyticsData], error) {
	return executeEndpoint[GetExtensionAnalyticsData](s.client, ctx, "get-extension-analytics", req)
}

func (s *AnalyticsService) GetExtensionAnalyticsPager(req GetExtensionAnalyticsRequest, opts ...PagerOption) (*Pager[GetExtensionAnalyticsData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	return newPager(func(ctx context.Context, cursor string) (*Response[GetExtensionAnalyticsData], error) {
		request := any(req)
		if cursor != "" {
			cursorRequest, cursorErr := plan.withCursor(req, cursor)
			if cursorErr != nil {
				return nil, cursorErr
			}
			request = cursorRequest
		}
		return s.GetExtensionAnalytics(ctx, request.(GetExtensionAnalyticsRequest))
	}, opts...)
}

func (s *AnalyticsService) GetGameAnalytics(ctx context.Context, req GetGameAnalyticsRequest) (*Response[GetGameAnalyticsData], error) {
	return executeEndpoint[GetGameAnalyticsData](s.client, ctx, "get-game-analytics", req)
}

func (s *AnalyticsService) GetGameAnalyticsPager(req GetGameAnalyticsRequest, opts ...PagerOption) (*Pager[GetGameAnalyticsData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	return newPager(func(ctx context.Context, cursor string) (*Response[GetGameAnalyticsData], error) {
		request := any(req)
		if cursor != "" {
			cursorRequest, cursorErr := plan.withCursor(req, cursor)
			if cursorErr != nil {
				return nil, cursorErr
			}
			request = cursorRequest
		}
		return s.GetGameAnalytics(ctx, request.(GetGameAnalyticsRequest))
	}, opts...)
}
