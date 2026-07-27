package helix_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

const task15RewardJSON = `{"data":[{"broadcaster_id":"123456","broadcaster_login":"streamer","broadcaster_name":"Streamer","id":"reward-1","title":"Hydrate","prompt":"","cost":500,"image":null,"default_image":{"url_1x":"one","url_2x":"two","url_4x":"four"},"background_color":"#00E5CB","is_enabled":true,"is_user_input_required":false,"max_per_stream_setting":{"is_enabled":false,"max_per_stream":0},"max_per_user_per_stream_setting":{"is_enabled":false,"max_per_user_per_stream":0},"global_cooldown_setting":{"is_enabled":false,"global_cooldown_seconds":0},"is_paused":false,"is_in_stock":true,"should_redemptions_skip_request_queue":false,"redemptions_redeemed_current_stream":null,"cooldown_expires_at":null}]}`

func TestChannelPointsCreateCustomRewards(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(task15Response(task15RewardJSON))
	client := task15Client(t, transport)
	result, err := client.ChannelPoints.CreateCustomRewards(context.Background(), helix.CreateCustomRewardsRequest{BroadcasterID: "123456", Title: "Hydrate", Cost: 500})
	fixture := task15Fixture(urlValues("broadcaster_id", "123456"), `{"title":"Hydrate","cost":500}`, task15Success(task15RewardJSON))
	task15Contract(t, "create-custom-rewards", fixture, transport, result.Meta, err)
	if len(result.Data) != 1 || result.Data[0].ID != "reward-1" || result.Data[0].Cost != 500 {
		t.Fatalf("custom rewards = %#v", result.Data)
	}
}

func TestChannelPointsDeleteCustomReward(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(task15NoContent()))
	client := task15Client(t, transport)
	result, err := client.ChannelPoints.DeleteCustomReward(context.Background(), helix.DeleteCustomRewardRequest{BroadcasterID: "123456", ID: "reward-1"})
	fixture := task15Fixture(urlValues("broadcaster_id", "123456", "id", "reward-1"), "", task15NoContent())
	task15Contract(t, "delete-custom-reward", fixture, transport, result.Meta, err)
}

func TestChannelPointsGetCustomReward(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(task15Response(task15RewardJSON))
	client := task15Client(t, transport)
	result, err := client.ChannelPoints.GetCustomReward(context.Background(), helix.GetCustomRewardRequest{BroadcasterID: "123456", IDs: []string{"reward-1", "reward-2"}, OnlyManageableRewards: boolPointer(true)})
	fixture := task15Fixture(urlValues("broadcaster_id", "123456", "id", "reward-1", "id", "reward-2", "only_manageable_rewards", "true"), "", task15Success(task15RewardJSON))
	task15Contract(t, "get-custom-reward", fixture, transport, result.Meta, err)
}

