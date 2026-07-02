package eventsub

type AutomodMessageHoldEvent struct {
	// The ID of the broadcaster specified in the request.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the broadcaster specified in the request.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The username of the broadcaster specified in the request.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The message sender’s user ID.
	UserId string `json:"user_id"`
	// The message sender’s login name.
	UserLogin string `json:"user_login"`
	// The message sender’s display name.
	UserName string `json:"user_name"`
	// The ID of the message that was flagged by automod.
	MessageId string `json:"message_id"`
	// The body of the message.
	Message AutomodMessageHoldEventMessage `json:"message"`
	// The category of the message.
	Category string `json:"category"`
	// The level of severity.
	Level AutomodSeverityLevel `json:"level"`
	// The timestamp of when automod saved the message.
	HeldAt Timestamp `json:"held_at"`
}

type AutomodMessageHoldEventV2 struct {
	// The ID of the broadcaster specified in the request.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the broadcaster specified in the request.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The username of the broadcaster specified in the request.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The message sender’s user ID.
	UserId string `json:"user_id"`
	// The message sender’s login name.
	UserLogin string `json:"user_login"`
	// The message sender’s display name.
	UserName string `json:"user_name"`
	// The ID of the message that was flagged by automod.
	MessageId string `json:"message_id"`
	// The body of the message.
	Message AutomodMessageHoldEventMessageV2 `json:"message"`
	// The timestamp of when automod saved the message.
	HeldAt Timestamp `json:"held_at"`
	// The reason why automod hold this message.
	Reason AutomodHoldReason `json:"reason"`
	// Optional. If the message was caught by automod, this will be populated.
	Automod *Automod `json:"automod,omitempty"`
	// Optional. If the message was caught due to a blocked term, this will be populated.
	BlockedTerm *BlockedTerm `json:"blocked_term,omitempty"`
}

type AutomodMessageUpdateEvent struct {
	// The ID of the broadcaster specified in the request.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the broadcaster specified in the request.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The username of the broadcaster specified in the request.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The message sender’s user ID.
	UserId string `json:"user_id"`
	// The message sender’s login name.
	UserLogin string `json:"user_login"`
	// The message sender’s display name.
	UserName string `json:"user_name"`
	// The ID of the moderator.
	ModeratorUserId string `json:"moderator_user_id"`
	// The moderator’s username.
	ModeratorUserName string `json:"moderator_user_name"`
	// The login of the moderator.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The ID of the message that was flagged by automod.
	MessageId string `json:"message_id"`
	// The body of the message.
	Message AutomodMessageUpdateEventMessage `json:"message"`
	// The category of the message.
	Category string `json:"category"`
	// The level of severity.
	Level AutomodSeverityLevel `json:"level"`
	// The message’s status.
	Status AutomodMessageStatus `json:"status"`
	// The timestamp of when automod saved the message.
	HeldAt Timestamp `json:"held_at"`
}

type AutomodMessageUpdateEventV2 struct {
	// The ID of the broadcaster specified in the request.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the broadcaster specified in the request.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The username of the broadcaster specified in the request.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The message sender’s user ID.
	UserId string `json:"user_id"`
	// The message sender’s login name.
	UserLogin string `json:"user_login"`
	// The message sender’s display name.
	UserName string `json:"user_name"`
	// The ID of the moderator.
	ModeratorUserId string `json:"moderator_user_id"`
	// The moderator’s username.
	ModeratorUserName string `json:"moderator_user_name"`
	// The login of the moderator.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The ID of the message that was flagged by automod.
	MessageId string `json:"message_id"`
	// The body of the message.
	Message AutomodMessageUpdateEventMessageV2 `json:"message"`
	// The message’s status.
	Status AutomodMessageStatus `json:"status"`
	// The timestamp of when automod saved the message.
	HeldAt Timestamp `json:"held_at"`
	// The reason why automod hold this message.
	Reason AutomodHoldReason `json:"reason"`
	// Optional. If the message was caught by automod, this will be populated.
	Automod *Automod `json:"automod,omitempty"`
	// Optional. If the message was caught due to a blocked term, this will be populated.
	BlockedTerm *BlockedTerm `json:"blocked_term,omitempty"`
}

type AutomodSettingsUpdateEvent struct {
	// The ID of the broadcaster specified in the request.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the broadcaster specified in the request.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The username of the broadcaster specified in the request.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The ID of the moderator.
	ModeratorUserId string `json:"moderator_user_id"`
	// The moderator’s username.
	ModeratorUserName string `json:"moderator_user_name"`
	// The login of the moderator.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The Automod level for hostility involving name-calling or insults.
	Bullying int `json:"bullying"`
	// The default AutoMod level for the broadcaster.
	// This field is empty if the broadcaster has set one or more of the individual settings.
	OverallLevel int `json:"overall_level,omitempty"`
	// The Automod level for discrimination against disability.
	Disability int `json:"disability"`
	// The Automod level for racial discrimination.
	RaceEthnicityOrReligion int `json:"race_ethnicity_or_religion"`
	// Misogyny is an Automod level for discrimination against women.
	Misogyny int `json:"misogyny"`
	// The AutoMod level for discrimination based on sexuality, sex, or gender.
	SexualitySexOrGender int `json:"sexuality_sex_or_gender"`
	// The Automod level for hostility involving aggression.
	Aggression int `json:"aggression"`
	// The Automod level for sexual content.
	SexBasedTerms int `json:"sex_based_terms"`
	// The Automod level for profanity.
	Swearing int `json:"swearing"`
}

type AutomodTermsUpdateEvent struct {
	// The ID of the broadcaster specified in the request.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the broadcaster specified in the request.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The username of the broadcaster specified in the request.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The ID of the moderator.
	ModeratorUserId string `json:"moderator_user_id"`
	// The moderator’s username.
	ModeratorUserName string `json:"moderator_user_name"`
	// The login of the moderator.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The status change applied to the terms.
	Action AutomodTermsAction `json:"action"`
	// Indicates whether this term was added due to an Automod message approve/deny action.
	FromAutomod bool `json:"from_automod"`
	// The list of terms that had a status change.
	Terms []string `json:"terms"`
}

type ChannelBitsUseEvent struct {
	// The user ID of the channel where the bits were redeemed.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the channel where the bits were used.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The display name of the channel where the bits were used.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user ID of the redeeming user.
	UserId string `json:"user_id"`
	// The login name of the redeeming user.
	UserLogin string `json:"user_login"`
	// The display name of the redeeming user.
	UserName string `json:"user_name"`
	// The number if bits used.
	Bits int `json:"bits"`
	// The type of bits usage.
	Type BitsUseType `json:"type"`
	// Optional. An object that contains the user message and emote information needed to recreate the message.
	Message *ChannelBitsUseEventMessage `json:"message,omitempty"`
	// Optional. Data about Power-up.
	PowerUp *PowerUp `json:"power_up,omitempty"`
}

type ChannelUpdateEvent struct {
	// The broadcaster’s user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster’s user login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s user display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The channel’s stream title.
	Title string `json:"title"`
	// The channel’s broadcast language.
	Language string `json:"language"`
	// The channel’s category ID.
	CategoryId string `json:"category_id"`
	// The category name.
	CategoryName string `json:"category_name"`
	// Content classification label IDs currently applied on the channel.
	// To retrieve a list of all possible IDs, use the Get Content Classification Labels API endpoint.
	ContentClassificationLabels []string `json:"content_classification_labels"`
}

type ChannelFollowEvent struct {
	// The user ID for the user now following the specified channel.
	UserId string `json:"user_id"`
	// The user login for the user now following the specified channel.
	UserLogin string `json:"user_login"`
	// The user display name for the user now following the specified channel.
	UserName string `json:"user_name"`
	// The broadcaster’s user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster’s user login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s user display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Timestamp of when the follow occurred.
	FollowedAt Timestamp `json:"followed_at"`
}

type ChannelAdBreakBeginEvent struct {
	// Length in seconds of the mid-roll ad break requested.
	DurationSeconds int `json:"duration_seconds"`
	// The UTC timestamp of when the ad break began, in RFC3339 format.
	// Note that there is potential delay between this event, when the streamer requested the ad break, and when the viewers will see ads.
	StartedAt TimestampUTC `json:"started_at"`
	// Indicates if the ad was automatically scheduled via Ads Manager.
	IsAutomatic bool `json:"is_automatic"`
	// The broadcaster’s user ID for the channel the ad was run on.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster’s user login for the channel the ad was run on.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s user display name for the channel the ad was run on.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The ID of the user that requested the ad. For automatic ads, this will be the ID of the broadcaster.
	RequesterUserId string `json:"requester_user_id"`
	// The login of the user that requested the ad.
	RequesterUserLogin string `json:"requester_user_login"`
	// The display name of the user that requested the ad.
	RequesterUserName string `json:"requester_user_name"`
}

type ChannelChatClearEvent struct {
	// The broadcaster’s user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster’s user login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s user display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
}

type ChannelChatClearUserMessagesEvent struct {
	// The broadcaster’s user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster’s user login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s user display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The ID of the user that was banned or put in a timeout. All of their messages are deleted.
	TargetUserId string `json:"target_user_id"`
	// The user login of the user that was banned or put in a timeout.
	TargetUserLogin string `json:"target_user_login"`
	// The username of the user that was banned or put in a timeout.
	TargetUserName string `json:"target_user_name"`
}

type ChannelChatMessageEvent struct {
	// The broadcaster’s user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster’s user login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s user display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user ID of the user that sent the message.
	ChatterUserId string `json:"chatter_user_id"`
	// The username of the user that sent the message.
	ChatterUserName string `json:"chatter_user_login"`
	// The user login of the user that sent the message.
	ChatterUserLogin string `json:"chatter_user_name"`
	// A UUID that identifies the message.
	MessageId string `json:"message_id"`
	// The structured chat message.
	Message ChannelChatMessageEventMessage `json:"message"`
	// The type of message.
	Type MessageType `json:"message_type"`
	// List of chat badges.
	Badges []Badge `json:"badges"`
	// Optional. Metadata if this message is a cheer.
	Cheer *Cheer `json:"cheer,omitempty"`
	// The color of the user’s name in the chat room.
	// This is a hexadecimal RGB color code in the form, #&lt;RGB&gt;.
	// This tag may be empty if it is never set.
	Color string `json:"color"`
	// Optional. Metadata if this message is a reply.
	Reply *Reply `json:"reply,omitempty"`
	// ChannelPointsCustomRewardId is an ID of a channel points custom reward that was redeemed.
	ChannelPointsCustomRewardId string `json:"channel_points_custom_reward_id,omitempty"`
	// Optional. The broadcaster user ID of the channel the message was sent from.
	// Not presented when the message happens in the same channel as the broadcaster.
	// Presented when in a shared chat session, and the action happens in the channel of a participant other than the broadcaster.
	SourceBroadcasterUserId string `json:"source_broadcaster_user_id,omitempty"`
	// Optional. The username of the broadcaster of the channel the message was sent from.
	// Not presented when the message happens in the same channel as the broadcaster.
	// Presented when in a shared chat session, and the action happens in the channel of a participant other than the broadcaster.
	SourceBroadcasterUserName string `json:"source_broadcaster_user_name,omitempty"`
	// Optional. The login of the broadcaster of the channel the message was sent from.
	// Not presented when the message happens in the same channel as the broadcaster.
	// Presented when in a shared chat session, and the action happens in the channel of a participant other than the broadcaster.
	SourceBroadcasterUserLogin string `json:"source_broadcaster_user_login,omitempty"`
	// Optional. The UUID that identifies the source message from the channel the message was sent from.
	// Not presented when the message happens in the same channel as the broadcaster.
	// Presented when in a shared chat session, and the action happens in the channel of a participant other than the broadcaster.
	SourceMessageId string `json:"source_message_id,omitempty"`
	// Optional. The list of chat badges for the chatter in the channel the message was sent from.
	// Not presented when the message happens in the same channel as the broadcaster.
	// Presented when in a shared chat session, and the action happens in the channel of a participant other than the broadcaster.
	SourceBadges []Badge `json:"source_badges,omitempty"`
	// Optional. Determines if a message delivered during a shared chat session is only sent to the source channel.
	// Has no effect if the message is not sent during a shared chat session.
	IsSourceOnly *bool `json:"is_source_only,omitempty"`
}

