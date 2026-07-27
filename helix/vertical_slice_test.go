package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

const verticalSliceRateReset = "4102444800"

func TestGamesGetGames(t *testing.T) {
	// Given a recording transport and an app credential allowed by get-games.
	fixture := testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query: urlValues("id", "33214", "name", "Fortnite", "igdb_id", "1905"),
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
			Body:    `{"data":[{"id":"33214","name":"Fortnite","box_art_url":"https://cdn.test/{width}x{height}.jpg","igdb_id":"1905"}]}`,
			Success: true,
		},
		Want: testkit.ContractExpectation{Attempts: 1, RateLimitValid: true},
	}
	fixture.Want.RateLimit.Limit = 8000
	fixture.Want.RateLimit.Remaining = 7999
	transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(fixture.Response))
	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp}),
	)
	if err != nil {
		t.Fatal(err)
	}

	// When the stable Games service is called with repeated filters.
	result, callErr := client.Games.GetGames(context.Background(), helix.GetGamesRequest{
		IDs:     []string{"33214"},
		Names:   []string{"Fortnite"},
		IGDBIDs: []string{"1905"},
	})

	// Then the typed game and response metadata are returned through the shared transport.
	if callErr != nil {
		t.Fatal(callErr)
	}
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, "get-games"), fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return result.Meta, nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 || result.Data[0].ID != "33214" || result.Data[0].IGDBID != "1905" {
		t.Fatalf("decoded games = %#v", result.Data)
	}
}

func TestExperimentalGetCustomPowerUp(t *testing.T) {
	// Given an experimental Bits call with a user credential and bits:read.
	fixture := testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query: urlValues("broadcaster_id", "274637212", "id", "power-up-id"),
			Headers: http.Header{
				"Authorization": {"Bearer user-token"},
				"Client-Id":     {"client-id"},
			},
		},
		Response: testkit.ContractResponse{
			Status: http.StatusOK,
			Headers: http.Header{
				"Ratelimit-Limit":     {"8000"},
				"Ratelimit-Remaining": {"7998"},
				"Ratelimit-Reset":     {verticalSliceRateReset},
			},
			Body:    `{"data":[{"broadcaster_id":"274637212","broadcaster_login":"torpedo09","broadcaster_name":"Torpedo09","id":"power-up-id","title":"game analysis","prompt":"","bits":100,"image":null,"default_image":{"url_1x":"https://cdn.test/1x.png","url_2x":"https://cdn.test/2x.png","url_4x":"https://cdn.test/4x.png"},"background_color":"#00FF00","is_enabled":true,"is_user_input_required":false,"max_per_stream_setting":{"is_enabled":false,"max_per_stream":0},"max_per_user_per_stream_setting":{"is_enabled":false,"max_per_user_per_stream":0},"global_cooldown_setting":{"is_enabled":false,"global_cooldown_seconds":0},"is_paused":false,"is_in_stock":true,"redemptions_redeemed_current_stream":null,"cooldown_expires_at":null}]}`,
			Success: true,
		},
		Want: testkit.ContractExpectation{Attempts: 1, RateLimitValid: true},
	}
	fixture.Want.RateLimit.Limit = 8000
	fixture.Want.RateLimit.Remaining = 7998
	transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(fixture.Response))
	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(helix.Credential{
			AccessToken: "user-token",
			ClientID:    "client-id",
			TokenClass:  helix.TokenClassUser,
			UserID:      "274637212",
			Scopes:      []helix.AuthorizationScope{helix.ScopeBitsRead},
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	// When the experimental namespace is called.
	result, callErr := client.Experimental.Bits.GetCustomPowerUp(context.Background(), helix.GetCustomPowerUpRequest{
		BroadcasterID: "274637212",
		IDs:           []string{"power-up-id"},
	})

	// Then the nested typed data and rate metadata are returned.
	if callErr != nil {
		t.Fatal(callErr)
	}
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, "get-custom-power-up"), fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return result.Meta, nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 || result.Data[0].Title != "game analysis" || result.Data[0].Image != nil {
		t.Fatalf("decoded custom power-ups = %#v", result.Data)
	}
}

