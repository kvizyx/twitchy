package helix

import "context"

type GetUserEmotesRequest struct {
	UserID        string  `query:"user_id"`
	After         *string `query:"after,omitempty"`
	BroadcasterID *string `query:"broadcaster_id,omitempty"`
}

type ChatUserEmote struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	EmoteType  string   `json:"emote_type"`
	EmoteSetID string   `json:"emote_set_id"`
	OwnerID    string   `json:"owner_id"`
	Format     []string `json:"format"`
	Scale      []string `json:"scale"`
	ThemeMode  []string `json:"theme_mode"`
}

type GetUserEmotesData []ChatUserEmote

type SendChatAnnouncementRequest struct {
	BroadcasterID string `query:"broadcaster_id" json:"-"`
	ModeratorID   string `query:"moderator_id" json:"-"`
	Message       string `query:"-" json:"message"`
	Color         string `query:"-" json:"color,omitempty"`
	ForSourceOnly *bool  `query:"-" json:"for_source_only,omitempty"`
}

type SendChatAnnouncementData struct{}

type SendShoutoutRequest struct {
	FromBroadcasterID string `query:"from_broadcaster_id"`
	ToBroadcasterID   string `query:"to_broadcaster_id"`
	ModeratorID       string `query:"moderator_id"`
}

type SendShoutoutData struct{}

type SendChatMessageRequest struct {
	BroadcasterID        string `query:"-" json:"broadcaster_id"`
	SenderID             string `query:"-" json:"sender_id"`
	Message              string `query:"-" json:"message"`
	ReplyParentMessageID string `query:"-" json:"reply_parent_message_id,omitempty"`
	ForSourceOnly        *bool  `query:"-" json:"for_source_only,omitempty"`
	Pin                  *bool  `query:"-" json:"pin,omitempty"`
}

type ChatMessageDropReason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SendChatMessageResult struct {
	MessageID  string                 `json:"message_id"`
	IsSent     bool                   `json:"is_sent"`
	DropReason *ChatMessageDropReason `json:"drop_reason"`
}

type SendChatMessageData []SendChatMessageResult

func (s *ChatService) GetUserEmotes(ctx context.Context, req GetUserEmotesRequest) (*Response[GetUserEmotesData], error) {
	return executeChatEndpoint[GetUserEmotesData](chatEndpointSpec{
		client: s.client,
		ctx:    ctx,
		anchor: "get-user-emotes",
		query:  req,
		auth: chatAuthorization{
			userScopeSets: chatReadScopes(ScopeUserReadEmotes),
			subjectID:     req.UserID,
		},
	})
}

func (s *ChatService) GetUserEmotesPager(req GetUserEmotesRequest, opts ...PagerOption) (*Pager[GetUserEmotesData], error) {
	initialCursor := ""
	if req.After != nil {
		initialCursor = *req.After
	}
	request := req
	request.After = nil
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetUserEmotesData], error) {
		pageRequest := request
		if cursor != "" {
			pageRequest.After = &cursor
		}
		return s.GetUserEmotes(ctx, pageRequest)
	}, initialCursor, opts...)
}

func (s *ChatService) SendChatAnnouncement(ctx context.Context, req SendChatAnnouncementRequest) (*Response[SendChatAnnouncementData], error) {
	return executeChatEndpoint[SendChatAnnouncementData](chatEndpointSpec{
		client: s.client,
		ctx:    ctx,
		anchor: "send-chat-announcement",
		query:  req,
		body:   req,
		auth: chatAuthorization{
			// App access tokens carry no scopes: Twitch authorizes them through the
			// moderator's moderator:manage:announcements + user:bot grants and the
			// broadcaster's channel:bot grant, so there is nothing to check locally.
			userScopeSets:           chatReadScopes(ScopeModeratorManageAnnouncements),
			subjectID:               req.ModeratorID,
			rejectForSourceOnlyUser: req.ForSourceOnly != nil,
		},
	})
}

func (s *ChatService) SendShoutout(ctx context.Context, req SendShoutoutRequest) (*Response[SendShoutoutData], error) {
	return executeChatEndpoint[SendShoutoutData](chatEndpointSpec{
		client: s.client,
		ctx:    ctx,
		anchor: "send-a-shoutout",
		query:  req,
		auth: chatAuthorization{
			userScopeSets: chatReadScopes(ScopeModeratorManageShoutouts),
			subjectID:     req.ModeratorID,
		},
	})
}

func (s *ChatService) SendChatMessage(ctx context.Context, req SendChatMessageRequest) (*Response[SendChatMessageData], error) {
	userScopeSets := chatReadScopes(ScopeUserWriteChat)
	if req.Pin != nil && *req.Pin {
		userScopeSets = [][]AuthorizationScope{{ScopeUserWriteChat, ScopeModeratorManageChatMessages}}
	}
	return executeChatEndpoint[SendChatMessageData](chatEndpointSpec{
		client: s.client,
		ctx:    ctx,
		anchor: "send-chat-message",
		body:   req,
		auth: chatAuthorization{
			// App access tokens carry no scopes: Twitch authorizes them through the
			// sender's user:write:chat + user:bot grants and the broadcaster's
			// channel:bot grant, so there is nothing to check locally.
			userScopeSets:           userScopeSets,
			subjectID:               req.SenderID,
			rejectForSourceOnlyUser: req.ForSourceOnly != nil,
		},
	})
}
