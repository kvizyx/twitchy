package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestExperimentalGuestStarSlotOperations(t *testing.T) {
	// Given a user token with the channel manage alternative and optional slot settings.
	credential := guestStarCredential(helix.ScopeChannelManageGuestStar)
	audioEnabled := false
	videoEnabled := true
	isLive := true
	volume := 75
	reinvite := "true"
	tests := []struct {
		name    string
		call    func(*helix.Client) (helix.ResponseMeta, error)
		fixture testkit.ContractFixture
	}{
		{
			name: "assign slot",
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				response, err := client.Experimental.GuestStar.AssignGuestStarSlot(context.Background(), helix.AssignGuestStarSlotRequest{BroadcasterID: "1234", ModeratorID: "5678", SessionID: "session-1", GuestID: "9012", SlotID: "1"})
				return guestStarMeta(response, err)
			},
			fixture: guestStarFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "session_id", "session-1", "guest_id", "9012", "slot_id", "1"), "", http.StatusNoContent, ""),
		},
		{
			name: "update slot",
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				response, err := client.Experimental.GuestStar.UpdateGuestStarSlot(context.Background(), helix.UpdateGuestStarSlotRequest{BroadcasterID: "1234", ModeratorID: "5678", SessionID: "session-1", SourceSlotID: "1", DestinationSlotID: guestStarString("2")})
				return guestStarMeta(response, err)
			},
			fixture: guestStarFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "session_id", "session-1", "source_slot_id", "1", "destination_slot_id", "2"), "", http.StatusNoContent, ""),
		},
		{
			name: "delete slot",
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				response, err := client.Experimental.GuestStar.DeleteGuestStarSlot(context.Background(), helix.DeleteGuestStarSlotRequest{BroadcasterID: "1234", ModeratorID: "5678", SessionID: "session-1", GuestID: "9012", SlotID: "1", ShouldReinviteGuest: &reinvite})
				return guestStarMeta(response, err)
			},
			fixture: guestStarFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "session_id", "session-1", "guest_id", "9012", "slot_id", "1", "should_reinvite_guest", "true"), "", http.StatusNoContent, ""),
		},
		{
			name: "update slot settings",
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				response, err := client.Experimental.GuestStar.UpdateGuestStarSlotSettings(context.Background(), helix.UpdateGuestStarSlotSettingsRequest{BroadcasterID: "1234", ModeratorID: "5678", SessionID: "session-1", SlotID: "1", IsAudioEnabled: &audioEnabled, IsVideoEnabled: &videoEnabled, IsLive: &isLive, Volume: &volume})
				return guestStarMeta(response, err)
			},
			fixture: guestStarFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "session_id", "session-1", "slot_id", "1", "is_audio_enabled", "false", "is_video_enabled", "true", "is_live", "true", "volume", "75"), "", http.StatusNoContent, ""),
		},
	}
	anchors := map[string]string{"assign slot": "assign-guest-star-slot", "update slot": "update-guest-star-slot", "delete slot": "delete-guest-star-slot", "update slot settings": "update-guest-star-slot-settings"}
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
