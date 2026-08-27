package helix

import "context"

type UnbanUserRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
	UserID        string `query:"user_id"`
}

type UnbanUserData struct{}

type GetUnbanRequestsRequest struct {
	BroadcasterID string  `query:"broadcaster_id"`
	ModeratorID   string  `query:"moderator_id"`
	Status        string  `query:"status"`
	UserID        string  `query:"user_id,omitempty"`
	After         *string `query:"after,omitempty"`
	First         *int    `query:"first,omitempty"`
}

type UnbanRequest struct {
	ID               string     `json:"id"`
	BroadcasterID    string     `json:"broadcaster_id"`
	BroadcasterLogin string     `json:"broadcaster_login"`
	BroadcasterName  string     `json:"broadcaster_name"`
	ModeratorID      string     `json:"moderator_id"`
	ModeratorLogin   string     `json:"moderator_login"`
	ModeratorName    string     `json:"moderator_name"`
	UserID           string     `json:"user_id"`
	UserLogin        string     `json:"user_login"`
	UserName         string     `json:"user_name"`
	Text             string     `json:"text"`
	Status           string     `json:"status"`
	CreatedAt        Timestamp  `json:"created_at"`
	ResolvedAt       *Timestamp `json:"resolved_at"`
	ResolutionText   *string    `json:"resolution_text"`
}

type GetUnbanRequestsData []UnbanRequest
type ResolveUnbanRequestsData []UnbanRequest

type ResolveUnbanRequestsRequest struct {
	BroadcasterID  string `query:"broadcaster_id"`
	ModeratorID    string `query:"moderator_id"`
	UnbanRequestID string `query:"unban_request_id"`
	Status         string `query:"status"`
	ResolutionText string `query:"resolution_text,omitempty"`
}

func (s *ModerationService) UnbanUser(ctx context.Context, req UnbanUserRequest) (*Response[UnbanUserData], error) {
	return executeModerationEndpoint[UnbanUserData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "unban-user", query: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModeratorManageBannedUsers), subjectIDs: []string{req.ModeratorID}}})
}

func (s *ModerationService) GetUnbanRequests(ctx context.Context, req GetUnbanRequestsRequest) (*Response[GetUnbanRequestsData], error) {
	scopes := [][]AuthorizationScope{{ScopeModeratorReadUnbanRequests}, {ScopeModeratorManageUnbanRequests}}
	return executeModerationEndpoint[GetUnbanRequestsData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "get-unban-requests", query: req, auth: moderationAuthorization{userScopeSets: scopes, subjectIDs: []string{req.ModeratorID}}})
}

func (s *ModerationService) GetUnbanRequestsPager(req GetUnbanRequestsRequest, opts ...PagerOption) (*Pager[GetUnbanRequestsData], error) {
	initialCursor := ""
	if req.After != nil {
		initialCursor = *req.After
	}
	request := req
	request.After = nil
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetUnbanRequestsData], error) {
		pageRequest := request
		if cursor != "" {
			pageRequest.After = &cursor
		}
		return s.GetUnbanRequests(ctx, pageRequest)
	}, initialCursor, opts...)
}

func (s *ModerationService) ResolveUnbanRequests(ctx context.Context, req ResolveUnbanRequestsRequest) (*Response[ResolveUnbanRequestsData], error) {
	return executeModerationEndpoint[ResolveUnbanRequestsData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "resolve-unban-requests", query: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModeratorManageUnbanRequests), subjectIDs: []string{req.ModeratorID}}})
}
