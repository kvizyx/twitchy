package oauth

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

// CoordinatedRegistry is an opt-in registry entry path that loads a durable
// credential pair while holding a distributed refresh lease.
type CoordinatedRegistry struct {
	registry    *Registry
	coordinator RefreshCoordinator
}

var _ helix.CredentialResolver = (*CoordinatedRegistry)(nil)

const coordinatedReleaseTimeout = time.Second

func NewCoordinatedRegistry(client *Client, coordinator RefreshCoordinator) (*CoordinatedRegistry, error) {
	if isNilCoordinatorValue(coordinator) {
		return nil, ErrInvalidOption
	}
	registry, err := NewRegistry(client)
	if err != nil {
		return nil, err
	}
	return &CoordinatedRegistry{registry: registry, coordinator: coordinator}, nil
}

// AddCoordinatedUser loads the authoritative persisted pair after it acquires
// the user's lease, then registers the user through the legacy registry path.
func (registry *CoordinatedRegistry) AddCoordinatedUser(
	ctx context.Context,
	userID string,
	loader CredentialLoader,
	hook CredentialHook,
	intents ...helix.Intent,
) (returnErr error) {
	if ctx == nil || userID == "" || loader == nil || hook == nil {
		return ErrInvalidOption
	}
	lease, err := registry.coordinator.Acquire(ctx, userID)
	if err != nil {
		return fmt.Errorf("acquire refresh lease: %w", err)
	}
	if isNilCoordinatorValue(lease) {
		return ErrInvalidOption
	}
	defer func() {
		releaseContext, cancel := context.WithTimeout(context.WithoutCancel(ctx), coordinatedReleaseTimeout)
		defer cancel()
		returnErr = errors.Join(returnErr, lease.Release(releaseContext))
	}()

	leaseContext := lease.Context()
	if leaseContext == nil {
		return ErrInvalidOption
	}
	pair, loaderErr := loader(leaseContext, userID)
	if leaseErr := lease.Err(); loaderErr != nil || leaseErr != nil {
		return errors.Join(loaderErr, leaseErr)
	}
	return errors.Join(registry.registry.AddUser(ctx, userID, pair, hook, intents...), lease.Err())
}

func (registry *CoordinatedRegistry) SourceForUser(userID string) (helix.TokenSource, error) {
	return registry.registry.SourceForUser(userID)
}

func (registry *CoordinatedRegistry) SourceForIntent(intents ...helix.Intent) (helix.TokenSource, error) {
	return registry.registry.SourceForIntent(intents...)
}

func (registry *CoordinatedRegistry) RemoveUser(userID string) error {
	return registry.registry.RemoveUser(userID)
}

func (registry *CoordinatedRegistry) Close() error {
	return registry.registry.Close()
}

func isNilCoordinatorValue(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	kind := reflected.Kind()
	if kind == reflect.Chan || kind == reflect.Func || kind == reflect.Interface ||
		kind == reflect.Map || kind == reflect.Pointer || kind == reflect.Slice {
		return reflected.IsNil()
	}
	return false
}
