package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestSearchCategoriesAndChannels_preserveExactQueryFields(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(
		task24Response(http.StatusOK, `{"data":[{"id":"game-1","name":"Game","box_art_url":"https://example.test/box.jpg"}],"pagination":{"cursor":"categories-next"}}`),
		task24Response(http.StatusOK, `{"data":[{"broadcaster_language":"en","broadcaster_login":"channel","display_name":"Channel","game_id":"game-1","game_name":"Game","id":"user-1","is_live":true,"tag_ids":[],"tags":["English"],"thumbnail_url":"https://example.test/profile.jpg","title":"Live","started_at":"2026-07-27T18:00:00Z"}],"pagination":{"cursor":"channels-next"}}`),
	)
	client := task24Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})

	categories, err := client.Search.SearchCategories(context.Background(), helix.SearchCategoriesRequest{Query: "love computer", First: task24Int(10), After: task24String("old")})
	if err != nil {
		t.Fatal(err)
	}
	channels, err := client.Search.SearchChannels(context.Background(), helix.SearchChannelsRequest{Query: "channel", LiveOnly: task24Bool(true), First: task24Int(10), After: task24String("old")})
	if err != nil {
		t.Fatal(err)
	}
	requests := transport.Requests()
	if requests[0].Path != "/helix/search/categories?after=old&first=10&query=love+computer" || requests[1].Path != "/helix/search/channels?after=old&first=10&live_only=true&query=channel" {
		t.Fatalf("search requests = %#v", requests)
	}
	if categories.Data[0].BoxArtURL == "" || categories.Pagination.Cursor() != "categories-next" || channels.Data[0].StartedAt == "" || channels.Pagination.Cursor() != "channels-next" {
		t.Fatalf("search results = %#v %#v", categories.Data, channels.Data)
	}
}

func TestSearchCategoriesPager_forwardsAfterCursor(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(
		task24Response(http.StatusOK, `{"data":[{"id":"game-1","name":"Game","box_art_url":"box"}],"pagination":{"cursor":"next"}}`),
		task24Response(http.StatusOK, `{"data":[{"id":"game-2","name":"Other","box_art_url":"box2"}],"pagination":{}}`),
	)
	client := task24Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	pager, err := client.Search.SearchCategoriesPager(helix.SearchCategoriesRequest{Query: "game"})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || !pager.Next(context.Background()) || pager.Next(context.Background()) || pager.Err() != nil {
		t.Fatalf("search pager state: page=%#v err=%v", pager.Page(), pager.Err())
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[1].Path != "/helix/search/categories?after=next&query=game" {
		t.Fatalf("search pager requests = %#v", requests)
	}
}
