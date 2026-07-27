package helix

import "context"

type GetUserChatColorRequest struct {
	UserIDs []string `query:"user_id"`
}

type ChatUserColor struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	Color     string `json:"color"`
}

type GetUserChatColorData []ChatUserColor

type UpdateUserChatColorRequest struct {
	UserID string `query:"user_id"`
	Color  string `query:"color"`
}

type UpdateUserChatColorData struct{}

func (s *ChatService) GetUserChatColor(ctx context.Context, req GetUserChatColorRequest) (*Response[GetUserChatColorData], error) {
	return executeChatEndpoint[GetUserChatColorData](chatEndpointSpec{client: s.client, ctx: ctx, anchor: "get-user-chat-color", query: req})
}

func (s *ChatService) UpdateUserChatColor(ctx context.Context, req UpdateUserChatColorRequest) (*Response[UpdateUserChatColorData], error) {
	return executeChatEndpoint[UpdateUserChatColorData](chatEndpointSpec{
		client: s.client,
		ctx:    ctx,
		anchor: "update-user-chat-color",
		query:  req,
		auth: chatAuthorization{
			userScopeSets: chatReadScopes(ScopeUserManageChatColor),
			subjectID:     req.UserID,
		},
	})
}
