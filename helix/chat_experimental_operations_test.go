package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
)

func TestChatExperimentalPinOperations(t *testing.T) {
	// Given app credentials carrying the documented bot and moderator scopes.
	credential := chatCredential(
		helix.TokenClassApp,
		"",
		helix.ScopeModeratorManageChatMessages,
		helix.ScopeModeratorReadChatMessages,
		helix.ScopeUserBot,
		helix.ScopeChannelBot,
	)
	duration := 300
	tests := []chatOperationCase{
		{
			anchor:     "get-pinned-chat-message",
			fixture:    chatSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusOK, `{"data":[]}`),
			credential: credential,
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Experimental.Chat.GetPinnedChatMessage(context.Background(), helix.GetPinnedChatMessageRequest{BroadcasterID: "1234", ModeratorID: "5678"}))
			},
		},
		{
			anchor:     "pin-chat-message",
			fixture:    chatSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "message_id", "message-1", "duration_seconds", "300"), http.StatusNoContent, ""),
			credential: credential,
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Experimental.Chat.PinChatMessage(context.Background(), helix.PinChatMessageRequest{BroadcasterID: "1234", ModeratorID: "5678", MessageID: "message-1", DurationSeconds: &duration}))
			},
		},
		{
			anchor:     "update-pinned-chat-message",
			fixture:    chatSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "message_id", "message-1", "duration_seconds", "600"), http.StatusNoContent, ""),
			credential: credential,
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				updatedDuration := 600
				return chatMeta(client.Experimental.Chat.UpdatePinnedChatMessage(context.Background(), helix.UpdatePinnedChatMessageRequest{BroadcasterID: "1234", ModeratorID: "5678", MessageID: "message-1", DurationSeconds: &updatedDuration}))
			},
		},
		{
			anchor:     "unpin-chat-message",
			fixture:    chatSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "message_id", "message-1"), http.StatusNoContent, ""),
			credential: credential,
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Experimental.Chat.UnpinChatMessage(context.Background(), helix.UnpinChatMessageRequest{BroadcasterID: "1234", ModeratorID: "5678", MessageID: "message-1"}))
			},
		},
	}

	// When each experimental pin method is called.
	for _, testCase := range tests {
		t.Run(testCase.anchor, func(t *testing.T) {
			runChatOperation(t, testCase)
		})
	}
}
