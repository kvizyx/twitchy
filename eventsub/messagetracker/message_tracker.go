package messagetracker

import (
	"context"
	"time"
)

// SafeMessageTTL is a safe TTL for tracked messages in the storage, allowing them to be treated as duplicates, since
// messages from event-sub server with an earlier creation timestamp will not be processed by eventsub.EventSub at all.
//
// In almost all use-cases this is what you want to use in your custom MessageTracker implementation (if needed).
// It's also used by default in ready-made MessageTracker implementations.
const SafeMessageTTL = 10 * time.Minute

// MessageTracker tracks messages sent from event-sub server to avoid handling duplicate messages multiple times.
type MessageTracker interface {
	// Track begins tracking of the message with the specified id and returns a result indicating whether the message is
	// a duplicate or not.
	Track(ctx context.Context, messageID string) (bool, error)
}
