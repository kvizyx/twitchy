package helix

import "context"

type CreateChannelStreamScheduleSegmentRequest struct {
	BroadcasterID string  `query:"broadcaster_id" json:"-"`
	StartTime     string  `query:"-" json:"start_time"`
	Timezone      string  `query:"-" json:"timezone"`
	Duration      string  `query:"-" json:"duration"`
	IsRecurring   *bool   `query:"-" json:"is_recurring,omitempty"`
	CategoryID    *string `query:"-" json:"category_id,omitempty"`
	Title         *string `query:"-" json:"title,omitempty"`
}

type UpdateChannelStreamScheduleSegmentRequest struct {
	BroadcasterID string  `query:"broadcaster_id" json:"-"`
	ID            string  `query:"id" json:"-"`
	StartTime     *string `query:"-" json:"start_time,omitempty"`
	Duration      *string `query:"-" json:"duration,omitempty"`
	CategoryID    *string `query:"-" json:"category_id,omitempty"`
	Title         *string `query:"-" json:"title,omitempty"`
	IsCanceled    *bool   `query:"-" json:"is_canceled,omitempty"`
	Timezone      *string `query:"-" json:"timezone,omitempty"`
}

type DeleteChannelStreamScheduleSegmentRequest struct {
	BroadcasterID string `query:"broadcaster_id" json:"-"`
	ID            string `query:"id" json:"-"`
}

func (s *ScheduleService) CreateChannelStreamScheduleSegment(ctx context.Context, req CreateChannelStreamScheduleSegmentRequest) (*Response[CreateChannelStreamScheduleSegmentData], error) {
	return executeTask24Endpoint[CreateChannelStreamScheduleSegmentData](s.client, ctx, "create-channel-stream-schedule-segment", req, req, req.BroadcasterID)
}

func (s *ScheduleService) UpdateChannelStreamScheduleSegment(ctx context.Context, req UpdateChannelStreamScheduleSegmentRequest) (*Response[UpdateChannelStreamScheduleSegmentData], error) {
	return executeTask24Endpoint[UpdateChannelStreamScheduleSegmentData](s.client, ctx, "update-channel-stream-schedule-segment", req, req, req.BroadcasterID)
}

func (s *ScheduleService) DeleteChannelStreamScheduleSegment(ctx context.Context, req DeleteChannelStreamScheduleSegmentRequest) (*Response[DeleteChannelStreamScheduleSegmentData], error) {
	return executeTask24Endpoint[DeleteChannelStreamScheduleSegmentData](s.client, ctx, "delete-channel-stream-schedule-segment", req, nil, req.BroadcasterID)
}
