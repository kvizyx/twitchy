package helix

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type GetExtensionConfigurationSegmentRequest struct {
	BroadcasterID string   `query:"broadcaster_id,omitempty"`
	ExtensionID   string   `query:"extension_id"`
	Segment       []string `query:"segment"`
}

type ExtensionConfigurationSegment struct {
	Segment       string `json:"segment"`
	BroadcasterID string `json:"broadcaster_id,omitempty"`
	Content       string `json:"content"`
	Version       string `json:"version"`
}

type GetExtensionConfigurationSegmentData []ExtensionConfigurationSegment

type SetExtensionConfigurationSegmentRequest struct {
	ExtensionID   string  `json:"extension_id"`
	Segment       string  `json:"segment"`
	BroadcasterID *string `json:"broadcaster_id,omitempty"`
	Content       *string `json:"content,omitempty"`
	Version       *string `json:"version,omitempty"`
}

type SetExtensionConfigurationSegmentData struct{}

type SetExtensionRequiredConfigurationRequest struct {
	BroadcasterID         string `query:"broadcaster_id" json:"-"`
	ExtensionID           string `json:"extension_id"`
	ExtensionVersion      string `json:"extension_version"`
	RequiredConfiguration string `json:"required_configuration"`
}

type SetExtensionRequiredConfigurationData struct{}

type ExtensionView struct {
	ViewerURL              string `json:"viewer_url"`
	Height                 int    `json:"height"`
	CanLinkExternalContent bool   `json:"can_link_external_content"`
	AspectRatioX           int    `json:"aspect_ratio_x"`
	AspectRatioY           int    `json:"aspect_ratio_y"`
	Autoscale              bool   `json:"autoscale"`
	ScalePixels            int    `json:"scale_pixels"`
	TargetHeight           int    `json:"target_height"`
}

type ExtensionViews struct {
	Mobile       *ExtensionView `json:"mobile"`
	Panel        *ExtensionView `json:"panel"`
	VideoOverlay *ExtensionView `json:"video_overlay"`
	Component    *ExtensionView `json:"component"`
	Config       *ExtensionView `json:"config"`
}

type Extension struct {
	AuthorName                string            `json:"author_name"`
	BitsEnabled               bool              `json:"bits_enabled"`
	CanInstall                bool              `json:"can_install"`
	ConfigurationLocation     string            `json:"configuration_location"`
	Description               string            `json:"description"`
	EULATOSURL                string            `json:"eula_tos_url"`
	HasChatSupport            bool              `json:"has_chat_support"`
	IconURL                   string            `json:"icon_url"`
	IconURLs                  map[string]string `json:"icon_urls"`
	ID                        string            `json:"id"`
	Name                      string            `json:"name"`
	PrivacyPolicyURL          string            `json:"privacy_policy_url"`
	RequestIdentityLink       bool              `json:"request_identity_link"`
	ScreenshotURLs            []string          `json:"screenshot_urls"`
	State                     string            `json:"state"`
	SubscriptionsSupportLevel string            `json:"subscriptions_support_level"`
	Summary                   string            `json:"summary"`
	SupportEmail              string            `json:"support_email"`
	Version                   string            `json:"version"`
	ViewerSummary             string            `json:"viewer_summary"`
	Views                     ExtensionViews    `json:"views"`
	AllowlistedConfigURLs     []string          `json:"allowlisted_config_urls"`
	AllowlistedPanelURLs      []string          `json:"allowlisted_panel_urls"`
}

type GetExtensionsRequest struct {
	ExtensionID      string  `query:"extension_id"`
	ExtensionVersion *string `query:"extension_version,omitempty"`
}

type GetExtensionsData []Extension

type GetReleasedExtensionsRequest struct {
	ExtensionID      string  `query:"extension_id"`
	ExtensionVersion *string `query:"extension_version,omitempty"`
}

type GetReleasedExtensionsData []Extension

type extensionAuthRoundTripper struct {
	base  http.RoundTripper
	token string
}

func (transport *extensionAuthRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	clone := request.Clone(request.Context())
	clone.Header.Set("Authorization", "Extension "+transport.token)
	return transport.base.RoundTrip(clone)
}

func executeExtensionEndpoint[T any](client *Client, ctx context.Context, anchor string, query, body any) (*Response[T], error) {
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
	if credential.TokenClass() != TokenClassExtension {
		return nil, localCredentialAuthError(operation.OperationID)
	}
	if err := validateCredentialForOperation(credential, operation, "", ""); err != nil {
		return nil, err
	}
	request, err := buildRequest(requestSpec{Context: ctx, Method: operation.Method, URL: client.endpointURL(operation.Path), Query: query, Body: body})
	if err != nil {
		return nil, fmt.Errorf("build %s request: %w", operation.OperationID, err)
	}
	request.Header.Set("User-Agent", client.userAgent)
	executor := *client.executor
	httpClient := *client.executor.httpClient
	baseTransport := httpClient.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}
	httpClient.Transport = &extensionAuthRoundTripper{base: baseTransport, token: credential.AccessToken()}
	executor.httpClient = &httpClient
	response, meta, err := executor.execute(ctx, request, operation, credential)
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

func (s *ExtensionsService) GetExtensionConfigurationSegment(ctx context.Context, req GetExtensionConfigurationSegmentRequest) (*Response[GetExtensionConfigurationSegmentData], error) {
	return executeExtensionEndpoint[GetExtensionConfigurationSegmentData](s.client, ctx, "get-extension-configuration-segment", req, nil)
}

func (s *ExtensionsService) SetExtensionConfigurationSegment(ctx context.Context, req SetExtensionConfigurationSegmentRequest) (*Response[SetExtensionConfigurationSegmentData], error) {
	return executeExtensionEndpoint[SetExtensionConfigurationSegmentData](s.client, ctx, "set-extension-configuration-segment", nil, req)
}

func (s *ExtensionsService) SetExtensionRequiredConfiguration(ctx context.Context, req SetExtensionRequiredConfigurationRequest) (*Response[SetExtensionRequiredConfigurationData], error) {
	query := struct {
		BroadcasterID string `query:"broadcaster_id"`
	}{BroadcasterID: req.BroadcasterID}
	return executeExtensionEndpoint[SetExtensionRequiredConfigurationData](s.client, ctx, "set-extension-required-configuration", query, req)
}

func (s *ExtensionsService) GetExtensions(ctx context.Context, req GetExtensionsRequest) (*Response[GetExtensionsData], error) {
	return executeExtensionEndpoint[GetExtensionsData](s.client, ctx, "get-extensions", req, nil)
}
