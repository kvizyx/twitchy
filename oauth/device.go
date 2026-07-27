package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

const defaultDeviceInterval = 5 * time.Second

func (client *Client) StartDeviceAuthorization(ctx context.Context, request DeviceAuthorizationRequest) (*DeviceAuthorization, error) {
	if err := client.validClient(); err != nil {
		return nil, err
	}
	if request.ClientID == "" {
		return nil, ErrInvalidOption
	}
	values := url.Values{"client_id": {request.ClientID}}
	if len(request.Scopes) > 0 {
		values.Set("scopes", joinScopes(request.Scopes))
	}
	body, err := client.post(ctx, "device_authorization", "device", values, request.ClientID)
	if err != nil {
		return nil, wrapDeviceError(err)
	}
	var wire struct {
		DeviceCode      string `json:"device_code"`
		ExpiresIn       *int64 `json:"expires_in"`
		Interval        *int64 `json:"interval"`
		UserCode        string `json:"user_code"`
		VerificationURI string `json:"verification_uri"`
	}
	if err := json.Unmarshal(body, &wire); err != nil {
		return nil, &ProtocolError{operation: "device_authorization", status: http.StatusOK, cause: err}
	}
	if wire.DeviceCode == "" || wire.ExpiresIn == nil || *wire.ExpiresIn <= 0 || wire.UserCode == "" || wire.VerificationURI == "" {
		return nil, &ProtocolError{operation: "device_authorization", status: http.StatusOK, cause: errors.New("missing device authorization fields")}
	}
	interval := defaultDeviceInterval
	if wire.Interval != nil && *wire.Interval > 0 {
		interval = time.Duration(*wire.Interval) * time.Second
	}
	return &DeviceAuthorization{
		DeviceCode:      wire.DeviceCode,
		ExpiresIn:       time.Duration(*wire.ExpiresIn) * time.Second,
		Interval:        interval,
		UserCode:        wire.UserCode,
		VerificationURI: wire.VerificationURI,
	}, nil
}

func (client *Client) PollDeviceToken(ctx context.Context, request PollDeviceTokenRequest) (*TokenPair, error) {
	if err := client.validClient(); err != nil {
		return nil, err
	}
	if request.ClientID == "" || request.DeviceCode == "" {
		return nil, ErrInvalidOption
	}
	interval := defaultDeviceInterval
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if !request.ExpiresAt.IsZero() && !client.clock.Now().Before(request.ExpiresAt) {
			return nil, ErrDeviceExpired
		}
		values := url.Values{
			"client_id":   {request.ClientID},
			"device_code": {request.DeviceCode},
			"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		}
		if len(request.Scopes) > 0 {
			values.Set("scopes", joinScopes(request.Scopes))
		}
		body, err := client.post(ctx, "device_token", "token", values, request.ClientID, request.DeviceCode)
		if err == nil {
			return decodeToken(body, "device_token")
		}
		var oauthErr *OAuthError
		if !errors.As(err, &oauthErr) || oauthErr.StatusCode() != http.StatusBadRequest {
			return nil, wrapDeviceError(err)
		}
		deviceErr := &DeviceAuthorizationError{oauthError: oauthErr}
		switch deviceErr.Code() {
		case "authorization_pending":
		case "slow_down":
			interval += 5 * time.Second
		default:
			return nil, deviceErr
		}
		if err := client.sleeper(ctx, interval); err != nil {
			return nil, err
		}
	}
}

func wrapDeviceError(err error) error {
	var oauthErr *OAuthError
	if errors.As(err, &oauthErr) {
		return &DeviceAuthorizationError{oauthError: oauthErr}
	}
	return err
}
