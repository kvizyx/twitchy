package manifest

import "net/http"

var predictionsOperations = operationsForGroup("Predictions",
	defineOperation("get-predictions", Operation{
		Name:           "Get Predictions",
		Method:         http.MethodGet,
		Path:           "/helix/predictions",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:read:predictions", "channel:manage:predictions"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "broadcaster_id", Type: "String", Required: true},
					{Name: "id", Type: "String"},
					{Name: "first", Type: "String"},
					{Name: "after", Type: "String"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401}},
		Pagination: PaginationSpec{Shape: "cursor", CursorParameter: "after"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:       "Client.Predictions",
			ServiceType:    "PredictionsService",
			Method:         "GetPredictions",
			Signature:      "func (s *PredictionsService) GetPredictions(ctx context.Context, req GetPredictionsRequest) (*Response[GetPredictionsData], error)",
			PagerSignature: pagerSignature("func (s *PredictionsService) GetPredictionsPager(req GetPredictionsRequest, opts ...PagerOption) (*Pager[GetPredictionsData], error)"),
			RequestType:    "GetPredictionsRequest",
			DataType:       "GetPredictionsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-predictions",
	}),
	defineOperation("create-prediction", Operation{
		Name:           "Create Prediction",
		Method:         http.MethodPost,
		Path:           "/helix/predictions",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:manage:predictions"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"body": {
					{Name: "broadcaster_id", Type: "String", Required: true},
					{Name: "title", Type: "String", Required: true},
					{Name: "outcomes", Type: "Object[]", Required: true},
					{Name: "title", Type: "String", Required: true},
					{Name: "prediction_window", Type: "Integer", Required: true},
				},
			},
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 429}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Implementation: ImplementationSpec{
			Selector:    "Client.Predictions",
			ServiceType: "PredictionsService",
			Method:      "CreatePrediction",
			Signature:   "func (s *PredictionsService) CreatePrediction(ctx context.Context, req CreatePredictionRequest) (*Response[CreatePredictionData], error)",
			RequestType: "CreatePredictionRequest",
			DataType:    "CreatePredictionData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#create-prediction",
	}),
	defineOperation("end-prediction", Operation{
		Name:           "End Prediction",
		Method:         http.MethodPatch,
		Path:           "/helix/predictions",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:manage:predictions"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"body": {
					{Name: "broadcaster_id", Type: "String", Required: true},
					{Name: "id", Type: "String", Required: true},
					{Name: "status", Type: "String", Required: true},
					{Name: "winning_outcome_id", Type: "String"},
				},
			},
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 404}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Predictions",
			ServiceType: "PredictionsService",
			Method:      "EndPrediction",
			Signature:   "func (s *PredictionsService) EndPrediction(ctx context.Context, req EndPredictionRequest) (*Response[EndPredictionData], error)",
			RequestType: "EndPredictionRequest",
			DataType:    "EndPredictionData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#end-prediction",
	}),
)
