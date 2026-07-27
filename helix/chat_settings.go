package helix

import "context"

type GetChatSettingsRequest struct {
	BroadcasterID string  `query:"broadcaster_id"`
	ModeratorID   *string `query:"moderator_id,omitempty"`
}

type ChatSettings struct {
	BroadcasterID                 string `json:"broadcaster_id"`
	ModeratorID                   string `json:"moderator_id"`
	EmoteMode                     bool   `json:"emote_mode"`
	FollowerMode                  bool   `json:"follower_mode"`
	FollowerModeDuration          *int   `json:"follower_mode_duration"`
	NonModeratorChatDelay         bool   `json:"non_moderator_chat_delay"`
	NonModeratorChatDelayDuration *int   `json:"non_moderator_chat_delay_duration"`
	SlowMode                      bool   `json:"slow_mode"`
	SlowModeWaitTime              *int   `json:"slow_mode_wait_time"`
	SubscriberMode                bool   `json:"subscriber_mode"`
	UniqueChatMode                bool   `json:"unique_chat_mode"`
}

type GetChatSettingsData []ChatSettings

type UpdateChatSettingsRequest struct {
	BroadcasterID                 string `query:"broadcaster_id" json:"-"`
	ModeratorID                   string `query:"moderator_id" json:"-"`
	EmoteMode                     *bool  `query:"-" json:"emote_mode,omitempty"`
	FollowerMode                  *bool  `query:"-" json:"follower_mode,omitempty"`
	FollowerModeDuration          *int   `query:"-" json:"follower_mode_duration,omitempty"`
	NonModeratorChatDelay         *bool  `query:"-" json:"non_moderator_chat_delay,omitempty"`
	NonModeratorChatDelayDuration *int   `query:"-" json:"non_moderator_chat_delay_duration,omitempty"`
	SlowMode                      *bool  `query:"-" json:"slow_mode,omitempty"`
	SlowModeWaitTime              *int   `query:"-" json:"slow_mode_wait_time,omitempty"`
	SubscriberMode                *bool  `query:"-" json:"subscriber_mode,omitempty"`
	UniqueChatMode                *bool  `query:"-" json:"unique_chat_mode,omitempty"`
}

type UpdateChatSettingsData []ChatSettings

type GetSharedChatSessionRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type SharedChatSession struct {
	SessionID         string                  `json:"session_id"`
	HostBroadcasterID string                  `json:"host_broadcaster_id"`
	Participants      []SharedChatParticipant `json:"participants"`
	CreatedAt         Timestamp               `json:"created_at"`
	UpdatedAt         Timestamp               `json:"updated_at"`
}

type SharedChatParticipant struct {
	BroadcasterID string `json:"broadcaster_id"`
}

type GetSharedChatSessionData []SharedChatSession

func (s *ChatService) GetChatSettings(ctx context.Context, req GetChatSettingsRequest) (*Response[GetChatSettingsData], error) {
	subjectID := ""
	if req.ModeratorID != nil {
		subjectID = *req.ModeratorID
	}
	return executeChatEndpoint[GetChatSettingsData](chatEndpointSpec{
		client: s.client,
		ctx:    ctx,
		anchor: "get-chat-settings",
		query:  req,
		auth:   chatAuthorization{subjectID: subjectID},
	})
}

func (s *ChatService) UpdateChatSettings(ctx context.Context, req UpdateChatSettingsRequest) (*Response[UpdateChatSettingsData], error) {
	return executeChatEndpoint[UpdateChatSettingsData](chatEndpointSpec{
		client: s.client,
		ctx:    ctx,
		anchor: "update-chat-settings",
		query:  req,
		body:   req,
		auth: chatAuthorization{
			userScopeSets: chatReadScopes(ScopeModeratorManageChatSettings),
			appScopeSets:  chatReadScopes(ScopeModeratorManageChatSettings),
			subjectID:     req.ModeratorID,
		},
	})
}

func (s *ChatService) GetSharedChatSession(ctx context.Context, req GetSharedChatSessionRequest) (*Response[GetSharedChatSessionData], error) {
	return executeChatEndpoint[GetSharedChatSessionData](chatEndpointSpec{client: s.client, ctx: ctx, anchor: "get-shared-chat-session", query: req})
}
