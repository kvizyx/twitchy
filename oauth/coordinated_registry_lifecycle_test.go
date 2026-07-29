package oauth

import (
	"context"
	"errors"
	"testing"

	"github.com/kvizyx/twitchy/helix"
)

type sequencedContractLease struct {
	ctx          context.Context
	errs         []error
	releaseErr   error
	errCalls     int
	releaseCalls int
}

func (lease *sequencedContractLease) Context() context.Context { return lease.ctx }

func (lease *sequencedContractLease) Err() error {
	if lease.errCalls >= len(lease.errs) {
		return nil
	}
	err := lease.errs[lease.errCalls]
	lease.errCalls++
	return err
}

func (lease *sequencedContractLease) Release(context.Context) error {
	lease.releaseCalls++
	return lease.releaseErr
}

func TestCoordinatedRegistry_rollsBackNewUserWhenLeaseLostAfterRegistration(t *testing.T) {
	lease := &sequencedContractLease{
		ctx:  context.Background(),
		errs: []error{nil, ErrRefreshLeaseLost},
	}
	registry := newSequencedCoordinatedRegistry(t, lease)

	err := registry.AddCoordinatedUser(context.Background(), "42", func(
		context.Context,
		string,
	) (TokenPair, error) {
		return registryPair("coordinated"), nil
	}, noopHook)
	if !errors.Is(err, ErrRefreshLeaseLost) {
		t.Fatalf("AddCoordinatedUser error = %v, want lease loss", err)
	}
	if _, err = registry.SourceForUser("42"); !errors.Is(err, helix.ErrUserNotFound) {
		t.Fatalf("SourceForUser error = %v, want rolled back user", err)
	}
	if lease.releaseCalls != 1 {
		t.Fatalf("release calls = %d, want 1", lease.releaseCalls)
	}
}

func TestCoordinatedRegistry_rollsBackNewUserWhenReleaseFails(t *testing.T) {
	releaseErr := errors.New("release unavailable")
	lease := &sequencedContractLease{ctx: context.Background(), releaseErr: releaseErr}
	registry := newSequencedCoordinatedRegistry(t, lease)

	err := registry.AddCoordinatedUser(context.Background(), "42", func(
		context.Context,
		string,
	) (TokenPair, error) {
		return registryPair("coordinated"), nil
	}, noopHook)
	if !errors.Is(err, releaseErr) {
		t.Fatalf("AddCoordinatedUser error = %v, want release failure", err)
	}
	if _, err = registry.SourceForUser("42"); !errors.Is(err, helix.ErrUserNotFound) {
		t.Fatalf("SourceForUser error = %v, want rolled back user", err)
	}
}

func newSequencedCoordinatedRegistry(t *testing.T, lease RefreshLease) *CoordinatedRegistry {
	t.Helper()
	server := registryServer(t)
	client, err := New(WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}
	registry, err := NewCoordinatedRegistry(client, contractCoordinator{lease: lease})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = registry.Close() })
	return registry
}
