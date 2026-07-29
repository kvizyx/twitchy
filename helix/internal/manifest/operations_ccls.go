package manifest

import "net/http"

var cclsOperations = operationsForGroup("CCLs",
	defineOperation("get-content-classification-labels", Operation{
		Name:           "Get Content Classification Labels",
		Method:         http.MethodGet,
		Path:           "/helix/content_classification_labels",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "locale", Type: "String"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", StatusUnknown: true},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.CCLs",
			ServiceType: "CCLsService",
			Method:      "GetContentClassificationLabels",
			Signature:   "func (s *CCLsService) GetContentClassificationLabels(ctx context.Context, req GetContentClassificationLabelsRequest) (*Response[GetContentClassificationLabelsData], error)",
			RequestType: "GetContentClassificationLabelsRequest",
			DataType:    "GetContentClassificationLabelsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-content-classification-labels",
	}),
)
