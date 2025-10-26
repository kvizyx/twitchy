package eventsub

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/internal/json"
)

type (
	websocketRawMessage struct {
		Metadata websocketRawMessageMetadata `json:"metadata"`
		Payload  json.RawMessage             `json:"payload"`
	}

	websocketRawMessageMetadata struct {
		MessageId           string       `json:"message_id"`
		MessageType         string       `json:"message_type"`
		MessageTimestamp    TimestampUTC `json:"message_timestamp"`
		SubscriptionType    EventType    `json:"subscription_type"`
		SubscriptionVersion string       `json:"subscription_version"`
	}

	// websocketRawEvent is a websocket duplicate for RawEvent to close internal JSON wrapper from user.
	websocketRawEvent struct {
		Subscription json.RawMessage `json:"subscription"`
		Event        json.RawMessage `json:"event"`
	}
)

// handleMessage handles raw websocket message.
//
// Reference: https://dev.twitch.tv/docs/eventsub/handling-websocket-events.
func (ws *Websocket) handleMessage(ctx context.Context, message websocketRawMessage) error {
	var (
		payload  = message.Payload
		metadata = message.Metadata
	)

	switch metadata.MessageType {
	case "session_welcome":
		if err := ws.handleWelcomeMessage(metadata, payload); err != nil {
			return fmt.Errorf("handle welcome message: %w", err)
		}
	case "session_keepalive":
		ws.handleKeepaliveMessage(metadata)
	case "session_reconnect":
		if err := ws.handleReconnectMessage(ctx, metadata, payload); err != nil {
			return fmt.Errorf("handle reconnect message: %w", err)
		}
	case "notification":
		if err := ws.handleNotificationMessage(metadata, payload); err != nil {
			return fmt.Errorf("handle notification message: %w", err)
		}
	case "revocation":
		if err := ws.handleRevocationMessage(metadata, payload); err != nil {
			return fmt.Errorf("handle revocation message: %w", err)
		}
	}

	return nil
}

// handleRevocationMessage handles revocation message that is sent if Twitch revokes an event subscription.
//
// Reference: https://dev.twitch.tv/docs/eventsub/handling-websocket-events/#revocation-message.
func (ws *Websocket) handleRevocationMessage(rawMetadata websocketRawMessageMetadata, rawPayload json.RawMessage) error {
	// TODO: implement me
	return nil
}

// handleWelcomeMessage handles welcome message that is sent when you connect to the server.
//
// Reference: https://dev.twitch.tv/docs/eventsub/handling-websocket-events/#welcome-message.
func (ws *Websocket) handleWelcomeMessage(rawMetadata websocketRawMessageMetadata, rawPayload json.RawMessage) error {
	var welcomePayload WebsocketWelcomePayload

	if err := json.Unmarshal(rawPayload, &welcomePayload); err != nil {
		return fmt.Errorf("unmarshal raw payload: %w", err)
	}

	welcomeMessage := WebsocketWelcomeMessage{
		Metadata: WebsocketWelcomeMetadata{
			MessageId:        rawMetadata.MessageId,
			MessageType:      rawMetadata.MessageType,
			MessageTimestamp: rawMetadata.MessageTimestamp,
		},
		Payload: welcomePayload,
	}

	select {
	case ws.welcome <- struct{}{}:
	default:
	}

	if ws.onWelcome != nil {
		go ws.onWelcome(welcomeMessage)
	}

	return nil
}

// handleKeepaliveMessage handles keepalive message that indicate that the websocket connection is healthy.
//
// Reference: https://dev.twitch.tv/docs/eventsub/handling-websocket-events/#keepalive-message.
func (ws *Websocket) handleKeepaliveMessage(rawMetadata websocketRawMessageMetadata) {
	keepaliveMessage := WebsocketKeepaliveMessage{
		Metadata: WebsocketKeepaliveMetadata{
			MessageId:        rawMetadata.MessageId,
			MessageType:      rawMetadata.MessageType,
			MessageTimestamp: rawMetadata.MessageTimestamp,
		},
		Payload: struct{}{},
	}

	ws.setLastKeepalive(keepaliveMessage.Metadata.MessageTimestamp)

	if ws.onKeepalive != nil {
		go ws.onKeepalive(keepaliveMessage)
	}
}

// handleReconnectMessage handles reconnect message that is sent if the edge server that you are connected to needs to be swapped.
//
// Reference: https://dev.twitch.tv/docs/eventsub/handling-websocket-events/#reconnect-message.
func (ws *Websocket) handleReconnectMessage(
	ctx context.Context,
	rawMetadata websocketRawMessageMetadata,
	rawPayload json.RawMessage,
) error {
	var reconnectPayload WebsocketReconnectPayload

	if err := json.Unmarshal(rawPayload, &reconnectPayload); err != nil {
		return fmt.Errorf("unmarshal raw payload: %w", err)
	}

	reconnectMessage := WebsocketReconnectMessage{
		Metadata: WebsocketReconnectMetadata{
			MessageId:        rawMetadata.MessageId,
			MessageType:      rawMetadata.MessageType,
			MessageTimestamp: rawMetadata.MessageTimestamp,
		},
		Payload: reconnectPayload,
	}

	go func() {
		err := ws.reconnect(ctx, reconnectMessage.Payload.Session.ReconnectURL)
		if err == nil {
			return
		}

		if ws.onError != nil {
			go ws.onError(err)
		}
	}()

	if ws.onReconnect != nil {
		go ws.onReconnect(reconnectMessage)
	}

	return nil
}

// handleNotificationMessage handles notification message that is sent when an event occurs.
//
// Reference: https://dev.twitch.tv/docs/eventsub/handling-websocket-events/#notification-message.
func (ws *Websocket) handleNotificationMessage(rawMetadata websocketRawMessageMetadata, rawPayload json.RawMessage) error {
	var wsRawEvent websocketRawEvent

	if err := json.Unmarshal(rawPayload, &wsRawEvent); err != nil {
		return fmt.Errorf("unmarshal raw payload: %w", err)
	}

	metadata := WebsocketNotificationMetadata{
		MessageId:           rawMetadata.MessageId,
		MessageType:         rawMetadata.MessageType,
		MessageTimestamp:    rawMetadata.MessageTimestamp,
		SubscriptionType:    rawMetadata.SubscriptionType,
		SubscriptionVersion: rawMetadata.SubscriptionVersion,
	}

	rawEvent := RawEvent{
		Subscription: wsRawEvent.Subscription,
		Event:        wsRawEvent.Event,
	}

	if err := ws.callback.runEventCallback(metadata.SubscriptionType, metadata.SubscriptionVersion, rawEvent, metadata); err != nil {
		return fmt.Errorf("run event callback: %w", err)
	}

	return nil
}
