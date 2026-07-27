package helix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type GetConduitsRequest struct{}

type Conduit struct {
	ID         string `json:"id"`
	ShardCount int    `json:"shard_count"`
}

type GetConduitsData []Conduit

type CreateConduitsRequest struct {
	ShardCount int `json:"shard_count"`
}

type CreateConduitsData []Conduit

type UpdateConduitsRequest struct {
	ID         string `json:"id"`
	ShardCount int    `json:"shard_count"`
}

type UpdateConduitsData []Conduit

type DeleteConduitRequest struct {
	ID string `query:"id"`
}

type DeleteConduitData struct{}

type GetConduitShardsRequest struct {
	ConduitID string `query:"conduit_id"`
	Status    string `query:"status,omitempty"`
	After     string `query:"after,omitempty"`
}

type ConduitShardTransport struct {
	Method         string        `json:"method"`
	Callback       string        `json:"callback,omitempty"`
	Secret         string        `json:"secret,omitempty"`
	SessionID      string        `json:"session_id,omitempty"`
	ConnectedAt    *TimestampUTC `json:"connected_at,omitempty"`
	DisconnectedAt *TimestampUTC `json:"disconnected_at,omitempty"`
}

type ConduitShard struct {
	ID        string                `json:"id"`
	Status    string                `json:"status"`
	Transport ConduitShardTransport `json:"transport"`
}

type GetConduitShardsData []ConduitShard

type UpdateConduitShard struct {
	ID        string                `json:"id"`
	Transport ConduitShardTransport `json:"transport"`
}

type UpdateConduitShardsRequest struct {
	ConduitID string               `json:"conduit_id"`
	Shards    []UpdateConduitShard `json:"shards"`
}

type ConduitShardError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

type UpdateConduitShardsData struct {
	Shards []ConduitShard      `json:"data"`
	Errors []ConduitShardError `json:"errors"`
}

func (s *ConduitsService) GetConduits(ctx context.Context, req GetConduitsRequest) (*Response[GetConduitsData], error) {
	return executeConduitEndpoint[GetConduitsData](s.client, ctx, "get-conduits", req, nil)
}

func (s *ConduitsService) CreateConduits(ctx context.Context, req CreateConduitsRequest) (*Response[CreateConduitsData], error) {
	return executeConduitEndpoint[CreateConduitsData](s.client, ctx, "create-conduits", nil, req)
}

func (s *ConduitsService) UpdateConduits(ctx context.Context, req UpdateConduitsRequest) (*Response[UpdateConduitsData], error) {
	return executeConduitEndpoint[UpdateConduitsData](s.client, ctx, "update-conduits", nil, req)
}

func (s *ConduitsService) DeleteConduit(ctx context.Context, req DeleteConduitRequest) (*Response[DeleteConduitData], error) {
	return executeConduitEndpoint[DeleteConduitData](s.client, ctx, "delete-conduit", req, nil)
}

func (s *ConduitsService) GetConduitShards(ctx context.Context, req GetConduitShardsRequest) (*Response[GetConduitShardsData], error) {
	return executeConduitEndpoint[GetConduitShardsData](s.client, ctx, "get-conduit-shards", req, nil)
}

func (s *ConduitsService) GetConduitShardsPager(req GetConduitShardsRequest, opts ...PagerOption) (*Pager[GetConduitShardsData], error) {
	return newPager(func(ctx context.Context, cursor string) (*Response[GetConduitShardsData], error) {
		req.After = cursor
		return s.GetConduitShards(ctx, req)
	}, opts...)
}

func (s *ConduitsService) UpdateConduitShards(ctx context.Context, req UpdateConduitShardsRequest) (*Response[UpdateConduitShardsData], error) {
	return executeConduitEnvelope(s.client, ctx, "update-conduit-shards", req)
}

func executeConduitEndpoint[T any](client *Client, ctx context.Context, anchor string, query, body any) (*Response[T], error) {
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

func executeConduitEnvelope(client *Client, ctx context.Context, anchor string, body any) (*Response[UpdateConduitShardsData], error) {
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
	request, err := buildRequest(requestSpec{Context: ctx, Method: operation.Method, URL: client.endpointURL(operation.Path), Body: body})
	if err != nil {
		return nil, fmt.Errorf("build %s request: %w", operation.OperationID, err)
	}
	request.Header.Set("User-Agent", client.userAgent)
	response, meta, err := client.executor.execute(ctx, request, operation, credential)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var wire UpdateConduitShardsData
	if response.StatusCode != http.StatusNoContent {
		data, readErr := io.ReadAll(response.Body)
		if readErr != nil {
			return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: readErr, secrets: []string{credential.AccessToken()}})
		}
		if err := json.Unmarshal(data, &wire); err != nil {
			return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: []string{credential.AccessToken()}})
		}
	}
	return &Response[UpdateConduitShardsData]{Data: wire, Meta: meta}, nil
}
