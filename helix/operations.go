package helix

import "context"

type emptyRequest struct{}
type emptyData struct{}

type StartCommercialRequest = emptyRequest
type StartCommercialData = emptyData
type GetAdScheduleRequest = emptyRequest
type GetAdScheduleData = emptyData
type SnoozeNextAdRequest = emptyRequest
type SnoozeNextAdData = emptyData
type GetExtensionAnalyticsRequest = emptyRequest
type GetExtensionAnalyticsData = emptyData
type GetGameAnalyticsRequest = emptyRequest
type GetGameAnalyticsData = emptyData
type GetBitsLeaderboardRequest = emptyRequest
type GetBitsLeaderboardData = emptyData
type GetCheermotesRequest = emptyRequest
type GetCheermotesData = emptyData
type GetCustomPowerUpRequest = emptyRequest
type GetCustomPowerUpData = emptyData
type GetExtensionTransactionsRequest = emptyRequest
type GetExtensionTransactionsData = emptyData
type GetChannelInformationRequest = emptyRequest
type GetChannelInformationData = emptyData
type ModifyChannelInformationRequest = emptyRequest
type ModifyChannelInformationData = emptyData
type GetChannelEditorsRequest = emptyRequest
type GetChannelEditorsData = emptyData
type GetFollowedChannelsRequest = emptyRequest
type GetFollowedChannelsData = emptyData
type GetChannelFollowersRequest = emptyRequest
type GetChannelFollowersData = emptyData
type CreateCustomRewardsRequest = emptyRequest
type CreateCustomRewardsData = emptyData
type DeleteCustomRewardRequest = emptyRequest
type DeleteCustomRewardData = emptyData
type GetCustomRewardRequest = emptyRequest
type GetCustomRewardData = emptyData
type GetCustomRewardRedemptionRequest = emptyRequest
type GetCustomRewardRedemptionData = emptyData
type UpdateCustomRewardRequest = emptyRequest
type UpdateCustomRewardData = emptyData
type UpdateRedemptionStatusRequest = emptyRequest
type UpdateRedemptionStatusData = emptyData
type GetCharityCampaignRequest = emptyRequest
type GetCharityCampaignData = emptyData
type GetCharityCampaignDonationsRequest = emptyRequest
type GetCharityCampaignDonationsData = emptyData
type GetChattersRequest = emptyRequest
type GetChattersData = emptyData
type GetChannelEmotesRequest = emptyRequest
type GetChannelEmotesData = emptyData
type GetGlobalEmotesRequest = emptyRequest
type GetGlobalEmotesData = emptyData
type GetEmoteSetsRequest = emptyRequest
type GetEmoteSetsData = emptyData
type GetChannelChatBadgesRequest = emptyRequest
type GetChannelChatBadgesData = emptyData
type GetGlobalChatBadgesRequest = emptyRequest
type GetGlobalChatBadgesData = emptyData
type GetChatSettingsRequest = emptyRequest
type GetChatSettingsData = emptyData
type GetSharedChatSessionRequest = emptyRequest
type GetSharedChatSessionData = emptyData
type GetUserEmotesRequest = emptyRequest
type GetUserEmotesData = emptyData
type UpdateChatSettingsRequest = emptyRequest
type UpdateChatSettingsData = emptyData
type SendChatAnnouncementRequest = emptyRequest
type SendChatAnnouncementData = emptyData
type SendShoutoutRequest = emptyRequest
type SendShoutoutData = emptyData
type SendChatMessageRequest = emptyRequest
type SendChatMessageData = emptyData
type GetPinnedChatMessageRequest = emptyRequest
type GetPinnedChatMessageData = emptyData
type PinChatMessageRequest = emptyRequest
type PinChatMessageData = emptyData
type UpdatePinnedChatMessageRequest = emptyRequest
type UpdatePinnedChatMessageData = emptyData
type UnpinChatMessageRequest = emptyRequest
type UnpinChatMessageData = emptyData
type GetUserChatColorRequest = emptyRequest
type GetUserChatColorData = emptyData
type UpdateUserChatColorRequest = emptyRequest
type UpdateUserChatColorData = emptyData
type CreateClipRequest = emptyRequest
type CreateClipData = emptyData
type CreateClipFromVODRequest = emptyRequest
type CreateClipFromVODData = emptyData
type GetClipsRequest = emptyRequest
type GetClipsData = emptyData
type GetClipsDownloadRequest = emptyRequest
type GetClipsDownloadData = emptyData
type GetConduitsRequest = emptyRequest
type GetConduitsData = emptyData
type CreateConduitsRequest = emptyRequest
type CreateConduitsData = emptyData
type UpdateConduitsRequest = emptyRequest
type UpdateConduitsData = emptyData
type DeleteConduitRequest = emptyRequest
type DeleteConduitData = emptyData
type GetConduitShardsRequest = emptyRequest
type GetConduitShardsData = emptyData
type UpdateConduitShardsRequest = emptyRequest
type UpdateConduitShardsData = emptyData
type GetContentClassificationLabelsRequest = emptyRequest
type GetContentClassificationLabelsData = emptyData
type GetDropsEntitlementsRequest = emptyRequest
type GetDropsEntitlementsData = emptyData
type UpdateDropsEntitlementsRequest = emptyRequest
type UpdateDropsEntitlementsData = emptyData
type GetExtensionConfigurationSegmentRequest = emptyRequest
type GetExtensionConfigurationSegmentData = emptyData
type SetExtensionConfigurationSegmentRequest = emptyRequest
type SetExtensionConfigurationSegmentData = emptyData
type SetExtensionRequiredConfigurationRequest = emptyRequest
type SetExtensionRequiredConfigurationData = emptyData
type SendExtensionPubSubMessageRequest = emptyRequest
type SendExtensionPubSubMessageData = emptyData
type GetExtensionLiveChannelsRequest = emptyRequest
type GetExtensionLiveChannelsData = emptyData
type GetExtensionSecretsRequest = emptyRequest
type GetExtensionSecretsData = emptyData
type CreateExtensionSecretRequest = emptyRequest
type CreateExtensionSecretData = emptyData
type SendExtensionChatMessageRequest = emptyRequest
type SendExtensionChatMessageData = emptyData
type GetExtensionsRequest = emptyRequest
type GetExtensionsData = emptyData
type GetReleasedExtensionsRequest = emptyRequest
type GetReleasedExtensionsData = emptyData
type GetExtensionBitsProductsRequest = emptyRequest
type GetExtensionBitsProductsData = emptyData
type UpdateExtensionBitsProductRequest = emptyRequest
type UpdateExtensionBitsProductData = emptyData
type CreateEventSubSubscriptionRequest = emptyRequest
type CreateEventSubSubscriptionData = emptyData
type DeleteEventSubSubscriptionRequest = emptyRequest
type DeleteEventSubSubscriptionData = emptyData
type GetEventSubSubscriptionsRequest = emptyRequest
type GetEventSubSubscriptionsData = emptyData
type GetTopGamesRequest = emptyRequest
type GetTopGamesData = emptyData
type GetGamesRequest = emptyRequest
type GetGamesData = emptyData
type GetCreatorGoalsRequest = emptyRequest
type GetCreatorGoalsData = emptyData
type GetChannelGuestStarSettingsRequest = emptyRequest
type GetChannelGuestStarSettingsData = emptyData
type UpdateChannelGuestStarSettingsRequest = emptyRequest
type UpdateChannelGuestStarSettingsData = emptyData
type GetGuestStarSessionRequest = emptyRequest
type GetGuestStarSessionData = emptyData
type CreateGuestStarSessionRequest = emptyRequest
type CreateGuestStarSessionData = emptyData
type EndGuestStarSessionRequest = emptyRequest
type EndGuestStarSessionData = emptyData
type GetGuestStarInvitesRequest = emptyRequest
type GetGuestStarInvitesData = emptyData
type SendGuestStarInviteRequest = emptyRequest
type SendGuestStarInviteData = emptyData
type DeleteGuestStarInviteRequest = emptyRequest
type DeleteGuestStarInviteData = emptyData
type AssignGuestStarSlotRequest = emptyRequest
type AssignGuestStarSlotData = emptyData
type UpdateGuestStarSlotRequest = emptyRequest
type UpdateGuestStarSlotData = emptyData
type DeleteGuestStarSlotRequest = emptyRequest
type DeleteGuestStarSlotData = emptyData
type UpdateGuestStarSlotSettingsRequest = emptyRequest
type UpdateGuestStarSlotSettingsData = emptyData
type GetHypeTrainStatusRequest = emptyRequest
type GetHypeTrainStatusData = emptyData
type CheckAutoModStatusRequest = emptyRequest
type CheckAutoModStatusData = emptyData
type ManageHeldAutoModMessagesRequest = emptyRequest
type ManageHeldAutoModMessagesData = emptyData
type GetAutoModSettingsRequest = emptyRequest
type GetAutoModSettingsData = emptyData
type UpdateAutoModSettingsRequest = emptyRequest
type UpdateAutoModSettingsData = emptyData
type GetBannedUsersRequest = emptyRequest
type GetBannedUsersData = emptyData
type BanUserRequest = emptyRequest
type BanUserData = emptyData
type UnbanUserRequest = emptyRequest
type UnbanUserData = emptyData
type GetUnbanRequestsRequest = emptyRequest
type GetUnbanRequestsData = emptyData
type ResolveUnbanRequestsRequest = emptyRequest
type ResolveUnbanRequestsData = emptyData
type GetBlockedTermsRequest = emptyRequest
type GetBlockedTermsData = emptyData
type AddBlockedTermRequest = emptyRequest
type AddBlockedTermData = emptyData
type RemoveBlockedTermRequest = emptyRequest
type RemoveBlockedTermData = emptyData
type DeleteChatMessagesRequest = emptyRequest
type DeleteChatMessagesData = emptyData
type GetModeratedChannelsRequest = emptyRequest
type GetModeratedChannelsData = emptyData
type GetModeratorsRequest = emptyRequest
type GetModeratorsData = emptyData
type AddChannelModeratorRequest = emptyRequest
type AddChannelModeratorData = emptyData
type RemoveChannelModeratorRequest = emptyRequest
type RemoveChannelModeratorData = emptyData
type GetVIPsRequest = emptyRequest
type GetVIPsData = emptyData
type AddChannelVIPRequest = emptyRequest
type AddChannelVIPData = emptyData
type RemoveChannelVIPRequest = emptyRequest
type RemoveChannelVIPData = emptyData
type UpdateShieldModeStatusRequest = emptyRequest
type UpdateShieldModeStatusData = emptyData
type GetShieldModeStatusRequest = emptyRequest
type GetShieldModeStatusData = emptyData
type WarnChatUserRequest = emptyRequest
type WarnChatUserData = emptyData
type AddSuspiciousStatusToChatUserRequest = emptyRequest
type AddSuspiciousStatusToChatUserData = emptyData
type RemoveSuspiciousStatusFromChatUserRequest = emptyRequest
type RemoveSuspiciousStatusFromChatUserData = emptyData
type GetPollsRequest = emptyRequest
type GetPollsData = emptyData
type CreatePollRequest = emptyRequest
type CreatePollData = emptyData
type EndPollRequest = emptyRequest
type EndPollData = emptyData
type GetPredictionsRequest = emptyRequest
type GetPredictionsData = emptyData
type CreatePredictionRequest = emptyRequest
type CreatePredictionData = emptyData
type EndPredictionRequest = emptyRequest
type EndPredictionData = emptyData
type StartRaidRequest = emptyRequest
type StartRaidData = emptyData
type CancelRaidRequest = emptyRequest
type CancelRaidData = emptyData
type GetChannelStreamScheduleRequest = emptyRequest
type GetChannelStreamScheduleData = emptyData
type GetChannelICalendarRequest = emptyRequest
type GetChannelICalendarData = emptyData
type GetChannelICalendar = emptyData
type UpdateChannelStreamScheduleRequest = emptyRequest
type UpdateChannelStreamScheduleData = emptyData
type CreateChannelStreamScheduleSegmentRequest = emptyRequest
type CreateChannelStreamScheduleSegmentData = emptyData
type UpdateChannelStreamScheduleSegmentRequest = emptyRequest
type UpdateChannelStreamScheduleSegmentData = emptyData
type DeleteChannelStreamScheduleSegmentRequest = emptyRequest
type DeleteChannelStreamScheduleSegmentData = emptyData
type SearchCategoriesRequest = emptyRequest
type SearchCategoriesData = emptyData
type SearchChannelsRequest = emptyRequest
type SearchChannelsData = emptyData
type GetStreamKeyRequest = emptyRequest
type GetStreamKeyData = emptyData
type GetStreamsRequest = emptyRequest
type GetStreamsData = emptyData
type GetFollowedStreamsRequest = emptyRequest
type GetFollowedStreamsData = emptyData
type CreateStreamMarkerRequest = emptyRequest
type CreateStreamMarkerData = emptyData
type GetStreamMarkersRequest = emptyRequest
type GetStreamMarkersData = emptyData
type GetBroadcasterSubscriptionsRequest = emptyRequest
type GetBroadcasterSubscriptionsData = emptyData
type CheckUserSubscriptionRequest = emptyRequest
type CheckUserSubscriptionData = emptyData
type GetAllStreamTagsRequest = emptyRequest
type GetAllStreamTagsData = emptyData
type GetStreamTagsRequest = emptyRequest
type GetStreamTagsData = emptyData
type GetChannelTeamsRequest = emptyRequest
type GetChannelTeamsData = emptyData
type GetTeamsRequest = emptyRequest
type GetTeamsData = emptyData
type GetUsersRequest = emptyRequest
type GetUsersData = emptyData
type UpdateUserRequest = emptyRequest
type UpdateUserData = emptyData
type GetAuthorizationByUserRequest = emptyRequest
type GetAuthorizationByUserData = emptyData
type GetUserBlockListRequest = emptyRequest
type GetUserBlockListData = emptyData
type BlockUserRequest = emptyRequest
type BlockUserData = emptyData
type UnblockUserRequest = emptyRequest
type UnblockUserData = emptyData
type GetUserExtensionsRequest = emptyRequest
type GetUserExtensionsData = emptyData
type GetUserActiveExtensionsRequest = emptyRequest
type GetUserActiveExtensionsData = emptyData
type UpdateUserExtensionsRequest = emptyRequest
type UpdateUserExtensionsData = emptyData
type GetVideosRequest = emptyRequest
type GetVideosData = emptyData
type DeleteVideosRequest = emptyRequest
type DeleteVideosData = emptyData
type SendWhisperRequest = emptyRequest
type SendWhisperData = emptyData

