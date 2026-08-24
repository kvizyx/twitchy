package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
)

func TestChatStableOperations(t *testing.T) {
	// Given one exact wire fixture for each stable Chat operation.
	falseValue := false
	tests := []chatOperationCase{
		{
			anchor:     "get-chatters",
			fixture:    chatSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "first", "2", "after", "cursor-a"), http.StatusOK, `{"data":[{"user_id":"1","user_login":"viewer","user_name":"Viewer"}],"pagination":{"cursor":"cursor-b"}}`),
			credential: chatCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorReadChatters),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.GetChatters(context.Background(), helix.GetChattersRequest{BroadcasterID: "1234", ModeratorID: "5678", First: chatInt(2), After: chatString("cursor-a")}))
			},
		},
		{
			anchor:     "get-channel-emotes",
			fixture:    chatSuccessFixture(urlValues("broadcaster_id", "1234"), http.StatusOK, `{"data":[{"id":"e1","name":"Wave","images":{"url_1x":"1","url_2x":"2","url_4x":"4"}}]}`),
			credential: chatCredential(helix.TokenClassApp, ""),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.GetChannelEmotes(context.Background(), helix.GetChannelEmotesRequest{BroadcasterID: "1234"}))
			},
		},
		{
			anchor:     "get-global-emotes",
			fixture:    chatSuccessFixture(urlValues(), http.StatusOK, `{"data":[]}`),
			credential: chatCredential(helix.TokenClassApp, ""),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.GetGlobalEmotes(context.Background(), helix.GetGlobalEmotesRequest{}))
			},
		},
		{
			anchor:     "get-emote-sets",
			fixture:    chatSuccessFixture(urlValues("emote_set_id", "set-a", "emote_set_id", "set-b"), http.StatusOK, `{"data":[]}`),
			credential: chatCredential(helix.TokenClassApp, ""),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.GetEmoteSets(context.Background(), helix.GetEmoteSetsRequest{EmoteSetIDs: []string{"set-a", "set-b"}}))
			},
		},
		{
			anchor:     "get-channel-chat-badges",
			fixture:    chatSuccessFixture(urlValues("broadcaster_id", "1234"), http.StatusOK, `{"data":[]}`),
			credential: chatCredential(helix.TokenClassApp, ""),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.GetChannelChatBadges(context.Background(), helix.GetChannelChatBadgesRequest{BroadcasterID: "1234"}))
			},
		},
		{
			anchor:     "get-global-chat-badges",
			fixture:    chatSuccessFixture(urlValues(), http.StatusOK, `{"data":[]}`),
			credential: chatCredential(helix.TokenClassApp, ""),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.GetGlobalChatBadges(context.Background(), helix.GetGlobalChatBadgesRequest{}))
			},
		},
		{
			anchor:     "get-chat-settings",
			fixture:    chatSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusOK, `{"data":[]}`),
			credential: chatCredential(helix.TokenClassApp, ""),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.GetChatSettings(context.Background(), helix.GetChatSettingsRequest{BroadcasterID: "1234", ModeratorID: chatString("5678")}))
			},
		},
		{
			anchor:     "get-shared-chat-session",
			fixture:    chatSuccessFixture(urlValues("broadcaster_id", "1234"), http.StatusOK, `{"data":[]}`),
			credential: chatCredential(helix.TokenClassApp, ""),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.GetSharedChatSession(context.Background(), helix.GetSharedChatSessionRequest{BroadcasterID: "1234"}))
			},
		},
		{
			anchor:     "get-user-emotes",
			fixture:    chatSuccessFixture(urlValues("user_id", "5678", "after", "cursor-a", "broadcaster_id", "1234"), http.StatusOK, `{"data":[],"pagination":{"cursor":"cursor-b"}}`),
			credential: chatCredential(helix.TokenClassUser, "5678", helix.ScopeUserReadEmotes),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.GetUserEmotes(context.Background(), helix.GetUserEmotesRequest{UserID: "5678", After: chatString("cursor-a"), BroadcasterID: chatString("1234")}))
			},
		},
		{
			anchor:     "update-chat-settings",
			fixture:    chatRequestBodyFixture(chatSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusOK, `{"data":[]}`), `{"follower_mode":false}`),
			credential: chatCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageChatSettings),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.UpdateChatSettings(context.Background(), helix.UpdateChatSettingsRequest{BroadcasterID: "1234", ModeratorID: "5678", FollowerMode: &falseValue}))
			},
		},
		{
			anchor:     "send-chat-announcement",
			fixture:    chatRequestBodyFixture(chatSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusNoContent, ""), `{"message":"hello","color":"purple"}`),
			credential: chatCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageAnnouncements),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.SendChatAnnouncement(context.Background(), helix.SendChatAnnouncementRequest{BroadcasterID: "1234", ModeratorID: "5678", Message: "hello", Color: "purple"}))
			},
		},
		{
			anchor:     "send-a-shoutout",
			fixture:    chatSuccessFixture(urlValues("from_broadcaster_id", "1234", "to_broadcaster_id", "9876", "moderator_id", "5678"), http.StatusNoContent, ""),
			credential: chatCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageShoutouts),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.SendShoutout(context.Background(), helix.SendShoutoutRequest{FromBroadcasterID: "1234", ToBroadcasterID: "9876", ModeratorID: "5678"}))
			},
		},
		{
			anchor:     "send-chat-message",
			fixture:    chatRequestBodyFixture(chatSuccessFixture(urlValues(), http.StatusOK, `{"data":[{"message_id":"message-1","is_sent":true}]}`), `{"broadcaster_id":"1234","sender_id":"5678","message":"hello","for_source_only":false}`),
			credential: chatCredential(helix.TokenClassApp, "", helix.ScopeUserWriteChat, helix.ScopeUserBot, helix.ScopeChannelBot),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.SendChatMessage(context.Background(), helix.SendChatMessageRequest{BroadcasterID: "1234", SenderID: "5678", Message: "hello", ForSourceOnly: &falseValue}))
			},
		},
		{
			anchor:     "get-user-chat-color",
			fixture:    chatSuccessFixture(urlValues("user_id", "1111", "user_id", "2222"), http.StatusOK, `{"data":[]}`),
			credential: chatCredential(helix.TokenClassApp, ""),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.GetUserChatColor(context.Background(), helix.GetUserChatColorRequest{UserIDs: []string{"1111", "2222"}}))
			},
		},
		{
			anchor:     "update-user-chat-color",
			fixture:    chatSuccessFixture(urlValues("user_id", "5678", "color", "blue"), http.StatusNoContent, ""),
			credential: chatCredential(helix.TokenClassUser, "5678", helix.ScopeUserManageChatColor),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return chatMeta(client.Chat.UpdateUserChatColor(context.Background(), helix.UpdateUserChatColorRequest{UserID: "5678", Color: "blue"}))
			},
		},
	}

	// When each stable service method is called.
	for _, testCase := range tests {
		t.Run(testCase.anchor, func(t *testing.T) {
			runChatOperation(t, testCase)
		})
	}
}

