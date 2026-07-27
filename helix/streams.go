package helix

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type GetStreamKeyRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type GetStreamKeyData []StreamKey

type StreamKey struct {
	StreamKey string `json:"stream_key"`
}

type GetStreamsRequest struct {
	UserID    []string `query:"user_id,omitempty"`
	UserLogin []string `query:"user_login,omitempty"`
	GameID    []string `query:"game_id,omitempty"`
	Type      *string  `query:"type,omitempty"`
	Language  []string `query:"language,omitempty"`
	First     *int     `query:"first,omitempty"`
	Before    *string  `query:"before,omitempty"`
	After     *string  `query:"after,omitempty"`
}

type StreamType string

const (
	StreamTypeAll  StreamType = "all"
	StreamTypeLive StreamType = "live"
)

type Stream struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	UserLogin    string     `json:"user_login"`
	UserName     string     `json:"user_name"`
	GameID       string     `json:"game_id"`
	GameName     string     `json:"game_name"`
	Type         StreamType `json:"type"`
	Title        string     `json:"title"`
	Tags         []string   `json:"tags"`
	ViewerCount  int        `json:"viewer_count"`
	StartedAt    Timestamp  `json:"started_at"`
	Language     string     `json:"language"`
	ThumbnailURL string     `json:"thumbnail_url"`
	TagIDs       []string   `json:"tag_ids"`
	IsMature     bool       `json:"is_mature"`
}

type GetStreamsData []Stream

type GetFollowedStreamsRequest struct {
	UserID string  `query:"user_id"`
	First  *int    `query:"first,omitempty"`
	After  *string `query:"after,omitempty"`
}

type GetFollowedStreamsData []Stream

func (s *StreamsService) GetStreamKey(ctx context.Context, req GetStreamKeyRequest) (*Response[GetStreamKeyData], error) {
	return executeTask24Endpoint[GetStreamKeyData](s.client, ctx, "get-stream-key", req, nil, req.BroadcasterID)
}

func (s *StreamsService) GetStreams(ctx context.Context, req GetStreamsRequest) (*Response[GetStreamsData], error) {
	return executeEndpoint[GetStreamsData](s.client, ctx, "get-streams", req)
}

func (s *StreamsService) GetStreamsPager(req GetStreamsRequest, opts ...PagerOption) (*Pager[GetStreamsData], error) {
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
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetStreamsData], error) {
		if cursor == "" {
			return s.GetStreams(ctx, baseRequest)
		}
		requestValue, err := plan.withCursor(baseRequest, cursor)
		if err != nil {
			return nil, err
		}
		request, ok := requestValue.(GetStreamsRequest)
		if !ok {
			return nil, fmt.Errorf("helix: streams pagination request has type %T", requestValue)
		}
		return s.GetStreams(ctx, request)
	}, initialCursor, opts...)
}

func (s *StreamsService) GetFollowedStreams(ctx context.Context, req GetFollowedStreamsRequest) (*Response[GetFollowedStreamsData], error) {
	return executeTask24Endpoint[GetFollowedStreamsData](s.client, ctx, "get-followed-streams", req, nil, req.UserID)
}

func (s *StreamsService) GetFollowedStreamsPager(req GetFollowedStreamsRequest, opts ...PagerOption) (*Pager[GetFollowedStreamsData], error) {
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
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetFollowedStreamsData], error) {
		if cursor == "" {
			return s.GetFollowedStreams(ctx, baseRequest)
		}
		requestValue, err := plan.withCursor(baseRequest, cursor)
		if err != nil {
			return nil, err
		}
		request, ok := requestValue.(GetFollowedStreamsRequest)
		if !ok {
			return nil, fmt.Errorf("helix: followed streams pagination request has type %T", requestValue)
		}
		return s.GetFollowedStreams(ctx, request)
	}, initialCursor, opts...)
}

func executeTask24Endpoint[T any](client *Client, ctx context.Context, anchor string, query, body any, subjectID string) (*Response[T], error) {
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
	if anchor == "get-stream-markers" {
		if !operationAllowsTokenClass(credential.TokenClass(), operation.TokenClasses) ||
			(!snapshotHasScope(credential, ScopeUserReadBroadcast) && !snapshotHasScope(credential, ScopeChannelManageBroadcast)) {
			return nil, localCredentialAuthError(operation.OperationID)
		}
	} else if err := validateCredentialForOperation(credential, operation, "", ""); err != nil {
		return nil, err
	}
	if subjectID != "" && credential.TokenClass() == TokenClassUser && credential.UserID() != subjectID {
		return nil, localCredentialAuthError(operation.OperationID)
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
