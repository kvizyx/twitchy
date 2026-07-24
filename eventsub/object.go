package eventsub

type AutomodSeverityLevel int

const (
	AutomodSeverityFirst AutomodSeverityLevel = iota + 1
	AutomodSeveritySecond
	AutomodSeverityThird
	AutomodSeverityFourth
)

type AutomodHoldReason string

const (
	AutomodHoldReasonAutomod     = "automod"
	AutomodHoldReasonBlockedTerm = "blocked_term"
)

type AutomodMessageStatus string

const (
	AutomodMessageStatusApproved AutomodMessageStatus = "approved"
	AutomodMessageStatusDenied   AutomodMessageStatus = "denied"
	AutomodMessageStatusExpired  AutomodMessageStatus = "expired"
)

type AutomodTermsAction string

const (
	AutomodTermsAddPermitted    AutomodTermsAction = "add_permitted"
	AutomodTermsRemovePermitted AutomodTermsAction = "remove_permitted"
	AutomodTermsAddBlocked      AutomodTermsAction = "add_blocked"
	AutomodTermsRemoveBlocked   AutomodTermsAction = "remove_blocked"
)

type MessageFragmentType string

const (
	MessageFragmentText      MessageFragmentType = "text"
	MessageFragmentEmote     MessageFragmentType = "emote"
	MessageFragmentCheermote MessageFragmentType = "cheermote"
	MessageFragmentMention   MessageFragmentType = "mention"
	MessageFragmentGif       MessageFragmentType = "gif"
)

type Mention struct {
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login"`
}

type SubscriptionTier string

const (
	SubscriptionTierOne   SubscriptionTier = "1000"
	SubscriptionTierTwo   SubscriptionTier = "2000"
	SubscriptionTierThree SubscriptionTier = "3000"
)

type EmoteFormat string

const (
	EmoteFormatAnimated EmoteFormat = "animated"
	EmoteFormatStatic   EmoteFormat = "static"
)

type RewardEmote struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Cheermote struct {
	Prefix string `json:"prefix"`
	Bits   int    `json:"bits"`
	Tier   int    `json:"tier"`
}

type AutomodMessageHoldEventMessage struct {
	Text      string                                   `json:"text"`
	Fragments []AutomodMessageHoldEventMessageFragment `json:"fragments"`
}

type AutomodMessageHoldEventMessageFragment struct {
	Text      string                               `json:"text"`
	Emote     *AutomodMessageHoldEventMessageEmote `json:"emote,omitempty"`
	Cheermote *Cheermote                           `json:"cheermote,omitempty"`
}

type AutomodMessageHoldEventMessageEmote struct {
	Id         string `json:"id"`
	EmoteSetId string `json:"emote_set_id"`
}

type AutomodMessageHoldEventMessageV2 struct {
	Text      string                                     `json:"text"`
	Fragments []AutomodMessageHoldEventMessageFragmentV2 `json:"fragments"`
}

type AutomodMessageHoldEventMessageFragmentV2 struct {
	Type      MessageFragmentType                    `json:"type"`
	Text      string                                 `json:"text"`
	Emote     *AutomodMessageHoldEventMessageEmoteV2 `json:"emote,omitempty"`
	Cheermote *Cheermote                             `json:"cheermote,omitempty"`
}

type AutomodMessageHoldEventMessageEmoteV2 struct {
	Id         string `json:"id"`
	EmoteSetId string `json:"emote_set_id"`
}

type AutomodMessageUpdateEventMessage struct {
	Text      string                                     `json:"text"`
	Fragments []AutomodMessageUpdateEventMessageFragment `json:"fragments"`
}

type AutomodMessageUpdateEventMessageFragment struct {
	Text      string                                 `json:"text"`
	Emote     *AutomodMessageUpdateEventMessageEmote `json:"emote,omitempty"`
	Cheermote *Cheermote                             `json:"cheermote,omitempty"`
}

type AutomodMessageUpdateEventMessageEmote struct {
	Id         string `json:"id"`
	EmoteSetId string `json:"emote_set_id"`
}

type AutomodMessageUpdateEventMessageV2 struct {
	Text      string                                       `json:"text"`
	Fragments []AutomodMessageUpdateEventMessageFragmentV2 `json:"fragments"`
}

type AutomodMessageUpdateEventMessageFragmentV2 struct {
	Type      MessageFragmentType                      `json:"type"`
	Text      string                                   `json:"text"`
	Emote     *AutomodMessageUpdateEventMessageEmoteV2 `json:"emote,omitempty"`
	Cheermote *Cheermote                               `json:"cheermote,omitempty"`
}

type AutomodMessageUpdateEventMessageEmoteV2 struct {
	Id         string `json:"id"`
	EmoteSetId string `json:"emote_set_id"`
}

type ChannelBitsUseEventMessage struct {
	Text      string                               `json:"text"`
	Fragments []ChannelBitsUseEventMessageFragment `json:"fragments"`
}

type ChannelBitsUseEventMessageFragment struct {
	Type      MessageFragmentType              `json:"type"`
	Text      string                           `json:"text"`
	Emote     *ChannelBitsUseEventMessageEmote `json:"emote,omitempty"`
	Cheermote *Cheermote                       `json:"cheermote,omitempty"`
}

type ChannelBitsUseEventMessageEmote struct {
	Id         string        `json:"id"`
	EmoteSetId string        `json:"emote_set_id"`
	OwnerId    string        `json:"owner_id"`
	Format     []EmoteFormat `json:"format"`
}

