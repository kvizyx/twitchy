package helix

import "context"

type GetExtensionSecretsRequest struct {
	ExtensionID string `query:"extension_id"`
}

type ExtensionSecret struct {
	Content   string    `json:"content"`
	ActiveAt  Timestamp `json:"active_at"`
	ExpiresAt Timestamp `json:"expires_at"`
}

type ExtensionSecrets struct {
	FormatVersion int               `json:"format_version"`
	Secrets       []ExtensionSecret `json:"secrets"`
}

type GetExtensionSecretsData []ExtensionSecrets

type CreateExtensionSecretRequest struct {
	ExtensionID string `query:"extension_id"`
	Delay       *int   `query:"delay,omitempty"`
}

type CreateExtensionSecretData []ExtensionSecrets

func (s *ExtensionsService) GetExtensionSecrets(ctx context.Context, req GetExtensionSecretsRequest) (*Response[GetExtensionSecretsData], error) {
	return executeExtensionEndpoint[GetExtensionSecretsData](s.client, ctx, "get-extension-secrets", req, nil)
}

func (s *ExtensionsService) CreateExtensionSecret(ctx context.Context, req CreateExtensionSecretRequest) (*Response[CreateExtensionSecretData], error) {
	return executeExtensionEndpoint[CreateExtensionSecretData](s.client, ctx, "create-extension-secret", req, nil)
}
