package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestClipsCreateClip(t *testing.T) {
	fixture := testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query: urlValues("broadcaster_id", "141981764", "duration", "45.5", "title", "Task 17 clip"),
			Headers: http.Header{
				"Authorization": {"Bearer user-token"},
				"Client-Id":     {"client-id"},
			},
		},
		Response: testkit.ContractResponse{
			Status:  http.StatusAccepted,
			Body:    `{"data":[{"id":"created-clip","edit_url":"https://www.twitch.tv/twitchdev/clip/created-clip"}]}`,
			Success: true,
		},
		Want: testkit.ContractExpectation{Attempts: 1},
	}
	transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(fixture.Response))
	client, err := newTask17Client(transport, helix.Credential{
		AccessToken: "user-token",
		ClientID:    "client-id",
		TokenClass:  helix.TokenClassUser,
		Scopes:      []helix.AuthorizationScope{helix.ScopeClipsEdit},
	})
	if err != nil {
		t.Fatal(err)
	}

	result, callErr := client.Clips.CreateClip(context.Background(), helix.CreateClipRequest{
		BroadcasterID: "141981764",
		Title:         stringPtr("Task 17 clip"),
		Duration:      float64Ptr(45.5),
	})
	if callErr != nil {
		t.Fatal(callErr)
	}
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, "create-clip"), fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return result.Meta, nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 || result.Data[0].ID != "created-clip" || result.Data[0].EditURL == "" {
		t.Fatalf("created clip = %#v", result.Data)
	}
}

func TestClipsGetClips_encodesRepeatedIDsAndFilters(t *testing.T) {
	fixture := testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query: urlValues(
				"id", "clip-a", "id", "clip-b",
				"started_at", "2026-01-01T00:00:00Z",
				"ended_at", "2026-01-08T00:00:00Z",
				"first", "2", "before", "previous-cursor", "is_featured", "false",
			),
			Headers: http.Header{
				"Authorization": {"Bearer app-token"},
				"Client-Id":     {"client-id"},
			},
		},
		Response: testkit.ContractResponse{
			Status: http.StatusOK,
			Headers: http.Header{
				"Ratelimit-Limit":     {"8000"},
				"Ratelimit-Remaining": {"7999"},
				"Ratelimit-Reset":     {verticalSliceRateReset},
			},
			Body:    `{"data":[{"id":"clip-a","url":"https://clips.twitch.tv/clip-a","embed_url":"https://clips.twitch.tv/embed?clip=clip-a","broadcaster_id":"141981764","broadcaster_name":"TwitchDev","creator_id":"creator-a","creator_name":"Creator A","video_id":"video-a","game_id":"33214","language":"en","title":"Featured clip","view_count":42,"created_at":"2026-01-02T03:04:05Z","thumbnail_url":"https://cdn.test/clip-a.jpg","duration":12.9,"vod_offset":1957,"is_featured":true}],"pagination":{"cursor":"next-cursor"}}`,
			Success: true,
		},
		Want: testkit.ContractExpectation{Attempts: 1, RateLimitValid: true},
	}
	fixture.Want.RateLimit.Limit = 8000
	fixture.Want.RateLimit.Remaining = 7999
	transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(fixture.Response))
	client, err := newTask17Client(transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	if err != nil {
		t.Fatal(err)
	}

	result, callErr := client.Clips.GetClips(context.Background(), helix.GetClipsRequest{
		IDs:        []string{"clip-a", "clip-b"},
		StartedAt:  stringPtr("2026-01-01T00:00:00Z"),
		EndedAt:    stringPtr("2026-01-08T00:00:00Z"),
		First:      intPtr(2),
		Before:     stringPtr("previous-cursor"),
		IsFeatured: boolPtr(false),
	})
	if callErr != nil {
		t.Fatal(callErr)
	}
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, "get-clips"), fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return result.Meta, nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 || result.Data[0].ID != "clip-a" || result.Data[0].VODOffset == nil || *result.Data[0].VODOffset != 1957 || !result.Data[0].CreatedAt.Time.Equal(timestamp("2026-01-02T03:04:05Z").Time) {
		t.Fatalf("decoded clips = %#v", result.Data)
	}
}

