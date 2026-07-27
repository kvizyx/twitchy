package oauth

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"sort"
	"strings"

	"github.com/kvizyx/twitchy/helix"
)

var responseQueryKeys = map[string]struct{}{
	"access_token": {}, "code": {}, "error": {}, "error_description": {},
	"scope": {}, "state": {}, "token_type": {},
}

func (client *Client) AuthorizationURL(request AuthorizationRequest) (Authorization, error) {
	if err := client.validClient(); err != nil {
		return Authorization{}, err
	}
	if request.ClientID == "" || request.RedirectURI == "" {
		return Authorization{}, ErrInvalidOption
	}
	if request.ResponseType != ResponseTypeCode && request.ResponseType != ResponseTypeToken {
		return Authorization{}, ErrInvalidOption
	}
	redirect, err := url.Parse(request.RedirectURI)
	if err != nil || redirect.Scheme == "" || redirect.Host == "" || redirect.User != nil || redirect.Fragment != "" {
		return Authorization{}, ErrInvalidOption
	}
	stateBytes := make([]byte, 32)
	if _, err := io.ReadFull(client.random, stateBytes); err != nil {
		return Authorization{}, fmt.Errorf("oauth: generate state: %w", err)
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes)
	query := url.Values{}
	query.Set("client_id", request.ClientID)
	query.Set("redirect_uri", request.RedirectURI)
	query.Set("response_type", string(request.ResponseType))
	query.Set("state", state)
	if len(request.Scopes) > 0 {
		query.Set("scope", joinScopes(request.Scopes))
	}
	if request.ForceVerify {
		query.Set("force_verify", "true")
	}
	endpoint := client.endpoint("authorize")
	return Authorization{
		URL:          endpoint + "?" + query.Encode(),
		State:        state,
		RedirectURI:  request.RedirectURI,
		ResponseType: request.ResponseType,
	}, nil
}

func (client *Client) ParseAuthorizationCallback(raw string, authorization Authorization) (AuthorizationResult, error) {
	if err := client.validClient(); err != nil {
		return AuthorizationResult{}, err
	}
	callback, err := url.Parse(raw)
	if err != nil {
		return AuthorizationResult{}, fmt.Errorf("oauth: parse authorization callback: %w", err)
	}
	registered, err := url.Parse(authorization.RedirectURI)
	if err != nil {
		return AuthorizationResult{}, fmt.Errorf("oauth: parse registered redirect: %w", err)
	}
	response := callback.Query()
	if authorization.ResponseType == ResponseTypeToken && callback.Fragment != "" {
		fragment, fragmentErr := url.ParseQuery(callback.Fragment)
		if fragmentErr != nil {
			return AuthorizationResult{}, fmt.Errorf("oauth: parse authorization fragment: %w", fragmentErr)
		}
		response = fragment
	}
	if authorization.State == "" || subtle.ConstantTimeCompare([]byte(response.Get("state")), []byte(authorization.State)) != 1 {
		return AuthorizationResult{}, ErrStateMismatch
	}
	if !sameRedirect(callback, registered) || !sameFixedQuery(callback.Query(), registered.Query()) {
		return AuthorizationResult{}, ErrRedirectMismatch
	}
	if code := response.Get("error"); code != "" {
		return AuthorizationResult{}, &OAuthError{
			operation:   "authorization callback",
			code:        sanitizeText(code),
			description: sanitizeText(response.Get("error_description")),
		}
	}
	result := AuthorizationResult{State: response.Get("state"), TokenType: response.Get("token_type")}
	result.Scopes = parseScopes(response.Get("scope"))
	switch authorization.ResponseType {
	case ResponseTypeCode:
		result.Code = response.Get("code")
		if result.Code == "" {
			return AuthorizationResult{}, &ProtocolError{operation: "authorization callback", cause: errors.New("missing code")}
		}
	case ResponseTypeToken:
		result.AccessToken = response.Get("access_token")
		if result.AccessToken == "" || result.TokenType == "" {
			return AuthorizationResult{}, &ProtocolError{operation: "authorization callback", cause: errors.New("missing implicit token fields")}
		}
	default:
		return AuthorizationResult{}, ErrInvalidOption
	}
	return result, nil
}

func sameRedirect(actual, registered *url.URL) bool {
	if actual.User != nil || registered.User != nil {
		return false
	}
	return normalizedEndpoint(actual) == normalizedEndpoint(registered)
}

func normalizedEndpoint(value *url.URL) string {
	scheme := strings.ToLower(value.Scheme)
	host := strings.ToLower(value.Hostname())
	port := value.Port()
	if (scheme == "http" && port == "80") || (scheme == "https" && port == "443") {
		port = ""
	}
	if port != "" {
		host += ":" + port
	}
	redirectPath := path.Clean(value.Path)
	if redirectPath == "." {
		redirectPath = "/"
	}
	return scheme + "://" + host + redirectPath
}

func sameFixedQuery(actual, registered url.Values) bool {
	actual = fixedQuery(actual)
	registered = fixedQuery(registered)
	if len(actual) != len(registered) {
		return false
	}
	for key, expected := range registered {
		got, ok := actual[key]
		if !ok || !sameValues(got, expected) {
			return false
		}
	}
	return true
}

func fixedQuery(values url.Values) url.Values {
	result := make(url.Values)
	for key, entries := range values {
		if _, responseKey := responseQueryKeys[key]; responseKey {
			continue
		}
		result[key] = append([]string(nil), entries...)
	}
	return result
}

func sameValues(left, right []string) bool {
	left = append([]string(nil), left...)
	right = append([]string(nil), right...)
	sort.Strings(left)
	sort.Strings(right)
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func joinScopes(scopes []helix.AuthorizationScope) string {
	values := make([]string, len(scopes))
	for index, scope := range scopes {
		values[index] = string(scope)
	}
	return strings.Join(values, " ")
}

func parseScopes(raw string) []helix.AuthorizationScope {
	if raw == "" {
		return nil
	}
	fields := strings.Fields(raw)
	scopes := make([]helix.AuthorizationScope, len(fields))
	for index, field := range fields {
		scopes[index] = helix.AuthorizationScope(field)
	}
	return scopes
}
