package oauth

import (
	"context"
	"net/url"
)

func (client *Client) Refresh(ctx context.Context, request RefreshRequest) (*TokenPair, error) {
	if err := client.validClient(); err != nil {
		return nil, err
	}
	if request.ClientID == "" || request.RefreshToken == "" {
		return nil, ErrInvalidOption
	}
	values := url.Values{
		"client_id":     {request.ClientID},
		"grant_type":    {"refresh_token"},
		"refresh_token": {request.RefreshToken},
	}
	if request.ClientSecret != "" {
		values.Set("client_secret", request.ClientSecret)
	}
	body, err := client.post(ctx, "refresh", "token", values, request.ClientID, request.ClientSecret, request.RefreshToken)
	if err != nil {
		return nil, err
	}
	return decodeToken(body, "refresh")
}
