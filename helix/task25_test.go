package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func task25Body(t *testing.T, name string) string {
	t.Helper()
	body, err := testkit.LoadText("testdata/task25", name)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func task25Client(t *testing.T, transport *testkit.RecordingRoundTripper, credential helix.Credential) *helix.Client {
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

func task25Response(status int, body string) testkit.RoundTripResponse {
	return testkit.RoundTripResponse{
		StatusCode: status,
		Header: http.Header{
			"Ratelimit-Limit":     {"8000"},
			"Ratelimit-Remaining": {"7999"},
			"Ratelimit-Reset":     {"4102444800"},
		},
		Body: body,
	}
}

func TestSubscriptionsGetBroadcasterSubscriptions_supportsExtensionManagerAppToken(t *testing.T) {
	first := task25Body(t, "subscriptions.json")
	second := `{"data":[],"pagination":{}}
`
	transport := testkit.NewRecordingRoundTripper(task25Response(http.StatusOK, first), task25Response(http.StatusOK, first), task25Response(http.StatusOK, second))
	client := task25Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})

	result, err := client.Subscriptions.GetBroadcasterSubscriptions(context.Background(), helix.GetBroadcasterSubscriptionsRequest{BroadcasterID: "broadcaster-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Data.Subscriptions) != 1 || result.Data.Subscriptions[0].Tier != helix.SubscriptionTier1000 || result.Data.Points == nil || *result.Data.Points != 13 || result.Data.Total == nil || *result.Data.Total != 12 {
		t.Fatalf("subscriptions = %#v", result.Data)
	}
	if result.Data.Subscriptions[0].GifterID != "gifter-1" || !result.Data.Subscriptions[0].IsGift {
		t.Fatalf("gifter fields = %#v", result.Data.Subscriptions[0])
	}

	pager, err := client.Subscriptions.GetBroadcasterSubscriptionsPager(helix.GetBroadcasterSubscriptionsRequest{BroadcasterID: "broadcaster-1"})
	if err != nil || !pager.Next(context.Background()) || !pager.Next(context.Background()) || pager.Next(context.Background()) || pager.Err() != nil {
		t.Fatalf("pager = page=%#v err=%v", pager.Page(), pager.Err())
	}
	requests := transport.Requests()
	if len(requests) != 3 || requests[0].Path != "/helix/subscriptions?broadcaster_id=broadcaster-1" || requests[2].Path != "/helix/subscriptions?after=next&broadcaster_id=broadcaster-1" {
		t.Fatalf("subscription requests = %#v", requests)
	}
}

func TestSubscriptionsCheckUserSubscription_requiresUserScopeAndPreservesGifterFields(t *testing.T) {
	body := task25Body(t, "check_subscription.json")
	transport := testkit.NewRecordingRoundTripper(task25Response(http.StatusOK, body))
	client := task25Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "user-1", Scopes: []helix.AuthorizationScope{helix.ScopeUserReadSubscriptions}})

	result, err := client.Subscriptions.CheckUserSubscription(context.Background(), helix.CheckUserSubscriptionRequest{BroadcasterID: "broadcaster-1", UserID: "user-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 || result.Data[0].Tier != helix.SubscriptionTier2000 || result.Data[0].GifterLogin != "gifter" || !result.Data[0].IsGift {
		t.Fatalf("subscription = %#v", result.Data)
	}
	if got := transport.Requests()[0].Path; got != "/helix/subscriptions/user?broadcaster_id=broadcaster-1&user_id=user-1" {
		t.Fatalf("request path = %q", got)
	}

	missingScopeTransport := testkit.NewRecordingRoundTripper()
	missingScope := task25Client(t, missingScopeTransport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "user-1"})
	_, err = missingScope.Subscriptions.CheckUserSubscription(context.Background(), helix.CheckUserSubscriptionRequest{BroadcasterID: "broadcaster-1", UserID: "user-1"})
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) || len(missingScopeTransport.Requests()) != 0 {
		t.Fatalf("missing scope error = %T %v", err, err)
	}

	mismatchTransport := testkit.NewRecordingRoundTripper()
	mismatch := task25Client(t, mismatchTransport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "different-user", Scopes: []helix.AuthorizationScope{helix.ScopeUserReadSubscriptions}})
	_, err = mismatch.Subscriptions.CheckUserSubscription(context.Background(), helix.CheckUserSubscriptionRequest{BroadcasterID: "broadcaster-1", UserID: "user-1"})
	if !errors.As(err, &authErr) || len(mismatchTransport.Requests()) != 0 {
		t.Fatalf("subject mismatch error = %T %v calls=%d", err, err, len(mismatchTransport.Requests()))
	}
}

