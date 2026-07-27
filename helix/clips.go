package helix

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type CreateClipRequest struct {
	BroadcasterID string   `query:"broadcaster_id"`
	Title         *string  `query:"title,omitempty"`
	Duration      *float64 `query:"duration,omitempty"`
}

type ClipCreation struct {
	ID      string `json:"id"`
	EditURL string `json:"edit_url"`
}

type CreateClipData []ClipCreation

type CreateClipFromVODRequest struct {
	EditorID      string   `query:"editor_id"`
	BroadcasterID string   `query:"broadcaster_id"`
	VODID         string   `query:"vod_id"`
	VODOffset     int      `query:"vod_offset"`
	Duration      *float64 `query:"duration,omitempty"`
	Title         string   `query:"title"`
}

type CreateClipFromVODData []ClipCreation

type GetClipsRequest struct {
	BroadcasterID string   `query:"broadcaster_id,omitempty"`
	GameID        string   `query:"game_id,omitempty"`
	IDs           []string `query:"id,omitempty"`
	StartedAt     *string  `query:"started_at,omitempty"`
	EndedAt       *string  `query:"ended_at,omitempty"`
	First         *int     `query:"first,omitempty"`
	Before        *string  `query:"before,omitempty"`
	After         *string  `query:"after,omitempty"`
	IsFeatured    *bool    `query:"is_featured,omitempty"`
}

type Clip struct {
	ID              string    `json:"id"`
	URL             string    `json:"url"`
	EmbedURL        string    `json:"embed_url"`
	BroadcasterID   string    `json:"broadcaster_id"`
	BroadcasterName string    `json:"broadcaster_name"`
	CreatorID       string    `json:"creator_id"`
	CreatorName     string    `json:"creator_name"`
	VideoID         string    `json:"video_id"`
	GameID          string    `json:"game_id"`
	Language        string    `json:"language"`
	Title           string    `json:"title"`
	ViewCount       int       `json:"view_count"`
	CreatedAt       Timestamp `json:"created_at"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	Duration        float64   `json:"duration"`
	VODOffset       *int      `json:"vod_offset"`
	IsFeatured      bool      `json:"is_featured"`
}

type GetClipsData []Clip

type GetClipsDownloadRequest struct {
	EditorID      string   `query:"editor_id"`
	BroadcasterID string   `query:"broadcaster_id"`
	ClipIDs       []string `query:"clip_id"`
}

type ClipDownload struct {
	ClipID               string  `json:"clip_id"`
	LandscapeDownloadURL *string `json:"landscape_download_url"`
	PortraitDownloadURL  *string `json:"portrait_download_url"`
}

type GetClipsDownloadData []ClipDownload

func (s *ClipsService) CreateClip(ctx context.Context, req CreateClipRequest) (*Response[CreateClipData], error) {
	return executeEndpoint[CreateClipData](s.client, ctx, "create-clip", req)
}

func (s *ClipsService) GetClips(ctx context.Context, req GetClipsRequest) (*Response[GetClipsData], error) {
	return executeEndpoint[GetClipsData](s.client, ctx, "get-clips", req)
}

func (s *ClipsService) GetClipsPager(req GetClipsRequest, opts ...PagerOption) (*Pager[GetClipsData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	initialCursor := ""
	baseRequest := req
	if req.Before != nil {
		initialCursor = *req.Before
		baseRequest.Before = nil
	} else if req.After != nil {
		initialCursor = *req.After
		baseRequest.After = nil
	}
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetClipsData], error) {
		if cursor == "" {
			return s.GetClips(ctx, baseRequest)
		}
		requestValue, err := plan.withCursor(baseRequest, cursor)
		if err != nil {
			return nil, err
		}
		request, ok := requestValue.(GetClipsRequest)
		if !ok {
			return nil, fmt.Errorf("helix: clips pagination request has type %T", requestValue)
		}
		return s.GetClips(ctx, request)
	}, initialCursor, opts...)
}

func (s *ExperimentalClipsService) CreateClipFromVOD(ctx context.Context, req CreateClipFromVODRequest) (*Response[CreateClipFromVODData], error) {
	return executeClipEndpoint[CreateClipFromVODData](s.client, ctx, "create-clip-from-vod", req)
}

func (s *ExperimentalClipsService) GetClipsDownload(ctx context.Context, req GetClipsDownloadRequest) (*Response[GetClipsDownloadData], error) {
	return executeClipEndpoint[GetClipsDownloadData](s.client, ctx, "get-clips-download", req)
}

func executeClipEndpoint[T any](client *Client, ctx context.Context, anchor string, query any) (*Response[T], error) {
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
	if err := validateClipCredential(credential, operation); err != nil {
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

func validateClipCredential(snapshot CredentialSnapshot, operation manifest.Operation) error {
	if !operationAllowsTokenClass(snapshot.TokenClass(), operation.TokenClasses) {
		return localCredentialAuthError(operation.OperationID)
	}
	for _, scope := range operation.Scopes {
		if scope == "" || scope == "unknown" || snapshotHasScope(snapshot, AuthorizationScope(scope)) {
			return nil
		}
	}
	if len(operation.Scopes) == 0 {
		return nil
	}
	return localCredentialAuthError(operation.OperationID)
}
