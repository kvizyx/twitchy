package helix

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type GuestStarGroupLayout StringEnum

const (
	GuestStarGroupLayoutTiled       GuestStarGroupLayout = "TILED_LAYOUT"
	GuestStarGroupLayoutScreenshare GuestStarGroupLayout = "SCREENSHARE_LAYOUT"
	GuestStarGroupLayoutHorizontal  GuestStarGroupLayout = "HORIZONTAL_LAYOUT"
	GuestStarGroupLayoutVertical    GuestStarGroupLayout = "VERTICAL_LAYOUT"

	GuestStarGroupLayoutTiledLayout       = GuestStarGroupLayoutTiled
	GuestStarGroupLayoutScreenshareLayout = GuestStarGroupLayoutScreenshare
	GuestStarGroupLayoutHorizontalLayout  = GuestStarGroupLayoutHorizontal
	GuestStarGroupLayoutVerticalLayout    = GuestStarGroupLayoutVertical
)

type GetChannelGuestStarSettingsRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
}

type GuestStarChannelSettings struct {
	IsModeratorSendLiveEnabled  bool                 `json:"is_moderator_send_live_enabled"`
	SlotCount                   int                  `json:"slot_count"`
	IsBrowserSourceAudioEnabled bool                 `json:"is_browser_source_audio_enabled"`
	GroupLayout                 GuestStarGroupLayout `json:"group_layout"`
	BrowserSourceToken          string               `json:"browser_source_token"`
}

type GetChannelGuestStarSettingsData []GuestStarChannelSettings

type UpdateChannelGuestStarSettingsRequest struct {
	BroadcasterID               string                `query:"broadcaster_id" json:"-"`
	IsModeratorSendLiveEnabled  *bool                 `query:"-" json:"is_moderator_send_live_enabled,omitempty"`
	SlotCount                   *int                  `query:"-" json:"slot_count,omitempty"`
	IsBrowserSourceAudioEnabled *bool                 `query:"-" json:"is_browser_source_audio_enabled,omitempty"`
	GroupLayout                 *GuestStarGroupLayout `query:"-" json:"group_layout,omitempty"`
	RegenerateBrowserSources    *bool                 `query:"-" json:"regenerate_browser_sources,omitempty"`
}

type UpdateChannelGuestStarSettingsData struct{}

type guestStarAuthorization struct {
	scopeSets [][]AuthorizationScope
	subjectID string
}

func executeGuestStarEndpoint[T any](client *Client, ctx context.Context, anchor string, query, body any, auth guestStarAuthorization) (*Response[T], error) {
	if err := client.validClient(); err != nil {
		return nil, err
	}
	operation, err := manifest.OperationByAnchor(anchor)
	if err != nil {
		return nil, err
	}
	credential, err := client.credential(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateGuestStarCredential(credential, operation, auth); err != nil {
		return nil, err
	}
	request, err := buildRequest(requestSpec{Context: ctx, Method: operation.Method, URL: client.endpointURL(operation.Path), Query: query, Body: body})
	if err != nil {
		return nil, fmt.Errorf("build %s request: %w", operation.OperationID, err)
	}
	request.Header.Set("User-Agent", client.userAgent)
	response, meta, err := client.executor.execute(ctx, request, operation, credential)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	result, err := decodeResponse[T](response.StatusCode, response.Body, DecodeOptions{})
	if err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: []string{credential.AccessToken()}})
	}
	result.Meta = meta
	return result, nil
}

func validateGuestStarCredential(snapshot CredentialSnapshot, operation manifest.Operation, auth guestStarAuthorization) error {
	if snapshot.TokenClass() != TokenClassUser || !guestStarHasScopeSet(snapshot, auth.scopeSets) {
		return localCredentialAuthError(operation.OperationID)
	}
	if auth.subjectID != "" && snapshot.UserID() != auth.subjectID {
		return localCredentialAuthError(operation.OperationID)
	}
	return nil
}

func guestStarHasScopeSet(snapshot CredentialSnapshot, scopeSets [][]AuthorizationScope) bool {
	for _, scopeSet := range scopeSets {
		complete := true
		for _, scope := range scopeSet {
			if !snapshotHasScope(snapshot, scope) {
				complete = false
				break
			}
		}
		if complete {
			return true
		}
	}
	return len(scopeSets) == 0
}

func guestStarReadScopes() [][]AuthorizationScope {
	return [][]AuthorizationScope{{ScopeChannelReadGuestStar}, {ScopeChannelManageGuestStar}, {ScopeModeratorReadGuestStar}, {ScopeModeratorManageGuestStar}}
}

func guestStarManageScopes() [][]AuthorizationScope {
	return [][]AuthorizationScope{{ScopeChannelManageGuestStar}, {ScopeModeratorManageGuestStar}}
}

func (s *ExperimentalGuestStarService) GetChannelGuestStarSettings(ctx context.Context, req GetChannelGuestStarSettingsRequest) (*Response[GetChannelGuestStarSettingsData], error) {
	return executeGuestStarEndpoint[GetChannelGuestStarSettingsData](s.client, ctx, "get-channel-guest-star-settings", req, nil, guestStarAuthorization{scopeSets: guestStarReadScopes(), subjectID: req.ModeratorID})
}

func (s *ExperimentalGuestStarService) UpdateChannelGuestStarSettings(ctx context.Context, req UpdateChannelGuestStarSettingsRequest) (*Response[UpdateChannelGuestStarSettingsData], error) {
	return executeGuestStarEndpoint[UpdateChannelGuestStarSettingsData](s.client, ctx, "update-channel-guest-star-settings", req, req, guestStarAuthorization{scopeSets: [][]AuthorizationScope{{ScopeChannelManageGuestStar}}, subjectID: req.BroadcasterID})
}
