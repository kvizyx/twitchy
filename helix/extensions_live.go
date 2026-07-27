package helix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type GetExtensionLiveChannelsRequest struct {
	ExtensionID string  `query:"extension_id"`
	First       *int    `query:"first,omitempty"`
	After       *string `query:"after,omitempty"`
}

type ExtensionLiveChannel struct {
	BroadcasterID   string `json:"broadcaster_id"`
	BroadcasterName string `json:"broadcaster_name"`
	GameName        string `json:"game_name"`
	GameID          string `json:"game_id"`
	Title           string `json:"title"`
}

type GetExtensionLiveChannelsData []ExtensionLiveChannel

type extensionLiveChannelsWire struct {
	Data       GetExtensionLiveChannelsData `json:"data"`
	Pagination string                       `json:"pagination"`
}

func (s *ExtensionsService) GetExtensionLiveChannels(ctx context.Context, req GetExtensionLiveChannelsRequest) (*Response[GetExtensionLiveChannelsData], error) {
	if err := s.client.validClient(); err != nil {
		return nil, err
	}
	operation, err := manifest.OperationByAnchor("get-extension-live-channels")
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
	if response.StatusCode == http.StatusNoContent {
		return &Response[GetExtensionLiveChannelsData]{Meta: meta}, nil
	}
	data, err := readBounded(response.Body, BodyLimits{}.responseLimit())
	if err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: []string{credential.AccessToken()}})
	}
	var wire extensionLiveChannelsWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, body: data, cause: err, secrets: []string{credential.AccessToken()}})
	}
	return &Response[GetExtensionLiveChannelsData]{Data: wire.Data, Pagination: Pagination{cursor: wire.Pagination}, Meta: meta}, nil
}

func (s *ExtensionsService) GetExtensionLiveChannelsPager(req GetExtensionLiveChannelsRequest, opts ...PagerOption) (*Pager[GetExtensionLiveChannelsData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	return newPager(func(ctx context.Context, cursor string) (*Response[GetExtensionLiveChannelsData], error) {
		request := any(req)
		if cursor != "" {
			cursorRequest, cursorErr := plan.withCursor(req, cursor)
			if cursorErr != nil {
				return nil, cursorErr
			}
			request = cursorRequest
		}
		return s.GetExtensionLiveChannels(ctx, request.(GetExtensionLiveChannelsRequest))
	}, opts...)
}
