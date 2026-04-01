package main

import (
	"context"
	"log"
	"time"

	"github.com/kvizyx/twitchy/eventsub"
)

func main() {
	// Create new eventsub instance (it's like a manager for webhook handler and websocket client instances).
	// You need only one instance of eventsub per application.
	manager := eventsub.New()

	// Create new eventsub websocket client.
	websocket := manager.Websocket()

	websocket.OnWelcome(func(message eventsub.WebsocketWelcomeMessage) {
		// Subscribe to events on welcome message (unless you don't use conduits).
		// We will use channel follow event in this example.
	})

	websocket.OnDisconnect(func() {
		// Disconnection callback may be helpful if you want to do some clean-up when client disconnects from the eventsub server.
		log.Println("websocket client has been disconnected")
	})

	websocket.OnChannelFollow(func(event eventsub.ChannelFollowEvent, metadata eventsub.WebsocketNotificationMetadata) {
		// Here you can handle event from eventsub. Each event handler is running in a separate goroutine so you will not block
		// other events even if you do some long-running tasks in handler.
		log.Printf("%s is now following %s!\n", event.UserName, event.BroadcasterUserName)
	})

	// Connect to the eventsub server in infinite loop to imitate reconnection logic if client is disconnected for some reason.
	for {
		if err := websocket.Connect(context.TODO()); err != nil {
			log.Printf("websocket client is failed to connect or listen: %s\n", err)
		}

		// Sleep for one second to prevent spamming with new connections.
		// Use some kind of retry mechanism here in real applications.
		time.Sleep(1 * time.Second)
	}
}
