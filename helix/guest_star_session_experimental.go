package helix

import "context"

type GetGuestStarSessionRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
}

type CreateGuestStarSessionRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type EndGuestStarSessionRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	SessionID     string `query:"session_id"`
}

type GuestStarMediaSettings struct {
	IsHostEnabled  bool `json:"is_host_enabled"`
	IsGuestEnabled bool `json:"is_guest_enabled"`
	IsAvailable    bool `json:"is_available"`
}

type GuestStarGuest struct {
	ID              string                 `json:"id"`
	SlotID          string                 `json:"slot_id"`
	UserID          string                 `json:"user_id"`
	UserDisplayName string                 `json:"user_display_name"`
	UserLogin       string                 `json:"user_login"`
	IsLive          bool                   `json:"is_live"`
	Volume          int                    `json:"volume"`
	AssignedAt      Timestamp              `json:"assigned_at"`
	AudioSettings   GuestStarMediaSettings `json:"audio_settings"`
	VideoSettings   GuestStarMediaSettings `json:"video_settings"`
}

type GuestStarSession struct {
	ID     string           `json:"id"`
	Guests []GuestStarGuest `json:"guests"`
}

type GetGuestStarSessionData []GuestStarSession
type CreateGuestStarSessionData []GuestStarSession
type EndGuestStarSessionData []GuestStarSession

func (s *ExperimentalGuestStarService) GetGuestStarSession(ctx context.Context, req GetGuestStarSessionRequest) (*Response[GetGuestStarSessionData], error) {
	return executeGuestStarEndpoint[GetGuestStarSessionData](s.client, ctx, "get-guest-star-session", req, nil, guestStarAuthorization{scopeSets: guestStarReadScopes()})
}

func (s *ExperimentalGuestStarService) CreateGuestStarSession(ctx context.Context, req CreateGuestStarSessionRequest) (*Response[CreateGuestStarSessionData], error) {
	return executeGuestStarEndpoint[CreateGuestStarSessionData](s.client, ctx, "create-guest-star-session", req, nil, guestStarAuthorization{scopeSets: [][]AuthorizationScope{{ScopeChannelManageGuestStar}}, subjectID: req.BroadcasterID})
}

func (s *ExperimentalGuestStarService) EndGuestStarSession(ctx context.Context, req EndGuestStarSessionRequest) (*Response[EndGuestStarSessionData], error) {
	return executeGuestStarEndpoint[EndGuestStarSessionData](s.client, ctx, "end-guest-star-session", req, nil, guestStarAuthorization{scopeSets: [][]AuthorizationScope{{ScopeChannelManageGuestStar}}, subjectID: req.BroadcasterID})
}
