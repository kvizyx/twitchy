package eventsub

// Handler is a handler for generic event.
type Handler[Event any, Metadata any] func(Event, Metadata)

// callback is a store for user's EventSub callbacks.
type callback[Metadata any] struct {
	onDuplicate      func(Metadata)
	onUndefinedEvent func(RawEvent, Metadata)

	onAutomodMessageHold                          Handler[AutomodMessageHoldEvent, Metadata]
	onAutomodMessageHoldV2                        Handler[AutomodMessageHoldEventV2, Metadata]
	onAutomodMessageUpdate                        Handler[AutomodMessageUpdateEvent, Metadata]
	onAutomodMessageUpdateV2                      Handler[AutomodMessageUpdateEventV2, Metadata]
	onAutomodSettingsUpdate                       Handler[AutomodSettingsUpdateEvent, Metadata]
	onAutomodTermsUpdate                          Handler[AutomodTermsUpdateEvent, Metadata]
	onChannelBitsUse                              Handler[ChannelBitsUseEvent, Metadata]
	onChannelUpdate                               Handler[ChannelUpdateEvent, Metadata]
	onChannelFollow                               Handler[ChannelFollowEvent, Metadata]
	onChannelAdBreakBegin                         Handler[ChannelAdBreakBeginEvent, Metadata]
	onChannelChatClear                            Handler[ChannelChatClearEvent, Metadata]
	onChannelChatClearUserMessages                Handler[ChannelChatClearUserMessagesEvent, Metadata]
	onChannelChatMessage                          Handler[ChannelChatMessageEvent, Metadata]
	onConduitShardDisabled                        Handler[ConduitShardDisabledEvent, Metadata]
	onChannelBan                                  Handler[ChannelBanEvent, Metadata]
	onChannelUnban                                Handler[ChannelUnbanEvent, Metadata]
	onChannelChatNotification                     Handler[ChannelChatNotificationEvent, Metadata]
	onChannelModeratorAdd                         Handler[ChannelModeratorAddEvent, Metadata]
	onChannelModeratorRemove                      Handler[ChannelModeratorRemoveEvent, Metadata]
	onChannelPollBegin                            Handler[ChannelPollBeginEvent, Metadata]
	onChannelPollProgress                         Handler[ChannelPollProgressEvent, Metadata]
	onChannelPollEnd                              Handler[ChannelPollEndEvent, Metadata]
	onChannelPredictionBegin                      Handler[ChannelPredictionBeginEvent, Metadata]
	onChannelPredictionProgress                   Handler[ChannelPredictionProgressEvent, Metadata]
	onChannelPredictionLock                       Handler[ChannelPredictionLockEvent, Metadata]
	onChannelPredictionEnd                        Handler[ChannelPredictionEndEvent, Metadata]
	onChannelRaid                                 Handler[ChannelRaidEvent, Metadata]
	onChannelPointsCustomRewardRedemptionAdd      Handler[ChannelPointsCustomRewardRedemptionAddEvent, Metadata]
	onChannelPointsCustomRewardRedemptionUpdate   Handler[ChannelPointsCustomRewardRedemptionUpdateEvent, Metadata]
	onChannelPointsAutomaticRewardRedemptionAdd   Handler[ChannelPointsAutomaticRewardRedemptionAddEvent, Metadata]
	onChannelPointsAutomaticRewardRedemptionAddV2 Handler[ChannelPointsAutomaticRewardRedemptionAddEventV2, Metadata]
	onChannelPointsCustomRewardAdd                Handler[ChannelPointsCustomRewardAddEvent, Metadata]
	onChannelPointsCustomRewardUpdate             Handler[ChannelPointsCustomRewardUpdateEvent, Metadata]
	onChannelPointsCustomRewardRemove             Handler[ChannelPointsCustomRewardRemoveEvent, Metadata]
	onStreamOffline                               Handler[StreamOfflineEvent, Metadata]
	onStreamOnline                                Handler[StreamOnlineEvent, Metadata]
	onChannelSubscribe                            Handler[ChannelSubscribeEvent, Metadata]
	onChannelSubscriptionEnd                      Handler[ChannelSubscriptionEndEvent, Metadata]
	onChannelSubscriptionMessage                  Handler[ChannelSubscriptionMessageEvent, Metadata]
	onChannelSubscriptionGift                     Handler[ChannelSubscriptionGiftEvent, Metadata]
	onChannelUnbanRequestCreate                   Handler[ChannelUnbanRequestCreateEvent, Metadata]
	onChannelUnbanRequestResolve                  Handler[ChannelUnbanRequestResolveEvent, Metadata]
	onUserUpdate                                  Handler[UserUpdateEvent, Metadata]
	onChannelVipAdd                               Handler[ChannelVipAddEvent, Metadata]
	onChannelVipRemove                            Handler[ChannelVipRemoveEvent, Metadata]
	onChannelChatMessageDelete                    Handler[ChannelChatMessageDeleteEvent, Metadata]
	onUserAuthorizationRevoke                     Handler[UserAuthorizationRevokeEvent, Metadata]
	onChannelModerateV2                           Handler[ChannelModerateEventV2, Metadata]
	onChannelChatSettingsUpdate                   Handler[ChannelChatSettingsUpdateEvent, Metadata]
	onChannelChatUserMessageHold                  Handler[ChannelChatUserMessageHoldEvent, Metadata]
	onChannelChatUserMessageUpdate                Handler[ChannelChatUserMessageUpdateEvent, Metadata]
	onChannelCustomPowerUpRedemptionAdd           Handler[ChannelCustomPowerUpRedemptionAddEvent, Metadata]
	onChannelSharedChatSessionBegin               Handler[ChannelSharedChatSessionBeginEvent, Metadata]
	onChannelSharedChatSessionUpdate              Handler[ChannelSharedChatSessionUpdateEvent, Metadata]
	onChannelSharedChatSessionEnd                 Handler[ChannelSharedChatSessionEndEvent, Metadata]
	onChannelSuspiciousUserMessage                Handler[ChannelSuspiciousUserMessageEvent, Metadata]
	onChannelSuspiciousUserUpdate                 Handler[ChannelSuspiciousUserUpdateEvent, Metadata]
	onChannelWarningAcknowledge                   Handler[ChannelWarningAcknowledgeEvent, Metadata]
	onChannelWarningSend                          Handler[ChannelWarningSendEvent, Metadata]
	onChannelCheer                                Handler[ChannelCheerEvent, Metadata]
	onCharityCampaignDonate                       Handler[CharityCampaignDonateEvent, Metadata]
	onCharityCampaignStart                        Handler[CharityCampaignStartEvent, Metadata]
	onCharityCampaignProgress                     Handler[CharityCampaignProgressEvent, Metadata]
	onCharityCampaignStop                         Handler[CharityCampaignStopEvent, Metadata]
	onDropEntitlementGrant                        Handler[DropEntitlementGrantEvent, Metadata]
	onExtensionBitsTransactionCreate              Handler[ExtensionBitsTransactionCreateEvent, Metadata]
	onChannelGoalBegin                            Handler[ChannelGoalBeginEvent, Metadata]
	onChannelGoalProgress                         Handler[ChannelGoalProgressEvent, Metadata]
	onChannelGoalEnd                              Handler[ChannelGoalEndEvent, Metadata]
	onHypeTrainBegin                              Handler[HypeTrainBeginEvent, Metadata]
	onHypeTrainProgress                           Handler[HypeTrainProgressEvent, Metadata]
	onHypeTrainEnd                                Handler[HypeTrainEndEvent, Metadata]
	onShieldModeBegin                             Handler[ShieldModeBeginEvent, Metadata]
	onShieldModeEnd                               Handler[ShieldModeEndEvent, Metadata]
	onShoutoutCreate                              Handler[ShoutoutCreateEvent, Metadata]
	onShoutoutReceive                             Handler[ShoutoutReceiveEvent, Metadata]
	onUserAuthorizationGrant                      Handler[UserAuthorizationGrantEvent, Metadata]
	onUserWhisperMessage                          Handler[WhisperReceivedEvent, Metadata]
	onGuestStarSessionBegin                       Handler[ChannelGuestStarSessionBeginEvent, Metadata]
	onGuestStarSessionEnd                         Handler[ChannelGuestStarSessionEndEvent, Metadata]
	onGuestStarGuestUpdate                        Handler[ChannelGuestStarGuestUpdateEvent, Metadata]
	onGuestStarSettingsUpdate                     Handler[ChannelGuestStarSettingsUpdateEvent, Metadata]
}

