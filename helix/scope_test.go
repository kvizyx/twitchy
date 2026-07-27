package helix_test

import (
	"encoding/json"
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

type scopeDescriptor struct {
	Scopes struct {
		RequiredValues   []string          `json:"required_values"`
		NewConstantNames map[string]string `json:"new_constant_names"`
	} `json:"scopes"`
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
	// Given the copied core/OAuth descriptor and the exported scope table.
	bytes, err := os.ReadFile(filepath.Join("internal", "manifest", "core-oauth-descriptor.json"))
	if err != nil {
		t.Fatal(err)
	}
	var descriptor scopeDescriptor
	if err := json.Unmarshal(bytes, &descriptor); err != nil {
		t.Fatal(err)
	}

	// When converting both tables to sets.
	actual := make([]string, len(declaredScopes))
	for index, scope := range declaredScopes {
		actual[index] = string(scope.value)
	}
	actualSet, actualDuplicate := uniqueValues(actual)
	requiredSet, requiredDuplicate := uniqueValues(descriptor.Scopes.RequiredValues)

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

	for value, name := range descriptor.Scopes.NewConstantNames {
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
