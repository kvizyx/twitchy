package helix

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestTimestamp_roundTripsRFC3339Nano(t *testing.T) {
	want := time.Date(2026, time.July, 27, 12, 34, 56, 123456789, time.FixedZone("UTC+2", 2*60*60))
	encoded, err := json.Marshal(Timestamp{Time: want})
	if err != nil {
		t.Fatalf("marshal timestamp: %v", err)
	}
	var got Timestamp
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatalf("unmarshal timestamp: %v", err)
	}
	if !got.Time.Equal(want) {
		t.Fatalf("timestamp: got %s, want %s", got.Time, want)
	}
}

func TestTimestampUTC_normalizesLocation(t *testing.T) {
	var timestamp TimestampUTC
	if err := json.Unmarshal([]byte(`"2026-07-27T12:34:56.123456789+02:00"`), &timestamp); err != nil {
		t.Fatalf("unmarshal UTC timestamp: %v", err)
	}
	if timestamp.Location() != time.UTC {
		t.Fatalf("location: got %s, want UTC", timestamp.Location())
	}
}

func TestTimestamp_rejectsInvalidValues(t *testing.T) {
	for _, input := range []string{"null", `"not-a-timestamp"`, `"2026-07-27"`} {
		var timestamp Timestamp
		if err := json.Unmarshal([]byte(input), &timestamp); err == nil {
			t.Fatalf("input %s unexpectedly decoded", input)
		}
	}
	if _, err := time.Parse(time.RFC3339Nano, strings.TrimSpace("not-a-timestamp")); err == nil {
		t.Fatal("invalid fixture unexpectedly parsed")
	}
}