type ChannelChatMessageEventMessage struct {
	Text      string                                   `json:"text"`
	Fragments []ChannelChatMessageEventMessageFragment `json:"fragments"`
}

type ChannelChatMessageEventMessageFragment struct {
	Type      MessageFragmentType                  `json:"type"`
	Text      string                               `json:"text"`
	Emote     *ChannelChatMessageEventMessageEmote `json:"emote,omitempty"`
	Cheermote *Cheermote                           `json:"cheermote,omitempty"`
	Mention   *Mention                             `json:"mention"`
	Gif       *ChannelChatMessageEventMessageGif   `json:"gif,omitempty"`
}

type ChannelChatMessageEventMessageGif struct {
	// GifId is an ID that uniquely identifies this GIF.
	GifId string `json:"gif_id"`
	// URL is the URL of the GIF asset. Applications rendering the GIF must use the full URL provided; it must not be modified.
	URL string `json:"url"`
}

type ChannelChatMessageEventMessageEmote struct {
	Id         string        `json:"id"`
	EmoteSetId string        `json:"emote_set_id"`
	OwnerId    string        `json:"owner_id"`
	Format     []EmoteFormat `json:"format"`
}

type MessageEmote struct {
	// The emote ID.
	Id string `json:"id"`
	// The index of where the Emote starts in the text.
	Begin int `json:"begin"`
	// The index of where the Emote ends in the text.
	End int `json:"end"`
}

type Automod struct {
	Category   string               `json:"category"`
	Level      AutomodSeverityLevel `json:"level"`
	Boundaries []AutomodBoundary    `json:"boundaries"`
}

type AutomodBoundary struct {
	StartPosition int `json:"start_position"`
	EndPosition   int `json:"end_position"`
}

type Term struct {
	TermId                    string          `json:"term_id"`
	Boundary                  AutomodBoundary `json:"boundary"`
	OwnerBroadcasterUserId    string          `json:"owner_broadcaster_user_id"`
	OwnerBroadcasterUserLogin string          `json:"owner_broadcaster_user_login"`
	OwnerBroadcasterUserName  string          `json:"owner_broadcaster_user_name"`
}

type BlockedTerm struct {
	TermsFound []Term `json:"terms_found"`
}

type BitsUseType string

const (
	BitsUseCheer         BitsUseType = "cheer"
	BitsUsePowerUp       BitsUseType = "power_up"
	BitsUseCustomPowerUp BitsUseType = "custom_power_up"
)

type BitsPowerUpType string

const (
	BitsPowerUpMessageEffect    BitsPowerUpType = "message_effect"
	BitsPowerUpCelebration      BitsPowerUpType = "celebration"
	BitsPowerUpGigantifyAnEmote BitsPowerUpType = "gigantify_an_emote"
)

type PowerUp struct {
	Type            BitsPowerUpType `json:"type"`
	Emote           *RewardEmote    `json:"emote,omitempty"`
	MessageEffectId string          `json:"message_effect_id,omitempty"`
}

type ChannelBitsUseCustomPowerUp struct {
	// The title of the custom Power-up.
	Title string `json:"title"`
	// The ID of the custom Power-up.
	RewardId string `json:"reward_id"`
}

type MessageType string

const (
	MessageText                     MessageType = "text"
	MessageChannelPointsHighlighted MessageType = "channel_points_highlighted"
	MessageChannelPointsSubOnly     MessageType = "channel_points_sub_only"
	MessageUserIntro                MessageType = "user_intro"
	MessagePowerUpsMessageEffect    MessageType = "power_ups_message_effect"
	MessagePowerUpsGigantifiedEmote MessageType = "power_ups_gigantified_emote"
)

type Badge struct {
	Id    string `json:"id"`
	SetId string `json:"set_id"`
	Info  string `json:"info"`
}

type Cheer struct {
	Bits int `json:"bits"`
}

type Reply struct {
	ParentMessageId   string `json:"parent_message_id"`
	ParentMessageBody string `json:"parent_message_body"`
	ParentUserId      string `json:"parent_user_id"`
	ParentUserName    string `json:"parent_user_name"`
	ParentUserLogin   string `json:"parent_user_login"`
	ThreadMessageId   string `json:"thread_message_id"`
	ThreadUserId      string `json:"thread_user_id"`
	ThreadUserName    string `json:"thread_user_name"`
	ThreadUserLogin   string `json:"thread_user_login"`
}

type ConduitShardDisabledEventTransport struct {
	// Method websocket or webhook
	Method string `json:"method"`
	// Optional. Webhook callback URL. Empty if method is set to websocket.
	Callback string `json:"callback,omitempty"`
	// Optional. WebSocket session ID. Empty if  method is set to webhook.
	SessionId string `json:"session_id,omitempty"`
	// Optional. Time that the WebSocket session connected. Empty if method is set to webhook.
	ConnectedAt *TimestampUTC `json:"connected_at,omitempty"`
	// Optional. Time that the WebSocket session disconnected. Empty if method is set to webhook.
	DisconnectAt *TimestampUTC `json:"disconnected_at,omitempty"`
}

type ChatNotificationEventMessage struct {
	// The chat message in plain text.
	Text string `json:"text"`
	// Ordered list of chat message fragments.
	Fragments []ChatNotificationEventMessageFragment `json:"fragments"`
}

