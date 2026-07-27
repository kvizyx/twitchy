package helix

import "context"

type GetBannedUsersRequest struct {
	BroadcasterID string   `query:"broadcaster_id"`
	UserIDs       []string `query:"user_id,omitempty"`
	First         *int     `query:"first,omitempty"`
	After         *string  `query:"after,omitempty"`
	Before        *string  `query:"before,omitempty"`
}

type BannedUser struct {
	UserID         string `json:"user_id"`
	UserLogin      string `json:"user_login"`
	UserName       string `json:"user_name"`
	ExpiresAt      string `json:"expires_at"`
	CreatedAt      string `json:"created_at"`
	Reason         string `json:"reason"`
	ModeratorID    string `json:"moderator_id"`
	ModeratorLogin string `json:"moderator_login"`
	ModeratorName  string `json:"moderator_name"`
}

type GetBannedUsersData []BannedUser

type BanUserBody struct {
	UserID   string `json:"user_id"`
	Duration *int   `json:"duration,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

type BanUser = BanUserBody

type BanUserRequest struct {
	BroadcasterID string      `query:"broadcaster_id" json:"-"`
	ModeratorID   string      `query:"moderator_id" json:"-"`
	Data          BanUserBody `query:"-" json:"data"`
}

type Ban struct {
	BroadcasterID string     `json:"broadcaster_id"`
	ModeratorID   string     `json:"moderator_id"`
	UserID        string     `json:"user_id"`
	CreatedAt     Timestamp  `json:"created_at"`
	EndTime       *Timestamp `json:"end_time"`
}

type BanUserData []Ban

func (s *ModerationService) GetBannedUsers(ctx context.Context, req GetBannedUsersRequest) (*Response[GetBannedUsersData], error) {
	scopes := [][]AuthorizationScope{{ScopeModerationRead}, {ScopeModeratorManageBannedUsers}}
	return executeModerationEndpoint[GetBannedUsersData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "get-banned-users", query: req, auth: moderationAuthorization{userScopeSets: scopes, appScopeSets: scopes, subjectIDs: []string{req.BroadcasterID}}})
}

func (s *ModerationService) GetBannedUsersPager(req GetBannedUsersRequest, opts ...PagerOption) (*Pager[GetBannedUsersData], error) {
	initialCursor := ""
	if req.After != nil {
		initialCursor = *req.After
	}
	request := req
	request.After = nil
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetBannedUsersData], error) {
		pageRequest := request
		pageRequest.Before = nil
		if cursor != "" {
			pageRequest.After = &cursor
		}
		return s.GetBannedUsers(ctx, pageRequest)
	}, initialCursor, opts...)
}

func (s *ModerationService) BanUser(ctx context.Context, req BanUserRequest) (*Response[BanUserData], error) {
	return executeModerationEndpoint[BanUserData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "ban-user", query: req, body: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModeratorManageBannedUsers), appScopeSets: [][]AuthorizationScope{{ScopeModeratorManageBannedUsers, ScopeUserBot}}, subjectIDs: []string{req.ModeratorID}}})
}
