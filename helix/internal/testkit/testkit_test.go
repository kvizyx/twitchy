package testkit

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

func TestFailingRoundTripperRejectsExternalHostBeforeDelegate(t *testing.T) {
	called := false
	delegate := roundTripperFunc(func(*http.Request) (*http.Response, error) {
		called = true
		return nil, nil
	})
	transport := NewFailingRoundTripper(delegate)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://reserved.external.test/helix", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = transport.RoundTrip(request)
	var externalErr *ExternalDialError
	if !errors.As(err, &externalErr) || externalErr.Host != "reserved.external.test" {
		t.Fatalf("RoundTrip() error = %v, want ExternalDialError for reserved host", err)
	}
	if called {
		t.Fatal("delegate received an external request")
	}
}

func TestFakeSleeperAdvancesFakeClockWithoutWallSleep(t *testing.T) {
	clock := NewFakeClock(time.Unix(1704067200, 0))
	sleeper := NewFakeSleeper(clock)
	start := time.Now()
	if err := sleeper.Sleep(context.Background(), 3*time.Second); err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
		t.Fatalf("fake sleep took wall time: %s", elapsed)
	}
	if got := clock.Now(); !got.Equal(time.Unix(1704067203, 0)) {
		t.Fatalf("fake clock = %s, want +3s", got)
	}
}

func TestFixtureLoaderReportsMalformedJSONPath(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "malformed.json"), []byte(`{"fixture":`), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := LoadJSON[map[string]string](root, "malformed.json")
	if err == nil || !strings.Contains(err.Error(), "malformed.json") || !strings.Contains(err.Error(), "decode JSON") {
		t.Fatalf("LoadJSON() error = %v, want fixture path and decode context", err)
	}
}

func TestManifestContractRejectsMalformedSuccessPayload(t *testing.T) {
	fixture := ContractFixture{
		Request:  ContractRequest{Headers: http.Header{"Authorization": {"Bearer test-token"}}},
		Response: ContractResponse{Status: http.StatusOK, Body: "not-json", Success: true},
		Want:     ContractExpectation{Attempts: 1},
	}
	operation := manifest.Operation{Method: http.MethodGet, Path: "/helix/test", Request: manifest.RequestSpec{Locations: map[string][]manifest.RequestField{"query": nil}}}
	transport := NewRecordingRoundTripper(ResponseFromFixture(fixture.Response))
	err := RunManifestContract(context.Background(), operation, fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		request, requestErr := http.NewRequest(http.MethodGet, "https://api.twitch.test/helix/test", nil)
		if requestErr != nil {
			return helix.ResponseMeta{}, requestErr
		}
		request.Header = fixture.Request.Headers.Clone()
		_, requestErr = transport.RoundTrip(request)
		if requestErr != nil {
			return helix.ResponseMeta{}, requestErr
		}
		return helix.ResponseMeta{}, nil
	})
	if err == nil || !strings.Contains(err.Error(), "decode") {
		t.Fatalf("RunManifestContract() error = %v, want malformed success decode failure", err)
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
