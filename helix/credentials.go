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

type CredentialSnapshot struct{}

type TokenSource interface {
	Token(context.Context) (CredentialSnapshot, error)
}