type ConduitShardDisabledEvent struct {
	// The ID of the conduit.
	ConduitId string `json:"conduit_id"`
	// The ID of the disabled shard.
	ShardId string `json:"shard_id"`
	// The new status of the transport.
	Status string `json:"status"`
	// The disabled transport.
	Transport ConduitShardDisabledEventTransport `json:"transport"`
}

type ChannelBanEvent struct {
	// The user ID for the user who was banned on the specified channel.
	UserId string `json:"user_id"`
	// The user login for the user who was banned on the specified channel.
	UserLogin string `json:"user_login"`
	// The user display name for the user who was banned on the specified channel.
	UserName string `json:"user_name"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user ID of the issuer of the ban.
	ModeratorUserID string `json:"moderator_user_id"`
	// The user login of the issuer of the ban.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The username of the issuer of the ban.
	ModeratorUserName string `json:"moderator_user_name"`
	// The reason behind the ban.
	Reason string `json:"reason"`
	// Will be not presented if permanent ban. If it is a timeout, this field shows when the timeout will end.
	EndsAt *TimestampUTC `json:"ends_at,omitempty"`
	// Indicates whether the ban is permanent (true) or a timeout (false). If true, ends_at will be not presented.
	IsPermanent bool `json:"is_permanent"`
}

type ChannelUnbanEvent struct {
	// The user id for the user who was unbanned on the specified channel.
	UserID string `json:"user_id"`
	// The user login for the user who was unbanned on the specified channel.
	UserLogin string `json:"user_login"`
	// The user display name for the user who was unbanned on the specified channel.
	UserName string `json:"user_name"`
	// The requested broadcaster ID.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user ID of the issuer of the unban.
	ModeratorUserID string `json:"moderator_user_id"`
	// The user login of the issuer of the unban.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The user name of the issuer of the unban.
	ModeratorUserName string `json:"moderator_user_name"`
}

type ChannelChatNotificationEvent struct {
	// The broadcaster user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// Optional. The broadcaster user ID of the channel the message was sent from.
	// Not presented when the message notification happens in the same channel as the broadcaster.
	// Presented when in a shared chat session, and the action happens in the channel of a participant other than the broadcaster.
	SourceBroadcasterUserId string `json:"source_broadcaster_user_id,omitempty"`
	// Optional. The username of the broadcaster of the channel the message was sent from.
	// Not presented when the message notification happens in the same channel as the broadcaster.
	// Presented when in a shared chat session, and the action happens in the channel of a participant other than the broadcaster.
	SourceBroadcasterUserName string `json:"source_broadcaster_user_name,omitempty"`
	// Optional. The login of the broadcaster of the channel the message was sent from
	// Not presented when the message notification happens in the same channel as the broadcaster.
	// Presented when in a shared chat session, and the action happens in the channel of a participant other than the broadcaster.
	SourceBroadcasterUserLogin string `json:"source_broadcaster_user_login,omitempty"`
	// Optional. The UUID that identifies the source message from the channel the message was sent from.
	// Not presented when the message happens in the same channel as the broadcaster.
	// Presented when in a shared chat session, and the action happens in the channel of a participant other than the broadcaster.
	SourceMessageId string `json:"source_message_id,omitempty"`
	// Optional. The list of chat badges for the chatter in the channel the message was sent from.
	// Not presented when the message happens in the same channel as the broadcaster.
	// Presented when in a shared chat session, and the action happens in the channel of a participant other than the broadcaster.
	SourceBadges []ChatNotificationEventBadge `json:"source_badges,omitempty"`
	// The user ID of the user that sent the message.
	ChatterUserId string `json:"chatter_user_id"`
	// The username of the user that sent the message.
	ChatterUserName string `json:"chatter_user_name"`
	// The user login of the user that sent the message.
	ChatterUserLogin string `json:"chatter_user_login"`
	// Whether the chatter is anonymous or not.
	ChatterIsAnonymous bool `json:"chatter_is_anonymous"`
	// The color of the user’s name in the chat room.
	Color string `json:"color"`
	// List of chat badges.
	Badges []ChatNotificationEventBadge `json:"badges"`
	// The message Twitch shows in the chat room for this notice.
	SystemMessage string `json:"system_message"`
	// A UUID that identifies the message.
	MessageId string `json:"message_id"`
	// The structured chat message
	Message ChatNotificationEventMessage `json:"message"`
	// The type of notice.
	NoticeType ChatNotificationEventNoticeType `json:"notice_type"`
	// Information about the sub event. Not presented if notice_type is not sub.
	Sub *ChatNotificationEventSubEvent `json:"sub,omitempty"`
	// Information about the re-sub event. Not presented if notice_type is not re-sub.
	ReSub *ChatNotificationEventReSubEvent `json:"resub,omitempty"`
	// Information about the gift sub event. Not presented if notice_type is not sub_gift.
	SubGift *ChatNotificationEventSubGiftEvent `json:"sub_gift,omitempty"`
	// Information about the community gift sub event. Not presented if notice_type is not community_sub_gift.
	CommunitySubGift *ChatNotificationEventCommunitySubGiftEvent `json:"community_sub_gift,omitempty"`
	// Information about the community gift paid upgrade event. Not presented if notice_type is not gift_paid_upgrade.
	GiftPaidUpgrade *ChatNotificationEventGiftPaidUpgradeEvent `json:"gift_paid_upgrade,omitempty"`
	// Information about the Prime gift paid upgrade event. Not presented if notice_type is not prime_paid_upgrade.
	PrimePaidUpgrade *ChatNotificationEventGiftPaidUpgradeEvent `json:"prime_paid_upgrade,omitempty"`
	// Information about the raid event. Not presented if notice_type is not raid.
	Raid *ChatNotificationEventRaidEvent `json:"raid,omitempty"`
	// Returns an empty payload if notice_type is un-raid, otherwise returns null.
	UnRaid *ChatNotificationEventUnRaidEvent `json:"unraid,omitempty"`
	// Information about the pay it forward event. Not presented if notice_type is not pay_it_forward.
	PayItForward *ChatNotificationEventPayItForwardEvent `json:"pay_it_forward,omitempty"`
	// Information about the announcement event. Not presented if notice_type is not announcement.
	Announcement *ChatNotificationEventAnnouncementEvent `json:"announcement,omitempty"`
	// Information about the charity donation event. Not presented if notice_type is not charity_donation.
	CharityDonation *ChatNotificationEventCharityDonationEvent `json:"charity_donation,omitempty"`
	// Information about the bits badge tier event. Not presented if notice_type is not bits_badge_tier.
	BitsBadgeTier *ChatNotificationEventBitsBadgeTierEvent `json:"bits_badge_tier,omitempty"`
	// Information about the shared_chat_sub event.
	// Not presented if notice_type is not shared_chat_sub.
	// This field has the same information as the sub-field but for a notice that happened for a channel in a shared chat
	// session other than the broadcaster in the subscription condition.
	SharedChatSub *ChatNotificationEventSubEvent `json:"shared_chat_sub,omitempty"`
	// Information about the shared_chat_resub event.
	// Not presented if notice_type is not shared_chat_resub.
	// This field has the same information as the re-sub field but for a notice that happened for a channel in a shared chat
	// session other than the broadcaster in the subscription condition.
	SharedChatReSub *ChatNotificationEventReSubEvent `json:"shared_chat_resub,omitempty"`
	// Information about the shared_chat_sub_gift event.
	// Not presented if notice_type is not shared_chat_sub_gift.
	// This field has the same information as the chat_sub_gift field but for a notice that happened for a channel in a
	// shared chat session other than the broadcaster in the subscription condition.
	SharedChatSubGift *ChatNotificationEventSubGiftEvent `json:"shared_chat_sub_gift,omitempty"`
	// Information about the shared_chat_community_sub_gift event.
	// Not presented if notice_type is not shared_chat_community_sub_gift.
	// This field has the same information as the community_sub_gift field but for a notice that happened for a channel
	// in a shared chat session other than the broadcaster in the subscription condition.
	SharedChatCommunitySubGift *ChatNotificationEventCommunitySubGiftEvent `json:"shared_chat_community_sub_gift,omitempty"`
	// Information about the shared_chat_gift_paid_upgrade event.
	// Not presented if notice_type is not shared_chat_gift_paid_upgrade.
	// This field has the same information as the gift_paid_upgrade field but for a notice that happened for a channel
	// in a shared chat session other than the broadcaster in the subscription condition.
	SharedChatGiftPaidUpgrade *ChatNotificationEventGiftPaidUpgradeEvent `json:"shared_chat_gift_paid_upgrade,omitempty"`
	// Information about the shared_chat_chat_prime_paid_upgrade event.
	// Not presented if notice_type is not shared_chat_prime_paid_upgrade.
	// This field has the same information as the prime_paid_upgrade field but for a notice that happened for a channel
	// in a shared chat session other than the broadcaster in the subscription condition.
	SharedChatPrimePaidUpgrade *ChatNotificationEventGiftPaidUpgradeEvent `json:"shared_chat_prime_paid_upgrade,omitempty"`
	// Information about the shared_chat_pay_it_forward event.
	// Not presented if notice_type is not shared_chat_pay_it_forward.
	// This field has the same information as the pay_it_forward field but for a notice that happened for a channel in a shared
	// chat session other than the broadcaster in the subscription condition.
	SharedChatPayItForward *ChatNotificationEventPayItForwardEvent `json:"shared_chat_pay_it_forward,omitempty"`
	// Information about the shared_chat_raid event.
	// Not presented if notice_type is not shared_chat_raid.
	// This field has the same information as the raid field but for a notice that happened for a channel in a shared chat
	// session other than the broadcaster in the subscription condition.
	SharedChatRaid *ChatNotificationEventRaidEvent `json:"shared_chat_raid,omitempty"`
	// Information about the shared_chat_announcement event.
	// Not presented if notice_type is not shared_chat_announcement.
	// This field has the same information as the announcement field but for a notice that happened for a channel in a shared
	// chat session other than the broadcaster in the subscription condition.
	SharedChatAnnouncement *ChatNotificationEventAnnouncementEvent `json:"shared_chat_announcement,omitempty"`
}

type ChannelModeratorAddEvent struct {
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user ID of the new moderator.
	UserId string `json:"user_id"`
	// The user login of the new moderator.
	UserLogin string `json:"user_login"`
	// The display name of the new moderator.
	UserName string `json:"user_name"`
}

type ChannelModeratorRemoveEvent struct {
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user ID of the removed moderator.
	UserId string `json:"user_id"`
	// The user login of the removed moderator.
	UserLogin string `json:"user_login"`
	// The display name of the removed moderator.
	UserName string `json:"user_name"`
}

type ChannelPollBeginEvent struct {
	// ID of the poll.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Question displayed for the poll.
	Title string `json:"title"`
	// An array of choices for the poll.
	Choices []ChannelPollEventChoice `json:"choices"`
	// The bits voting settings for the poll.
	BitsVoting ChannelPollEventBitsVoting `json:"bits_voting"`
	// The voting settings for the poll.
	ChannelPointsVoting ChannelPollEventChannelPointsVoting `json:"channel_points_voting"`
	// The time the poll started.
	StartedAt TimestampUTC `json:"started_at"`
	// The time the poll will end.
	EndsAt TimestampUTC `json:"ends_at"`
}

type ChannelPollProgressEvent struct {
	// ID of the poll.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Question displayed for the poll.
	Title string `json:"title"`
	// An array of choices for the poll. Includes vote counts.
	Choices []ChannelPollEventChoice `json:"choices"`
	// The bits voting settings for the poll.
	BitsVoting ChannelPollEventBitsVoting `json:"bits_voting"`
	// The voting settings for the poll.
	ChannelPointsVoting ChannelPollEventChannelPointsVoting `json:"channel_points_voting"`
	// The time the poll started.
	StartedAt TimestampUTC `json:"started_at"`
	// The time the poll will end.
	EndsAt TimestampUTC `json:"ends_at"`
}

type ChannelPollEndEvent struct {
	// ID of the poll.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Question displayed for the poll.
	Title string `json:"title"`
	// An array of choices for the poll. Includes vote counts.
	Choices []ChannelPollEventChoice `json:"choices"`
	// The bits voting settings for the poll.
	BitsVoting ChannelPollEventBitsVoting `json:"bits_voting"`
	// The voting settings for the poll.
	ChannelPointsVoting ChannelPollEventChannelPointsVoting `json:"channel_points_voting"`
	// The status of the poll.
	Status ChannelPollEndEventStatus `json:"status"`
	// The time the poll started.
	StartedAt TimestampUTC `json:"started_at"`
	// The time the poll ended.
	EndedAt TimestampUTC `json:"ended_at"`
}

type ChannelPredictionBeginEvent struct {
	// The prediction ID.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Title for the prediction.
	Title string `json:"title"`
	// An array of outcomes for the prediction.
	Outcomes []ChannelPredictionEventOutcome `json:"outcomes"`
	// The time the prediction started.
	StartedAt TimestampUTC `json:"started_at"`
	// The time the prediction will automatically lock.
	LocksAt TimestampUTC `json:"locks_at"`
}

type ChannelPredictionProgressEvent struct {
	// The prediction ID.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Title for the Channel Points Prediction.
	Title string `json:"title"`
	// An array of outcomes for the prediction. Includes ChannelPredictionEventOutcome.TopPredictors.
	Outcomes []ChannelPredictionEventOutcome `json:"outcomes"`
	// The time the Channel prediction started.
	StartedAt TimestampUTC `json:"started_at"`
	// The time the Channel prediction will automatically lock.
	LocksAt TimestampUTC `json:"locks_at"`
}

type ChannelPredictionLockEvent struct {
	// The prediction ID.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Title for the Channel Points Prediction.
	Title string `json:"title"`
	// An array of outcomes for the prediction. Includes ChannelPredictionEventOutcome.TopPredictors.
	Outcomes []ChannelPredictionEventOutcome `json:"outcomes"`
	// The time the prediction started.
	StartedAt TimestampUTC `json:"started_at"`
	// The time the prediction was locked.
	LockedAt TimestampUTC `json:"locked_at"`
}

type ChannelPredictionEndEvent struct {
	// The prediction ID.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Title for the prediction.
	Title string `json:"title"`
	// ID of the winning outcome.
	WinningOutcomeId string `json:"winning_outcome_id"`
	// An array of outcomes for the prediction. Includes ChannelPredictionEventOutcome.TopPredictors.
	Outcomes []ChannelPredictionEventOutcome `json:"outcomes"`
	// The status of the prediction.
	Status ChannelPredictionEndEventStatus `json:"status"`
	// The time the prediction started.
	StartedAt TimestampUTC `json:"started_at"`
	// The time the prediction ended.
	EndedAt TimestampUTC `json:"ended_at"`
}

type ChannelRaidEvent struct {
	// The broadcaster ID that created the raid.
	FromBroadcasterUserId string `json:"from_broadcaster_user_id"`
	// The broadcaster login that created the raid.
	FromBroadcasterUserLogin string `json:"from_broadcaster_user_login"`
	// The broadcaster display name that created the raid.
	FromBroadcasterUserName string `json:"from_broadcaster_user_name"`
	// The broadcaster ID that received the raid.
	ToBroadcasterUserId string `json:"to_broadcaster_user_id"`
	// The broadcaster login that received the raid.
	ToBroadcasterUserLogin string `json:"to_broadcaster_user_login"`
	// The broadcaster display name that received the raid.
	ToBroadcasterUserName string `json:"to_broadcaster_user_name"`
	// The number of viewers in the raid.
	Viewers int `json:"viewers"`
}

type ChannelPointsCustomRewardRedemptionAddEvent struct {
	// The redemption identifier.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// User ID of the user that redeemed the reward.
	UserId string `json:"user_id"`
	// Login of the user that redeemed the reward.
	UserLogin string `json:"user_login"`
	// Display name of the user that redeemed the reward.
	UserName string `json:"user_name"`
	// The user input provided. Empty string if not provided.
	UserInput string `json:"user_input"`
	// Defaults to CustomRewardRedemptionStatusUnfulfilled.
	Status CustomRewardRedemptionStatus `json:"status"`
	// Basic information about the reward that was redeemed, at the time it was redeemed.
	Reward ChannelPointsCustomEventReward `json:"reward"`
	// RFC3339 timestamp of when the reward was redeemed.
	RedeemedAt TimestampUTC `json:"redeemed_at"`
}

type ChannelPointsCustomRewardRedemptionUpdateEvent struct {
	// The redemption identifier.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// User ID of the user that redeemed the reward.
	UserId string `json:"user_id"`
	// Login of the user that redeemed the reward.
	UserLogin string `json:"user_login"`
	// Display name of the user that redeemed the reward.
	UserName string `json:"user_name"`
	// The user input provided. Empty string if not provided.
	UserInput string `json:"user_input"`
	// Will be fulfilled or canceled.
	Status CustomRewardRedemptionStatus `json:"status"`
	// Basic information about the reward that was redeemed, at the time it was redeemed.
	Reward ChannelPointsCustomEventReward `json:"reward"`
	// RFC3339 timestamp of when the reward was redeemed.
	RedeemedAt TimestampUTC `json:"redeemed_at"`
}

type ChannelPointsAutomaticRewardRedemptionAddEvent struct {
	// The ID of the channel where the reward was redeemed.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the channel where the reward was redeemed.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The display name of the channel where the reward was redeemed.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The ID of the redeeming user.
	UserId string `json:"user_id"`
	// The login of the redeeming user.
	UserLogin string `json:"user_login"`
	// The display name of the redeeming user.
	UserName string `json:"user_name"`
	// The ID of the Redemption.
	Id string `json:"id"`
	// An object that contains the reward information.
	Reward ChannelPointsAutomaticRewardEventReward `json:"reward"`
	// An object that contains the user message and emote information needed to recreate the message.
	Message ChannelPointsAutomaticRewardEventRewardMessage `json:"message"`
	// Optional. A string that the user entered if the reward requires to be input.
	UserInput string `json:"user_input,omitempty"`
	// The UTC date and time (in RFC3339 format) of when the reward was
	RedeemedAt TimestampUTC `json:"redeemed_at"`
}

type ChannelPointsAutomaticRewardRedemptionAddEventV2 struct {
	// The ID of the channel where the reward was redeemed.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the channel where the reward was redeemed.
	BroadcasterUserLogin string `json:"broadcaster_user~_login"`
	// The display name of the channel where the reward was redeemed.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The ID of the redeeming user.
	UserId string `json:"user_id"`
	// The login of the redeeming user.
	UserLogin string `json:"user_login"`
	// The display name of the redeeming user.
	UserName string `json:"user_name"`
	// The ID of the Redemption.
	Id string `json:"id"`
	// An object that contains the reward information.
	Reward ChannelPointsAutomaticRewardEventRewardV2 `json:"reward"`
	// An object that contains the user message and emote information needed to recreate the message.
	Message ChannelPointsAutomaticRewardEventRewardMessageV2 `json:"message"`
	// The UTC date and time (in RFC3339 format) of when the reward was redeemed.
	RedeemedAt TimestampUTC `json:"redeemed_at"`
}

type UserAuthorizationRevokeEvent struct {
	// The client_id of the application with revoked user access.
	ClientId string `json:"client_id"`
	// The user id for the user who has revoked authorization for your client id.
	UserId string `json:"user_id"`
	// The user login for the user who has revoked authorization for your client id. Not presented if the user no longer exists.
	UserLogin string `json:"user_login,omitempty"`
	// The user display name for the user who has revoked authorization for your client id. Not presented if the user no longer exists.
	UserName string `json:"user_name,omitempty"`
}

type ChannelPointsCustomRewardAddEvent struct {
	// The reward identifier.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Is the reward currently enabled. If false, the reward won’t show up to viewers.
	IsEnabled bool `json:"is_enabled"`
	// Is the reward currently paused. If true, viewers can’t redeem.
	IsPaused bool `json:"is_paused"`
	// Is the reward currently in stock. If false, viewers can’t redeem.
	IsInStock bool `json:"is_in_stock"`
	// The reward title.
	Title string `json:"title"`
	// The reward cost.
	Cost int `json:"cost"`
	// The reward description.
	Prompt string `json:"prompt"`
	// Does the viewer need to enter information when redeeming the reward.
	IsUserInputRequired bool `json:"is_user_input_required"`
	// Should redemptions be set to fulfilled status immediately when redeemed and skip the request queue instead of the normal unfulfilled status.
	ShouldRedemptionsSkipRequestQueue bool `json:"should_redemptions_skip_request_queue"`
	// Whether a maximum per stream is enabled and what the maximum is.
	MaxPerStream ChannelPointsRewardEventMaxPerStream `json:"max_per_stream"`
	// Whether a maximum per user per stream is enabled and what the maximum is.
	MaxPerUserPerStream ChannelPointsRewardEventMaxPerUserPerStream `json:"max_per_user_per_stream"`
	// Custom background color for the reward. Format: Hex with # prefix. Example: #FA1ED2.
	BackgroundColor string `json:"background_color"`
	// Set of custom images of 1x, 2x and 4x sizes for the reward. Can be Not presented if no images have been uploaded.
	Image *ChannelPointsRewardEventImage `json:"image,omitempty"`
	// Set of default images of 1x, 2x and 4x sizes for the reward.
	DefaultImage ChannelPointsRewardEventImage `json:"default_image"`
	// Whether a cooldown is enabled and what the cooldown is in seconds.
	GlobalCooldown ChannelPointsRewardEventGlobalCooldown `json:"global_cooldown"`
	// Timestamp of the cooldown expiration. Not presented if the reward isn’t on cooldown.
	CooldownExpiresAt *TimestampUTC `json:"cooldown_expires_at,omitempty"`
	// The number of redemptions redeemed during the current live stream. Counts against the max_per_stream limit.
	// Not presented if the broadcasters stream isn’t live or max_per_stream isn’t enabled.
	RedemptionsRedeemedCurrentStream int `json:"redemptions_redeemed_current_stream"`
}

type ChannelPointsCustomRewardUpdateEvent struct {
	// The reward identifier.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Is the reward currently enabled. If false, the reward won’t show up to viewers.
	IsEnabled bool `json:"is_enabled"`
	// Is the reward currently paused. If true, viewers can’t redeem.
	IsPaused bool `json:"is_paused"`
	// Is the reward currently in stock. If false, viewers can’t redeem.
	IsInStock bool `json:"is_in_stock"`
	// The reward title.
	Title string `json:"title"`
	// The reward cost.
	Cost int `json:"cost"`
	// The reward description.
	Prompt string `json:"prompt"`
	// Does the viewer need to enter information when redeeming the reward.
	IsUserInputRequired bool `json:"is_user_input_required"`
	// Should redemptions be set to fulfilled status immediately when redeemed and skip the request queue instead of the normal unfulfilled status.
	ShouldRedemptionsSkipRequestQueue bool `json:"should_redemptions_skip_request_queue"`
	// Whether a maximum per stream is enabled and what the maximum is.
	MaxPerStream ChannelPointsRewardEventMaxPerStream `json:"max_per_stream"`
	// Whether a maximum per user per stream is enabled and what the maximum is.
	MaxPerUserPerStream ChannelPointsRewardEventMaxPerUserPerStream `json:"max_per_user_per_stream"`
	// Custom background color for the reward. Format: Hex with # prefix. Example: #FA1ED2.
	BackgroundColor string `json:"background_color"`
	// Set of custom images of 1x, 2x and 4x sizes for the reward. Can be Not presented if no images have been uploaded.
	Image *ChannelPointsRewardEventImage `json:"image,omitempty"`
	// Set of default images of 1x, 2x and 4x sizes for the reward.
	DefaultImage ChannelPointsRewardEventImage `json:"default_image"`
	// Whether a cooldown is enabled and what the cooldown is in seconds.
	GlobalCooldown ChannelPointsRewardEventGlobalCooldown `json:"global_cooldown"`
	// Timestamp of the cooldown expiration. Not presented if the reward isn’t on cooldown.
	CooldownExpiresAt *TimestampUTC `json:"cooldown_expires_at,omitempty"`
	// The number of redemptions redeemed during the current live stream. Counts against the max_per_stream limit.
	// Not presented if the broadcasters stream isn’t live or max_per_stream isn’t enabled.
	RedemptionsRedeemedCurrentStream int `json:"redemptions_redeemed_current_stream,omitempty"`
}

type ChannelPointsCustomRewardRemoveEvent struct {
	// The reward identifier.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Is the reward currently enabled. If false, the reward won’t show up to viewers.
	IsEnabled bool `json:"is_enabled"`
	// Is the reward currently paused. If true, viewers can’t redeem.
	IsPaused bool `json:"is_paused"`
	// Is the reward currently in stock. If false, viewers can’t redeem.
	IsInStock bool `json:"is_in_stock"`
	// The reward title.
	Title string `json:"title"`
	// The reward cost.
	Cost int `json:"cost"`
	// The reward description.
	Prompt string `json:"prompt"`
	// Does the viewer need to enter information when redeeming the reward.
	IsUserInputRequired bool `json:"is_user_input_required"`
	// Should redemptions be set to fulfilled status immediately when redeemed and skip the request queue instead of the normal unfulfilled status.
	ShouldRedemptionsSkipRequestQueue bool `json:"should_redemptions_skip_request_queue"`
	// Whether a maximum per stream is enabled and what the maximum is.
	MaxPerStream ChannelPointsRewardEventMaxPerStream `json:"max_per_stream"`
	// Whether a maximum per user per stream is enabled and what the maximum is.
	MaxPerUserPerStream ChannelPointsRewardEventMaxPerUserPerStream `json:"max_per_user_per_stream"`
	// Custom background color for the reward. Format: Hex with # prefix. Example: #FA1ED2.
	BackgroundColor string `json:"background_color"`
	// Set of custom images of 1x, 2x and 4x sizes for the reward. Not presented if no images have been uploaded.
	Image *ChannelPointsRewardEventImage `json:"image,omitempty"`
	// Set of default images of 1x, 2x and 4x sizes for the reward.
	DefaultImage ChannelPointsRewardEventImage `json:"default_image"`
	// Whether a cooldown is enabled and what the cooldown is in seconds.
	GlobalCooldown ChannelPointsRewardEventGlobalCooldown `json:"global_cooldown"`
	// Timestamp of the cooldown expiration. Not presented if the reward isn’t on cooldown.
	CooldownExpiresAt *TimestampUTC `json:"cooldown_expires_at,omitempty"`
	// The number of redemptions redeemed during the current live stream.
	// Counts against the max_per_stream limit.
	// Not presented if the broadcasters stream isn’t live or max_per_stream isn’t enabled.
	RedemptionsRedeemedCurrentStream int `json:"redemptions_redeemed_current_stream,omitempty"`
}

type StreamOnlineEvent struct {
	// The id of the stream.
	Id string `json:"id"`
	// The broadcaster’s user id.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster’s user login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s user display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The stream type.
	Type StreamType `json:"type"`
	// The timestamp at which the stream went online at.
	StartedAt TimestampUTC `json:"started_at"`
}

type StreamOfflineEvent struct {
	// The broadcaster’s user id.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster’s user login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s user display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
}

type ChannelSubscribeEvent struct {
	// The user ID for the user who subscribed to the specified channel.
	UserId string `json:"user_id"`
	// The user login for the user who subscribed to the specified channel.
	UserLogin string `json:"user_login"`
	// The user display name for the user who subscribed to the specified channel.
	UserName string `json:"user_name"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The tier of the subscription.
	Tier SubscriptionTier `json:"tier"`
	// Whether the subscription is a gift.
	IsGift bool `json:"is_gift"`
}