type ChatNotificationEventMessageFragment struct {
	// The type of message fragment.
	Type MessageFragmentType `json:"type"`
	// Message text in fragment
	Text string `json:"text"`
	// Optional. Metadata pertaining to the cheermote.
	Cheermote *ChatNotificationEventMessageFragmentCheermote `json:"cheermote,omitempty"`
	// Optional. Metadata pertaining to the emote.
	Emote *ChatNotificationEventMessageFragmentEmote `json:"emote,omitempty"`
	// Optional. Metadata pertaining to the mention.
	Mention *ChatNotificationEventMessageFragmentMention `json:"mention,omitempty"`
}

type ChatNotificationEventMessageFragmentCheermote struct {
	// The name portion of the ChatNotificationMessageFragmentCheermote string that you use in chat to cheer Bits.
	// The full ChatNotificationMessageFragmentCheermote string is the concatenation of {prefix} + {number of Bits}.
	// For example, if the prefix is “Cheer” and you want to cheer 100 Bits, the full ChatNotificationMessageFragmentCheermote
	// string is Cheer100. When the ChatNotificationMessageFragmentCheermote string is entered in chat, Twitch converts it to
	// the image associated with the Bits tier that was cheered.
	Prefix string `json:"prefix"`
	// The amount of bits cheered.
	Bits int `json:"bits"`
	// The tier level of the cheermote.
	Tier int `json:"tier"`
}

type ChatNotificationEventMessageFragmentEmote struct {
	// An ID that uniquely identifies this emote.
	Id string `json:"id"`
	// An ID that identifies the emote set that the emote belongs to.
	EmoteSetId string `json:"emote_set_id"`
	// The ID of the broadcaster who owns the emote.
	OwnerId string `json:"owner_id"`
	// The formats that the emote is available in.
	// For example, if the emote is available only as a static PNG, the array contains only EmoteFormatStatic, but if the
	// emote is available as a static PNG and an animated GIF, the array contains EmoteFormatStatic and EmoteFormatAnimated.
	Format []EmoteFormat `json:"format"`
}

type ChatNotificationEventMessageFragmentMention struct {
	// The user ID of the mentioned user.
	UserId string `json:"user_id"`
	// The username of the mentioned user.
	UserName string `json:"user_name"`
	// The user login of the mentioned user.
	UserLogin string `json:"user_login"`
}

type ChatNotificationEventBadge struct {
	// An ID that identifies this set of chat badges. For example, Bits or Subscriber.
	SetId string `json:"set_id"`
	// An ID that identifies this version of the badge.
	// The ID can be any value.
	// For example, for Bits, the ID is the Bits tier level, but for World of Warcraft, it could be Alliance or Horde.
	Id string `json:"id"`
	// Contains metadata related to the chat badges in the badges tag.
	// Currently, this tag contains metadata only for subscriber badges, to indicate the number of months the user has been a subscriber.
	Info string `json:"info"`
}

type ChatNotificationEventSubEvent struct {
	// The type of subscription plan being used.
	SubTier SubscriptionTier `json:"sub_tier"`
	// Indicates if the subscription was obtained through Amazon Prime.
	IsPrime bool `json:"is_prime"`
	// The number of months the subscription is for.
	DurationMonths int `json:"duration_months"`
}

type ChatNotificationEventReSubEvent struct {
	// The total number of months the user has subscribed.
	CumulativeMonths int `json:"cumulative_months"`
	// The number of months the subscription is for.
	DurationMonths int `json:"duration_months"`
	// Optional. The number of consecutive months the user has subscribed.
	StreakMonths int `json:"streak_months,omitempty"`
	// The type of subscription plan being used.
	SubTier SubscriptionTier `json:"sub_tier"`
	// Indicates if the re-sub was obtained through Amazon Prime.
	IsPrime bool `json:"is_prime"`
	// Whether the re-sub was a result of a gift or not.
	IsGift bool `json:"is_gift"`
	// Optional. Whether the gift was anonymous or not.
	GifterIsAnonymous bool `json:"gifter_is_anonymous,omitempty"`
	// Optional. The user ID of the subscription gifter. Not presented if anonymous.
	GifterUserId string `json:"gifter_user_id,omitempty"`
	// Optional. The username of the subscription gifter. Not presented if anonymous.
	GifterUserName string `json:"gifter_user_name,omitempty"`
	// Optional. The user login of the subscription gifter. Not presented if anonymous.
	GifterUserLogin string `json:"gifter_user_login,omitempty"`
}

type ChatNotificationEventSubGiftEvent struct {
	// The number of months the subscription is for.
	DurationMonths int `json:"duration_months"`
	// Optional. The amount of gifts the gifter has given in this channel. Not presented if anonymous.
	CumulativeTotal int `json:"cumulative_total,omitempty"`
	// The user ID of the subscription gift recipient.
	RecipientUserId string `json:"recipient_user_id"`
	// The username of the subscription gift recipient.
	RecipientUserName string `json:"recipient_user_name"`
	// The user login of the subscription gift recipient.
	RecipientUserLogin string `json:"recipient_user_login"`
	// The type of subscription plan being used.
	SubTier SubscriptionTier `json:"sub_tier"`
	// Optional. The ID of the associated community gift. Not presented if not associated with a community gift.
	CommunityGiftId string `json:"community_gift_id,omitempty"`
}

