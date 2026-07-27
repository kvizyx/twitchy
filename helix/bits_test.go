package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestBitsLeaderboardAndCheermotes_decodeData(t *testing.T) {
	// Given user-token leaderboard and app-token-compatible Cheermote responses.
	transport := testkit.NewRecordingRoundTripper(
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[{"user_id":"u-1","user_login":"viewer","user_name":"Viewer","rank":1,"score":12543}],"date_range":{"started_at":"2026-07-01T00:00:00Z","ended_at":"2026-07-08T00:00:00Z"},"total":1}`},
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[{"prefix":"Cheer","tiers":[{"min_bits":1,"id":"1","color":"#979797","images":{"dark":{"animated":{"1":"https://cdn.test/cheer.gif"}},"light":{}},"can_cheer":true,"show_in_bits_card":true}],"type":"global_first_party","order":1,"last_updated":"2026-07-01T00:00:00Z","is_charitable":false}]}`},
	)
	client := newTask14Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, Scopes: []helix.AuthorizationScope{helix.ScopeBitsRead}})

	// When both Bits operations are called.
	startedAt := "2026-07-01T00:00:00Z"
	period := "week"
	count := 1
	leaderboard, err := client.Bits.GetBitsLeaderboard(context.Background(), helix.GetBitsLeaderboardRequest{Count: &count, Period: &period, StartedAt: &startedAt})
	if err != nil {
		t.Fatal(err)
	}
	cheermotes, err := client.Bits.GetCheermotes(context.Background(), helix.GetCheermotesRequest{BroadcasterID: "broadcaster-1"})

	// Then nested wire data and exact query fields are decoded.
	if err != nil || len(leaderboard.Data) != 1 || leaderboard.Data[0].Score != 12543 || cheermotes.Data[0].Tiers[0].Images["dark"]["animated"]["1"] == "" {
		t.Fatalf("Bits results = %#v, %v", leaderboard, err)
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[0].Path != "/helix/bits/leaderboard?count=1&period=week&started_at=2026-07-01T00%3A00%3A00Z" || requests[1].Path != "/helix/bits/cheermotes?broadcaster_id=broadcaster-1" {
		t.Fatalf("Bits requests = %#v", requests)
	}
}

func TestBitsExtensionTransactions_usesAppTokenAndPager(t *testing.T) {
	// Given two transaction pages and an app access token.
	transport := testkit.NewRecordingRoundTripper(
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[{"id":"tx-1","timestamp":"2026-07-01T00:00:00Z","broadcaster_id":"b-1","broadcaster_login":"channel","broadcaster_name":"Channel","user_id":"u-1","user_login":"viewer","user_name":"Viewer","product_type":"BITS_IN_EXTENSION","product_data":{"sku":"sku-1","domain":"twitch.ext.ext-1","cost":{"amount":100,"type":"bits"},"inDevelopment":false,"displayName":"Product","expiration":"","broadcast":false}}],"pagination":{"cursor":"next"}}`},
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[{"id":"tx-2","timestamp":"2026-07-02T00:00:00Z","broadcaster_id":"b-1","broadcaster_login":"channel","broadcaster_name":"Channel","user_id":"u-2","user_login":"viewer2","user_name":"Viewer 2","product_type":"BITS_IN_EXTENSION","product_data":{"sku":"sku-2","domain":"twitch.ext.ext-1","cost":{"amount":200,"type":"bits"},"inDevelopment":false,"displayName":"Product 2","expiration":"","broadcast":true}}]}`},
	)
	client := newTask14Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})

	// When both transaction pages are consumed.
	pager, err := client.Bits.GetExtensionTransactionsPager(helix.GetExtensionTransactionsRequest{ExtensionID: "ext-1", IDs: []string{"tx-1"}})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || pager.Page().Data[0].ProductData.Cost.Amount != 100 || !pager.Next(context.Background()) || pager.Page().Data[0].ProductData.Broadcast != true {
		t.Fatalf("transaction pages missing")
	}

	// Then app authentication and forward pagination are visible on the wire.
	requests := transport.Requests()
	if len(requests) != 2 || requests[0].Header.Get("Authorization") != "Bearer app-token" || requests[0].Path != "/helix/extensions/transactions?extension_id=ext-1&id=tx-1" || requests[1].Path != "/helix/extensions/transactions?after=next&extension_id=ext-1&id=tx-1" {
		t.Fatalf("transaction requests = %#v", requests)
	}
}

func TestBitsExtensionTransactions_rejectsUserTokenBeforeNetwork(t *testing.T) {
	// Given a user token for the app-token-only transaction endpoint.
	transport := testkit.NewRecordingRoundTripper()
	client := newTask14Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser})

	// When transactions are requested.
	_, err := client.Bits.GetExtensionTransactions(context.Background(), helix.GetExtensionTransactionsRequest{ExtensionID: "ext-1"})

	// Then auth fails without network access.
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) || len(transport.Requests()) != 0 {
		t.Fatalf("error = %T %v, requests = %d", err, err, len(transport.Requests()))
	}
}