// OnDuplicate invokes when duplicate message is caught (this is not necessarily an event).
func (c *callback[Metadata]) OnDuplicate(onDuplicate func(Metadata)) {
	c.onDuplicate = onDuplicate
}

// OnUndefinedEvent invokes when event with type that is not defined in library is caught.
//
// If this callback handler is set then ErrUndefinedEventType error will not be returned from event callback runner.
func (c *callback[Metadata]) OnUndefinedEvent(onUndefinedEvent func(RawEvent, Metadata)) {
	c.onUndefinedEvent = onUndefinedEvent
}

// OnAutomodMessageHold invokes when message is caught by automod for review.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#automodmessagehold.
func (c *callback[Metadata]) OnAutomodMessageHold(onAutomodMessageHold Handler[AutomodMessageHoldEvent, Metadata]) {
	c.onAutomodMessageHold = onAutomodMessageHold
}

// OnAutomodMessageHoldV2 invokes when message is caught by automod for review.
// Only public blocked terms trigger notifications, not private ones.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#automodmessagehold-v2.
func (c *callback[Metadata]) OnAutomodMessageHoldV2(onAutomodMessageHoldV2 Handler[AutomodMessageHoldEventV2, Metadata]) {
	c.onAutomodMessageHoldV2 = onAutomodMessageHoldV2
}