func TestClipsGetClipsPager_usesAfterCursor(t *testing.T) {
	fixture := testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query: urlValues("broadcaster_id", "141981764", "first", "1", "after", "initial-cursor"),
			Headers: http.Header{
				"Authorization": {"Bearer app-token"},
				"Client-Id":     {"client-id"},
			},
		},
		Response: testkit.ContractResponse{
			Status:  http.StatusOK,
			Body:    `{"data":[{"id":"page-clip","created_at":"2026-01-02T03:04:05Z"}],"pagination":{}}`,
			Success: true,
		},
		Want: testkit.ContractExpectation{Attempts: 1},
	}
	transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(fixture.Response))
	client, err := newTask17Client(transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	if err != nil {
		t.Fatal(err)
	}

	pager, err := client.Clips.GetClipsPager(helix.GetClipsRequest{
		BroadcasterID: "141981764",
		First:         intPtr(1),
		After:         stringPtr("initial-cursor"),
	}, helix.WithPageLimit(1))
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) {
		t.Fatalf("pager.Next() = false, error = %v", pager.Err())
	}
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, "get-clips"), fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return pager.Page().Meta, nil
	}); err != nil {
		t.Fatal(err)
	}
	if pager.Page().Data[0].ID != "page-clip" {
		t.Fatalf("paged clip = %#v", pager.Page().Data)
	}
}

func TestClipsGetClipsDownload_returnsMetadataWithoutFetchingURLs(t *testing.T) {
	fixture := testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query: urlValues("broadcaster_id", "141981764", "editor_id", "141981764", "clip_id", "clip-a", "clip_id", "clip-b"),
			Headers: http.Header{
				"Authorization": {"Bearer user-token"},
				"Client-Id":     {"client-id"},
			},
		},
		Response: testkit.ContractResponse{
			Status:  http.StatusOK,
			Body:    `{"data":[{"clip_id":"clip-a","landscape_download_url":"https://download.trap.test/landscape-a","portrait_download_url":null},{"clip_id":"clip-b","landscape_download_url":null,"portrait_download_url":"https://download.trap.test/portrait-b"}]}`,
			Success: true,
		},
		Want: testkit.ContractExpectation{Attempts: 1},
	}
	recording := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(fixture.Response))
	transport := testkit.NewOfflineRoundTripper(recording, "api.twitch.test")
	client, err := newTask17Client(transport, helix.Credential{
		AccessToken: "user-token",
		ClientID:    "client-id",
		TokenClass:  helix.TokenClassUser,
		UserID:      "141981764",
		Scopes:      []helix.AuthorizationScope{helix.ScopeEditorManageClips},
	})
	if err != nil {
		t.Fatal(err)
	}

	result, callErr := client.Experimental.Clips.GetClipsDownload(context.Background(), helix.GetClipsDownloadRequest{
		EditorID:      "141981764",
		BroadcasterID: "141981764",
		ClipIDs:       []string{"clip-a", "clip-b"},
	})
	if callErr != nil {
		t.Fatal(callErr)
	}
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, "get-clips-download"), fixture, recording, func(context.Context) (helix.ResponseMeta, error) {
		return result.Meta, nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 2 || result.Data[0].LandscapeDownloadURL == nil || *result.Data[0].LandscapeDownloadURL != "https://download.trap.test/landscape-a" || result.Data[0].PortraitDownloadURL != nil {
		t.Fatalf("download metadata = %#v", result.Data)
	}
	if len(recording.Requests()) != 1 {
		t.Fatalf("requests = %d, want only Helix request", len(recording.Requests()))
	}
}
