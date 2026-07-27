package helix

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type GetChannelInformationRequest struct {
	BroadcasterIDs []string `query:"broadcaster_id"`
}

type ChannelInformation struct {
	BroadcasterID               string   `json:"broadcaster_id"`
	BroadcasterLogin            string   `json:"broadcaster_login"`
	BroadcasterName             string   `json:"broadcaster_name"`
	BroadcasterLanguage         string   `json:"broadcaster_language"`
	GameName                    string   `json:"game_name"`
	GameID                      string   `json:"game_id"`
	Title                       string   `json:"title"`
	Delay                       int      `json:"delay"`
	Tags                        []string `json:"tags"`
	ContentClassificationLabels []string `json:"content_classification_labels"`
	IsBrandedContent            bool     `json:"is_branded_content"`
}

type GetChannelInformationData []ChannelInformation

type ChannelClassificationLabel struct {
	ID        string `json:"id"`
	IsEnabled bool   `json:"is_enabled"`
}

type ModifyChannelInformationRequest struct {
	BroadcasterID               string                        `query:"broadcaster_id" json:"-"`
	GameID                      *string                       `query:"-" json:"game_id,omitempty"`
	BroadcasterLanguage         *string                       `query:"-" json:"broadcaster_language,omitempty"`
	Title                       *string                       `query:"-" json:"title,omitempty"`
	Delay                       *int                          `query:"-" json:"delay,omitempty"`
	Tags                        *[]string                     `query:"-" json:"tags,omitempty"`
	ContentClassificationLabels *[]ChannelClassificationLabel `query:"-" json:"content_classification_labels,omitempty"`
	IsBrandedContent            *bool                         `query:"-" json:"is_branded_content,omitempty"`
}

type ModifyChannelInformationData struct{}

type GetChannelEditorsRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type ChannelEditor struct {
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	CreatedAt Timestamp `json:"created_at"`
}

type GetChannelEditorsData []ChannelEditor

type GetFollowedChannelsRequest struct {
	UserID        string  `query:"user_id"`
	BroadcasterID *string `query:"broadcaster_id,omitempty"`
	First         *int    `query:"first,omitempty"`
	After         *string `query:"after,omitempty"`
}

type FollowedChannel struct {
	BroadcasterID    string    `json:"broadcaster_id"`
	BroadcasterLogin string    `json:"broadcaster_login"`
	BroadcasterName  string    `json:"broadcaster_name"`
	FollowedAt       Timestamp `json:"followed_at"`
}

type GetFollowedChannelsData []FollowedChannel

type GetChannelFollowersRequest struct {
	UserID        *string `query:"user_id,omitempty"`
	BroadcasterID string  `query:"broadcaster_id"`
	First         *int    `query:"first,omitempty"`
	After         *string `query:"after,omitempty"`
}

type ChannelFollower struct {
	FollowedAt Timestamp `json:"followed_at"`
	UserID     string    `json:"user_id"`
	UserLogin  string    `json:"user_login"`
	UserName   string    `json:"user_name"`
}

type GetChannelFollowersData []ChannelFollower

func (s *ChannelsService) GetChannelInformation(ctx context.Context, req GetChannelInformationRequest) (*Response[GetChannelInformationData], error) {
	return executeEndpointRequest[GetChannelInformationData](s.client, ctx, "get-channel-information", req, nil, "")
}

func (s *ChannelsService) ModifyChannelInformation(ctx context.Context, req ModifyChannelInformationRequest) (*Response[ModifyChannelInformationData], error) {
	return executeEndpointRequest[ModifyChannelInformationData](s.client, ctx, "modify-channel-information", req, req, "")
}

func (s *ChannelsService) GetChannelEditors(ctx context.Context, req GetChannelEditorsRequest) (*Response[GetChannelEditorsData], error) {
	return executeEndpointRequest[GetChannelEditorsData](s.client, ctx, "get-channel-editors", req, nil, "")
}

func (s *ChannelsService) GetFollowedChannels(ctx context.Context, req GetFollowedChannelsRequest) (*Response[GetFollowedChannelsData], error) {
	return executeEndpointRequest[GetFollowedChannelsData](s.client, ctx, "get-followed-channels", req, nil, "")
}

func (s *ChannelsService) GetFollowedChannelsPager(req GetFollowedChannelsRequest, opts ...PagerOption) (*Pager[GetFollowedChannelsData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	return newPager(func(ctx context.Context, cursor string) (*Response[GetFollowedChannelsData], error) {
		next, err := plan.withCursor(req, cursor)
		if err != nil {
			return nil, err
		}
		return s.GetFollowedChannels(ctx, next.(GetFollowedChannelsRequest))
	}, opts...)
}

func (s *ChannelsService) GetChannelFollowers(ctx context.Context, req GetChannelFollowersRequest) (*Response[GetChannelFollowersData], error) {
	subjectID := ""
	if req.UserID != nil {
		subjectID = *req.UserID
	}
	return executeEndpointRequest[GetChannelFollowersData](s.client, ctx, "get-channel-followers", req, nil, subjectID)
}

func (s *ChannelsService) GetChannelFollowersPager(req GetChannelFollowersRequest, opts ...PagerOption) (*Pager[GetChannelFollowersData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	return newPager(func(ctx context.Context, cursor string) (*Response[GetChannelFollowersData], error) {
		next, err := plan.withCursor(req, cursor)
		if err != nil {
			return nil, err
		}
		return s.GetChannelFollowers(ctx, next.(GetChannelFollowersRequest))
	}, opts...)
}

func executeEndpointRequest[T any](client *Client, ctx context.Context, anchor string, query, body any, subjectID string) (*Response[T], error) {
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
	if err := validateCredentialForOperation(credential, operation, "", subjectID); err != nil {
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
