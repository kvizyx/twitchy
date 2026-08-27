package helix_test

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestModerationAuthorizationRejectsBeforeNetwork(t *testing.T) {
	tests := []struct {
		name       string
		credential helix.Credential
		call       func(*helix.Client) error
	}{
		{
			name:       "missing moderator subject",
			credential: moderationCredential(helix.TokenClassUser, "wrong", helix.ScopeModeratorManageWarnings),
			call: func(client *helix.Client) error {
				_, err := client.Moderation.WarnChatUser(context.Background(), helix.WarnChatUserRequest{BroadcasterID: "1234", ModeratorID: "5678"})
				return err
			},
		},
		{
			name:       "missing scope",
			credential: moderationCredential(helix.TokenClassUser, "5678"),
			call: func(client *helix.Client) error {
				_, err := client.Moderation.DeleteChatMessages(context.Background(), helix.DeleteChatMessagesRequest{BroadcasterID: "1234", ModeratorID: "5678"})
				return err
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper()
			client, err := helix.New(helix.WithBaseURL("https://api.twitch.test/helix"), helix.WithHTTPClient(&http.Client{Transport: transport}), helix.WithStaticToken(testCase.credential))
			if err != nil {
				t.Fatal(err)
			}
			callErr := testCase.call(client)
			var authErr *helix.AuthError
			if !errors.As(callErr, &authErr) || len(transport.Requests()) != 0 {
				t.Fatalf("error = %T %v, requests = %d", callErr, callErr, len(transport.Requests()))
			}
		})
	}
}

func TestModerationMutations_doNotReplayAfterTransientFailures(t *testing.T) {
	tests := []struct {
		name string
		call func(*helix.Client) error
	}{
		{"ban", func(client *helix.Client) error {
			_, err := client.Moderation.BanUser(context.Background(), helix.BanUserRequest{BroadcasterID: "1234", ModeratorID: "5678", Data: helix.BanUserBody{UserID: "9876"}})
			return err
		}},
		{"unban", func(client *helix.Client) error {
			_, err := client.Moderation.UnbanUser(context.Background(), helix.UnbanUserRequest{BroadcasterID: "1234", ModeratorID: "5678", UserID: "9876"})
			return err
		}},
		{"delete chat", func(client *helix.Client) error {
			_, err := client.Moderation.DeleteChatMessages(context.Background(), helix.DeleteChatMessagesRequest{BroadcasterID: "1234", ModeratorID: "5678"})
			return err
		}},
		{"suspicious", func(client *helix.Client) error {
			_, err := client.Experimental.Moderation.RemoveSuspiciousStatusFromChatUser(context.Background(), helix.RemoveSuspiciousStatusFromChatUserRequest{BroadcasterID: "1234", ModeratorID: "5678", UserID: "9876"})
			return err
		}},
	}
	for _, status := range []int{http.StatusServiceUnavailable, http.StatusUnauthorized} {
		for _, testCase := range tests {
			t.Run(testCase.name+"/"+http.StatusText(status), func(t *testing.T) {
				credential := moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageBannedUsers, helix.ScopeModeratorManageChatMessages, helix.ScopeModeratorManageSuspiciousUsers)
				body := `{"error":"failure","status":` + strconv.Itoa(status) + `,"message":"retry"}`
				client, transport := moderationClient(t, credential, moderationResponse(status, body), moderationResponse(http.StatusOK, `{"data":[]}`))
				_ = testCase.call(client)
				if len(transport.Requests()) != 1 {
					t.Fatalf("mutation attempts = %d, want 1", len(transport.Requests()))
				}
			})
		}
	}
}

func TestModerationTypedErrors_preserveStatusMetadata(t *testing.T) {
	// Given a structured API failure and malformed success response.
	tests := []struct {
		name string
		resp testkit.RoundTripResponse
		call func(*helix.Client) error
		want func(error) bool
	}{
		{
			name: "api error",
			resp: moderationResponse(http.StatusBadRequest, `{"error":"Bad Request","status":400,"message":"invalid"}`),
			call: func(client *helix.Client) error {
				_, err := client.Moderation.GetShieldModeStatus(context.Background(), helix.GetShieldModeStatusRequest{BroadcasterID: "1234", ModeratorID: "5678"})
				return err
			},
			want: func(err error) bool { var typed *helix.APIError; return errors.As(err, &typed) },
		},
		{
			name: "protocol error",
			resp: moderationResponse(http.StatusOK, "not-json"),
			call: func(client *helix.Client) error {
				_, err := client.Moderation.GetShieldModeStatus(context.Background(), helix.GetShieldModeStatusRequest{BroadcasterID: "1234", ModeratorID: "5678"})
				return err
			},
			want: func(err error) bool { var typed *helix.ProtocolError; return errors.As(err, &typed) },
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			client, _ := moderationClient(t, moderationCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorReadShieldMode), testCase.resp)
			callErr := testCase.call(client)
			if !testCase.want(callErr) {
				t.Fatalf("error = %T %v", callErr, callErr)
			}
			var metaErr interface{ Meta() helix.ResponseMeta }
			if !errors.As(callErr, &metaErr) || metaErr.Meta().StatusCode() != testCase.resp.StatusCode {
				t.Fatalf("error metadata = %#v", callErr)
			}
		})
	}
}
