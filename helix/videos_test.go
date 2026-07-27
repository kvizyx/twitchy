package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestVideosGetVideos(t *testing.T) {
	// Given published videos selected by repeated IDs and all documented filters.
	body := task26Body(t, "videos.json")
	transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusOK, body))
	client := task26Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "viewer"})
	period := helix.VideoPeriodWeek
	sortOrder := helix.VideoSortViews
	videoType := helix.VideoTypeArchive

	// When videos are requested.
	result, err := client.Videos.GetVideos(context.Background(), helix.GetVideosRequest{
		IDs:      []string{"video-a", "video-b"},
		Language: task26Ptr("en"),
		Period:   &period,
		Sort:     &sortOrder,
		Type:     &videoType,
		First:    task26Ptr(2),
		After:    task26Ptr("cursor-a"),
	})

	// Then repeated IDs, enums, pagination, and nullable stream IDs decode exactly.
	if err != nil {
		t.Fatal(err)
	}
	task26Contract(t, "get-videos", urlValues("id", "video-a", "id", "video-b", "language", "en", "period", "week", "sort", "views", "type", "archive", "first", "2", "after", "cursor-a"), "", task26Success(body), transport, result.Meta)
	if len(result.Data) != 1 || result.Data[0].StreamID != nil || result.Data[0].Type != helix.VideoTypeArchive || len(result.Data[0].MutedSegments) != 1 {
		t.Fatalf("videos = %#v", result.Data)
	}
}

func TestVideosDeleteVideosDoesNotReplay(t *testing.T) {
	// Given a delete mutation and a transient upstream failure.
	transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusServiceUnavailable, `{"error":"Unavailable","status":503,"message":"try again"}`))
	client := task26Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "owner", Scopes: []helix.AuthorizationScope{helix.ScopeChannelManageVideos}})

	// When multiple videos are deleted.
	_, err := client.Videos.DeleteVideos(context.Background(), helix.DeleteVideosRequest{IDs: []string{"video-a", "video-b"}})

	// Then the mutation is not automatically replayed.
	var apiErr *helix.APIError
	if !errors.As(err, &apiErr) || len(transport.Requests()) != 1 || transport.Requests()[0].Path != "/helix/videos?id=video-a&id=video-b" {
		t.Fatalf("error = %T %v, requests = %#v", err, err, transport.Requests())
	}
}

func TestVideosGetVideosPager(t *testing.T) {
	// Given two pages of videos.
	firstBody := task26Body(t, "videos.json")
	secondBody := `{"data":[],"pagination":{}}`
	transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusOK, firstBody), task26Response(http.StatusOK, secondBody))
	client := task26Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})

	// When the pager advances twice.
	pager, err := client.Videos.GetVideosPager(helix.GetVideosRequest{IDs: []string{"video-a"}})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || !pager.Next(context.Background()) || pager.Next(context.Background()) {
		t.Fatal("unexpected pager state")
	}

	// Then the second request carries the response cursor.
	requests := transport.Requests()
	if len(requests) != 2 || requests[0].Path != "/helix/videos?id=video-a" || requests[1].Path != "/helix/videos?after=cursor-b&id=video-a" {
		t.Fatalf("pager requests = %#v", requests)
	}
}
