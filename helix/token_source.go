package helix

import (
	"context"
	"errors"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type RefreshReason string

const (
	RefreshReasonExpired      RefreshReason = "expired"
	RefreshReasonUnauthorized RefreshReason = "unauthorized"
)

type RefreshableTokenSource interface {
	Token(context.Context) (CredentialSnapshot, error)
	Refresh(context.Context, CredentialSnapshot, RefreshReason) (CredentialSnapshot, error)
}

type StaticTokenSource struct {
	snapshot CredentialSnapshot
}

func NewStaticTokenSource(credential Credential) *StaticTokenSource {
	return &StaticTokenSource{snapshot: NewCredentialSnapshot(credential)}
}

func (source *StaticTokenSource) Token(ctx context.Context) (CredentialSnapshot, error) {
	if err := ctx.Err(); err != nil {
		return CredentialSnapshot{}, err
	}
	return source.snapshot, nil
}

var errLocalCredentialRejected = errors.New("credential rejected by local authorization constraints")

func validateCredentialForOperation(snapshot CredentialSnapshot, operation manifest.Operation, clientID, subjectID string) error {
	if clientID != "" && snapshot.ClientID() != clientID {
		return localCredentialAuthError(operation.OperationID)
	}
	if !operationAllowsTokenClass(snapshot.TokenClass(), operation.TokenClasses) {
		return localCredentialAuthError(operation.OperationID)
	}
	for _, scope := range operation.Scopes {
		if scope == "" || scope == "unknown" {
			continue
		}
		if !snapshotHasScope(snapshot, AuthorizationScope(scope)) {
			return localCredentialAuthError(operation.OperationID)
		}
	}
	if operation.SubjectBinding != "" && operation.SubjectBinding != "unknown" && subjectID != "" && snapshot.UserID() != subjectID {
		return localCredentialAuthError(operation.OperationID)
	}
	return nil
}

func operationAllowsTokenClass(tokenClass TokenClass, allowed []manifest.TokenClass) bool {
	for _, candidate := range allowed {
		if candidate == manifest.TokenClassUnknown {
			return true
		}
		if TokenClass(candidate) == tokenClass {
			return true
		}
	}
	return len(allowed) == 0
}

func snapshotHasScope(snapshot CredentialSnapshot, wanted AuthorizationScope) bool {
	for _, scope := range snapshot.scopes {
		if scope == wanted {
			return true
		}
	}
	return false
}

func localCredentialAuthError(operationID string) error {
	return &AuthError{details: errorDetails{operation: operationID, cause: errLocalCredentialRejected}}
}