// OnAutomodMessageUpdate invokes when a message in the automod queue had its status changed.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#automodmessageupdate.
func (c *callback[Metadata]) OnAutomodMessageUpdate(onAutomodMessageUpdate Handler[AutomodMessageUpdateEvent, Metadata]) {
	c.onAutomodMessageUpdate = onAutomodMessageUpdate
}

// OnAutomodMessageUpdateV2 invokes when a message in the automod queue had its status changed. Only public blocked terms
// trigger notifications, not private ones.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#automodmessageupdate-v2.
func (c *callback[Metadata]) OnAutomodMessageUpdateV2(onAutomodMessageUpdateV2 Handler[AutomodMessageUpdateEventV2, Metadata]) {
	c.onAutomodMessageUpdateV2 = onAutomodMessageUpdateV2
}

// OnAutomodSettingsUpdate invokes when  a broadcaster’s automod settings are updated.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#automodsettingsupdate.
func (c *callback[Metadata]) OnAutomodSettingsUpdate(onAutomodSettingsUpdate Handler[AutomodSettingsUpdateEvent, Metadata]) {
	c.onAutomodSettingsUpdate = onAutomodSettingsUpdate
}

// OnAutomodTermsUpdate invokes when a broadcaster’s automod terms are updated. Changes to private terms are not sent.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#automodtermsupdate.
func (c *callback[Metadata]) OnAutomodTermsUpdate(onAutomodTermsUpdate Handler[AutomodTermsUpdateEvent, Metadata]) {
	c.onAutomodTermsUpdate = onAutomodTermsUpdate
}

// OnChannelBitsUse invokes when bits are used on a channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelbitsuse.
func (c *callback[Metadata]) OnChannelBitsUse(onChannelBitsUse Handler[ChannelBitsUseEvent, Metadata]) {
	c.onChannelBitsUse = onChannelBitsUse
}

// OnChannelUpdate invokes when a broadcaster updates the category, title, content classification labels, or broadcast
// language for their channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelupdate.
func (c *callback[Metadata]) OnChannelUpdate(onChannelUpdate Handler[ChannelUpdateEvent, Metadata]) {
	c.onChannelUpdate = onChannelUpdate
}

// OnChannelFollow invokes when a specified channel receives a follow.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelfollow
func (c *callback[Metadata]) OnChannelFollow(onChannelFollow Handler[ChannelFollowEvent, Metadata]) {
	c.onChannelFollow = onChannelFollow
}

// OnChannelAdBreakBegin invokes when a user runs a midroll commercial break, either manually or automatically via ads manager.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelad_breakbegin.
func (c *callback[Metadata]) OnChannelAdBreakBegin(onChannelAdBreakBegin Handler[ChannelAdBreakBeginEvent, Metadata]) {
	c.onChannelAdBreakBegin = onChannelAdBreakBegin
}

// OnChannelChatClear invokes when a moderator or bot clears all messages from the chat room.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchatclear.
func (c *callback[Metadata]) OnChannelChatClear(onChannelChatClear Handler[ChannelChatClearEvent, Metadata]) {
	c.onChannelChatClear = onChannelChatClear
}

// OnChannelChatClearUserMessages invokes when a moderator or bot clears all messages for a specific user.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchatclear_user_messages.
func (c *callback[Metadata]) OnChannelChatClearUserMessages(onChannelChatClearUserMessages Handler[ChannelChatClearUserMessagesEvent, Metadata]) {
	c.onChannelChatClearUserMessages = onChannelChatClearUserMessages
}

// OnChannelChatMessage invokes when any user sends a message to a channel’s chat room.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchatmessage.
func (c *callback[Metadata]) OnChannelChatMessage(onChannelChatMessage Handler[ChannelChatMessageEvent, Metadata]) {
	c.onChannelChatMessage = onChannelChatMessage
}

// OnConduitShardDisabled invokes when any shard of conduit becomes disabled.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#conduitsharddisabled.
func (c *callback[Metadata]) OnConduitShardDisabled(onConduitShardDisabled Handler[ConduitShardDisabledEvent, Metadata]) {
	c.onConduitShardDisabled = onConduitShardDisabled
}

// OnChannelBan invokes when a user is banned from a channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelban.
func (c *callback[Metadata]) OnChannelBan(onChannelBan Handler[ChannelBanEvent, Metadata]) {
	c.onChannelBan = onChannelBan
}

// OnChannelUnban sends a notification when a viewer is unbanned from the specified channel.
//
// https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelunban.
func (c *callback[Metadata]) OnChannelUnban(onChannelUnban Handler[ChannelUnbanEvent, Metadata]) {
	c.onChannelUnban = onChannelUnban
}