type ChannelSubscriptionMessageEvent struct {
	// The user ID of the user who sent a resubscription chat message.
	UserId string `json:"user_id"`
	// The user login of the user who sent a resubscription chat message.
	UserLogin string `json:"user_login"`
	// The user display name of the user who a resubscription chat message.
	UserName string `json:"user_name"`
	// The broadcaster user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The tier of the user’s subscription.
	Tier SubscriptionTier `json:"tier"`
	// An object that contains the resubscription message and emote information needed to recreate the message.
	Message ChannelSubscriptionMessageEventMessage `json:"message"`
	// The total number of months the user has been subscribed to the channel.
	CumulativeTotal int `json:"cumulative_total"`
	// The number of consecutive months the user’s current subscription has been active.
	// This value is not presented if the user has opted out of sharing this information.
	StreakMonths int `json:"streak_months,omitempty"`
	// The month duration of the subscription.
	DurationMonths int `json:"duration_months"`
}

type ChannelSubscriptionEndEvent struct {
	// The user ID for the user whose subscription ended.
	UserId string `json:"user_id"`
	// The user login for the user whose subscription ended.
	UserLogin string `json:"user_login"`
	// The user display name for the user whose subscription ended.
	UserName string `json:"user_name"`
	// The broadcaster user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The tier of the subscription that ended.
	Tier SubscriptionTier `json:"tier"`
	// Whether the subscription was a gift.
	IsGift bool `json:"is_gift"`
}

