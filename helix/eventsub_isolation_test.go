package helix_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/kvizyx/twitchy/helix"
	internaljson "github.com/kvizyx/twitchy/internal/json"
)

func TestEventSubIsolation_clientConstructionDoesNotMutateSharedJSONOrEventSubState(t *testing.T) {
	beforeUnmarshal := reflect.ValueOf(internaljson.Unmarshal).Pointer()
	beforeTransport, err := json.Marshal(eventsub.WebhookTransport{Method: "webhook", Callback: "https://callback.test"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := helix.New(); err != nil {
		t.Fatal(err)
	}
	afterUnmarshal := reflect.ValueOf(internaljson.Unmarshal).Pointer()
	afterTransport, err := json.Marshal(eventsub.WebhookTransport{Method: "webhook", Callback: "https://callback.test"})
	if err != nil {
		t.Fatal(err)
	}
	if beforeUnmarshal != afterUnmarshal || string(beforeTransport) != string(afterTransport) {
		t.Fatalf("shared state changed: unmarshal %x/%x transport %s/%s", beforeUnmarshal, afterUnmarshal, beforeTransport, afterTransport)
	}
}
