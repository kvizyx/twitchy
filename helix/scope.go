package helix

// AuthorizationScope is a Twitch scope that gives your app permission to perform the specified action on the user’s behalf.
type AuthorizationScope string

const (
	// ScopeAnalyticsReadExtensions view analytics data for Twitch extensions owned by the authenticated account.
	ScopeAnalyticsReadExtensions AuthorizationScope = "analytics:read:extensions"
	// ScopeAnalyticsReadGames view analytics data for games owned by the authenticated account.
	ScopeAnalyticsReadGames AuthorizationScope = "analytics:read:games"
	// ScopeBitsRead view bits information for a channel.
	ScopeBitsRead AuthorizationScope = "bits:read"
	// ScopeChannelBot join channel’s chatroom as bot user and perform chat-related actions.
	ScopeChannelBot AuthorizationScope = "channel:bot"
	// ScopeChannelManageAds manage ads schedule on a channel.
	ScopeChannelManageAds AuthorizationScope = "channel:manage:ads"
	// ScopeChannelReadAds read the ads schedule and details on your channel.
	ScopeChannelReadAds AuthorizationScope = "channel:read:ads"
	// ScopeChannelManageBroadcast manage a channel’s broadcast configuration and stream tags.
	ScopeChannelManageBroadcast AuthorizationScope = "channel:manage:broadcast"
	// ScopeChannelReadCharity read charity campaign details and donations on your channel.
	ScopeChannelReadCharity AuthorizationScope = "channel:read:charity"
	// ScopeChannelManageClips create and manage Clips for a channel.
	ScopeChannelManageClips AuthorizationScope = "channel:manage:clips"
	// ScopeChannelEditCommercial run commercials on a channel.
	ScopeChannelEditCommercial AuthorizationScope = "channel:edit:commercial"
	// ScopeChannelReadEditors view users with editor role for a channel.
	ScopeChannelReadEditors AuthorizationScope = "channel:read:editors"
	// ScopeChannelManageExtensions manage a channel’s Extension configuration.
	ScopeChannelManageExtensions AuthorizationScope = "channel:manage:extensions"
	// ScopeChannelReadGoals view Creator Goals for a channel.
	ScopeChannelReadGoals AuthorizationScope = "channel:read:goals"
	// ScopeChannelReadGuestStar read Guest Star details for your channel.
	ScopeChannelReadGuestStar AuthorizationScope = "channel:read:guest_star"
	// ScopeChannelManageGuestStar manage Guest Star for your channel.
	ScopeChannelManageGuestStar AuthorizationScope = "channel:manage:guest_star"
	// ScopeChannelReadHypeTrain view Hype Train information for a channel.
	ScopeChannelReadHypeTrain AuthorizationScope = "channel:read:hype_train"
	// ScopeChannelManageModerators add or remove the moderator role from users in your channel.
	ScopeChannelManageModerators AuthorizationScope = "channel:manage:moderators"
	// ScopeChannelReadPolls view a channel’s polls.
	ScopeChannelReadPolls AuthorizationScope = "channel:read:polls"
	// ScopeChannelManagePolls manage a channel’s polls.
	ScopeChannelManagePolls AuthorizationScope = "channel:manage:polls"
	// ScopeChannelReadPredictions view a channel’s predictions.
	ScopeChannelReadPredictions AuthorizationScope = "channel:read:predictions"
	// ScopeChannelManagePredictions manage a channel’s predictions.
	ScopeChannelManagePredictions AuthorizationScope = "channel:manage:predictions"
	// ScopeChannelManageRaids manage channel raids.
	ScopeChannelManageRaids AuthorizationScope = "channel:manage:raids"
	// ScopeChannelReadRedemptions view Channel Points custom rewards and redemptions.
	ScopeChannelReadRedemptions AuthorizationScope = "channel:read:redemptions"
	// ScopeChannelManageRedemptions manage Channel Points custom rewards and redemptions.
	ScopeChannelManageRedemptions AuthorizationScope = "channel:manage:redemptions"
	// ScopeChannelManageSchedule manage a channel’s stream schedule.
	ScopeChannelManageSchedule AuthorizationScope = "channel:manage:schedule"
	// ScopeChannelReadStreamKey view a channel’s stream key.
	ScopeChannelReadStreamKey AuthorizationScope = "channel:read:stream_key"
	// ScopeChannelReadSubscriptions view subscribers of a channel.
	ScopeChannelReadSubscriptions AuthorizationScope = "channel:read:subscriptions"
	// ScopeChannelManageVideos manage channel’s videos.
	ScopeChannelManageVideos AuthorizationScope = "channel:manage:videos"
	// ScopeChannelReadVIPs view list of VIPs in a channel.
	ScopeChannelReadVIPs AuthorizationScope = "channel:read:vips"
	// ScopeChannelManageVIPs add or remove VIP role in a channel.
	ScopeChannelManageVIPs AuthorizationScope = "channel:manage:vips"
	// ScopeChannelModerate perform moderation actions in a channel.
	ScopeChannelModerate AuthorizationScope = "channel:moderate"
	// ScopeClipsEdit manage Clips for a channel.
	ScopeClipsEdit AuthorizationScope = "clips:edit"
	// ScopeEditorManageClips manage Clips as an editor.
	ScopeEditorManageClips AuthorizationScope = "editor:manage:clips"
	// ScopeModerationRead view moderation data like moderators, bans, etc.
	ScopeModerationRead AuthorizationScope = "moderation:read"
	// ScopeModeratorManageAnnouncements send announcements as a moderator.
	ScopeModeratorManageAnnouncements AuthorizationScope = "moderator:manage:announcements"
	// ScopeModeratorManageAutoMod manage messages held by AutoMod.
	ScopeModeratorManageAutoMod AuthorizationScope = "moderator:manage:automod"
	// ScopeModeratorReadAutoModSettings view AutoMod settings.
	ScopeModeratorReadAutoModSettings AuthorizationScope = "moderator:read:automod_settings"
	// ScopeModeratorManageAutoModSettings manage AutoMod settings.
	ScopeModeratorManageAutoModSettings AuthorizationScope = "moderator:manage:automod_settings"
	// ScopeModeratorReadBannedUsers read list of banned users.
	ScopeModeratorReadBannedUsers AuthorizationScope = "moderator:read:banned_users"
	// ScopeModeratorManageBannedUsers ban/unban users.
	ScopeModeratorManageBannedUsers AuthorizationScope = "moderator:manage:banned_users"
	// ScopeModeratorReadBlockedTerms view blocked terms list.
	ScopeModeratorReadBlockedTerms AuthorizationScope = "moderator:read:blocked_terms"
	// ScopeModeratorManageBlockedTerms manage blocked terms list.
	ScopeModeratorManageBlockedTerms AuthorizationScope = "moderator:manage:blocked_terms"
	// ScopeModeratorReadChatMessages read deleted chat messages.
	ScopeModeratorReadChatMessages AuthorizationScope = "moderator:read:chat_messages"
	// ScopeModeratorManageChatMessages delete chat messages.
	ScopeModeratorManageChatMessages AuthorizationScope = "moderator:manage:chat_messages"
	// ScopeModeratorReadChatSettings view chat settings.
	ScopeModeratorReadChatSettings AuthorizationScope = "moderator:read:chat_settings"
	// ScopeModeratorManageChatSettings manage chat settings.
	ScopeModeratorManageChatSettings AuthorizationScope = "moderator:manage:chat_settings"
	// ScopeModeratorReadChatters view chatters in chat room.
	ScopeModeratorReadChatters AuthorizationScope = "moderator:read:chatters"
	// ScopeModeratorReadFollowers view followers of a broadcaster.
	ScopeModeratorReadFollowers AuthorizationScope = "moderator:read:followers"
	// ScopeModeratorReadGuestStar read Guest Star data as a moderator.
	ScopeModeratorReadGuestStar AuthorizationScope = "moderator:read:guest_star"
	// ScopeModeratorManageGuestStar manage Guest Star as a moderator.
	ScopeModeratorManageGuestStar AuthorizationScope = "moderator:manage:guest_star"
	// ScopeModeratorReadModerators read list of moderators.
	ScopeModeratorReadModerators AuthorizationScope = "moderator:read:moderators"
	// ScopeModeratorReadShieldMode view Shield Mode status.
	ScopeModeratorReadShieldMode AuthorizationScope = "moderator:read:shield_mode"
	// ScopeModeratorManageShieldMode manage Shield Mode status.
	ScopeModeratorManageShieldMode AuthorizationScope = "moderator:manage:shield_mode"
	// ScopeModeratorReadShoutouts view shoutouts.
	ScopeModeratorReadShoutouts AuthorizationScope = "moderator:read:shoutouts"
	// ScopeModeratorManageShoutouts manage shoutouts.
	ScopeModeratorManageShoutouts AuthorizationScope = "moderator:manage:shoutouts"
	// ScopeModeratorReadSuspiciousUsers read suspicious users and messages.
	ScopeModeratorReadSuspiciousUsers AuthorizationScope = "moderator:read:suspicious_users"
	// ScopeModeratorManageSuspiciousUsers manage suspicious users and messages.
	ScopeModeratorManageSuspiciousUsers AuthorizationScope = "moderator:manage:suspicious_users"
	// ScopeModeratorReadUnbanRequests view unban requests.
	ScopeModeratorReadUnbanRequests AuthorizationScope = "moderator:read:unban_requests"
	// ScopeModeratorManageUnbanRequests manage unban requests.
	ScopeModeratorManageUnbanRequests AuthorizationScope = "moderator:manage:unban_requests"
	// ScopeModeratorReadVIPs view list of VIPs (as moderator).
	ScopeModeratorReadVIPs AuthorizationScope = "moderator:read:vips"
	// ScopeModeratorReadWarnings view warnings issued.
	ScopeModeratorReadWarnings AuthorizationScope = "moderator:read:warnings"
	// ScopeModeratorManageWarnings send warnings in chat.
	ScopeModeratorManageWarnings AuthorizationScope = "moderator:manage:warnings"
	// ScopeUserBot join a chatroom and appear as bot (as user).
	ScopeUserBot AuthorizationScope = "user:bot"
	// ScopeUserEdit manage user profile.
	ScopeUserEdit AuthorizationScope = "user:edit"
	// ScopeUserEditBroadcast view/edit broadcasting configuration.
	ScopeUserEditBroadcast AuthorizationScope = "user:edit:broadcast"
	// ScopeUserReadBlockedUsers view block list of a user.
	ScopeUserReadBlockedUsers AuthorizationScope = "user:read:blocked_users"
	// ScopeUserManageBlockedUsers manage block list of a user.
	ScopeUserManageBlockedUsers AuthorizationScope = "user:manage:blocked_users"
	// ScopeUserReadBroadcast view broadcasting configuration.
	ScopeUserReadBroadcast AuthorizationScope = "user:read:broadcast"
	// ScopeUserReadChat receive chat messages and notifications.
	ScopeUserReadChat AuthorizationScope = "user:read:chat"
	// ScopeUserManageChatColor update name color in chat.
	ScopeUserManageChatColor AuthorizationScope = "user:manage:chat_color"
	// ScopeUserReadEmail view user’s email.
	ScopeUserReadEmail AuthorizationScope = "user:read:email"
	// ScopeUserReadEmotes view user emotes.
	ScopeUserReadEmotes AuthorizationScope = "user:read:emotes"
	// ScopeUserReadFollows view followed channels.
	ScopeUserReadFollows AuthorizationScope = "user:read:follows"
	// ScopeUserReadModeratedChannels view moderated channels.
	ScopeUserReadModeratedChannels AuthorizationScope = "user:read:moderated_channels"
	// ScopeUserReadSubscriptions view user’s channel subscriptions.
	ScopeUserReadSubscriptions AuthorizationScope = "user:read:subscriptions"
	// ScopeUserReadWhispers receive whispers.
	ScopeUserReadWhispers AuthorizationScope = "user:read:whispers"
	// ScopeUserManageWhispers manage whispers.
	ScopeUserManageWhispers AuthorizationScope = "user:manage:whispers"
	// ScopeUserWriteChat send chat messages to a chatroom.
	ScopeUserWriteChat AuthorizationScope = "user:write:chat"

	// ScopeChatEdit send chat messages to a chatroom using an IRC connection.
	ScopeChatEdit AuthorizationScope = "chat:edit"
	// ScopeChatRead view chat messages sent in a chatroom using an IRC connection.
	ScopeChatRead AuthorizationScope = "chat:read"
	// ScopeWhispersRead read whispers.
	ScopeWhispersRead AuthorizationScope = "whispers:read"
)
