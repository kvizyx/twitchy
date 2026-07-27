package helix

import "context"

type RedemptionStatus string

const (
	RedemptionStatusCanceled    RedemptionStatus = "CANCELED"
	RedemptionStatusFulfilled   RedemptionStatus = "FULFILLED"
	RedemptionStatusUnfulfilled RedemptionStatus = "UNFULFILLED"
)

type RedemptionSort string

const (
	RedemptionSortOldest RedemptionSort = "OLDEST"
	RedemptionSortNewest RedemptionSort = "NEWEST"
)

type CustomRewardImage struct {
	URL1x string `json:"url_1x"`
	URL2x string `json:"url_2x"`
	URL4x string `json:"url_4x"`
}

type CustomRewardLimitSetting struct {
	IsEnabled    bool  `json:"is_enabled"`
	MaxPerStream int64 `json:"max_per_stream"`
}

type CustomRewardUserLimitSetting struct {
	IsEnabled           bool  `json:"is_enabled"`
	MaxPerUserPerStream int64 `json:"max_per_user_per_stream"`
}

type CustomRewardCooldownSetting struct {
	IsEnabled             bool  `json:"is_enabled"`
	GlobalCooldownSeconds int64 `json:"global_cooldown_seconds"`
}

type CustomReward struct {
	BroadcasterID                     string                       `json:"broadcaster_id"`
	BroadcasterLogin                  string                       `json:"broadcaster_login"`
	BroadcasterName                   string                       `json:"broadcaster_name"`
	ID                                string                       `json:"id"`
	Title                             string                       `json:"title"`
	Prompt                            string                       `json:"prompt"`
	Cost                              int                          `json:"cost"`
	Image                             *CustomRewardImage           `json:"image"`
	DefaultImage                      CustomRewardImage            `json:"default_image"`
	BackgroundColor                   string                       `json:"background_color"`
	IsEnabled                         bool                         `json:"is_enabled"`
	IsUserInputRequired               bool                         `json:"is_user_input_required"`
	MaxPerStreamSetting               CustomRewardLimitSetting     `json:"max_per_stream_setting"`
	MaxPerUserPerStreamSetting        CustomRewardUserLimitSetting `json:"max_per_user_per_stream_setting"`
	GlobalCooldownSetting             CustomRewardCooldownSetting  `json:"global_cooldown_setting"`
	IsPaused                          bool                         `json:"is_paused"`
	IsInStock                         bool                         `json:"is_in_stock"`
	ShouldRedemptionsSkipRequestQueue bool                         `json:"should_redemptions_skip_request_queue"`
	RedemptionsRedeemedCurrentStream  *int                         `json:"redemptions_redeemed_current_stream"`
	CooldownExpiresAt                 *Timestamp                   `json:"cooldown_expires_at"`
}

type CustomRedemptionReward struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Prompt string `json:"prompt"`
	Cost   int64  `json:"cost"`
}

type CustomRewardRedemption struct {
	BroadcasterID    string                 `json:"broadcaster_id"`
	BroadcasterLogin string                 `json:"broadcaster_login"`
	BroadcasterName  string                 `json:"broadcaster_name"`
	ID               string                 `json:"id"`
	UserLogin        string                 `json:"user_login"`
	UserID           string                 `json:"user_id"`
	UserName         string                 `json:"user_name"`
	UserInput        string                 `json:"user_input"`
	Status           RedemptionStatus       `json:"status"`
	RedeemedAt       Timestamp              `json:"redeemed_at"`
	Reward           CustomRedemptionReward `json:"reward"`
}

type CreateCustomRewardsRequest struct {
	BroadcasterID                     string  `query:"broadcaster_id" json:"-"`
	Title                             string  `query:"-" json:"title"`
	Cost                              int64   `query:"-" json:"cost"`
	Prompt                            *string `query:"-" json:"prompt,omitempty"`
	IsEnabled                         *bool   `query:"-" json:"is_enabled,omitempty"`
	BackgroundColor                   *string `query:"-" json:"background_color,omitempty"`
	IsUserInputRequired               *bool   `query:"-" json:"is_user_input_required,omitempty"`
	IsMaxPerStreamEnabled             *bool   `query:"-" json:"is_max_per_stream_enabled,omitempty"`
	MaxPerStream                      *int    `query:"-" json:"max_per_stream,omitempty"`
	IsMaxPerUserPerStreamEnabled      *bool   `query:"-" json:"is_max_per_user_per_stream_enabled,omitempty"`
	MaxPerUserPerStream               *int    `query:"-" json:"max_per_user_per_stream,omitempty"`
	IsGlobalCooldownEnabled           *bool   `query:"-" json:"is_global_cooldown_enabled,omitempty"`
	GlobalCooldownSeconds             *int    `query:"-" json:"global_cooldown_seconds,omitempty"`
	ShouldRedemptionsSkipRequestQueue *bool   `query:"-" json:"should_redemptions_skip_request_queue,omitempty"`
}

type CreateCustomRewardsData []CustomReward

type DeleteCustomRewardRequest struct {
	BroadcasterID string `query:"broadcaster_id" json:"-"`
	ID            string `query:"id" json:"-"`
}

