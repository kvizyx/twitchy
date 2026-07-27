package helix

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestStringEnum_preservesUnknownWireValues(t *testing.T) {
	var value StringEnum
	if err := json.Unmarshal([]byte(`"future_value"`), &value); err != nil {
		t.Fatalf("unmarshal enum: %v", err)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal enum: %v", err)
	}
	if string(encoded) != `"future_value"` {
		t.Fatalf("enum round trip: got %s", encoded)
	}
}

func TestResponseMeta_copiesHeadersAndPaginationExposesCursor(t *testing.T) {
	meta := newResponseMetaAt(200, http.Header{"X-Request-Id": {"abc"}}, 1, time.Now().Add(-time.Second))
	headers := meta.Header()
	headers.Set("X-Request-ID", "mutated")
	if meta.RequestID() != "abc" || meta.Header().Get("X-Request-ID") != "abc" {
		t.Fatalf("metadata was not immutable: %#v", meta.Header())
	}

	pagination := Pagination{cursor: "next"}
	if pagination.Cursor() != "next" {
		t.Fatalf("cursor: got %q", pagination.Cursor())
	}
}
