package manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSnapshot20260726(t *testing.T) {
	manifest, err := LoadManifest()
	if err != nil {
		t.Fatal(err)
	}
	if err := manifest.Validate(); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Operations) != 149 {
		t.Fatalf("got %d operations, want 149", len(manifest.Operations))
	}
	if got := len(expectedGroups); got != 30 {
		t.Fatalf("got %d groups, want 30", got)
	}
	labels := map[Stability]int{}
	for _, operation := range manifest.Operations {
		labels[operation.Stability]++
		expectedReplayable := (operation.Method == "GET" || operation.Method == "HEAD") && operation.Request.BodyReconstructible
		if operation.Replay.Replayable != expectedReplayable {
			t.Errorf("%s: replayable=%t, want %t", operation.Anchor, operation.Replay.Replayable, expectedReplayable)
		}
		for _, status := range operation.Response.Status {
			if status == 429 && operation.Replay.BucketWaitable {
				t.Errorf("%s: endpoint-specific 429 marked bucket-waitable", operation.Anchor)
			}
		}
	}
	if labels[StabilityStable] != 127 || labels[StabilityNew] != 10 || labels[StabilityBeta] != 12 {
		t.Fatalf("stability partition: %#v", labels)
	}

	t.Run("duplicate anchor", func(t *testing.T) {
		root := copyManifest(t)
		path := filepath.Join(root, "groups", "ads.json")
		group := decodeGroup(t, path)
		group.Operations[1].Anchor = group.Operations[0].Anchor
		group.Operations[1].OperationID = group.Operations[0].OperationID
		writeGroup(t, path, group)
		_, err := LoadManifestFrom(root)
		if err == nil || !strings.Contains(err.Error(), "duplicate anchor") {
			t.Fatalf("got %v, want duplicate anchor error", err)
		}
	})
	t.Run("missing auth", func(t *testing.T) {
		root := copyManifest(t)
		path := filepath.Join(root, "groups", "ads.json")
		group := decodeGroup(t, path)
		group.Operations[0].TokenClasses = nil
		writeGroup(t, path, group)
		_, err := LoadManifestFrom(root)
		if err == nil || !strings.Contains(err.Error(), "token_classes is required") {
			t.Fatalf("got %v, want token_classes error", err)
		}
	})
	t.Run("148 rows", func(t *testing.T) {
		root := copyManifest(t)
		path := filepath.Join(root, "groups", "ads.json")
		group := decodeGroup(t, path)
		group.Operations = group.Operations[:len(group.Operations)-1]
		writeGroup(t, path, group)
		_, err := LoadManifestFrom(root)
		if err == nil || !strings.Contains(err.Error(), "got 148 operations, want 149") {
			t.Fatalf("got %v, want 148-row error", err)
		}
	})
}

func copyManifest(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "groups"), 0o755); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir("groups")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		bytes, err := os.ReadFile(filepath.Join("groups", entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "groups", entry.Name()), bytes, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func decodeGroup(t *testing.T, path string) GroupFile {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var group GroupFile
	if err := json.Unmarshal(bytes, &group); err != nil {
		t.Fatal(err)
	}
	return group
}

func writeGroup(t *testing.T, path string, group GroupFile) {
	t.Helper()
	bytes, err := json.MarshalIndent(group, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(bytes, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}
