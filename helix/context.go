package helix

import "errors"

// Intent is an opaque label attached to a managed user credential in a
// CredentialResolver (e.g. "chat" or "eventsub"). Intents let callers ask
// for any user credential that covers a capability without naming the user.
type Intent string

var (
	ErrNoCredentialResolver = errors.New("helix: client has no credential resolver")
	ErrUserNotFound         = errors.New("helix: no credential registered for user")
	ErrIntentNotCovered     = errors.New("helix: no registered credential covers the requested intents")
)

// CredentialResolver resolves a user or a set of intents to a TokenSource.
// Implemented by oauth.Registry; custom multi-tenant stores can implement
// it directly.
type CredentialResolver interface {
	SourceForUser(userID string) (TokenSource, error)
	SourceForIntent(intents ...Intent) (TokenSource, error)
}

// WithCredentialResolver enables multi-tenant context switching on the
// client via AsUser and AsIntent.
func WithCredentialResolver(resolver CredentialResolver) Option {
	return func(client *Client) error {
		if resolver == nil {
			return ErrInvalidOption
		}
		client.resolver = resolver
		return nil
	}
}

// AsUser returns a derived client that executes every operation with the
// credential registered for userID. It performs no network I/O.
func (c *Client) AsUser(userID string) (*Client, error) {
	if err := c.validClient(); err != nil {
		return nil, err
	}
	if c.resolver == nil {
		return nil, ErrNoCredentialResolver
	}
	source, err := c.resolver.SourceForUser(userID)
	if err != nil {
		return nil, err
	}
	return c.derive(source), nil
}

// AsIntent returns a derived client that executes every operation with any
// registered credential whose intents cover all the requested ones. It
// performs no network I/O.
func (c *Client) AsIntent(intents ...Intent) (*Client, error) {
	if err := c.validClient(); err != nil {
		return nil, err
	}
	if c.resolver == nil {
		return nil, ErrNoCredentialResolver
	}
	source, err := c.resolver.SourceForIntent(intents...)
	if err != nil {
		return nil, err
	}
	return c.derive(source), nil
}

// derive shallow-copies the client, sharing the HTTP client, base URL,
// user agent and rate-limit policy with the root, and swaps only the token
// source. Derived clients do not own the credential session; the resolver
// controls its lifecycle.
func (c *Client) derive(source TokenSource) *Client {
	derived := *c
	derived.tokenSource = source
	derived.staticToken = nil
	derived.executor = newTransportExecutor(c.executor.httpClient, source, c.rateLimitPolicy, nil, nil)
	initializeServices(&derived)
	derived.valid = true
	return &derived
}
