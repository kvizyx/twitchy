package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestScheduleGetChannelStreamSchedule_preservesScheduleFields(t *testing.T) {
	body := `{"data":{"segments":[{"id":"segment-1","start_time":"2026-07-27T18:00:00Z","end_time":"2026-07-27T19:00:00Z","title":"Weekly","canceled_until":null,"category":{"id":"game-1","name":"Game"},"is_recurring":true}],"broadcaster_id":"broadcaster-1","broadcaster_name":"Broadcaster","broadcaster_login":"broadcaster","vacation":{"start_time":"2026-08-01T00:00:00Z","end_time":"2026-08-08T00:00:00Z"}},"pagination":{"cursor":"next"}}`
	transport := testkit.NewRecordingRoundTripper(task24Response(http.StatusOK, body))
	client := task24Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})

	result, err := client.Schedule.GetChannelStreamSchedule(context.Background(), helix.GetChannelStreamScheduleRequest{
		BroadcasterID: "broadcaster-1",
		IDs:           []string{"segment-1", "segment-2"},
		StartTime:     task24String("2026-07-27T00:00:00Z"),
		UTCOffset:     task24String("-04:00"),
		First:         task24Int(25),
		After:         task24String("prior"),
	})
	if err != nil {
		t.Fatal(err)
	}

	requests := transport.Requests()
	if len(requests) != 1 || requests[0].Path != "/helix/schedule?after=prior&broadcaster_id=broadcaster-1&first=25&id=segment-1&id=segment-2&start_time=2026-07-27T00%3A00%3A00Z&utc_offset=-04%3A00" {
		t.Fatalf("schedule request = %#v", requests)
	}
	if len(result.Data.Segments) != 1 || !result.Data.Segments[0].IsRecurring || result.Data.Segments[0].Category.Name != "Game" {
		t.Fatalf("schedule data = %#v", result.Data)
	}
	if result.Data.Segments[0].CanceledUntil != nil || result.Data.Vacation == nil || result.Pagination.Cursor() != "next" {
		t.Fatalf("schedule nullable data = %#v", result.Data)
	}
}

func TestScheduleGetChannelStreamSchedulePager_forwardsAfterCursor(t *testing.T) {
	first := `{"data":{"segments":[],"broadcaster_id":"broadcaster-1","broadcaster_name":"Broadcaster","broadcaster_login":"broadcaster","vacation":null},"pagination":{"cursor":"older"}}`
	second := `{"data":{"segments":[],"broadcaster_id":"broadcaster-1","broadcaster_name":"Broadcaster","broadcaster_login":"broadcaster","vacation":null},"pagination":{}}`
	transport := testkit.NewRecordingRoundTripper(task24Response(http.StatusOK, first), task24Response(http.StatusOK, second))
	client := task24Client(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	pager, err := client.Schedule.GetChannelStreamSchedulePager(helix.GetChannelStreamScheduleRequest{BroadcasterID: "broadcaster-1"})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) {
		t.Fatal("first schedule page failed")
	}
	if !pager.Next(context.Background()) {
		t.Fatal("second schedule page failed")
	}
	if pager.Next(context.Background()) || pager.Err() != nil {
		t.Fatalf("pager state: page=%#v err=%v", pager.Page(), pager.Err())
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[0].Path != "/helix/schedule?broadcaster_id=broadcaster-1" || requests[1].Path != "/helix/schedule?after=older&broadcaster_id=broadcaster-1" {
		t.Fatalf("schedule pager requests = %#v", requests)
	}
}

func TestScheduleGetChannelICalendar_isUnauthenticatedRawBytesWithMetadata(t *testing.T) {
	calendar := []byte("BEGIN:VCALENDAR\r\nX-RAW:{not-json}\r\nEND:VCALENDAR\r\n")
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {"text/calendar; charset=utf-8"}, "X-Request-ID": {"calendar-request"}},
		Body:       string(calendar),
	})
	client := task24Client(t, transport, helix.Credential{AccessToken: "must-not-be-sent", ClientID: "must-not-be-sent", TokenClass: helix.TokenClassUser})

	result, err := client.Schedule.GetChannelICalendar(context.Background(), helix.GetChannelICalendarRequest{BroadcasterID: "broadcaster-1"})
	if err != nil {
		t.Fatal(err)
	}
	requests := transport.Requests()
	if len(requests) != 1 || requests[0].Path != "/helix/schedule/icalendar?broadcaster_id=broadcaster-1" {
		t.Fatalf("calendar request = %#v", requests)
	}
	if requests[0].Header.Get("Authorization") != "" || requests[0].Header.Get("Client-Id") != "" {
		t.Fatalf("calendar auth headers = %#v", requests[0].Header)
	}
	if string(result.Data) != string(calendar) || result.Meta.Header().Get("Content-Type") != "text/calendar; charset=utf-8" || result.Meta.RequestID() != "calendar-request" {
		t.Fatalf("calendar result = %#v", result)
	}
}

