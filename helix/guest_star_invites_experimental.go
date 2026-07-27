package helix

import "context"

type GetGuestStarInvitesRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
	SessionID     string `query:"session_id"`
}

type GuestStarInviteStatus StringEnum

const (
	GuestStarInviteStatusInvited  GuestStarInviteStatus = "INVITED"
	GuestStarInviteStatusAccepted GuestStarInviteStatus = "ACCEPTED"
	GuestStarInviteStatusReady    GuestStarInviteStatus = "READY"
)

type GuestStarInvite struct {
	UserID           string                `json:"user_id"`
	InvitedAt        Timestamp             `json:"invited_at"`
	Status           GuestStarInviteStatus `json:"status"`
	IsVideoEnabled   bool                  `json:"is_video_enabled"`
	IsAudioEnabled   bool                  `json:"is_audio_enabled"`
	IsVideoAvailable bool                  `json:"is_video_available"`
	IsAudioAvailable bool                  `json:"is_audio_available"`
}

type GetGuestStarInvitesData []GuestStarInvite

type SendGuestStarInviteRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
	SessionID     string `query:"session_id"`
	GuestID       string `query:"guest_id"`
}

type SendGuestStarInviteData struct{}

type DeleteGuestStarInviteRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
	SessionID     string `query:"session_id"`
	GuestID       string `query:"guest_id"`
}

type DeleteGuestStarInviteData struct{}

func (s *ExperimentalGuestStarService) GetGuestStarInvites(ctx context.Context, req GetGuestStarInvitesRequest) (*Response[GetGuestStarInvitesData], error) {
	return executeGuestStarEndpoint[GetGuestStarInvitesData](s.client, ctx, "get-guest-star-invites", req, nil, guestStarAuthorization{scopeSets: guestStarReadScopes(), subjectID: req.ModeratorID})
}

func (s *ExperimentalGuestStarService) SendGuestStarInvite(ctx context.Context, req SendGuestStarInviteRequest) (*Response[SendGuestStarInviteData], error) {
	return executeGuestStarEndpoint[SendGuestStarInviteData](s.client, ctx, "send-guest-star-invite", req, nil, guestStarAuthorization{scopeSets: guestStarManageScopes(), subjectID: req.ModeratorID})
}

func (s *ExperimentalGuestStarService) DeleteGuestStarInvite(ctx context.Context, req DeleteGuestStarInviteRequest) (*Response[DeleteGuestStarInviteData], error) {
	return executeGuestStarEndpoint[DeleteGuestStarInviteData](s.client, ctx, "delete-guest-star-invite", req, nil, guestStarAuthorization{scopeSets: guestStarManageScopes(), subjectID: req.ModeratorID})
}
