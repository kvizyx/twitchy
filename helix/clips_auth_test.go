package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestClipsCreateClipFromVOD_scopeAlternativesAllowEitherScope(t *testing.T) {
	for _, test := range []struct {
		name  string
		scope helix.AuthorizationScope
	}{
		{name: "channel scope", scope: helix.ScopeChannelManageClips},
		{name: "editor scope", scope: helix.ScopeEditorManageClips},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := testkit.ContractFixture{
				Request: testkit.ContractRequest{
					Query: urlValues("broadcaster_id", "141981764", "editor_id", "12826", "vod_id", "2277656159", "vod_offset", "434", "duration", "30.5", "title", "VOD clip"),
					Headers: http.Header{
						"Authorization": {"Bearer user-token"},
						"Client-Id":     {"client-id"},
					},
				},
				Response: testkit.ContractResponse{
					Status:  http.StatusAccepted,
					Body:    `{"data":[{"id":"vod-clip","edit_url":"https://www.twitch.tv/twitchdev/clip/vod-clip"}]}`,
					Success: true,
				},
				Want: testkit.ContractExpectation{Attempts: 1},
			}
			transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(fixture.Response))
			client, err := newTask17Client(transport, helix.Credential{
				AccessToken: "user-token",
				ClientID:    "client-id",
				TokenClass:  helix.TokenClassUser,
				UserID:      "12826",
				Scopes:      []helix.AuthorizationScope{test.scope},
			})
			if err != nil {
				t.Fatal(err)
			}

			result, callErr := client.Experimental.Clips.CreateClipFromVOD(context.Background(), helix.CreateClipFromVODRequest{
				EditorID:      "12826",
				BroadcasterID: "141981764",
				VODID:         "2277656159",
				VODOffset:     434,
				Duration:      float64Ptr(30.5),
				Title:         "VOD clip",
			})
			if callErr != nil {
				t.Fatal(callErr)
			}
			if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, "create-clip-from-vod"), fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
				return result.Meta, nil
			}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestClipsGetClipsDownload_rejectsMissingScopeBeforeNetwork(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper()
	client, err := newTask17Client(transport, helix.Credential{
		AccessToken: "user-token",
		ClientID:    "client-id",
		TokenClass:  helix.TokenClassUser,
		UserID:      "141981764",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, callErr := client.Experimental.Clips.GetClipsDownload(context.Background(), helix.GetClipsDownloadRequest{
		EditorID:      "141981764",
		BroadcasterID: "141981764",
		ClipIDs:       []string{"clip-a"},
	})
	var authErr *helix.AuthError
	if !errors.As(callErr, &authErr) {
		t.Fatalf("error = %T %v, want AuthError", callErr, callErr)
	}
	if len(transport.Requests()) != 0 {
		t.Fatal("missing clip scope reached the network")
	}
}

func TestClips_typedFailureIncludesStatusMetadata(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{
		StatusCode: http.StatusForbidden,
		Header: http.Header{
			"Ratelimit-Limit":     {"8000"},
			"Ratelimit-Remaining": {"7998"},
			"Ratelimit-Reset":     {verticalSliceRateReset},
		},
		Body: `{"error":"Forbidden","status":403,"message":"not allowed"}`,
	})
	client, err := newTask17Client(transport, helix.Credential{
		AccessToken: "user-token",
		ClientID:    "client-id",
		TokenClass:  helix.TokenClassUser,
		Scopes:      []helix.AuthorizationScope{helix.ScopeClipsEdit},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, callErr := client.Clips.CreateClip(context.Background(), helix.CreateClipRequest{BroadcasterID: "141981764"})
	var authErr *helix.AuthError
	if !errors.As(callErr, &authErr) || authErr.StatusCode() != http.StatusForbidden || !authErr.Meta().RateLimit().Valid() {
		t.Fatalf("error = %T %v, want forbidden AuthError with rate metadata", callErr, callErr)
	}
}
