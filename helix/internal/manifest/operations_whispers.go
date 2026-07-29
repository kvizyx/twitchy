package manifest

import "net/http"

var whispersOperations = operationsForGroup("Whispers",
	defineOperation("send-whisper", Operation{
		Name:           "Send Whisper",
		Method:         http.MethodPost,
		Path:           "/helix/whispers",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"user:manage:whispers"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"body_parameters": {
					{Name: "message", Type: "String", Required: true},
				},
				"query_parameters": {
					{Name: "from_user_id", Type: "String", Required: true},
					{Name: "to_user_id", Type: "String", Required: true},
				},
			},
		},
		Response:   ResponseSpec{Format: "unknown", Status: []int{204, 400, 401, 403, 404, 429}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Implementation: ImplementationSpec{
			Selector:    "Client.Whispers",
			ServiceType: "WhispersService",
			Method:      "SendWhisper",
			Signature:   "func (s *WhispersService) SendWhisper(ctx context.Context, req SendWhisperRequest) (*Response[SendWhisperData], error)",
			RequestType: "SendWhisperRequest",
			DataType:    "SendWhisperData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#send-whisper",
	}),
)
