package helix

import "context"

type ExtensionBitsProductCost struct {
	Amount int    `json:"amount"`
	Type   string `json:"type"`
}

type ExtensionBitsProduct struct {
	SKU           string                   `json:"sku"`
	Cost          ExtensionBitsProductCost `json:"cost"`
	InDevelopment bool                     `json:"in_development"`
	DisplayName   string                   `json:"display_name"`
	Expiration    string                   `json:"expiration"`
	IsBroadcast   bool                     `json:"is_broadcast"`
}

type GetExtensionBitsProductsRequest struct {
	ShouldIncludeAll *bool `query:"should_include_all,omitempty"`
}

type GetExtensionBitsProductsData []ExtensionBitsProduct

type UpdateExtensionBitsProductRequest struct {
	SKU           string                   `json:"sku"`
	Cost          ExtensionBitsProductCost `json:"cost"`
	DisplayName   string                   `json:"display_name"`
	InDevelopment *bool                    `json:"in_development,omitempty"`
	Expiration    *string                  `json:"expiration,omitempty"`
	IsBroadcast   *bool                    `json:"is_broadcast,omitempty"`
}

type UpdateExtensionBitsProductData []ExtensionBitsProduct

func (s *ExtensionsService) GetReleasedExtensions(ctx context.Context, req GetReleasedExtensionsRequest) (*Response[GetReleasedExtensionsData], error) {
	return executeEndpoint[GetReleasedExtensionsData](s.client, ctx, "get-released-extensions", req)
}

func (s *ExtensionsService) GetExtensionBitsProducts(ctx context.Context, req GetExtensionBitsProductsRequest) (*Response[GetExtensionBitsProductsData], error) {
	return executeEndpoint[GetExtensionBitsProductsData](s.client, ctx, "get-extension-bits-products", req)
}

func (s *ExtensionsService) UpdateExtensionBitsProduct(ctx context.Context, req UpdateExtensionBitsProductRequest) (*Response[UpdateExtensionBitsProductData], error) {
	return executeEndpointWithBody[UpdateExtensionBitsProductData](s.client, ctx, "update-extension-bits-product", nil, req)
}
