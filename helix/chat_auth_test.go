package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestChatAuthorizationRejectsUnsupportedCredentialsBeforeNetwork(t *testing.T) {
	// Given Chat credentials that violate token, scope, subject, or omission rules.
	tests := []struct {
		name       string
		credential helix.Credential
		call       func(*helix.Client) error
	}{
		{
			name:       "chatters subject mismatch",
			credential: chatCredential(helix.TokenClassUser, "wrong-user", helix.ScopeModeratorReadChatters),
			call: func(client *helix.Client) error {
				_, err := client.Chat.GetChatters(context.Background(), helix.GetChattersRequest{BroadcasterID: "1234", ModeratorID: "5678"})
				return err
			},
		},
		{
			name:       "user emotes missing scope",
			credential: chatCredential(helix.TokenClassUser, "5678"),
			call: func(client *helix.Client) error {
				_, err := client.Chat.GetUserEmotes(context.Background(), helix.GetUserEmotesRequest{UserID: "5678"})
				return err
			},
		},
		{
			name:       "user announcement for source only is forbidden",
			credential: chatCredential(helix.TokenClassUser, "5678", helix.ScopeModeratorManageAnnouncements),
			call: func(client *helix.Client) error {
				_, err := client.Chat.SendChatAnnouncement(context.Background(), helix.SendChatAnnouncementRequest{BroadcasterID: "1234", ModeratorID: "5678", Message: "hello", ForSourceOnly: chatBool(false)})
				return err
			},
		},
		{
			name:       "pinned message rejects extension token",
			credential: chatCredential(helix.TokenClassExtension, "5678", helix.ScopeModeratorManageChatMessages),
			call: func(client *helix.Client) error {
				_, err := client.Experimental.Chat.PinChatMessage(context.Background(), helix.PinChatMessageRequest{BroadcasterID: "1234", ModeratorID: "5678", MessageID: "message-1"})
				return err
			},
		},
		{
			name:       "pinned send missing manage scope",
			credential: chatCredential(helix.TokenClassUser, "5678", helix.ScopeUserWriteChat),
			call: func(client *helix.Client) error {
				_, err := client.Chat.SendChatMessage(context.Background(), helix.SendChatMessageRequest{BroadcasterID: "1234", SenderID: "5678", Message: "hello", Pin: chatBool(true)})
				return err
			},
		},
	}

	// When each invalid call is made.
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper()
			client, err := helix.New(
				helix.WithBaseURL("https://api.twitch.test/helix"),
				helix.WithHTTPClient(&http.Client{Transport: transport}),
				helix.WithStaticToken(testCase.credential),
			)
			if err != nil {
				t.Fatal(err)
			}
			callErr := testCase.call(client)
			var authErr *helix.AuthError
			if !errors.As(callErr, &authErr) {
				t.Fatalf("error = %T %v, want AuthError", callErr, callErr)
			}
			if len(transport.Requests()) != 0 {
				t.Fatal("invalid Chat credential reached the network")
			}
		})
	}
}

func TestChatTypedFailuresPreserveResponseMetadata(t *testing.T) {
	// Given malformed success and a structured API failure from Chat.
	tests := []struct {
		name       string
		response   testkit.RoundTripResponse
		call       func(*helix.Client) error
		wantStatus int
		wantType   func(error) bool
	}{
		{
			name:     "malformed success",
			response: testkit.RoundTripResponse{StatusCode: http.StatusOK, Header: chatRateHeaders(), Body: "not-json"},
			call: func(client *helix.Client) error {
				_, err := client.Chat.GetGlobalEmotes(context.Background(), helix.GetGlobalEmotesRequest{})
				return err
			},
			wantStatus: http.StatusOK,
			wantType: func(err error) bool {
				var typed *helix.ProtocolError
				return errors.As(err, &typed)
			},
		},
		{
			name:     "api failure",
			response: testkit.RoundTripResponse{StatusCode: http.StatusBadRequest, Header: chatRateHeaders(), Body: `{"error":"Bad Request","status":400,"message":"invalid"}`},
			call: func(client *helix.Client) error {
				_, err := client.Chat.GetChannelEmotes(context.Background(), helix.GetChannelEmotesRequest{BroadcasterID: "1234"})
				return err
			},
			wantStatus: http.StatusBadRequest,
			wantType: func(err error) bool {
				var typed *helix.APIError
				return errors.As(err, &typed)
			},
		},
	}

	// When the typed failure is returned.
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper(testCase.response)
			client, err := helix.New(
				helix.WithBaseURL("https://api.twitch.test/helix"),
				helix.WithHTTPClient(&http.Client{Transport: transport}),
				helix.WithStaticToken(chatCredential(helix.TokenClassApp, "")),
			)
			if err != nil {
				t.Fatal(err)
			}
			callErr := testCase.call(client)
			if !testCase.wantType(callErr) {
				t.Fatalf("error = %T %v, want typed Chat error", callErr, callErr)
			}
			var metaErr interface{ Meta() helix.ResponseMeta }
			if !errors.As(callErr, &metaErr) || metaErr.Meta().StatusCode() != testCase.wantStatus || !metaErr.Meta().RateLimit().Valid() {
				t.Fatalf("error metadata = %#v, want status and rate metadata", callErr)
			}
		})
	}
}
