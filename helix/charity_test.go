package helix_test

import (
	"context"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestCharityGetCharityCampaign(t *testing.T) {
	body := `{"data":[{"id":"campaign-1","broadcaster_id":"123456","broadcaster_login":"streamer","broadcaster_name":"Streamer","charity_name":"Example Charity","charity_description":"Help","charity_logo":"https://logo.test/logo.png","charity_website":"https://charity.test","current_amount":{"value":500,"decimal_places":2,"currency":"USD"},"target_amount":null}]}`
	transport := testkit.NewRecordingRoundTripper(task15Response(body))
	client := task15Client(t, transport)
	result, err := client.Charity.GetCharityCampaign(context.Background(), helix.GetCharityCampaignRequest{BroadcasterID: "123456"})
	fixture := task15Fixture(urlValues("broadcaster_id", "123456"), "", task15Success(body))
	task15Contract(t, "get-charity-campaign", fixture, transport, result.Meta, err)
	if len(result.Data) != 1 || result.Data[0].ID != "campaign-1" || result.Data[0].TargetAmount != nil {
		t.Fatalf("campaigns = %#v", result.Data)
	}
}

func TestCharityGetCharityCampaignDonationsPager(t *testing.T) {
	first := `{"data":[{"id":"donation-1","campaign_id":"campaign-1","user_id":"9","user_login":"viewer","user_name":"Viewer","amount":{"value":500,"decimal_places":2,"currency":"USD"}}],"pagination":{"cursor":"next"}}`
	second := `{"data":[{"id":"donation-2","campaign_id":"campaign-1","user_id":"10","user_login":"viewer2","user_name":"Viewer2","amount":{"value":1000,"decimal_places":2,"currency":"USD"}}]}`
	transport := testkit.NewRecordingRoundTripper(task15Response(first), task15Response(second))
	client := task15Client(t, transport)
	pager, err := client.Charity.GetCharityCampaignDonationsPager(helix.GetCharityCampaignDonationsRequest{BroadcasterID: "123456", First: intPointer(1)})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || pager.Page().Data[0].ID != "donation-1" || !pager.Next(context.Background()) || pager.Page().Data[0].ID != "donation-2" || pager.Next(context.Background()) || pager.Err() != nil {
		t.Fatalf("donation pager state: page=%#v err=%v", pager.Page(), pager.Err())
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[1].Path != "/helix/charity/donations?after=next&broadcaster_id=123456&first=1" {
		t.Fatalf("donation pager requests = %#v", requests)
	}
}