type ChatNotificationEventCommunitySubGiftEvent struct {
	// The ID of the associated community gift.
	Id string `json:"id"`
	// Number of subscriptions being gifted.
	Total int `json:"total"`
	// The type of subscription plan being used.
	SubTier SubscriptionTier `json:"sub_tier"`
	// Optional. The amount of gifts the gifter has given in this channel. Not presented if anonymous.
	CumulativeTotal int `json:"cumulative_total,omitempty"`
}

type ChatNotificationEventGiftPaidUpgradeEvent struct {
	// Whether the gift was given anonymously.
	GifterIsAnonymous bool `json:"gifter_is_anonymous"`
	// Optional. The user ID of the user who gifted the subscription. Not presented if anonymous.
	GifterUserId string `json:"gifter_user_id,omitempty"`
	// Optional. The username of the user who gifted the subscription. Not presented if anonymous.
	GifterUserName string `json:"gifter_user_name,omitempty"`
	// Optional. The user login of the user who gifted the subscription. Not presented if anonymous.
	GifterUserLogin string `json:"gifter_user_login,omitempty"`
}

// ChatNotificationEventRaidEvent represents information about the raid event.
type ChatNotificationEventRaidEvent struct {
	// The user ID of the broadcaster raiding this channel.
	UserId string `json:"user_id"`
	// The username of the broadcaster raiding this channel.
	UserName string `json:"user_name"`
	// The login name of the broadcaster raiding this channel.
	UserLogin string `json:"user_login"`
	// The number of viewers raiding this channel from the broadcaster’s channel.
	ViewerCount int `json:"viewer_count"`
	// Profile image URL of the broadcaster raiding this channel.
	ProfileImageURL string `json:"profile_image_url"`
}

type ChatNotificationEventUnRaidEvent struct{}

type ChatNotificationEventPayItForwardEvent struct {
	// Whether the gift was given anonymously.
	GifterIsAnonymous bool `json:"gifter_is_anonymous"`
	// Optional. The user ID of the user who gifted the subscription. Not presented if anonymous.
	GifterUserId string `json:"gifter_user_id,omitempty"`
	// Optional. The username of the user who gifted the subscription. Not presented if anonymous.
	GifterUserName string `json:"gifter_user_name,omitempty"`
	// Optional. The user login of the user who gifted the subscription. Not presented if anonymous.
	GifterUserLogin string `json:"gifter_user_login,omitempty"`
}

type ChatNotificationEventAnnouncementEvent struct {
	// Color of the announcement.
	Color string `json:"color"`
}

type ChatNotificationEventCharityDonationEvent struct {
	// Name of the charity.
	CharityName string `json:"charity_name"`
	// An object that contains the amount of money that the user paid.
	Amount ChatNotificationEventCharityDonationEventDonationAmount `json:"amount"`
}

type ChatNotificationEventCharityDonationEventDonationAmount struct {
	// The monetary amount. The amount is specified in the currency’s minor unit.
	// For example, the minor units for USD is cents, so if the amount is $5.50 USD, value is set to 550.
	Value int `json:"value"`
	// The number of decimal places used by the currency.
	// For example, USD uses two decimal places.
	DecimalPlaces int `json:"decimal_places"`
	// The ISO-4217 three-letter currency code that identifies the type of currency in value.
	Currency string `json:"currency"`
}

// ChatNotificationEventBitsBadgeTierEvent represents information about the bits badge tier event.
type ChatNotificationEventBitsBadgeTierEvent struct {
	// The tier of the Bits badge the user just earned. For example, 100, 1000, or 10000.
	Tier int `json:"tier"`
}

// ChatNotificationEventWatchStreakEvent represents information about the watch streak event.
type ChatNotificationEventWatchStreakEvent struct {
	// The number of consecutive broadcasts for which the user has been watching.
	StreakCount int `json:"streak_count"`
	// The number of channel points awarded for the Watch Streak milestone.
	ChannelPointsAwarded int `json:"channel_points_awarded"`
}

// ChatNotificationEventModiversaryEvent represents information about the modiversary event.
type ChatNotificationEventModiversaryEvent struct {
	// The number of months the user has been a moderator in this channel.
	Months int `json:"months"`
}

type ChatNotificationEventNoticeType string

