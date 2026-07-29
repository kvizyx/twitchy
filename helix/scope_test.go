package helix_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kvizyx/twitchy/helix"
)

type scopeDefinition struct {
	name  string
	value helix.AuthorizationScope
}

var requiredScopeValues = []string{
	"analytics:read:extensions",
	"analytics:read:games",
	"bits:read",
	"channel:bot",
	"channel:manage:ads",
	"channel:read:ads",
	"channel:manage:broadcast",
	"channel:read:charity",
	"channel:manage:clips",
	"channel:edit:commercial",
	"channel:read:editors",
	"channel:manage:extensions",
	"channel:read:goals",
	"channel:read:guest_star",
	"channel:manage:guest_star",
	"channel:read:hype_train",
	"channel:manage:moderators",
	"channel:read:polls",
	"channel:manage:polls",
	"channel:read:predictions",
	"channel:manage:predictions",
	"channel:manage:raids",
	"channel:read:redemptions",
	"channel:manage:redemptions",
	"channel:manage:schedule",
	"channel:read:stream_key",
	"channel:read:subscriptions",
	"channel:manage:videos",
	"channel:read:vips",
	"channel:manage:vips",
	"channel:moderate",
	"clips:edit",
	"editor:manage:clips",
	"moderation:read",
	"moderator:manage:announcements",
	"moderator:manage:automod",
	"moderator:read:automod_settings",
	"moderator:manage:automod_settings",
	"moderator:read:banned_users",
	"moderator:manage:banned_users",
	"moderator:read:blocked_terms",
	"moderator:manage:blocked_terms",
	"moderator:read:chat_messages",
	"moderator:manage:chat_messages",
	"moderator:read:chat_settings",
	"moderator:manage:chat_settings",
	"moderator:read:chatters",
	"moderator:read:followers",
	"moderator:read:guest_star",
	"moderator:manage:guest_star",
	"moderator:read:moderators",
	"moderator:read:shield_mode",
	"moderator:manage:shield_mode",
	"moderator:read:shoutouts",
	"moderator:manage:shoutouts",
	"moderator:read:suspicious_users",
	"moderator:manage:suspicious_users",
	"moderator:read:unban_requests",
	"moderator:manage:unban_requests",
	"moderator:read:vips",
	"moderator:read:warnings",
	"moderator:manage:warnings",
	"user:bot",
	"user:edit",
	"user:edit:broadcast",
	"user:read:blocked_users",
	"user:manage:blocked_users",
	"user:read:broadcast",
	"user:read:chat",
	"user:manage:chat_color",
	"user:read:email",
	"user:read:emotes",
	"user:read:follows",
	"user:read:moderated_channels",
	"user:read:subscriptions",
	"user:read:whispers",
	"user:manage:whispers",
	"user:write:chat",
	"chat:edit",
	"chat:read",
	"whispers:read",
}

var scopeConstantNames = map[string]string{
	"channel:manage:clips":              "ScopeChannelManageClips",
	"editor:manage:clips":               "ScopeEditorManageClips",
	"moderator:manage:suspicious_users": "ScopeModeratorManageSuspiciousUsers",
	"whispers:read":                     "ScopeWhispersRead",
}

