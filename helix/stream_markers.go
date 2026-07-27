package helix

import (
	"context"
	"fmt"
)

type CreateStreamMarkerRequest struct {
	UserID      string  `json:"user_id"`
	Description *string `json:"description,omitempty"`
}

type StreamMarker struct {
	ID              string    `json:"id"`
	CreatedAt       Timestamp `json:"created_at"`
	Description     string    `json:"description"`
	PositionSeconds int       `json:"position_seconds"`
	URL             string    `json:"URL"`
}

type CreateStreamMarkerData []StreamMarker

type GetStreamMarkersRequest struct {
	UserID  *string `query:"user_id,omitempty"`
	VideoID *string `query:"video_id,omitempty"`
	First   *string `query:"first,omitempty"`
	Before  *string `query:"before,omitempty"`
	After   *string `query:"after,omitempty"`
}

type StreamMarkerVideo struct {
	VideoID string         `json:"video_id"`
	Markers []StreamMarker `json:"markers"`
}

type StreamMarkerGroup struct {
	UserID    string              `json:"user_id"`
	UserName  string              `json:"user_name"`
	UserLogin string              `json:"user_login"`
	Videos    []StreamMarkerVideo `json:"videos"`
}

type GetStreamMarkersData []StreamMarkerGroup

func (s *StreamsService) CreateStreamMarker(ctx context.Context, req CreateStreamMarkerRequest) (*Response[CreateStreamMarkerData], error) {
	return executeEndpointWithBody[CreateStreamMarkerData](s.client, ctx, "create-stream-marker", nil, req)
}

func (s *StreamsService) GetStreamMarkers(ctx context.Context, req GetStreamMarkersRequest) (*Response[GetStreamMarkersData], error) {
	if req.UserID == nil && req.VideoID == nil {
		return nil, &RequestEncodingError{Reason: "user_id or video_id is required"}
	}
	if req.UserID != nil && req.VideoID != nil {
		return nil, &ExclusiveParametersError{First: "user_id", Second: "video_id"}
	}
	return executeTask24Endpoint[GetStreamMarkersData](s.client, ctx, "get-stream-markers", req, nil, "")
}

func (s *StreamsService) GetStreamMarkersPager(req GetStreamMarkersRequest, opts ...PagerOption) (*Pager[GetStreamMarkersData], error) {
	if req.UserID == nil && req.VideoID == nil {
		return nil, &RequestEncodingError{Reason: "user_id or video_id is required"}
	}
	if req.UserID != nil && req.VideoID != nil {
		return nil, &ExclusiveParametersError{First: "user_id", Second: "video_id"}
	}
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
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetStreamMarkersData], error) {
		if cursor == "" {
			return s.GetStreamMarkers(ctx, baseRequest)
		}
		requestValue, err := plan.withCursor(baseRequest, cursor)
		if err != nil {
			return nil, err
		}
		request, ok := requestValue.(GetStreamMarkersRequest)
		if !ok {
			return nil, fmt.Errorf("helix: stream markers pagination request has type %T", requestValue)
		}
		return s.GetStreamMarkers(ctx, request)
	}, initialCursor, opts...)
}
