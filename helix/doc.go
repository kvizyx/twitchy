// Package helix provides a concurrent-safe, typed client for the Twitch Helix
// API. Stable operations are grouped by domain service on Client; NEW and BETA
// operations are grouped under Client.Experimental.
//
// Anchor: architecture
//
// A Client is built with New and functional options. Options are applied in
// order, and the last option for a setting wins. The default endpoint is
// Twitch's Helix endpoint. Every operation accepts a context.Context and
// returns a Response[T] or a typed error. Service values are initialized by
// the constructor and are safe to share between callers when their contexts
// are independent.
//
// Anchor: custom-http-client
//
// WithHTTPClient accepts a caller-owned *http.Client. Twitchy makes a shallow
// clone of that client and does not mutate the caller's client, Transport, Jar,
// or connection pool. The clone rejects redirects. The caller remains the
// owner of the Transport and Jar and must not expect the Helix client to close
// borrowed idle connections.
//
// Anchor: token-source
//
// WithStaticToken installs a copied Credential for simple app, user, or
// extension calls. WithTokenSource accepts the TokenSource contract:
// Token(context.Context) (CredentialSnapshot, error). A source may also
// implement RefreshableTokenSource, whose Refresh method receives the current
// snapshot and one of the two RefreshReason values. Snapshots are immutable
// views; Scopes returns a copy. Custom sources may read AccessToken from a
// snapshot, but callers must not log or persist it through this package.
//
// Anchor: response-and-errors
//
// Successful calls return Response[T], including endpoint data, pagination,
// and ResponseMeta. ResponseMeta contains status, request ID, attempts, a
// defensive header copy, and parsed rate-limit metadata. TransportError,
// ProtocolError, APIError, AuthError, and RateLimitError preserve their typed
// operation/status/metadata accessors through wrapping; use errors.As rather
// than matching error strings.
//
// Anchor: rate-limit-policy
//
// RateLimitPolicy is opt-in. WithRateLimitPolicy(RateLimitPolicy{Wait: true,
// MaxWait: ...}) permits cancellable waiting for a bucket-exhaustion 429 when
// the endpoint is marked bucket-waitable and its rate metadata is valid. The
// wait is bounded by MaxWait. The default does not wait, and endpoint-specific
// cooldowns are returned as errors instead of being treated as bucket waits.
//
// Anchor: replay-policy
//
// Replay is bounded and manifest-driven. At most one replay is attempted for a
// given retry cause within the transport attempt budget, and a request body is
// replayed only when it can be reconstructed. Mutating operations are replayed
// only after a 401, which Twitch sends before applying the request; other
// transient failures never replay mutations. A 401 is not a blanket retry:
// refresh is used only for a refreshable user credential or for an app
// credential backed by a RefreshableTokenSource.
//
// Anchor: pager-semantics
//
// A normal endpoint method performs one request. Its Pager companion is lazy:
// construction performs no I/O, and each Next call fetches at most one page.
// Page returns the last successful page and Err returns a stable terminal
// error. WithPageLimit bounds a pager (the default is 100 pages and the maximum
// is 10000); the cap is reported after the final successful page. Cursor cycles
// and context cancellation stop the pager without discarding the last page.
//
// Anchor: experimental-compatibility
//
// Client.Experimental is the only access path for operations labeled NEW or
// BETA. Stable services contain stable operations only. Experimental operations
// are exposed for early adoption, but carry no stable compatibility promise:
// their request, response, behavior, and availability may change with the
// upstream Twitch surface.
//
// Anchor: special-formats
//
// Special response formats remain typed and are not fetched or interpreted by
// the client. GetChannelICalendar returns Response[GetChannelICalendarData],
// where the data is a copied []byte iCalendar document. Analytics methods
// return CSV URLs as data; clip download methods return clip URLs as data. The
// client does not download those URLs, parse external files, or persist them.
package helix