const (
	ChatNotificationEventNoticeTypeSub                        ChatNotificationEventNoticeType = "sub"
	ChatNotificationEventNoticeTypeReSub                      ChatNotificationEventNoticeType = "resub"
	ChatNotificationEventNoticeTypeSubGift                    ChatNotificationEventNoticeType = "sub_gift"
	ChatNotificationEventNoticeTypeCommunitySubGift           ChatNotificationEventNoticeType = "community_sub_gift"
	ChatNotificationEventNoticeTypeGiftPaidUpgrade            ChatNotificationEventNoticeType = "gift_paid_upgrade"
	ChatNotificationEventNoticeTypePrimePaidUpgrade           ChatNotificationEventNoticeType = "prime_paid_upgrade"
	ChatNotificationEventNoticeTypeRaid                       ChatNotificationEventNoticeType = "raid"
	ChatNotificationEventNoticeTypeUnRaid                     ChatNotificationEventNoticeType = "unraid"
	ChatNotificationEventNoticeTypePayItForward               ChatNotificationEventNoticeType = "pay_it_forward"
	ChatNotificationEventNoticeTypeAnnouncement               ChatNotificationEventNoticeType = "announcement"
	ChatNotificationEventNoticeTypeBitsBadgeTier              ChatNotificationEventNoticeType = "bits_badge_tier"
	ChatNotificationEventNoticeTypeCharityDonation            ChatNotificationEventNoticeType = "charity_donation"
	ChatNotificationEventNoticeTypeWatchStreak                ChatNotificationEventNoticeType = "watch_streak"
	ChatNotificationEventNoticeTypeModiversary                ChatNotificationEventNoticeType = "modiversary"
	ChatNotificationEventNoticeTypeSharedChatSub              ChatNotificationEventNoticeType = "shared_chat_sub"
	ChatNotificationEventNoticeTypeSharedChatReSub            ChatNotificationEventNoticeType = "shared_chat_resub"
	ChatNotificationEventNoticeTypeSharedChatSubGift          ChatNotificationEventNoticeType = "shared_chat_sub_gift"
	ChatNotificationEventNoticeTypeSharedChatCommunitySubGift ChatNotificationEventNoticeType = "shared_chat_community_sub_gift"
	ChatNotificationEventNoticeTypeSharedChatGiftPaidUpgrade  ChatNotificationEventNoticeType = "shared_chat_gift_paid_upgrade"
	ChatNotificationEventNoticeTypeSharedChatPrimePaidUpgrade ChatNotificationEventNoticeType = "shared_chat_prime_paid_upgrade"
	ChatNotificationEventNoticeTypeSharedChatRaid             ChatNotificationEventNoticeType = "shared_chat_raid"
	ChatNotificationEventNoticeTypeSharedChatPayItForward     ChatNotificationEventNoticeType = "shared_chat_pay_it_forward"
	ChatNotificationEventNoticeTypeSharedChatAnnouncement     ChatNotificationEventNoticeType = "shared_chat_announcement"
	ChatNotificationEventNoticeTypeSharedChatModiversary      ChatNotificationEventNoticeType = "shared_chat_modiversary"
	ChatNotificationEventNoticeTypeUnknown                    ChatNotificationEventNoticeType = "unknown"
)

func (c ChatNotificationEventNoticeType) String() string {
	return string(c)
}

type ChannelPollEventChoice struct {
	// ID for the choice.
	Id string `json:"id"`
	// Text displayed for the choice.
	Title string `json:"title"`
	// Number of votes received via Bits.
	BitsVotes int `json:"bits_votes"`
	// Number of votes received via Channel Points.
	ChannelPointsVotes int `json:"channel_points_votes"`
	// Total number of votes received for the choice across all methods of voting.
	Votes int `json:"votes"`
}

type ChannelPollEventBitsVoting struct {
	// Indicates if Bits can be used for voting.
	IsEnabled bool `json:"is_enabled"`
	// Number of Bits required to vote once with Bits.
	AmountPerVote int `json:"amount_per_vote"`
}

type ChannelPollEventChannelPointsVoting struct {
	// Indicates if Channel Points can be used for voting.
	IsEnabled bool `json:"is_enabled"`
	// Number of Channel Points required to vote once with Channel Points.
	AmountPerVote int `json:"amount_per_vote"`
}

type ChannelPollEndEventStatus string

const (
	ChannelPollEndEventStatusCompleted  ChannelPollEndEventStatus = "completed"
	ChannelPollEndEventStatusArchived   ChannelPollEndEventStatus = "archived"
	ChannelPollEndEventStatusTerminated ChannelPollEndEventStatus = "terminated"
)

func (c ChannelPollEndEventStatus) String() string {
	return string(c)
}

type ChannelPredictionEventOutcome struct {
	// The outcome ID.
	Id string `json:"id"`
	// The outcome title.
	Title string `json:"title"`
	// The color for the outcome. Valid values are pink and blue.
	Color string `json:"color"`
	// The number of users who used Channel Points on this outcome.
	Users int `json:"users"`
	// The total number of Channel Points used on this outcome.
	ChannelPoints int `json:"channel_points"`
	// An array of users who used the most Channel Points on this outcome.
	TopPredictors []ChannelPredictionEventOutcomeTopPredictor `json:"top_predictors"`
}

type ChannelPredictionEventOutcomeTopPredictor struct {
	// The ID of the user.
	UserId string `json:"user_id"`
	// The login of the user.
	UserLogin string `json:"user_login"`
	// The display name of the user.
	UserName string `json:"user_name"`
	// The number of Channel Points won.
	// This value is always null in the event payload for Prediction progress and Prediction lock.
	// This value is 0 if the outcome did not win or if the Prediction was canceled and Channel Points were refunded.
	ChannelPointsWon int `json:"channel_points_won"`
	// The number of Channel Points used to participate in the Prediction.
	ChannelPointsUsed int `json:"channel_points_used"`
}

type ChannelPredictionEndEventStatus string

const (
	ChannelPredictionEndStatusResolved ChannelPredictionEndEventStatus = "resolved"
	ChannelPredictionEndStatusCanceled ChannelPredictionEndEventStatus = "canceled"
)

type ChannelPointsCustomEventReward struct {
	// The reward identifier.
	Id string `json:"id"`
	// The reward name.
	Title string `json:"title"`
	// The reward cost.
	Cost int `json:"cost"`
	// The reward description.
	Prompt string `json:"prompt"`
}

type ChannelPointsAutomaticRewardEventRewardType string

