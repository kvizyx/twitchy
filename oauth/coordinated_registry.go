package oauth

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kvizyx/twitchy/helix"
)

// CoordinatedRegistry is an opt-in registry entry path that loads a durable
// credential pair while holding a distributed refresh lease.
type CoordinatedRegistry struct {
	registry    *Registry
	coordinator RefreshCoordinator
}

var _ helix.CredentialResolver = (*CoordinatedRegistry)(nil)

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
	if isNilCoordinatorValue(lease) || lease.Context() == nil {
		return ErrInvalidOption
	}
	defer func() {
		releaseErr := lease.Release(context.WithoutCancel(ctx))
		if returnErr == nil && releaseErr != nil {
			returnErr = fmt.Errorf("release refresh lease: %w", releaseErr)
		}
	}()

	pair, err := loader(lease.Context(), userID)
	if err := lease.Err(); err != nil {
		return err
	}
	if err != nil {
		return fmt.Errorf("load persisted credential: %w", err)
	}
	return registry.registry.AddUser(ctx, userID, pair, hook, intents...)
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
