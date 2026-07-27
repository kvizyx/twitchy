package helix

import "context"

type SendExtensionPubSubMessageRequest struct {
	Target            []string `json:"target"`
	BroadcasterID     string   `json:"broadcaster_id"`
	IsGlobalBroadcast *bool    `json:"is_global_broadcast,omitempty"`
	Message           string   `json:"message"`
}

type SendExtensionPubSubMessageData struct{}

type SendExtensionChatMessageRequest struct {
	BroadcasterID    string `query:"broadcaster_id" json:"-"`
	Text             string `json:"text"`
	ExtensionID      string `json:"extension_id"`
	ExtensionVersion string `json:"extension_version"`
}

type SendExtensionChatMessageData struct{}

func (s *ExtensionsService) SendExtensionPubSubMessage(ctx context.Context, req SendExtensionPubSubMessageRequest) (*Response[SendExtensionPubSubMessageData], error) {
	return executeExtensionEndpoint[SendExtensionPubSubMessageData](s.client, ctx, "send-extension-pubsub-message", nil, req)
}

func (s *ExtensionsService) SendExtensionChatMessage(ctx context.Context, req SendExtensionChatMessageRequest) (*Response[SendExtensionChatMessageData], error) {
	query := struct {
		BroadcasterID string `query:"broadcaster_id"`
	}{BroadcasterID: req.BroadcasterID}
	return executeExtensionEndpoint[SendExtensionChatMessageData](s.client, ctx, "send-extension-chat-message", query, req)
}
