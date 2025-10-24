# Twitchy

Twitchy is a comprehensive Twitch client library that provides an easy way to communicate with all major Twitch APIs in
Go applications.

## Packages

Each scope of the Twitch API is divided into separate packages with corresponding names, so use the API name (e.g.,```eventsub```)
as a package name to access the functionality associated with that API (unless it's not a sub-package)

- [x] [EventSub](https://dev.twitch.tv/docs/eventsub) (Websocket client and Webhook handler)
- [ ] [Helix](https://dev.twitch.tv/docs/api)

## Contributing

You are more than welcome to contribute! Where it's possible, please include unit-tests for any code that is introduced
by your contribution. It's also helpful if you can include usage examples in the documentation.

## License

This library is distributed under the MIT license.
