package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestChannelsGetChannelInformation(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(channelResponse(`{"data":[{"broadcaster_id":"111","broadcaster_login":"alpha","broadcaster_name":"Alpha","broadcaster_language":"en","game_name":"Game","game_id":"42","title":"Live","delay":0,"tags":["tag"],"content_classification_labels":["Gambling"],"is_branded_content":false}]}`))
	client := channelClient(t, transport)

	result, err := client.Channels.GetChannelInformation(context.Background(), helix.GetChannelInformationRequest{BroadcasterIDs: []string{"111", "222"}})
	fixture := channelFixture(urlValues("broadcaster_id", "111", "broadcaster_id", "222"), "", channelSuccess(`{"data":[{"broadcaster_id":"111","broadcaster_login":"alpha","broadcaster_name":"Alpha","broadcaster_language":"en","game_name":"Game","game_id":"42","title":"Live","delay":0,"tags":["tag"],"content_classification_labels":["Gambling"],"is_branded_content":false}]}`))
	if err != nil {
		t.Fatal(err)
	}
	channelContract(t, "get-channel-information", fixture, transport, result.Meta, nil)
	if len(result.Data) != 1 || result.Data[0].GameID != "42" || result.Data[0].Delay != 0 {
		t.Fatalf("channel information = %#v", result.Data)
	}
}

func TestChannelsModifyChannelInformation_preservesExplicitZeroFalseAndEmpty(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(channelNoContent()))
	client := channelClient(t, transport)
	tags := []string{}
	labels := []helix.ChannelClassificationLabel{}
	request := helix.ModifyChannelInformationRequest{
		BroadcasterID:               "123456",
		GameID:                      stringPointer(""),
		BroadcasterLanguage:         stringPointer(""),
		Title:                       stringPointer("title"),
		Delay:                       intPointer(0),
		Tags:                        &tags,
		ContentClassificationLabels: &labels,
		IsBrandedContent:            boolPointer(false),
	}

	result, err := client.Channels.ModifyChannelInformation(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	requests := transport.Requests()
	if len(requests) != 1 || string(requests[0].Body) != `{"game_id":"","broadcaster_language":"","title":"title","delay":0,"tags":[],"content_classification_labels":[],"is_branded_content":false}` {
		t.Fatalf("PATCH body = %q", requests[0].Body)
	}
	if result.Meta.StatusCode() != http.StatusNoContent {
		t.Fatalf("status = %d", result.Meta.StatusCode())
	}
}

func TestChannelsModifyChannelInformation_omitsUnsetFields(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(channelNoContent()))
	client := channelClient(t, transport)
	_, err := client.Channels.ModifyChannelInformation(context.Background(), helix.ModifyChannelInformationRequest{BroadcasterID: "123456"})
	if err != nil {
		t.Fatal(err)
	}
	if got := string(transport.Requests()[0].Body); got != `{}` {
		t.Fatalf("omitted PATCH body = %q", got)
	}
}

func TestChannelsGetChannelEditors(t *testing.T) {
	body := `{"data":[{"user_id":"9","user_name":"editor","created_at":"2024-01-02T03:04:05Z"}]}`
	transport := testkit.NewRecordingRoundTripper(channelResponse(body))
	client := channelClient(t, transport)
	result, err := client.Channels.GetChannelEditors(context.Background(), helix.GetChannelEditorsRequest{BroadcasterID: "123456"})
	fixture := channelFixture(urlValues("broadcaster_id", "123456"), "", channelSuccess(body))
	channelContract(t, "get-channel-editors", fixture, transport, result.Meta, err)
	if len(result.Data) != 1 || result.Data[0].UserID != "9" {
		t.Fatalf("editors = %#v", result.Data)
	}
}

func TestChannelsGetFollowedChannelsPager(t *testing.T) {
	first := `{"data":[{"broadcaster_id":"9","broadcaster_login":"alpha","broadcaster_name":"Alpha","followed_at":"2024-01-02T03:04:05Z"}],"pagination":{"cursor":"next"}}`
	second := `{"data":[{"broadcaster_id":"10","broadcaster_login":"beta","broadcaster_name":"Beta","followed_at":"2024-01-03T03:04:05Z"}]}`
	transport := testkit.NewRecordingRoundTripper(channelResponse(first), channelResponse(second))
	client := channelClient(t, transport)
	pager, err := client.Channels.GetFollowedChannelsPager(helix.GetFollowedChannelsRequest{UserID: "123456", First: intPointer(1)})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || pager.Page().Data[0].BroadcasterID != "9" || !pager.Next(context.Background()) || pager.Page().Data[0].BroadcasterID != "10" || pager.Next(context.Background()) || pager.Err() != nil {
		t.Fatalf("pager state: page=%#v err=%v", pager.Page(), pager.Err())
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[1].Path != "/helix/channels/followed?after=next&first=1&user_id=123456" {
		t.Fatalf("pager requests = %#v", requests)
	}
}

func TestChannelsGetChannelFollowers_rejectsMismatchedUserIDBeforeNetwork(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper()
	client := channelClient(t, transport)
	_, err := client.Channels.GetChannelFollowers(context.Background(), helix.GetChannelFollowersRequest{UserID: stringPointer("wrong"), BroadcasterID: "123456"})
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("error = %T %v, want AuthError", err, err)
	}
	if len(transport.Requests()) != 0 {
		t.Fatal("mismatched subject reached network")
	}
}
