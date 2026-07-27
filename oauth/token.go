package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

const maxOAuthBody = 1 << 20

type tokenWire struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    *int64   `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

func (client *Client) ClientCredentials(ctx context.Context, request ClientCredentialsRequest) (*TokenPair, error) {
	if err := client.validClient(); err != nil {
		return nil, err
	}
	if request.ClientID == "" || request.ClientSecret == "" {
		return nil, ErrInvalidOption
	}
	values := url.Values{
		"client_id":     {request.ClientID},
		"client_secret": {request.ClientSecret},
		"grant_type":    {"client_credentials"},
	}
	body, err := client.post(ctx, "client_credentials", "token", values, request.ClientID, request.ClientSecret)
	if err != nil {
		return nil, err
	}
	return decodeToken(body, "client_credentials")
}

func (client *Client) ExchangeCode(ctx context.Context, request ExchangeCodeRequest) (*TokenPair, error) {
	if err := client.validClient(); err != nil {
		return nil, err
	}
	if request.ClientID == "" || request.ClientSecret == "" || request.Code == "" || request.RedirectURI == "" {
		return nil, ErrInvalidOption
	}
	values := url.Values{
		"client_id":     {request.ClientID},
		"client_secret": {request.ClientSecret},
		"code":          {request.Code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {request.RedirectURI},
	}
	body, err := client.post(ctx, "exchange_code", "token", values, request.ClientID, request.ClientSecret, request.Code)
	if err != nil {
		return nil, err
	}
	return decodeToken(body, "exchange_code")
}

func (client *Client) post(ctx context.Context, operation, endpoint string, values url.Values, secrets ...string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.endpoint(endpoint), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, fmt.Errorf("oauth: create %s request: %w", operation, err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.httpClient.Do(request)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, &TransportError{operation: operation, cause: err}
	}
	defer response.Body.Close()
	body, err := readOAuthBody(response.Body)
	if err != nil {
		return nil, &ProtocolError{operation: operation, status: response.StatusCode, cause: err}
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		oauthErr := decodeOAuthError(operation, response.StatusCode, body, secrets...)
		return nil, oauthErr
	}
	return body, nil
}

func decodeToken(body []byte, operation string) (*TokenPair, error) {
	var wire tokenWire
	if err := json.Unmarshal(body, &wire); err != nil {
		return nil, &ProtocolError{operation: operation, status: http.StatusOK, cause: err}
	}
	if wire.AccessToken == "" || wire.TokenType == "" || wire.ExpiresIn == nil || *wire.ExpiresIn <= 0 {
		return nil, &ProtocolError{operation: operation, status: http.StatusOK, cause: errors.New("missing token fields")}
	}
	return &TokenPair{
		AccessToken:  wire.AccessToken,
		RefreshToken: wire.RefreshToken,
		ExpiresIn:    time.Duration(*wire.ExpiresIn) * time.Second,
		Scopes:       parseWireScopes(wire.Scope),
		TokenType:    wire.TokenType,
	}, nil
}

func readOAuthBody(reader io.Reader) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(reader, maxOAuthBody+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxOAuthBody {
		return nil, errors.New("oauth response body exceeds limit")
	}
	return body, nil
}

func decodeOAuthError(operation string, status int, body []byte, secrets ...string) *OAuthError {
	var wire struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
		Message          string `json:"message"`
	}
	_ = json.Unmarshal(body, &wire)
	code := wire.Error
	description := wire.ErrorDescription
	if code == "" && wire.Message != "" {
		code = wire.Message
	}
	if description == "" && wire.Error != "" && wire.Message != "" {
		description = wire.Message
	}
	return &OAuthError{
		operation:   operation,
		status:      status,
		code:        sanitizeText(code, secrets...),
		description: sanitizeText(description, secrets...),
		retryable:   status == http.StatusTooManyRequests || status >= 500,
	}
}

func parseWireScopes(scopes []string) []helix.AuthorizationScope {
	if len(scopes) == 0 {
		return nil
	}
	result := make([]helix.AuthorizationScope, len(scopes))
	for index, scope := range scopes {
		result[index] = helix.AuthorizationScope(scope)
	}
	return result
}
