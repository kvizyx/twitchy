package helix

import "context"

type GetPinnedChatMessageRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
}

type ChatCheermote struct {
	Prefix string `json:"prefix"`
	Bits   int    `json:"bits"`
	Tier   int    `json:"tier"`
}

type ChatMessageEmote struct {
	ID         string   `json:"id"`
	EmoteSetID string   `json:"emote_set_id"`
	OwnerID    string   `json:"owner_id"`
	Format     []string `json:"format"`
}

type ChatMessageMention struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type ChatMessageFragment struct {
	Type      string              `json:"type"`
	Text      string              `json:"text"`
	Cheermote *ChatCheermote      `json:"cheermote"`
	Emote     *ChatMessageEmote   `json:"emote"`
	Mention   *ChatMessageMention `json:"mention"`
}

type ChatMessageContent struct {
	Text      string                `json:"text"`
	Fragments []ChatMessageFragment `json:"fragments"`
}

type PinnedChatMessage struct {
	MessageID         string             `json:"message_id"`
	BroadcasterID     string             `json:"broadcaster_id"`
	SenderUserID      string             `json:"sender_user_id"`
	SenderUserLogin   string             `json:"sender_user_login"`
	SenderUserName    string             `json:"sender_user_name"`
	PinnedByUserID    string             `json:"pinned_by_user_id"`
	PinnedByUserLogin string             `json:"pinned_by_user_login"`
	PinnedByUserName  string             `json:"pinned_by_user_name"`
	Message           ChatMessageContent `json:"message"`
	StartsAt          Timestamp          `json:"starts_at"`
	EndsAt            *Timestamp         `json:"ends_at"`
	UpdatedAt         Timestamp          `json:"updated_at"`
}

type GetPinnedChatMessageData []PinnedChatMessage

type PinChatMessageRequest struct {
	BroadcasterID   string `query:"broadcaster_id"`
	ModeratorID     string `query:"moderator_id"`
	MessageID       string `query:"message_id"`
	DurationSeconds *int   `query:"duration_seconds,omitempty"`
}

type PinChatMessageData struct{}

type UpdatePinnedChatMessageRequest struct {
	BroadcasterID   string `query:"broadcaster_id"`
	ModeratorID     string `query:"moderator_id"`
	MessageID       string `query:"message_id"`
	DurationSeconds *int   `query:"duration_seconds,omitempty"`
}

type UpdatePinnedChatMessageData struct{}

type UnpinChatMessageRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
	MessageID     string `query:"message_id"`
}

type UnpinChatMessageData struct{}

func (s *ExperimentalChatService) GetPinnedChatMessage(ctx context.Context, req GetPinnedChatMessageRequest) (*Response[GetPinnedChatMessageData], error) {
	return executeChatEndpoint[GetPinnedChatMessageData](chatEndpointSpec{
		client: s.client,
		ctx:    ctx,
		anchor: "get-pinned-chat-message",
		query:  req,
		auth: chatAuthorization{
			userScopeSets: [][]AuthorizationScope{{ScopeModeratorManageChatMessages}, {ScopeModeratorReadChatMessages}},
			subjectID:     req.ModeratorID,
		},
	})
}

func (s *ExperimentalChatService) PinChatMessage(ctx context.Context, req PinChatMessageRequest) (*Response[PinChatMessageData], error) {
	return executeChatEndpoint[PinChatMessageData](chatEndpointSpec{
		client: s.client,
		ctx:    ctx,
		anchor: "pin-chat-message",
		query:  req,
		auth: chatAuthorization{
			userScopeSets: chatReadScopes(ScopeModeratorManageChatMessages),
			subjectID:     req.ModeratorID,
		},
	})
}

func (s *ExperimentalChatService) UpdatePinnedChatMessage(ctx context.Context, req UpdatePinnedChatMessageRequest) (*Response[UpdatePinnedChatMessageData], error) {
	return executeChatEndpoint[UpdatePinnedChatMessageData](chatEndpointSpec{
		client: s.client,
		ctx:    ctx,
		anchor: "update-pinned-chat-message",
		query:  req,
		auth: chatAuthorization{
			userScopeSets: chatReadScopes(ScopeModeratorManageChatMessages),
			subjectID:     req.ModeratorID,
		},
	})
}

func (s *ExperimentalChatService) UnpinChatMessage(ctx context.Context, req UnpinChatMessageRequest) (*Response[UnpinChatMessageData], error) {
	return executeChatEndpoint[UnpinChatMessageData](chatEndpointSpec{
		client: s.client,
		ctx:    ctx,
		anchor: "unpin-chat-message",
		query:  req,
		auth: chatAuthorization{
			userScopeSets: chatReadScopes(ScopeModeratorManageChatMessages),
			subjectID:     req.ModeratorID,
		},
	})
}