var declaredScopes = []scopeDefinition{
	{"ScopeAnalyticsReadExtensions", helix.ScopeAnalyticsReadExtensions},
	{"ScopeAnalyticsReadGames", helix.ScopeAnalyticsReadGames},
	{"ScopeBitsRead", helix.ScopeBitsRead},
	{"ScopeChannelBot", helix.ScopeChannelBot},
	{"ScopeChannelManageAds", helix.ScopeChannelManageAds},
	{"ScopeChannelReadAds", helix.ScopeChannelReadAds},
	{"ScopeChannelManageBroadcast", helix.ScopeChannelManageBroadcast},
	{"ScopeChannelReadCharity", helix.ScopeChannelReadCharity},
	{"ScopeChannelManageClips", helix.ScopeChannelManageClips},
	{"ScopeChannelEditCommercial", helix.ScopeChannelEditCommercial},
	{"ScopeChannelReadEditors", helix.ScopeChannelReadEditors},
	{"ScopeChannelManageExtensions", helix.ScopeChannelManageExtensions},
	{"ScopeChannelReadGoals", helix.ScopeChannelReadGoals},
	{"ScopeChannelReadGuestStar", helix.ScopeChannelReadGuestStar},
	{"ScopeChannelManageGuestStar", helix.ScopeChannelManageGuestStar},
	{"ScopeChannelReadHypeTrain", helix.ScopeChannelReadHypeTrain},
	{"ScopeChannelManageModerators", helix.ScopeChannelManageModerators},
	{"ScopeChannelReadPolls", helix.ScopeChannelReadPolls},
	{"ScopeChannelManagePolls", helix.ScopeChannelManagePolls},
	{"ScopeChannelReadPredictions", helix.ScopeChannelReadPredictions},
	{"ScopeChannelManagePredictions", helix.ScopeChannelManagePredictions},
	{"ScopeChannelManageRaids", helix.ScopeChannelManageRaids},
	{"ScopeChannelReadRedemptions", helix.ScopeChannelReadRedemptions},
	{"ScopeChannelManageRedemptions", helix.ScopeChannelManageRedemptions},
	{"ScopeChannelManageSchedule", helix.ScopeChannelManageSchedule},
	{"ScopeChannelReadStreamKey", helix.ScopeChannelReadStreamKey},
	{"ScopeChannelReadSubscriptions", helix.ScopeChannelReadSubscriptions},
	{"ScopeChannelManageVideos", helix.ScopeChannelManageVideos},
	{"ScopeChannelReadVIPs", helix.ScopeChannelReadVIPs},
	{"ScopeChannelManageVIPs", helix.ScopeChannelManageVIPs},
	{"ScopeChannelModerate", helix.ScopeChannelModerate},
	{"ScopeClipsEdit", helix.ScopeClipsEdit},
	{"ScopeEditorManageClips", helix.ScopeEditorManageClips},
	{"ScopeModerationRead", helix.ScopeModerationRead},
	{"ScopeModeratorManageAnnouncements", helix.ScopeModeratorManageAnnouncements},
	{"ScopeModeratorManageAutoMod", helix.ScopeModeratorManageAutoMod},
	{"ScopeModeratorReadAutoModSettings", helix.ScopeModeratorReadAutoModSettings},
	{"ScopeModeratorManageAutoModSettings", helix.ScopeModeratorManageAutoModSettings},
	{"ScopeModeratorReadBannedUsers", helix.ScopeModeratorReadBannedUsers},
	{"ScopeModeratorManageBannedUsers", helix.ScopeModeratorManageBannedUsers},
	{"ScopeModeratorReadBlockedTerms", helix.ScopeModeratorReadBlockedTerms},
	{"ScopeModeratorManageBlockedTerms", helix.ScopeModeratorManageBlockedTerms},
	{"ScopeModeratorReadChatMessages", helix.ScopeModeratorReadChatMessages},
	{"ScopeModeratorManageChatMessages", helix.ScopeModeratorManageChatMessages},
	{"ScopeModeratorReadChatSettings", helix.ScopeModeratorReadChatSettings},
	{"ScopeModeratorManageChatSettings", helix.ScopeModeratorManageChatSettings},
	{"ScopeModeratorReadChatters", helix.ScopeModeratorReadChatters},
	{"ScopeModeratorReadFollowers", helix.ScopeModeratorReadFollowers},
	{"ScopeModeratorReadGuestStar", helix.ScopeModeratorReadGuestStar},
	{"ScopeModeratorManageGuestStar", helix.ScopeModeratorManageGuestStar},
	{"ScopeModeratorReadModerators", helix.ScopeModeratorReadModerators},
	{"ScopeModeratorReadShieldMode", helix.ScopeModeratorReadShieldMode},
	{"ScopeModeratorManageShieldMode", helix.ScopeModeratorManageShieldMode},
	{"ScopeModeratorReadShoutouts", helix.ScopeModeratorReadShoutouts},
	{"ScopeModeratorManageShoutouts", helix.ScopeModeratorManageShoutouts},
	{"ScopeModeratorReadSuspiciousUsers", helix.ScopeModeratorReadSuspiciousUsers},
	{"ScopeModeratorManageSuspiciousUsers", helix.ScopeModeratorManageSuspiciousUsers},
	{"ScopeModeratorReadUnbanRequests", helix.ScopeModeratorReadUnbanRequests},
	{"ScopeModeratorManageUnbanRequests", helix.ScopeModeratorManageUnbanRequests},
	{"ScopeModeratorReadVIPs", helix.ScopeModeratorReadVIPs},
	{"ScopeModeratorReadWarnings", helix.ScopeModeratorReadWarnings},
	{"ScopeModeratorManageWarnings", helix.ScopeModeratorManageWarnings},
	{"ScopeUserBot", helix.ScopeUserBot},
	{"ScopeUserEdit", helix.ScopeUserEdit},
	{"ScopeUserEditBroadcast", helix.ScopeUserEditBroadcast},
	{"ScopeUserReadBlockedUsers", helix.ScopeUserReadBlockedUsers},
	{"ScopeUserManageBlockedUsers", helix.ScopeUserManageBlockedUsers},
	{"ScopeUserReadBroadcast", helix.ScopeUserReadBroadcast},
	{"ScopeUserReadChat", helix.ScopeUserReadChat},
	{"ScopeUserManageChatColor", helix.ScopeUserManageChatColor},
	{"ScopeUserReadEmail", helix.ScopeUserReadEmail},
	{"ScopeUserReadEmotes", helix.ScopeUserReadEmotes},
	{"ScopeUserReadFollows", helix.ScopeUserReadFollows},
	{"ScopeUserReadModeratedChannels", helix.ScopeUserReadModeratedChannels},
	{"ScopeUserReadSubscriptions", helix.ScopeUserReadSubscriptions},
	{"ScopeUserReadWhispers", helix.ScopeUserReadWhispers},
	{"ScopeUserManageWhispers", helix.ScopeUserManageWhispers},
	{"ScopeUserWriteChat", helix.ScopeUserWriteChat},
	{"ScopeChatEdit", helix.ScopeChatEdit},
	{"ScopeChatRead", helix.ScopeChatRead},
	{"ScopeWhispersRead", helix.ScopeWhispersRead},
}