// OnChannelChatNotification invokes when a user sends a chat notification to a channel’s chat room.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchatnotification.
func (c *callback[Metadata]) OnChannelChatNotification(onChannelChatNotification Handler[ChannelChatNotificationEvent, Metadata]) {
	c.onChannelChatNotification = onChannelChatNotification
}

// OnChannelModeratorAdd invokes when a user is added as a moderator to a channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelmoderatoradd.
func (c *callback[Metadata]) OnChannelModeratorAdd(onChannelModeratorAdd Handler[ChannelModeratorAddEvent, Metadata]) {
	c.onChannelModeratorAdd = onChannelModeratorAdd
}

// OnChannelModeratorRemove invokes when a user is removed as a moderator from a channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelmoderatorremove.
func (c *callback[Metadata]) OnChannelModeratorRemove(onChannelModeratorRemove Handler[ChannelModeratorRemoveEvent, Metadata]) {
	c.onChannelModeratorRemove = onChannelModeratorRemove
}

// OnChannelPollBegin invokes when a broadcaster starts a poll in their channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelpollbegin.
func (c *callback[Metadata]) OnChannelPollBegin(onChannelPollBegin Handler[ChannelPollBeginEvent, Metadata]) {
	c.onChannelPollBegin = onChannelPollBegin
}

// OnChannelPollProgress invokes when a broadcaster updates a poll in their channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelpollprogress.
func (c *callback[Metadata]) OnChannelPollProgress(onChannelPollProgress Handler[ChannelPollProgressEvent, Metadata]) {
	c.onChannelPollProgress = onChannelPollProgress
}

// OnChannelPollEnd invokes when a broadcaster ends a poll in their channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelpollend.
func (c *callback[Metadata]) OnChannelPollEnd(onChannelPollEnd Handler[ChannelPollEndEvent, Metadata]) {
	c.onChannelPollEnd = onChannelPollEnd
}

// OnChannelPredictionBegin invokes when a broadcaster starts a prediction in their channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelpredictionbegin.
func (c *callback[Metadata]) OnChannelPredictionBegin(onChannelPredictionBegin Handler[ChannelPredictionBeginEvent, Metadata]) {
	c.onChannelPredictionBegin = onChannelPredictionBegin
}

// OnChannelPredictionProgress invokes when a broadcaster updates a prediction in their channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelpredictionprogress.
func (c *callback[Metadata]) OnChannelPredictionProgress(onChannelPredictionProgress Handler[ChannelPredictionProgressEvent, Metadata]) {
	c.onChannelPredictionProgress = onChannelPredictionProgress
}

// OnChannelPredictionLock invokes when a broadcaster locks a prediction in their channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelpredictionlock.
func (c *callback[Metadata]) OnChannelPredictionLock(onChannelPredictionLock Handler[ChannelPredictionLockEvent, Metadata]) {
	c.onChannelPredictionLock = onChannelPredictionLock
}

// OnChannelPredictionEnd invokes when a broadcaster ends a prediction in their channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelpredictionend.
func (c *callback[Metadata]) OnChannelPredictionEnd(onChannelPredictionEnd Handler[ChannelPredictionEndEvent, Metadata]) {
	c.onChannelPredictionEnd = onChannelPredictionEnd
}

// OnChannelRaid invokes when a channel raids another channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelraid.
func (c *callback[Metadata]) OnChannelRaid(onChannelRaid Handler[ChannelRaidEvent, Metadata]) {
	c.onChannelRaid = onChannelRaid
}

// OnChannelPointsCustomRewardRedemptionAdd invokes when a user redeems a custom channel points reward.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchannel_points_custom_reward_redemptionadd.
func (c *callback[Metadata]) OnChannelPointsCustomRewardRedemptionAdd(onChannelPointsCustomRewardRedemptionAdd Handler[ChannelPointsCustomRewardRedemptionAddEvent, Metadata]) {
	c.onChannelPointsCustomRewardRedemptionAdd = onChannelPointsCustomRewardRedemptionAdd
}

// OnChannelPointsCustomRewardRedemptionUpdate invokes when a user updates a custom channel points reward redemption.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchannel_points_custom_reward_redemptionupdate.
func (c *callback[Metadata]) OnChannelPointsCustomRewardRedemptionUpdate(onChannelPointsCustomRewardRedemptionUpdate Handler[ChannelPointsCustomRewardRedemptionUpdateEvent, Metadata]) {
	c.onChannelPointsCustomRewardRedemptionUpdate = onChannelPointsCustomRewardRedemptionUpdate
}