type ChannelSubscriptionGiftEvent struct {
	// The user ID of the user who sent the subscription gift. Not presented if it was an anonymous subscription gift.
	UserId string `json:"user_id,omitempty"`
	// The user login of the user who sent the gift. Not presented if it was an anonymous subscription gift.
	UserLogin string `json:"user_login,omitempty"`
	// The user display name of the user who sent the gift. Not presented if it was an anonymous subscription gift.
	UserName string `json:"user_name,omitempty"`
	// The broadcaster user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The number of subscriptions in the subscription gift.
	Total int `json:"total"`
	// The tier of subscriptions in the subscription gift.
	Tier string `json:"tier"`
	// The number of subscriptions gifted by this user in the channel.
	// Not presented for anonymous gifts or if the gifter has opted out of sharing this information.
	CumulativeTotal int `json:"cumulative_total,omitempty"`
	// Whether the subscription gift was anonymous.
	IsAnonymous bool `json:"is_anonymous"`
}

type ChannelUnbanRequestCreateEvent struct {
	// The ID of the unban request.
	Id string `json:"id"`
	// The broadcaster’s user ID for the channel the unban request was created for.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster’s login name.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// User ID of user that is requesting to be unbanned.
	UserId string `json:"user_id"`
	// The user’s login name.
	UserLogin string `json:"user_login"`
	// The user’s display name.
	UserName string `json:"user_name"`
	// Message sent in the unban request.
	Text string `json:"text"`
	// The UTC timestamp (in RFC3339 format) of when the unban request was created.
	CreatedAt TimestampUTC `json:"created_at"`
}

