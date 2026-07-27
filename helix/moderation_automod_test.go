package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
)

func TestModerationAutoModOperations(t *testing.T) {
	level := 3
	tests := []moderationOperationCase{
		{
			anchor:     "check-automod-status",
			fixture:    moderationBodyFixture(moderationSuccessFixture(urlValues("broadcaster_id", "1234"), http.StatusOK, `{"data":[{"msg_id":"m1","is_permitted":true}]}`), `{"data":[{"msg_id":"m1","msg_text":"hello"}]}`),
			credential: moderationCredential(helix.TokenClassUser, "1234", helix.ScopeModerationRead),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.CheckAutoModStatus(context.Background(), helix.CheckAutoModStatusRequest{BroadcasterID: "1234", Data: []helix.CheckAutoModStatusMessage{{MsgID: "m1", MsgText: "hello"}}}))
			},
		},
		{
			anchor:     "manage-held-automod-messages",
			fixture:    moderationBodyFixture(moderationSuccessFixture(urlValues(), http.StatusNoContent, ""), `{"user_id":"5678","msg_id":"m1","action":"ALLOW"}`),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageAutoMod),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.ManageHeldAutoModMessages(context.Background(), helix.ManageHeldAutoModMessagesRequest{UserID: "5678", MsgID: "m1", Action: "ALLOW"}))
			},
		},
		{
			anchor:     "get-automod-settings",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusOK, `{"data":[{"broadcaster_id":"1234","moderator_id":"5678","overall_level":null,"aggression":0}]}`),
			credential: moderationCredential(helix.TokenClassApp, "", helix.ScopeModeratorReadAutoModSettings),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.GetAutoModSettings(context.Background(), helix.GetAutoModSettingsRequest{BroadcasterID: "1234", ModeratorID: "5678"}))
			},
		},
		{
			anchor:     "update-automod-settings",
			fixture:    moderationBodyFixture(moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusOK, `{"data":[{"overall_level":3}]}`), `{"overall_level":3}`),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageAutoModSettings),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.UpdateAutoModSettings(context.Background(), helix.UpdateAutoModSettingsRequest{BroadcasterID: "1234", ModeratorID: "5678", OverallLevel: &level}))
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.anchor, func(t *testing.T) { runModerationOperation(t, testCase) })
	}
}
