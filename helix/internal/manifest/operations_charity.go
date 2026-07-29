package manifest

import "net/http"

var charityOperations = operationsForGroup("Charity",
	defineOperation("get-charity-campaign", Operation{
		Name:           "Get Charity Campaign",
		Method:         http.MethodGet,
		Path:           "/helix/charity/campaigns",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:read:charity"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "broadcaster_id", Type: "String", Required: true},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 403}},
		Pagination: PaginationSpec{Shape: "none", CursorParameter: "unknown"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:    "Client.Charity",
			ServiceType: "CharityService",
			Method:      "GetCharityCampaign",
			Signature:   "func (s *CharityService) GetCharityCampaign(ctx context.Context, req GetCharityCampaignRequest) (*Response[GetCharityCampaignData], error)",
			RequestType: "GetCharityCampaignRequest",
			DataType:    "GetCharityCampaignData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-charity-campaign",
	}),
	defineOperation("get-charity-campaign-donations", Operation{
		Name:           "Get Charity Campaign Donations",
		Method:         http.MethodGet,
		Path:           "/helix/charity/donations",
		Stability:      StabilityStable,
		TokenClasses:   []TokenClass{TokenClassUser},
		Scopes:         []string{"channel:read:charity"},
		SubjectBinding: "unknown",
		Request: RequestSpec{
			Locations: map[string][]RequestField{
				"query_parameters": {
					{Name: "broadcaster_id", Type: "String", Required: true},
					{Name: "first", Type: "Integer"},
					{Name: "after", Type: "String"},
				},
			},
			BodyReconstructible: true,
		},
		Response:   ResponseSpec{Format: "json", Status: []int{200, 400, 401, 403}},
		Pagination: PaginationSpec{Shape: "cursor", CursorParameter: "after"},
		Replay:     ReplaySpec{BucketWaitable: true},
		Implementation: ImplementationSpec{
			Selector:       "Client.Charity",
			ServiceType:    "CharityService",
			Method:         "GetCharityCampaignDonations",
			Signature:      "func (s *CharityService) GetCharityCampaignDonations(ctx context.Context, req GetCharityCampaignDonationsRequest) (*Response[GetCharityCampaignDonationsData], error)",
			PagerSignature: pagerSignature("func (s *CharityService) GetCharityCampaignDonationsPager(req GetCharityCampaignDonationsRequest, opts ...PagerOption) (*Pager[GetCharityCampaignDonationsData], error)"),
			RequestType:    "GetCharityCampaignDonationsRequest",
			DataType:       "GetCharityCampaignDonationsData",
		},
		Source: "https://dev.twitch.tv/docs/api/reference/#get-charity-campaign-donations",
	}),
)
