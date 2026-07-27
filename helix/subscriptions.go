package helix

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type SubscriptionTier string

const (
	SubscriptionTier1000 SubscriptionTier = "1000"
	SubscriptionTier2000 SubscriptionTier = "2000"
	SubscriptionTier3000 SubscriptionTier = "3000"
)

type GetBroadcasterSubscriptionsRequest struct {
	BroadcasterID string   `query:"broadcaster_id"`
	UserIDs       []string `query:"user_id,omitempty"`
	First         *string  `query:"first,omitempty"`
	After         *string  `query:"after,omitempty"`
	Before        *string  `query:"before,omitempty"`
}

type BroadcasterSubscription struct {
	BroadcasterID    string           `json:"broadcaster_id"`
	BroadcasterLogin string           `json:"broadcaster_login"`
	BroadcasterName  string           `json:"broadcaster_name"`
	GifterID         string           `json:"gifter_id"`
	GifterLogin      string           `json:"gifter_login"`
	GifterName       string           `json:"gifter_name"`
	IsGift           bool             `json:"is_gift"`
	PlanName         string           `json:"plan_name"`
	Tier             SubscriptionTier `json:"tier"`
	UserID           string           `json:"user_id"`
	UserName         string           `json:"user_name"`
	UserLogin        string           `json:"user_login"`
}

type GetBroadcasterSubscriptionsData struct {
	Subscriptions []BroadcasterSubscription `json:"data"`
	Points        *int                      `json:"points"`
	Total         *int                      `json:"total"`
}

type CheckUserSubscriptionRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	UserID        string `query:"user_id"`
}

type UserSubscription struct {
	BroadcasterID    string           `json:"broadcaster_id"`
	BroadcasterLogin string           `json:"broadcaster_login"`
	BroadcasterName  string           `json:"broadcaster_name"`
	GifterID         string           `json:"gifter_id"`
	GifterLogin      string           `json:"gifter_login"`
	GifterName       string           `json:"gifter_name"`
	IsGift           bool             `json:"is_gift"`
	Tier             SubscriptionTier `json:"tier"`
}

type CheckUserSubscriptionData []UserSubscription

type broadcasterSubscriptionsWire struct {
	Data       []BroadcasterSubscription `json:"data"`
	Pagination Pagination                `json:"pagination"`
	Points     *int                      `json:"points"`
	Total      *int                      `json:"total"`
}

func (s *SubscriptionsService) GetBroadcasterSubscriptions(ctx context.Context, req GetBroadcasterSubscriptionsRequest) (*Response[GetBroadcasterSubscriptionsData], error) {
	if err := s.client.validClient(); err != nil {
		return nil, err
	}
	operation, err := manifest.OperationByAnchor("get-broadcaster-subscriptions")
	if err != nil {
		return nil, err
	}
	credential, err := s.client.credential(ctx)
	if err != nil {
		return nil, err
	}
	if credential.TokenClass() == TokenClassApp {
		operation.Scopes = nil
	}
	if err := validateCredentialForOperation(credential, operation, "", ""); err != nil {
		return nil, err
	}
	request, err := buildRequest(requestSpec{Context: ctx, Method: operation.Method, URL: s.client.endpointURL(operation.Path), Query: req})
	if err != nil {
		return nil, fmt.Errorf("build %s request: %w", operation.OperationID, err)
	}
	request.Header.Set("User-Agent", s.client.userAgent)
	response, meta, err := s.client.executor.execute(ctx, request, operation, credential)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := readBounded(response.Body, BodyLimits{}.responseLimit())
	if err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: []string{credential.AccessToken()}})
	}
	var wire broadcasterSubscriptionsWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, body: data, cause: err, secrets: []string{credential.AccessToken()}})
	}
	return &Response[GetBroadcasterSubscriptionsData]{Data: GetBroadcasterSubscriptionsData{Subscriptions: wire.Data, Points: wire.Points, Total: wire.Total}, Pagination: wire.Pagination, Meta: meta}, nil
}

func (s *SubscriptionsService) GetBroadcasterSubscriptionsPager(req GetBroadcasterSubscriptionsRequest, opts ...PagerOption) (*Pager[GetBroadcasterSubscriptionsData], error) {
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
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetBroadcasterSubscriptionsData], error) {
		if cursor == "" {
			return s.GetBroadcasterSubscriptions(ctx, baseRequest)
		}
		requestValue, err := plan.withCursor(baseRequest, cursor)
		if err != nil {
			return nil, err
		}
		request, ok := requestValue.(GetBroadcasterSubscriptionsRequest)
		if !ok {
			return nil, fmt.Errorf("helix: subscriptions pagination request has type %T", requestValue)
		}
		return s.GetBroadcasterSubscriptions(ctx, request)
	}, initialCursor, opts...)
}

func (s *SubscriptionsService) CheckUserSubscription(ctx context.Context, req CheckUserSubscriptionRequest) (*Response[CheckUserSubscriptionData], error) {
	if err := s.client.validClient(); err != nil {
		return nil, err
	}
	operation, err := manifest.OperationByAnchor("check-user-subscription")
	if err != nil {
		return nil, err
	}
	credential, err := s.client.credential(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateCredentialForOperation(credential, operation, "", ""); err != nil {
		return nil, err
	}
	if credential.UserID() != req.UserID {
		return nil, localCredentialAuthError(operation.OperationID)
	}
	request, err := buildRequest(requestSpec{Context: ctx, Method: operation.Method, URL: s.client.endpointURL(operation.Path), Query: req})
	if err != nil {
		return nil, fmt.Errorf("build %s request: %w", operation.OperationID, err)
	}
	request.Header.Set("User-Agent", s.client.userAgent)
	response, meta, err := s.client.executor.execute(ctx, request, operation, credential)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	result, err := decodeResponse[CheckUserSubscriptionData](response.StatusCode, response.Body, DecodeOptions{})
	if err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: []string{credential.AccessToken()}})
	}
	result.Meta = meta
	return result, nil
}
