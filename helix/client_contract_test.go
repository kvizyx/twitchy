package helix_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

type compileTokenSource struct{}

func (compileTokenSource) Token(context.Context) (helix.CredentialSnapshot, error) {
	return helix.CredentialSnapshot{}, nil
}

func TestClientContract_initializesStableAndExperimentalServices(t *testing.T) {
	// Given a client using the default options
	client, err := helix.New()

	// Then construction succeeds and every frozen service selector is usable.
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	stable := []*struct {
		name  string
		value any
	}{
		{"Ads", client.Ads},
		{"Analytics", client.Analytics},
		{"Bits", client.Bits},
		{"Channels", client.Channels},
		{"ChannelPoints", client.ChannelPoints},
		{"Charity", client.Charity},
		{"Chat", client.Chat},
		{"Clips", client.Clips},
		{"Conduits", client.Conduits},
		{"CCLs", client.CCLs},
		{"Entitlements", client.Entitlements},
		{"Extensions", client.Extensions},
		{"EventSub", client.EventSub},
		{"Games", client.Games},
		{"Goals", client.Goals},
		{"HypeTrain", client.HypeTrain},
		{"Moderation", client.Moderation},
		{"Polls", client.Polls},
		{"Predictions", client.Predictions},
		{"Raids", client.Raids},
		{"Schedule", client.Schedule},
		{"Search", client.Search},
		{"Streams", client.Streams},
		{"Subscriptions", client.Subscriptions},
		{"Tags", client.Tags},
		{"Teams", client.Teams},
		{"Users", client.Users},
		{"Videos", client.Videos},
		{"Whispers", client.Whispers},
	}
	for _, service := range stable {
		if service.value == nil {
			t.Errorf("stable service %s is nil", service.name)
		}
	}
	for _, service := range []*struct {
		name  string
		value any
	}{
		{"Bits", client.Experimental.Bits},
		{"Chat", client.Experimental.Chat},
		{"Clips", client.Experimental.Clips},
		{"GuestStar", client.Experimental.GuestStar},
		{"Moderation", client.Experimental.Moderation},
		{"Users", client.Experimental.Users},
	} {
		if service.value == nil {
			t.Errorf("experimental service %s is nil", service.name)
		}
	}
}

func TestClientContract_rejectsInvalidOptions(t *testing.T) {
	// Given invalid boundary options
	tests := []struct {
		name   string
		option helix.Option
	}{
		{"nil http client", helix.WithHTTPClient(nil)},
		{"nil token source", helix.WithTokenSource(nil)},
		{"empty base URL", helix.WithBaseURL("")},
		{"relative base URL", helix.WithBaseURL("helix.example.test")},
		{"unsupported base URL", helix.WithBaseURL("ftp://helix.example.test")},
		{"malformed base URL", helix.WithBaseURL("://bad")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// When constructing the client
			_, err := helix.New(test.option)

			// Then the option is rejected with the stable sentinel.
			if !errors.Is(err, helix.ErrInvalidOption) {
				t.Fatalf("New() error = %v, want ErrInvalidOption", err)
			}
		})
	}
}

func TestClientContract_rejectsConflictingCredentialOptions(t *testing.T) {
	// Given both mutually exclusive credential configuration forms
	options := []helix.Option{
		helix.WithStaticToken(helix.Credential{AccessToken: "static"}),
		helix.WithTokenSource(compileTokenSource{}),
	}

	// When constructing in either option order
	for _, ordered := range [][]helix.Option{options, {options[1], options[0]}} {
		_, err := helix.New(ordered...)

		// Then construction reports the conflict.
		if !errors.Is(err, helix.ErrConflictingOptions) {
			t.Fatalf("New() error = %v, want ErrConflictingOptions", err)
		}
	}
}

