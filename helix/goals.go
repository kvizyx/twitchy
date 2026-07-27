package helix

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type GetCreatorGoalsRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type CreatorGoal struct {
	ID               string    `json:"id"`
	BroadcasterID    string    `json:"broadcaster_id"`
	BroadcasterName  string    `json:"broadcaster_name"`
	BroadcasterLogin string    `json:"broadcaster_login"`
	Type             string    `json:"type"`
	Description      string    `json:"description"`
	CurrentAmount    int       `json:"current_amount"`
	TargetAmount     int       `json:"target_amount"`
	CreatedAt        Timestamp `json:"created_at"`
}

type GetCreatorGoalsData []CreatorGoal

func (s *GoalsService) GetCreatorGoals(ctx context.Context, req GetCreatorGoalsRequest) (*Response[GetCreatorGoalsData], error) {
	return executeBroadcasterEndpoint[GetCreatorGoalsData](s.client, ctx, "get-creator-goals", req, req.BroadcasterID)
}

func executeBroadcasterEndpoint[T any](client *Client, ctx context.Context, anchor string, query any, broadcasterID string) (*Response[T], error) {
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
	if credential.UserID() != broadcasterID {
		return nil, localCredentialAuthError(operation.OperationID)
	}
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
	result, err := decodeResponse[T](response.StatusCode, response.Body, DecodeOptions{})
	if err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: []string{credential.AccessToken()}})
	}
	result.Meta = meta
	return result, nil
}
