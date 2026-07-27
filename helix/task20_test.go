package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

const task20RateReset = "4102444800"

func task20Body(t *testing.T, name string) string {
	t.Helper()
	body, err := testkit.LoadText("testdata/task20", name)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func task20Response(status int, body string) testkit.RoundTripResponse {
	return testkit.RoundTripResponse{
		StatusCode: status,
		Header: http.Header{
			"Ratelimit-Limit":     {"8000"},
			"Ratelimit-Remaining": {"7999"},
			"Ratelimit-Reset":     {task20RateReset},
		},
		Body: body,
	}
}

func task20Client(t *testing.T, transport *testkit.RecordingRoundTripper, credential helix.Credential) *helix.Client {
	t.Helper()
	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(credential),
	)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func task20Success(body string) testkit.ContractResponse {
	return testkit.ContractResponse{
		Status: http.StatusOK,
		Headers: http.Header{
			"Ratelimit-Limit":     {"8000"},
			"Ratelimit-Remaining": {"7999"},
			"Ratelimit-Reset":     {task20RateReset},
		},
		Body:    body,
		Success: true,
	}
}

func task20Contract(t *testing.T, anchor string, query map[string][]string, body string, response testkit.ContractResponse, transport *testkit.RecordingRoundTripper, meta helix.ResponseMeta, callErr error) {
	t.Helper()
	fixture := testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query:   query,
			Headers: http.Header{"Authorization": {"Bearer user-token"}, "Client-Id": {"client-id"}},
			Body:    body,
		},
		Response: response,
		Want: testkit.ContractExpectation{
			Attempts:       1,
			RateLimitValid: true,
		},
	}
	fixture.Want.RateLimit.Limit = 8000
	fixture.Want.RateLimit.Remaining = 7999
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, anchor), fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return meta, callErr
	}); err != nil {
		t.Fatal(err)
	}
}

func TestGamesGetTopGames(t *testing.T) {
	// Given a top-games response and an app credential.
	body := task20Body(t, "top_games_page_1.json")
	transport := testkit.NewRecordingRoundTripper(task20Response(http.StatusOK, body))
	client := task20Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "141981764"})
	first := 2

	// When the top-games endpoint is requested with an explicit page size.
	result, err := client.Games.GetTopGames(context.Background(), helix.GetTopGamesRequest{First: &first})

	// Then the exact request contract and game fields are preserved.
	if err != nil {
		t.Fatal(err)
	}
	task20Contract(t, "get-top-games", urlValues("first", "2"), "", task20Success(body), transport, result.Meta, nil)
	if len(result.Data) != 1 || result.Data[0].ID != "493057" || result.Data[0].IGDBID != "27789" || result.Pagination.Cursor() != "games-next" {
		t.Fatalf("top games = %#v, cursor = %q", result.Data, result.Pagination.Cursor())
	}
}

func TestGamesGetTopGamesPager_forward(t *testing.T) {
	// Given two pages returned by a forward cursor.
	firstBody := task20Body(t, "top_games_page_1.json")
	secondBody := task20Body(t, "top_games_page_2.json")
	transport := testkit.NewRecordingRoundTripper(task20Response(http.StatusOK, firstBody), task20Response(http.StatusOK, secondBody))
	client := task20Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	first := 1

	// When the pager advances twice.
	pager, err := client.Games.GetTopGamesPager(helix.GetTopGamesRequest{First: &first})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || !pager.Next(context.Background()) || pager.Next(context.Background()) {
		t.Fatal("unexpected forward pager state")
	}

	// Then the second request uses the response cursor as after.
	requests := transport.Requests()
	if len(requests) != 2 || requests[0].Path != "/helix/games/top?first=1" || requests[1].Path != "/helix/games/top?after=games-next&first=1" {
		t.Fatalf("forward pager requests = %#v", requests)
	}
	if pager.Err() != nil || pager.Page().Data[0].ID != "33214" {
		t.Fatalf("forward pager page = %#v, err = %v", pager.Page(), pager.Err())
	}
}

func TestGamesGetTopGamesPager_backward(t *testing.T) {
	// Given a backward cursor and two pages returned by the endpoint.
	firstBody := task20Body(t, "top_games_page_1.json")
	secondBody := task20Body(t, "top_games_page_2.json")
	transport := testkit.NewRecordingRoundTripper(task20Response(http.StatusOK, firstBody), task20Response(http.StatusOK, secondBody))
	client := task20Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	first := 1
	before := "games-before"

	// When the backward pager advances twice.
	pager, err := client.Games.GetTopGamesPager(helix.GetTopGamesRequest{First: &first, Before: &before})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || !pager.Next(context.Background()) || pager.Next(context.Background()) {
		t.Fatal("unexpected backward pager state")
	}

	// Then both requests use the before cursor parameter.
	requests := transport.Requests()
	if len(requests) != 2 || requests[0].Path != "/helix/games/top?before=games-before&first=1" || requests[1].Path != "/helix/games/top?before=games-next&first=1" {
		t.Fatalf("backward pager requests = %#v", requests)
	}
}

