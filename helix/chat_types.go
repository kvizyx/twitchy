package helix

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type GetChattersRequest struct {
	BroadcasterID string  `query:"broadcaster_id" json:"-"`
	ModeratorID   string  `query:"moderator_id" json:"-"`
	First         *int    `query:"first,omitempty" json:"-"`
	After         *string `query:"after,omitempty" json:"-"`
}

type ChatChatter struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type GetChattersData []ChatChatter

func (s *ChatService) GetChatters(ctx context.Context, req GetChattersRequest) (*Response[GetChattersData], error) {
	return executeChatEndpoint[GetChattersData](chatEndpointSpec{
		client: s.client,
		ctx:    ctx,
		anchor: "get-chatters",
		query:  req,
		auth: chatAuthorization{
			userScopeSets: chatReadScopes(ScopeModeratorReadChatters),
			subjectID:     req.ModeratorID,
		},
	})
}

func (s *ChatService) GetChattersPager(req GetChattersRequest, opts ...PagerOption) (*Pager[GetChattersData], error) {
	initialCursor := ""
	if req.After != nil {
		initialCursor = *req.After
	}
	request := req
	request.After = nil
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetChattersData], error) {
		pageRequest := request
		if cursor != "" {
			pageRequest.After = &cursor
		}
		return s.GetChatters(ctx, pageRequest)
	}, initialCursor, opts...)
}

type chatAuthorization struct {
	userScopeSets           [][]AuthorizationScope
	subjectID               string
	rejectForSourceOnlyUser bool
}

type chatEndpointSpec struct {
	client *Client
	ctx    context.Context
	anchor string
	query  any
	body   any
	auth   chatAuthorization
}

func executeChatEndpoint[T any](spec chatEndpointSpec) (*Response[T], error) {
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
	if err := validateChatCredential(credential, operation, spec.auth); err != nil {
		return nil, err
	}
	request, err := buildRequest(requestSpec{
		Context: spec.ctx,
		Method:  operation.Method,
		URL:     spec.client.endpointURL(operation.Path),
		Query:   spec.query,
		Body:    spec.body,
	})
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
		return nil, newProtocolError(errorInput{
			operation:  operation.OperationID,
			statusCode: response.StatusCode,
			meta:       meta,
			cause:      err,
			secrets:    []string{credential.AccessToken()},
		})
	}
	result.Meta = meta
	return result, nil
}

func validateChatCredential(snapshot CredentialSnapshot, operation manifest.Operation, auth chatAuthorization) error {
	if !operationAllowsTokenClass(snapshot.TokenClass(), operation.TokenClasses) {
		return localCredentialAuthError(operation.OperationID)
	}
	if snapshot.TokenClass() != TokenClassUser && snapshot.TokenClass() != TokenClassApp {
		return localCredentialAuthError(operation.OperationID)
	}
	// App access tokens carry no scopes: Twitch authorizes them through prior
	// user grants, which cannot be verified client-side, so scope checks only
	// apply to user access tokens.
	if snapshot.TokenClass() == TokenClassUser {
		if !chatHasScopeSet(snapshot, auth.userScopeSets) {
			return localCredentialAuthError(operation.OperationID)
		}
		if auth.rejectForSourceOnlyUser {
			return localCredentialAuthError(operation.OperationID)
		}
		if auth.subjectID != "" && snapshot.UserID() != auth.subjectID {
			return localCredentialAuthError(operation.OperationID)
		}
	}
	return nil
}

func chatHasScopeSet(snapshot CredentialSnapshot, scopeSets [][]AuthorizationScope) bool {
	if len(scopeSets) == 0 {
		return true
	}
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

func chatReadScopes(scope AuthorizationScope) [][]AuthorizationScope {
	return [][]AuthorizationScope{{scope}}
}
