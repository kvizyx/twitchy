package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestExperimentalGuestStarSessionOperations(t *testing.T) {
	// Given a user token and the session response fixture.
	credential := guestStarCredential(helix.ScopeChannelManageGuestStar)
	tests := []struct {
		name    string
		call    func(*helix.Client) (helix.ResponseMeta, error)
		fixture testkit.ContractFixture
	}{
		{
			name: "get session",
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				response, err := client.Experimental.GuestStar.GetGuestStarSession(context.Background(), helix.GetGuestStarSessionRequest{BroadcasterID: "1234", ModeratorID: "5678"})
				return guestStarMeta(response, err)
			},
			fixture: guestStarFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), "", http.StatusOK, guestStarBody(t, "session.json")),
		},
		{
			name: "create session",
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				response, err := client.Experimental.GuestStar.CreateGuestStarSession(context.Background(), helix.CreateGuestStarSessionRequest{BroadcasterID: "5678"})
				return guestStarMeta(response, err)
			},
			fixture: guestStarFixture(urlValues("broadcaster_id", "5678"), "", http.StatusOK, guestStarBody(t, "session.json")),
		},
		{
			name: "end session",
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				response, err := client.Experimental.GuestStar.EndGuestStarSession(context.Background(), helix.EndGuestStarSessionRequest{BroadcasterID: "5678", SessionID: "session-1"})
				return guestStarMeta(response, err)
			},
			fixture: guestStarFixture(urlValues("broadcaster_id", "5678", "session_id", "session-1"), "", http.StatusOK, guestStarBody(t, "session.json")),
		},
	}
	anchors := map[string]string{"get session": "get-guest-star-session", "create session": "create-guest-star-session", "end session": "end-guest-star-session"}
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
