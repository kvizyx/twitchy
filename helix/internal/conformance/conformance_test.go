package conformance

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/manifest"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

const canonicalSummary = "149 operations, 30 groups, 127 stable, 10 NEW, 12 BETA, 0 missing, 0 extra, 0 unclassified, 0 duplicate mappings"

type publicDescriptor struct {
	Operations []publicOperation `json:"operations"`
}

type publicOperation struct {
	Anchor         string  `json:"anchor"`
	Selector       string  `json:"selector"`
	ServiceType    string  `json:"service_type"`
	Method         string  `json:"method"`
	Signature      string  `json:"signature"`
	PagerSignature *string `json:"pager_signature"`
	RequestType    string  `json:"request_type"`
	DataType       string  `json:"data_type"`
	Stability      string  `json:"stability"`
}

func TestManifestConformance(t *testing.T) {
	loaded, err := manifest.LoadManifestFrom("../manifest")
	if err != nil {
		t.Fatal(err)
	}
	descriptor := loadDescriptor(t)
	if err := verifyManifestRows(loaded, descriptor); err != nil {
		t.Fatal(err)
	}
	if err := verifyServiceSurface(loaded); err != nil {
		t.Fatal(err)
	}

	for _, operation := range loaded.Operations {
		t.Run(operation.Implementation.TestIDs[0], func(t *testing.T) {
			executeHappyContract(t, operation)
		})
		t.Run(operation.Implementation.TestIDs[1], func(t *testing.T) {
			executeNegativeContract(t, operation)
		})
		if operation.Replay.Replayable {
			t.Run("replay/"+operation.Anchor, func(t *testing.T) {
				executeReplayContract(t, operation)
			})
		}
		if operation.Pagination.Shape != "none" {
			t.Run("pagination/"+operation.Anchor, func(t *testing.T) {
				executePagerContract(t, operation)
			})
		}
	}

	fmt.Fprintln(os.Stdout, canonicalSummary)
}

func loadDescriptor(t *testing.T) publicDescriptor {
	t.Helper()
	data, err := os.ReadFile("../manifest/public-descriptor.json")
	if err != nil {
		t.Fatal(err)
	}
	var descriptor publicDescriptor
	if err := json.Unmarshal(data, &descriptor); err != nil {
		t.Fatal(err)
	}
	return descriptor
}

