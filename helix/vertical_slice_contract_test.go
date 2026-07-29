package helix_test

import (
	"encoding/json"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

func TestVerticalSlice_unknownEnumRoundTrips(t *testing.T) {
	// Given a forward-compatible named enum value unknown to this client.
	value := struct {
		State helix.StringEnum `json:"state"`
	}{State: helix.StringEnum("future-state")}

	// When it is encoded and decoded through the standard JSON boundary.
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		State helix.StringEnum `json:"state"`
	}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}

	// Then the unknown wire value survives unchanged.
	if decoded.State != value.State {
		t.Fatalf("state = %q, want %q", decoded.State, value.State)
	}
}

func TestVerticalSlice_manifestRowsMapOneToOne(t *testing.T) {
	// Given the frozen manifest and public implementation descriptors.
	games := manifestOperation(t, "get-games")
	powerUp := manifestOperation(t, "get-custom-power-up")

	// Then each row points to exactly one public symbol and test surface.
	if games.Implementation.Selector != "Client.Games" || games.Implementation.Method != "GetGames" || games.Implementation.RequestType != "GetGamesRequest" || games.Implementation.DataType != "GetGamesData" {
		t.Fatalf("get-games implementation = %#v", games.Implementation)
	}
	if powerUp.Implementation.Selector != "Client.Experimental.Bits" || powerUp.Implementation.Method != "GetCustomPowerUp" || powerUp.Implementation.RequestType != "GetCustomPowerUpRequest" || powerUp.Implementation.DataType != "GetCustomPowerUpData" {
		t.Fatalf("get-custom-power-up implementation = %#v", powerUp.Implementation)
	}
}

func manifestOperation(t *testing.T, anchor string) manifest.Operation {
	t.Helper()
	operation, err := manifest.OperationByAnchor(anchor)
	if err != nil {
		t.Fatal(err)
	}
	return operation
}

func urlValues(items ...string) map[string][]string {
	values := make(map[string][]string, len(items)/2)
	for index := 0; index < len(items); index += 2 {
		values[items[index]] = append(values[items[index]], items[index+1])
	}
	return values
}
