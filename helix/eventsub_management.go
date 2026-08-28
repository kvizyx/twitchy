package helix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type EventSubTransportMethod string

const (
	EventSubTransportWebhook   EventSubTransportMethod = "webhook"
	EventSubTransportWebsocket EventSubTransportMethod = "websocket"
	EventSubTransportConduit   EventSubTransportMethod = "conduit"
)

type EventSubCondition map[string]string

type EventSubTransport struct {
	Method    EventSubTransportMethod `json:"method"`
	Callback  string                  `json:"callback,omitempty"`
	Secret    string                  `json:"secret,omitempty"`
	SessionID string                  `json:"session_id,omitempty"`
	ConduitID string                  `json:"conduit_id,omitempty"`
}

type EventSubSubscription struct {
	ID        string            `json:"id"`
	Status    string            `json:"status"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition EventSubCondition `json:"condition"`
	Transport EventSubTransport `json:"transport"`
	CreatedAt TimestampUTC      `json:"created_at"`
	Cost      int               `json:"cost"`
}

type EventSubSubscriptionsData struct {
	Subscriptions []EventSubSubscription `json:"data"`
	Total         int                    `json:"total"`
	TotalCost     int                    `json:"total_cost"`
	MaxTotalCost  int                    `json:"max_total_cost"`
}

type CreateEventSubSubscriptionRequest struct {
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition EventSubCondition `json:"condition"`
	Transport EventSubTransport `json:"transport"`
}

type CreateEventSubSubscriptionData = EventSubSubscriptionsData

type DeleteEventSubSubscriptionRequest struct {
	ID              string                  `query:"id"`
	Transport       EventSubTransport       `query:"-"`
	TransportMethod EventSubTransportMethod `query:"-"`
}

type DeleteEventSubSubscriptionData struct{}

type GetEventSubSubscriptionsRequest struct {
	Status          string                  `query:"status,omitempty"`
	Type            string                  `query:"type,omitempty"`
	UserID          string                  `query:"user_id,omitempty"`
	SubscriptionID  string                  `query:"subscription_id,omitempty"`
	ConduitID       string                  `query:"conduit_id,omitempty"`
	After           string                  `query:"after,omitempty"`
	First           int                     `query:"-"`
	Transport       EventSubTransport       `query:"-"`
	TransportMethod EventSubTransportMethod `query:"-"`
}

type GetEventSubSubscriptionsData = EventSubSubscriptionsData

func (s *EventSubService) CreateEventSubSubscription(ctx context.Context, req CreateEventSubSubscriptionRequest) (*Response[CreateEventSubSubscriptionData], error) {
	operation, credential, err := s.eventSubCredential(ctx, "create-eventsub-subscription", req.Transport.Method)
	if err != nil {
		return nil, err
	}
	result, err := executeEventSubEnvelope[CreateEventSubSubscriptionData](s.client, ctx, operation, credential, req, []string{string(req.Transport.Secret)})
	return result, redactEventSubError(err, req.Transport.Secret)
}

func (s *EventSubService) DeleteEventSubSubscription(ctx context.Context, req DeleteEventSubSubscriptionRequest) (*Response[DeleteEventSubSubscriptionData], error) {
	method := req.Transport.Method
	if method == "" {
		method = req.TransportMethod
	}
	operation, credential, err := s.eventSubCredential(ctx, "delete-eventsub-subscription", method)
	if err != nil {
		return nil, err
	}
	return executeConduitLikeNoBody[DeleteEventSubSubscriptionData](s.client, ctx, operation, credential, req)
}

func (s *EventSubService) GetEventSubSubscriptions(ctx context.Context, req GetEventSubSubscriptionsRequest) (*Response[GetEventSubSubscriptionsData], error) {
	if req.First != 0 {
		return nil, &RequestEncodingError{Reason: "first is not supported; use after"}
	}
	if filterCount(req.Status, req.Type, req.UserID, req.SubscriptionID, req.ConduitID) > 1 {
		return nil, &RequestEncodingError{Reason: "EventSub subscription filters are mutually exclusive"}
	}
	method := req.Transport.Method
	if method == "" {
		method = req.TransportMethod
	}
	operation, credential, err := s.eventSubCredential(ctx, "get-eventsub-subscriptions", method)
	if err != nil {
		return nil, err
	}
	return executeEventSubEnvelope[GetEventSubSubscriptionsData](s.client, ctx, operation, credential, req, nil)
}

func (s *EventSubService) GetEventSubSubscriptionsPager(req GetEventSubSubscriptionsRequest, opts ...PagerOption) (*Pager[GetEventSubSubscriptionsData], error) {
	return newPager(func(ctx context.Context, cursor string) (*Response[GetEventSubSubscriptionsData], error) {
		req.After = cursor
		return s.GetEventSubSubscriptions(ctx, req)
	}, opts...)
}

func (s *EventSubService) eventSubCredential(ctx context.Context, anchor string, method EventSubTransportMethod) (manifest.Operation, CredentialSnapshot, error) {
	if method != "" && method != EventSubTransportWebhook && method != EventSubTransportWebsocket && method != EventSubTransportConduit {
		return manifest.Operation{}, CredentialSnapshot{}, &RequestEncodingError{Reason: fmt.Sprintf("unsupported EventSub transport method %q", method)}
	}
	operation, err := manifest.OperationByAnchor(anchor)
	if err != nil {
		return manifest.Operation{}, CredentialSnapshot{}, err
	}
	credential, err := s.client.credential(ctx)
	if err != nil {
		return manifest.Operation{}, CredentialSnapshot{}, err
	}
	if err := validateCredentialForOperation(credential, operation, "", ""); err != nil {
		return manifest.Operation{}, CredentialSnapshot{}, err
	}
	wanted := TokenClassApp
	if method == EventSubTransportWebsocket {
		wanted = TokenClassUser
	}
	if method != "" && credential.TokenClass() != wanted {
		return manifest.Operation{}, CredentialSnapshot{}, localCredentialAuthError(operation.OperationID)
	}
	return operation, credential, nil
}

func executeEventSubEnvelope[T any](client *Client, ctx context.Context, operation manifest.Operation, credential CredentialSnapshot, body any, secrets []string) (*Response[T], error) {
	requestBody := body
	if operation.Method == "GET" {
		requestBody = nil
	}
	request, err := buildRequest(requestSpec{Context: ctx, Method: operation.Method, URL: client.endpointURL(operation.Path), Body: requestBody, Query: bodyQuery(body)})
	if err != nil {
		return nil, fmt.Errorf("build %s request: %w", operation.OperationID, err)
	}
	request.Header.Set("User-Agent", client.userAgent)
	response, meta, err := client.executor.execute(ctx, request, operation, credential)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: append([]string{credential.AccessToken()}, secrets...)})
	}
	var envelope EventSubSubscriptionsData
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: append([]string{credential.AccessToken()}, secrets...)})
	}
	var wire struct {
		Pagination Pagination `json:"pagination"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: append([]string{credential.AccessToken()}, secrets...)})
	}
	resultData, ok := any(envelope).(T)
	if !ok {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: fmt.Errorf("decode %s response data type", operation.OperationID), secrets: append([]string{credential.AccessToken()}, secrets...)})
	}
	result := &Response[T]{Data: resultData, Pagination: wire.Pagination, Meta: meta}
	return result, nil
}