type ChannelUnbanRequestResolveEvent struct {
	// The ID of the unban request.
	Id string `json:"id"`
	// The broadcaster’s user ID for the channel the unban request was created for.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster’s login name.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Optional. User ID of moderator who approved/denied the request.
	ModeratorUserId string `json:"moderator_user_id,omitempty"`
	// Optional. The moderator’s login name
	ModeratorUserLogin string `json:"moderator_user_login,omitempty"`
	// Optional. The moderator’s display name
	ModeratorUserName string `json:"moderator_user_name,omitempty"`
	// User ID of user that requested to be unbanned.
	UserId string `json:"user_id"`
	// 	The user’s login name.
	UserLogin string `json:"user_login"`
	// The user’s display name.
	UserName string `json:"user_name"`
	// Optional. Resolution text supplied by the mod/broadcaster upon approval/denial of the request.
	ResolutionText string `json:"resolution_text,omitempty"`
	// Dictates whether the unban request was approved or denied.
	Status ChannelUnbanRequestResolveEventStatus `json:"status"`
}

type UserUpdateEvent struct {
	// The user’s user id.
	UserId string `json:"user_id"`
	// The user’s user login.
	UserLogin string `json:"user_login"`
	// The user’s user display name.
	UserName string `json:"user_name"`
	// The user’s email.
	// Only included if you have the user:read:email scope for the user.
	Email string `json:"email,omitempty"`
	// Determines whether Twitch has verified the user’s email address.
	// Is true if Twitch has verified the email address; otherwise, false.
	// Ignore this field if the email field contains an empty string.
	EmailVerified bool `json:"email_verified"`
	// The user’s description.
	Description string `json:"description"`
}

type ChannelVipAddEvent struct {
	// The user ID of the user that was added as a VIP.
	UserId string `json:"user_id"`
	// The user login of the user that was added as a VIP.
	UserLogin string `json:"user_login"`
	// The user display name of the user that was added as a VIP.
	UserName string `json:"user_name"`
	// The broadcaster user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
}

type ChannelVipRemoveEvent struct {
	// The user ID of the user that was removed as a VIP.
	UserId string `json:"user_id"`
	// The user login of the user that was removed as a VIP.
	UserLogin string `json:"user_login"`
	// The user display name of the user that was removed as a VIP.
	UserName string `json:"user_name"`
	// The broadcaster user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
}

type ChannelChatMessageDeleteEvent struct {
	// The broadcaster user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The ID of the user whose message was deleted.
	TargetUserId string `json:"target_user_id"`
	// The username of the user whose message was deleted.
	TargetUserName string `json:"target_user_name"`
	// The user login of the user whose message was deleted.
	TargetUserLogin string `json:"target_user_login"`
	// A UUID that identifies the message that was removed.
	MessageId string `json:"message_id"`
}

type ChannelChatSettingsUpdateEvent struct {
	// The ID of the broadcaster specified in the request.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the broadcaster specified in the request.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The user name of the broadcaster specified in the request.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// A Boolean value that determines whether chat messages must contain only emotes.
	EmoteMode bool `json:"emote_mode"`
	// A Boolean value that determines whether the broadcaster restricts the chat room to followers only.
	FollowerMode bool `json:"follower_mode"`
	// The length of time, in minutes, that the followers must have followed the broadcaster to participate in the chat room.
	// Null if follower_mode is false.
	FollowerModeDurationMinutes int `json:"follower_mode_duration_minutes"`
	// A Boolean value that determines whether the broadcaster limits how often users in the chat room are allowed to send messages.
	SlowMode bool `json:"slow_mode"`
	// The amount of time, in seconds, that users need to wait between sending messages. Null if slow_mode is false.
	SlowModeWaitTimeSeconds int `json:"slow_mode_wait_time_seconds"`
	// A Boolean value that determines whether only users that subscribe to the broadcaster's channel can talk in the chat room.
	SubscriberMode bool `json:"subscriber_mode"`
	// A Boolean value that determines whether the broadcaster requires users to post only unique messages in the chat room.
	UniqueChatMode bool `json:"unique_chat_mode"`
}

type ChannelChatUserMessageHoldEvent struct {
	// The ID of the broadcaster specified in the request.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the broadcaster specified in the request.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The user name of the broadcaster specified in the request.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The User ID of the message sender.
	UserId string `json:"user_id"`
	// The message sender's login.
	UserLogin string `json:"user_login"`
	// The message sender's display name.
	UserName string `json:"user_name"`
	// The ID of the message that was flagged by automod.
	MessageId string `json:"message_id"`
	// The body of the message.
	Message ChannelChatMessageEventMessage `json:"message"`
}

type ChannelChatUserMessageUpdateEvent struct {
	// The ID of the broadcaster specified in the request.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the broadcaster specified in the request.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The user name of the broadcaster specified in the request.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The User ID of the message sender.
	UserId string `json:"user_id"`
	// The message sender's login.
	UserLogin string `json:"user_login"`
	// The message sender's user name.
	UserName string `json:"user_name"`
	// The message's status. Possible values are: approved, denied, invalid.
	Status UserMessageUpdateStatus `json:"status"`
	// The ID of the message that was flagged by automod.
	MessageId string `json:"message_id"`
	// The body of the message.
	Message ChannelChatMessageEventMessage `json:"message"`
}

type ChannelCustomPowerUpRedemptionAddEvent struct {
	// The redemption identifier.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// User ID of the user that redeemed the custom Power-up.
	UserId string `json:"user_id"`
	// Login of the user that redeemed the custom Power-up.
	UserLogin string `json:"user_login"`
	// Display name of the user that redeemed the custom Power-up.
	UserName string `json:"user_name"`
	// The user input provided. Empty string if not provided.
	UserInput string `json:"user_input"`
	// Defaults to unfulfilled. Possible values are unknown, unfulfilled, fulfilled, and canceled.
	Status CustomRewardRedemptionStatus `json:"status"`
	// Basic information about the custom Power-up that was redeemed.
	CustomPowerUp CustomPowerUp `json:"custom_power_up"`
	// RFC3339 timestamp of when the custom Power-up was redeemed.
	RedeemedAt TimestampUTC `json:"redeemed_at"`
}

type ChannelSharedChatSessionBeginEvent struct {
	// The unique identifier for the shared chat session.
	SessionId string `json:"session_id"`
	// The User ID of the channel in the subscription condition which is now active in the shared chat session.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The display name of the channel in the subscription condition which is now active in the shared chat session.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user login of the channel in the subscription condition which is now active in the shared chat session.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The User ID of the host channel.
	HostBroadcasterUserId string `json:"host_broadcaster_user_id"`
	// The display name of the host channel.
	HostBroadcasterUserName string `json:"host_broadcaster_user_name"`
	// The user login of the host channel.
	HostBroadcasterUserLogin string `json:"host_broadcaster_user_login"`
	// The list of participants in the session.
	Participants []SharedChatParticipant `json:"participants"`
}

type ChannelSharedChatSessionUpdateEvent struct {
	// The unique identifier for the shared chat session.
	SessionId string `json:"session_id"`
	// The User ID of the channel in the subscription condition.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The display name of the channel in the subscription condition.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user login of the channel in the subscription condition.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The User ID of the host channel.
	HostBroadcasterUserId string `json:"host_broadcaster_user_id"`
	// The display name of the host channel.
	HostBroadcasterUserName string `json:"host_broadcaster_user_name"`
	// The user login of the host channel.
	HostBroadcasterUserLogin string `json:"host_broadcaster_user_login"`
	// The list of participants in the session.
	Participants []SharedChatParticipant `json:"participants"`
}

type ChannelSharedChatSessionEndEvent struct {
	// The unique identifier for the shared chat session.
	SessionId string `json:"session_id"`
	// The User ID of the channel in the subscription condition which is no longer active in the shared chat session.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The display name of the channel in the subscription condition which is no longer active in the shared chat session.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user login of the channel in the subscription condition which is no longer active in the shared chat session.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The User ID of the host channel.
	HostBroadcasterUserId string `json:"host_broadcaster_user_id"`
	// The display name of the host channel.
	HostBroadcasterUserName string `json:"host_broadcaster_user_name"`
	// The user login of the host channel.
	HostBroadcasterUserLogin string `json:"host_broadcaster_user_login"`
}

type ChannelSuspiciousUserMessageEvent struct {
	// The ID of the channel where the treatment for a suspicious user was updated.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The display name of the channel where the treatment for a suspicious user was updated.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The login of the channel where the treatment for a suspicious user was updated.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The user ID of the user that sent the message.
	UserId string `json:"user_id"`
	// The user name of the user that sent the message.
	UserName string `json:"user_name"`
	// The user login of the user that sent the message.
	UserLogin string `json:"user_login"`
	// The status set for the suspicious user. Can be none, active_monitoring, or restricted.
	LowTrustStatus SuspiciousUserLowTrustStatus `json:"low_trust_status"`
	// A list of channel IDs where the suspicious user is also banned.
	SharedBanChannelIds []string `json:"shared_ban_channel_ids"`
	// User types (if any) that apply to the suspicious user.
	Types []SuspiciousUserType `json:"types"`
	// A ban evasion likelihood value (if any) that has been applied to the user automatically by Twitch.
	BanEvasionEvaluation SuspiciousUserBanEvasionEvaluation `json:"ban_evasion_evaluation"`
	// The structured chat message.
	Message SuspiciousUserMessageMessage `json:"message"`
}

type ChannelSuspiciousUserUpdateEvent struct {
	// The ID of the channel where the treatment for a suspicious user was updated.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The display name of the channel where the treatment for a suspicious user was updated.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The login of the channel where the treatment for a suspicious user was updated.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The ID of the moderator that updated the treatment for a suspicious user.
	ModeratorUserId string `json:"moderator_user_id"`
	// The display name of the moderator that updated the treatment for a suspicious user.
	ModeratorUserName string `json:"moderator_user_name"`
	// The login of the moderator that updated the treatment for a suspicious user.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The ID of the suspicious user whose treatment was updated.
	UserId string `json:"user_id"`
	// The display name of the suspicious user whose treatment was updated.
	UserName string `json:"user_name"`
	// The login of the suspicious user whose treatment was updated.
	UserLogin string `json:"user_login"`
	// The status set for the suspicious user. Can be none, active_monitoring, or restricted.
	LowTrustStatus SuspiciousUserLowTrustStatus `json:"low_trust_status"`
}

type ChannelWarningAcknowledgeEvent struct {
	// The user ID of the broadcaster.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the broadcaster.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The user name of the broadcaster.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The ID of the user that has acknowledged their warning.
	UserId string `json:"user_id"`
	// The login of the user that has acknowledged their warning.
	UserLogin string `json:"user_login"`
	// The user name of the user that has acknowledged their warning.
	UserName string `json:"user_name"`
}

type ChannelWarningSendEvent struct {
	// The user ID of the broadcaster.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The login of the broadcaster.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The user name of the broadcaster.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The user ID of the moderator who sent the warning.
	ModeratorUserId string `json:"moderator_user_id"`
	// The login of the moderator.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The user name of the moderator.
	ModeratorUserName string `json:"moderator_user_name"`
	// The ID of the user being warned.
	UserId string `json:"user_id"`
	// The login of the user being warned.
	UserLogin string `json:"user_login"`
	// The user name of the user being warned.
	UserName string `json:"user_name"`
	// Optional. The reason given for the warning.
	Reason string `json:"reason,omitempty"`
	// Optional. The chat rules cited for the warning.
	ChatRulesCited []string `json:"chat_rules_cited,omitempty"`
}

