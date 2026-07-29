# Twitchy

Twitchy is a comprehensive Twitch client library that provides an easy way to communicate with all major Twitch APIs in
Go applications.

## Packages

Each scope of the Twitch API is divided into separate packages with corresponding names, so use the API name (e.g.,```eventsub```)
as a package name to access the functionality associated with that API (unless it's not a sub-package)

- [x] [EventSub](https://dev.twitch.tv/docs/eventsub) (Websocket client and Webhook handler)
- [x] [Helix](https://dev.twitch.tv/docs/api)
- [x] OAuth lifecycle support (the `oauth` package)

## Helix quick start

Construct a client with functional options, call a typed domain service, and
branch on typed errors. The examples under `examples/helix-*` use local fake
servers, so they never need Twitch credentials or network access.

```go
client, err := helix.New(
	helix.WithStaticToken(helix.Credential{
		AccessToken: "placeholder",
		TokenClass:  helix.TokenClassApp,
	}),
)
if err != nil {
	return err
}

response, err := client.Users.GetUsers(ctx, helix.GetUsersRequest{Logins: []string{"example-user"}})
if err != nil {
	var apiErr *helix.APIError
	if errors.As(err, &apiErr) {
		return fmt.Errorf("Helix request failed with status %d: %w", apiErr.StatusCode(), err)
	}
	return err
}
fmt.Println(len(response.Data), response.Meta.RateLimit().Remaining())
```

`WithHTTPClient` shallow-clones caller-owned clients, `WithTokenSource` accepts
the `TokenSource` contract, and `oauth.NewRefreshingTokenSource` plus
`oauth.NewManagedSession` provide rotation hooks and explicit lifecycle Close.
429 waiting is opt-in, mutations are never replayed, pagers are lazy and
bounded with `helix.WithPageLimit`, and NEW/BETA operations are available only
through `Client.Experimental` without a stable compatibility promise. iCalendar
data is returned as `[]byte`; analytics CSV and clip download URLs are returned
without fetching them.

## OAuth lifecycle

Twitch user access tokens expire in about four hours, and every refresh
rotates the refresh token itself — the old one becomes invalid immediately.
The `oauth` package exists so you don't have to babysit that: it performs the
grant flows, refreshes tokens proactively, notifies you on every rotation so
you can persist the new pair, and keeps the token validated in the background.

```go
hook := func(_ context.Context, pair oauth.TokenPair) error {
	return saveToDatabase(pair) // persist the rotated refresh token
}

source, err := oauth.NewRefreshingTokenSource(client, oauth.TokenPair{
	AccessToken:  "access-token",
	RefreshToken: "refresh-token",
	TokenType:    "bearer",
}, hook)

// Validates the token on a managed hourly interval; Close stops it.
session, err := oauth.NewManagedSession(ctx, source)
defer session.Close()

// Hand the source straight to Helix — every request always uses a valid token.
helixClient, err := helix.New(helix.WithTokenSource(source))
```

See [`examples/helix-oauth`](examples/helix-oauth) for a runnable offline
version. Note that `helix.WithStaticToken` is only suitable for app access
tokens that you refresh yourself — user tokens will eventually stop working
without a `TokenSource`.

## Pagination

List operations expose two forms: a single-page request and a lazy pager.
The pager performs no I/O on construction and fetches the next page only when
you call `Next`, so a dropped `ctx` or an early `break` never wastes requests.

```go
pager, err := client.Streams.GetStreamsPager(
	helix.GetStreamsRequest{},
	helix.WithPageLimit(100), // hard bound on total pages fetched
)
if err != nil {
	return err
}

for pager.Next(ctx) {
	for _, stream := range pager.Page().Data {
		fmt.Println(stream.ID)
	}
}
if err := pager.Err(); err != nil { // distinguish iteration failure from EOF
	return err
}
```

See [`examples/helix-pagination`](examples/helix-pagination) for a runnable
offline version.

## Experimental operations

Twitch ships NEW and BETA endpoints that can change or disappear without
notice. To keep the stable surface trustworthy, these 22 operations are
isolated under `client.Experimental` — using them is an explicit opt-in, and
they carry no compatibility promise across releases.

```go
response, err := client.Experimental.Bits.GetCustomPowerUp(ctx,
	helix.GetCustomPowerUpRequest{BroadcasterID: broadcasterID})
```

See [`examples/helix-experimental`](examples/helix-experimental) for a runnable
offline version. When Twitch promotes an endpoint out of NEW/BETA, it moves to
the corresponding stable service in the next release.

## Multiple users and bots

A bot serving many channels keeps one credential per user. `oauth.Registry`
runs a managed session for every registered user (proactive refresh, hourly
validation, rotation hooks) and implements `helix.CredentialResolver`, so a
single root client can switch context explicitly:

```go
registry, err := oauth.NewRegistry(oauthClient)
if err != nil {
	return err
}
defer registry.Close() // closes every user session

err = registry.AddUser(ctx, broadcasterID, pair, hook,
	helix.Intent("chat"), helix.Intent("eventsub"))

client, err := helix.New(helix.WithCredentialResolver(registry))

// A derived client that always acts as one specific user:
asBroadcaster, err := client.AsUser(broadcasterID)

// A derived client acting as any registered user covering the intents:
asChatBot, err := client.AsIntent(helix.Intent("chat"))
```

Derived clients share the root HTTP transport, perform no network I/O when
created, and keep all pre-network credential checks (token class, scopes,
subject binding). Intents are opaque labels — the registry resolves the first
user (by sorted ID) whose intents cover the request, skipping terminated
sessions.

## Coordinated registries across processes

`oauth.Registry` coordinates refreshes within one process only. When several
replicas share the same credentials, build a `CoordinatedRegistry` with a
`RefreshCoordinator` so only one process rotates a refresh token at a time:

```go
redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
coordinator, err := oauth.NewRedisRefreshCoordinator(redisClient,
	func(userID string) string { return "myapp:twitch:refresh:" + userID })

registry, err := oauth.NewCoordinatedRegistry(oauthClient, coordinator)
defer registry.Close()

loader := func(ctx context.Context, userID string) (oauth.TokenPair, error) {
	return loadPairFromDatabase(ctx, userID) // ExpiresIn = remaining lifetime
}
hook := func(ctx context.Context, pair oauth.TokenPair) error {
	return savePairToDatabase(ctx, pair) // persist the rotated pair
}

err = registry.AddCoordinatedUser(ctx, broadcasterID, loader, hook,
	helix.Intent("chat"))
```

While a process holds the lease it reloads the durable pair, adopts it when
another process already rotated it, refreshes otherwise, commits through the
hook, and only then activates the credential. Waiters take the lease next,
reload, and adopt — one remote refresh per rotation across the fleet. The
loader must return the pair's *remaining* lifetime in `ExpiresIn`. Lease keys
must contain user IDs only, never token material. The Redis implementation
derives a bounded, context-aware client from the supplied one (shared pool,
finite lease I/O timeouts) and never closes the caller's client.

Failure semantics matter because Twitch rotates the refresh token itself:

- A failed durable commit keeps the rotated pair pending; `RetryCommit`
  reacquires the lease, reloads, adopts a pair committed by another process,
  or retries the same hook.
- If the process dies after Twitch rotated the token but before the hook
  persisted it, the credential is unrecoverable and the user must
  reauthorize — `CredentialRotationError` marks this terminal state.
- Caller cancellation and definitive OAuth rejections are safe: the source
  stays usable and the next attempt reloads the durable pair.

## Contributing

You are more than welcome to contribute! Where it's possible, please include unit-tests for any code that is introduced
by your contribution. It's also helpful if you can include usage examples in the documentation.

## License

This library is distributed under the MIT license.