func TestVerticalSlice_rejectsWrongTokenClassBeforeNetwork(t *testing.T) {
	// Given credentials that violate each operation's manifest token class.
	tests := []struct {
		name string
		call func(*helix.Client) error
	}{
		{
			name: "games rejects extension token",
			call: func(client *helix.Client) error {
				_, err := client.Games.GetGames(context.Background(), helix.GetGamesRequest{IDs: []string{"33214"}})
				return err
			},
		},
		{
			name: "custom power-up rejects app token",
			call: func(client *helix.Client) error {
				_, err := client.Experimental.Bits.GetCustomPowerUp(context.Background(), helix.GetCustomPowerUpRequest{BroadcasterID: "274637212"})
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper()
			client, err := helix.New(
				helix.WithBaseURL("https://api.twitch.test/helix"),
				helix.WithHTTPClient(&http.Client{Transport: transport}),
				helix.WithStaticToken(helix.Credential{AccessToken: "wrong-token", ClientID: "client-id", TokenClass: helix.TokenClassExtension}),
			)
			if test.name == "custom power-up rejects app token" {
				client, err = helix.New(
					helix.WithBaseURL("https://api.twitch.test/helix"),
					helix.WithHTTPClient(&http.Client{Transport: transport}),
					helix.WithStaticToken(helix.Credential{AccessToken: "wrong-token", ClientID: "client-id", TokenClass: helix.TokenClassApp}),
				)
			}
			if err != nil {
				t.Fatal(err)
			}

			// When the operation is invoked.
			callErr := test.call(client)

			// Then local auth validation rejects it before the transport is touched.
			var authErr *helix.AuthError
			if !errors.As(callErr, &authErr) {
				t.Fatalf("error = %T %v, want AuthError", callErr, callErr)
			}
			if len(transport.Requests()) != 0 {
				t.Fatal("wrong token class reached the network")
			}
		})
	}
}

func TestVerticalSlice_typedFailuresAndRateMetadata(t *testing.T) {
	// Given malformed success and well-formed 404 responses from the shared transport.
	tests := []struct {
		name       string
		status     int
		body       string
		wantErr    any
		wantStatus int
	}{
		{name: "malformed success", status: http.StatusOK, body: "not-json", wantErr: &helix.ProtocolError{}, wantStatus: http.StatusOK},
		{name: "not found", status: http.StatusNotFound, body: `{"error":"Not Found","status":404,"message":"missing"}`, wantErr: &helix.APIError{}, wantStatus: http.StatusNotFound},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{
				StatusCode: test.status,
				Header: http.Header{
					"Ratelimit-Limit":     {"8000"},
					"Ratelimit-Remaining": {"7997"},
					"Ratelimit-Reset":     {verticalSliceRateReset},
				},
				Body: test.body,
			})
			client, err := helix.New(
				helix.WithBaseURL("https://api.twitch.test/helix"),
				helix.WithHTTPClient(&http.Client{Transport: transport}),
				helix.WithStaticToken(helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp}),
			)
			if err != nil {
				t.Fatal(err)
			}

			// When the stable endpoint is invoked.
			_, callErr := client.Games.GetGames(context.Background(), helix.GetGamesRequest{IDs: []string{"33214"}})

			// Then the declared typed failure includes status/rate metadata and no replay occurs.
			if !errors.As(callErr, &test.wantErr) {
				t.Fatalf("error = %T %v, want %T", callErr, callErr, test.wantErr)
			}
			var metaErr interface{ Meta() helix.ResponseMeta }
			if !errors.As(callErr, &metaErr) || metaErr.Meta().StatusCode() != test.wantStatus || !metaErr.Meta().RateLimit().Valid() {
				t.Fatalf("error metadata = %#v, want status/rate metadata", callErr)
			}
			if len(transport.Requests()) != 1 {
				t.Fatalf("attempts = %d, want 1", len(transport.Requests()))
			}
		})
	}
}
