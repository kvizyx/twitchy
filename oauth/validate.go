package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (client *Client) Validate(ctx context.Context, request ValidateRequest) (*Validation, error) {
	if err := client.validClient(); err != nil {
		return nil, err
	}
	if request.AccessToken == "" {
		return nil, ErrInvalidOption
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, client.endpoint("validate"), nil)
	if err != nil {
		return nil, fmt.Errorf("oauth: create validate request: %w", err)
	}
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Authorization", "OAuth "+request.AccessToken)
	response, err := client.httpClient.Do(httpRequest)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, &TransportError{operation: "validate", cause: err}
	}
	defer response.Body.Close()
	body, err := readOAuthBody(response.Body)
	if err != nil {
		return nil, &ProtocolError{operation: "validate", status: response.StatusCode, cause: err}
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, decodeOAuthError("validate", response.StatusCode, body, request.AccessToken)
	}
	var wire struct {
		ClientID  string   `json:"client_id"`
		Login     string   `json:"login"`
		Scopes    []string `json:"scopes"`
		UserID    string   `json:"user_id"`
		ExpiresIn *int64   `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &wire); err != nil {
		return nil, &ProtocolError{operation: "validate", status: http.StatusOK, cause: err}
	}
	if wire.ClientID == "" || wire.Login == "" || wire.UserID == "" || wire.ExpiresIn == nil || *wire.ExpiresIn <= 0 {
		return nil, &ProtocolError{operation: "validate", status: http.StatusOK, cause: errors.New("missing validation fields")}
	}
	return &Validation{
		ClientID:  wire.ClientID,
		Login:     wire.Login,
		Scopes:    parseWireScopes(wire.Scopes),
		UserID:    wire.UserID,
		ExpiresIn: time.Duration(*wire.ExpiresIn) * time.Second,
	}, nil
}