const (
	ChannelPointsAutomaticRewardEventRewardTypeSingleMessageBypassSubMode   ChannelPointsAutomaticRewardEventRewardType = "single_message_bypass_sub_mode"
	ChannelPointsAutomaticRewardEventRewardTypeSendHighlightedMessage       ChannelPointsAutomaticRewardEventRewardType = "send_highlighted_message"
	ChannelPointsAutomaticRewardEventRewardTypeRandomSubEmoteUnlock         ChannelPointsAutomaticRewardEventRewardType = "random_sub_emote_unlock"
	ChannelPointsAutomaticRewardEventRewardTypeChosenSubEmoteUnlock         ChannelPointsAutomaticRewardEventRewardType = "chosen_sub_emote_unlock"
	ChannelPointsAutomaticRewardEventRewardTypeChosenModifiedSubEmoteUnlock ChannelPointsAutomaticRewardEventRewardType = "chosen_modified_sub_emote_unlock"
	ChannelPointsAutomaticRewardEventRewardTypeMessageEffect                ChannelPointsAutomaticRewardEventRewardType = "message_effect"
	ChannelPointsAutomaticRewardEventRewardTypeGigantifyAnEmote             ChannelPointsAutomaticRewardEventRewardType = "gigantify_an_emote"
	ChannelPointsAutomaticRewardEventRewardTypeCelebration                  ChannelPointsAutomaticRewardEventRewardType = "celebration"
)

type ChannelPointsAutomaticRewardEventReward struct {
	// The type of reward.
	Type ChannelPointsAutomaticRewardEventRewardType `json:"type"`
	// The reward cost.
	Cost int
	// Optional. Emote that was unlocked.
	UnlockedEmote *ChannelPointsAutomaticRewardEventRewardUnlockedEmote `json:"unlocked_emote,omitempty"`
}

type ChannelPointsAutomaticRewardEventRewardUnlockedEmote struct {
	// The emote ID.
	Id string `json:"id"`
	// The human-readable emote token.
	Name string `json:"name"`
}

type ChannelPointsAutomaticRewardEventRewardMessage struct {
	// The text of the chat message.
	Text string `json:"text"`
	// An array that includes the emote ID and start and end positions for where the emote appears in the text.
	Emotes []MessageEmote `json:"emotes,omitempty"`
}

type ChannelPointsAutomaticRewardEventRewardTypeV2 string

const (
	ChannelPointsAutomaticRewardEventRewardTypeV2SingleMessageBypassSubMode   ChannelPointsAutomaticRewardEventRewardTypeV2 = "single_message_bypass_sub_mode"
	ChannelPointsAutomaticRewardEventRewardTypeV2SendHighlightedMessage       ChannelPointsAutomaticRewardEventRewardTypeV2 = "send_highlighted_message"
	ChannelPointsAutomaticRewardEventRewardTypeV2RandomSubEmoteUnlock         ChannelPointsAutomaticRewardEventRewardTypeV2 = "random_sub_emote_unlock"
	ChannelPointsAutomaticRewardEventRewardTypeV2ChosenSubEmoteUnlock         ChannelPointsAutomaticRewardEventRewardTypeV2 = "chosen_sub_emote_unlock"
	ChannelPointsAutomaticRewardEventRewardTypeV2ChosenModifiedSubEmoteUnlock ChannelPointsAutomaticRewardEventRewardTypeV2 = "chosen_modified_sub_emote_unlock"
)

func (c ChannelPointsAutomaticRewardEventRewardTypeV2) String() string {
	return string(c)
}

type ChannelPointsAutomaticRewardEventRewardV2 struct {
	// The type of reward.
	Type ChannelPointsAutomaticRewardEventRewardTypeV2 `json:"type"`
	// Number of channel points used.
	ChannelPoints int `json:"channel_points"`
	// Optional. Emote associated with the reward.
	Emote *ChannelPointsAutomaticRewardEventRewardEmoteV2 `json:"emote"`
}

type ChannelPointsAutomaticRewardEventRewardEmoteV2 struct {
	// The emote ID.
	Id string `json:"id"`
	// The human-readable emote token.
	Name string `json:"name"`
}

type ChannelPointsAutomaticRewardEventRewardMessageV2 struct {
	// The chat message in plain text.
	Text string `json:"text"`
	// The ordered list of chat message fragments.
	Fragments []ChannelPointsAutomaticRewardEventRewardMessageFragmentV2 `json:"fragments,omitempty"`
}

type ChannelPointsAutomaticRewardEventRewardMessageFragmentTypeV2 string

const (
	ChannelPointsAutomaticRewardEventRewardMessageFragmentTypeV2Text  ChannelPointsAutomaticRewardEventRewardMessageFragmentTypeV2 = "text"
	ChannelPointsAutomaticRewardEventRewardMessageFragmentTypeV2Emote ChannelPointsAutomaticRewardEventRewardMessageFragmentTypeV2 = "emote"
)

type ChannelPointsAutomaticRewardEventRewardMessageFragmentV2 struct {
	// The message text in fragment.
	Text string `json:"text"`
	// The type of message fragment.
	Type ChannelPointsAutomaticRewardEventRewardMessageFragmentTypeV2 `json:"type"`
	// Optional. The metadata pertaining to the emote.
	Emote *ChannelPointsAutomaticRewardEventRewardMessageFragmentV2Emote `json:"emote,omitempty"`
}