func TestSubscriptionsGetBroadcasterSubscriptions_returnsTypedAuthError(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(task25Response(http.StatusUnauthorized, task25Body(t, "unauthorized.json")))
	client := task25Client(t, transport, helix.Credential{AccessToken: "secret-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	_, err := client.Subscriptions.GetBroadcasterSubscriptions(context.Background(), helix.GetBroadcasterSubscriptionsRequest{BroadcasterID: "broadcaster-1"})
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) || authErr.StatusCode() != http.StatusUnauthorized {
		t.Fatalf("auth error = %T %v", err, err)
	}
}

func TestTagsGetAllStreamTags_repeatsTagIDsAndPaginates(t *testing.T) {
	first := task25Body(t, "stream_tags.json")
	second := `{"data":[],"pagination":{}}
`
	transport := testkit.NewRecordingRoundTripper(task25Response(http.StatusOK, first), task25Response(http.StatusOK, second))
	client := task25Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	ids := []string{"tag-1", "tag-2"}
	pager, err := client.Tags.GetAllStreamTagsPager(helix.GetAllStreamTagsRequest{TagIDs: ids, First: intPointer25(1)})
	if err != nil || !pager.Next(context.Background()) || !pager.Next(context.Background()) || pager.Next(context.Background()) || pager.Err() != nil {
		t.Fatalf("tag pager = page=%#v err=%v", pager.Page(), pager.Err())
	}
	if len(pager.Page().Data) != 0 {
		t.Fatalf("last tag page = %#v", pager.Page().Data)
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[0].Path != "/helix/tags/streams?first=1&tag_id=tag-1&tag_id=tag-2" || requests[1].Path != "/helix/tags/streams?after=next&first=1&tag_id=tag-1&tag_id=tag-2" {
		t.Fatalf("tag requests = %#v", requests)
	}
}

func TestTagsGetStreamTags_preservesLocalizationMaps(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(task25Response(http.StatusOK, task25Body(t, "stream_tags.json")))
	client := task25Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	result, err := client.Tags.GetStreamTags(context.Background(), helix.GetStreamTagsRequest{BroadcasterID: "broadcaster-1"})
	if err != nil || len(result.Data) != 1 || result.Data[0].LocalizationNames["en-us"] != "English" || result.Data[0].LocalizationDescriptions["en-us"] == "" {
		t.Fatalf("stream tags = %#v err=%v", result.Data, err)
	}
}

func TestTeamsGetChannelTeams_preservesNullableImageURLs(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(task25Response(http.StatusOK, task25Body(t, "channel_teams.json")))
	client := task25Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	result, err := client.Teams.GetChannelTeams(context.Background(), helix.GetChannelTeamsRequest{BroadcasterID: "broadcaster-1"})
	if err != nil || len(result.Data) != 1 || result.Data[0].BackgroundImageURL != nil || result.Data[0].TeamName != "livecoders" {
		t.Fatalf("channel teams = %#v err=%v", result.Data, err)
	}
}

func TestTeamsGetTeams_enforcesExclusiveSelectors(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(task25Response(http.StatusOK, task25Body(t, "teams.json")))
	client := task25Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	name := "livecoders"
	result, err := client.Teams.GetTeams(context.Background(), helix.GetTeamsRequest{Name: name})
	if err != nil || len(result.Data) != 1 || len(result.Data[0].Users) != 1 || result.Data[0].Users[0].UserID != "user-1" {
		t.Fatalf("teams = %#v err=%v", result.Data, err)
	}
	for _, request := range []helix.GetTeamsRequest{{}, {Name: name, ID: "6358"}} {
		before := len(transport.Requests())
		_, err := client.Teams.GetTeams(context.Background(), request)
		var requestErr *helix.RequestEncodingError
		if !errors.As(err, &requestErr) || len(transport.Requests()) != before {
			t.Fatalf("selector request=%#v error=%T %v calls=%d", request, err, err, len(transport.Requests()))
		}
	}
}

func intPointer25(value int) *int { return &value }
