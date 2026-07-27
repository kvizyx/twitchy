package helix

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type CheckAutoModStatusMessage struct {
	MsgID   string `json:"msg_id"`
	MsgText string `json:"msg_text"`
}

type AutoModMessage = CheckAutoModStatusMessage

type CheckAutoModStatusResult struct {
	MsgID       string `json:"msg_id"`
	IsPermitted bool   `json:"is_permitted"`
}

type AutoModStatus = CheckAutoModStatusResult
type CheckAutoModStatusData []CheckAutoModStatusResult

type CheckAutoModStatusRequest struct {
	BroadcasterID string                      `query:"broadcaster_id" json:"-"`
	Data          []CheckAutoModStatusMessage `query:"-" json:"data"`
}

type ManageHeldAutoModMessagesRequest struct {
	UserID string `query:"-" json:"user_id"`
	MsgID  string `query:"-" json:"msg_id"`
	Action string `query:"-" json:"action"`
}

type ManageHeldAutoModMessagesData struct{}

type AutoModSettings struct {
	BroadcasterID           string `json:"broadcaster_id"`
	ModeratorID             string `json:"moderator_id"`
	OverallLevel            *int   `json:"overall_level"`
	Disability              int    `json:"disability"`
	Aggression              int    `json:"aggression"`
	SexualitySexOrGender    int    `json:"sexuality_sex_or_gender"`
	Misogyny                int    `json:"misogyny"`
	Bullying                int    `json:"bullying"`
	Swearing                int    `json:"swearing"`
	RaceEthnicityOrReligion int    `json:"race_ethnicity_or_religion"`
	SexBasedTerms           int    `json:"sex_based_terms"`
}

type GetAutoModSettingsRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
}

type GetAutoModSettingsData []AutoModSettings
type UpdateAutoModSettingsData []AutoModSettings

type UpdateAutoModSettingsRequest struct {
	BroadcasterID           string `query:"broadcaster_id" json:"-"`
	ModeratorID             string `query:"moderator_id" json:"-"`
	Aggression              *int   `query:"-" json:"aggression,omitempty"`
	Bullying                *int   `query:"-" json:"bullying,omitempty"`
	Disability              *int   `query:"-" json:"disability,omitempty"`
	Misogyny                *int   `query:"-" json:"misogyny,omitempty"`
	OverallLevel            *int   `query:"-" json:"overall_level,omitempty"`
	RaceEthnicityOrReligion *int   `query:"-" json:"race_ethnicity_or_religion,omitempty"`
	SexBasedTerms           *int   `query:"-" json:"sex_based_terms,omitempty"`
	SexualitySexOrGender    *int   `query:"-" json:"sexuality_sex_or_gender,omitempty"`
	Swearing                *int   `query:"-" json:"swearing,omitempty"`
}

type moderationAuthorization struct {
	userScopeSets [][]AuthorizationScope
	appScopeSets  [][]AuthorizationScope
	subjectIDs    []string
}

type moderationEndpointSpec struct {
	client *Client
	ctx    context.Context
	anchor string
	query  any
	body   any
	auth   moderationAuthorization
}

func executeModerationEndpoint[T any](spec moderationEndpointSpec) (*Response[T], error) {
	if err := spec.client.validClient(); err != nil {
		return nil, err
	}
	operation, err := manifest.OperationByAnchor(spec.anchor)
	if err != nil {
		return nil, err
	}
	credential, err := spec.client.credential(spec.ctx)
	if err != nil {
		return nil, err
	}
	if err := validateModerationCredential(credential, operation, spec.auth); err != nil {
		return nil, err
	}
	request, err := buildRequest(requestSpec{Context: spec.ctx, Method: operation.Method, URL: spec.client.endpointURL(operation.Path), Query: spec.query, Body: spec.body})
	if err != nil {
		return nil, fmt.Errorf("build %s request: %w", operation.OperationID, err)
	}
	request.Header.Set("User-Agent", spec.client.userAgent)
	response, meta, err := spec.client.executor.execute(spec.ctx, request, operation, credential)
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

func validateModerationCredential(snapshot CredentialSnapshot, operation manifest.Operation, auth moderationAuthorization) error {
	if !operationAllowsTokenClass(snapshot.TokenClass(), operation.TokenClasses) {
		return localCredentialAuthError(operation.OperationID)
	}
	scopeSets := auth.userScopeSets
	if snapshot.TokenClass() == TokenClassApp {
		scopeSets = auth.appScopeSets
	}
	if snapshot.TokenClass() != TokenClassUser && snapshot.TokenClass() != TokenClassApp {
		return localCredentialAuthError(operation.OperationID)
	}
	if !chatHasScopeSet(snapshot, scopeSets) {
		return localCredentialAuthError(operation.OperationID)
	}
	if snapshot.TokenClass() == TokenClassUser && len(auth.subjectIDs) > 0 {
		for _, subjectID := range auth.subjectIDs {
			if subjectID != "" && snapshot.UserID() == subjectID {
				return nil
			}
		}
		return localCredentialAuthError(operation.OperationID)
	}
	return nil
}

func (s *ModerationService) CheckAutoModStatus(ctx context.Context, req CheckAutoModStatusRequest) (*Response[CheckAutoModStatusData], error) {
	return executeModerationEndpoint[CheckAutoModStatusData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "check-automod-status", query: req, body: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModerationRead), appScopeSets: chatReadScopes(ScopeModerationRead), subjectIDs: []string{req.BroadcasterID}}})
}

func (s *ModerationService) ManageHeldAutoModMessages(ctx context.Context, req ManageHeldAutoModMessagesRequest) (*Response[ManageHeldAutoModMessagesData], error) {
	return executeModerationEndpoint[ManageHeldAutoModMessagesData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "manage-held-automod-messages", body: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModeratorManageAutoMod), appScopeSets: chatReadScopes(ScopeModeratorManageAutoMod), subjectIDs: []string{req.UserID}}})
}

func (s *ModerationService) GetAutoModSettings(ctx context.Context, req GetAutoModSettingsRequest) (*Response[GetAutoModSettingsData], error) {
	scopes := [][]AuthorizationScope{{ScopeModeratorReadAutoModSettings}, {ScopeModeratorManageAutoModSettings}}
	return executeModerationEndpoint[GetAutoModSettingsData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "get-automod-settings", query: req, auth: moderationAuthorization{userScopeSets: scopes, appScopeSets: scopes, subjectIDs: []string{req.ModeratorID}}})
}

func (s *ModerationService) UpdateAutoModSettings(ctx context.Context, req UpdateAutoModSettingsRequest) (*Response[UpdateAutoModSettingsData], error) {
	return executeModerationEndpoint[UpdateAutoModSettingsData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "update-automod-settings", query: req, body: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModeratorManageAutoModSettings), appScopeSets: chatReadScopes(ScopeModeratorManageAutoModSettings), subjectIDs: []string{req.ModeratorID}}})
}
