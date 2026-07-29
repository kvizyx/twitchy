package oauth

import (
	"context"
	"errors"
	"testing"

	"github.com/kvizyx/twitchy/helix"
)

type contractCoordinator struct {
	lease RefreshLease
	err   error
}

type nilContractCoordinator struct{}

func (*nilContractCoordinator) Acquire(context.Context, string) (RefreshLease, error) {
	return nil, nil
}

func (coordinator contractCoordinator) Acquire(context.Context, string) (RefreshLease, error) {
	return coordinator.lease, coordinator.err
}

type contractLease struct {
	ctx      context.Context
	released bool
}

func (lease *contractLease) Context() context.Context {
	return lease.ctx
}

func (lease *contractLease) Err() error {
	return lease.ctx.Err()
}

func (lease *contractLease) Release(context.Context) error {
	lease.released = true
	return nil
}

func TestLegacyRegistryFunctionTypes(t *testing.T) {
	var _ func(*Client) (*Registry, error) = NewRegistry
	var _ func(*Registry, context.Context, string, TokenPair, CredentialHook, ...helix.Intent) error = (*Registry).AddUser
}

func TestCoordinatedRegistryContracts(t *testing.T) {
	if _, err := NewCoordinatedRegistry(nil, contractCoordinator{}); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("nil client error = %v, want ErrInvalidOption", err)
	}

	server := registryServer(t)
	client, err := New(WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = NewCoordinatedRegistry(client, nil); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("nil coordinator error = %v, want ErrInvalidOption", err)
	}
	var nilCoordinator *nilContractCoordinator
	if _, err = NewCoordinatedRegistry(client, nilCoordinator); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("typed nil coordinator error = %v, want ErrInvalidOption", err)
	}

	leaseContext, cancel := context.WithCancel(context.Background())
	defer cancel()
	lease := &contractLease{ctx: leaseContext}
	registry, err := NewCoordinatedRegistry(client, contractCoordinator{lease: lease})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = registry.Close() })

	loaded := false
	err = registry.AddCoordinatedUser(context.Background(), "111", func(
		ctx context.Context,
		userID string,
	) (TokenPair, error) {
		if ctx != leaseContext {
			t.Fatal("loader did not receive the lease context")
		}
		if userID != "111" {
			t.Fatalf("loader userID = %q, want 111", userID)
		}
		loaded = true
		return registryPair("one"), nil
	}, noopHook, "chat")
	if err != nil {
		t.Fatal(err)
	}
	if !loaded || !lease.released {
		t.Fatalf("loader called = %t, lease released = %t", loaded, lease.released)
	}

	for _, test := range []struct {
		name   string
		userID string
		loader CredentialLoader
		hook   CredentialHook
	}{
		{
			name:   "empty user",
			loader: func(context.Context, string) (TokenPair, error) { return TokenPair{}, nil },
			hook:   noopHook,
		},
		{name: "nil loader", userID: "222", hook: noopHook},
		{
			name:   "nil hook",
			userID: "333",
			loader: func(context.Context, string) (TokenPair, error) { return TokenPair{}, nil },
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := registry.AddCoordinatedUser(context.Background(), test.userID, test.loader, test.hook)
			if !errors.Is(err, ErrInvalidOption) {
				t.Fatalf("AddCoordinatedUser error = %v, want ErrInvalidOption", err)
			}
		})
	}
	invalidRegistry, err := NewCoordinatedRegistry(client, contractCoordinator{lease: (*contractLease)(nil)})
	if err != nil {
		t.Fatal(err)
	}
	if err = invalidRegistry.AddCoordinatedUser(context.Background(), "444", func(
		context.Context,
		string,
	) (TokenPair, error) {
		return registryPair("four"), nil
	}, noopHook); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("nil lease error = %v, want ErrInvalidOption", err)
	}
}
