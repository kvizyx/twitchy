package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
)

func TestModerationTermsAndChatOperations(t *testing.T) {
	tests := []moderationOperationCase{
		{
			anchor:     "get-blocked-terms",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "first", "2", "after", "cursor-a"), http.StatusOK, `{"data":[{"broadcaster_id":"1234","moderator_id":"5678","id":"term-1","text":"bad phrase","created_at":"2026-01-02T03:04:05Z","updated_at":"2026-01-02T03:04:05Z","expires_at":null}],"pagination":{"cursor":"cursor-b"}}`),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorReadBlockedTerms),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.GetBlockedTerms(context.Background(), helix.GetBlockedTermsRequest{BroadcasterID: "1234", ModeratorID: "5678", First: moderationInt(2), After: moderationString("cursor-a")}))
			},
		},
		{
			anchor:     "add-blocked-term",
			fixture:    moderationBodyFixture(moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678"), http.StatusOK, `{"data":[{"id":"term-1","text":"bad phrase"}]}`), `{"text":"bad phrase"}`),
			credential: moderationCredential(helix.TokenClassApp, "", helix.ScopeModeratorManageBlockedTerms),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.AddBlockedTerm(context.Background(), helix.AddBlockedTermRequest{BroadcasterID: "1234", ModeratorID: "5678", Text: "bad phrase"}))
			},
		},
		{
			anchor:     "remove-blocked-term",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "id", "term-1"), http.StatusNoContent, ""),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageBlockedTerms),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.RemoveBlockedTerm(context.Background(), helix.RemoveBlockedTermRequest{BroadcasterID: "1234", ModeratorID: "5678", ID: "term-1"}))
			},
		},
		{
			anchor:     "delete-chat-messages",
			fixture:    moderationSuccessFixture(urlValues("broadcaster_id", "1234", "moderator_id", "5678", "message_id", "message-1"), http.StatusNoContent, ""),
			credential: moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageChatMessages),
			call: func(client *helix.Client) (helix.ResponseMeta, error) {
				return moderationMeta(client.Moderation.DeleteChatMessages(context.Background(), helix.DeleteChatMessagesRequest{BroadcasterID: "1234", ModeratorID: "5678", MessageID: "message-1"}))
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.anchor, func(t *testing.T) { runModerationOperation(t, testCase) })
	}
}

func TestModerationTermsPager_forwardsAfterCursor(t *testing.T) {
	// Given two pages of blocked terms with a cursor on the first page.
	client, transport := moderationClient(t, moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorReadBlockedTerms),
		moderationResponse(http.StatusOK, `{"data":[{"id":"term-1"}],"pagination":{"cursor":"cursor-b"}}`),
		moderationResponse(http.StatusOK, `{"data":[{"id":"term-2"}],"pagination":{}}`))

	// When the pager is consumed.
	pager, err := client.Moderation.GetBlockedTermsPager(helix.GetBlockedTermsRequest{BroadcasterID: "1234", ModeratorID: "5678"})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || !pager.Next(context.Background()) || pager.Next(context.Background()) {
		t.Fatal("pager did not yield exactly two pages")
	}

	// Then the second request carries the response cursor.
	requests := transport.Requests()
	if len(requests) != 2 || requests[1].Path != "/helix/moderation/blocked_terms?after=cursor-b&broadcaster_id=1234&moderator_id=5678" {
		t.Fatalf("pager requests = %#v", requests)
	}
}
