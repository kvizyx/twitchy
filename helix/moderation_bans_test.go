package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestModerationBanOperations(t *testing.T) {
	duration := 60
	tests := []moderationOperationCase{
		{
			anchor:     "get-banned-users",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "user_id", "9876", "first", "2", "after", "cursor-a"), http.StatusOK, `{"data":[{"user_id":"9876","user_login":"viewer","user_name":"Viewer","expires_at":"","created_at":"2026-01-02T03:04:05Z","reason":"reason","moderator_id":"5678","moderator_login":"mod","moderator_name":"Mod"}],"pagination":{"cursor":"cursor-b"}}`),
			credential: moderationCredential(helix.TokenClassUser, "1234", helix.ScopeModerationRead),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.GetBannedUsers(context.Background(), helix.GetBannedUsersRequest{BroadcasterID: "1234", UserIDs: []string{"9876"}, First: moderationInt(2), After: moderationString("cursor-a")}))
			},
		},
		{
			anchor:     "ban-user",
			fixture:    moderationBodyFixture(moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusOK, `{"data":[{"broadcaster_id":"1234","moderator_id":"5678","user_id":"9876","created_at":"2026-01-02T03:04:05Z","end_time":null}]}`), `{"data":{"user_id":"9876","duration":60,"reason":"reason"}}`),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageBannedUsers),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.BanUser(context.Background(), helix.BanUserRequest{BroadcasterID: "1234", ModeratorID: "5678", Data: helix.BanUserBody{UserID: "9876", Duration: &duration, Reason: "reason"}}))
			},
		},
		{
			anchor:     "unban-user",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "user_id", "9876"), http.StatusNoContent, ""),
			credential: moderationCredential(helix.TokenClassApp, "", helix.ScopeModeratorManageBannedUsers, helix.ScopeUserBot),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.UnbanUser(context.Background(), helix.UnbanUserRequest{BroadcasterID: "1234", ModeratorID: "5678", UserID: "9876"}))
			},
		},
		{
			anchor:     "get-unban-requests",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "status", "pending", "user_id", "9876", "after", "cursor-a", "first", "2"), http.StatusOK, `{"data":[{"id":"request-1","broadcaster_id":"1234","broadcaster_name":"Broadcaster","broadcaster_login":"broadcaster","moderator_id":"5678","moderator_login":"mod","moderator_name":"Mod","user_id":"9876","user_login":"viewer","user_name":"Viewer","text":"please unban","status":"pending","created_at":"2026-01-02T03:04:05Z","resolved_at":null,"resolution_text":null}],"pagination":{"cursor":"cursor-b"}}`),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorReadUnbanRequests),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.GetUnbanRequests(context.Background(), helix.GetUnbanRequestsRequest{BroadcasterID: "1234", ModeratorID: "5678", Status: "pending", UserID: "9876", After: moderationString("cursor-a"), First: moderationInt(2)}))
			},
		},
		{
			anchor:     "resolve-unban-requests",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "unban_request_id", "request-1", "status", "approved", "resolution_text", "welcome back"), http.StatusOK, `{"data":[{"id":"request-1","status":"approved","resolution_text":"welcome back"}]}`),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageUnbanRequests),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.ResolveUnbanRequests(context.Background(), helix.ResolveUnbanRequestsRequest{BroadcasterID: "1234", ModeratorID: "5678", UnbanRequestID: "request-1", Status: "approved", ResolutionText: "welcome back"}))
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.anchor, func(t *testing.T) { runModerationOperation(t, testCase) })
	}
}

func TestModerationBannedUsersPager_forwardsAfterCursor(t *testing.T) {
	// Given two pages of banned users with a cursor on the first page.
	client, transport := moderationClient(t, moderationCredential(helix.TokenClassUser, "1234", helix.ScopeModerationRead),
		moderationResponse(http.StatusOK, `{"data":[{"user_id":"u1"}],"pagination":{"cursor":"cursor-b"}}`),
		moderationResponse(http.StatusOK, `{"data":[{"user_id":"u2"}],"pagination":{}}`))

	// When the pager is consumed.
	pager, err := client.Moderation.GetBannedUsersPager(helix.GetBannedUsersRequest{BroadcasterID: "1234", After: moderationString("cursor-a")})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || !pager.Next(context.Background()) || pager.Next(context.Background()) {
		t.Fatal("pager did not yield exactly two pages")
	}

	// Then the second request carries the response cursor.
	requests := transport.Requests()
	if len(requests) != 2 || requests[1].Path != "/helix/moderation/banned?after=cursor-b&broadcaster_id=1234" {
		t.Fatalf("pager requests = %#v", requests)
	}
}

func moderationResponse(status int, body string) testkit.RoundTripResponse {
	return testkit.RoundTripResponse{StatusCode: status, Header: http.Header{"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {verticalSliceRateReset}}, Body: body}
}