func executeConduitLikeNoBody[T any](client *Client, ctx context.Context, operation manifest.Operation, credential CredentialSnapshot, query any) (*Response[T], error) {
	request, err := buildRequest(requestSpec{Context: ctx, Method: operation.Method, URL: client.endpointURL(operation.Path), Query: query})
	if err != nil {
		return nil, fmt.Errorf("build %s request: %w", operation.OperationID, err)
	}
	request.Header.Set("User-Agent", client.userAgent)
	response, meta, err := client.executor.execute(ctx, request, operation, credential)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return &Response[T]{Meta: meta}, nil
}

func bodyQuery(value any) any {
	switch request := value.(type) {
	case DeleteEventSubSubscriptionRequest:
		return request
	case GetEventSubSubscriptionsRequest:
		return request
	default:
		return nil
	}
}

func filterCount(values ...string) int {
	count := 0
	for _, value := range values {
		if value != "" {
			count++
		}
	}
	return count
}

func redactEventSubError(err error, secrets ...string) error {
	if err == nil {
		return nil
	}
	switch typed := err.(type) {
	case *TransportError:
		typed.secrets = append(typed.secrets, secrets...)
	case *ProtocolError:
		typed.details.secrets = append(typed.details.secrets, secrets...)
	case *APIError:
		typed.details.secrets = append(typed.details.secrets, secrets...)
	case *AuthError:
		typed.details.secrets = append(typed.details.secrets, secrets...)
	case *RateLimitError:
		typed.details.secrets = append(typed.details.secrets, secrets...)
	}
	return err
}
