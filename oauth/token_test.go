package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func TestToken_codeCredentialsAndExchangeUseExactForm(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests.Add(1)
		if request.Method != http.MethodPost || request.URL.Path != "/oauth2/token" || request.URL.RawQuery != "" {
			t.Errorf("request = %s %s?%s", request.Method, request.URL.Path, request.URL.RawQuery)
		}
		if request.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("content type = %q", request.Header.Get("Content-Type"))
		}
		if err := request.ParseForm(); err != nil {
			t.Errorf("ParseForm() error = %v", err)
		}
		if request.Form.Get("client_secret") == "" {
			t.Errorf("client_secret missing from form")
		}
		_, _ = io.WriteString(writer, `{"access_token":"access","expires_in":3600,"scope":["user:read:email"],"token_type":"bearer","refresh_token":"refresh"}`)
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	appToken, err := client.ClientCredentials(context.Background(), ClientCredentialsRequest{ClientID: "client", ClientSecret: "secret"})
	if err != nil {
		t.Fatalf("ClientCredentials() error = %v", err)
	}
	codeToken, err := client.ExchangeCode(context.Background(), ExchangeCodeRequest{ClientID: "client", ClientSecret: "secret", Code: "code", RedirectURI: "https://app.example/callback"})
	if err != nil {
		t.Fatalf("ExchangeCode() error = %v", err)
	}
	if requests.Load() != 2 || appToken.AccessToken != "access" || appToken.ExpiresIn != time.Hour || codeToken.RefreshToken != "refresh" {
		t.Fatalf("requests/tokens = %d/%+v/%+v", requests.Load(), appToken, codeToken)
	}
}

func TestToken_malformedSuccessIsProtocolError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) { _, _ = io.WriteString(writer, `{}`) }))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	_, err = client.ClientCredentials(context.Background(), ClientCredentialsRequest{ClientID: "client", ClientSecret: "secret"})
	var protocolErr *ProtocolError
	if !errors.As(err, &protocolErr) {
		t.Fatalf("error = %v, want ProtocolError", err)
	}
}

func TestDevice_pollPendingThenSuccessAndPublicDeviceSendsNoSecret(t *testing.T) {
	var polls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if err := request.ParseForm(); err != nil {
			t.Errorf("ParseForm() error = %v", err)
		}
		switch request.URL.Path {
		case "/oauth2/device":
			if request.Form.Get("client_secret") != "" || request.Form.Get("scopes") == "" {
				t.Errorf("device form = %v", request.Form)
			}
			_, _ = io.WriteString(writer, `{"device_code":"device","expires_in":600,"interval":1,"user_code":"ABCD","verification_uri":"https://twitch.tv/activate"}`)
		case "/oauth2/token":
			if polls.Add(1) == 1 {
				writer.WriteHeader(http.StatusBadRequest)
				_, _ = io.WriteString(writer, `{"error":"authorization_pending","error_description":"wait"}`)
				return
			}
			_, _ = io.WriteString(writer, `{"access_token":"device-access","expires_in":600,"scope":[],"token_type":"bearer"}`)
		}
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	client.sleeper = func(context.Context, time.Duration) error { return nil }
	device, err := client.StartDeviceAuthorization(context.Background(), DeviceAuthorizationRequest{ClientID: "public", Scopes: []helix.AuthorizationScope{helix.ScopeUserReadEmail}})
	if err != nil {
		t.Fatalf("StartDeviceAuthorization() error = %v", err)
	}
	token, err := client.PollDeviceToken(context.Background(), PollDeviceTokenRequest{ClientID: "public", DeviceCode: device.DeviceCode, ExpiresAt: time.Now().Add(time.Hour)})
	if err != nil || token.AccessToken != "device-access" || polls.Load() != 2 {
		t.Fatalf("token/error/polls = %+v/%v/%d", token, err, polls.Load())
	}
}

func TestDevice_unknownErrorAndExpiryAreTerminal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(writer, `{"error":"access_denied"}`)
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	_, err = client.PollDeviceToken(context.Background(), PollDeviceTokenRequest{ClientID: "client", DeviceCode: "device", ExpiresAt: time.Now().Add(time.Hour)})
	var deviceErr *DeviceAuthorizationError
	if !errors.As(err, &deviceErr) || deviceErr.Code() != "access_denied" {
		t.Fatalf("device error = %v", err)
	}
	_, err = client.PollDeviceToken(context.Background(), PollDeviceTokenRequest{ClientID: "client", DeviceCode: "device", ExpiresAt: time.Now().Add(-time.Second)})
	if !errors.Is(err, ErrDeviceExpired) {
		t.Fatalf("expiry error = %v", err)
	}
}

func TestDevice_responseShapeDoesNotUseArbitraryJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = json.Marshal(request.URL.Path)
		_, _ = io.WriteString(writer, `{}`)
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	_, err = client.StartDeviceAuthorization(context.Background(), DeviceAuthorizationRequest{ClientID: "client"})
	var protocolErr *ProtocolError
	if !errors.As(err, &protocolErr) || strings.Contains(err.Error(), "client") {
		t.Fatalf("device protocol error = %v", err)
	}
}
