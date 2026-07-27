package helix_test

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

const task26RateReset = "4102444800"

var _ interface {
	GetAuthorizationByUser(context.Context, helix.GetAuthorizationByUserRequest) (*helix.Response[helix.GetAuthorizationByUserData], error)
} = (*helix.ExperimentalUsersService)(nil)

func task26Ptr[T any](value T) *T {
	return &value
}

func task26Body(t *testing.T, name string) string {
	t.Helper()
	body, err := testkit.LoadText("testdata/task26", name)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func task26Response(status int, body string) testkit.RoundTripResponse {
	return testkit.RoundTripResponse{
		StatusCode: status,
		Header: http.Header{
			"Ratelimit-Limit":     {"8000"},
			"Ratelimit-Remaining": {"7999"},
			"Ratelimit-Reset":     {task26RateReset},
		},
		Body: body,
	}
}

func task26Client(t *testing.T, transport *testkit.RecordingRoundTripper, credential helix.Credential) *helix.Client {
	t.Helper()
	client, err := helix.New(
		helix.WithBaseURL("https://api.twitch.test/helix"),
		helix.WithHTTPClient(&http.Client{Transport: transport}),
		helix.WithStaticToken(credential),
	)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func task26Success(body string) testkit.ContractResponse {
	return testkit.ContractResponse{
		Status: http.StatusOK,
		Headers: http.Header{
			"Ratelimit-Limit":     {"8000"},
			"Ratelimit-Remaining": {"7999"},
			"Ratelimit-Reset":     {task26RateReset},
		},
		Body:    body,
		Success: true,
	}
}

func task26Contract(t *testing.T, anchor string, query map[string][]string, body string, response testkit.ContractResponse, transport *testkit.RecordingRoundTripper, meta helix.ResponseMeta) {
	t.Helper()
	fixture := testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query:   query,
			Headers: http.Header{"Authorization": {"Bearer user-token"}, "Client-Id": {"client-id"}},
			Body:    body,
		},
		Response: response,
		Want: testkit.ContractExpectation{
			Attempts:       1,
			RateLimitValid: true,
		},
	}
	fixture.Want.RateLimit.Limit = 8000
	fixture.Want.RateLimit.Remaining = 7999
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, anchor), fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return meta, nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestUsersGetUsers(t *testing.T) {
	// Given users selected by repeated IDs and login names.
	body := task26Body(t, "users.json")
	transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusOK, body))
	client := task26Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "viewer"})

	// When the users endpoint is requested.
	result, err := client.Users.GetUsers(context.Background(), helix.GetUsersRequest{IDs: []string{"1", "2"}, Logins: []string{"alpha", "beta"}})

	// Then repeated query values and optional email data are preserved.
	if err != nil {
		t.Fatal(err)
	}
	task26Contract(t, "get-users", urlValues("id", "1", "id", "2", "login", "alpha", "login", "beta"), "", task26Success(body), transport, result.Meta)
	if len(result.Data) != 2 || result.Data[0].Email != "alpha@example.test" || result.Data[1].CreatedAt.IsZero() {
		t.Fatalf("users = %#v", result.Data)
	}
}

func TestUsersUpdateUser(t *testing.T) {
	// Given a user token with the user:edit scope.
	body := task26Body(t, "users.json")
	transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusOK, body))
	client := task26Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "viewer", Scopes: []helix.AuthorizationScope{helix.ScopeUserEdit}})

	// When the description is updated.
	result, err := client.Users.UpdateUser(context.Background(), helix.UpdateUserRequest{Description: task26Ptr("new description")})

	// Then the optional description is encoded as a query field.
	if err != nil {
		t.Fatal(err)
	}
	task26Contract(t, "update-user", urlValues("description", "new description"), "", task26Success(body), transport, result.Meta)
}

func TestUsersAuthorizationByUserIsExperimentalOnly(t *testing.T) {
	// Given an app token and the compile-boundary method on Experimental.Users.
	body := task26Body(t, "authorization_by_user.json")
	transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusOK, body))
	client := task26Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})

	// When authorization is requested for repeated user IDs.
	result, err := client.Experimental.Users.GetAuthorizationByUser(context.Background(), helix.GetAuthorizationByUserRequest{UserIDs: []string{"1", "2"}})

	// Then the NEW operation is reachable only through Experimental.Users.
	if err != nil {
		t.Fatal(err)
	}
	task26Contract(t, "get-authorization-by-user", urlValues("user_id", "1", "user_id", "2"), "", task26Success(body), transport, result.Meta)
	if _, exists := reflect.TypeOf(client.Users).MethodByName("GetAuthorizationByUser"); exists {
		t.Fatal("GetAuthorizationByUser leaked onto stable UsersService")
	}
}

func TestUsersAuthorizationByUserRequiresAppToken(t *testing.T) {
	// Given a user token for an app-only operation.
	transport := testkit.NewRecordingRoundTripper()
	client := task26Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "viewer"})

	// When authorization is requested.
	_, err := client.Experimental.Users.GetAuthorizationByUser(context.Background(), helix.GetAuthorizationByUserRequest{UserIDs: []string{"1"}})

	// Then local auth rejects it before network I/O.
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) || len(transport.Requests()) != 0 {
		t.Fatalf("error = %T %v, requests = %d", err, err, len(transport.Requests()))
	}
}

