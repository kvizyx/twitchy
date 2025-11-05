package eventsub

import (
	"net/http"
	"time"
)

// WebsocketOption is an optional setting for Websocket.
type WebsocketOption func(*Websocket)

// WebsocketWithClient sets HTTP client that will be used to connect to the eventsub websocket server.
//
// Default value is http.DefaultClient.
func WebsocketWithClient(client *http.Client) WebsocketOption {
	return func(ws *Websocket) {
		ws.client = client
	}
}

// WebsocketWithServerURL sets the URL that will be used as the URL address of the websocket eventsub server to which the client will connect.
// It's helpful when you want to test your Websocket client with custom server and therefore need to set a custom URL.
// See https://dev.twitch.tv/docs/cli/websocket-event-command for more information on how to test Websocket client.
//
// Default value is "wss://eventsub.wss.twitch.tv/ws".
func WebsocketWithServerURL(serverURL string) WebsocketOption {
	return func(ws *Websocket) {
		ws.serverURL = serverURL
	}
}

// WebsocketWithKeepalive sets timeout in seconds for Twitch to send keepalive messages to ensure that connection to server is healthy.
// See https://dev.twitch.tv/docs/eventsub/handling-websocket-events/#keepalive-message for more information on keepalive timeout.
//
// Default value is 300.
func WebsocketWithKeepalive(seconds uint) WebsocketOption {
	return func(ws *Websocket) {
		ws.keepaliveSeconds = seconds
	}
}

// WebsocketWithNoRetry specifies that Websocket will not do retries when connecting to the eventsub server fails.
// Usually you don't want to use this option.
//
// Default retry attempts number is 5.
func WebsocketWithNoRetry() WebsocketOption {
	return func(ws *Websocket) {
		ws.retryAttempts = 0
	}
}

// WebsocketWithRetryAttempts sets number of retry attempts to perform when connecting to the eventsub server fails.
//
// Default value is 5.
func WebsocketWithRetryAttempts(attempts uint) WebsocketOption {
	return func(ws *Websocket) {
		ws.retryAttempts = attempts
	}
}

// WebsocketWithRetryDelay sets delay between retries when connecting to the eventsub server fails.
//
// Default value is 1 second.
func WebsocketWithRetryDelay(delay time.Duration) WebsocketOption {
	return func(ws *Websocket) {
		ws.retryDelay = delay
	}
}