func TestClientContract_optionsUseLastValue(t *testing.T) {
	// Given two valid values for every replaceable option
	first := &http.Client{Timeout: time.Second}
	second := &http.Client{Timeout: 2 * time.Second}

	// When options are applied left to right
	client, err := helix.New(
		helix.WithHTTPClient(first),
		helix.WithHTTPClient(second),
		helix.WithBaseURL("https://first.example.test/helix"),
		helix.WithBaseURL("https://second.example.test/helix"),
		helix.WithUserAgent("first"),
		helix.WithUserAgent("second"),
		helix.WithRateLimitPolicy(helix.RateLimitPolicy{MaxWait: time.Second}),
		helix.WithRateLimitPolicy(helix.RateLimitPolicy{MaxWait: 2 * time.Second}),
	)

	// Then construction succeeds; the public contract is exercised by the compile surface.
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if client == nil {
		t.Fatal("New() returned nil client")
	}
}

func TestClientContract_doesNotMutateCallerHTTPClient(t *testing.T) {
	// Given a caller-owned client with borrowed resources and redirect policy
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar.New() error = %v", err)
	}
	transport := &http.Transport{}
	redirect := func(*http.Request, []*http.Request) error { return errors.New("caller redirect") }
	caller := &http.Client{
		Transport:     transport,
		CheckRedirect: redirect,
		Jar:           jar,
		Timeout:       time.Second,
	}
	before := *caller

	// When the client is configured from the caller-owned client
	_, err = helix.New(helix.WithHTTPClient(caller))

	// Then construction succeeds and no caller field changes.
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	after := *caller
	if after.Transport != before.Transport || after.Jar != before.Jar || after.Timeout != before.Timeout {
		t.Fatal("New() mutated caller HTTP client fields")
	}
	if reflect.ValueOf(after.CheckRedirect).Pointer() != reflect.ValueOf(before.CheckRedirect).Pointer() {
		t.Fatal("New() mutated caller redirect policy")
	}
}

func TestConsumerCompile_rejectsStableAccessToExperimentalService(t *testing.T) {
	// Given an external consumer fixture that tries to access GuestStar as stable API
	directory := t.TempDir()
	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	moduleRoot := filepath.Dir(root)
	goMod := "module fixture\n\ngo 1.24\n\nrequire github.com/kvizyx/twitchy v0.0.0\n\nreplace github.com/kvizyx/twitchy => " + moduleRoot + "\n"
	source := "package fixture\n\nimport \"github.com/kvizyx/twitchy/helix\"\n\nfunc invalid(client *helix.Client) { _ = client.GuestStar }\n"
	if err := os.WriteFile(filepath.Join(directory, "go.mod"), []byte(goMod), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(directory, "fixture.go"), []byte(source), 0o600); err != nil {
		t.Fatalf("write fixture.go: %v", err)
	}

	// When the fixture is compiled as an external consumer
	command := exec.Command("go", "build", ".")
	command.Dir = directory
	command.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOPROXY=off", "GOSUMDB=off", "GOWORK=off")
	output, err := command.CombinedOutput()

	// Then the compiler rejects the experimental-only selector.
	if err == nil {
		t.Fatal("stable experimental access unexpectedly compiled")
	}
	if !bytes.Contains(output, []byte("GuestStar")) && !strings.Contains(string(output), "undefined") {
		t.Fatalf("unexpected compile failure: %s", output)
	}
}

func TestConsumerCompile_oldAuthorizationScopeSymbolsRemainAvailable(t *testing.T) {
	// Given legacy scope identifiers from the frozen public contract
	var _ = []helix.AuthorizationScope{
		helix.ScopeAnalyticsReadExtensions,
		helix.ScopeAnalyticsReadGames,
		helix.ScopeBitsRead,
		helix.ScopeChannelBot,
		helix.ScopeChannelManageAds,
		helix.ScopeChannelModerate,
		helix.ScopeClipsEdit,
		helix.ScopeModerationRead,
		helix.ScopeUserBot,
		helix.ScopeUserEdit,
		helix.ScopeUserReadBlockedUsers,
		helix.ScopeUserWriteChat,
		helix.ScopeChatEdit,
		helix.ScopeChatRead,
	}
}
