package helix_test

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestPollsCreatePoll_preservesExactWireAndNestedChoices(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: http.StatusOK, Header: task23RateHeaders(), Body: task23Fixture(t, "poll.json")})
	client := task23Client(t, transport)
	result, err := client.Polls.CreatePoll(context.Background(), helix.CreatePollRequest{
		BroadcasterID:              "123456",
		Title:                      "Heads or tails?",
		Choices:                    []helix.PollChoice{{Title: "Heads"}, {Title: "Tails"}},
		Duration:                   600,
		ChannelPointsVotingEnabled: boolPointer(true),
		ChannelPointsPerVote:       intPointer(100),
	})
	if err != nil {
		t.Fatal(err)
	}
	request := transport.Requests()[0]
	if request.Method != http.MethodPost || request.Path != "/helix/polls" || string(request.Body) != `{"broadcaster_id":"123456","title":"Heads or tails?","choices":[{"title":"Heads"},{"title":"Tails"}],"duration":600,"channel_points_voting_enabled":true,"channel_points_per_vote":100}` {
		t.Fatalf("poll request = %#v", request)
	}
	if result.Data[0].Choices[0].Title != "Heads" || result.Data[0].ChannelPointsPerVote != 100 || result.Data[0].Status != helix.PollStatusActive {
		t.Fatalf("poll response = %#v", result.Data[0])
	}
}

func TestPollsGetPollsPager_sendsCursorAndPreservesUnknownStatus(t *testing.T) {
	first := `{"data":[{"id":"poll-1","broadcaster_id":"123456","title":"Question","choices":[{"id":"choice-1","title":"Yes","votes":2,"channel_points_votes":1,"bits_votes":0}],"status":"FUTURE_STATUS","duration":60,"started_at":"2024-01-02T03:04:05Z","ended_at":null}],"pagination":{"cursor":"next"}}`
	second := `{"data":[{"id":"poll-2","status":"ARCHIVED","started_at":"2024-01-02T03:04:05Z","ended_at":"2024-01-02T03:05:05Z"}]}`
	transport := testkit.NewRecordingRoundTripper(task23Response(first), task23Response(second))
	client := task23Client(t, transport)
	pager, err := client.Polls.GetPollsPager(helix.GetPollsRequest{BroadcasterID: "123456", First: intPointer(1)})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || pager.Page().Data[0].Status != helix.PollStatus("FUTURE_STATUS") || !pager.Next(context.Background()) || pager.Page().Data[0].Status != helix.PollStatusArchived || pager.Next(context.Background()) || pager.Err() != nil {
		t.Fatalf("poll pager state: page=%#v err=%v", pager.Page(), pager.Err())
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[1].Path != "/helix/polls?after=next&broadcaster_id=123456&first=1" {
		t.Fatalf("poll pager requests = %#v", requests)
	}
}

func TestPollsEndPoll_isOneAttemptOn503(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: http.StatusServiceUnavailable, Header: task23RateHeaders(), Body: `{"error":"Unavailable","status":503,"message":"try later"}`})
	client := task23Client(t, transport)
	_, err := client.Polls.EndPoll(context.Background(), helix.EndPollRequest{BroadcasterID: "123456", ID: "poll-1", Status: helix.PollStatusTerminated})
	var apiErr *helix.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode() != http.StatusServiceUnavailable || len(transport.Requests()) != 1 {
		t.Fatalf("poll 503 error=%T %v requests=%d", err, err, len(transport.Requests()))
	}
}

func TestPollsRequireUserTokenBoundToBroadcaster(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper()
	client := task23ClientWithUser(t, transport, "different-user")
	_, err := client.Polls.GetPolls(context.Background(), helix.GetPollsRequest{BroadcasterID: "123456"})
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) || len(transport.Requests()) != 0 {
		t.Fatalf("poll auth error=%T %v requests=%d", err, err, len(transport.Requests()))
	}
}

func TestPollsGetPolls_acceptsReadScopeWithoutManageScope(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: http.StatusOK, Header: task23RateHeaders(), Body: task23Fixture(t, "poll.json")})
	client := task23ClientWithScopes(t, transport, "123456", helix.ScopeChannelReadPolls)
	result, err := client.Polls.GetPolls(context.Background(), helix.GetPollsRequest{BroadcasterID: "123456"})
	if err != nil || len(result.Data) != 1 {
		t.Fatalf("read-only poll result=%#v err=%v", result, err)
	}
}

func task23Fixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile("testdata/task23/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