// OnChannelPointsAutomaticRewardRedemptionAdd invokes when a user redeems an automatic channel points reward.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchannel_points_automatic_reward_redemptionadd.
func (c *callback[Metadata]) OnChannelPointsAutomaticRewardRedemptionAdd(onChannelPointsAutomaticRewardRedemptionAdd Handler[ChannelPointsAutomaticRewardRedemptionAddEvent, Metadata]) {
	c.onChannelPointsAutomaticRewardRedemptionAdd = onChannelPointsAutomaticRewardRedemptionAdd
}

// OnChannelPointsAutomaticRewardRedemptionAddV2 invokes when a user redeems an automatic channel points reward.
// Only public rewards trigger notifications, not private ones.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchannel_points_automatic_reward_redemptionadd-v2.
func (c *callback[Metadata]) OnChannelPointsAutomaticRewardRedemptionAddV2(onChannelPointsAutomaticRewardRedemptionAddV2 Handler[ChannelPointsAutomaticRewardRedemptionAddEventV2, Metadata]) {
	c.onChannelPointsAutomaticRewardRedemptionAddV2 = onChannelPointsAutomaticRewardRedemptionAddV2
}

// OnChannelPointsCustomRewardAdd invokes when a broadcaster adds a custom channel points reward.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchannel_points_rewardadd.
func (c *callback[Metadata]) OnChannelPointsCustomRewardAdd(onChannelPointsCustomRewardAdd Handler[ChannelPointsCustomRewardAddEvent, Metadata]) {
	c.onChannelPointsCustomRewardAdd = onChannelPointsCustomRewardAdd
}

// OnChannelPointsCustomRewardUpdate invokes when a broadcaster updates a custom channel points reward.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchannel_points_rewardupdate.
func (c *callback[Metadata]) OnChannelPointsCustomRewardUpdate(onChannelPointsCustomRewardUpdate Handler[ChannelPointsCustomRewardUpdateEvent, Metadata]) {
	c.onChannelPointsCustomRewardUpdate = onChannelPointsCustomRewardUpdate
}

// OnChannelPointsCustomRewardRemove invokes when a broadcaster removes a custom channel points reward.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchannel_points_rewardremove.
func (c *callback[Metadata]) OnChannelPointsCustomRewardRemove(onChannelPointsCustomRewardRemove Handler[ChannelPointsCustomRewardRemoveEvent, Metadata]) {
	c.onChannelPointsCustomRewardRemove = onChannelPointsCustomRewardRemove
}

// OnStreamOffline invokes when a channel goes offline.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#streamoffline.
func (c *callback[Metadata]) OnStreamOffline(onStreamOffline Handler[StreamOfflineEvent, Metadata]) {
	c.onStreamOffline = onStreamOffline
}

// OnStreamOnline invokes when a channel goes online.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#streamonline.
func (c *callback[Metadata]) OnStreamOnline(onStreamOnline Handler[StreamOnlineEvent, Metadata]) {
	c.onStreamOnline = onStreamOnline
}

// OnChannelSubscribe invokes when a user subscribes to a channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelsubscribe.
func (c *callback[Metadata]) OnChannelSubscribe(onChannelSubscribe Handler[ChannelSubscribeEvent, Metadata]) {
	c.onChannelSubscribe = onChannelSubscribe
}

// OnChannelSubscriptionEnd invokes when a user’s subscription to a channel ends.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelsubscriptionend.
func (c *callback[Metadata]) OnChannelSubscriptionEnd(onChannelSubscriptionEnd Handler[ChannelSubscriptionEndEvent, Metadata]) {
	c.onChannelSubscriptionEnd = onChannelSubscriptionEnd
}

// OnChannelSubscriptionMessage invokes when a user sends a subscription message to a channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelsubscriptionmessage.
func (c *callback[Metadata]) OnChannelSubscriptionMessage(onChannelSubscriptionMessage Handler[ChannelSubscriptionMessageEvent, Metadata]) {
	c.onChannelSubscriptionMessage = onChannelSubscriptionMessage
}

// OnChannelSubscriptionGift invokes when a user gifts a subscription to a channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelsubscriptiongift.
func (c *callback[Metadata]) OnChannelSubscriptionGift(onChannelSubscriptionGift Handler[ChannelSubscriptionGiftEvent, Metadata]) {
	c.onChannelSubscriptionGift = onChannelSubscriptionGift
}

// OnChannelUnbanRequestCreate invokes when a user creates an unban request for a channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelunbanrequestcreate.
func (c *callback[Metadata]) OnChannelUnbanRequestCreate(onChannelUnbanRequestCreate Handler[ChannelUnbanRequestCreateEvent, Metadata]) {
	c.onChannelUnbanRequestCreate = onChannelUnbanRequestCreate
}

// OnChannelUnbanRequestResolve invokes when a user resolves an unban request for a channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelunbanrequestresolve.
func (c *callback[Metadata]) OnChannelUnbanRequestResolve(onChannelUnbanRequestResolve Handler[ChannelUnbanRequestResolveEvent, Metadata]) {
	c.onChannelUnbanRequestResolve = onChannelUnbanRequestResolve
}

