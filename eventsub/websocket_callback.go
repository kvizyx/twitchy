package eventsub

// OnWelcome invokes when eventsub sends welcome message to let you subscribe to events.
//
// Reference: https://dev.twitch.tv/docs/eventsub/handling-websocket-events/#welcome-message.
func (ws *Websocket) OnWelcome(onWelcome func(WebsocketWelcomeMessage)) {
	ws.onWelcome = onWelcome
}

// OnKeepalive invokes when eventsub sends keepalive message which indicates that the websocket connection is healthy.
//
// Reference: https://dev.twitch.tv/docs/eventsub/handling-websocket-events/#keepalive-message.
func (ws *Websocket) OnKeepalive(onKeepalive func(WebsocketKeepaliveMessage)) {
	ws.onKeepalive = onKeepalive
}

// OnPing invokes when eventsub sends websocket ping frame (ping-pong cycle).
//
// Reference: https://dev.twitch.tv/docs/eventsub/handling-websocket-events/#ping-message.
func (ws *Websocket) OnPing(onPing func()) {
	ws.onPing = onPing
}

// OnReconnect invokes when eventsub sends reconnect message which indicates that the edge server that the client is connected
// to needs to be swapped.
//
// Reference: https://dev.twitch.tv/docs/eventsub/handling-websocket-events/#reconnect-message.
func (ws *Websocket) OnReconnect(onReconnect func(WebsocketReconnectMessage)) {
	ws.onReconnect = onReconnect
}

// OnError invokes when some internal error occurs that cannot be returned at start (e.g. reconnection error).
func (ws *Websocket) OnError(onError func(error)) {
	ws.onError = onError
}

// OnDisconnect invokes when client instance is being disconnected from eventsub server.
func (ws *Websocket) OnDisconnect(onDisconnect func()) {
	ws.onDisconnect = onDisconnect
}
