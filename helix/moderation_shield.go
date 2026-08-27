package helix

import "context"

type UpdateShieldModeStatusRequest struct {
	BroadcasterID string `query:"broadcaster_id" json:"-"`
	ModeratorID   string `query:"moderator_id" json:"-"`
	IsActive      bool   `query:"-" json:"is_active"`
}

type ShieldModeStatus struct {
	IsActive        bool      `json:"is_active"`
	ModeratorID     string    `json:"moderator_id"`
	ModeratorLogin  string    `json:"moderator_login"`
	ModeratorName   string    `json:"moderator_name"`
	LastActivatedAt Timestamp `json:"last_activated_at"`
}

type ShieldMode = ShieldModeStatus

type UpdateShieldModeStatusData []ShieldModeStatus
type GetShieldModeStatusData []ShieldModeStatus

type GetShieldModeStatusRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
}

func (s *ModerationService) UpdateShieldModeStatus(ctx context.Context, req UpdateShieldModeStatusRequest) (*Response[UpdateShieldModeStatusData], error) {
	return executeModerationEndpoint[UpdateShieldModeStatusData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "update-shield-mode-status", query: req, body: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModeratorManageShieldMode), subjectIDs: []string{req.ModeratorID}}})
}

func (s *ModerationService) GetShieldModeStatus(ctx context.Context, req GetShieldModeStatusRequest) (*Response[GetShieldModeStatusData], error) {
	scopes := [][]AuthorizationScope{{ScopeModeratorReadShieldMode}, {ScopeModeratorManageShieldMode}}
	return executeModerationEndpoint[GetShieldModeStatusData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "get-shield-mode-status", query: req, auth: moderationAuthorization{userScopeSets: scopes, subjectIDs: []string{req.ModeratorID}}})
}
