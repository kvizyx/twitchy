package helix

import "context"

type GetChannelEmotesRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type GetGlobalEmotesRequest struct{}

type GetEmoteSetsRequest struct {
	EmoteSetIDs []string `query:"emote_set_id"`
}

type ChatEmoteImages struct {
	URL1x string `json:"url_1x"`
	URL2x string `json:"url_2x"`
	URL4x string `json:"url_4x"`
}

type ChatEmote struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Images     ChatEmoteImages `json:"images"`
	Tier       string          `json:"tier"`
	EmoteType  string          `json:"emote_type"`
	EmoteSetID string          `json:"emote_set_id"`
	OwnerID    string          `json:"owner_id"`
	Format     []string        `json:"format"`
	Scale      []string        `json:"scale"`
	ThemeMode  []string        `json:"theme_mode"`
}

type GetChannelEmotesData []ChatEmote
type GetGlobalEmotesData []ChatEmote
type GetEmoteSetsData []ChatEmote

type GetChannelChatBadgesRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type GetGlobalChatBadgesRequest struct{}

type ChatBadgeSet struct {
	SetID    string             `json:"set_id"`
	Versions []ChatBadgeVersion `json:"versions"`
}

type ChatBadgeVersion struct {
	ID          string  `json:"id"`
	ImageURL1x  string  `json:"image_url_1x"`
	ImageURL2x  string  `json:"image_url_2x"`
	ImageURL4x  string  `json:"image_url_4x"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ClickAction *string `json:"click_action"`
	ClickURL    *string `json:"click_url"`
}

type GetChannelChatBadgesData []ChatBadgeSet
type GetGlobalChatBadgesData []ChatBadgeSet

func (s *ChatService) GetChannelEmotes(ctx context.Context, req GetChannelEmotesRequest) (*Response[GetChannelEmotesData], error) {
	return executeChatEndpoint[GetChannelEmotesData](chatEndpointSpec{client: s.client, ctx: ctx, anchor: "get-channel-emotes", query: req})
}

func (s *ChatService) GetGlobalEmotes(ctx context.Context, req GetGlobalEmotesRequest) (*Response[GetGlobalEmotesData], error) {
	return executeChatEndpoint[GetGlobalEmotesData](chatEndpointSpec{client: s.client, ctx: ctx, anchor: "get-global-emotes", query: req})
}

func (s *ChatService) GetEmoteSets(ctx context.Context, req GetEmoteSetsRequest) (*Response[GetEmoteSetsData], error) {
	return executeChatEndpoint[GetEmoteSetsData](chatEndpointSpec{client: s.client, ctx: ctx, anchor: "get-emote-sets", query: req})
}

func (s *ChatService) GetChannelChatBadges(ctx context.Context, req GetChannelChatBadgesRequest) (*Response[GetChannelChatBadgesData], error) {
	return executeChatEndpoint[GetChannelChatBadgesData](chatEndpointSpec{client: s.client, ctx: ctx, anchor: "get-channel-chat-badges", query: req})
}

func (s *ChatService) GetGlobalChatBadges(ctx context.Context, req GetGlobalChatBadgesRequest) (*Response[GetGlobalChatBadgesData], error) {
	return executeChatEndpoint[GetGlobalChatBadgesData](chatEndpointSpec{client: s.client, ctx: ctx, anchor: "get-global-chat-badges", query: req})
}
