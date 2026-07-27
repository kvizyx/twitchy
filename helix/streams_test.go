package helix_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestStreamsGetStreamKey_requiresUserScopeAndRedactsTokenFromErrors(t *testing.T) {
	token := "secret-stream-token"
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: http.StatusForbidden, Body: `{"error":"Forbidden","status":403,"message":"secret-stream-token"}`})
	client := task24Client(t, transport, helix.Credential{AccessToken: token, ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "broadcaster-1", Scopes: []helix.AuthorizationScope{helix.ScopeChannelReadStreamKey}})

	_, err := client.Streams.GetStreamKey(context.Background(), helix.GetStreamKeyRequest{BroadcasterID: "broadcaster-1"})
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) || strings.Contains(err.Error(), token) {
		t.Fatalf("stream key error = %T %v", err, err)
	}
	if len(transport.Requests()) != 1 || transport.Requests()[0].Header.Get("Authorization") != "Bearer "+token {
		t.Fatalf("stream key request = %#v", transport.Requests())
	}
}

func TestStreamsGetStreamKey_rejectsMismatchedBroadcasterBeforeNetwork(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper()
	client := task24Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "different-user", Scopes: []helix.AuthorizationScope{helix.ScopeChannelReadStreamKey}})

	_, err := client.Streams.GetStreamKey(context.Background(), helix.GetStreamKeyRequest{BroadcasterID: "broadcaster-1"})
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) || len(transport.Requests()) != 0 {
		t.Fatalf("stream key mismatch = %T %v requests=%d", err, err, len(transport.Requests()))
	}
}

func TestStreamsGetStreamsPager_preservesBackwardCursorDirection(t *testing.T) {
	body := `{"data":[{"id":"stream-1","user_id":"user-1","user_login":"channel","user_name":"Channel","game_id":"game-1","game_name":"Game","type":"live","title":"Live","tags":["English"],"viewer_count":42,"started_at":"2026-07-27T18:00:00Z","language":"en","thumbnail_url":"thumb/{width}x{height}","tag_ids":[],"is_mature":false}],"pagination":{"cursor":"older"}}`
	second := `{"data":[],"pagination":{}}`
	transport := testkit.NewRecordingRoundTripper(task24Response(http.StatusOK, body), task24Response(http.StatusOK, second))
	client := task24Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	before := "before-cursor"
	pager, err := client.Streams.GetStreamsPager(helix.GetStreamsRequest{Before: &before, UserID: []string{"user-1", "user-2"}, Language: []string{"en", "de"}})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) {
		t.Fatal("first streams page failed")
	}
	if !pager.Next(context.Background()) {
		t.Fatal("second streams page failed")
	}
	if pager.Next(context.Background()) || pager.Err() != nil {
		t.Fatalf("streams pager state: page=%#v err=%v", pager.Page(), pager.Err())
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[0].Path != "/helix/streams?before=before-cursor&language=en&language=de&user_id=user-1&user_id=user-2" || requests[1].Path != "/helix/streams?before=older&language=en&language=de&user_id=user-1&user_id=user-2" {
		t.Fatalf("streams pager requests = %#v", requests)
	}
}

func TestStreamsGetFollowedStreams_requiresSubjectAndScope(t *testing.T) {
	body := `{"data":[{"id":"stream-1","user_id":"user-1","user_login":"channel","user_name":"Channel","game_id":"game-1","game_name":"Game","type":"live","title":"Live","viewer_count":42,"started_at":"2026-07-27T18:00:00Z","language":"en","thumbnail_url":"thumb","tag_ids":[],"tags":[]}],"pagination":{}}`
	transport := testkit.NewRecordingRoundTripper(task24Response(http.StatusOK, body))
	client := task24Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "user-1", Scopes: []helix.AuthorizationScope{helix.ScopeUserReadFollows}})

	result, err := client.Streams.GetFollowedStreams(context.Background(), helix.GetFollowedStreamsRequest{UserID: "user-1", First: task24Int(1), After: task24String("next")})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 || transport.Requests()[0].Path != "/helix/streams/followed?after=next&first=1&user_id=user-1" {
		t.Fatalf("followed streams = %#v requests=%#v", result.Data, transport.Requests())
	}
}

func TestStreamsGetFollowedStreamsPager_forwardsAfterCursor(t *testing.T) {
	body := `{"data":[],"pagination":{"cursor":"next"}}`
	transport := testkit.NewRecordingRoundTripper(task24Response(http.StatusOK, body), task24Response(http.StatusOK, `{"data":[],"pagination":{}}`))
	client := task24Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "user-1", Scopes: []helix.AuthorizationScope{helix.ScopeUserReadFollows}})
	pager, err := client.Streams.GetFollowedStreamsPager(helix.GetFollowedStreamsRequest{UserID: "user-1"})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || !pager.Next(context.Background()) || pager.Next(context.Background()) {
		t.Fatalf("followed streams pager state: page=%#v err=%v", pager.Page(), pager.Err())
	}
	if requests := transport.Requests(); len(requests) != 2 || requests[1].Path != "/helix/streams/followed?after=next&user_id=user-1" {
		t.Fatalf("followed streams pager requests = %#v", requests)
	}
}

func TestStreamsMarkers_validateSelectorBeforeNetwork(t *testing.T) {
	for _, request := range []helix.GetStreamMarkersRequest{{}, {UserID: task24String("user-1"), VideoID: task24String("video-1")}} {
		transport := testkit.NewRecordingRoundTripper()
		client := task24Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, Scopes: []helix.AuthorizationScope{helix.ScopeUserReadBroadcast}})

		_, err := client.Streams.GetStreamMarkers(context.Background(), request)
		if err == nil || len(transport.Requests()) != 0 {
			t.Fatalf("marker selector request=%#v error=%v network=%d", request, err, len(transport.Requests()))
		}
	}
}

func TestStreamsMarkers_createAndGetPreserveNestedWireFields(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(
		task24Response(http.StatusOK, `{"data":[{"id":"marker-1","created_at":"2026-07-27T18:00:00Z","description":"hello","position_seconds":42}]}`),
		task24Response(http.StatusOK, `{"data":[{"user_id":"user-1","user_name":"User","user_login":"user","videos":[{"video_id":"video-1","markers":[{"id":"marker-1","created_at":"2026-07-27T18:00:00Z","description":"hello","position_seconds":42,"URL":"https://twitch.tv/marker"}]}]}],"pagination":{"cursor":"next"}}`),
	)
	client := task24Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "user-1", Scopes: []helix.AuthorizationScope{helix.ScopeChannelManageBroadcast, helix.ScopeUserReadBroadcast}})

	created, err := client.Streams.CreateStreamMarker(context.Background(), helix.CreateStreamMarkerRequest{UserID: "user-1", Description: task24String("hello")})
	if err != nil {
		t.Fatal(err)
	}
	markers, err := client.Streams.GetStreamMarkers(context.Background(), helix.GetStreamMarkersRequest{VideoID: task24String("video-1"), First: task24String("5")})
	if err != nil {
		t.Fatal(err)
	}
	requests := transport.Requests()
	if len(requests) != 2 || string(requests[0].Body) != `{"user_id":"user-1","description":"hello"}` || requests[1].Path != "/helix/streams/markers?first=5&video_id=video-1" {
		t.Fatalf("marker requests = %#v", requests)
	}
	if len(created.Data) != 1 || created.Data[0].PositionSeconds != 42 || markers.Data[0].Videos[0].Markers[0].URL == "" || markers.Pagination.Cursor() != "next" {
		t.Fatalf("marker results = %#v %#v", created.Data, markers.Data)
	}
}