func TestGoalsGetCreatorGoals(t *testing.T) {
	// Given a creator-goals response and a scoped broadcaster user token.
	body := task20Body(t, "creator_goals.json")
	transport := testkit.NewRecordingRoundTripper(task20Response(http.StatusOK, body))
	client := task20Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "141981764", Scopes: []helix.AuthorizationScope{helix.ScopeChannelReadGoals}})

	// When the broadcaster's goals are requested.
	result, err := client.Goals.GetCreatorGoals(context.Background(), helix.GetCreatorGoalsRequest{BroadcasterID: "141981764"})

	// Then the scoped request and RFC3339 goal data are decoded exactly.
	if err != nil {
		t.Fatal(err)
	}
	task20Contract(t, "get-creator-goals", urlValues("broadcaster_id", "141981764"), "", task20Success(body), transport, result.Meta, nil)
	if len(result.Data) != 1 || result.Data[0].ID != "goal-1" || result.Data[0].CurrentAmount != 27062 || result.Data[0].CreatedAt.IsZero() {
		t.Fatalf("creator goals = %#v", result.Data)
	}
}

func TestGoalsGetCreatorGoals_rejectsMissingScopeAndMismatchedSubject(t *testing.T) {
	// Given user credentials that violate the local scope or broadcaster binding.
	tests := []struct {
		name       string
		credential helix.Credential
	}{
		{name: "missing scope", credential: helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "141981764"}},
		{name: "mismatched broadcaster", credential: helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "different-user", Scopes: []helix.AuthorizationScope{helix.ScopeChannelReadGoals}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper()
			client := task20Client(t, transport, test.credential)

			// When goals are requested for the broadcaster.
			_, err := client.Goals.GetCreatorGoals(context.Background(), helix.GetCreatorGoalsRequest{BroadcasterID: "141981764"})

			// Then local authorization rejects the request before network I/O.
			var authErr *helix.AuthError
			if !errors.As(err, &authErr) || len(transport.Requests()) != 0 {
				t.Fatalf("error = %T %v, requests = %d", err, err, len(transport.Requests()))
			}
		})
	}
}

func TestGoalsGetCreatorGoals_returnsTypedAPIError(t *testing.T) {
	// Given a scoped user request and an upstream unauthorized response.
	body := task20Body(t, "api_unauthorized.json")
	transport := testkit.NewRecordingRoundTripper(task20Response(http.StatusUnauthorized, body))
	client := task20Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "141981764", Scopes: []helix.AuthorizationScope{helix.ScopeChannelReadGoals}})

	// When the endpoint returns its declared error shape.
	_, err := client.Goals.GetCreatorGoals(context.Background(), helix.GetCreatorGoalsRequest{BroadcasterID: "141981764"})

	// Then the shared transport exposes a typed authorization error.
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("error = %T %v", err, err)
	}
}

func TestHypeTrainGetHypeTrainStatus(t *testing.T) {
	// Given an active Hype Train response and a scoped broadcaster user token.
	body := task20Body(t, "hype_train_status.json")
	transport := testkit.NewRecordingRoundTripper(task20Response(http.StatusOK, body))
	client := task20Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "1337", Scopes: []helix.AuthorizationScope{helix.ScopeChannelReadHypeTrain}})

	// When the broadcaster's Hype Train status is requested.
	result, err := client.HypeTrain.GetHypeTrainStatus(context.Background(), helix.GetHypeTrainStatusRequest{BroadcasterID: "1337"})

	// Then nullable nested records, timestamps, and contribution fields decode.
	if err != nil {
		t.Fatal(err)
	}
	task20Contract(t, "get-hype-train-status", urlValues("broadcaster_id", "1337"), "", task20Success(body), transport, result.Meta, nil)
	if len(result.Data) != 1 || result.Data[0].Current == nil || result.Data[0].Current.Level != 2 || len(result.Data[0].Current.TopContributions) != 2 || result.Data[0].AllTimeHigh == nil || result.Data[0].SharedAllTimeHigh == nil {
		t.Fatalf("hype train status = %#v", result.Data)
	}
}

func TestHypeTrainGetHypeTrainStatus_rejectsMissingScopeAndMismatchedSubject(t *testing.T) {
	// Given user credentials that violate the local scope or broadcaster binding.
	tests := []struct {
		name       string
		credential helix.Credential
	}{
		{name: "missing scope", credential: helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "1337"}},
		{name: "mismatched broadcaster", credential: helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "different-user", Scopes: []helix.AuthorizationScope{helix.ScopeChannelReadHypeTrain}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper()
			client := task20Client(t, transport, test.credential)

			// When Hype Train status is requested for the broadcaster.
			_, err := client.HypeTrain.GetHypeTrainStatus(context.Background(), helix.GetHypeTrainStatusRequest{BroadcasterID: "1337"})

			// Then local authorization rejects the request before network I/O.
			var authErr *helix.AuthError
			if !errors.As(err, &authErr) || len(transport.Requests()) != 0 {
				t.Fatalf("error = %T %v, requests = %d", err, err, len(transport.Requests()))
			}
		})
	}
}

func TestHypeTrainGetHypeTrainStatus_returnsTypedAPIError(t *testing.T) {
	// Given a scoped user request and an upstream internal error.
	body := task20Body(t, "api_internal_error.json")
	transport := testkit.NewRecordingRoundTripper(task20Response(http.StatusInternalServerError, body))
	client := task20Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "1337", Scopes: []helix.AuthorizationScope{helix.ScopeChannelReadHypeTrain}})

	// When the endpoint returns its declared error shape.
	_, err := client.HypeTrain.GetHypeTrainStatus(context.Background(), helix.GetHypeTrainStatusRequest{BroadcasterID: "1337"})

	// Then the shared transport exposes a typed API error.
	var apiErr *helix.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode() != http.StatusInternalServerError {
		t.Fatalf("error = %T %v", err, err)
	}
}
