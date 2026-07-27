package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

type stubResolverEntry struct {
	source  helix.TokenSource
	intents []helix.Intent
}

type stubResolver struct {
	users map[string]stubResolverEntry
}

func (r *stubResolver) SourceForUser(userID string) (helix.TokenSource, error) {
	entry, ok := r.users[userID]
	if !ok {
		return nil, helix.ErrUserNotFound
	}
	return entry.source, nil
}

func (r *stubResolver) SourceForIntent(intents ...helix.Intent) (helix.TokenSource, error) {
	for _, id := range []string{"111", "222"} {
		entry, ok := r.users[id]
		if !ok {
			continue
		}
		if stubCovers(entry.intents, intents) {
			return entry.source, nil
		}
	}
	return nil, helix.ErrIntentNotCovered
}

func stubCovers(have, want []helix.Intent) bool {
	set := make(map[helix.Intent]struct{}, len(have))
	for _, intent := range have {
		set[intent] = struct{}{}
	}
	for _, intent := range want {
		if _, ok := set[intent]; !ok {
			return false
		}
	}
	return true
}

func staticSource(accessToken, userID string) helix.TokenSource {
	return helix.NewStaticTokenSource(helix.Credential{
		AccessToken: accessToken,
		ClientID:    "client-id",
		TokenClass:  helix.TokenClassUser,
		UserID:      userID,
		Scopes:      []helix.AuthorizationScope{helix.ScopeModeratorReadFollowers},
	})
}

func TestClientAs_requiresResolver(t *testing.T) {
	client, err := helix.New(helix.WithStaticToken(helix.Credential{AccessToken: "app", TokenClass: helix.TokenClassApp}))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.AsUser("111"); !errors.Is(err, helix.ErrNoCredentialResolver) {
		t.Fatalf("AsUser error = %v, want ErrNoCredentialResolver", err)
	}
	if _, err := client.AsIntent("chat"); !errors.Is(err, helix.ErrNoCredentialResolver) {
		t.Fatalf("AsIntent error = %v, want ErrNoCredentialResolver", err)
	}
}

func TestClientAsUser_derivesClientWithUserCredential(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[]}`},
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[]}`},
		testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: `{"data":[]}`},
	)
	resolver := &stubResolver{users: map[string]stubResolverEntry{
		"111": {source: staticSource("token-one", "111"), intents: []helix.Intent{"chat"}},
		"222": {source: staticSource("token-two", "222"), intents: []helix.Intent{"chat", "eventsub"}},
	}}
	root, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(helix.Credential{AccessToken: "root-token", ClientID: "client-id", TokenClass: helix.TokenClassApp}),
		helix.WithCredentialResolver(resolver),
	)
	if err != nil {
		t.Fatal(err)
	}

	first, err := root.AsUser("111")
	if err != nil {
		t.Fatal(err)
	}
	second, err := root.AsIntent("eventsub")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if _, err := first.Users.GetUsers(ctx, helix.GetUsersRequest{}); err != nil {
		t.Fatal(err)
	}
	if _, err := second.Users.GetUsers(ctx, helix.GetUsersRequest{}); err != nil {
		t.Fatal(err)
	}
	if _, err := root.Users.GetUsers(ctx, helix.GetUsersRequest{}); err != nil {
		t.Fatal(err)
	}

	requests := transport.Requests()
	if len(requests) != 3 {
		t.Fatalf("requests = %d, want 3", len(requests))
	}
	wantAuth := []string{"Bearer token-one", "Bearer token-two", "Bearer root-token"}
	for index, want := range wantAuth {
		if got := requests[index].Header.Get("Authorization"); got != want {
			t.Errorf("request %d Authorization = %q, want %q", index, got, want)
		}
	}
}

func TestClientAs_resolutionFailure(t *testing.T) {
	resolver := &stubResolver{users: map[string]stubResolverEntry{}}
	root, err := helix.New(
		helix.WithStaticToken(helix.Credential{AccessToken: "root", TokenClass: helix.TokenClassApp}),
		helix.WithCredentialResolver(resolver),
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := root.AsUser("missing"); !errors.Is(err, helix.ErrUserNotFound) {
		t.Fatalf("AsUser error = %v, want ErrUserNotFound", err)
	}
	if _, err := root.AsIntent("chat"); !errors.Is(err, helix.ErrIntentNotCovered) {
		t.Fatalf("AsIntent error = %v, want ErrIntentNotCovered", err)
	}
}

func TestClientAsUser_keepsSubjectBindingValidation(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper()
	resolver := &stubResolver{users: map[string]stubResolverEntry{
		"111": {source: staticSource("token-one", "111")},
	}}
	root, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(helix.Credential{AccessToken: "root", TokenClass: helix.TokenClassApp}),
		helix.WithCredentialResolver(resolver),
	)
	if err != nil {
		t.Fatal(err)
	}
	derived, err := root.AsUser("111")
	if err != nil {
		t.Fatal(err)
	}
	_, err = derived.Channels.GetChannelFollowers(context.Background(), helix.GetChannelFollowersRequest{UserID: stringPointer("999"), BroadcasterID: "111"})
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("GetChannelFollowers error = %v, want *helix.AuthError", err)
	}
	if got := len(transport.Requests()); got != 0 {
		t.Fatalf("requests = %d, want 0 (rejected before network)", got)
	}
}

func TestWithCredentialResolver_rejectsNil(t *testing.T) {
	if _, err := helix.New(helix.WithCredentialResolver(nil)); !errors.Is(err, helix.ErrInvalidOption) {
		t.Fatalf("New error = %v, want ErrInvalidOption", err)
	}
}
