package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kvizyx/twitchy/eventsub"
)

func main() {
	// Create new eventsub instance (it's like a manager for webhook handler and websocket client instances).
	// You need only one instance of eventsub per application.
	manager := eventsub.New()

	// Your webhook secret.
	secret := os.Getenv("WEBHOOK_SECRET")

	// Create new eventsub webhook HTTP handler.
	webhook, err := manager.Webhook([]byte(secret), true)
	if err != nil {
		log.Printf("failed to create webhook handler: %s\n", err)
		return
	}

	webhook.OnChannelFollow(func(event eventsub.ChannelFollowEvent, metadata eventsub.WebhookNotificationMetadata) {
		// Here you can handle event from eventsub. Each event handler is running in a separate goroutine so you will not block
		// other events even if you do some long-running tasks in handler.
		log.Printf("%s is now following %s!\n", event.UserName, event.BroadcasterUserName)
	})

	// Webhook is just an HTTP handler, so you can use it with any server or framework that accepts http.Handler interface.
	if err = http.ListenAndServe(":8080", webhook); err != nil {
		log.Printf("http server is failed to listen: %s\n", err)
	}
}
