package helix

import "context"

type DeleteChatMessagesRequest struct {
	BroadcasterID string  `query:"broadcaster_id"`
	ModeratorID   string  `query:"moderator_id"`
	MessageID     string  `query:"message_id,omitempty"`
}

type DeleteChatMessagesData struct{}

func (s *ModerationService) DeleteChatMessages(ctx context.Context, req DeleteChatMessagesRequest) (*Response[DeleteChatMessagesData], error) {
	return executeModerationEndpoint[DeleteChatMessagesData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "delete-chat-messages", query: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModeratorManageChatMessages), appScopeSets: chatReadScopes(ScopeModeratorManageChatMessages), subjectIDs: []string{req.ModeratorID}}})
}
