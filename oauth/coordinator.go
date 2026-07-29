package oauth

import (
	"context"
	"errors"
)

var (
	ErrRefreshCoordinator = errors.New("oauth: refresh coordinator failure")
	ErrRefreshLeaseLost   = errors.New("oauth: refresh lease lost")
)

// CredentialLoader loads the persisted credential pair for a user. ExpiresIn
// must be the pair's remaining lifetime, not its original lifetime.
type CredentialLoader func(context.Context, string) (TokenPair, error)

// RefreshCoordinator serializes credential work for one user across processes.
type RefreshCoordinator interface {
	Acquire(context.Context, string) (RefreshLease, error)
}

// RefreshLease owns one user's coordinated credential work until it is released
// or its context is canceled after ownership is lost.
type RefreshLease interface {
	Context() context.Context
	Err() error
	Release(context.Context) error
}