func TestChannelPointsGetCustomRewardRedemptionPager(t *testing.T) {
	first := `{"data":[{"broadcaster_id":"123456","broadcaster_login":"streamer","broadcaster_name":"Streamer","id":"redemption-1","user_login":"viewer","user_id":"9","user_name":"Viewer","user_input":"hello","status":"UNFULFILLED","redeemed_at":"2024-01-02T03:04:05Z","reward":{"id":"reward-1","title":"Hydrate","prompt":"","cost":500}}],"pagination":{"cursor":"next"}}`
	second := `{"data":[{"broadcaster_id":"123456","id":"redemption-2","status":"FULFILLED","redeemed_at":"2024-01-03T03:04:05Z","reward":{"id":"reward-1","title":"Hydrate","prompt":"","cost":500}}]}`
	transport := testkit.NewRecordingRoundTripper(task15Response(first), task15Response(second))
	client := task15Client(t, transport)
	pager, err := client.ChannelPoints.GetCustomRewardRedemptionPager(helix.GetCustomRewardRedemptionRequest{BroadcasterID: "123456", RewardID: "reward-1", Status: helix.RedemptionStatusUnfulfilled, IDs: []string{"redemption-1", "redemption-2"}, First: intPointer(1)})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || pager.Page().Data[0].ID != "redemption-1" || !pager.Next(context.Background()) || pager.Page().Data[0].ID != "redemption-2" || pager.Next(context.Background()) || pager.Err() != nil {
		t.Fatalf("redemption pager state: page=%#v err=%v", pager.Page(), pager.Err())
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[1].Path != "/helix/channel_points/custom_rewards/redemptions?after=next&broadcaster_id=123456&first=1&id=redemption-1&id=redemption-2&reward_id=reward-1&status=UNFULFILLED" {
		t.Fatalf("redemption pager requests = %#v", requests)
	}
}

func TestChannelPointsUpdateCustomReward_preservesExplicitZeroFalseAndEmpty(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(task15Response(task15RewardJSON))
	client := task15Client(t, transport)
	result, err := client.ChannelPoints.UpdateCustomReward(context.Background(), helix.UpdateCustomRewardRequest{
		BroadcasterID: "123456",
		ID:            "reward-1",
		Title:         stringPointer(""),
		Cost:          int64Pointer(0),
		IsEnabled:     boolPointer(false),
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := string(transport.Requests()[0].Body); got != `{"title":"","cost":0,"is_enabled":false}` {
		t.Fatalf("PATCH body = %q", got)
	}
	if len(result.Data) != 1 {
		t.Fatalf("updated rewards = %#v", result.Data)
	}
}

func TestChannelPointsUpdateCustomReward_omitsUnsetFields(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(task15Response(task15RewardJSON))
	client := task15Client(t, transport)
	_, err := client.ChannelPoints.UpdateCustomReward(context.Background(), helix.UpdateCustomRewardRequest{BroadcasterID: "123456", ID: "reward-1"})
	if err != nil {
		t.Fatal(err)
	}
	if got := string(transport.Requests()[0].Body); got != `{}` {
		t.Fatalf("omitted PATCH body = %q", got)
	}
}

func TestChannelPointsUpdateRedemptionStatus(t *testing.T) {
	body := `{"data":[{"broadcaster_id":"123456","broadcaster_login":"streamer","broadcaster_name":"Streamer","id":"redemption-1","user_id":"9","user_name":"Viewer","user_login":"viewer","reward":{"id":"reward-1","title":"Hydrate","prompt":"","cost":500},"user_input":"","status":"FULFILLED","redeemed_at":"2024-01-02T03:04:05Z"}]}`
	transport := testkit.NewRecordingRoundTripper(task15Response(body))
	client := task15Client(t, transport)
	result, err := client.ChannelPoints.UpdateRedemptionStatus(context.Background(), helix.UpdateRedemptionStatusRequest{IDs: []string{"redemption-1", "redemption-2"}, BroadcasterID: "123456", RewardID: "reward-1", Status: helix.RedemptionStatusFulfilled})
	fixture := task15Fixture(urlValues("id", "redemption-1", "id", "redemption-2", "broadcaster_id", "123456", "reward_id", "reward-1"), `{"status":"FULFILLED"}`, task15Success(body))
	task15Contract(t, "update-redemption-status", fixture, transport, result.Meta, err)
}

func TestChannelPoints_userOperationsRejectAppTokenBeforeNetwork(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper()
	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp}),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.ChannelPoints.GetCustomReward(context.Background(), helix.GetCustomRewardRequest{BroadcasterID: "123456"})
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("error = %T %v, want AuthError", err, err)
	}
	if len(transport.Requests()) != 0 {
		t.Fatal("app token reached channel points network")
	}
}

func TestChannelPoints_typedProtocolAndAPIErrors(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		body       string
		wantStatus int
		wantType   any
	}{
		{name: "malformed success", status: http.StatusOK, body: "not-json", wantStatus: http.StatusOK, wantType: &helix.ProtocolError{}},
		{name: "api error", status: http.StatusNotFound, body: `{"error":"Not Found","status":404,"message":"missing"}`, wantStatus: http.StatusNotFound, wantType: &helix.APIError{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: test.status, Header: task15RateHeaders(), Body: test.body})
			client := task15Client(t, transport)
			_, err := client.ChannelPoints.GetCustomReward(context.Background(), helix.GetCustomRewardRequest{BroadcasterID: "123456"})
			if !errors.As(err, &test.wantType) {
				t.Fatalf("error = %T %v, want %T", err, err, test.wantType)
			}
		})
	}
}

func TestChannelPoints_updateRequestJSONIsObject(t *testing.T) {
	value := helix.UpdateCustomRewardRequest{BroadcasterID: "123456", ID: "reward-1", Cost: int64Pointer(0)}
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != `{"cost":0}` {
		t.Fatalf("request JSON = %s", encoded)
	}
}
