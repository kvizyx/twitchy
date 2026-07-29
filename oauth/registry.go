package oauth

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/kvizyx/twitchy/helix"
)

var (
	ErrUserExists     = errors.New("oauth: user already registered")
	ErrRegistryClosed = errors.New("oauth: registry is closed")
)

type registryEntry struct {
	source  *RefreshingTokenSource
	session *ManagedSession
	intents map[helix.Intent]struct{}
}

type registryUserRegistration struct {
	ctx           context.Context
	userID        string
	pair          TokenPair
	hook          CredentialHook
	intents       []helix.Intent
	sourceOptions []SourceOption
}

// Registry is a multi-tenant credential store. Each registered user gets a
// RefreshingTokenSource with a ManagedSession (proactive refresh, hourly
// validation), and the registry as a whole implements
// helix.CredentialResolver so a root helix.Client can derive per-user or
// per-intent clients via AsUser and AsIntent.
type Registry struct {
	client *Client

	mu      sync.RWMutex
	users   map[string]*registryEntry
	closed  bool
	closeEr error
}

func NewRegistry(client *Client) (*Registry, error) {
	if client == nil {
		return nil, ErrInvalidOption
	}
	if err := client.validClient(); err != nil {
		return nil, err
	}
	return &Registry{client: client, users: make(map[string]*registryEntry)}, nil
}

// AddUser registers a user credential with its initial token pair and
// rotation hook, starts a managed session for it, and tags it with intents.
// The first validation happens synchronously, so an invalid pair fails here.
func (r *Registry) AddUser(ctx context.Context, userID string, pair TokenPair, hook CredentialHook, intents ...helix.Intent) error {
	_, err := r.addUser(registryUserRegistration{
		ctx:     ctx,
		userID:  userID,
		pair:    pair,
		hook:    hook,
		intents: intents,
	})
	return err
}

func (r *Registry) addUser(registration registryUserRegistration) (*registryEntry, error) {
	if registration.ctx == nil || registration.userID == "" || registration.hook == nil {
		return nil, ErrInvalidOption
	}
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return nil, ErrRegistryClosed
	}
	if _, exists := r.users[registration.userID]; exists {
		r.mu.Unlock()
		return nil, ErrUserExists
	}
	r.mu.Unlock()

	source, err := NewRefreshingTokenSource(
		r.client,
		registration.pair,
		registration.hook,
		registration.sourceOptions...,
	)
	if err != nil {
		return nil, err
	}
	session, err := NewManagedSession(registration.ctx, source)
	if err != nil {
		_ = source.Close()
		return nil, err
	}
	set := make(map[helix.Intent]struct{}, len(registration.intents))
	for _, intent := range registration.intents {
		set[intent] = struct{}{}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		_ = session.Close()
		return nil, ErrRegistryClosed
	}
	if _, exists := r.users[registration.userID]; exists {
		_ = session.Close()
		return nil, ErrUserExists
	}
	entry := &registryEntry{source: source, session: session, intents: set}
	r.users[registration.userID] = entry
	return entry, nil
}

func (r *Registry) removeUserEntry(userID string, expected *registryEntry) error {
	r.mu.Lock()
	entry, exists := r.users[userID]
	if !exists || entry != expected {
		r.mu.Unlock()
		return nil
	}
	delete(r.users, userID)
	r.mu.Unlock()
	return entry.session.Close()
}

// RemoveUser closes the user's session and drops the credential.
func (r *Registry) RemoveUser(userID string) error {
	r.mu.Lock()
	entry, exists := r.users[userID]
	if exists {
		delete(r.users, userID)
	}
	r.mu.Unlock()
	if !exists {
		return helix.ErrUserNotFound
	}
	return entry.session.Close()
}

// SourceForUser implements helix.CredentialResolver.
func (r *Registry) SourceForUser(userID string) (helix.TokenSource, error) {
	r.mu.RLock()
	entry, exists := r.users[userID]
	r.mu.RUnlock()
	if !exists {
		return nil, helix.ErrUserNotFound
	}
	if err := entry.source.lifecycleError(); err != nil {
		return nil, err
	}
	return entry.source, nil
}

// SourceForIntent implements helix.CredentialResolver. Resolution is
// deterministic: the first user (by sorted user ID) whose intents cover all
// the requested ones wins.
func (r *Registry) SourceForIntent(intents ...helix.Intent) (helix.TokenSource, error) {
	r.mu.RLock()
	ids := make([]string, 0, len(r.users))
	for id := range r.users {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	candidates := make([]*registryEntry, 0, len(ids))
	for _, id := range ids {
		entry := r.users[id]
		if covers(entry.intents, intents) {
			candidates = append(candidates, entry)
		}
	}
	r.mu.RUnlock()
	for _, entry := range candidates {
		if err := entry.source.lifecycleError(); err == nil {
			return entry.source, nil
		}
	}
	return nil, helix.ErrIntentNotCovered
}

func covers(have map[helix.Intent]struct{}, want []helix.Intent) bool {
	for _, intent := range want {
		if _, ok := have[intent]; !ok {
			return false
		}
	}
	return true
}

// Close closes every managed session. It is idempotent.
func (r *Registry) Close() error {
	r.mu.Lock()
	if r.closed {
		err := r.closeEr
		r.mu.Unlock()
		return err
	}
	r.closed = true
	entries := make([]*registryEntry, 0, len(r.users))
	for _, entry := range r.users {
		entries = append(entries, entry)
	}
	r.users = make(map[string]*registryEntry)
	r.mu.Unlock()

	var closeErr error
	for _, entry := range entries {
		if err := entry.session.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
	}
	r.mu.Lock()
	r.closeEr = closeErr
	r.mu.Unlock()
	return closeErr
}