func TestScopes(t *testing.T) {
	// Given the required scope table and the exported scope constants.

	// When converting both tables to sets.
	actual := make([]string, len(declaredScopes))
	for index, scope := range declaredScopes {
		actual[index] = string(scope.value)
	}
	actualSet, actualDuplicate := uniqueValues(actual)
	requiredSet, requiredDuplicate := uniqueValues(requiredScopeValues)

	// Then every descriptor value is represented exactly once and no extra value exists.
	if actualDuplicate != "" {
		t.Fatalf("scope value %q is duplicated in exported constants", actualDuplicate)
	}
	if requiredDuplicate != "" {
		t.Fatalf("scope value %q is duplicated in descriptor", requiredDuplicate)
	}
	if len(actualSet) != len(requiredSet) {
		t.Fatalf("scope count = %d, want descriptor count %d", len(actualSet), len(requiredSet))
	}
	for value := range requiredSet {
		if _, ok := actualSet[value]; !ok {
			t.Errorf("descriptor scope %q is missing from exported constants", value)
		}
	}
	for value := range actualSet {
		if _, ok := requiredSet[value]; !ok {
			t.Errorf("exported scope %q is absent from descriptor", value)
		}
	}

	for value, name := range scopeConstantNames {
		found := false
		for _, scope := range declaredScopes {
			if scope.name == name {
				if string(scope.value) != value {
					t.Errorf("%s = %q, want %q", name, scope.value, value)
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("descriptor constant %s is not exported", name)
		}
	}
}

func TestLegacyScopeCompatibility(t *testing.T) {
	// Given an external consumer fixture importing every pre-existing scope symbol.
	directory := t.TempDir()
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	moduleRoot := filepath.Dir(root)
	legacy := make([]string, 0, len(declaredScopes)-4)
	for _, scope := range declaredScopes {
		switch scope.name {
		case "ScopeChannelManageClips", "ScopeEditorManageClips", "ScopeModeratorManageSuspiciousUsers", "ScopeWhispersRead":
			continue
		default:
			legacy = append(legacy, "helix."+scope.name)
		}
	}
	goMod := fmt.Sprintf("module fixture\n\ngo 1.24\n\nrequire github.com/kvizyx/twitchy v0.0.0\n\nreplace github.com/kvizyx/twitchy => %s\n", moduleRoot)
	source := "package fixture\n\nimport \"github.com/kvizyx/twitchy/helix\"\n\nvar scopes = []helix.AuthorizationScope{\n\t" + strings.Join(legacy, ",\n\t") + ",\n}\n"
	if err := os.WriteFile(filepath.Join(directory, "go.mod"), []byte(goMod), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(directory, "fixture.go"), []byte(source), 0o600); err != nil {
		t.Fatalf("write fixture.go: %v", err)
	}

	// When the fixture is compiled as an external module.
	command := exec.Command("go", "build", ".")
	command.Dir = directory
	command.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOPROXY=off", "GOSUMDB=off", "GOWORK=off")
	output, err := command.CombinedOutput()

	// Then every pre-existing exported symbol remains source-compatible.
	if err != nil {
		t.Fatalf("legacy scope fixture failed to compile: %v\n%s", err, output)
	}
}

func uniqueValues(values []string) (map[string]struct{}, string) {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		if _, exists := set[value]; exists {
			return set, value
		}
		set[value] = struct{}{}
	}
	return set, ""
}
