package helix

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type StartCommercialRequest struct {
	BroadcasterID *string `json:"broadcaster_id,omitempty"`
	Length        *int    `json:"length,omitempty"`
}

type StartCommercial struct {
	Length     int    `json:"length"`
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after"`
}

type StartCommercialData []StartCommercial

type GetAdScheduleRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type AdSchedule struct {
	SnoozeCount     int       `json:"snooze_count"`
	SnoozeRefreshAt Timestamp `json:"snooze_refresh_at"`
	NextAdAt        Timestamp `json:"next_ad_at"`
	Duration        int       `json:"duration"`
	LastAdAt        Timestamp `json:"last_ad_at"`
	PrerollFreeTime int       `json:"preroll_free_time"`
}

type GetAdScheduleData []AdSchedule

type SnoozeNextAdRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type SnoozeNextAd struct {
	SnoozeCount     int       `json:"snooze_count"`
	SnoozeRefreshAt Timestamp `json:"snooze_refresh_at"`
	NextAdAt        Timestamp `json:"next_ad_at"`
}

type SnoozeNextAdData []SnoozeNextAd

type CooldownError struct {
	operation string
	status    int
	meta      ResponseMeta
	cause     error
}

type AdCooldownError = CooldownError

func (e *CooldownError) Error() string {
	if e == nil {
		return "helix cooldown error: <nil>"
	}
	return fmt.Sprintf("helix cooldown error: operation=%s status=%d", e.operation, e.status)
}

func (e *CooldownError) Operation() string  { return e.operation }
func (e *CooldownError) StatusCode() int    { return e.status }
func (e *CooldownError) Meta() ResponseMeta { return e.meta }
func (e *CooldownError) Unwrap() error      { return e.cause }

func (e *CooldownError) RetryAfter() time.Duration {
	if e == nil {
		return 0
	}
	seconds, err := strconv.Atoi(e.meta.Header().Get("Retry-After"))
	if err != nil || seconds < 0 {
		return 0
	}
	return time.Duration(seconds) * time.Second
}

func (s *AdsService) StartCommercial(ctx context.Context, req StartCommercialRequest) (*Response[StartCommercialData], error) {
	result, err := executeEndpointWithBody[StartCommercialData](s.client, ctx, "start-commercial", nil, req)
	return result, adCooldownError(err)
}

func (s *AdsService) GetAdSchedule(ctx context.Context, req GetAdScheduleRequest) (*Response[GetAdScheduleData], error) {
	return executeEndpoint[GetAdScheduleData](s.client, ctx, "get-ad-schedule", req)
}

func (s *AdsService) SnoozeNextAd(ctx context.Context, req SnoozeNextAdRequest) (*Response[SnoozeNextAdData], error) {
	result, err := executeEndpoint[SnoozeNextAdData](s.client, ctx, "snooze-next-ad", req)
	return result, adCooldownError(err)
}

func adCooldownError(err error) error {
	if err == nil {
		return nil
	}
	var rateLimitErr *RateLimitError
	if !errors.As(err, &rateLimitErr) {
		return err
	}
	return &CooldownError{operation: rateLimitErr.Operation(), status: rateLimitErr.StatusCode(), meta: rateLimitErr.Meta(), cause: err}
}

func executeEndpointWithBody[T any](client *Client, ctx context.Context, anchor string, query, body any) (*Response[T], error) {
	if err := client.validClient(); err != nil {
		return nil, err
	}
	operation, err := manifest.OperationByAnchor(anchor)
	if err != nil {
		return nil, err
	}
	credential, err := client.credential(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateCredentialForOperation(credential, operation, "", ""); err != nil {
		return nil, err
	}
	request, err := buildRequest(requestSpec{Context: ctx, Method: operation.Method, URL: client.endpointURL(operation.Path), Query: query, Body: body})
	if err != nil {
		return nil, fmt.Errorf("build %s request: %w", operation.OperationID, err)
	}
	request.Header.Set("User-Agent", client.userAgent)
	response, meta, err := client.executor.execute(ctx, request, operation, credential)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	result, err := decodeResponse[T](response.StatusCode, response.Body, DecodeOptions{})
	if err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: []string{credential.AccessToken()}})
	}
	result.Meta = meta
	return result, nil
}
