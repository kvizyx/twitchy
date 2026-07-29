package manifest

import "net/http"

var subscriptionsOperations = operationsForGroup("Subscriptions",
	defineOperation("get-broadcaster-subscriptions", Operation{
		Name:           "Get Broadcaster Subscriptions",
		Method:         http.MethodGet,
		Path:           "/helix/subscriptions",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassApp, TokenClassUser},
		Scopes:         []string{"channel:read:subscriptions"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameter": {
					{Name: "broadcaster_id", Type: "String", Required: true},
					{Name: "user_id", Type: "String"},
					{Name: "first", Type: "String"},
					{Name: "after", Type: "String"},
					{Name: "before", Type: "String"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401}},
		Pagination: PaginationSpec{Shape: "cursor", CursorParameter: "after"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:       "Client.Subscriptions",
			ServiceType:    "SubscriptionsService",
			Method:         "GetBroadcasterSubscriptions",
			Signature:      "func (s *SubscriptionsService) GetBroadcasterSubscriptions(ctx context.Context, req GetBroadcasterSubscriptionsRequest) (*Response[GetBroadcasterSubscriptionsData], error)",
			PagerSignature: pagerSignature("func (s *SubscriptionsService) GetBroadcasterSubscriptionsPager(req GetBroadcasterSubscriptionsRequest, opts ...PagerOption) (*Pager[GetBroadcasterSubscriptionsData], error)"),
			RequestType:    "GetBroadcasterSubscriptionsRequest",
			DataType:       "GetBroadcasterSubscriptionsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-broadcaster-subscriptions",
	}),
	defineOperation("check-user-subscription", Operation{
		Name:           "Check User Subscription",
		Method:         http.MethodGet,
		Path:           "/helix/subscriptions/user",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"user:read:subscriptions"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "broadcaster_id", Type: "String", Required: true},
					{Name: "user_id", Type: "String", Required: true},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 404}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Subscriptions",
			ServiceType: "SubscriptionsService",
			Method:      "CheckUserSubscription",
			Signature:   "func (s *SubscriptionsService) CheckUserSubscription(ctx context.Context, req CheckUserSubscriptionRequest) (*Response[CheckUserSubscriptionData], error)",
			RequestType: "CheckUserSubscriptionRequest",
			DataType:    "CheckUserSubscriptionData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#check-user-subscription",
	}),
)