func TestChatSendChatMessageAllowsScopelessAppToken(t *testing.T) {
	// Given a message request authorized by a scope-less app access token.
	testCase := chatOperationCase{
		anchor:     "send-chat-message",
		fixture:    chatRequestBodyFixture(chatSuccessFixture(urlValues(), http.StatusOK, `{"data":[{"message_id":"message-1","is_sent":true}]}`), `{"broadcaster_id":"1234","sender_id":"5678","message":"hello"}`),
		credential: chatCredential(helix.TokenClassApp, ""),
		call: func(client *helix.Client) (helix.ResponseMeta, error) {
			return chatMeta(client.Chat.SendChatMessage(context.Background(), helix.SendChatMessageRequest{BroadcasterID: "1234", SenderID: "5678", Message: "hello"}))
		},
	}

	// When the message is sent.
	runChatOperation(t, testCase)
}

func TestChatSendChatAnnouncementAllowsScopelessAppToken(t *testing.T) {
	// Given an announcement request authorized by a scope-less app access token.
	testCase := chatOperationCase{
		anchor:     "send-chat-announcement",
		fixture:    chatRequestBodyFixture(chatSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusNoContent, ""), `{"message":"hello"}`),
		credential: chatCredential(helix.TokenClassApp, ""),
		call: func(client *helix.Client) (helix.ResponseMeta, error) {
			return chatMeta(client.Chat.SendChatAnnouncement(context.Background(), helix.SendChatAnnouncementRequest{BroadcasterID: "1234", ModeratorID: "5678", Message: "hello"}))
		},
	}

	// When the announcement is sent.
	runChatOperation(t, testCase)
}

func TestChatSendChatMessageOmitsUnsetSourceAndPin(t *testing.T) {
	// Given a message request without the optional source or pin flags.
	testCase := chatOperationCase{
		anchor:     "send-chat-message",
		fixture:    chatRequestBodyFixture(chatSuccessFixture(urlValues(), http.StatusOK, `{"data":[{"message_id":"message-1","is_sent":true}]}`), `{"broadcaster_id":"1234","sender_id":"5678","message":"hello"}`),
		credential: chatCredential(helix.TokenClassUser, "5678", helix.ScopeUserWriteChat),
		call: func(client *helix.Client) (helix.ResponseMeta, error) {
			return chatMeta(client.Chat.SendChatMessage(context.Background(), helix.SendChatMessageRequest{BroadcasterID: "1234", SenderID: "5678", Message: "hello"}))
		},
	}

	// When the message is sent.
	runChatOperation(t, testCase)
}
