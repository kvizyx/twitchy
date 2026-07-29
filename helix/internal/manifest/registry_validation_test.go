package manifest

import (
	"strings"
	"testing"
)

func TestRegistryValidation(t *testing.T) {
	if err := validateOperations(Operations()); err != nil {
		t.Fatal(err)
	}
}

func TestRegistryFrozenSnapshot(t *testing.T) {
	operations := Operations()
	if len(operations) != 149 {
		t.Fatalf("got %d operations, want 149", len(operations))
	}
	if got := len(expectedGroups); got != 30 {
		t.Fatalf("got %d groups, want 30", got)
	}
	labels := map[Stability]int{}
	for _, operation := range operations {
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
		mutated := Operations()
		mutated[1].Anchor = mutated[0].Anchor
		if _, err := buildOperationIndex(mutated); err == nil || !strings.Contains(err.Error(), "duplicate anchor") {
			t.Fatalf("got %v, want duplicate anchor error", err)
		}
	})
	t.Run("missing auth", func(t *testing.T) {
		mutated := Operations()
		mutated[0].TokenClasses = nil
		if err := validateOperation(mutated[0]); err == nil || !strings.Contains(err.Error(), "token_classes is required") {
			t.Fatalf("got %v, want token_classes error", err)
		}
	})
	t.Run("148 rows", func(t *testing.T) {
		mutated := Operations()[:148]
		if err := validateOperations(mutated); err == nil || !strings.Contains(err.Error(), "got 148 operations, want 149") {
			t.Fatalf("got %v, want 148-row error", err)
		}
	})
}