// OnUserUpdate invokes when a user updates their profile information.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#userupdate.
func (c *callback[Metadata]) OnUserUpdate(onUserUpdate Handler[UserUpdateEvent, Metadata]) {
	c.onUserUpdate = onUserUpdate
}

// OnChannelVipAdd invokes when a user is added as a VIP to a channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelvipadd.
func (c *callback[Metadata]) OnChannelVipAdd(onChannelVipAdd Handler[ChannelVipAddEvent, Metadata]) {
	c.onChannelVipAdd = onChannelVipAdd
}

// OnChannelVipRemove invokes when a user is removed as a VIP from a channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelvipremove.
func (c *callback[Metadata]) OnChannelVipRemove(onChannelVipRemove Handler[ChannelVipRemoveEvent, Metadata]) {
	c.onChannelVipRemove = onChannelVipRemove
}

// OnChannelChatMessageDelete invokes when a user deletes a message in a channel's chat room.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchatmessage_delete.
func (c *callback[Metadata]) OnChannelChatMessageDelete(onChannelChatMessageDelete Handler[ChannelChatMessageDeleteEvent, Metadata]) {
	c.onChannelChatMessageDelete = onChannelChatMessageDelete
}

// OnUserAuthorizationRevoke invokes when a user revokes authorization for an application.
//
// Note: This subscription type is only supported by webhooks, and cannot be used with websockets.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#userauthorizationrevoke.
func (c *callback[Metadata]) OnUserAuthorizationRevoke(onUserAuthorizationRevoke Handler[UserAuthorizationRevokeEvent, Metadata]) {
	c.onUserAuthorizationRevoke = onUserAuthorizationRevoke
}

// OnChannelModerateV2 subscription type sends a notification when a moderator performs a moderation action in a channel. Some of these actions affect chatters in other channels during Shared Chat.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelmoderate-v2.
func (c *callback[Metadata]) OnChannelModerateV2(onChannelModerateV2 Handler[ChannelModerateEventV2, Metadata]) {
	c.onChannelModerateV2 = onChannelModerateV2
}

// OnChannelChatSettingsUpdate invokes when a broadcaster's chat settings are updated.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchat_settingsupdate.
func (c *callback[Metadata]) OnChannelChatSettingsUpdate(onChannelChatSettingsUpdate Handler[ChannelChatSettingsUpdateEvent, Metadata]) {
	c.onChannelChatSettingsUpdate = onChannelChatSettingsUpdate
}

// OnChannelChatUserMessageHold invokes when a user is notified if their message is caught by automod.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchatuser_message_hold.
func (c *callback[Metadata]) OnChannelChatUserMessageHold(onChannelChatUserMessageHold Handler[ChannelChatUserMessageHoldEvent, Metadata]) {
	c.onChannelChatUserMessageHold = onChannelChatUserMessageHold
}

// OnChannelChatUserMessageUpdate invokes when a user is notified if their message's automod status is updated.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelchatuser_message_update.
func (c *callback[Metadata]) OnChannelChatUserMessageUpdate(onChannelChatUserMessageUpdate Handler[ChannelChatUserMessageUpdateEvent, Metadata]) {
	c.onChannelChatUserMessageUpdate = onChannelChatUserMessageUpdate
}

// OnChannelCustomPowerUpRedemptionAdd invokes when a viewer has redeemed a custom Power-up on the specified channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelcustom_power_up_redemptionadd.
func (c *callback[Metadata]) OnChannelCustomPowerUpRedemptionAdd(onChannelCustomPowerUpRedemptionAdd Handler[ChannelCustomPowerUpRedemptionAddEvent, Metadata]) {
	c.onChannelCustomPowerUpRedemptionAdd = onChannelCustomPowerUpRedemptionAdd
}

// OnChannelSharedChatSessionBegin invokes when a channel becomes active in an active shared chat session.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelshared_chatbegin.
func (c *callback[Metadata]) OnChannelSharedChatSessionBegin(onChannelSharedChatSessionBegin Handler[ChannelSharedChatSessionBeginEvent, Metadata]) {
	c.onChannelSharedChatSessionBegin = onChannelSharedChatSessionBegin
}

// OnChannelSharedChatSessionUpdate invokes when the active shared chat session the channel is in changes.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelshared_chatupdate.
func (c *callback[Metadata]) OnChannelSharedChatSessionUpdate(onChannelSharedChatSessionUpdate Handler[ChannelSharedChatSessionUpdateEvent, Metadata]) {
	c.onChannelSharedChatSessionUpdate = onChannelSharedChatSessionUpdate
}

