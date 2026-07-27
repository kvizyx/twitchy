package helix

import "context"

type AddSuspiciousStatusToChatUserRequest struct {
	BroadcasterID string `query:"broadcaster_id" json:"-"`
	ModeratorID   string `query:"moderator_id" json:"-"`
	UserID        string `query:"-" json:"user_id"`
	Status        string `query:"-" json:"status"`
}

type SuspiciousUserStatus struct {
	UserID        string    `json:"user_id"`
	BroadcasterID string    `json:"broadcaster_id"`
	ModeratorID   string    `json:"moderator_id"`
	UpdatedAt     Timestamp `json:"updated_at"`
	Status        string    `json:"status"`
	Types         []string  `json:"types"`
}

type SuspiciousUser = SuspiciousUserStatus

type AddSuspiciousStatusToChatUserData []SuspiciousUserStatus

type RemoveSuspiciousStatusFromChatUserRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
	UserID        string `query:"user_id"`
}

type RemoveSuspiciousStatusFromChatUserData []SuspiciousUserStatus

func (s *ExperimentalModerationService) AddSuspiciousStatusToChatUser(ctx context.Context, req AddSuspiciousStatusToChatUserRequest) (*Response[AddSuspiciousStatusToChatUserData], error) {
	return executeModerationEndpoint[AddSuspiciousStatusToChatUserData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "add-suspicious-status-to-chat-user", query: req, body: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModeratorManageSuspiciousUsers), appScopeSets: chatReadScopes(ScopeModeratorManageSuspiciousUsers), subjectIDs: []string{req.ModeratorID}}})
}

func (s *ExperimentalModerationService) RemoveSuspiciousStatusFromChatUser(ctx context.Context, req RemoveSuspiciousStatusFromChatUserRequest) (*Response[RemoveSuspiciousStatusFromChatUserData], error) {
	return executeModerationEndpoint[RemoveSuspiciousStatusFromChatUserData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "remove-suspicious-status-from-chat-user", query: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModeratorManageSuspiciousUsers), appScopeSets: chatReadScopes(ScopeModeratorManageSuspiciousUsers), subjectIDs: []string{req.ModeratorID}}})
}
