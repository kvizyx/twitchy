package helix

import "context"

type GetCharityCampaignRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type CharityAmount struct {
	Value         int    `json:"value"`
	DecimalPlaces int    `json:"decimal_places"`
	Currency      string `json:"currency"`
}

type CharityCampaign struct {
	ID                 string         `json:"id"`
	BroadcasterID      string         `json:"broadcaster_id"`
	BroadcasterLogin   string         `json:"broadcaster_login"`
	BroadcasterName    string         `json:"broadcaster_name"`
	CharityName        string         `json:"charity_name"`
	CharityDescription string         `json:"charity_description"`
	CharityLogo        string         `json:"charity_logo"`
	CharityWebsite     string         `json:"charity_website"`
	CurrentAmount      CharityAmount  `json:"current_amount"`
	TargetAmount       *CharityAmount `json:"target_amount"`
}

type GetCharityCampaignData []CharityCampaign

type GetCharityCampaignDonationsRequest struct {
	BroadcasterID string  `query:"broadcaster_id"`
	First         *int    `query:"first,omitempty"`
	After         *string `query:"after,omitempty"`
}

type CharityDonation struct {
	ID         string        `json:"id"`
	CampaignID string        `json:"campaign_id"`
	UserID     string        `json:"user_id"`
	UserLogin  string        `json:"user_login"`
	UserName   string        `json:"user_name"`
	Amount     CharityAmount `json:"amount"`
}

type GetCharityCampaignDonationsData []CharityDonation

func (s *CharityService) GetCharityCampaign(ctx context.Context, req GetCharityCampaignRequest) (*Response[GetCharityCampaignData], error) {
	return executeEndpointRequest[GetCharityCampaignData](s.client, ctx, "get-charity-campaign", req, nil, "")
}

func (s *CharityService) GetCharityCampaignDonations(ctx context.Context, req GetCharityCampaignDonationsRequest) (*Response[GetCharityCampaignDonationsData], error) {
	return executeEndpointRequest[GetCharityCampaignDonationsData](s.client, ctx, "get-charity-campaign-donations", req, nil, "")
}

func (s *CharityService) GetCharityCampaignDonationsPager(req GetCharityCampaignDonationsRequest, opts ...PagerOption) (*Pager[GetCharityCampaignDonationsData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	return newPager(func(ctx context.Context, cursor string) (*Response[GetCharityCampaignDonationsData], error) {
		next, err := plan.withCursor(req, cursor)
		if err != nil {
			return nil, err
		}
		return s.GetCharityCampaignDonations(ctx, next.(GetCharityCampaignDonationsRequest))
	}, opts...)
}
