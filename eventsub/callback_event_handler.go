package eventsub

import (
	"errors"
	"fmt"

	"github.com/kvizyx/twitchy/internal/json"
)

// ErrUndefinedEventType indicates that eventsub sent an event that is not defined in the library, and the user has not
// set a callback handler for undefined events.
var ErrUndefinedEventType = errors.New("undefined event type")

// RawEvent is an event with raw payload that is not depends on the underlying transport.
type RawEvent struct {
	Subscription []byte
	Event        []byte
}

// runEventCallback runs an event callback from the callback store if it is set by the user, or skips this run without
// error otherwise.
//
// If the provided EventType is not defined in the library, ErrUndefinedEventType will be returned or OnUndefinedEvent
// will be triggered, if it is set by the user.
func (c *callback[Metadata]) runEventCallback(
	eventType EventType,
	eventVersion string,
	rawEvent RawEvent,
	metadata Metadata,
) error {
	event := rawEvent.Event

	switch eventType {
	case EventTypeAutomodMessageHold:
		switch eventVersion {
		case "1":
			return runEventCallbackHandler(c.onAutomodMessageHold, event, metadata)
		case "2":
			return runEventCallbackHandler(c.onAutomodMessageHoldV2, event, metadata)
		}
	case EventTypeAutomodMessageUpdate:
		switch eventVersion {
		case "1":
			return runEventCallbackHandler(c.onAutomodMessageUpdate, event, metadata)
		case "2":
			return runEventCallbackHandler(c.onAutomodMessageUpdateV2, event, metadata)
		}
	case EventTypeAutomodSettingsUpdate:
		return runEventCallbackHandler(c.onAutomodSettingsUpdate, event, metadata)
	case EventTypeAutomodTermsUpdate:
		return runEventCallbackHandler(c.onAutomodTermsUpdate, event, metadata)
	case EventTypeChannelBitsUse:
		return runEventCallbackHandler(c.onChannelBitsUse, event, metadata)
	case EventTypeChannelUpdate:
		return runEventCallbackHandler(c.onChannelUpdate, event, metadata)
	case EventTypeChannelFollow:
		return runEventCallbackHandler(c.onChannelFollow, event, metadata)
	case EventTypeChannelAdBreakBegin:
		return runEventCallbackHandler(c.onChannelAdBreakBegin, event, metadata)
	case EventTypeChannelChatClear:
		return runEventCallbackHandler(c.onChannelChatClear, event, metadata)
	case EventTypeChannelChatClearUserMessages:
		return runEventCallbackHandler(c.onChannelChatClearUserMessages, event, metadata)
	case EventTypeChannelChatMessage:
		return runEventCallbackHandler(c.onChannelChatMessage, event, metadata)
	case EventTypeConduitShardDisabled:
		return runEventCallbackHandler(c.onConduitShardDisabled, event, metadata)
	case EventTypeChannelBan:
		return runEventCallbackHandler(c.onChannelBan, event, metadata)
	case EventTypeChannelChatNotification:
		return runEventCallbackHandler(c.onChannelChatNotification, event, metadata)
	case EventTypeChannelUnban:
		return runEventCallbackHandler(c.onChannelUnban, event, metadata)
	case EventTypeChannelModeratorAdd:
		return runEventCallbackHandler(c.onChannelModeratorAdd, event, metadata)
	case EventTypeChannelModeratorRemove:
		return runEventCallbackHandler(c.onChannelModeratorRemove, event, metadata)
	case EventTypeChannelPollBegin:
		return runEventCallbackHandler(c.onChannelPollBegin, event, metadata)
	case EventTypeChannelPollProgress:
		return runEventCallbackHandler(c.onChannelPollProgress, event, metadata)
	case EventTypeChannelPollEnd:
		return runEventCallbackHandler(c.onChannelPollEnd, event, metadata)
	case EventTypeChannelPredictionBegin:
		return runEventCallbackHandler(c.onChannelPredictionBegin, event, metadata)
	case EventTypeChannelPredictionProgress:
		return runEventCallbackHandler(c.onChannelPredictionProgress, event, metadata)
	case EventTypeChannelPredictionLock:
		return runEventCallbackHandler(c.onChannelPredictionLock, event, metadata)
	case EventTypeChannelPredictionEnd:
		return runEventCallbackHandler(c.onChannelPredictionEnd, event, metadata)
	case EventTypeChannelRaid:
		return runEventCallbackHandler(c.onChannelRaid, event, metadata)
	case EventTypeChannelPointsCustomRewardRedemptionAdd:
		return runEventCallbackHandler(c.onChannelPointsCustomRewardRedemptionAdd, event, metadata)
	case EventTypeChannelPointsCustomRewardRedemptionUpdate:
		return runEventCallbackHandler(c.onChannelPointsCustomRewardRedemptionUpdate, event, metadata)
	case EventTypeChannelPointsAutomaticRewardRedemptionAdd:
		switch eventVersion {
		case "1":
			return runEventCallbackHandler(c.onChannelPointsAutomaticRewardRedemptionAdd, event, metadata)
		case "2":
			return runEventCallbackHandler(c.onChannelPointsAutomaticRewardRedemptionAddV2, event, metadata)
		}
	case EventTypeUserAuthorizationRevoke:
		return runEventCallbackHandler(c.onUserAuthorizationRevoke, event, metadata)
	case EventTypeChannelPointsRewardAdd:
		return runEventCallbackHandler(c.onChannelPointsCustomRewardAdd, event, metadata)
	case EventTypeChannelPointsRewardUpdate:
		return runEventCallbackHandler(c.onChannelPointsCustomRewardUpdate, event, metadata)
	case EventTypeChannelPointsRewardRemove:
		return runEventCallbackHandler(c.onChannelPointsCustomRewardRemove, event, metadata)
	case EventTypeStreamOffline:
		return runEventCallbackHandler(c.onStreamOffline, event, metadata)
	case EventTypeStreamOnline:
		return runEventCallbackHandler(c.onStreamOnline, event, metadata)
	case EventTypeChannelSubscribe:
		return runEventCallbackHandler(c.onChannelSubscribe, event, metadata)
	case EventTypeChannelSubscriptionEnd:
		return runEventCallbackHandler(c.onChannelSubscriptionEnd, event, metadata)
	case EventTypeChannelSubscriptionMessage:
		return runEventCallbackHandler(c.onChannelSubscriptionMessage, event, metadata)
	case EventTypeChannelSubscriptionGift:
		return runEventCallbackHandler(c.onChannelSubscriptionGift, event, metadata)
	case EventTypeChannelUnbanRequestCreate:
		return runEventCallbackHandler(c.onChannelUnbanRequestCreate, event, metadata)
	case EventTypeChannelUnbanRequestResolve:
		return runEventCallbackHandler(c.onChannelUnbanRequestResolve, event, metadata)
	case EventTypeUserUpdate:
		return runEventCallbackHandler(c.onUserUpdate, event, metadata)
	case EventTypeChannelVipAdd:
		return runEventCallbackHandler(c.onChannelVipAdd, event, metadata)
	case EventTypeChannelVipRemove:
		return runEventCallbackHandler(c.onChannelVipRemove, event, metadata)
	case EventTypeChannelMessageDelete:
		return runEventCallbackHandler(c.onChannelChatMessageDelete, event, metadata)
	case EventTypeChannelModerate:
		switch eventVersion {
		case "2":
			return runEventCallbackHandler(c.onChannelModerateV2, event, metadata)
		}
	case EventTypeChannelChatSettingsUpdate:
		return runEventCallbackHandler(c.onChannelChatSettingsUpdate, event, metadata)
	case EventTypeChannelChatUserMessageHold:
		return runEventCallbackHandler(c.onChannelChatUserMessageHold, event, metadata)
	case EventTypeChannelChatUserMessageUpdate:
		return runEventCallbackHandler(c.onChannelChatUserMessageUpdate, event, metadata)
	case EventTypeChannelCustomPowerUpRedemptionAdd:
		return runEventCallbackHandler(c.onChannelCustomPowerUpRedemptionAdd, event, metadata)
	case EventTypeChannelSharedChatSessionBegin:
		return runEventCallbackHandler(c.onChannelSharedChatSessionBegin, event, metadata)
	case EventTypeChannelSharedChatSessionUpdate:
		return runEventCallbackHandler(c.onChannelSharedChatSessionUpdate, event, metadata)
	case EventTypeChannelSharedChatSessionEnd:
		return runEventCallbackHandler(c.onChannelSharedChatSessionEnd, event, metadata)
	case EventTypeChannelSuspiciousUserMessage:
		return runEventCallbackHandler(c.onChannelSuspiciousUserMessage, event, metadata)
	case EventTypeChannelSuspiciousUserUpdate:
		return runEventCallbackHandler(c.onChannelSuspiciousUserUpdate, event, metadata)
	case EventTypeChannelWarningAcknowledge:
		return runEventCallbackHandler(c.onChannelWarningAcknowledge, event, metadata)
	case EventTypeChannelWarningSend:
		return runEventCallbackHandler(c.onChannelWarningSend, event, metadata)
	case EventTypeChannelCheer:
		return runEventCallbackHandler(c.onChannelCheer, event, metadata)
	case EventTypeCharityCampaignDonate:
		return runEventCallbackHandler(c.onCharityCampaignDonate, event, metadata)
	case EventTypeCharityCampaignStart:
		return runEventCallbackHandler(c.onCharityCampaignStart, event, metadata)
	case EventTypeCharityCampaignProgress:
		return runEventCallbackHandler(c.onCharityCampaignProgress, event, metadata)
	case EventTypeCharityCampaignStop:
		return runEventCallbackHandler(c.onCharityCampaignStop, event, metadata)
	case EventTypeDropEntitlementGrant:
		return runEventCallbackHandler(c.onDropEntitlementGrant, event, metadata)
	case EventTypeExtensionBitsTransactionCreate:
		return runEventCallbackHandler(c.onExtensionBitsTransactionCreate, event, metadata)
	case EventTypeChannelGoalBegin:
		return runEventCallbackHandler(c.onChannelGoalBegin, event, metadata)
	case EventTypeChannelGoalProgress:
		return runEventCallbackHandler(c.onChannelGoalProgress, event, metadata)
	case EventTypeChannelGoalEnd:
		return runEventCallbackHandler(c.onChannelGoalEnd, event, metadata)
	case EventTypeHypeTrainBegin:
		return runEventCallbackHandler(c.onHypeTrainBegin, event, metadata)
	case EventTypeHypeTrainProgress:
		return runEventCallbackHandler(c.onHypeTrainProgress, event, metadata)
	case EventTypeHypeTrainEnd:
		return runEventCallbackHandler(c.onHypeTrainEnd, event, metadata)
	case EventTypeShieldModeBegin:
		return runEventCallbackHandler(c.onShieldModeBegin, event, metadata)
	case EventTypeShieldModeEnd:
		return runEventCallbackHandler(c.onShieldModeEnd, event, metadata)
	case EventTypeShoutoutCreate:
		return runEventCallbackHandler(c.onShoutoutCreate, event, metadata)
	case EventTypeShoutoutReceive:
		return runEventCallbackHandler(c.onShoutoutReceive, event, metadata)
	case EventTypeUserAuthorizationGrant:
		return runEventCallbackHandler(c.onUserAuthorizationGrant, event, metadata)
	case EventTypeUserWhisperMessage:
		return runEventCallbackHandler(c.onUserWhisperMessage, event, metadata)
	case EventTypeGuestStarSessionBegin:
		return runEventCallbackHandler(c.onGuestStarSessionBegin, event, metadata)
	case EventTypeGuestStarSessionEnd:
		return runEventCallbackHandler(c.onGuestStarSessionEnd, event, metadata)
	case EventTypeGuestStarGuestUpdate:
		return runEventCallbackHandler(c.onGuestStarGuestUpdate, event, metadata)
	case EventTypeGuestStarSettingsUpdate:
		return runEventCallbackHandler(c.onGuestStarSettingsUpdate, event, metadata)
	default:
		if c.onUndefinedEvent != nil {
			go c.onUndefinedEvent(rawEvent, metadata)
			return nil
		}

		return ErrUndefinedEventType
	}

	return nil
}

// runEventCallbackHandler parses provided payload as JSON data to generic event and runs handler in separate go-routine
// with this event and metadata if handler is defined by user, otherwise returns without error.
func runEventCallbackHandler[Event, Metadata any](
	handler Handler[Event, Metadata],
	eventPayload []byte,
	eventMetadata Metadata,
) error {
	var event Event

	if handler != nil {
		if err := json.Unmarshal(eventPayload, &event); err != nil {
			return fmt.Errorf("unmarshal event payload: %w", err)
		}

		go handler(event, eventMetadata)
	}

	return nil
}
