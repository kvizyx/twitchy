package helix_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGuestStarStableSurfaceExcludesExperimentalOperations(t *testing.T) {
	// Given an external package that attempts to call Guest Star through a stable service.
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	moduleRoot := filepath.Dir(filepath.Dir(sourceFile))
	moduleDir := t.TempDir()
	goMod := "module guest-star-boundary\n\ngo 1.24\n\nrequire github.com/kvizyx/twitchy v0.0.0\n\nreplace github.com/kvizyx/twitchy => " + moduleRoot + "\n"
	mainGo := `package boundary

import (
	"context"
	"github.com/kvizyx/twitchy/helix"
)

func useStableGuestStar(client *helix.Client) {
	_, _ = client.Channels.GetChannelGuestStarSettings(context.Background(), helix.GetChannelGuestStarSettingsRequest{})
}
`
	if err := os.WriteFile(filepath.Join(moduleDir, "go.mod"), []byte(goMod), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "boundary.go"), []byte(mainGo), 0o600); err != nil {
		t.Fatal(err)
	}

	// When the external package is compiled.
	command := exec.Command("go", "test", "./...")
	command.Dir = moduleDir
	command.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOPROXY=off", "GOSUMDB=off")
	output, err := command.CombinedOutput()

	// Then the stable namespace rejects the experimental Guest Star symbol.
	if err == nil {
		t.Fatal("stable Guest Star call unexpectedly compiled")
	}
	if !strings.Contains(string(output), "GetChannelGuestStarSettings undefined") {
		t.Fatalf("compile error = %s", output)
	}
}