func TestScheduleUpdateChannelStreamSchedule_omitsUnsetVacationFields(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(task24Response(http.StatusNoContent, ""))
	client := task24Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "broadcaster-1", Scopes: []helix.AuthorizationScope{helix.ScopeChannelManageSchedule}})

	_, err := client.Schedule.UpdateChannelStreamSchedule(context.Background(), helix.UpdateChannelStreamScheduleRequest{
		BroadcasterID:     "broadcaster-1",
		IsVacationEnabled: task24Bool(false),
	})
	if err != nil {
		t.Fatal(err)
	}
	requests := transport.Requests()
	if len(requests) != 1 || requests[0].Method != http.MethodPatch || requests[0].Path != "/helix/schedule/settings?broadcaster_id=broadcaster-1&is_vacation_enabled=false" {
		t.Fatalf("schedule settings request = %#v", requests)
	}
}

func TestScheduleSegmentRequests_preserveRecurringAndCancellationWireFields(t *testing.T) {
	responses := []testkit.RoundTripResponse{
		task24Response(http.StatusOK, `{"data":{"segments":[],"broadcaster_id":"broadcaster-1","broadcaster_name":"Broadcaster","broadcaster_login":"broadcaster","vacation":null}}`),
		task24Response(http.StatusOK, `{"data":{"segments":[],"broadcaster_id":"broadcaster-1","broadcaster_name":"Broadcaster","broadcaster_login":"broadcaster","vacation":null}}`),
		task24Response(http.StatusNoContent, ""),
	}
	transport := testkit.NewRecordingRoundTripper(responses...)
	client := task24Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "broadcaster-1", Scopes: []helix.AuthorizationScope{helix.ScopeChannelManageSchedule}})

	_, err := client.Schedule.CreateChannelStreamScheduleSegment(context.Background(), helix.CreateChannelStreamScheduleSegmentRequest{
		BroadcasterID: "broadcaster-1",
		StartTime:     "2026-07-27T18:00:00Z",
		Timezone:      "America/New_York",
		Duration:      "60",
		IsRecurring:   task24Bool(true),
		CategoryID:    task24String("game-1"),
		Title:         task24String("Weekly"),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Schedule.UpdateChannelStreamScheduleSegment(context.Background(), helix.UpdateChannelStreamScheduleSegmentRequest{
		BroadcasterID: "broadcaster-1",
		ID:            "segment-1",
		IsCanceled:    task24Bool(false),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Schedule.DeleteChannelStreamScheduleSegment(context.Background(), helix.DeleteChannelStreamScheduleSegmentRequest{BroadcasterID: "broadcaster-1", ID: "segment-1"})
	if err != nil {
		t.Fatal(err)
	}
	requests := transport.Requests()
	if len(requests) != 3 {
		t.Fatalf("segment request count = %d", len(requests))
	}
	if string(requests[0].Body) != `{"start_time":"2026-07-27T18:00:00Z","timezone":"America/New_York","duration":"60","is_recurring":true,"category_id":"game-1","title":"Weekly"}` {
		t.Fatalf("create segment body = %q", requests[0].Body)
	}
	if string(requests[1].Body) != `{"is_canceled":false}` {
		t.Fatalf("false cancellation body = %q", requests[1].Body)
	}
	if requests[2].Method != http.MethodDelete || requests[2].Path != "/helix/schedule/segment?broadcaster_id=broadcaster-1&id=segment-1" {
		t.Fatalf("delete segment request = %#v", requests[2])
	}
}

func TestScheduleGetChannelICalendar_doesNotRequireCredential(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: http.StatusOK, Body: "calendar"})
	client, err := helix.New(helix.WithBaseURL("https://api.twitch.test/helix"), helix.WithHTTPClient(&http.Client{Transport: transport}))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Schedule.GetChannelICalendar(context.Background(), helix.GetChannelICalendarRequest{BroadcasterID: "broadcaster-1"})
	if err != nil {
		t.Fatal(err)
	}
}

func task24Client(t *testing.T, transport *testkit.RecordingRoundTripper, credential helix.Credential) *helix.Client {
	t.Helper()
	client, err := helix.New(helix.WithBaseURL("https://api.twitch.test/helix"), helix.WithHTTPClient(&http.Client{Transport: transport}), helix.WithStaticToken(credential))
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func task24Response(status int, body string) testkit.RoundTripResponse {
	return testkit.RoundTripResponse{StatusCode: status, Header: http.Header{"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {"4102444800"}}, Body: body}
}

func task24String(value string) *string { return &value }
func task24Int(value int) *int          { return &value }
func task24Bool(value bool) *bool       { return &value }
