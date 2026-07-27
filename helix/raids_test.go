package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestRaidsStartRaid_preservesQueryAndResponse(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: http.StatusOK, Header: task23RateHeaders(), Body: task23Fixture(t, "raid.json")})
	client := task23Client(t, transport)
	result, err := client.Raids.StartRaid(context.Background(), helix.StartRaidRequest{FromBroadcasterID: "123456", ToBroadcasterID: "654321"})
	if err != nil {
		t.Fatal(err)
	}
	request := transport.Requests()[0]
	if request.Method != http.MethodPost || request.Path != "/helix/raids?from_broadcaster_id=123456&to_broadcaster_id=654321" || len(request.Body) != 0 || result.Data[0].IsMature {
		t.Fatalf("raid request=%#v response=%#v", request, result.Data)
	}
}

func TestRaidsCancelRaid_returnsEmptySuccessWithoutReplay(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: http.StatusNoContent, Header: task23RateHeaders()})
	client := task23Client(t, transport)
	result, err := client.Raids.CancelRaid(context.Background(), helix.CancelRaidRequest{BroadcasterID: "123456"})
	if err != nil || result == nil || result.Data != (helix.CancelRaidData{}) || len(transport.Requests()) != 1 || transport.Requests()[0].Path != "/helix/raids?broadcaster_id=123456" {
		t.Fatalf("cancel raid result=%#v err=%v requests=%#v", result, err, transport.Requests())
	}
}

func TestRaidsStartRaid_requiresManageScopeAndBroadcasterBinding(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper()
	client := task23ClientWithScopes(t, transport, "different-user", helix.ScopeChannelReadPolls)
	_, err := client.Raids.StartRaid(context.Background(), helix.StartRaidRequest{FromBroadcasterID: "123456", ToBroadcasterID: "654321"})
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) || len(transport.Requests()) != 0 {
		t.Fatalf("raid auth error=%T %v requests=%d", err, err, len(transport.Requests()))
	}
}
