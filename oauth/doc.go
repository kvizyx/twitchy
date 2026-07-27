// Package oauth implements Twitch OAuth authorization, token exchange,
// validation, refresh, revocation, and opt-in managed validation sessions.
// It is deliberately separate from helix so applications can choose how and
// where credentials are stored.
//
// Anchor: oauth-architecture
//
// New constructs a client with functional options. OAuth operations use the
// supplied context and return typed errors. WithHTTPClient shallow-clones the
// caller-owned *http.Client, rejects redirects on the clone, and leaves the
// caller's Transport, Jar, and connection ownership unchanged. OAuth does not
// persist credentials, put secrets in URLs, or log token material.
//
// Anchor: oauth-lifecycle
//
// A TokenPair contains an access token, optional refresh token, expiry, scopes,
// and token type. NewRefreshingTokenSource turns a user TokenPair into a
// helix.TokenSource. Token returns the current immutable snapshot and refreshes
// a user credential when it is within the configured skew (one minute by
// default). Refresh operations are generation-keyed and concurrent callers
// share one refresh flight.
//
// Anchor: refresh-rotation-hooks
//
// Twitch may rotate both tokens. CredentialHook runs with the complete rotated
// TokenPair before the new generation is committed.
// The hook should complete durable credential replacement before returning;
// token values must not be exposed in logs.
// A hook failure leaves the source in ErrCredentialCommit state; RetryCommit
// retries the hook and commits the pending pair. Close is idempotent and stops
// future token work.
//
// Anchor: managed-session-close
//
// NewManagedSession validates immediately, then schedules validation every
// hour by default. A non-terminal validation failure can be observed through
// WithValidationErrorHook and is retried at the next hourly tick. Invalid
// sessions terminalize the source. ManagedSession.Close cancels its scheduler
// and closes the source; callers must call Close explicitly, and may call it
// more than once.
//
// Anchor: token-source-contract
//
// The source implements helix.TokenSource.Token(context.Context) and the
// optional helix.RefreshableTokenSource.Refresh(context.Context,
// helix.CredentialSnapshot, helix.RefreshReason). Refresh reasons are limited
// to expiry and unauthorized responses. App and non-refreshable credentials
// are never refreshed.
//
// Anchor: oauth-errors
//
// TransportError, ProtocolError, OAuthError, and DeviceAuthorizationError
// preserve structured operation/status/code information. Use errors.Is and
// errors.As for branching; error strings are not a compatibility contract.
package oauth