type DeleteCustomRewardData struct{}

type GetCustomRewardRequest struct {
	BroadcasterID         string   `query:"broadcaster_id"`
	IDs                   []string `query:"id,omitempty"`
	OnlyManageableRewards *bool    `query:"only_manageable_rewards,omitempty"`
}

type GetCustomRewardData []CustomReward

type GetCustomRewardRedemptionRequest struct {
	BroadcasterID string           `query:"broadcaster_id"`
	RewardID      string           `query:"reward_id"`
	Status        RedemptionStatus `query:"status"`
	IDs           []string         `query:"id,omitempty"`
	Sort          *RedemptionSort  `query:"sort,omitempty"`
	After         *string          `query:"after,omitempty"`
	First         *int             `query:"first,omitempty"`
}

type GetCustomRewardRedemptionData []CustomRewardRedemption

type UpdateCustomRewardRequest struct {
	BroadcasterID                     string  `query:"broadcaster_id" json:"-"`
	ID                                string  `query:"id" json:"-"`
	Title                             *string `query:"-" json:"title,omitempty"`
	Prompt                            *string `query:"-" json:"prompt,omitempty"`
	Cost                              *int64  `query:"-" json:"cost,omitempty"`
	BackgroundColor                   *string `query:"-" json:"background_color,omitempty"`
	IsEnabled                         *bool   `query:"-" json:"is_enabled,omitempty"`
	IsUserInputRequired               *bool   `query:"-" json:"is_user_input_required,omitempty"`
	IsMaxPerStreamEnabled             *bool   `query:"-" json:"is_max_per_stream_enabled,omitempty"`
	MaxPerStream                      *int64  `query:"-" json:"max_per_stream,omitempty"`
	IsMaxPerUserPerStreamEnabled      *bool   `query:"-" json:"is_max_per_user_per_stream_enabled,omitempty"`
	MaxPerUserPerStream               *int64  `query:"-" json:"max_per_user_per_stream,omitempty"`
	IsGlobalCooldownEnabled           *bool   `query:"-" json:"is_global_cooldown_enabled,omitempty"`
	GlobalCooldownSeconds             *int64  `query:"-" json:"global_cooldown_seconds,omitempty"`
	IsPaused                          *bool   `query:"-" json:"is_paused,omitempty"`
	ShouldRedemptionsSkipRequestQueue *bool   `query:"-" json:"should_redemptions_skip_request_queue,omitempty"`
}

type UpdateCustomRewardData []CustomReward

type UpdateRedemptionStatusRequest struct {
	IDs           []string         `query:"id" json:"-"`
	BroadcasterID string           `query:"broadcaster_id" json:"-"`
	RewardID      string           `query:"reward_id" json:"-"`
	Status        RedemptionStatus `query:"-" json:"status"`
}

type UpdateRedemptionStatusData []CustomRewardRedemption

func (s *ChannelPointsService) CreateCustomRewards(ctx context.Context, req CreateCustomRewardsRequest) (*Response[CreateCustomRewardsData], error) {
	return executeEndpointRequest[CreateCustomRewardsData](s.client, ctx, "create-custom-rewards", req, req, "")
}

func (s *ChannelPointsService) DeleteCustomReward(ctx context.Context, req DeleteCustomRewardRequest) (*Response[DeleteCustomRewardData], error) {
	return executeEndpointRequest[DeleteCustomRewardData](s.client, ctx, "delete-custom-reward", req, nil, "")
}

func (s *ChannelPointsService) GetCustomReward(ctx context.Context, req GetCustomRewardRequest) (*Response[GetCustomRewardData], error) {
	return executeEndpointRequest[GetCustomRewardData](s.client, ctx, "get-custom-reward", req, nil, "")
}

func (s *ChannelPointsService) GetCustomRewardRedemption(ctx context.Context, req GetCustomRewardRedemptionRequest) (*Response[GetCustomRewardRedemptionData], error) {
	return executeEndpointRequest[GetCustomRewardRedemptionData](s.client, ctx, "get-custom-reward-redemption", req, nil, "")
}

func (s *ChannelPointsService) GetCustomRewardRedemptionPager(req GetCustomRewardRedemptionRequest, opts ...PagerOption) (*Pager[GetCustomRewardRedemptionData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	return newPager(func(ctx context.Context, cursor string) (*Response[GetCustomRewardRedemptionData], error) {
		next, err := plan.withCursor(req, cursor)
		if err != nil {
			return nil, err
		}
		return s.GetCustomRewardRedemption(ctx, next.(GetCustomRewardRedemptionRequest))
	}, opts...)
}

func (s *ChannelPointsService) UpdateCustomReward(ctx context.Context, req UpdateCustomRewardRequest) (*Response[UpdateCustomRewardData], error) {
	return executeEndpointRequest[UpdateCustomRewardData](s.client, ctx, "update-custom-reward", req, req, "")
}

func (s *ChannelPointsService) UpdateRedemptionStatus(ctx context.Context, req UpdateRedemptionStatusRequest) (*Response[UpdateRedemptionStatusData], error) {
	return executeEndpointRequest[UpdateRedemptionStatusData](s.client, ctx, "update-redemption-status", req, req, "")
}
