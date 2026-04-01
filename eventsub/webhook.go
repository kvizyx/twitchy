package eventsub

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"unicode"

	"github.com/kvizyx/twitchy/eventsub/messagetracker"
)

// ErrInvalidWebhookSecret indicates that your webhook secret doesn't match with Twitch's webhook secret requirements.
//
// See note section: https://dev.twitch.tv/docs/eventsub/handling-webhook-events/#verifying-the-event-message.
var ErrInvalidWebhookSecret = errors.New("webhook secret is not valid")

// Webhook is an EventSub webhook HTTP handler.
type Webhook struct {
	eventTracker              messagetracker.MessageTracker
	secret                    []byte
	withSignatureVerification bool

	onRevocation   func(WebhookRevocationNotification)
	onVerification func(WebhookCallbackVerificationNotification)

	callback[WebhookNotificationMetadata]
}

var _ http.Handler = (*Webhook)(nil)

func newWebhook(secret []byte, eventTracker messagetracker.MessageTracker, verifySignature bool) (*Webhook, error) {
	if !isValidWebhookSecret(secret) {
		return nil, ErrInvalidWebhookSecret
	}

	return &Webhook{
		eventTracker:              eventTracker,
		secret:                    secret,
		withSignatureVerification: verifySignature,
	}, nil
}

func (wh *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	metadata, err := getWebhookNotificationMetadata(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if wh.withSignatureVerification {
		if len(metadata.MessageSignature) == 0 {
			http.Error(w, "missing signature header", http.StatusBadRequest)
			return
		}

		if !wh.isValidSignature(metadata.MessageSignature, body, metadata.MessageId, metadata.MessageTimestamp) {
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}
	}

	err, safe := isSafeMessage(
		r.Context(),
		wh.onDuplicate,
		wh.eventTracker,
		metadata,
		metadata.MessageId,
		metadata.MessageTimestamp,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !safe {
		http.Error(w, "message is not safe to process", http.StatusBadRequest)
		return
	}

	wh.handleNotification(w, metadata, body)
}

// isValidSignature validates that provided request signature is valid based on the request body, message id and timestamp.
//
// Reference: https://dev.twitch.tv/docs/eventsub/handling-webhook-events/#verifying-the-event-message.
func (wh *Webhook) isValidSignature(signature string, body []byte, messageId string, messageTimestamp TimestampUTC) bool {
	mac := hmac.New(sha256.New, wh.secret)

	mac.Write([]byte(messageId))
	mac.Write([]byte(messageTimestamp.String()))
	mac.Write(body)

	computedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(computedSignature), []byte(signature))
}

// isValidWebhookSecret validates that secret meets Twitch's webhook secret requirements.
func isValidWebhookSecret(secret []byte) bool {
	if len(secret) < 10 || len(secret) > 100 {
		return false
	}

	// Check that secret contains only ASCII characters.
	for i := range len(secret) {
		if secret[i] > unicode.MaxASCII {
			return false
		}
	}

	return true
}

// getWebhookNotificationMetadata extracts Twitch-specific headers from webhook request and compose them as WebhookNotificationMetadata.
func getWebhookNotificationMetadata(header http.Header) (WebhookNotificationMetadata, error) {
	var (
		rawMessageRetry     = header.Get("Twitch-Eventsub-Message-Retry")
		rawMessageTimestamp = header.Get("Twitch-Eventsub-Message-Timestamp")
	)

	messageRetry, err := strconv.Atoi(rawMessageRetry)
	if err != nil {
		return WebhookNotificationMetadata{}, fmt.Errorf("parse retry header: %w", err)
	}

	messageTimestamp, err := timestampUTCFromString(rawMessageTimestamp)
	if err != nil {
		return WebhookNotificationMetadata{}, fmt.Errorf("parse message timestamp header: %w", err)
	}

	subscriptionType := header.Get("Twitch-Eventsub-Subscription-Type")

	return WebhookNotificationMetadata{
		MessageId:           header.Get("Twitch-Eventsub-Message-Id"),
		MessageRetry:        messageRetry,
		MessageType:         header.Get("Twitch-Eventsub-Message-Type"),
		MessageSignature:    header.Get("Twitch-Eventsub-Message-Signature"),
		MessageTimestamp:    messageTimestamp,
		SubscriptionType:    EventType(subscriptionType),
		SubscriptionVersion: header.Get("Twitch-Eventsub-Subscription-Version"),
	}, nil
}