func TestUsersGetUserBlockList(t *testing.T) {
	// Given a scoped user token for the requested broadcaster.
	body := task26Body(t, "user_blocks.json")
	transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusOK, body))
	client := task26Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "broadcaster", Scopes: []helix.AuthorizationScope{helix.ScopeUserReadBlockedUsers}})

	// When the block list is requested with pagination fields.
	result, err := client.Users.GetUserBlockList(context.Background(), helix.GetUserBlockListRequest{BroadcasterID: "broadcaster", First: task26Ptr(2), After: task26Ptr("cursor-a")})

	// Then the user subject, cursor, and blocked-user fields are preserved.
	if err != nil {
		t.Fatal(err)
	}
	task26Contract(t, "get-user-block-list", urlValues("broadcaster_id", "broadcaster", "first", "2", "after", "cursor-a"), "", task26Success(body), transport, result.Meta)
	if len(result.Data) != 1 || result.Data[0].DisplayName != "Blocked" || result.Pagination.Cursor() != "cursor-b" {
		t.Fatalf("blocked users = %#v", result.Data)
	}
}

func TestUsersBlockAndUnblockUserPreserveEnums(t *testing.T) {
	// Given a user token with blocked-user management scope.
	for _, testCase := range []struct {
		name  string
		call  func(*helix.Client) error
		query map[string][]string
	}{
		{
			name: "block",
			call: func(client *helix.Client) error {
				_, err := client.Users.BlockUser(context.Background(), helix.BlockUserRequest{TargetUserID: "target", SourceContext: task26Ptr(helix.BlockUserSourceContextWhisper), Reason: task26Ptr(helix.BlockUserReasonSpam)})
				return err
			},
			query: urlValues("target_user_id", "target", "source_context", "whisper", "reason", "spam"),
		},
		{
			name: "unblock",
			call: func(client *helix.Client) error {
				_, err := client.Users.UnblockUser(context.Background(), helix.UnblockUserRequest{TargetUserID: "target"})
				return err
			},
			query: urlValues("target_user_id", "target"),
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusNoContent, ""))
			client := task26Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "broadcaster", Scopes: []helix.AuthorizationScope{helix.ScopeUserManageBlockedUsers}})

			// When the mutation is requested.
			err := testCase.call(client)

			// Then the exact query is sent without a body.
			if err != nil {
				t.Fatal(err)
			}
			if len(transport.Requests()) != 1 || transport.Requests()[0].Body != nil {
				t.Fatalf("requests = %#v", transport.Requests())
			}
			request := transport.Requests()[0]
			if request.Path != "/helix/users/blocks?"+urlValuesToQuery(testCase.query) {
				t.Fatalf("path = %q", request.Path)
			}
		})
	}
}

func TestUsersExtensions(t *testing.T) {
	// Given installed extensions and active panel, overlay, and component maps.
	installedBody := task26Body(t, "user_extensions.json")
	activeBody := task26Body(t, "active_extensions.json")
	transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusOK, installedBody), task26Response(http.StatusOK, activeBody))
	client := task26Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "broadcaster", Scopes: []helix.AuthorizationScope{helix.ScopeUserReadBroadcast}})

	// When installed extensions are requested.
	installed, err := client.Users.GetUserExtensions(context.Background(), helix.GetUserExtensionsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(installed.Data) != 1 || installed.Data[0].Type[0] != helix.UserExtensionTypePanel {
		t.Fatalf("installed extensions = %#v", installed.Data)
	}

	// When active extensions are requested for the token subject.
	active, err := client.Users.GetUserActiveExtensions(context.Background(), helix.GetUserActiveExtensionsRequest{UserID: task26Ptr("broadcaster")})
	if err != nil {
		t.Fatal(err)
	}
	if active.Data.Panel["1"].ID != "panel-id" || active.Data.Overlay["1"].Version != "2.0" || active.Data.Component["1"].X != 3 {
		t.Fatalf("active extensions = %#v", active.Data)
	}
}

func urlValuesToQuery(values map[string][]string) string {
	query := make([]string, 0, len(values))
	for key, items := range values {
		for _, item := range items {
			query = append(query, key+"="+item)
		}
	}
	sort.Strings(query)
	return strings.Join(query, "&")
}

func TestUsersUpdateUserExtensions(t *testing.T) {
	// Given a user token with broadcast-edit scope and nested extension maps.
	body := task26Body(t, "active_extensions.json")
	requestBody := `{"data":{"component":{"1":{"active":false}},"panel":{"1":{"active":true,"id":"panel-id","version":"1.0"}}}}`
	transport := testkit.NewRecordingRoundTripper(task26Response(http.StatusOK, body))
	client := task26Client(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser, UserID: "broadcaster", Scopes: []helix.AuthorizationScope{helix.ScopeUserEditBroadcast}})

	// When installed extension state is updated.
	result, err := client.Users.UpdateUserExtensions(context.Background(), helix.UpdateUserExtensionsRequest{Data: map[string]map[string]helix.UserExtensionUpdate{
		"panel":     {"1": {Active: true, ID: task26Ptr("panel-id"), Version: task26Ptr("1.0")}},
		"component": {"1": {Active: false}},
	}})

	// Then the exact nested JSON body is sent.
	if err != nil {
		t.Fatal(err)
	}
	task26Contract(t, "update-user-extensions", nil, requestBody, task26Success(body), transport, result.Meta)
}
