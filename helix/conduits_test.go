package helix_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestConduits_verticalSlice(t *testing.T) {
	tests := []struct {
		name string
		call func(context.Context, *helix.Client) error
		body string
	}{
		{
			name: "get",
			call: func(ctx context.Context, client *helix.Client) error {
				result, err := client.Conduits.GetConduits(ctx, helix.GetConduitsRequest{})
				if err == nil && (len(result.Data) != 1 || result.Data[0].ShardCount != 5) {
					t.Fatalf("conduits = %#v", result.Data)
				}
				return err
			},
			body: `{"data":[{"id":"conduit-1","shard_count":5}]}`,
		},
		{
			name: "create",
			call: func(ctx context.Context, client *helix.Client) error {
				_, err := client.Conduits.CreateConduits(ctx, helix.CreateConduitsRequest{ShardCount: 5})
				return err
			},
			body: `{"data":[{"id":"conduit-1","shard_count":5}]}`,
		},
		{
			name: "update",
			call: func(ctx context.Context, client *helix.Client) error {
				_, err := client.Conduits.UpdateConduits(ctx, helix.UpdateConduitsRequest{ID: "conduit-1", ShardCount: 5})
				return err
			},
			body: `{"data":[{"id":"conduit-1","shard_count":5}]}`,
		},
		{
			name: "delete",
			call: func(ctx context.Context, client *helix.Client) error {
				_, err := client.Conduits.DeleteConduit(ctx, helix.DeleteConduitRequest{ID: "conduit-1"})
				return err
			},
			body: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {verticalSliceRateReset}},
				Body:       test.body,
			})
			if test.name == "delete" {
				transport = testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: http.StatusNoContent})
			}
			client, err := helix.New(
				helix.WithBaseURL("https://api.twitch.test/helix"),
				helix.WithHTTPClient(&http.Client{Transport: transport}),
				helix.WithStaticToken(helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp}),
			)
			if err != nil {
				t.Fatal(err)
			}
			if err := test.call(context.Background(), client); err != nil {
				t.Fatal(err)
			}
			if len(transport.Requests()) != 1 {
				t.Fatalf("requests = %d, want 1", len(transport.Requests()))
			}
		})
	}
}

func TestConduits_shardsWireAndPagination(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {verticalSliceRateReset}},
		Body:       `{"data":[{"id":"0","status":"enabled","transport":{"method":"websocket","session_id":"session-1","connected_at":"2020-11-10T14:32:18.730260295Z"}}],"pagination":{"cursor":"next"}}`,
	})
	client, err := helix.New(helix.WithBaseURL("https://api.twitch.test/helix"), helix.WithHTTPClient(&http.Client{Transport: transport}), helix.WithStaticToken(helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp}))
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.Conduits.GetConduitShards(context.Background(), helix.GetConduitShardsRequest{ConduitID: "conduit-1", Status: "enabled", After: "cursor-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Data[0].Transport.SessionID != "session-1" || result.Pagination.Cursor() != "next" {
		t.Fatalf("shards = %#v, pagination = %q", result.Data, result.Pagination.Cursor())
	}
	request := transport.Requests()[0]
	if request.Path != "/helix/eventsub/conduits/shards?after=cursor-1&conduit_id=conduit-1&status=enabled" {
		t.Fatalf("path = %q", request.Path)
	}
}

func TestConduits_updateShardsWire(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{
		StatusCode: http.StatusAccepted,
		Header:     http.Header{"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {verticalSliceRateReset}},
		Body:       `{"data":[{"id":"0","status":"enabled","transport":{"method":"webhook","callback":"https://callback.test"}}],"errors":[{"id":"1","message":"bad shard","code":"invalid_parameter"}]}`,
	})
	client, err := helix.New(helix.WithBaseURL("https://api.twitch.test/helix"), helix.WithHTTPClient(&http.Client{Transport: transport}), helix.WithStaticToken(helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp}))
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.Conduits.UpdateConduitShards(context.Background(), helix.UpdateConduitShardsRequest{
		ConduitID: "conduit-1",
		Shards:    []helix.UpdateConduitShard{{ID: "0", Transport: helix.ConduitShardTransport{Method: "webhook", Callback: "https://callback.test", Secret: "secret-value"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Data.Shards) != 1 || len(result.Data.Errors) != 1 {
		t.Fatalf("update result = %#v", result.Data)
	}
	if got := string(transport.Requests()[0].Body); got != `{"conduit_id":"conduit-1","shards":[{"id":"0","transport":{"method":"webhook","callback":"https://callback.test","secret":"secret-value"}}]}` {
		t.Fatalf("body = %s", got)
	}
}
