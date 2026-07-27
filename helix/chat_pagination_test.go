package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestChatGetChattersPagerFollowsCursor(t *testing.T) {
	// Given two Chatters pages linked by the response cursor.
	transportResponses := []testkit.RoundTripResponse{
		{StatusCode: http.StatusOK, Body: `{"data":[{"user_id":"1"}],"pagination":{"cursor":"next"}}`},
		{StatusCode: http.StatusOK, Body: `{"data":[{"user_id":"2"}]}`},
	}
	client, transport := chatClientForCredential(t, chatCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorReadChatters), transportResponses...)
	pager, err := client.Chat.GetChattersPager(helix.GetChattersRequest{BroadcasterID: "1234", ModeratorID: "5678"}, helix.WithPageLimit(2))
	if err != nil {
		t.Fatal(err)
	}

	// When both pages are consumed.
	pageCount := 0
	for pager.Next(context.Background()) {
		pageCount++
	}

	// Then the pager follows after and stops at the empty cursor.
	if pageCount != 2 || pager.Err() != nil {
		t.Fatalf("pages = %d, err = %v", pageCount, pager.Err())
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[0].Path != "/helix/chat/chatters?broadcaster_id=1234&moderator_id=5678" || requests[1].Path != "/helix/chat/chatters?after=next&broadcaster_id=1234&moderator_id=5678" {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestChatGetUserEmotesPagerPreservesInitialCursor(t *testing.T) {
	// Given a user-emote request that starts from an existing cursor.
	client, transport := chatClientForCredential(t, chatCredential(helix.TokenClassUser, "5678", helix.ScopeUserReadEmotes),
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[]}`},
	)
	pager, err := client.Chat.GetUserEmotesPager(helix.GetUserEmotesRequest{UserID: "5678", After: chatString("initial")})
	if err != nil {
		t.Fatal(err)
	}

	// When the first page is requested.
	if !pager.Next(context.Background()) {
		t.Fatalf("pager.Next() = false, err = %v", pager.Err())
	}

	// Then the initial cursor is sent on the wire without mutating the caller's request.
	requests := transport.Requests()
	if len(requests) != 1 || requests[0].Path != "/helix/chat/emotes/user?after=initial&user_id=5678" {
		t.Fatalf("requests = %#v", requests)
	}
}
