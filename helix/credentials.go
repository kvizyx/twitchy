package helix

import (
	"context"
	"time"
)

type TokenClass string

const (
	TokenClassApp       TokenClass = "app"
	TokenClassUser      TokenClass = "user"
	TokenClassExtension TokenClass = "extension"
)

type Credential struct {
	AccessToken string
	ClientID    string
	TokenClass  TokenClass
	UserID      string
	Scopes      []AuthorizationScope
	ExpiresAt   time.Time
	Refreshable bool
	Generation  uint64
}

func (Credential) String() string   { return "helix.Credential{redacted}" }
func (Credential) GoString() string { return "helix.Credential{redacted}" }

type CredentialSnapshot struct {
	accessToken string
	clientID    string
	tokenClass  TokenClass
	userID      string
	scopes      []AuthorizationScope
	expiresAt   time.Time
	refreshable bool
	generation  uint64
}

func NewCredentialSnapshot(credential Credential) CredentialSnapshot {
	return CredentialSnapshot{
		accessToken: credential.AccessToken,
		clientID:    credential.ClientID,
		tokenClass:  credential.TokenClass,
		userID:      credential.UserID,
		scopes:      append([]AuthorizationScope(nil), credential.Scopes...),
		expiresAt:   credential.ExpiresAt,
		refreshable: credential.Refreshable,
		generation:  credential.Generation,
	}
}

func (snapshot CredentialSnapshot) AccessToken() string { return snapshot.accessToken }
func (snapshot CredentialSnapshot) ClientID() string    { return snapshot.clientID }
func (snapshot CredentialSnapshot) TokenClass() TokenClass {
	return snapshot.tokenClass
}
func (snapshot CredentialSnapshot) UserID() string { return snapshot.userID }
func (snapshot CredentialSnapshot) Scopes() []AuthorizationScope {
	return append([]AuthorizationScope(nil), snapshot.scopes...)
}
func (snapshot CredentialSnapshot) ExpiresAt() time.Time { return snapshot.expiresAt }
func (snapshot CredentialSnapshot) Refreshable() bool    { return snapshot.refreshable }
func (snapshot CredentialSnapshot) Generation() uint64   { return snapshot.generation }

func (CredentialSnapshot) String() string   { return "helix.CredentialSnapshot{redacted}" }
func (CredentialSnapshot) GoString() string { return "helix.CredentialSnapshot{redacted}" }

type TokenSource interface {
	Token(context.Context) (CredentialSnapshot, error)
}
