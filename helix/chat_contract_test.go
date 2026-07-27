package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestChatGetChatters(t *testing.T) {
	// Given a user credential authorized as the requested moderator.
	fixture := testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query: urlValues("broadcaster_id", "1234", "moderator_id", "5678", "first", "2", "after", "cursor-a"),
			Headers: http.Header{
				"Authorization": {"Bearer user-token"},
				"Client-Id":     {"client-id"},
			},
		},
		Response: testkit.ContractResponse{
			Status: http.StatusOK,
			Headers: http.Header{
				"Ratelimit-Limit":     {"8000"},
				"Ratelimit-Remaining": {"7999"},
				"Ratelimit-Reset":     {verticalSliceRateReset},
			},
			Body:    `{"data":[{"user_id":"1","user_login":"viewer","user_name":"Viewer"}],"pagination":{"cursor":"cursor-b"},"total":1}`,
			Success: true,
		},
		Want: testkit.ContractExpectation{Attempts: 1, RateLimitValid: true},
	}
	fixture.Want.RateLimit.Limit = 8000
	fixture.Want.RateLimit.Remaining = 7999
	transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(fixture.Response))
	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(helix.Credential{
			AccessToken: "user-token",
			ClientID:    "client-id",
			TokenClass:  helix.TokenClassUser,
			UserID:      "5678",
			Scopes:      []helix.AuthorizationScope{helix.ScopeModeratorReadChatters},
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	// When Chatters is requested with every documented query field.
	result, callErr := client.Chat.GetChatters(context.Background(), helix.GetChattersRequest{
		BroadcasterID: "1234",
		ModeratorID:   "5678",
		First:         chatInt(2),
		After:         chatString("cursor-a"),
	})

	// Then the exact request and typed response are returned.
	if callErr != nil {
		t.Fatal(callErr)
	}
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, "get-chatters"), fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return result.Meta, nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 || result.Data[0].UserID != "1" || result.Pagination.Cursor() != "cursor-b" {
		t.Fatalf("decoded chatters = %#v", result.Data)
	}
}

func TestChatManifestRowsMapOneToOne(t *testing.T) {
	// Given the frozen descriptor rows for all Chat operations.
	tests := []struct {
		anchor, selector, method, requestType, dataType string
	}{
		{"get-chatters", "Client.Chat", "GetChatters", "GetChattersRequest", "GetChattersData"},
		{"get-channel-emotes", "Client.Chat", "GetChannelEmotes", "GetChannelEmotesRequest", "GetChannelEmotesData"},
		{"get-global-emotes", "Client.Chat", "GetGlobalEmotes", "GetGlobalEmotesRequest", "GetGlobalEmotesData"},
		{"get-emote-sets", "Client.Chat", "GetEmoteSets", "GetEmoteSetsRequest", "GetEmoteSetsData"},
		{"get-channel-chat-badges", "Client.Chat", "GetChannelChatBadges", "GetChannelChatBadgesRequest", "GetChannelChatBadgesData"},
		{"get-global-chat-badges", "Client.Chat", "GetGlobalChatBadges", "GetGlobalChatBadgesRequest", "GetGlobalChatBadgesData"},
		{"get-chat-settings", "Client.Chat", "GetChatSettings", "GetChatSettingsRequest", "GetChatSettingsData"},
		{"get-shared-chat-session", "Client.Chat", "GetSharedChatSession", "GetSharedChatSessionRequest", "GetSharedChatSessionData"},
		{"get-user-emotes", "Client.Chat", "GetUserEmotes", "GetUserEmotesRequest", "GetUserEmotesData"},
		{"update-chat-settings", "Client.Chat", "UpdateChatSettings", "UpdateChatSettingsRequest", "UpdateChatSettingsData"},
		{"send-chat-announcement", "Client.Chat", "SendChatAnnouncement", "SendChatAnnouncementRequest", "SendChatAnnouncementData"},
		{"send-a-shoutout", "Client.Chat", "SendShoutout", "SendShoutoutRequest", "SendShoutoutData"},
		{"send-chat-message", "Client.Chat", "SendChatMessage", "SendChatMessageRequest", "SendChatMessageData"},
		{"get-pinned-chat-message", "Client.Experimental.Chat", "GetPinnedChatMessage", "GetPinnedChatMessageRequest", "GetPinnedChatMessageData"},
		{"pin-chat-message", "Client.Experimental.Chat", "PinChatMessage", "PinChatMessageRequest", "PinChatMessageData"},
		{"update-pinned-chat-message", "Client.Experimental.Chat", "UpdatePinnedChatMessage", "UpdatePinnedChatMessageRequest", "UpdatePinnedChatMessageData"},
		{"unpin-chat-message", "Client.Experimental.Chat", "UnpinChatMessage", "UnpinChatMessageRequest", "UnpinChatMessageData"},
		{"get-user-chat-color", "Client.Chat", "GetUserChatColor", "GetUserChatColorRequest", "GetUserChatColorData"},
		{"update-user-chat-color", "Client.Chat", "UpdateUserChatColor", "UpdateUserChatColorRequest", "UpdateUserChatColorData"},
	}

	// When each row is resolved from the manifest.
	for _, testCase := range tests {
		t.Run(testCase.anchor, func(t *testing.T) {
			operation := manifestOperation(t, testCase.anchor)
			implementation := operation.Implementation
			if implementation.Selector != testCase.selector || implementation.Method != testCase.method || implementation.RequestType != testCase.requestType || implementation.DataType != testCase.dataType {
				t.Fatalf("implementation = %#v", implementation)
			}
		})
	}
}
