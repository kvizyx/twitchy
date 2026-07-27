package conformance

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/manifest"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func executeReplayContract(t *testing.T, operation manifest.Operation) {
	t.Helper()
	status := successfulStatus(operation)
	responses := []testkit.RoundTripResponse{
		{StatusCode: http.StatusServiceUnavailable, Body: `{"error":"Unavailable","status":503,"message":"retry"}`},
		{StatusCode: status, Body: responseBody(operation, status)},
	}
	transport := testkit.NewRecordingRoundTripper(responses...)
	client, err := newClient(operation, transport)
	if err != nil {
		t.Fatal(err)
	}
	method, err := resolveMethod(operation, client)
	if err != nil {
		t.Fatal(err)
	}
	result, callErr := invoke(method, operationRequest(operation, method))
	if callErr != nil {
		t.Fatalf("replayable call: %v", callErr)
	}
	meta, err := responseMeta(result)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Attempts() != 2 || len(transport.Requests()) != 2 {
		t.Fatalf("replay attempts = %d/%d, want 2/2", meta.Attempts(), len(transport.Requests()))
	}
	assertContractFixture(t, operation, transport, testkit.ContractResponse{Status: status, Body: responseBody(operation, status), Success: true}, meta, nil, 2)
}

func executePagerContract(t *testing.T, operation manifest.Operation) {
	t.Helper()
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: successfulStatus(operation), Body: responseBody(operation, successfulStatus(operation))})
	client, err := newClient(operation, transport)
	if err != nil {
		t.Fatal(err)
	}
	method, err := resolveMethod(operation, client)
	if err != nil {
		t.Fatal(err)
	}
	service, err := resolveService(client, operation.Implementation.Selector)
	if err != nil {
		t.Fatal(err)
	}
	pagerMethod := service.MethodByName(operation.Implementation.Method + "Pager")
	if !pagerMethod.IsValid() {
		t.Fatalf("pagination shape %q has no pager method", operation.Pagination.Shape)
	}
	results := pagerMethod.Call([]reflect.Value{operationRequest(operation, method)})
	if errValue := results[1].Interface(); errValue != nil {
		t.Fatalf("pager construction: %v", errValue)
	}
	pager := results[0]
	next := pager.MethodByName("Next").Call([]reflect.Value{reflect.ValueOf(context.Background())})
	if len(next) != 1 || !next[0].Bool() {
		t.Fatal("pager did not return its first page")
	}
	if errValue := pager.MethodByName("Err").Call(nil)[0].Interface(); errValue != nil {
		t.Fatalf("pager error: %v", errValue)
	}
	if len(transport.Requests()) != 1 {
		t.Fatalf("pager requests = %d, want 1", len(transport.Requests()))
	}
}

func assertContractFixture(t *testing.T, operation manifest.Operation, transport *testkit.RecordingRoundTripper, response testkit.ContractResponse, meta helix.ResponseMeta, callErr error, attempts int) {
	t.Helper()
	fixture, err := requestFixture(operation, transport, response, attempts)
	if err != nil {
		t.Fatal(err)
	}
	if operation.Response.Format == "text" {
		return
	}
	if err := testkit.RunManifestContract(context.Background(), operation, fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return meta, callErr
	}); err != nil {
		t.Fatal(err)
	}
}

func TestManifestConformance_rejectsExternalRequests(t *testing.T) {
	// Given the offline transport used by every conformance call.
	transport := testkit.NewOfflineRoundTripper(testkit.NewRecordingRoundTripper())
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.twitch.tv/helix", nil)
	if err != nil {
		t.Fatal(err)
	}

	// When a non-loopback request reaches the transport.
	_, err = transport.RoundTrip(request)

	// Then it is rejected before any delegate can receive it.
	var externalErr *testkit.ExternalDialError
	if !errors.As(err, &externalErr) || externalErr.Host != "api.twitch.tv" {
		t.Fatalf("offline transport error = %v, want external dial rejection", err)
	}
}

func TestManifestConformance_rechecksSecurityBoundaries(t *testing.T) {
	// Given a bounded response body and a credential containing a secret.
	if helix.DefaultResponseBodyLimit < 1 || helix.DefaultErrorExcerptLimit < 1 {
		t.Fatal("body bounds are not positive")
	}
	credential := helix.Credential{AccessToken: "secret-token"}
	if got := credential.String(); strings.Contains(got, credential.AccessToken) {
		t.Fatalf("credential redaction leaked access token: %s", got)
	}
	if got := fmt.Sprintf("%#v", credential); strings.Contains(got, credential.AccessToken) {
		t.Fatalf("credential GoString redaction leaked access token: %s", got)
	}
}
