package eventsub

// EventType is a type of EventSub event.
//
// Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types
type EventType string

const (
	EventTypeAutomodMessageUpdate                      EventType = "automod.message.update"
	EventTypeAutomodMessageHold                        EventType = "automod.message.hold"
	EventTypeAutomodSettingsUpdate                     EventType = "automod.settings.update"
	EventTypeAutomodTermsUpdate                        EventType = "automod.terms.update"
	EventTypeChannelBitsUse                            EventType = "channel.bits.use"
	EventTypeChannelUpdate                             EventType = "channel.update"
	EventTypeChannelFollow                             EventType = "channel.follow"
	EventTypeChannelAdBreakBegin                       EventType = "channel.ad_break.begin"
	EventTypeChannelChatClear                          EventType = "channel.chat.clear"
	EventTypeChannelChatClearUserMessages              EventType = "channel.chat.clear_user_messages"
	EventTypeChannelChatMessage                        EventType = "channel.chat.message"
	EventTypeConduitShardDisabled                      EventType = "conduit.shard.disabled"
	EventTypeChannelBan                                EventType = "channel.ban"
	EventTypeChannelChatNotification                   EventType = "channel.chat.notification"
	EventTypeChannelUnban                              EventType = "channel.unban"
	EventTypeChannelModeratorAdd                       EventType = "channel.moderator.add"
	EventTypeChannelModeratorRemove                    EventType = "channel.moderator.remove"
	EventTypeChannelPollBegin                          EventType = "channel.poll.begin"
	EventTypeChannelPollProgress                       EventType = "channel.poll.progress"
	EventTypeChannelPollEnd                            EventType = "channel.poll.end"
	EventTypeChannelPredictionBegin                    EventType = "channel.prediction.begin"
	EventTypeChannelPredictionProgress                 EventType = "channel.prediction.progress"
	EventTypeChannelPredictionLock                     EventType = "channel.prediction.lock"
	EventTypeChannelPredictionEnd                      EventType = "channel.prediction.end"
	EventTypeChannelRaid                               EventType = "channel.raid"
	EventTypeChannelPointsCustomRewardRedemptionAdd    EventType = "channel.channel_points_custom_reward_redemption.add"
	EventTypeChannelPointsCustomRewardRedemptionUpdate EventType = "channel.channel_points_custom_reward_redemption.update"
	EventTypeChannelPointsAutomaticRewardRedemptionAdd EventType = "channel.channel_points_automatic_reward_redemption.add"
	EventTypeUserAuthorizationRevoke                   EventType = "user.authorization.revoke"
	EventTypeChannelPointsRewardAdd                    EventType = "channel.channel_points_custom_reward.add"
	EventTypeChannelPointsRewardUpdate                 EventType = "channel.channel_points_custom_reward.update"
	EventTypeChannelPointsRewardRemove                 EventType = "channel.channel_points_custom_reward.remove"
	EventTypeStreamOffline                             EventType = "stream.offline"
	EventTypeStreamOnline                              EventType = "stream.online"
	EventTypeChannelSubscribe                          EventType = "channel.subscribe"
	EventTypeChannelSubscriptionEnd                    EventType = "channel.subscription.end"
	EventTypeChannelSubscriptionMessage                EventType = "channel.subscription.message"
	EventTypeChannelSubscriptionGift                   EventType = "channel.subscription.gift"
	EventTypeChannelUnbanRequestCreate                 EventType = "channel.unban_request.create"
	EventTypeChannelUnbanRequestResolve                EventType = "channel.unban_request.resolve"
	EventTypeUserUpdate                                EventType = "user.update"
	EventTypeChannelVipAdd                             EventType = "channel.vip.add"
	EventTypeChannelVipRemove                          EventType = "channel.vip.remove"
	EventTypeChannelMessageDelete                      EventType = "channel.chat.message_delete"
	EventTypeChannelModerate                           EventType = "channel.moderate"
)

func (et EventType) String() string {
	return string(et)
}