// OnChannelSharedChatSessionEnd invokes when a channel leaves a shared chat session or the session ends.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelshared_chatend.
func (c *callback[Metadata]) OnChannelSharedChatSessionEnd(onChannelSharedChatSessionEnd Handler[ChannelSharedChatSessionEndEvent, Metadata]) {
	c.onChannelSharedChatSessionEnd = onChannelSharedChatSessionEnd
}

// OnChannelSuspiciousUserMessage invokes when a chat message has been sent by a suspicious user.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelsuspicious_usermessage.
func (c *callback[Metadata]) OnChannelSuspiciousUserMessage(onChannelSuspiciousUserMessage Handler[ChannelSuspiciousUserMessageEvent, Metadata]) {
	c.onChannelSuspiciousUserMessage = onChannelSuspiciousUserMessage
}

// OnChannelSuspiciousUserUpdate invokes when a suspicious user has been updated.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelsuspicious_userupdate.
func (c *callback[Metadata]) OnChannelSuspiciousUserUpdate(onChannelSuspiciousUserUpdate Handler[ChannelSuspiciousUserUpdateEvent, Metadata]) {
	c.onChannelSuspiciousUserUpdate = onChannelSuspiciousUserUpdate
}

// OnChannelWarningAcknowledge invokes when a user acknowledges a warning.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelwarningacknowledge.
func (c *callback[Metadata]) OnChannelWarningAcknowledge(onChannelWarningAcknowledge Handler[ChannelWarningAcknowledgeEvent, Metadata]) {
	c.onChannelWarningAcknowledge = onChannelWarningAcknowledge
}

// OnChannelWarningSend invokes when a user is sent a warning.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelwarningsend.
func (c *callback[Metadata]) OnChannelWarningSend(onChannelWarningSend Handler[ChannelWarningSendEvent, Metadata]) {
	c.onChannelWarningSend = onChannelWarningSend
}

// OnChannelCheer invokes when a user cheers on the specified channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelcheer.
func (c *callback[Metadata]) OnChannelCheer(onChannelCheer Handler[ChannelCheerEvent, Metadata]) {
	c.onChannelCheer = onChannelCheer
}

// OnCharityCampaignDonate invokes when a user donates to the broadcaster's charity campaign.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelcharity_campaigndonate.
func (c *callback[Metadata]) OnCharityCampaignDonate(onCharityCampaignDonate Handler[CharityCampaignDonateEvent, Metadata]) {
	c.onCharityCampaignDonate = onCharityCampaignDonate
}

// OnCharityCampaignStart invokes when the broadcaster starts a charity campaign.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelcharity_campaignstart.
func (c *callback[Metadata]) OnCharityCampaignStart(onCharityCampaignStart Handler[CharityCampaignStartEvent, Metadata]) {
	c.onCharityCampaignStart = onCharityCampaignStart
}

// OnCharityCampaignProgress invokes when progress is made towards the campaign's goal or when the broadcaster changes the fundraising goal.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelcharity_campaignprogress.
func (c *callback[Metadata]) OnCharityCampaignProgress(onCharityCampaignProgress Handler[CharityCampaignProgressEvent, Metadata]) {
	c.onCharityCampaignProgress = onCharityCampaignProgress
}

// OnCharityCampaignStop invokes when the broadcaster stops a charity campaign.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelcharity_campaignstop.
func (c *callback[Metadata]) OnCharityCampaignStop(onCharityCampaignStop Handler[CharityCampaignStopEvent, Metadata]) {
	c.onCharityCampaignStop = onCharityCampaignStop
}

// OnDropEntitlementGrant invokes when an entitlement for a Drop is granted to a user.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#dropentitlementgrant.
func (c *callback[Metadata]) OnDropEntitlementGrant(onDropEntitlementGrant Handler[DropEntitlementGrantEvent, Metadata]) {
	c.onDropEntitlementGrant = onDropEntitlementGrant
}

// OnExtensionBitsTransactionCreate invokes when a Bits transaction occurred for a specified Twitch Extension.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#extensionbits_transactioncreate.
func (c *callback[Metadata]) OnExtensionBitsTransactionCreate(onExtensionBitsTransactionCreate Handler[ExtensionBitsTransactionCreateEvent, Metadata]) {
	c.onExtensionBitsTransactionCreate = onExtensionBitsTransactionCreate
}

// OnChannelGoalBegin invokes when a broadcaster begins a goal.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelgoalbegin.
func (c *callback[Metadata]) OnChannelGoalBegin(onChannelGoalBegin Handler[ChannelGoalBeginEvent, Metadata]) {
	c.onChannelGoalBegin = onChannelGoalBegin
}

// OnChannelGoalProgress invokes when progress is made towards a broadcaster's goal.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelgoalprogress.
func (c *callback[Metadata]) OnChannelGoalProgress(onChannelGoalProgress Handler[ChannelGoalProgressEvent, Metadata]) {
	c.onChannelGoalProgress = onChannelGoalProgress
}

