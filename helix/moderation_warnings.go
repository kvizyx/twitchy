package helix

import "context"

type WarnChatUserBody struct {
	UserID string `json:"user_id"`
	Reason string `json:"reason"`
}

type WarnChatUser = WarnChatUserBody

type WarnChatUserRequest struct {
	BroadcasterID string           `query:"broadcaster_id" json:"-"`
	ModeratorID   string           `query:"moderator_id" json:"-"`
	Data          WarnChatUserBody `query:"-" json:"data"`
}

type ChatWarning struct {
	BroadcasterID string `json:"broadcaster_id"`
	UserID        string `json:"user_id"`
	ModeratorID   string `json:"moderator_id"`
	Reason        string `json:"reason"`
}

type WarnChatUserData []ChatWarning

func (s *ModerationService) WarnChatUser(ctx context.Context, req WarnChatUserRequest) (*Response[WarnChatUserData], error) {
	return executeModerationEndpoint[WarnChatUserData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "warn-chat-user", query: req, body: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModeratorManageWarnings), appScopeSets: chatReadScopes(ScopeModeratorManageWarnings), subjectIDs: []string{req.ModeratorID}}})
}