type ChannelCheerEvent struct {
	// Whether the user cheered anonymously or not.
	IsAnonymous bool `json:"is_anonymous"`
	// The user ID for the user who cheered on the specified channel. Null if is_anonymous is true.
	UserId string `json:"user_id,omitempty"`
	// The user login for the user who cheered on the specified channel. Null if is_anonymous is true.
	UserLogin string `json:"user_login,omitempty"`
	// The user display name for the user who cheered on the specified channel. Null if is_anonymous is true.
	UserName string `json:"user_name,omitempty"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The message sent with the cheer.
	Message string `json:"message"`
	// The number of Bits cheered.
	Bits int `json:"bits"`
}

type CharityCampaignDonateEvent struct {
	// An ID that identifies the donation. The ID is unique across campaigns.
	Id string `json:"id"`
	// An ID that identifies the charity campaign.
	CampaignId string `json:"campaign_id"`
	// An ID that identifies the broadcaster that's running the campaign.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster's login name.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster's display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// An ID that identifies the user that donated to the campaign.
	UserId string `json:"user_id"`
	// The user's login name.
	UserLogin string `json:"user_login"`
	// The user's display name.
	UserName string `json:"user_name"`
	// The charity's name.
	CharityName string `json:"charity_name"`
	// A description of the charity.
	CharityDescription string `json:"charity_description"`
	// A URL to an image of the charity's logo.
	CharityLogo string `json:"charity_logo"`
	// A URL to the charity's website.
	CharityWebsite string `json:"charity_website"`
	// An object that contains the amount of money that the user donated.
	Amount CharityAmount `json:"amount"`
}

type CharityCampaignStartEvent struct {
	// An ID that identifies the charity campaign.
	Id string `json:"id"`
	// An ID that identifies the broadcaster that's running the campaign.
	BroadcasterUserId string `json:"broadcaster_id"`
	// The broadcaster's login name.
	BroadcasterUserLogin string `json:"broadcaster_login"`
	// The broadcaster's display name.
	BroadcasterUserName string `json:"broadcaster_name"`
	// The charity's name.
	CharityName string `json:"charity_name"`
	// A description of the charity.
	CharityDescription string `json:"charity_description"`
	// A URL to an image of the charity's logo.
	CharityLogo string `json:"charity_logo"`
	// A URL to the charity's website.
	CharityWebsite string `json:"charity_website"`
	// An object that contains the current amount of donations that the campaign has received.
	CurrentAmount CharityAmount `json:"current_amount"`
	// An object that contains the campaign's target fundraising goal.
	TargetAmount CharityAmount `json:"target_amount"`
	// The UTC timestamp (in RFC3339 format) of when the broadcaster started the campaign.
	StartedAt TimestampUTC `json:"started_at"`
}

type CharityCampaignProgressEvent struct {
	// An ID that identifies the charity campaign.
	Id string `json:"id"`
	// An ID that identifies the broadcaster that's running the campaign.
	BroadcasterUserId string `json:"broadcaster_id"`
	// The broadcaster's login name.
	BroadcasterUserLogin string `json:"broadcaster_login"`
	// The broadcaster's display name.
	BroadcasterUserName string `json:"broadcaster_name"`
	// The charity's name.
	CharityName string `json:"charity_name"`
	// A description of the charity.
	CharityDescription string `json:"charity_description"`
	// A URL to an image of the charity's logo.
	CharityLogo string `json:"charity_logo"`
	// A URL to the charity's website.
	CharityWebsite string `json:"charity_website"`
	// An object that contains the current amount of donations that the campaign has received.
	CurrentAmount CharityAmount `json:"current_amount"`
	// An object that contains the campaign's target fundraising goal.
	TargetAmount CharityAmount `json:"target_amount"`
}

type CharityCampaignStopEvent struct {
	// An ID that identifies the charity campaign.
	Id string `json:"id"`
	// An ID that identifies the broadcaster that ran the campaign.
	BroadcasterUserId string `json:"broadcaster_id"`
	// The broadcaster's login name.
	BroadcasterUserLogin string `json:"broadcaster_login"`
	// The broadcaster's display name.
	BroadcasterUserName string `json:"broadcaster_name"`
	// The charity's name.
	CharityName string `json:"charity_name"`
	// A description of the charity.
	CharityDescription string `json:"charity_description"`
	// A URL to an image of the charity's logo.
	CharityLogo string `json:"charity_logo"`
	// A URL to the charity's website.
	CharityWebsite string `json:"charity_website"`
	// An object that contains the final amount of donations that the campaign received.
	CurrentAmount CharityAmount `json:"current_amount"`
	// An object that contains the campaign's target fundraising goal.
	TargetAmount CharityAmount `json:"target_amount"`
	// The UTC timestamp (in RFC3339 format) of when the broadcaster stopped the campaign.
	StoppedAt TimestampUTC `json:"stopped_at"`
}

type DropEntitlementGrantEvent struct {
	// Individual event ID, as assigned by EventSub.
	Id string `json:"id"`
	// Entitlement object.
	Data []DropEntitlementGrantData `json:"data"`
}

type ExtensionBitsTransactionCreateEvent struct {
	// ExtensionClientId is the client ID of the extension.
	ExtensionClientId string `json:"extension_client_id"`
	// Id is the transaction ID.
	Id string `json:"id"`
	// The transaction's broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The transaction's broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The transaction's broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The transaction's user ID.
	UserId string `json:"user_id"`
	// The transaction's user login.
	UserLogin string `json:"user_login"`
	// The transaction's user display name.
	UserName string `json:"user_name"`
	// Additional extension product information.
	Product ExtensionBitsTransactionProduct `json:"product"`
}

type ChannelGoalBeginEvent struct {
	// An ID that identifies this event.
	Id string `json:"id"`
	// An ID that uniquely identifies the broadcaster.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster's display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The broadcaster's user handle.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The type of goal.
	Type GoalType `json:"type"`
	// A description of the goal, if specified.
	Description string `json:"description"`
	// The goal's current value.
	CurrentAmount int `json:"current_amount"`
	// The goal's target value.
	TargetAmount int `json:"target_amount"`
	// The UTC timestamp of when the broadcaster created the goal.
	StartedAt TimestampUTC `json:"started_at"`
}

type ChannelGoalProgressEvent struct {
	// An ID that identifies this event.
	Id string `json:"id"`
	// An ID that uniquely identifies the broadcaster.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster's display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The broadcaster's user handle.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The type of goal.
	Type GoalType `json:"type"`
	// A description of the goal, if specified.
	Description string `json:"description"`
	// The goal's current value.
	CurrentAmount int `json:"current_amount"`
	// The goal's target value.
	TargetAmount int `json:"target_amount"`
	// The UTC timestamp of when the broadcaster created the goal.
	StartedAt TimestampUTC `json:"started_at"`
}

type ChannelGoalEndEvent struct {
	// An ID that identifies this event.
	Id string `json:"id"`
	// An ID that uniquely identifies the broadcaster.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster's display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The broadcaster's user handle.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The type of goal.
	Type GoalType `json:"type"`
	// A description of the goal, if specified.
	Description string `json:"description"`
	// A Boolean value that indicates whether the broadcaster achieved their goal.
	IsAchieved bool `json:"is_achieved"`
	// The goal's current value.
	CurrentAmount int `json:"current_amount"`
	// The goal's target value.
	TargetAmount int `json:"target_amount"`
	// The UTC timestamp of when the broadcaster created the goal.
	StartedAt TimestampUTC `json:"started_at"`
	// The UTC timestamp of when the broadcaster ended the goal.
	EndedAt TimestampUTC `json:"ended_at"`
}

type HypeTrainBeginEvent struct {
	// The Hype Train ID.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Total points contributed to the Hype Train.
	Total int `json:"total"`
	// The number of points contributed to the Hype Train at the current level.
	Progress int `json:"progress"`
	// The number of points required to reach the next level.
	Goal int `json:"goal"`
	// The contributors with the most points contributed.
	TopContributions []HypeTrainContribution `json:"top_contributions"`
	// The current level of the Hype Train.
	Level int `json:"level"`
	// The all-time high level this type of Hype Train has reached for this broadcaster.
	AllTimeHighLevel int `json:"all_time_high_level"`
	// The all-time high total this type of Hype Train has reached for this broadcaster.
	AllTimeHighTotal int `json:"all_time_high_total"`
	// Optional. Non-null for a shared Hype Train.
	SharedTrainParticipants []SharedChatParticipant `json:"shared_train_participants,omitempty"`
	// The time when the Hype Train started.
	StartedAt TimestampUTC `json:"started_at"`
	// The time when the Hype Train expires.
	ExpiresAt TimestampUTC `json:"expires_at"`
	// The type of the Hype Train.
	Type HypeTrainType `json:"type"`
	// Indicates if the Hype Train is shared.
	IsSharedTrain bool `json:"is_shared_train"`
}

type HypeTrainProgressEvent struct {
	// The Hype Train ID.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Total points contributed to the Hype Train.
	Total int `json:"total"`
	// The number of points contributed to the Hype Train at the current level.
	Progress int `json:"progress"`
	// The number of points required to reach the next level.
	Goal int `json:"goal"`
	// The contributors with the most points contributed.
	TopContributions []HypeTrainContribution `json:"top_contributions"`
	// The current level of the Hype Train.
	Level int `json:"level"`
	// Optional. Non-null for a shared Hype Train.
	SharedTrainParticipants []SharedChatParticipant `json:"shared_train_participants,omitempty"`
	// The time when the Hype Train started.
	StartedAt TimestampUTC `json:"started_at"`
	// The time when the Hype Train expires.
	ExpiresAt TimestampUTC `json:"expires_at"`
	// The type of the Hype Train.
	Type HypeTrainType `json:"type"`
	// Indicates if the Hype Train is shared.
	IsSharedTrain bool `json:"is_shared_train"`
}

type HypeTrainEndEvent struct {
	// The Hype Train ID.
	Id string `json:"id"`
	// The requested broadcaster ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The requested broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The requested broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Total points contributed to the Hype Train.
	Total int `json:"total"`
	// The contributors with the most points contributed.
	TopContributions []HypeTrainContribution `json:"top_contributions"`
	// The current level of the Hype Train.
	Level int `json:"level"`
	// Optional. Non-null for a shared Hype Train.
	SharedTrainParticipants []SharedChatParticipant `json:"shared_train_participants,omitempty"`
	// The time when the Hype Train started.
	StartedAt TimestampUTC `json:"started_at"`
	// The time when the Hype Train cooldown ends so that the next Hype Train can start.
	CooldownEndsAt TimestampUTC `json:"cooldown_ends_at"`
	// The time when the Hype Train ended.
	EndedAt TimestampUTC `json:"ended_at"`
	// The type of the Hype Train.
	Type HypeTrainType `json:"type"`
	// Indicates if the Hype Train is shared.
	IsSharedTrain bool `json:"is_shared_train"`
}

type ShieldModeBeginEvent struct {
	// The user ID of the broadcaster.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster's display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The broadcaster's login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The user ID of the moderator who activated Shield Mode.
	ModeratorUserId string `json:"moderator_user_id"`
	// The moderator's display name.
	ModeratorUserName string `json:"moderator_user_name"`
	// The moderator's login.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The UTC timestamp of when Shield Mode was activated.
	StartedAt TimestampUTC `json:"started_at"`
}

type ShieldModeEndEvent struct {
	// The user ID of the broadcaster.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster's display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The broadcaster's login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The user ID of the moderator who deactivated Shield Mode.
	ModeratorUserId string `json:"moderator_user_id"`
	// The moderator's display name.
	ModeratorUserName string `json:"moderator_user_name"`
	// The moderator's login.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The UTC timestamp of when Shield Mode was deactivated.
	EndedAt TimestampUTC `json:"ended_at"`
}

type ShoutoutCreateEvent struct {
	// The user ID of the broadcaster that sent the Shoutout.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster's display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The broadcaster's login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The user ID of the moderator that sent the Shoutout.
	ModeratorUserId string `json:"moderator_user_id"`
	// The moderator's display name.
	ModeratorUserName string `json:"moderator_user_name"`
	// The moderator's login.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The user ID of the broadcaster that received the Shoutout.
	ToBroadcasterUserId string `json:"to_broadcaster_user_id"`
	// The display name of the broadcaster that received the Shoutout.
	ToBroadcasterUserName string `json:"to_broadcaster_user_name"`
	// The login of the broadcaster that received the Shoutout.
	ToBroadcasterUserLogin string `json:"to_broadcaster_user_login"`
	// The UTC timestamp of when the Shoutout was sent.
	StartedAt TimestampUTC `json:"started_at"`
	// The number of users that were watching the broadcaster's stream at the time of the Shoutout.
	ViewerCount int `json:"viewer_count"`
	// The UTC timestamp of when the broadcaster may send a new Shoutout.
	CooldownEndsAt TimestampUTC `json:"cooldown_ends_at"`
	// The UTC timestamp of when the broadcaster may send a Shoutout to a specified broadcaster.
	TargetCooldownEndsAt TimestampUTC `json:"target_cooldown_ends_at"`
}

type ShoutoutReceiveEvent struct {
	// The user ID of the broadcaster that received the Shoutout.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster's display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The broadcaster's login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The user ID of the broadcaster that sent the Shoutout.
	FromBroadcasterUserId string `json:"from_broadcaster_user_id"`
	// The display name of the broadcaster that sent the Shoutout.
	FromBroadcasterUserName string `json:"from_broadcaster_user_name"`
	// The login of the broadcaster that sent the Shoutout.
	FromBroadcasterUserLogin string `json:"from_broadcaster_user_login"`
	// The number of users that were watching the from-broadcaster's stream at the time of the Shoutout.
	ViewerCount int `json:"viewer_count"`
	// The UTC timestamp of when the Shoutout was received.
	StartedAt TimestampUTC `json:"started_at"`
}

type UserAuthorizationGrantEvent struct {
	// The client_id of the application that was granted user access.
	ClientId string `json:"client_id"`
	// The user id for the user who has granted authorization for your client id.
	UserId string `json:"user_id"`
	// The user login for the user who has granted authorization for your client id.
	UserLogin string `json:"user_login"`
	// The user display name for the user who has granted authorization for your client id.
	UserName string `json:"user_name"`
}

type WhisperReceivedEvent struct {
	// The ID of the user sending the message.
	FromUserId string `json:"from_user_id"`
	// The name of the user sending the message.
	FromUserName string `json:"from_user_name"`
	// The login of the user sending the message.
	FromUserLogin string `json:"from_user_login"`
	// The ID of the user receiving the message.
	ToUserId string `json:"to_user_id"`
	// The name of the user receiving the message.
	ToUserName string `json:"to_user_name"`
	// The login of the user receiving the message.
	ToUserLogin string `json:"to_user_login"`
	// The whisper ID.
	WhisperId string `json:"whisper_id"`
	// Object containing whisper information.
	Whisper WhisperMessageBody `json:"whisper"`
}

type ChannelGuestStarSessionBeginEvent struct {
	// The broadcaster user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// ID representing the unique session that was started.
	SessionId string `json:"session_id"`
	// RFC3339 timestamp indicating the time the session began.
	StartedAt TimestampUTC `json:"started_at"`
}

type ChannelGuestStarSessionEndEvent struct {
	// The non-host broadcaster user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The non-host broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The non-host broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// ID representing the unique session that was started.
	SessionId string `json:"session_id"`
	// RFC3339 timestamp indicating the time the session began.
	StartedAt TimestampUTC `json:"started_at"`
	// RFC3339 timestamp indicating the time the session ended.
	EndedAt TimestampUTC `json:"ended_at"`
	// User ID of the host channel.
	HostUserId string `json:"host_user_id"`
	// The host display name.
	HostUserName string `json:"host_user_name"`
	// The host login.
	HostUserLogin string `json:"host_user_login"`
}

type ChannelGuestStarGuestUpdateEvent struct {
	// The non-host broadcaster user ID.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The non-host broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The non-host broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// ID representing the unique session that was started.
	SessionId string `json:"session_id"`
	// The user ID of the moderator who updated the guest's state. Null if the update was performed by the guest.
	ModeratorUserId string `json:"moderator_user_id"`
	// The moderator display name. Null if the update was performed by the guest.
	ModeratorUserName string `json:"moderator_user_name"`
	// The moderator login. Null if the update was performed by the guest.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The user ID of the guest who transitioned states in the session. Null if the slot is now empty.
	GuestUserId string `json:"guest_user_id"`
	// The guest display name. Null if the slot is now empty.
	GuestUserName string `json:"guest_user_name"`
	// The guest login. Null if the slot is now empty.
	GuestUserLogin string `json:"guest_user_login"`
	// The ID of the slot assignment the guest is assigned to. Null if the guest is in the INVITED, REMOVED, READY, or ACCEPTED state.
	SlotId string `json:"slot_id"`
	// The current state of the user after the update has taken place. Null if the slot is now empty.
	State GuestStarGuestState `json:"state"`
	// User ID of the host channel.
	HostUserId string `json:"host_user_id"`
	// The host display name.
	HostUserName string `json:"host_user_name"`
	// The host login.
	HostUserLogin string `json:"host_user_login"`
	// Flag that signals whether the host is allowing the slot's video to be seen by participants.
	HostVideoEnabled bool `json:"host_video_enabled"`
	// Flag that signals whether the host is allowing the slot's audio to be heard by participants.
	HostAudioEnabled bool `json:"host_audio_enabled"`
	// Value between 0-100 that represents the slot's audio level as heard by participants.
	HostVolume int `json:"host_volume"`
}

type ChannelGuestStarSettingsUpdateEvent struct {
	// User ID of the host channel.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// The broadcaster display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// The broadcaster login.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// Flag determining if Guest Star moderators have access to control whether a guest is live once assigned to a slot.
	IsModeratorSendLiveEnabled bool `json:"is_moderator_send_live_enabled"`
	// Number of slots the Guest Star call interface will allow the host to add to a call.
	SlotCount int `json:"slot_count"`
	// Flag determining if browser sources subscribed to sessions on this channel should output audio.
	IsBrowserSourceAudioEnabled bool `json:"is_browser_source_audio_enabled"`
	// This setting determines how the guests within a session should be laid out within a group browser source.
	GroupLayout GuestStarGroupLayout `json:"group_layout"`
}

type ChannelModerateEventV2 struct {
	// The ID of the broadcaster.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The login of the broadcaster.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The user name of the broadcaster.
	BroadcasterUserName string `json:"broadcaster_user_name"`

	// The channel in which the action originally occurred. Is the same as the broadcaster_user_id if not in shared chat.
	SourceBroadcasterUserID string `json:"source_broadcaster_user_id"`
	// The channel in which the action originally occurred. Is the same as the broadcaster_user_login if not in shared chat.
	SourceBroadcasterUserLogin string `json:"source_broadcaster_user_login"`
	// The channel in which the action originally occurred. Is null when the moderator action happens in the same channel as the broadcaster. Is not null when in a shared chat session, and the action happens in the channel of a participant other than the broadcaster.
	SourceBroadcasterUserName string `json:"source_broadcaster_user_name"`

	// The ID of the moderator who performed the action.
	ModeratorUserID string `json:"moderator_user_id"`
	// The login of the moderator.
	ModeratorUserLogin string `json:"moderator_user_login"`
	// The user name of the moderator.
	ModeratorUserName string `json:"moderator_user_name"`

	// The action performed. Possible values are
	// ban, timeout, unban, untimeout, clear, emoteonly, emoteonlyoff, followers, followersoff, uniquechat, uniquechatoff, slow, slowoff, subscribers, subscribersoff, unraid, delete, unvip, vip, raid, add_blocked_term, add_permitted_term, remove_blocked_term, remove_permitted_term, mod, unmod, approve_unban_request, deny_unban_request, warn, shared_chat_ban, shared_chat_timeout, shared_chat_unban, shared_chat_untimeout, shared_chat_delete
	Action ChannelModerateEventV2Action `json:"action"`

	// Optional. Metadata associated with the followers command.
	Followers *ChannelModerateEventV2FollowersMetadata `json:"followers,omitempty"`
	// Optional. Metadata associated with the slow command.
	Slow *ChannelModerateEventV2SlowMetadata `json:"slow,omitempty"`
	// Optional. Metadata associated with the vip command.
	Vip *ChannelModerateEventV2VipMetadata `json:"vip,omitempty"`
	// Optional. Metadata associated with the unvip command.
	Unvip *ChannelModerateEventV2UnvipMetadata `json:"unvip,omitempty"`
	// Optional. Metadata associated with the mod command.
	Mod *ChannelModerateEventV2ModMetadata `json:"mod,omitempty"`
	// Optional. Metadata associated with the unmod command.
	Unmod *ChannelModerateEventV2UnmodMetadata `json:"unmod,omitempty"`
	// Optional. Metadata associated with the ban command.
	Ban *ChannelModerateEventV2BanMetadata `json:"ban,omitempty"`
	// Optional. Metadata associated with the unban command.
	Unban *ChannelModerateEventV2UnbanMetadata `json:"unban,omitempty"`
	// Optional. Metadata associated with the timeout command.
	Timeout *ChannelModerateEventV2TimeoutMetadata `json:"timeout,omitempty"`
	// Optional. Metadata associated with the untimeout command.
	Untimeout *ChannelModerateEventV2UntimeoutMetadata `json:"untimeout,omitempty"`
	// Optional. Metadata associated with the raid command.
	Raid *ChannelModerateEventV2RaidMetadata `json:"raid,omitempty"`
	// Optional. Metadata associated with the unraid command.
	Unraid *ChannelModerateEventV2UnraidMetadata `json:"unraid,omitempty"`
	// Optional. Metadata associated with the delete command.
	Delete *ChannelModerateEventV2DeleteMetadata `json:"delete,omitempty"`
	// Optional. Metadata associated with the automod terms changes.
	AutomodTerms *ChannelModerateEventV2AutomodTermsMetadata `json:"automod_terms,omitempty"`
	// Optional. Metadata associated with an unban request.
	UnbanRequest *ChannelModerateEventV2UnbanRequestMetadata `json:"unban_request,omitempty"`
	// Optional. Metadata associated with the warn command.
	Warn *ChannelModerateEventV2WarnMetadata `json:"warn,omitempty"`
	// Optional. Information about the shared_chat_ban event.
	SharedChatBan *ChannelModerateEventV2BanMetadata `json:"shared_chat_ban,omitempty"`
	// Optional. Information about the shared_chat_unban event.
	SharedChatUnban *ChannelModerateEventV2UnbanMetadata `json:"shared_chat_unban,omitempty"`
	// Optional. Information about the shared_chat_timeout event.
	SharedChatTimeout *ChannelModerateEventV2TimeoutMetadata `json:"shared_chat_timeout,omitempty"`
	// Optional. Information about the shared_chat_untimeout event.
	SharedChatUntimeout *ChannelModerateEventV2UntimeoutMetadata `json:"shared_chat_untimeout,omitempty"`
	// Optional. Information about the shared_chat_delete event.
	SharedChatDelete *ChannelModerateEventV2DeleteMetadata `json:"shared_chat_delete,omitempty"`
}

type ChannelModerateEventV2Action string

const (
	ChannelModerateEventV2ActionBan                 ChannelModerateEventV2Action = "ban"
	ChannelModerateEventV2ActionTimeout             ChannelModerateEventV2Action = "timeout"
	ChannelModerateEventV2ActionUnban               ChannelModerateEventV2Action = "unban"
	ChannelModerateEventV2ActionUntimeout           ChannelModerateEventV2Action = "untimeout"
	ChannelModerateEventV2ActionClear               ChannelModerateEventV2Action = "clear"
	ChannelModerateEventV2ActionEmoteOnly           ChannelModerateEventV2Action = "emoteonly"
	ChannelModerateEventV2ActionEmoteOnlyOff        ChannelModerateEventV2Action = "emoteonlyoff"
	ChannelModerateEventV2ActionFollowers           ChannelModerateEventV2Action = "followers"
	ChannelModerateEventV2ActionFollowersOff        ChannelModerateEventV2Action = "followersoff"
	ChannelModerateEventV2ActionUniqueChat          ChannelModerateEventV2Action = "uniquechat"
	ChannelModerateEventV2ActionUniqueChatOff       ChannelModerateEventV2Action = "uniquechatoff"
	ChannelModerateEventV2ActionSlow                ChannelModerateEventV2Action = "slow"
	ChannelModerateEventV2ActionSlowOff             ChannelModerateEventV2Action = "slowoff"
	ChannelModerateEventV2ActionSubscribers         ChannelModerateEventV2Action = "subscribers"
	ChannelModerateEventV2ActionSubscribersOff      ChannelModerateEventV2Action = "subscribersoff"
	ChannelModerateEventV2ActionUnraid              ChannelModerateEventV2Action = "unraid"
	ChannelModerateEventV2ActionDelete              ChannelModerateEventV2Action = "delete"
	ChannelModerateEventV2ActionUnvip               ChannelModerateEventV2Action = "unvip"
	ChannelModerateEventV2ActionVip                 ChannelModerateEventV2Action = "vip"
	ChannelModerateEventV2ActionRaid                ChannelModerateEventV2Action = "raid"
	ChannelModerateEventV2ActionAddBlockedTerm      ChannelModerateEventV2Action = "add_blocked_term"
	ChannelModerateEventV2ActionAddPermittedTerm    ChannelModerateEventV2Action = "add_permitted_term"
	ChannelModerateEventV2ActionRemoveBlockedTerm   ChannelModerateEventV2Action = "remove_blocked_term"
	ChannelModerateEventV2ActionRemovePermittedTerm ChannelModerateEventV2Action = "remove_permitted_term"
	ChannelModerateEventV2ActionMod                 ChannelModerateEventV2Action = "mod"
	ChannelModerateEventV2ActionUnmod               ChannelModerateEventV2Action = "unmod"
	ChannelModerateEventV2ActionApproveUnbanRequest ChannelModerateEventV2Action = "approve_unban_request"
	ChannelModerateEventV2ActionDenyUnbanRequest    ChannelModerateEventV2Action = "deny_unban_request"
	ChannelModerateEventV2ActionWarn                ChannelModerateEventV2Action = "warn"
	ChannelModerateEventV2ActionSharedChatBan       ChannelModerateEventV2Action = "shared_chat_ban"
	ChannelModerateEventV2ActionSharedChatTimeout   ChannelModerateEventV2Action = "shared_chat_timeout"
	ChannelModerateEventV2ActionSharedChatUnban     ChannelModerateEventV2Action = "shared_chat_unban"
	ChannelModerateEventV2ActionSharedChatUntimeout ChannelModerateEventV2Action = "shared_chat_untimeout"
	ChannelModerateEventV2ActionSharedChatDelete    ChannelModerateEventV2Action = "shared_chat_delete"
)

type ChannelModerateEventV2FollowersMetadata struct {
	// The length of time, in minutes, that the followers must have followed the broadcaster to participate in the chat room.
	FollowDurationMinutes int `json:"follow_duration_minutes"`
}

type ChannelModerateEventV2SlowMetadata struct {
	// The amount of time, in seconds, that users need to wait between sending messages.
	WaitTimeSeconds int `json:"wait_time_seconds"`
}

type ChannelModerateEventV2VipMetadata struct {
	// The ID of the user gaining VIP status.
	UserID string `json:"user_id"`
	// The login of the user gaining VIP status.
	UserLogin string `json:"user_login"`
	// The user name of the user gaining VIP status.
	UserName string `json:"user_name"`
}

type ChannelModerateEventV2UnvipMetadata struct {
	// The ID of the user losing VIP status.
	UserID string `json:"user_id"`
	// The login of the user losing VIP status.
	UserLogin string `json:"user_login"`
	// The user name of the user losing VIP status.
	UserName string `json:"user_name"`
}

type ChannelModerateEventV2ModMetadata struct {
	// The ID of the user gaining mod status.
	UserID string `json:"user_id"`
	// The login of the user gaining mod status.
	UserLogin string `json:"user_login"`
	// The user name of the user gaining mod status.
	UserName string `json:"user_name"`
}

type ChannelModerateEventV2UnmodMetadata struct {
	// The ID of the user losing mod status.
	UserID string `json:"user_id"`
	// The login of the user losing mod status.
	UserLogin string `json:"user_login"`
	// The user name of the user losing mod status.
	UserName string `json:"user_name"`
}

type ChannelModerateEventV2BanMetadata struct {
	// The ID of the user being banned.
	UserID string `json:"user_id"`
	// The login of the user being banned.
	UserLogin string `json:"user_login"`
	// The user name of the user being banned.
	UserName string `json:"user_name"`
	// Optional. Reason given for the ban.
	Reason string `json:"reason,omitempty"`
}

type ChannelModerateEventV2UnbanMetadata struct {
	// The ID of the user being unbanned.
	UserID string `json:"user_id"`
	// The login of the user being unbanned.
	UserLogin string `json:"user_login"`
	// The user name of the user being unbanned.
	UserName string `json:"user_name"`
}

type ChannelModerateEventV2TimeoutMetadata struct {
	// The ID of the user being timed out.
	UserID string `json:"user_id"`
	// The login of the user being timed out.
	UserLogin string `json:"user_login"`
	// The user name of the user being timed out.
	UserName string `json:"user_name"`
	// Optional. The reason given for the timeout.
	Reason string `json:"reason,omitempty"`
	// The UTC timestamp of when the ad break began, in RFC3339 format.
	// The time at which the timeout ends.
	ExpiresAt TimestampUTC `json:"expires_at"`
}

type ChannelModerateEventV2UntimeoutMetadata struct {
	// The ID of the user being untimed out.
	UserID string `json:"user_id"`
	// The login of the user being untimed out.
	UserLogin string `json:"user_login"`
	// The user name of the user untimed out.
	UserName string `json:"user_name"`
}

type ChannelModerateEventV2RaidMetadata struct {
	// The ID of the user being raided.
	UserID string `json:"user_id"`
	// The login of the user being raided.
	UserLogin string `json:"user_login"`
	// The user name of the user raided.
	UserName string `json:"user_name"`
	// The viewer count.
	ViewerCount int `json:"viewer_count"`
}

type ChannelModerateEventV2UnraidMetadata struct {
	// The ID of the user no longer being raided.
	UserID string `json:"user_id"`
	// The login of the user no longer being raided.
	UserLogin string `json:"user_login"`
	// The user name of the no longer user raided.
	UserName string `json:"user_name"`
}

type ChannelModerateEventV2DeleteMetadata struct {
	// The ID of the user whose message is being deleted.
	UserID string `json:"user_id"`
	// The login of the user.
	UserLogin string `json:"user_login"`
	// The user name of the user.
	UserName string `json:"user_name"`
	// The ID of the message being deleted.
	MessageID string `json:"message_id"`
	// The message body of the message being deleted.
	MessageBody string `json:"message_body"`
}

type ChannelModerateEventV2AutomodTermsMetadata struct {
	// Either “add” or “remove”.
	Action string `json:"action"`
	// Either “blocked” or “permitted”.
	List string `json:"list"`
	// Terms being added or removed.
	Terms []string `json:"terms"`
	// Whether the terms were added due to an Automod message approve/deny action.
	FromAutomod bool `json:"from_automod"`
}

type ChannelModerateEventV2UnbanRequestMetadata struct {
	// Whether or not the unban request was approved or denied.
	IsApproved bool `json:"is_approved"`
	// The ID of the banned user.
	UserID string `json:"user_id"`
	// The login of the user.
	UserLogin string `json:"user_login"`
	// The user name of the user.
	UserName string `json:"user_name"`
	// The message included by the moderator explaining their approval or denial.
	ModeratorMessage string `json:"moderator_message"`
}

type ChannelModerateEventV2WarnMetadata struct {
	// The ID of the user being warned.
	UserID string `json:"user_id"`
	// The login of the user being warned.
	UserLogin string `json:"user_login"`
	// The user name of the user being warned.
	UserName string `json:"user_name"`
	// Optional. Reason given for the warning.
	Reason string `json:"reason,omitempty"`
	// Optional. Chat rules cited for the warning.
	ChatRulesCited []string `json:"chat_rules_cited,omitempty"`
}
