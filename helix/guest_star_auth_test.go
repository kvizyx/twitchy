package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestExperimentalGuestStarRejectsAppTokensBeforeNetwork(t *testing.T) {
	// Given an app token that cannot authenticate Guest Star.
	transport := testkit.NewRecordingRoundTripper()
	client, err := guestStarClient(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "guest-star-client", TokenClass: helix.TokenClassApp})
	if err != nil {
		t.Fatal(err)
	}

	// When a Guest Star operation is called with the app token.
	_, callErr := client.Experimental.GuestStar.CreateGuestStarSession(context.Background(), helix.CreateGuestStarSessionRequest{BroadcasterID: "5678"})

	// Then local authorization rejects the call without a network request.
	requireGuestStarAuthError(t, callErr)
	if len(transport.Requests()) != 0 {
		t.Fatal("app-token Guest Star call reached the network")
	}
}

func TestExperimentalGuestStarAcceptsScopeAlternatives(t *testing.T) {
	// Given each documented read scope alternative.
	for _, scope := range []helix.AuthorizationScope{helix.ScopeChannelReadGuestStar, helix.ScopeChannelManageGuestStar, helix.ScopeModeratorReadGuestStar, helix.ScopeModeratorManageGuestStar} {
		t.Run(string(scope), func(t *testing.T) {
			response := testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: guestStarBody(t, "invites.json"), Header: http.Header{"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {guestStarRateReset}}}
			transport := testkit.NewRecordingRoundTripper(response)
			client, err := guestStarClient(t, transport, guestStarCredential(scope))
			if err != nil {
				t.Fatal(err)
			}

			// When the read operation is called with exactly one alternative scope.
			_, callErr := client.Experimental.GuestStar.GetGuestStarInvites(context.Background(), helix.GetGuestStarInvitesRequest{BroadcasterID: "1234", ModeratorID: "5678", SessionID: "session-1"})

			// Then the alternative is accepted and sent once.
			if callErr != nil {
				t.Fatal(callErr)
			}
			if len(transport.Requests()) != 1 {
				t.Fatalf("requests = %d, want 1", len(transport.Requests()))
			}
		})
	}
}

func TestExperimentalGuestStarRejectsWrongModeratorSubjectBeforeNetwork(t *testing.T) {
	// Given a user token whose subject does not match moderator_id.
	transport := testkit.NewRecordingRoundTripper()
	client, err := guestStarClient(t, transport, guestStarCredential(helix.ScopeModeratorManageGuestStar))
	if err != nil {
		t.Fatal(err)
	}

	// When a moderator-bound operation is called with a different moderator.
	_, callErr := client.Experimental.GuestStar.SendGuestStarInvite(context.Background(), helix.SendGuestStarInviteRequest{BroadcasterID: "1234", ModeratorID: "wrong", SessionID: "session-1", GuestID: "9012"})

	// Then the mismatch is rejected before transport execution.
	requireGuestStarAuthError(t, callErr)
	if len(transport.Requests()) != 0 {
		t.Fatal("wrong moderator subject reached the network")
	}
}

func TestExperimentalGuestStarTypedFailurePreservesResponseMetadata(t *testing.T) {
	// Given a structured Guest Star API failure with rate-limit headers.
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: http.StatusForbidden, Header: http.Header{"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7998"}, "Ratelimit-Reset": {guestStarRateReset}}, Body: `{"error":"Forbidden","status":403,"message":"not permitted"}`})
	client, err := guestStarClient(t, transport, guestStarCredential(helix.ScopeModeratorManageGuestStar))
	if err != nil {
		t.Fatal(err)
	}

	// When the API returns the typed failure.
	_, callErr := client.Experimental.GuestStar.DeleteGuestStarInvite(context.Background(), helix.DeleteGuestStarInviteRequest{BroadcasterID: "1234", ModeratorID: "5678", SessionID: "session-1", GuestID: "9012"})

	// Then the typed auth error retains status and rate metadata.
	var authErr *helix.AuthError
	if !errors.As(callErr, &authErr) || authErr.StatusCode() != http.StatusForbidden || !authErr.Meta().RateLimit().Valid() {
		t.Fatalf("error = %T %v, want forbidden AuthError with rate metadata", callErr, callErr)
	}
}
