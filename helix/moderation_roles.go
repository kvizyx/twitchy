package helix

import "context"

type GetModeratedChannelsRequest struct {
	UserID string  `query:"user_id"`
	After  *string `query:"after,omitempty"`
	First  *int    `query:"first,omitempty"`
}

type ModeratedChannel struct {
	BroadcasterID    string `json:"broadcaster_id"`
	BroadcasterLogin string `json:"broadcaster_login"`
	BroadcasterName  string `json:"broadcaster_name"`
}

type GetModeratedChannelsData []ModeratedChannel

type GetModeratorsRequest struct {
	BroadcasterID string   `query:"broadcaster_id"`
	UserIDs       []string `query:"user_id,omitempty"`
	First         *int     `query:"first,omitempty"`
	After         *string  `query:"after,omitempty"`
}

type Moderator struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type GetModeratorsData []Moderator

type AddChannelModeratorRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	UserID        string `query:"user_id"`
}

type AddChannelModeratorData struct{}

type RemoveChannelModeratorRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	UserID        string `query:"user_id"`
}

type RemoveChannelModeratorData struct{}

type GetVIPsRequest struct {
	UserIDs       []string `query:"user_id,omitempty"`
	BroadcasterID string   `query:"broadcaster_id"`
	First         *int     `query:"first,omitempty"`
	After         *string  `query:"after,omitempty"`
}

type VIP struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login"`
}

type ChannelVIP = VIP

type GetVIPsData []VIP

type AddChannelVIPRequest struct {
	UserID        string `query:"user_id"`
	BroadcasterID string `query:"broadcaster_id"`
}

type AddChannelVIPData struct{}

type RemoveChannelVIPRequest struct {
	UserID        string `query:"user_id"`
	BroadcasterID string `query:"broadcaster_id"`
}

type RemoveChannelVIPData struct{}

func (s *ModerationService) GetModeratedChannels(ctx context.Context, req GetModeratedChannelsRequest) (*Response[GetModeratedChannelsData], error) {
	return executeModerationEndpoint[GetModeratedChannelsData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "get-moderated-channels", query: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeUserReadModeratedChannels), subjectIDs: []string{req.UserID}}})
}

func (s *ModerationService) GetModeratedChannelsPager(req GetModeratedChannelsRequest, opts ...PagerOption) (*Pager[GetModeratedChannelsData], error) {
	initialCursor := ""
	if req.After != nil {
		initialCursor = *req.After
	}
	request := req
	request.After = nil
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetModeratedChannelsData], error) {
		pageRequest := request
		if cursor != "" {
			pageRequest.After = &cursor
		}
		return s.GetModeratedChannels(ctx, pageRequest)
	}, initialCursor, opts...)
}

func (s *ModerationService) GetModerators(ctx context.Context, req GetModeratorsRequest) (*Response[GetModeratorsData], error) {
	scopes := [][]AuthorizationScope{{ScopeModerationRead}, {ScopeChannelManageModerators}}
	return executeModerationEndpoint[GetModeratorsData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "get-moderators", query: req, auth: moderationAuthorization{userScopeSets: scopes, subjectIDs: []string{req.BroadcasterID}}})
}

func (s *ModerationService) GetModeratorsPager(req GetModeratorsRequest, opts ...PagerOption) (*Pager[GetModeratorsData], error) {
	initialCursor := ""
	if req.After != nil {
		initialCursor = *req.After
	}
	request := req
	request.After = nil
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetModeratorsData], error) {
		pageRequest := request
		if cursor != "" {
			pageRequest.After = &cursor
		}
		return s.GetModerators(ctx, pageRequest)
	}, initialCursor, opts...)
}

func (s *ModerationService) AddChannelModerator(ctx context.Context, req AddChannelModeratorRequest) (*Response[AddChannelModeratorData], error) {
	return executeModerationEndpoint[AddChannelModeratorData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "add-channel-moderator", query: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeChannelManageModerators), subjectIDs: []string{req.BroadcasterID}}})
}

func (s *ModerationService) RemoveChannelModerator(ctx context.Context, req RemoveChannelModeratorRequest) (*Response[RemoveChannelModeratorData], error) {
	return executeModerationEndpoint[RemoveChannelModeratorData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "remove-channel-moderator", query: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeChannelManageModerators), subjectIDs: []string{req.BroadcasterID}}})
}

func (s *ModerationService) GetVIPs(ctx context.Context, req GetVIPsRequest) (*Response[GetVIPsData], error) {
	scopes := [][]AuthorizationScope{{ScopeChannelReadVIPs}, {ScopeChannelManageVIPs}}
	return executeModerationEndpoint[GetVIPsData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "get-vips", query: req, auth: moderationAuthorization{userScopeSets: scopes, subjectIDs: []string{req.BroadcasterID}}})
}

func (s *ModerationService) GetVIPsPager(req GetVIPsRequest, opts ...PagerOption) (*Pager[GetVIPsData], error) {
	initialCursor := ""
	if req.After != nil {
		initialCursor = *req.After
	}
	request := req
	request.After = nil
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetVIPsData], error) {
		pageRequest := request
		if cursor != "" {
			pageRequest.After = &cursor
		}
		return s.GetVIPs(ctx, pageRequest)
	}, initialCursor, opts...)
}

func (s *ModerationService) AddChannelVIP(ctx context.Context, req AddChannelVIPRequest) (*Response[AddChannelVIPData], error) {
	return executeModerationEndpoint[AddChannelVIPData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "add-channel-vip", query: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeChannelManageVIPs), subjectIDs: []string{req.BroadcasterID}}})
}

func (s *ModerationService) RemoveChannelVIP(ctx context.Context, req RemoveChannelVIPRequest) (*Response[RemoveChannelVIPData], error) {
	return executeModerationEndpoint[RemoveChannelVIPData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "remove-channel-vip", query: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeChannelManageVIPs), subjectIDs: []string{req.BroadcasterID, req.UserID}}})
}
