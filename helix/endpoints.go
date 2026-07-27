package helix

import (
	"context"
	"fmt"
	"strings"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

// Endpoint convention: define named request/data types from the descriptor,
// resolve the manifest anchor, validate its credential rules, build the exact
// wire request, execute through the shared transport, and decode into Response.
// Add one fixture-backed contract test for the manifest row before each batch.

type GetGamesRequest struct {
	IDs     []string `query:"id"`
	Names   []string `query:"name"`
	IGDBIDs []string `query:"igdb_id"`
}

type Game struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"box_art_url"`
	IGDBID    string `json:"igdb_id"`
}

type GetGamesData []Game

type GetCustomPowerUpRequest struct {
	BroadcasterID string   `query:"broadcaster_id"`
	IDs           []string `query:"id,omitempty"`
}

type CustomPowerUpImage struct {
	URL1x string `json:"url_1x"`
	URL2x string `json:"url_2x"`
	URL4x string `json:"url_4x"`
}

type CustomPowerUpLimitSetting struct {
	IsEnabled    bool  `json:"is_enabled"`
	MaxPerStream int64 `json:"max_per_stream"`
}

type CustomPowerUpUserLimitSetting struct {
	IsEnabled           bool  `json:"is_enabled"`
	MaxPerUserPerStream int64 `json:"max_per_user_per_stream"`
}

type CustomPowerUpCooldownSetting struct {
	IsEnabled             bool `json:"is_enabled"`
	GlobalCooldownSeconds int  `json:"global_cooldown_seconds"`
}

type CustomPowerUp struct {
	BroadcasterID                    string                        `json:"broadcaster_id"`
	BroadcasterLogin                 string                        `json:"broadcaster_login"`
	BroadcasterName                  string                        `json:"broadcaster_name"`
	ID                               string                        `json:"id"`
	Title                            string                        `json:"title"`
	Prompt                           string                        `json:"prompt"`
	Bits                             int64                         `json:"bits"`
	Image                            *CustomPowerUpImage           `json:"image"`
	DefaultImage                     CustomPowerUpImage            `json:"default_image"`
	BackgroundColor                  string                        `json:"background_color"`
	IsEnabled                        bool                          `json:"is_enabled"`
	IsUserInputRequired              bool                          `json:"is_user_input_required"`
	MaxPerStreamSetting              CustomPowerUpLimitSetting     `json:"max_per_stream_setting"`
	MaxPerUserPerStreamSetting       CustomPowerUpUserLimitSetting `json:"max_per_user_per_stream_setting"`
	GlobalCooldownSetting            CustomPowerUpCooldownSetting  `json:"global_cooldown_setting"`
	IsPaused                         bool                          `json:"is_paused"`
	IsInStock                        bool                          `json:"is_in_stock"`
	RedemptionsRedeemedCurrentStream *int                          `json:"redemptions_redeemed_current_stream"`
	CooldownExpiresAt                *Timestamp                    `json:"cooldown_expires_at"`
}

type GetCustomPowerUpData []CustomPowerUp

func (s *GamesService) GetGames(ctx context.Context, req GetGamesRequest) (*Response[GetGamesData], error) {
	return executeEndpoint[GetGamesData](s.client, ctx, "get-games", req)
}

func (s *ExperimentalBitsService) GetCustomPowerUp(ctx context.Context, req GetCustomPowerUpRequest) (*Response[GetCustomPowerUpData], error) {
	return executeEndpoint[GetCustomPowerUpData](s.client, ctx, "get-custom-power-up", req)
}

func executeEndpoint[T any](client *Client, ctx context.Context, anchor string, query any) (*Response[T], error) {
	if err := client.validClient(); err != nil {
		return nil, err
	}
	operation, err := manifest.OperationByAnchor(anchor)
	if err != nil {
		return nil, err
	}
	credential, err := client.credential(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateCredentialForOperation(credential, operation, "", ""); err != nil {
		return nil, err
	}
	request, err := buildRequest(requestSpec{Context: ctx, Method: operation.Method, URL: client.endpointURL(operation.Path), Query: query})
	if err != nil {
		return nil, fmt.Errorf("build %s request: %w", operation.OperationID, err)
	}
	request.Header.Set("User-Agent", client.userAgent)
	response, meta, err := client.executor.execute(ctx, request, operation, credential)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	result, err := decodeResponse[T](response.StatusCode, response.Body, DecodeOptions{})
	if err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: []string{credential.AccessToken()}})
	}
	result.Meta = meta
	return result, nil
}

func (client *Client) credential(ctx context.Context) (CredentialSnapshot, error) {
	if client.executor.source == nil {
		if err := ctx.Err(); err != nil {
			return CredentialSnapshot{}, err
		}
		return CredentialSnapshot{}, nil
	}
	credential, err := client.executor.source.Token(ctx)
	if err != nil {
		return CredentialSnapshot{}, fmt.Errorf("read Helix credential: %w", err)
	}
	return credential, nil
}

func (client *Client) endpointURL(operationPath string) string {
	target := *client.baseURL
	basePath := strings.TrimSuffix(target.Path, "/")
	if strings.HasSuffix(basePath, "/helix") && strings.HasPrefix(operationPath, "/helix/") {
		target.Path = basePath + strings.TrimPrefix(operationPath, "/helix")
	} else {
		target.Path = basePath + operationPath
	}
	target.RawPath = ""
	return target.String()
}
