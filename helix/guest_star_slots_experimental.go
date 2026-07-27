package helix

import "context"

type AssignGuestStarSlotRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
	SessionID     string `query:"session_id"`
	GuestID       string `query:"guest_id"`
	SlotID        string `query:"slot_id"`
}

type AssignGuestStarSlotData struct {
	Code string `json:"code"`
}

type UpdateGuestStarSlotRequest struct {
	BroadcasterID     string  `query:"broadcaster_id"`
	ModeratorID       string  `query:"moderator_id"`
	SessionID         string  `query:"session_id"`
	SourceSlotID      string  `query:"source_slot_id"`
	DestinationSlotID *string `query:"destination_slot_id,omitempty"`
}

type UpdateGuestStarSlotData struct{}

type DeleteGuestStarSlotRequest struct {
	BroadcasterID       string  `query:"broadcaster_id"`
	ModeratorID         string  `query:"moderator_id"`
	SessionID           string  `query:"session_id"`
	GuestID             string  `query:"guest_id"`
	SlotID              string  `query:"slot_id"`
	ShouldReinviteGuest *string `query:"should_reinvite_guest,omitempty"`
}

type DeleteGuestStarSlotData struct{}

type UpdateGuestStarSlotSettingsRequest struct {
	BroadcasterID  string `query:"broadcaster_id"`
	ModeratorID    string `query:"moderator_id"`
	SessionID      string `query:"session_id"`
	SlotID         string `query:"slot_id"`
	IsAudioEnabled *bool  `query:"is_audio_enabled,omitempty"`
	IsVideoEnabled *bool  `query:"is_video_enabled,omitempty"`
	IsLive         *bool  `query:"is_live,omitempty"`
	Volume         *int   `query:"volume,omitempty"`
}

type UpdateGuestStarSlotSettingsData struct{}

func (s *ExperimentalGuestStarService) AssignGuestStarSlot(ctx context.Context, req AssignGuestStarSlotRequest) (*Response[AssignGuestStarSlotData], error) {
	return executeGuestStarEndpoint[AssignGuestStarSlotData](s.client, ctx, "assign-guest-star-slot", req, nil, guestStarAuthorization{scopeSets: guestStarManageScopes(), subjectID: req.ModeratorID})
}

func (s *ExperimentalGuestStarService) UpdateGuestStarSlot(ctx context.Context, req UpdateGuestStarSlotRequest) (*Response[UpdateGuestStarSlotData], error) {
	return executeGuestStarEndpoint[UpdateGuestStarSlotData](s.client, ctx, "update-guest-star-slot", req, nil, guestStarAuthorization{scopeSets: guestStarManageScopes(), subjectID: req.ModeratorID})
}

func (s *ExperimentalGuestStarService) DeleteGuestStarSlot(ctx context.Context, req DeleteGuestStarSlotRequest) (*Response[DeleteGuestStarSlotData], error) {
	return executeGuestStarEndpoint[DeleteGuestStarSlotData](s.client, ctx, "delete-guest-star-slot", req, nil, guestStarAuthorization{scopeSets: guestStarManageScopes(), subjectID: req.ModeratorID})
}

func (s *ExperimentalGuestStarService) UpdateGuestStarSlotSettings(ctx context.Context, req UpdateGuestStarSlotSettingsRequest) (*Response[UpdateGuestStarSlotSettingsData], error) {
	return executeGuestStarEndpoint[UpdateGuestStarSlotSettingsData](s.client, ctx, "update-guest-star-slot-settings", req, nil, guestStarAuthorization{scopeSets: guestStarManageScopes(), subjectID: req.ModeratorID})
}
