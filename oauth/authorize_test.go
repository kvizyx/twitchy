package oauth

import (
	"bytes"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/kvizyx/twitchy/helix"
)

func TestAuthorizeCode_buildsRandomStateAndParsesCallback(t *testing.T) {
	client, err := New(WithRandom(bytes.NewReader(bytes.Repeat([]byte{7}, 32))))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	authorization, err := client.AuthorizationURL(AuthorizationRequest{
		ClientID: "client-id", RedirectURI: "HTTPS://APP.EXAMPLE/callback?tenant=blue",
		ResponseType: ResponseTypeCode, Scopes: []helix.AuthorizationScope{helix.ScopeUserReadEmail},
	})
	if err != nil {
		t.Fatalf("AuthorizationURL() error = %v", err)
	}
	parsed, err := url.Parse(authorization.URL)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	state := parsed.Query().Get("state")
	decoded, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil || len(decoded) != 32 || authorization.State != state {
		t.Fatalf("state = %q, decoded length = %d, error = %v", state, len(decoded), err)
	}
	callback := "https://app.example/callback?tenant=blue&code=one-time-code&state=" + url.QueryEscape(state)
	result, err := client.ParseAuthorizationCallback(callback, authorization)
	if err != nil || result.Code != "one-time-code" || result.State != state {
		t.Fatalf("callback result/error = %+v/%v", result, err)
	}
}

func TestAuthorizeImplicit_rejectsStateAndRedirectMismatch(t *testing.T) {
	client, err := New(WithRandom(bytes.NewReader(bytes.Repeat([]byte{9}, 32))))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	authorization, err := client.AuthorizationURL(AuthorizationRequest{
		ClientID: "client-id", RedirectURI: "https://app.example/callback?tenant=blue", ResponseType: ResponseTypeToken,
	})
	if err != nil {
		t.Fatalf("AuthorizationURL() error = %v", err)
	}
	if subtle.ConstantTimeCompare([]byte(authorization.State), []byte("wrong")) == 1 {
		t.Fatal("test state unexpectedly matched")
	}
	_, err = client.ParseAuthorizationCallback("https://app.example/callback?tenant=blue&state=wrong", authorization)
	if !errors.Is(err, ErrStateMismatch) {
		t.Fatalf("wrong state error = %v", err)
	}
	_, err = client.ParseAuthorizationCallback("https://app.example/callback?tenant=red&state="+url.QueryEscape(authorization.State), authorization)
	if !errors.Is(err, ErrRedirectMismatch) {
		t.Fatalf("wrong redirect error = %v", err)
	}
	callback := "https://app.example/callback?tenant=blue#access_token=access&token_type=Bearer&scope=user%3Aread%3Aemail&state=" + url.QueryEscape(authorization.State)
	result, err := client.ParseAuthorizationCallback(callback, authorization)
	if err != nil || result.AccessToken != "access" || result.TokenType != "Bearer" || len(result.Scopes) != 1 {
		t.Fatalf("implicit result/error = %+v/%v", result, err)
	}
}

func TestAuthorizeCallback_rejectsUnknownFixedQueryAndUserInfo(t *testing.T) {
	client, err := New(WithRandom(bytes.NewReader(bytes.Repeat([]byte{1}, 32))))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	authorization, err := client.AuthorizationURL(AuthorizationRequest{ClientID: "client", RedirectURI: "https://app.example/callback?tenant=blue", ResponseType: ResponseTypeCode})
	if err != nil {
		t.Fatalf("AuthorizationURL() error = %v", err)
	}
	callback := "https://user@app.example/callback?tenant=blue&extra=1&state=" + url.QueryEscape(authorization.State)
	_, err = client.ParseAuthorizationCallback(callback, authorization)
	if !errors.Is(err, ErrRedirectMismatch) || strings.Contains(err.Error(), authorization.State) {
		t.Fatalf("hostile callback error = %v", err)
	}
}