// OnChannelGoalEnd invokes when a broadcaster ends a goal.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelgoalend.
func (c *callback[Metadata]) OnChannelGoalEnd(onChannelGoalEnd Handler[ChannelGoalEndEvent, Metadata]) {
	c.onChannelGoalEnd = onChannelGoalEnd
}

// OnHypeTrainBegin invokes when a Hype Train begins on the specified channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelhype_trainbegin.
func (c *callback[Metadata]) OnHypeTrainBegin(onHypeTrainBegin Handler[HypeTrainBeginEvent, Metadata]) {
	c.onHypeTrainBegin = onHypeTrainBegin
}

// OnHypeTrainProgress invokes when a Hype Train makes progress on the specified channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelhype_trainprogress.
func (c *callback[Metadata]) OnHypeTrainProgress(onHypeTrainProgress Handler[HypeTrainProgressEvent, Metadata]) {
	c.onHypeTrainProgress = onHypeTrainProgress
}

// OnHypeTrainEnd invokes when a Hype Train ends on the specified channel.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelhype_trainend.
func (c *callback[Metadata]) OnHypeTrainEnd(onHypeTrainEnd Handler[HypeTrainEndEvent, Metadata]) {
	c.onHypeTrainEnd = onHypeTrainEnd
}

// OnShieldModeBegin invokes when the broadcaster activates Shield Mode.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelshield_modebegin.
func (c *callback[Metadata]) OnShieldModeBegin(onShieldModeBegin Handler[ShieldModeBeginEvent, Metadata]) {
	c.onShieldModeBegin = onShieldModeBegin
}

// OnShieldModeEnd invokes when the broadcaster deactivates Shield Mode.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelshield_modeend.
func (c *callback[Metadata]) OnShieldModeEnd(onShieldModeEnd Handler[ShieldModeEndEvent, Metadata]) {
	c.onShieldModeEnd = onShieldModeEnd
}

// OnShoutoutCreate invokes when the specified broadcaster sends a Shoutout.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelshoutoutcreate.
func (c *callback[Metadata]) OnShoutoutCreate(onShoutoutCreate Handler[ShoutoutCreateEvent, Metadata]) {
	c.onShoutoutCreate = onShoutoutCreate
}

// OnShoutoutReceive invokes when the specified broadcaster receives a Shoutout.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelshoutoutreceive.
func (c *callback[Metadata]) OnShoutoutReceive(onShoutoutReceive Handler[ShoutoutReceiveEvent, Metadata]) {
	c.onShoutoutReceive = onShoutoutReceive
}

// OnUserAuthorizationGrant invokes when a user's authorization has been granted to your client id.
//
// Note: This subscription type is only supported by webhooks, and cannot be used with websockets.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#userauthorizationgrant.
func (c *callback[Metadata]) OnUserAuthorizationGrant(onUserAuthorizationGrant Handler[UserAuthorizationGrantEvent, Metadata]) {
	c.onUserAuthorizationGrant = onUserAuthorizationGrant
}

// OnUserWhisperMessage invokes when a user receives a whisper.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#userwhispermessage.
func (c *callback[Metadata]) OnUserWhisperMessage(onUserWhisperMessage Handler[WhisperReceivedEvent, Metadata]) {
	c.onUserWhisperMessage = onUserWhisperMessage
}

// OnGuestStarSessionBegin invokes when the host began a new Guest Star session.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelguest_star_sessionbegin.
func (c *callback[Metadata]) OnGuestStarSessionBegin(onGuestStarSessionBegin Handler[ChannelGuestStarSessionBeginEvent, Metadata]) {
	c.onGuestStarSessionBegin = onGuestStarSessionBegin
}

// OnGuestStarSessionEnd invokes when a running Guest Star session has ended.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelguest_star_sessionend.
func (c *callback[Metadata]) OnGuestStarSessionEnd(onGuestStarSessionEnd Handler[ChannelGuestStarSessionEndEvent, Metadata]) {
	c.onGuestStarSessionEnd = onGuestStarSessionEnd
}

// OnGuestStarGuestUpdate invokes when a guest or a slot is updated in an active Guest Star session.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelguest_star_guestupdate.
func (c *callback[Metadata]) OnGuestStarGuestUpdate(onGuestStarGuestUpdate Handler[ChannelGuestStarGuestUpdateEvent, Metadata]) {
	c.onGuestStarGuestUpdate = onGuestStarGuestUpdate
}

// OnGuestStarSettingsUpdate invokes when the host preferences for Guest Star have been updated.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelguest_star_settingsupdate.
func (c *callback[Metadata]) OnGuestStarSettingsUpdate(onGuestStarSettingsUpdate Handler[ChannelGuestStarSettingsUpdateEvent, Metadata]) {
	c.onGuestStarSettingsUpdate = onGuestStarSettingsUpdate
}
