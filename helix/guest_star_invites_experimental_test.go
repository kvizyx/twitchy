package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestExperimentalGuestStarInviteOperations(t *testing.T) {
	// Given a user token with the moderator manage scope alternative.
	credential := guestStarCredential(helix.ScopeModeratorManageGuestStar)
	tests := []struct {
		name    string
		call    func(*helix.Client) (helix.ResponseMeta, error)
		fixture testkit.ContractFixture
	}{
		{
			name: "get invites",
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				response, err := client.Experimental.GuestStar.GetGuestStarInvites(context.Background(), helix.GetGuestStarInvitesRequest{BroadcasterID: "1234", ModeratorID: "5678", SessionID: "session-1"})
				return guestStarMeta(response, err)
			},
			fixture: guestStarFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "session_id", "session-1"), "", http.StatusOK, guestStarBody(t, "invites.json")),
		},
		{
			name: "send invite",
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				response, err := client.Experimental.GuestStar.SendGuestStarInvite(context.Background(), helix.SendGuestStarInviteRequest{BroadcasterID: "1234", ModeratorID: "5678", SessionID: "session-1", GuestID: "9012"})
				return guestStarMeta(response, err)
			},
			fixture: guestStarFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "session_id", "session-1", "guest_id", "9012"), "", http.StatusNoContent, ""),
		},
		{
			name: "delete invite",
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				response, err := client.Experimental.GuestStar.DeleteGuestStarInvite(context.Background(), helix.DeleteGuestStarInviteRequest{BroadcasterID: "1234", ModeratorID: "5678", SessionID: "session-1", GuestID: "9012"})
				return guestStarMeta(response, err)
			},
			fixture: guestStarFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "session_id", "session-1", "guest_id", "9012"), "", http.StatusNoContent, ""),
		},
	}
	anchors := map[string]string{"get invites": "get-guest-star-invites", "send invite": "send-guest-star-invite", "delete invite": "delete-guest-star-invite"}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(testCase.fixture.Response))
			client, err := guestStarClient(t, transport, credential)
			if err != nil {
				t.Fatal(err)
			}
			meta, callErr := testCase.call(client)
			if err := runGuestStarContract(t, testCase.fixture, anchors[testCase.name], transport, meta, callErr); err != nil {
				t.Fatal(err)
			}
		})
	}
}
