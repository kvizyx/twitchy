package helix

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type GetUsersRequest struct {
	IDs    []string `query:"id,omitempty"`
	Logins []string `query:"login,omitempty"`
}

type UserType string

const (
	UserTypeAdmin     UserType = "admin"
	UserTypeGlobalMod UserType = "global_mod"
	UserTypeStaff     UserType = "staff"
	UserTypeNormal    UserType = ""
)

type BroadcasterType string

const (
	BroadcasterTypeAffiliate BroadcasterType = "affiliate"
	BroadcasterTypePartner   BroadcasterType = "partner"
	BroadcasterTypeNormal    BroadcasterType = ""
)

type User struct {
	ID              string          `json:"id"`
	Login           string          `json:"login"`
	DisplayName     string          `json:"display_name"`
	Type            UserType        `json:"type"`
	BroadcasterType BroadcasterType `json:"broadcaster_type"`
	Description     string          `json:"description"`
	ProfileImageURL string          `json:"profile_image_url"`
	OfflineImageURL string          `json:"offline_image_url"`
	ViewCount       int64           `json:"view_count"`
	Email           string          `json:"email"`
	CreatedAt       Timestamp       `json:"created_at"`
}

type GetUsersData []User

type UpdateUserRequest struct {
	Description *string `query:"description,omitempty" json:"-"`
}

type UpdateUserData []User

type UserExtensionType string

const (
	UserExtensionTypeComponent UserExtensionType = "component"
	UserExtensionTypeMobile    UserExtensionType = "mobile"
	UserExtensionTypeOverlay   UserExtensionType = "overlay"
	UserExtensionTypePanel     UserExtensionType = "panel"
)

type InstalledUserExtension struct {
	ID          string              `json:"id"`
	Version     string              `json:"version"`
	Name        string              `json:"name"`
	CanActivate bool                `json:"can_activate"`
	Type        []UserExtensionType `json:"type"`
}

type GetUserExtensionsRequest struct{}
type GetUserExtensionsData []InstalledUserExtension

func (s *UsersService) GetUsers(ctx context.Context, req GetUsersRequest) (*Response[GetUsersData], error) {
	return executeTask26Endpoint[GetUsersData](s.client, ctx, "get-users", req, nil, nil, "")
}

func (s *UsersService) UpdateUser(ctx context.Context, req UpdateUserRequest) (*Response[UpdateUserData], error) {
	return executeTask26Endpoint[UpdateUserData](s.client, ctx, "update-user", req, nil, nil, "")
}

func (s *UsersService) GetUserExtensions(ctx context.Context, req GetUserExtensionsRequest) (*Response[GetUserExtensionsData], error) {
	return executeTask26Endpoint[GetUserExtensionsData](s.client, ctx, "get-user-extensions", req, nil, [][]AuthorizationScope{{ScopeUserReadBroadcast}, {ScopeUserEditBroadcast}}, "")
}

type task26Authorization struct {
	scopeSets          [][]AuthorizationScope
	subjectID          string
	appSubjectRequired bool
}

func executeTask26Endpoint[T any](client *Client, ctx context.Context, anchor string, query, body any, scopeSets [][]AuthorizationScope, subjectID string) (*Response[T], error) {
	return executeTask26EndpointWithSecrets[T](client, ctx, anchor, query, body, task26Authorization{scopeSets: scopeSets, subjectID: subjectID}, nil)
}

func executeTask26EndpointWithSecrets[T any](client *Client, ctx context.Context, anchor string, query, body any, auth task26Authorization, secrets []string) (*Response[T], error) {
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
	if err := validateTask26Credential(credential, operation, auth); err != nil {
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
		redactionSecrets := append([]string{credential.AccessToken()}, secrets...)
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: redactionSecrets})
	}
	result.Meta = meta
	return result, nil
}

func validateTask26Credential(snapshot CredentialSnapshot, operation manifest.Operation, auth task26Authorization) error {
	if !operationAllowsTokenClass(snapshot.TokenClass(), operation.TokenClasses) {
		return localCredentialAuthError(operation.OperationID)
	}
	// App access tokens carry no scopes: Twitch authorizes them through prior
	// user grants, which cannot be verified client-side, so scope checks only
	// apply to user access tokens.
	if snapshot.TokenClass() != TokenClassApp {
		if len(auth.scopeSets) == 0 {
			if err := validateCredentialForOperation(snapshot, operation, "", ""); err != nil {
				return err
			}
		} else if !task26HasScopeSet(snapshot, auth.scopeSets) {
			return localCredentialAuthError(operation.OperationID)
		}
	}
	if snapshot.TokenClass() == TokenClassApp && auth.appSubjectRequired && auth.subjectID == "" {
		return localCredentialAuthError(operation.OperationID)
	}
	if snapshot.TokenClass() == TokenClassUser && auth.subjectID != "" && snapshot.UserID() != auth.subjectID {
		return localCredentialAuthError(operation.OperationID)
	}
	return nil
}

func task26HasScopeSet(snapshot CredentialSnapshot, scopeSets [][]AuthorizationScope) bool {
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
	return false
}
