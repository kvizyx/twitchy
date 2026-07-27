package helix_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestUsersAuthorization_isExperimentalOnlyAtCompileBoundary(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	moduleRoot := filepath.Dir(filepath.Dir(sourceFile))
	validSource := `package boundary

import (
    "context"
    "github.com/kvizyx/twitchy/helix"
)

func useExperimental(client *helix.Client) {
    _, _ = client.Experimental.Users.GetAuthorizationByUser(context.Background(), helix.GetAuthorizationByUserRequest{})
}
`
	if output, err := compileUsersBoundary(t, moduleRoot, validSource); err != nil {
		t.Fatalf("experimental users source failed to compile: %v\n%s", err, output)
	}

	invalidSource := `package boundary

import (
    "context"
    "github.com/kvizyx/twitchy/helix"
)

func useStable(client *helix.Client) {
    _, _ = client.Users.GetAuthorizationByUser(context.Background(), helix.GetAuthorizationByUserRequest{})
}
`
	output, err := compileUsersBoundary(t, moduleRoot, invalidSource)
	if err == nil {
		t.Fatal("stable users authorization source unexpectedly compiled")
	}
	if !strings.Contains(string(output), "GetAuthorizationByUser undefined") {
		t.Fatalf("compile error = %s", output)
	}
}

func compileUsersBoundary(t *testing.T, moduleRoot, source string) ([]byte, error) {
	t.Helper()
	moduleDir := t.TempDir()
	goMod := "module users-boundary\n\ngo 1.24\n\nrequire github.com/kvizyx/twitchy v0.0.0\n\nreplace github.com/kvizyx/twitchy => " + moduleRoot + "\n"
	if err := os.WriteFile(filepath.Join(moduleDir, "go.mod"), []byte(goMod), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "boundary.go"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	command := exec.Command("go", "test", "./...")
	command.Dir = moduleDir
	command.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOPROXY=off", "GOSUMDB=off")
	return command.CombinedOutput()
}
