package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestFrozenSourceDigest(t *testing.T) {
	want := map[string]string{
		"expected-operations.json":   "16f361ee282409ab4b5a94ab4a00a585a7a57a76f70855e7c73d93ef8b51615f",
		"public-descriptor.json":     "eb230fb36e9f98f06a71a34cc4d2b05c9a39098c51bc032c8f0b6f919c459f73",
		"core-oauth-descriptor.json": "5830f46d453ec3490eb323982fc3a9075c2a16bf8177410dd6e3bc350342ebed",
		"independent-verifier.mjs":   "486e735ea8e1f89c56a06607f3b809f794c85d4b9c8118dd27c1fd04612adcb7",
	}
	for name, expected := range want {
		path := filepath.Join(".", name)
		info, err := os.Lstat(path)
		if err != nil {
			t.Fatalf("stat %s: %v", name, err)
		}
		if !info.Mode().IsRegular() {
			t.Fatalf("%s is not a regular file", name)
		}
		bytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		actual := sha256.Sum256(bytes)
		if got := hex.EncodeToString(actual[:]); got != expected {
			t.Errorf("%s: got %s, want %s", name, got, expected)
		}
	}
}
