package helix

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type GetChannelStreamScheduleRequest struct {
	BroadcasterID string   `query:"broadcaster_id"`
	IDs           []string `query:"id,omitempty"`
	StartTime     *string  `query:"start_time,omitempty"`
	UTCOffset     *string  `query:"utc_offset,omitempty"`
	First         *int     `query:"first,omitempty"`
	After         *string  `query:"after,omitempty"`
}

type GetChannelICalendarRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type UpdateChannelStreamScheduleRequest struct {
	BroadcasterID     string  `query:"broadcaster_id"`
	IsVacationEnabled *bool   `query:"is_vacation_enabled,omitempty"`
	VacationStartTime *string `query:"vacation_start_time,omitempty"`
	VacationEndTime   *string `query:"vacation_end_time,omitempty"`
	Timezone          *string `query:"timezone,omitempty"`
}

type ScheduleCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ScheduleSegment struct {
	ID            string            `json:"id"`
	StartTime     Timestamp         `json:"start_time"`
	EndTime       Timestamp         `json:"end_time"`
	Title         string            `json:"title"`
	CanceledUntil *Timestamp        `json:"canceled_until"`
	Category      *ScheduleCategory `json:"category"`
	IsRecurring   bool              `json:"is_recurring"`
}

type ScheduleVacation struct {
	StartTime Timestamp `json:"start_time"`
	EndTime   Timestamp `json:"end_time"`
}

type scheduleData struct {
	Segments         []ScheduleSegment `json:"segments"`
	BroadcasterID    string            `json:"broadcaster_id"`
	BroadcasterName  string            `json:"broadcaster_name"`
	BroadcasterLogin string            `json:"broadcaster_login"`
	Vacation         *ScheduleVacation `json:"vacation"`
}

type GetChannelStreamScheduleData scheduleData
type UpdateChannelStreamScheduleData struct{}
type CreateChannelStreamScheduleSegmentData scheduleData
type UpdateChannelStreamScheduleSegmentData scheduleData
type DeleteChannelStreamScheduleSegmentData struct{}
type GetChannelICalendarData []byte

func (s *ScheduleService) GetChannelStreamSchedule(ctx context.Context, req GetChannelStreamScheduleRequest) (*Response[GetChannelStreamScheduleData], error) {
	return executeEndpoint[GetChannelStreamScheduleData](s.client, ctx, "get-channel-stream-schedule", req)
}

func (s *ScheduleService) GetChannelStreamSchedulePager(req GetChannelStreamScheduleRequest, opts ...PagerOption) (*Pager[GetChannelStreamScheduleData], error) {
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
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetChannelStreamScheduleData], error) {
		if cursor == "" {
			return s.GetChannelStreamSchedule(ctx, baseRequest)
		}
		requestValue, err := plan.withCursor(baseRequest, cursor)
		if err != nil {
			return nil, err
		}
		request, ok := requestValue.(GetChannelStreamScheduleRequest)
		if !ok {
			return nil, fmt.Errorf("helix: schedule pagination request has type %T", requestValue)
		}
		return s.GetChannelStreamSchedule(ctx, request)
	}, initialCursor, opts...)
}

func (s *ScheduleService) GetChannelICalendar(ctx context.Context, req GetChannelICalendarRequest) (*Response[GetChannelICalendarData], error) {
	if err := s.client.validClient(); err != nil {
		return nil, err
	}
	operation, err := manifest.OperationByAnchor("get-channel-icalendar")
	if err != nil {
		return nil, err
	}
	request, err := buildRequest(requestSpec{Context: ctx, Method: operation.Method, URL: s.client.endpointURL(operation.Path), Query: req})
	if err != nil {
		return nil, fmt.Errorf("build %s request: %w", operation.OperationID, err)
	}
	request.Header.Set("User-Agent", s.client.userAgent)
	response, meta, err := s.client.executor.execute(ctx, request, operation, CredentialSnapshot{})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := readBounded(response.Body, BodyLimits{}.responseLimit())
	if err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err})
	}
	return &Response[GetChannelICalendarData]{Data: data, Meta: meta}, nil
}

func (s *ScheduleService) UpdateChannelStreamSchedule(ctx context.Context, req UpdateChannelStreamScheduleRequest) (*Response[UpdateChannelStreamScheduleData], error) {
	return executeTask24Endpoint[UpdateChannelStreamScheduleData](s.client, ctx, "update-channel-stream-schedule", req, nil, req.BroadcasterID)
}
