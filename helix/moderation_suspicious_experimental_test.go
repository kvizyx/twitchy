package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
)

func TestModerationSuspiciousExperimentalOperations(t *testing.T) {
	tests := []moderationOperationCase{
		{
			anchor:     "add-suspicious-status-to-chat-user",
			fixture:    moderationBodyFixture(moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusOK, `{"data":[{"user_id":"9876","broadcaster_id":"1234","moderator_id":"5678","updated_at":"2026-01-02T03:04:05Z","status":"RESTRICTED","types":["MANUALLY_ADDED"]}]}`), `{"user_id":"9876","status":"RESTRICTED"}`),
			credential: moderationCredential(helix.TokenClassApp, "", helix.ScopeModeratorManageSuspiciousUsers),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Experimental.Moderation.AddSuspiciousStatusToChatUser(context.Background(), helix.AddSuspiciousStatusToChatUserRequest{BroadcasterID: "1234", ModeratorID: "5678", UserID: "9876", Status: "RESTRICTED"}))
			},
		},
		{
			anchor:     "remove-suspicious-status-from-chat-user",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "user_id", "9876"), http.StatusOK, `{"data":[{"user_id":"9876","broadcaster_id":"1234","moderator_id":"5678","updated_at":"2026-01-02T03:04:05Z","status":"NO_TREATMENT","types":["MANUALLY_ADDED"]}]}`),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageSuspiciousUsers),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Experimental.Moderation.RemoveSuspiciousStatusFromChatUser(context.Background(), helix.RemoveSuspiciousStatusFromChatUserRequest{BroadcasterID: "1234", ModeratorID: "5678", UserID: "9876"}))
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.anchor, func(t *testing.T) { runModerationOperation(t, testCase) })
	}
}
