package eventsub

import (
	"testing"
	"time"
)

func TestCallback_runEventCallback_ChannelChatMessageParsesGifFragment(t *testing.T) {
	payload := []byte(`{
		"broadcaster_user_id": "1971641",
		"broadcaster_user_login": "streamer",
		"broadcaster_user_name": "streamer",
		"chatter_user_id": "4145994",
		"chatter_user_login": "viewer32",
		"chatter_user_name": "viewer32",
		"message_id": "cc106a89-1814-919d-454c-f4f2f970aae7",
		"message": {
			"text": "Look at this gif",
			"fragments": [
				{
					"type": "text",
					"text": "Look at this ",
					"cheermote": null,
					"emote": null,
					"mention": null
				},
				{
					"type": "gif",
					"text": "",
					"cheermote": null,
					"emote": null,
					"mention": null,
					"gif": {
						"gif_id": "8a4f9d2c-1e3b-4c5d-9f6a-7b8c9d0e1f2a",
						"url": "https://static-cdn.jtvnw.net/emoticons/v2/gif/8a4f9d2c"
					}
				}
			]
		},
		"message_type": "text",
		"color": "#00FF7F"
	}`)

	events := make(chan ChannelChatMessageEvent, 1)
	callback := callback[WebhookNotificationMetadata]{}
	callback.OnChannelChatMessage(func(event ChannelChatMessageEvent, metadata WebhookNotificationMetadata) {
		events <- event
	})

	err := callback.runEventCallback(
		EventTypeChannelChatMessage,
		"1",
		RawEvent{Event: payload},
		WebhookNotificationMetadata{},
	)
	if err != nil {
		t.Fatalf("run event callback: %v", err)
	}

	captured := receiveEvent(t, events)
	fragments := captured.Message.Fragments
	if len(fragments) != 2 {
		t.Fatalf("fragments: got %d, want 2", len(fragments))
	}

	text := fragments[0]
	if text.Type != MessageFragmentText {
		t.Fatalf("first fragment type: got %q, want %q", text.Type, MessageFragmentText)
	}
	if text.Text != "Look at this " {
		t.Fatalf("first fragment text: got %q, want %q", text.Text, "Look at this ")
	}
	if text.Gif != nil {
		t.Fatalf("first fragment gif: got %+v, want nil", text.Gif)
	}

	gif := fragments[1]
	if gif.Type != MessageFragmentGif {
		t.Fatalf("second fragment type: got %q, want %q", gif.Type, MessageFragmentGif)
	}
	if gif.Gif == nil {
		t.Fatal("second fragment gif: got nil, want parsed gif object")
	}
	if gif.Gif.GifId != "8a4f9d2c-1e3b-4c5d-9f6a-7b8c9d0e1f2a" {
		t.Fatalf("gif id: got %q, want %q", gif.Gif.GifId, "8a4f9d2c-1e3b-4c5d-9f6a-7b8c9d0e1f2a")
	}
	if gif.Gif.URL != "https://static-cdn.jtvnw.net/emoticons/v2/gif/8a4f9d2c" {
		t.Fatalf("gif url: got %q, want %q", gif.Gif.URL, "https://static-cdn.jtvnw.net/emoticons/v2/gif/8a4f9d2c")
	}
}

func TestCallback_runEventCallback_ChannelChatMessageGifOnlyMessage(t *testing.T) {
	payload := []byte(`{
		"broadcaster_user_id": "1971641",
		"broadcaster_user_login": "streamer",
		"broadcaster_user_name": "streamer",
		"chatter_user_id": "4145994",
		"chatter_user_login": "viewer32",
		"chatter_user_name": "viewer32",
		"message_id": "0f1e2d3c-4b5a-6978-8796-a5b4c3d2e1f0",
		"message": {
			"text": "",
			"fragments": [
				{
					"type": "gif",
					"text": "",
					"gif": {
						"gif_id": "gif-42",
						"url": "https://static-cdn.jtvnw.net/gif/42"
					}
				}
			]
		},
		"message_type": "text",
		"color": "#00FF7F"
	}`)

	events := make(chan ChannelChatMessageEvent, 1)
	callback := callback[WebhookNotificationMetadata]{}
	callback.OnChannelChatMessage(func(event ChannelChatMessageEvent, metadata WebhookNotificationMetadata) {
		events <- event
	})

	err := callback.runEventCallback(
		EventTypeChannelChatMessage,
		"1",
		RawEvent{Event: payload},
		WebhookNotificationMetadata{},
	)
	if err != nil {
		t.Fatalf("run event callback: %v", err)
	}

	captured := receiveEvent(t, events)
	fragments := captured.Message.Fragments
	if len(fragments) != 1 {
		t.Fatalf("fragments: got %d, want 1", len(fragments))
	}

	gif := fragments[0]
	if gif.Type != MessageFragmentGif {
		t.Fatalf("fragment type: got %q, want %q", gif.Type, MessageFragmentGif)
	}
	if gif.Gif == nil {
		t.Fatal("fragment gif: got nil, want parsed gif object")
	}
	if gif.Gif.GifId != "gif-42" {
		t.Fatalf("gif id: got %q, want %q", gif.Gif.GifId, "gif-42")
	}
	if gif.Gif.URL != "https://static-cdn.jtvnw.net/gif/42" {
		t.Fatalf("gif url: got %q, want %q", gif.Gif.URL, "https://static-cdn.jtvnw.net/gif/42")
	}
}

func receiveEvent(t *testing.T, events chan ChannelChatMessageEvent) ChannelChatMessageEvent {
	t.Helper()

	select {
	case event := <-events:
		return event
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for chat message handler")
		return ChannelChatMessageEvent{}
	}
}
