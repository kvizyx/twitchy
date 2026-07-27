package oauth

import (
	"context"
	"net/url"
)

func (client *Client) Revoke(ctx context.Context, request RevokeRequest) error {
	if err := client.validClient(); err != nil {
		return err
	}
	if request.ClientID == "" || request.AccessToken == "" {
		return ErrInvalidOption
	}
	values := url.Values{
		"client_id": {request.ClientID},
		"token":     {request.AccessToken},
	}
	_, err := client.post(ctx, "revoke", "revoke", values, request.ClientID, request.AccessToken)
	return err
}
