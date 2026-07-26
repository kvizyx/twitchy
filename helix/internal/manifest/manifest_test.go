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
		"sources/receipt.json":        "637790fac19c7cfebf9aaa94d735b34d4bcf203050594832cd5dc66f3a4f1e37",
		"sources/reference.html":      "3c5ed44ffac743027b96897a07cca4e4b1d4225dc8fa941a2aec6acc7795e926",
		"sources/guide.html":          "d4053e1645c82622053e1280abd4a11960caa177a9453b8ce3d639accee9b2e3",
		"sources/scopes.html":         "123a411f410edf057c2eb7aa142b7f08e4335c13964164a2a0a71396a615d55e",
		"sources/oauth.html":          "d20843e754c6a22b295d634df99718e6f62d23210e4c6fbdd8cf6e750f609a59",
		"sources/refresh.html":        "a0f61ea5226dadb4c31b1f4b915540c246ae05df2d48ad601fcb1e3eeddb9eac",
		"sources/validate.html":       "26dbaa4684586580ff0281725505e4491960e9745933be3152d66f4048b75f2a",
		"sources/revoke.html":         "2264dd8fe2b79107906a484dc39f2436db8c2b33a9b4fa767aa4ced4108cd3cf",
		"sources/authentication.html": "f361ed78d1e54a392042f554753439a6916cf9e213023b6a5fd1f00eb11750e5",
		"expected-operations.json":    "16f361ee282409ab4b5a94ab4a00a585a7a57a76f70855e7c73d93ef8b51615f",
		"public-descriptor.json":      "eb230fb36e9f98f06a71a34cc4d2b05c9a39098c51bc032c8f0b6f919c459f73",
		"core-oauth-descriptor.json":  "5830f46d453ec3490eb323982fc3a9075c2a16bf8177410dd6e3bc350342ebed",
		"independent-verifier.mjs":    "2d2153fa7906d60c2e47a72de0f7f6c12279dfa50377200e5f767d5f8ff0dfb6",
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
