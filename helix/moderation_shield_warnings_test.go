package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
)

func TestModerationShieldAndWarningOperations(t *testing.T) {
	tests := []moderationOperationCase{
		{
			anchor:     "update-shield-mode-status",
			fixture:    moderationBodyFixture(moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusOK, `{"data":[{"is_active":false,"moderator_id":"5678","moderator_login":"mod","moderator_name":"Mod","last_activated_at":"2026-01-02T03:04:05Z"}]}`), `{"is_active":false}`),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageShieldMode),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.UpdateShieldModeStatus(context.Background(), helix.UpdateShieldModeStatusRequest{BroadcasterID: "1234", ModeratorID: "5678", IsActive: false}))
			},
		},
		{
			anchor:     "get-shield-mode-status",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusOK, `{"data":[{"is_active":true,"moderator_id":"5678","moderator_login":"mod","moderator_name":"Mod","last_activated_at":"2026-01-02T03:04:05Z"}]}`),
			credential: moderationCredential(helix.TokenClassApp, "", helix.ScopeModeratorReadShieldMode),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.GetShieldModeStatus(context.Background(), helix.GetShieldModeStatusRequest{BroadcasterID: "1234", ModeratorID: "5678"}))
			},
		},
		{
			anchor:     "warn-chat-user",
			fixture:    moderationBodyFixture(moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusOK, `{"data":[{"broadcaster_id":"1234","user_id":"9876","moderator_id":"5678","reason":"stop"}]}`), `{"data":{"user_id":"9876","reason":"stop"}}`),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageWarnings),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.WarnChatUser(context.Background(), helix.WarnChatUserRequest{BroadcasterID: "1234", ModeratorID: "5678", Data: helix.WarnChatUserBody{UserID: "9876", Reason: "stop"}}))
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.anchor, func(t *testing.T) { runModerationOperation(t, testCase) })
	}
}
