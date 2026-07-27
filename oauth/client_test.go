package oauth

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRefreshValidateRevoke_andTypedHTTPFailures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/oauth2/token":
			if err := request.ParseForm(); err != nil || request.Form.Get("grant_type") != "refresh_token" {
				t.Errorf("refresh form/error = %v/%v", request.Form, err)
			}
			_, _ = io.WriteString(writer, `{"access_token":"new-access","expires_in":60,"scope":[],"token_type":"bearer"}`)
		case "/oauth2/validate":
			if request.Method != http.MethodGet || request.Header.Get("Authorization") != "OAuth access" {
				t.Errorf("validate request = %s %q", request.Method, request.Header.Get("Authorization"))
			}
			_, _ = io.WriteString(writer, `{"client_id":"client","login":"streamer","scopes":["user:read:email"],"user_id":"42","expires_in":60}`)
		case "/oauth2/revoke":
			writer.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	refreshed, err := client.Refresh(context.Background(), RefreshRequest{ClientID: "client", ClientSecret: "secret", RefreshToken: "refresh"})
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}
	validation, err := client.Validate(context.Background(), ValidateRequest{AccessToken: "access"})
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if err := client.Revoke(context.Background(), RevokeRequest{ClientID: "client", AccessToken: "access"}); err != nil {
		t.Fatalf("Revoke() error = %v", err)
	}
	if refreshed.AccessToken != "new-access" || validation.UserID != "42" || validation.ExpiresIn != time.Minute {
		t.Fatalf("refresh/validation = %+v/%+v", refreshed, validation)
	}
}

func TestRevoke_returnsTypedBadRequestAndRedactsToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(writer, `{"status":400,"message":"invalid token","error":"invalid_token"}`)
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	err = client.Revoke(context.Background(), RevokeRequest{ClientID: "client", AccessToken: "secret-token"})
	var oauthErr *OAuthError
	if !errors.As(err, &oauthErr) || oauthErr.StatusCode() != http.StatusBadRequest || oauthErr.Code() != "invalid_token" || strings.Contains(err.Error(), "secret-token") {
		t.Fatalf("error = %v", err)
	}
}

func TestRevoke_returnsTypedNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(writer, `{"status":404,"message":"client does not exist"}`)
	}))
	defer server.Close()
	client, err := New(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	err = client.Revoke(context.Background(), RevokeRequest{ClientID: "client", AccessToken: "access"})
	var oauthErr *OAuthError
	if !errors.As(err, &oauthErr) || oauthErr.StatusCode() != http.StatusNotFound {
		t.Fatalf("error = %v, want typed 404", err)
	}
}

func TestNew_rejectsRedirectsWithoutMutatingCallerClient(t *testing.T) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar.New() error = %v", err)
	}
	transport := &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}}
	redirect := func(*http.Request, []*http.Request) error { return errors.New("caller redirect") }
	caller := &http.Client{Transport: transport, Jar: jar, CheckRedirect: redirect, Timeout: time.Second}
	before := *caller
	client, err := New(WithHTTPClient(caller))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	after := *caller
	if after.Transport != before.Transport || after.Jar != before.Jar || after.Timeout != before.Timeout || reflect.ValueOf(after.CheckRedirect).Pointer() != reflect.ValueOf(before.CheckRedirect).Pointer() {
		t.Fatal("New() mutated caller HTTP client")
	}
	if client.httpClient == caller || client.httpClient.Transport != caller.Transport || client.httpClient.Jar != caller.Jar {
		t.Fatal("New() did not shallow-clone caller HTTP client")
	}
}

func TestClient_redirectCannotCarryOAuthSecret(t *testing.T) {
	var targetHits int
	target := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) { targetHits++ }))
	defer target.Close()
	source := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, target.URL, http.StatusFound)
	}))
	defer source.Close()
	client, err := New(WithBaseURL(source.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	_, err = client.ClientCredentials(context.Background(), ClientCredentialsRequest{ClientID: "client", ClientSecret: "super-secret"})
	if targetHits != 0 || strings.Contains(err.Error(), "super-secret") {
		t.Fatalf("redirect/secrecy = %d/%v", targetHits, err)
	}
}

func TestNew_rejectsNonLoopbackHTTP(t *testing.T) {
	_, err := New(WithBaseURL("http://example.test"))
	if !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("non-loopback HTTP error = %v", err)
	}
}
