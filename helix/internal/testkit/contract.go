package testkit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type ContractRequest struct {
	Path    string      `json:"path"`
	Query   url.Values  `json:"query"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
}

type ContractResponse struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
	Success bool        `json:"success"`
}

type ContractExpectation struct {
	Attempts       int  `json:"attempts"`
	RateLimitValid bool `json:"rate_limit_valid"`
	RateLimit      struct {
		Limit     int `json:"limit"`
		Remaining int `json:"remaining"`
	} `json:"rate_limit"`
}

type ContractFixture struct {
	Request  ContractRequest     `json:"request"`
	Response ContractResponse    `json:"response"`
	Want     ContractExpectation `json:"want"`
}

func RunManifestContract(ctx context.Context, operation manifest.Operation, fixture ContractFixture, transport *RecordingRoundTripper, call func(context.Context) (helix.ResponseMeta, error)) error {
	if transport == nil || call == nil {
		return errors.New("testkit: contract transport and call are required")
	}
	if operation.Method == "" || operation.Path == "" {
		return fmt.Errorf("testkit: manifest operation has incomplete method/path")
	}
	meta, callErr := call(ctx)
	requests := transport.Requests()
	if len(requests) != fixture.Want.Attempts {
		return fmt.Errorf("testkit: attempts got %d, want %d", len(requests), fixture.Want.Attempts)
	}
	if len(requests) == 0 {
		return errors.New("testkit: contract call recorded no request")
	}
	request := requests[len(requests)-1]
	if request.Method != operation.Method {
		return fmt.Errorf("testkit: method got %s, want %s", request.Method, operation.Method)
	}
	wantPath := fixture.Request.Path
	if wantPath == "" {
		wantPath = operation.Path
	}
	if request.Path != appendQuery(wantPath, fixture.Request.Query) {
		return fmt.Errorf("testkit: path/query got %q, want %q", request.Path, appendQuery(wantPath, fixture.Request.Query))
	}
	for key, values := range fixture.Request.Headers {
		if !reflect.DeepEqual(request.Header.Values(key), values) {
			return fmt.Errorf("testkit: header %s got %v, want %v", key, request.Header.Values(key), values)
		}
	}
	if string(request.Body) != fixture.Request.Body {
		return fmt.Errorf("testkit: body got %q, want %q", request.Body, fixture.Request.Body)
	}
	if err := validateResponseFixture(fixture.Response); err != nil {
		return err
	}
	if meta.StatusCode() != 0 && meta.StatusCode() != fixture.Response.Status {
		return fmt.Errorf("testkit: response status got %d, want %d", meta.StatusCode(), fixture.Response.Status)
	}
	if fixture.Response.Success == (callErr != nil) {
		return fmt.Errorf("testkit: decode success=%t with error=%v", fixture.Response.Success, callErr)
	}
	if meta.Attempts() != fixture.Want.Attempts {
		return fmt.Errorf("testkit: metadata attempts got %d, want %d", meta.Attempts(), fixture.Want.Attempts)
	}
	rate := meta.RateLimit()
	if rate.Valid() != fixture.Want.RateLimitValid || (rate.Valid() && (rate.Limit() != fixture.Want.RateLimit.Limit || rate.Remaining() != fixture.Want.RateLimit.Remaining)) {
		return fmt.Errorf("testkit: rate metadata mismatch")
	}
	return nil
}

func validateResponseFixture(response ContractResponse) error {
	if response.Status >= 200 && response.Status < 300 {
		if response.Status == http.StatusNoContent || strings.TrimSpace(response.Body) == "" {
			return nil
		}
		var payload map[string]json.RawMessage
		if err := json.Unmarshal([]byte(response.Body), &payload); err != nil {
			return fmt.Errorf("testkit: decode successful response: %w", err)
		}
		if _, ok := payload["data"]; !ok {
			return errors.New("testkit: decode successful response: missing data envelope")
		}
		return nil
	}
	var payload struct {
		Error   string `json:"error"`
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(response.Body), &payload); err != nil || payload.Error == "" || payload.Status == 0 || payload.Message == "" {
		return fmt.Errorf("testkit: decode error response: malformed payload")
	}
	return nil
}

func appendQuery(path string, query url.Values) string {
	if len(query) == 0 {
		return path
	}
	keys := make([]string, 0, len(query))
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	encoded := make(url.Values, len(query))
	for _, key := range keys {
		encoded[key] = append([]string(nil), query[key]...)
	}
	separator := "?"
	if strings.Contains(path, "?") {
		separator = "&"
	}
	return path + separator + encoded.Encode()
}
