package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestExperimentalGuestStarSettingsOperations(t *testing.T) {
	// Given a user token with a channel Guest Star scope and exact settings wire data.
	enabled := true
	slotCount := 4
	regenerate := true
	credential := guestStarCredential(helix.ScopeChannelManageGuestStar)
	settingsBody := `{"is_moderator_send_live_enabled":true,"slot_count":4,"is_browser_source_audio_enabled":true,"group_layout":"SCREENSHARE_LAYOUT","regenerate_browser_sources":true}`
	tests := []struct {
		name    string
		call    func(*helix.Client) (helix.ResponseMeta, error)
		fixture testkit.ContractFixture
	}{
		{
			name: "get channel settings",
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				response, err := client.Experimental.GuestStar.GetChannelGuestStarSettings(context.Background(), helix.GetChannelGuestStarSettingsRequest{BroadcasterID: "1234", ModeratorID: "5678"})
				return guestStarMeta(response, err)
			},
			fixture: guestStarFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), "", http.StatusOK, guestStarBody(t, "settings.json")),
		},
		{
			name: "update channel settings",
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				response, err := client.Experimental.GuestStar.UpdateChannelGuestStarSettings(context.Background(), helix.UpdateChannelGuestStarSettingsRequest{
					BroadcasterID:               "5678",
					IsModeratorSendLiveEnabled:  &enabled,
					SlotCount:                   &slotCount,
					IsBrowserSourceAudioEnabled: &enabled,
					GroupLayout:                 guestStarGroupLayoutPointer(helix.GuestStarGroupLayoutScreenshare),
					RegenerateBrowserSources:    &regenerate,
				})
				return guestStarMeta(response, err)
			},
			fixture: guestStarFixture(urlValues("broadcaster_id", "5678"), settingsBody, http.StatusNoContent, ""),
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(testCase.fixture.Response))
			client, err := guestStarClient(t, transport, credential)
			if err != nil {
				t.Fatal(err)
			}
			meta, callErr := testCase.call(client)
			if err := runGuestStarContract(t, testCase.fixture, map[string]string{"get channel settings": "get-channel-guest-star-settings", "update channel settings": "update-channel-guest-star-settings"}[testCase.name], transport, meta, callErr); err != nil {
				t.Fatal(err)
			}
		})
	}
}
