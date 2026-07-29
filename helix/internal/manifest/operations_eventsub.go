package manifest

import "net/http"

var eventsubOperations = operationsForGroup("EventSub",
	defineOperation("create-eventsub-subscription", Operation{
		Name:           "Create EventSub Subscription",
		Method:         http.MethodPost,
		Path:           "/helix/eventsub/subscriptions",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"channel:read:subscriptions"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"body": {
					{Name: "type", Type: "String", Required: true},
					{Name: "version", Type: "String", Required: true},
					{Name: "condition", Type: "Object", Required: true},
					{Name: "transport", Type: "Object", Required: true},
					{Name: "method", Type: "String", Required: true},
					{Name: "callback", Type: "String"},
					{Name: "secret", Type: "String"},
					{Name: "session_id", Type: "String"},
					{Name: "conduit_id", Type: "String"},
				},
			},
		},
		Response:   ResponseSpec{Format: "json", Status: []int{202, 400, 401, 403, 409, 410, 429}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Implementation: ImplementationSpec{
			Selector:    "Client.EventSub",
			ServiceType: "EventSubService",
			Method:      "CreateEventSubSubscription",
			Signature:   "func (s *EventSubService) CreateEventSubSubscription(ctx context.Context, req CreateEventSubSubscriptionRequest) (*Response[CreateEventSubSubscriptionData], error)",
			RequestType: "CreateEventSubSubscriptionRequest",
			DataType:    "CreateEventSubSubscriptionData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#create-eventsub-subscription",
	}),
	defineOperation("delete-eventsub-subscription", Operation{
		Name:           "Delete EventSub Subscription",
		Method:         http.MethodDelete,
		Path:           "/helix/eventsub/subscriptions",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "id", Type: "String", Required: true},
				},
			},
		},
		Response:   ResponseSpec{Format: "unknown", Status: []int{204, 400, 401, 404}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.EventSub",
			ServiceType: "EventSubService",
			Method:      "DeleteEventSubSubscription",
			Signature:   "func (s *EventSubService) DeleteEventSubSubscription(ctx context.Context, req DeleteEventSubSubscriptionRequest) (*Response[DeleteEventSubSubscriptionData], error)",
			RequestType: "DeleteEventSubSubscriptionRequest",
			DataType:    "DeleteEventSubSubscriptionData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#delete-eventsub-subscription",
	}),
	defineOperation("get-eventsub-subscriptions", Operation{
		Name:           "Get EventSub Subscriptions",
		Method:         http.MethodGet,
		Path:           "/helix/eventsub/subscriptions",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "status", Type: "String"},
					{Name: "type", Type: "String"},
					{Name: "user_id", Type: "String"},
					{Name: "subscription_id", Type: "String"},
					{Name: "conduit_id", Type: "String"},
					{Name: "after", Type: "String"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401}},
		Pagination: PaginationSpec{Shape: "cursor", CursorParameter: "after"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:       "Client.EventSub",
			ServiceType:    "EventSubService",
			Method:         "GetEventSubSubscriptions",
			Signature:      "func (s *EventSubService) GetEventSubSubscriptions(ctx context.Context, req GetEventSubSubscriptionsRequest) (*Response[GetEventSubSubscriptionsData], error)",
			PagerSignature: pagerSignature("func (s *EventSubService) GetEventSubSubscriptionsPager(req GetEventSubSubscriptionsRequest, opts ...PagerOption) (*Pager[GetEventSubSubscriptionsData], error)"),
			RequestType:    "GetEventSubSubscriptionsRequest",
			DataType:       "GetEventSubSubscriptionsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-eventsub-subscriptions",
	}),
)
