package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
)

func TestModerationRolesOperations(t *testing.T) {
	tests := []moderationOperationCase{
		{
			anchor:     "get-moderated-channels",
			fixture:    moderationSuccessFixture(urlValues("user_id", "5678", "first", "2", "after", "cursor-a"), http.StatusOK, `{"data":[{"broadcaster_id":"1234","broadcaster_login":"channel","broadcaster_name":"Channel"}],"pagination":{"cursor":"cursor-b"}}`),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeUserReadModeratedChannels),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.GetModeratedChannels(context.Background(), helix.GetModeratedChannelsRequest{UserID: "5678", First: moderationInt(2), After: moderationString("cursor-a")}))
			},
		},
		{
			anchor:     "get-moderators",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "user_id", "5678", "first", "2", "after", "cursor-a"), http.StatusOK, `{"data":[{"user_id":"5678","user_login":"mod","user_name":"Mod"}],"pagination":{"cursor":"cursor-b"}}`),
			credential: moderationCredential(helix.TokenClassUser, "1234", helix.ScopeModerationRead),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.GetModerators(context.Background(), helix.GetModeratorsRequest{BroadcasterID: "1234", UserIDs: []string{"5678"}, First: moderationInt(2), After: moderationString("cursor-a")}))
			},
		},
		{
			anchor:     "add-channel-moderator",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "user_id", "5678"), http.StatusNoContent, ""),
			credential: moderationCredential(helix.TokenClassUser, "1234", helix.ScopeChannelManageModerators),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.AddChannelModerator(context.Background(), helix.AddChannelModeratorRequest{BroadcasterID: "1234", UserID: "5678"}))
			},
		},
		{
			anchor:     "remove-channel-moderator",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "user_id", "5678"), http.StatusNoContent, ""),
			credential: moderationCredential(helix.TokenClassUser, "1234", helix.ScopeChannelManageModerators),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.RemoveChannelModerator(context.Background(), helix.RemoveChannelModeratorRequest{BroadcasterID: "1234", UserID: "5678"}))
			},
		},
		{
			anchor:     "get-vips",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "user_id", "5678", "first", "2", "after", "cursor-a"), http.StatusOK, `{"data":[{"user_id":"5678","user_name":"VIP","user_login":"vip"}],"pagination":{"cursor":"cursor-b"}}`),
			credential: moderationCredential(helix.TokenClassUser, "1234", helix.ScopeChannelReadVIPs),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.GetVIPs(context.Background(), helix.GetVIPsRequest{BroadcasterID: "1234", UserIDs: []string{"5678"}, First: moderationInt(2), After: moderationString("cursor-a")}))
			},
		},
		{
			anchor:     "add-channel-vip",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "user_id", "5678"), http.StatusNoContent, ""),
			credential: moderationCredential(helix.TokenClassUser, "1234", helix.ScopeChannelManageVIPs),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.AddChannelVIP(context.Background(), helix.AddChannelVIPRequest{BroadcasterID: "1234", UserID: "5678"}))
			},
		},
		{
			anchor:     "remove-channel-vip",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "user_id", "5678"), http.StatusNoContent, ""),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeChannelManageVIPs),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.RemoveChannelVIP(context.Background(), helix.RemoveChannelVIPRequest{BroadcasterID: "1234", UserID: "5678"}))
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.anchor, func(t *testing.T) { runModerationOperation(t, testCase) })
	}
}