func (s *AdsService) StartCommercial(context.Context, StartCommercialRequest) (*Response[StartCommercialData], error) {
	return nil, ErrInvalidClient
}
func (s *AdsService) GetAdSchedule(context.Context, GetAdScheduleRequest) (*Response[GetAdScheduleData], error) {
	return nil, ErrInvalidClient
}
func (s *AdsService) SnoozeNextAd(context.Context, SnoozeNextAdRequest) (*Response[SnoozeNextAdData], error) {
	return nil, ErrInvalidClient
}
func (s *AnalyticsService) GetExtensionAnalytics(context.Context, GetExtensionAnalyticsRequest) (*Response[GetExtensionAnalyticsData], error) {
	return nil, ErrInvalidClient
}
func (s *AnalyticsService) GetExtensionAnalyticsPager(GetExtensionAnalyticsRequest, ...PagerOption) (*Pager[GetExtensionAnalyticsData], error) {
	return nil, ErrInvalidClient
}
func (s *AnalyticsService) GetGameAnalytics(context.Context, GetGameAnalyticsRequest) (*Response[GetGameAnalyticsData], error) {
	return nil, ErrInvalidClient
}
func (s *AnalyticsService) GetGameAnalyticsPager(GetGameAnalyticsRequest, ...PagerOption) (*Pager[GetGameAnalyticsData], error) {
	return nil, ErrInvalidClient
}
func (s *BitsService) GetBitsLeaderboard(context.Context, GetBitsLeaderboardRequest) (*Response[GetBitsLeaderboardData], error) {
	return nil, ErrInvalidClient
}
func (s *BitsService) GetCheermotes(context.Context, GetCheermotesRequest) (*Response[GetCheermotesData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalBitsService) GetCustomPowerUp(context.Context, GetCustomPowerUpRequest) (*Response[GetCustomPowerUpData], error) {
	return nil, ErrInvalidClient
}
func (s *BitsService) GetExtensionTransactions(context.Context, GetExtensionTransactionsRequest) (*Response[GetExtensionTransactionsData], error) {
	return nil, ErrInvalidClient
}
func (s *BitsService) GetExtensionTransactionsPager(GetExtensionTransactionsRequest, ...PagerOption) (*Pager[GetExtensionTransactionsData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelsService) GetChannelInformation(context.Context, GetChannelInformationRequest) (*Response[GetChannelInformationData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelsService) ModifyChannelInformation(context.Context, ModifyChannelInformationRequest) (*Response[ModifyChannelInformationData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelsService) GetChannelEditors(context.Context, GetChannelEditorsRequest) (*Response[GetChannelEditorsData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelsService) GetFollowedChannels(context.Context, GetFollowedChannelsRequest) (*Response[GetFollowedChannelsData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelsService) GetFollowedChannelsPager(GetFollowedChannelsRequest, ...PagerOption) (*Pager[GetFollowedChannelsData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelsService) GetChannelFollowers(context.Context, GetChannelFollowersRequest) (*Response[GetChannelFollowersData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelsService) GetChannelFollowersPager(GetChannelFollowersRequest, ...PagerOption) (*Pager[GetChannelFollowersData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelPointsService) CreateCustomRewards(context.Context, CreateCustomRewardsRequest) (*Response[CreateCustomRewardsData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelPointsService) DeleteCustomReward(context.Context, DeleteCustomRewardRequest) (*Response[DeleteCustomRewardData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelPointsService) GetCustomReward(context.Context, GetCustomRewardRequest) (*Response[GetCustomRewardData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelPointsService) GetCustomRewardRedemption(context.Context, GetCustomRewardRedemptionRequest) (*Response[GetCustomRewardRedemptionData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelPointsService) GetCustomRewardRedemptionPager(GetCustomRewardRedemptionRequest, ...PagerOption) (*Pager[GetCustomRewardRedemptionData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelPointsService) UpdateCustomReward(context.Context, UpdateCustomRewardRequest) (*Response[UpdateCustomRewardData], error) {
	return nil, ErrInvalidClient
}
func (s *ChannelPointsService) UpdateRedemptionStatus(context.Context, UpdateRedemptionStatusRequest) (*Response[UpdateRedemptionStatusData], error) {
	return nil, ErrInvalidClient
}
func (s *CharityService) GetCharityCampaign(context.Context, GetCharityCampaignRequest) (*Response[GetCharityCampaignData], error) {
	return nil, ErrInvalidClient
}
func (s *CharityService) GetCharityCampaignDonations(context.Context, GetCharityCampaignDonationsRequest) (*Response[GetCharityCampaignDonationsData], error) {
	return nil, ErrInvalidClient
}
func (s *CharityService) GetCharityCampaignDonationsPager(GetCharityCampaignDonationsRequest, ...PagerOption) (*Pager[GetCharityCampaignDonationsData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) GetChatters(context.Context, GetChattersRequest) (*Response[GetChattersData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) GetChattersPager(GetChattersRequest, ...PagerOption) (*Pager[GetChattersData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) GetChannelEmotes(context.Context, GetChannelEmotesRequest) (*Response[GetChannelEmotesData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) GetGlobalEmotes(context.Context, GetGlobalEmotesRequest) (*Response[GetGlobalEmotesData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) GetEmoteSets(context.Context, GetEmoteSetsRequest) (*Response[GetEmoteSetsData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) GetChannelChatBadges(context.Context, GetChannelChatBadgesRequest) (*Response[GetChannelChatBadgesData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) GetGlobalChatBadges(context.Context, GetGlobalChatBadgesRequest) (*Response[GetGlobalChatBadgesData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) GetChatSettings(context.Context, GetChatSettingsRequest) (*Response[GetChatSettingsData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) GetSharedChatSession(context.Context, GetSharedChatSessionRequest) (*Response[GetSharedChatSessionData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) GetUserEmotes(context.Context, GetUserEmotesRequest) (*Response[GetUserEmotesData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) GetUserEmotesPager(GetUserEmotesRequest, ...PagerOption) (*Pager[GetUserEmotesData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) UpdateChatSettings(context.Context, UpdateChatSettingsRequest) (*Response[UpdateChatSettingsData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) SendChatAnnouncement(context.Context, SendChatAnnouncementRequest) (*Response[SendChatAnnouncementData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) SendShoutout(context.Context, SendShoutoutRequest) (*Response[SendShoutoutData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) SendChatMessage(context.Context, SendChatMessageRequest) (*Response[SendChatMessageData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalChatService) GetPinnedChatMessage(context.Context, GetPinnedChatMessageRequest) (*Response[GetPinnedChatMessageData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalChatService) PinChatMessage(context.Context, PinChatMessageRequest) (*Response[PinChatMessageData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalChatService) UpdatePinnedChatMessage(context.Context, UpdatePinnedChatMessageRequest) (*Response[UpdatePinnedChatMessageData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalChatService) UnpinChatMessage(context.Context, UnpinChatMessageRequest) (*Response[UnpinChatMessageData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) GetUserChatColor(context.Context, GetUserChatColorRequest) (*Response[GetUserChatColorData], error) {
	return nil, ErrInvalidClient
}
func (s *ChatService) UpdateUserChatColor(context.Context, UpdateUserChatColorRequest) (*Response[UpdateUserChatColorData], error) {
	return nil, ErrInvalidClient
}
func (s *ClipsService) CreateClip(context.Context, CreateClipRequest) (*Response[CreateClipData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalClipsService) CreateClipFromVOD(context.Context, CreateClipFromVODRequest) (*Response[CreateClipFromVODData], error) {
	return nil, ErrInvalidClient
}
func (s *ClipsService) GetClips(context.Context, GetClipsRequest) (*Response[GetClipsData], error) {
	return nil, ErrInvalidClient
}
func (s *ClipsService) GetClipsPager(GetClipsRequest, ...PagerOption) (*Pager[GetClipsData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalClipsService) GetClipsDownload(context.Context, GetClipsDownloadRequest) (*Response[GetClipsDownloadData], error) {
	return nil, ErrInvalidClient
}
func (s *ConduitsService) GetConduits(context.Context, GetConduitsRequest) (*Response[GetConduitsData], error) {
	return nil, ErrInvalidClient
}
func (s *ConduitsService) CreateConduits(context.Context, CreateConduitsRequest) (*Response[CreateConduitsData], error) {
	return nil, ErrInvalidClient
}
func (s *ConduitsService) UpdateConduits(context.Context, UpdateConduitsRequest) (*Response[UpdateConduitsData], error) {
	return nil, ErrInvalidClient
}
func (s *ConduitsService) DeleteConduit(context.Context, DeleteConduitRequest) (*Response[DeleteConduitData], error) {
	return nil, ErrInvalidClient
}
func (s *ConduitsService) GetConduitShards(context.Context, GetConduitShardsRequest) (*Response[GetConduitShardsData], error) {
	return nil, ErrInvalidClient
}
func (s *ConduitsService) GetConduitShardsPager(GetConduitShardsRequest, ...PagerOption) (*Pager[GetConduitShardsData], error) {
	return nil, ErrInvalidClient
}
func (s *ConduitsService) UpdateConduitShards(context.Context, UpdateConduitShardsRequest) (*Response[UpdateConduitShardsData], error) {
	return nil, ErrInvalidClient
}
func (s *CCLsService) GetContentClassificationLabels(context.Context, GetContentClassificationLabelsRequest) (*Response[GetContentClassificationLabelsData], error) {
	return nil, ErrInvalidClient
}
func (s *EntitlementsService) GetDropsEntitlements(context.Context, GetDropsEntitlementsRequest) (*Response[GetDropsEntitlementsData], error) {
	return nil, ErrInvalidClient
}
func (s *EntitlementsService) GetDropsEntitlementsPager(GetDropsEntitlementsRequest, ...PagerOption) (*Pager[GetDropsEntitlementsData], error) {
	return nil, ErrInvalidClient
}
func (s *EntitlementsService) UpdateDropsEntitlements(context.Context, UpdateDropsEntitlementsRequest) (*Response[UpdateDropsEntitlementsData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) GetExtensionConfigurationSegment(context.Context, GetExtensionConfigurationSegmentRequest) (*Response[GetExtensionConfigurationSegmentData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) SetExtensionConfigurationSegment(context.Context, SetExtensionConfigurationSegmentRequest) (*Response[SetExtensionConfigurationSegmentData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) SetExtensionRequiredConfiguration(context.Context, SetExtensionRequiredConfigurationRequest) (*Response[SetExtensionRequiredConfigurationData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) SendExtensionPubSubMessage(context.Context, SendExtensionPubSubMessageRequest) (*Response[SendExtensionPubSubMessageData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) GetExtensionLiveChannels(context.Context, GetExtensionLiveChannelsRequest) (*Response[GetExtensionLiveChannelsData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) GetExtensionLiveChannelsPager(GetExtensionLiveChannelsRequest, ...PagerOption) (*Pager[GetExtensionLiveChannelsData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) GetExtensionSecrets(context.Context, GetExtensionSecretsRequest) (*Response[GetExtensionSecretsData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) CreateExtensionSecret(context.Context, CreateExtensionSecretRequest) (*Response[CreateExtensionSecretData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) SendExtensionChatMessage(context.Context, SendExtensionChatMessageRequest) (*Response[SendExtensionChatMessageData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) GetExtensions(context.Context, GetExtensionsRequest) (*Response[GetExtensionsData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) GetReleasedExtensions(context.Context, GetReleasedExtensionsRequest) (*Response[GetReleasedExtensionsData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) GetExtensionBitsProducts(context.Context, GetExtensionBitsProductsRequest) (*Response[GetExtensionBitsProductsData], error) {
	return nil, ErrInvalidClient
}
func (s *ExtensionsService) UpdateExtensionBitsProduct(context.Context, UpdateExtensionBitsProductRequest) (*Response[UpdateExtensionBitsProductData], error) {
	return nil, ErrInvalidClient
}
func (s *EventSubService) CreateEventSubSubscription(context.Context, CreateEventSubSubscriptionRequest) (*Response[CreateEventSubSubscriptionData], error) {
	return nil, ErrInvalidClient
}
func (s *EventSubService) DeleteEventSubSubscription(context.Context, DeleteEventSubSubscriptionRequest) (*Response[DeleteEventSubSubscriptionData], error) {
	return nil, ErrInvalidClient
}
func (s *EventSubService) GetEventSubSubscriptions(context.Context, GetEventSubSubscriptionsRequest) (*Response[GetEventSubSubscriptionsData], error) {
	return nil, ErrInvalidClient
}
func (s *EventSubService) GetEventSubSubscriptionsPager(GetEventSubSubscriptionsRequest, ...PagerOption) (*Pager[GetEventSubSubscriptionsData], error) {
	return nil, ErrInvalidClient
}
func (s *GamesService) GetTopGames(context.Context, GetTopGamesRequest) (*Response[GetTopGamesData], error) {
	return nil, ErrInvalidClient
}
func (s *GamesService) GetTopGamesPager(GetTopGamesRequest, ...PagerOption) (*Pager[GetTopGamesData], error) {
	return nil, ErrInvalidClient
}
func (s *GamesService) GetGames(context.Context, GetGamesRequest) (*Response[GetGamesData], error) {
	return nil, ErrInvalidClient
}
func (s *GoalsService) GetCreatorGoals(context.Context, GetCreatorGoalsRequest) (*Response[GetCreatorGoalsData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalGuestStarService) GetChannelGuestStarSettings(context.Context, GetChannelGuestStarSettingsRequest) (*Response[GetChannelGuestStarSettingsData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalGuestStarService) UpdateChannelGuestStarSettings(context.Context, UpdateChannelGuestStarSettingsRequest) (*Response[UpdateChannelGuestStarSettingsData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalGuestStarService) GetGuestStarSession(context.Context, GetGuestStarSessionRequest) (*Response[GetGuestStarSessionData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalGuestStarService) CreateGuestStarSession(context.Context, CreateGuestStarSessionRequest) (*Response[CreateGuestStarSessionData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalGuestStarService) EndGuestStarSession(context.Context, EndGuestStarSessionRequest) (*Response[EndGuestStarSessionData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalGuestStarService) GetGuestStarInvites(context.Context, GetGuestStarInvitesRequest) (*Response[GetGuestStarInvitesData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalGuestStarService) SendGuestStarInvite(context.Context, SendGuestStarInviteRequest) (*Response[SendGuestStarInviteData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalGuestStarService) DeleteGuestStarInvite(context.Context, DeleteGuestStarInviteRequest) (*Response[DeleteGuestStarInviteData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalGuestStarService) AssignGuestStarSlot(context.Context, AssignGuestStarSlotRequest) (*Response[AssignGuestStarSlotData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalGuestStarService) UpdateGuestStarSlot(context.Context, UpdateGuestStarSlotRequest) (*Response[UpdateGuestStarSlotData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalGuestStarService) DeleteGuestStarSlot(context.Context, DeleteGuestStarSlotRequest) (*Response[DeleteGuestStarSlotData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalGuestStarService) UpdateGuestStarSlotSettings(context.Context, UpdateGuestStarSlotSettingsRequest) (*Response[UpdateGuestStarSlotSettingsData], error) {
	return nil, ErrInvalidClient
}
func (s *HypeTrainService) GetHypeTrainStatus(context.Context, GetHypeTrainStatusRequest) (*Response[GetHypeTrainStatusData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) CheckAutoModStatus(context.Context, CheckAutoModStatusRequest) (*Response[CheckAutoModStatusData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) ManageHeldAutoModMessages(context.Context, ManageHeldAutoModMessagesRequest) (*Response[ManageHeldAutoModMessagesData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetAutoModSettings(context.Context, GetAutoModSettingsRequest) (*Response[GetAutoModSettingsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) UpdateAutoModSettings(context.Context, UpdateAutoModSettingsRequest) (*Response[UpdateAutoModSettingsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetBannedUsers(context.Context, GetBannedUsersRequest) (*Response[GetBannedUsersData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetBannedUsersPager(GetBannedUsersRequest, ...PagerOption) (*Pager[GetBannedUsersData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) BanUser(context.Context, BanUserRequest) (*Response[BanUserData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) UnbanUser(context.Context, UnbanUserRequest) (*Response[UnbanUserData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetUnbanRequests(context.Context, GetUnbanRequestsRequest) (*Response[GetUnbanRequestsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetUnbanRequestsPager(GetUnbanRequestsRequest, ...PagerOption) (*Pager[GetUnbanRequestsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) ResolveUnbanRequests(context.Context, ResolveUnbanRequestsRequest) (*Response[ResolveUnbanRequestsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetBlockedTerms(context.Context, GetBlockedTermsRequest) (*Response[GetBlockedTermsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetBlockedTermsPager(GetBlockedTermsRequest, ...PagerOption) (*Pager[GetBlockedTermsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) AddBlockedTerm(context.Context, AddBlockedTermRequest) (*Response[AddBlockedTermData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) RemoveBlockedTerm(context.Context, RemoveBlockedTermRequest) (*Response[RemoveBlockedTermData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) DeleteChatMessages(context.Context, DeleteChatMessagesRequest) (*Response[DeleteChatMessagesData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetModeratedChannels(context.Context, GetModeratedChannelsRequest) (*Response[GetModeratedChannelsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetModeratedChannelsPager(GetModeratedChannelsRequest, ...PagerOption) (*Pager[GetModeratedChannelsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetModerators(context.Context, GetModeratorsRequest) (*Response[GetModeratorsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetModeratorsPager(GetModeratorsRequest, ...PagerOption) (*Pager[GetModeratorsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) AddChannelModerator(context.Context, AddChannelModeratorRequest) (*Response[AddChannelModeratorData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) RemoveChannelModerator(context.Context, RemoveChannelModeratorRequest) (*Response[RemoveChannelModeratorData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetVIPs(context.Context, GetVIPsRequest) (*Response[GetVIPsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetVIPsPager(GetVIPsRequest, ...PagerOption) (*Pager[GetVIPsData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) AddChannelVIP(context.Context, AddChannelVIPRequest) (*Response[AddChannelVIPData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) RemoveChannelVIP(context.Context, RemoveChannelVIPRequest) (*Response[RemoveChannelVIPData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) UpdateShieldModeStatus(context.Context, UpdateShieldModeStatusRequest) (*Response[UpdateShieldModeStatusData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) GetShieldModeStatus(context.Context, GetShieldModeStatusRequest) (*Response[GetShieldModeStatusData], error) {
	return nil, ErrInvalidClient
}
func (s *ModerationService) WarnChatUser(context.Context, WarnChatUserRequest) (*Response[WarnChatUserData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalModerationService) AddSuspiciousStatusToChatUser(context.Context, AddSuspiciousStatusToChatUserRequest) (*Response[AddSuspiciousStatusToChatUserData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalModerationService) RemoveSuspiciousStatusFromChatUser(context.Context, RemoveSuspiciousStatusFromChatUserRequest) (*Response[RemoveSuspiciousStatusFromChatUserData], error) {
	return nil, ErrInvalidClient
}
func (s *PollsService) GetPolls(context.Context, GetPollsRequest) (*Response[GetPollsData], error) {
	return nil, ErrInvalidClient
}
func (s *PollsService) GetPollsPager(GetPollsRequest, ...PagerOption) (*Pager[GetPollsData], error) {
	return nil, ErrInvalidClient
}
func (s *PollsService) CreatePoll(context.Context, CreatePollRequest) (*Response[CreatePollData], error) {
	return nil, ErrInvalidClient
}
func (s *PollsService) EndPoll(context.Context, EndPollRequest) (*Response[EndPollData], error) {
	return nil, ErrInvalidClient
}
func (s *PredictionsService) GetPredictions(context.Context, GetPredictionsRequest) (*Response[GetPredictionsData], error) {
	return nil, ErrInvalidClient
}
func (s *PredictionsService) GetPredictionsPager(GetPredictionsRequest, ...PagerOption) (*Pager[GetPredictionsData], error) {
	return nil, ErrInvalidClient
}
func (s *PredictionsService) CreatePrediction(context.Context, CreatePredictionRequest) (*Response[CreatePredictionData], error) {
	return nil, ErrInvalidClient
}
func (s *PredictionsService) EndPrediction(context.Context, EndPredictionRequest) (*Response[EndPredictionData], error) {
	return nil, ErrInvalidClient
}
func (s *RaidsService) StartRaid(context.Context, StartRaidRequest) (*Response[StartRaidData], error) {
	return nil, ErrInvalidClient
}
func (s *RaidsService) CancelRaid(context.Context, CancelRaidRequest) (*Response[CancelRaidData], error) {
	return nil, ErrInvalidClient
}
func (s *ScheduleService) GetChannelStreamSchedule(context.Context, GetChannelStreamScheduleRequest) (*Response[GetChannelStreamScheduleData], error) {
	return nil, ErrInvalidClient
}
func (s *ScheduleService) GetChannelStreamSchedulePager(GetChannelStreamScheduleRequest, ...PagerOption) (*Pager[GetChannelStreamScheduleData], error) {
	return nil, ErrInvalidClient
}
func (s *ScheduleService) GetChannelICalendar(context.Context, GetChannelICalendarRequest) (*Response[GetChannelICalendar], error) {
	return nil, ErrInvalidClient
}
func (s *ScheduleService) UpdateChannelStreamSchedule(context.Context, UpdateChannelStreamScheduleRequest) (*Response[UpdateChannelStreamScheduleData], error) {
	return nil, ErrInvalidClient
}
func (s *ScheduleService) CreateChannelStreamScheduleSegment(context.Context, CreateChannelStreamScheduleSegmentRequest) (*Response[CreateChannelStreamScheduleSegmentData], error) {
	return nil, ErrInvalidClient
}
func (s *ScheduleService) UpdateChannelStreamScheduleSegment(context.Context, UpdateChannelStreamScheduleSegmentRequest) (*Response[UpdateChannelStreamScheduleSegmentData], error) {
	return nil, ErrInvalidClient
}
func (s *ScheduleService) DeleteChannelStreamScheduleSegment(context.Context, DeleteChannelStreamScheduleSegmentRequest) (*Response[DeleteChannelStreamScheduleSegmentData], error) {
	return nil, ErrInvalidClient
}
func (s *SearchService) SearchCategories(context.Context, SearchCategoriesRequest) (*Response[SearchCategoriesData], error) {
	return nil, ErrInvalidClient
}
func (s *SearchService) SearchCategoriesPager(SearchCategoriesRequest, ...PagerOption) (*Pager[SearchCategoriesData], error) {
	return nil, ErrInvalidClient
}
func (s *SearchService) SearchChannels(context.Context, SearchChannelsRequest) (*Response[SearchChannelsData], error) {
	return nil, ErrInvalidClient
}
func (s *SearchService) SearchChannelsPager(SearchChannelsRequest, ...PagerOption) (*Pager[SearchChannelsData], error) {
	return nil, ErrInvalidClient
}
func (s *StreamsService) GetStreamKey(context.Context, GetStreamKeyRequest) (*Response[GetStreamKeyData], error) {
	return nil, ErrInvalidClient
}
func (s *StreamsService) GetStreams(context.Context, GetStreamsRequest) (*Response[GetStreamsData], error) {
	return nil, ErrInvalidClient
}
func (s *StreamsService) GetStreamsPager(GetStreamsRequest, ...PagerOption) (*Pager[GetStreamsData], error) {
	return nil, ErrInvalidClient
}
func (s *StreamsService) GetFollowedStreams(context.Context, GetFollowedStreamsRequest) (*Response[GetFollowedStreamsData], error) {
	return nil, ErrInvalidClient
}
func (s *StreamsService) GetFollowedStreamsPager(GetFollowedStreamsRequest, ...PagerOption) (*Pager[GetFollowedStreamsData], error) {
	return nil, ErrInvalidClient
}
func (s *StreamsService) CreateStreamMarker(context.Context, CreateStreamMarkerRequest) (*Response[CreateStreamMarkerData], error) {
	return nil, ErrInvalidClient
}
func (s *StreamsService) GetStreamMarkers(context.Context, GetStreamMarkersRequest) (*Response[GetStreamMarkersData], error) {
	return nil, ErrInvalidClient
}
func (s *StreamsService) GetStreamMarkersPager(GetStreamMarkersRequest, ...PagerOption) (*Pager[GetStreamMarkersData], error) {
	return nil, ErrInvalidClient
}
func (s *SubscriptionsService) GetBroadcasterSubscriptions(context.Context, GetBroadcasterSubscriptionsRequest) (*Response[GetBroadcasterSubscriptionsData], error) {
	return nil, ErrInvalidClient
}
func (s *SubscriptionsService) GetBroadcasterSubscriptionsPager(GetBroadcasterSubscriptionsRequest, ...PagerOption) (*Pager[GetBroadcasterSubscriptionsData], error) {
	return nil, ErrInvalidClient
}
func (s *SubscriptionsService) CheckUserSubscription(context.Context, CheckUserSubscriptionRequest) (*Response[CheckUserSubscriptionData], error) {
	return nil, ErrInvalidClient
}
func (s *TagsService) GetAllStreamTags(context.Context, GetAllStreamTagsRequest) (*Response[GetAllStreamTagsData], error) {
	return nil, ErrInvalidClient
}
func (s *TagsService) GetAllStreamTagsPager(GetAllStreamTagsRequest, ...PagerOption) (*Pager[GetAllStreamTagsData], error) {
	return nil, ErrInvalidClient
}
func (s *TagsService) GetStreamTags(context.Context, GetStreamTagsRequest) (*Response[GetStreamTagsData], error) {
	return nil, ErrInvalidClient
}
func (s *TeamsService) GetChannelTeams(context.Context, GetChannelTeamsRequest) (*Response[GetChannelTeamsData], error) {
	return nil, ErrInvalidClient
}
func (s *TeamsService) GetTeams(context.Context, GetTeamsRequest) (*Response[GetTeamsData], error) {
	return nil, ErrInvalidClient
}
func (s *UsersService) GetUsers(context.Context, GetUsersRequest) (*Response[GetUsersData], error) {
	return nil, ErrInvalidClient
}
func (s *UsersService) UpdateUser(context.Context, UpdateUserRequest) (*Response[UpdateUserData], error) {
	return nil, ErrInvalidClient
}
func (s *ExperimentalUsersService) GetAuthorizationByUser(context.Context, GetAuthorizationByUserRequest) (*Response[GetAuthorizationByUserData], error) {
	return nil, ErrInvalidClient
}
func (s *UsersService) GetUserBlockList(context.Context, GetUserBlockListRequest) (*Response[GetUserBlockListData], error) {
	return nil, ErrInvalidClient
}
func (s *UsersService) GetUserBlockListPager(GetUserBlockListRequest, ...PagerOption) (*Pager[GetUserBlockListData], error) {
	return nil, ErrInvalidClient
}
func (s *UsersService) BlockUser(context.Context, BlockUserRequest) (*Response[BlockUserData], error) {
	return nil, ErrInvalidClient
}
func (s *UsersService) UnblockUser(context.Context, UnblockUserRequest) (*Response[UnblockUserData], error) {
	return nil, ErrInvalidClient
}
func (s *UsersService) GetUserExtensions(context.Context, GetUserExtensionsRequest) (*Response[GetUserExtensionsData], error) {
	return nil, ErrInvalidClient
}
func (s *UsersService) GetUserActiveExtensions(context.Context, GetUserActiveExtensionsRequest) (*Response[GetUserActiveExtensionsData], error) {
	return nil, ErrInvalidClient
}
func (s *UsersService) UpdateUserExtensions(context.Context, UpdateUserExtensionsRequest) (*Response[UpdateUserExtensionsData], error) {
	return nil, ErrInvalidClient
}
func (s *VideosService) GetVideos(context.Context, GetVideosRequest) (*Response[GetVideosData], error) {
	return nil, ErrInvalidClient
}
func (s *VideosService) GetVideosPager(GetVideosRequest, ...PagerOption) (*Pager[GetVideosData], error) {
	return nil, ErrInvalidClient
}
func (s *VideosService) DeleteVideos(context.Context, DeleteVideosRequest) (*Response[DeleteVideosData], error) {
	return nil, ErrInvalidClient
}
func (s *WhispersService) SendWhisper(context.Context, SendWhisperRequest) (*Response[SendWhisperData], error) {
	return nil, ErrInvalidClient
}