type ChannelPointsAutomaticRewardEventRewardMessageFragmentV2Emote struct {
	// The ID that uniquely identifies this emote.
	Id string `json:"id"`
}

type ChannelPointsRewardEventMaxPerStream struct {
	// Is the setting enabled.
	IsEnabled bool `json:"is_enabled"`
	// The max per stream limit.
	Value int `json:"value"`
}

type ChannelPointsRewardEventMaxPerUserPerStream struct {
	// Is the setting enabled.
	IsEnabled bool `json:"is_enabled"`
	// The max per user per stream limit.
	Value int `json:"value"`
}

type ChannelPointsRewardEventGlobalCooldown struct {
	// Is the setting enabled.
	IsEnabled bool `json:"is_enabled"`
	// The global cooldown in seconds.
	Seconds int `json:"seconds"`
}

type ChannelPointsRewardEventImage struct {
	// URL for the image at 1x size.
	URL1x string `json:"url_1x"`
	// URL for the image at 2x size.
	URL2x string `json:"url_2x"`
	// URL for the image at 4x size.
	URL4x string `json:"url_4x"`
}

type ChannelSubscriptionMessageEventMessage struct {
	// The text of the resubscription chat message.
	Text string `json:"text"`
	// An array that includes the emote ID and start and end positions for where the emote appears in the text.
	Emotes []MessageEmote `json:"emotes,omitempty"`
}

type ChannelUnbanRequestResolveEventStatus string

const (
	ChannelUnbanRequestResolveEventStatusApproved ChannelUnbanRequestResolveEventStatus = "approved"
	ChannelUnbanRequestResolveEventStatusCanceled ChannelUnbanRequestResolveEventStatus = "canceled"
	ChannelUnbanRequestResolveEventStatusDenied   ChannelUnbanRequestResolveEventStatus = "denied"
)

func (c ChannelUnbanRequestResolveEventStatus) String() string {
	return string(c)
}

type StreamType string

const (
	StreamTypeLive       StreamType = "live"
	StreamTypePlaylist   StreamType = "playlist"
	StreamTypeWatchParty StreamType = "watch_party"
	StreamTypePremiere   StreamType = "premiere"
	StreamTypeRerun      StreamType = "rerun"
)

func (s StreamType) String() string {
	return string(s)
}

type CustomRewardRedemptionStatus string

type CharityAmount struct {
	// The monetary amount. The amount is specified in the currency's minor unit.
	// For example, the minor units for USD is cents, so if the amount is $5.50 USD, value is set to 550.
	Value int `json:"value"`
	// The number of decimal places used by the currency. For example, USD uses two decimal places.
	DecimalPlaces int `json:"decimal_places"`
	// The ISO-4217 three-letter currency code that identifies the type of currency in value.
	Currency string `json:"currency"`
}

type CharityDescription struct {
	// CharityName is the charity's name.
	CharityName string `json:"charity_name"`
	// CharityDescription is a description of the charity.
	CharityDescription string `json:"charity_description"`
	// CharityLogo is a URL to an image of the charity's logo.
	CharityLogo string `json:"charity_logo"`
	// CharityWebsite is a URL to the charity's website.
	CharityWebsite string `json:"charity_website"`
}

type GoalType string

const (
	GoalTypeFollow              GoalType = "follow"
	GoalTypeSubscription        GoalType = "subscription"
	GoalTypeSubscriptionCount   GoalType = "subscription_count"
	GoalTypeNewSubscription     GoalType = "new_subscription"
	GoalTypeNewSubscriptionCount GoalType = "new_subscription_count"
	GoalTypeNewBit              GoalType = "new_bit"
	GoalTypeNewCheerer          GoalType = "new_cheerer"
)

func (g GoalType) String() string {
	return string(g)
}

type HypeTrainContributionType string

const (
	HypeTrainContributionBits         HypeTrainContributionType = "bits"
	HypeTrainContributionSubscription HypeTrainContributionType = "subscription"
	HypeTrainContributionOther        HypeTrainContributionType = "other"
)

func (h HypeTrainContributionType) String() string {
	return string(h)
}

type HypeTrainContribution struct {
	// UserId is the ID of the user that made the contribution.
	UserId string `json:"user_id"`
	// UserLogin is the user's login name.
	UserLogin string `json:"user_login"`
	// UserName is the user's display name.
	UserName string `json:"user_name"`
	// Type is the contribution method used.
	Type HypeTrainContributionType `json:"type"`
	// Total is the total amount contributed.
	Total int `json:"total"`
}

type HypeTrainType string

const (
	HypeTrainTypeTreasure    HypeTrainType = "treasure"
	HypeTrainTypeGoldenKappa HypeTrainType = "golden_kappa"
	HypeTrainTypeRegular     HypeTrainType = "regular"
)

func (h HypeTrainType) String() string {
	return string(h)
}

type SharedChatParticipant struct {
	// BroadcasterUserId is the User ID of the participant channel.
	BroadcasterUserId string `json:"broadcaster_user_id"`
	// BroadcasterUserName is the display name of the participant channel.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// BroadcasterUserLogin is the user login of the participant channel.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
}

type GuestStarGuestState string

const (
	GuestStarGuestStateInvited   GuestStarGuestState = "invited"
	GuestStarGuestStateAccepted  GuestStarGuestState = "accepted"
	GuestStarGuestStateReady     GuestStarGuestState = "ready"
	GuestStarGuestStateBackstage GuestStarGuestState = "backstage"
	GuestStarGuestStateLive      GuestStarGuestState = "live"
	GuestStarGuestStateRemoved   GuestStarGuestState = "removed"
)

