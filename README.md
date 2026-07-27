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

## Contributing

You are more than welcome to contribute! Where it's possible, please include unit-tests for any code that is introduced
by your contribution. It's also helpful if you can include usage examples in the documentation.

## License

This library is distributed under the MIT license.
