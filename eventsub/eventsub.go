package eventsub

import (
	"context"
	"fmt"
	"time"

	"github.com/kvizyx/twitchy/eventsub/messagetracker"
	"github.com/kvizyx/twitchy/internal/json"
)

// EventSub is a Twitch eventsub module.
//
// Reference: https://dev.twitch.tv/docs/eventsub.
type EventSub struct {
	eventTracker messagetracker.MessageTracker
	unmarshal    json.UnMarshaller
}

func New(options ...Option) EventSub {
	var es EventSub

	for _, option := range options {
		option(&es)
	}

	if es.unmarshal != nil {
		json.Unmarshal = es.unmarshal
	}

	return es
}

// Webhook returns new eventsub Webhook handler that implements http.Handler.
func (es *EventSub) Webhook(secret []byte, verifySignature bool) (*Webhook, error) {
	return newWebhook(secret, es.eventTracker, verifySignature)
}

// Websocket returns new eventsub Websocket client.
func (es *EventSub) Websocket(options ...WebsocketOption) *Websocket {
	return newWebsocket(es.eventTracker, options...)
}

// isExpiredMessage returns does eventsub message with provided timestamp is too old (expired) to process or not.
//
// According to the Twitch's documentation we should not process messages that are older than 10 minutes from the moment
// they were create to guard against replay attacks.
func isExpiredMessage(messageTimestamp TimestampUTC) bool {
	const expirationSeconds = 600

	now, _ := timestampUTCFromString(time.Now().String())

	return now.Sub(messageTimestamp.Time).Seconds() > expirationSeconds
}

// isSafeMessage validates eventsub message with it's metadata to ensure that it's not expired and not
// processed yet so it's safe to process it now.
func isSafeMessage[Metadata any](
	ctx context.Context,
	onDuplicate func(Metadata),
	eventTracker messagetracker.MessageTracker,
	metadata Metadata,
	messageID string,
	messageTimestamp TimestampUTC,
) (error, bool) {
	if isExpiredMessage(messageTimestamp) {
		return nil, false
	}

	if eventTracker != nil {
		isDuplicate, err := eventTracker.Track(ctx, messageID)
		if err != nil {
			return fmt.Errorf("track: %w", err), false
		}

		if isDuplicate {
			if onDuplicate != nil {
				go onDuplicate(metadata)
			}

			return nil, false
		}
	}

	return nil, true
}
