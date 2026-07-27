package helix

import "context"

type GetContentClassificationLabelsRequest struct {
	Locale *string `query:"locale,omitempty"`
}

type ContentClassificationLabel struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type GetContentClassificationLabelsData []ContentClassificationLabel

func (s *CCLsService) GetContentClassificationLabels(ctx context.Context, req GetContentClassificationLabelsRequest) (*Response[GetContentClassificationLabelsData], error) {
	return executeEndpoint[GetContentClassificationLabelsData](s.client, ctx, "get-content-classification-labels", req)
}