func verifyManifestRows(loaded manifest.Manifest, descriptor publicDescriptor) error {
	if len(loaded.Operations) != 149 || len(descriptor.Operations) != 149 {
		return fmt.Errorf("operation count: manifest=%d descriptor=%d", len(loaded.Operations), len(descriptor.Operations))
	}
	byAnchor := make(map[string]manifest.Operation, len(loaded.Operations))
	for _, operation := range loaded.Operations {
		if _, exists := byAnchor[operation.Anchor]; exists {
			return fmt.Errorf("duplicate operation mapping %q", operation.Anchor)
		}
		byAnchor[operation.Anchor] = operation
		if len(operation.Implementation.TestIDs) != 2 {
			return fmt.Errorf("%s: want one happy and one negative test ID", operation.Anchor)
		}
		if operation.Implementation.TestIDs[0] != "TestManifestConformance/happy/"+operation.Anchor || operation.Implementation.TestIDs[1] != "TestManifestConformance/negative/"+operation.Anchor {
			return fmt.Errorf("%s: invalid test IDs %q", operation.Anchor, operation.Implementation.TestIDs)
		}
		if operation.Implementation.Stability != operation.Stability {
			return fmt.Errorf("%s: implementation stability %q does not match row %q", operation.Anchor, operation.Implementation.Stability, operation.Stability)
		}
		if operation.Implementation.Selector == "" || operation.Implementation.ServiceType == "" || operation.Implementation.Method == "" || operation.Implementation.RequestType == "" || operation.Implementation.DataType == "" {
			return fmt.Errorf("%s: incomplete implementation mapping", operation.Anchor)
		}
		if operation.Stability == manifest.StabilityStable && strings.HasPrefix(operation.Implementation.Selector, "Client.Experimental.") {
			return fmt.Errorf("%s: stable operation is experimental-only", operation.Anchor)
		}
		if operation.Stability != manifest.StabilityStable && !strings.HasPrefix(operation.Implementation.Selector, "Client.Experimental.") {
			return fmt.Errorf("%s: experimental operation is reachable on stable service", operation.Anchor)
		}
		if operation.Replay.Replayable != (operation.Method == http.MethodGet || operation.Method == http.MethodHead) || operation.Request.BodyReconstructible != operation.Replay.Replayable {
			return fmt.Errorf("%s: unsafe replay classification", operation.Anchor)
		}
		if operation.Response.Format != "json" && operation.Response.Format != "text" && operation.Response.Format != "unknown" {
			return fmt.Errorf("%s: unclassified response format %q", operation.Anchor, operation.Response.Format)
		}
		if operation.Pagination.Shape == "" || operation.Pagination.CursorParameter == "" {
			return fmt.Errorf("%s: unclassified pagination", operation.Anchor)
		}
	}
	if len(byAnchor) != len(descriptor.Operations) {
		return fmt.Errorf("descriptor mappings: got %d manifest anchors", len(byAnchor))
	}
	seenMethods := make(map[string]string, len(loaded.Operations))
	for _, expected := range descriptor.Operations {
		operation, ok := byAnchor[expected.Anchor]
		if !ok {
			return fmt.Errorf("descriptor operation %q is missing", expected.Anchor)
		}
		implementation := operation.Implementation
		if implementation.Selector != expected.Selector || implementation.ServiceType != expected.ServiceType || implementation.Method != expected.Method || implementation.Signature != expected.Signature || implementation.RequestType != expected.RequestType || implementation.DataType != expected.DataType || string(implementation.Stability) != expected.Stability || !samePagerSignature(implementation.PagerSignature, expected.PagerSignature) {
			return fmt.Errorf("%s: implementation does not match public descriptor", expected.Anchor)
		}
		mapping := implementation.ServiceType + "." + implementation.Method
		if previous, exists := seenMethods[mapping]; exists {
			return fmt.Errorf("duplicate mappings %q and %q", previous, expected.Anchor)
		}
		seenMethods[mapping] = expected.Anchor
	}
	return nil
}

func samePagerSignature(got, want *string) bool {
	if got == nil || want == nil {
		return got == nil && want == nil
	}
	return *got == *want
}

func responseMeta(value reflect.Value) (helix.ResponseMeta, error) {
	if !value.IsValid() || value.IsNil() {
		return helix.ResponseMeta{}, errors.New("conformance: operation returned nil response")
	}
	field := value.Elem().FieldByName("Meta")
	if !field.IsValid() {
		return helix.ResponseMeta{}, errors.New("conformance: response has no metadata")
	}
	meta, ok := field.Interface().(helix.ResponseMeta)
	if !ok {
		return helix.ResponseMeta{}, errors.New("conformance: response metadata has wrong type")
	}
	return meta, nil
}

func operationRequest(operation manifest.Operation, method reflect.Value) reflect.Value {
	return populatedValue(method.Type().In(1), operation)
}

func requestFixture(operation manifest.Operation, transport *testkit.RecordingRoundTripper, response testkit.ContractResponse, attempts int) (testkit.ContractFixture, error) {
	records := transport.Requests()
	if len(records) == 0 {
		return testkit.ContractFixture{}, errors.New("conformance: no loopback request recorded")
	}
	request := records[len(records)-1]
	parsed, err := url.Parse("http://127.0.0.1" + request.Path)
	if err != nil {
		return testkit.ContractFixture{}, err
	}
	if parsed.Path != operation.Path || request.Method != operation.Method {
		return testkit.ContractFixture{}, fmt.Errorf("request got %s %s, want %s %s", request.Method, parsed.Path, operation.Method, operation.Path)
	}
	return testkit.ContractFixture{
		Request:  testkit.ContractRequest{Path: operation.Path, Query: parsed.Query(), Body: string(request.Body)},
		Response: response,
		Want:     testkit.ContractExpectation{Attempts: attempts},
	}, nil
}