func (g GuestStarGuestState) String() string {
	return string(g)
}

type GuestStarGroupLayout string

const (
	GuestStarGroupLayoutTiled            GuestStarGroupLayout = "tiled"
	GuestStarGroupLayoutScreenshare      GuestStarGroupLayout = "screenshare"
	GuestStarGroupLayoutHorizontalTop    GuestStarGroupLayout = "horizontal_top"
	GuestStarGroupLayoutHorizontalBottom GuestStarGroupLayout = "horizontal_bottom"
	GuestStarGroupLayoutVerticalLeft     GuestStarGroupLayout = "vertical_left"
	GuestStarGroupLayoutVerticalRight    GuestStarGroupLayout = "vertical_right"
)

func (g GuestStarGroupLayout) String() string {
	return string(g)
}

type SuspiciousUserLowTrustStatus string

const (
	SuspiciousUserLowTrustStatusNone             SuspiciousUserLowTrustStatus = "none"
	SuspiciousUserLowTrustStatusActiveMonitoring SuspiciousUserLowTrustStatus = "active_monitoring"
	SuspiciousUserLowTrustStatusRestricted       SuspiciousUserLowTrustStatus = "restricted"
)

func (s SuspiciousUserLowTrustStatus) String() string {
	return string(s)
}

type SuspiciousUserType string

const (
	SuspiciousUserTypeManuallyAdded         SuspiciousUserType = "manually_added"
	SuspiciousUserTypeBanEvader             SuspiciousUserType = "ban_evader"
	SuspiciousUserTypeBannedInSharedChannel SuspiciousUserType = "banned_in_shared_channel"
)

type SuspiciousUserBanEvasionEvaluation string

const (
	SuspiciousUserBanEvasionEvaluationUnknown  SuspiciousUserBanEvasionEvaluation = "unknown"
	SuspiciousUserBanEvasionEvaluationPossible SuspiciousUserBanEvasionEvaluation = "possible"
	SuspiciousUserBanEvasionEvaluationLikely   SuspiciousUserBanEvasionEvaluation = "likely"
)

type UserMessageUpdateStatus string

const (
	UserMessageUpdateStatusApproved UserMessageUpdateStatus = "approved"
	UserMessageUpdateStatusDenied   UserMessageUpdateStatus = "denied"
	UserMessageUpdateStatusInvalid  UserMessageUpdateStatus = "invalid"
)

func (u UserMessageUpdateStatus) String() string {
	return string(u)
}

type SuspiciousUserMessageMessage struct {
	// MessageId is the UUID that identifies the message.
	MessageId string `json:"message_id"`
	// Text is the chat message in plain text.
	Text string `json:"text"`
	// Fragments is an ordered list of chat message fragments.
	Fragments []ChannelChatMessageEventMessageFragment `json:"fragments"`
}

type ExtensionBitsTransactionProduct struct {
	// Name is the product name.
	Name string `json:"name"`
	// Bits is the number of Bits the product cost.
	Bits int `json:"bits"`
	// Sku is the product's SKU.
	Sku string `json:"sku"`
}

type WhisperMessageBody struct {
	// Text is the body of the whisper message.
	Text string `json:"text"`
}

type CustomPowerUp struct {
	// Id is the unique ID for this Custom Power-up.
	Id string `json:"id"`
	// Title is the user-viewable name of this Custom Power-up.
	Title string `json:"title"`
	// Bits is the cost of the Custom Power-up to redeem.
	Bits int `json:"bits"`
	// Prompt is the creator-provided description for this Power-up.
	Prompt string `json:"prompt"`
}

type DropEntitlementGrantData struct {
	// OrganizationId is the ID of the organization that owns the game that has Drops enabled.
	OrganizationId string `json:"organization_id"`
	// CategoryId is the Twitch category ID of the game that was being played when this benefit was entitled.
	CategoryId string `json:"category_id"`
	// CategoryName is the category name.
	CategoryName string `json:"category_name"`
	// CampaignId is the campaign this entitlement is associated with.
	CampaignId string `json:"campaign_id"`
	// UserId is the Twitch user ID of the user who was granted the entitlement.
	UserId string `json:"user_id"`
	// UserName is the user display name of the user who was granted the entitlement.
	UserName string `json:"user_name"`
	// UserLogin is the user login of the user who was granted the entitlement.
	UserLogin string `json:"user_login"`
	// EntitlementId is the unique identifier of the entitlement.
	EntitlementId string `json:"entitlement_id"`
	// BenefitId is the identifier of the Benefit.
	BenefitId string `json:"benefit_id"`
	// CreatedAt is the UTC timestamp in ISO format when this entitlement was granted on Twitch.
	CreatedAt TimestampUTC `json:"created_at"`
}

const (
	CustomRewardRedemptionStatusUnknown     CustomRewardRedemptionStatus = "unknown"
	CustomRewardRedemptionStatusUnfulfilled CustomRewardRedemptionStatus = "unfulfilled"
	CustomRewardRedemptionStatusFulfilled   CustomRewardRedemptionStatus = "fulfilled"
	CustomRewardRedemptionStatusCanceled    CustomRewardRedemptionStatus = "canceled"
)
