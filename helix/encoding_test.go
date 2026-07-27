package helix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestEncodeQuery_repeatsSlicesAndPreservesPointers(t *testing.T) {
	type query struct {
		IDs    []string `query:"id"`
		After  *string  `query:"after,omitempty"`
		Active *bool    `query:"active,omitempty"`
		Count  *int     `query:"count,omitempty"`
		Empty  *string  `query:"empty,omitempty"`
	}
	active := false
	count := 0
	empty := ""

	values, err := encodeQuery(query{IDs: []string{"a", "b"}, Active: &active, Count: &count, Empty: &empty})
	if err != nil {
		t.Fatalf("encode query: %v", err)
	}

	if got, want := values["id"], []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("id values: got %v, want %v", got, want)
	}
	if got, want := values.Encode(), "active=false&count=0&empty=&id=a&id=b"; got != want {
		t.Fatalf("encoded query: got %q, want %q", got, want)
	}
}

func TestEncodeQuery_rejectsExclusiveCursors(t *testing.T) {
	after, before := "after", "before"
	_, err := encodeQuery(struct {
		After  *string `query:"after,omitempty"`
		Before *string `query:"before,omitempty"`
	}{After: &after, Before: &before})
	var exclusive *ExclusiveParametersError
	if !errors.As(err, &exclusive) {
		t.Fatalf("error: got %T %v, want *ExclusiveParametersError", err, err)
	}
}

func TestEncodeQuery_rejectsUnsupportedTagOption(t *testing.T) {
	_, err := encodeQuery(struct {
		ID string `query:"id,comma"`
	}{ID: "x"})
	var unsupported *UnsupportedTagError
	if !errors.As(err, &unsupported) {
		t.Fatalf("error: got %T %v, want *UnsupportedTagError", err, err)
	}
}

func TestEncodeBody_supportsJSONAndFormAndNoBody(t *testing.T) {
	type payload struct {
		ID      string `json:"id"`
		Enabled bool   `json:"enabled"`
	}

	jsonBody, err := encodeJSONBody(payload{ID: "9007199254740993", Enabled: false})
	if err != nil {
		t.Fatalf("encode JSON: %v", err)
	}
	if got, want := string(jsonBody), `{"id":"9007199254740993","enabled":false}`; got != want {
		t.Fatalf("JSON body: got %q, want %q", got, want)
	}

	formBody, err := encodeFormBody(struct {
		Scope []string `form:"scope"`
	}{Scope: []string{"a", "b"}})
	if err != nil {
		t.Fatalf("encode form: %v", err)
	}
	if got, want := string(formBody), "scope=a&scope=b"; got != want {
		t.Fatalf("form body: got %q, want %q", got, want)
	}

	request, err := buildRequest(requestSpec{Context: context.Background(), Method: http.MethodGet, URL: "https://example.test/helix", Body: nil})
	if err != nil {
		t.Fatalf("build no-body request: %v", err)
	}
	if request.Body != nil {
		t.Fatal("GET request unexpectedly has a body")
	}
}

func TestEncodeBody_allowsJSONMutationMethodsOnly(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPatch, http.MethodPut} {
		request, err := buildRequest(requestSpec{
			Context: context.Background(),
			Method:  method,
			URL:     "https://example.test/helix",
			Body: struct {
				Enabled bool `json:"enabled"`
			}{Enabled: false},
		})
		if err != nil {
			t.Fatalf("build %s request: %v", method, err)
		}
		if got := request.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("%s content type: got %q", method, got)
		}
	}

	_, err := buildRequest(requestSpec{
		Context: context.Background(),
		Method:  http.MethodGet,
		URL:     "https://example.test/helix",
		Body: struct {
			ID string `json:"id"`
		}{ID: "x"},
	})
	var requestErr *RequestEncodingError
	if !errors.As(err, &requestErr) {
		t.Fatalf("GET body error: got %T %v, want *RequestEncodingError", err, err)
	}
}

func TestDecodeResponse_ignoresUnknownFieldsAndHandlesEmpty(t *testing.T) {
	type data struct {
		ID string `json:"id"`
	}
	response, err := decodeResponse[data](200, strings.NewReader(`{"data":{"id":"12345678901234567890"},"new_field":true}`), DecodeOptions{})
	if err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Data.ID != "12345678901234567890" {
		t.Fatalf("ID was coerced or lost: %q", response.Data.ID)
	}

	empty, err := decodeResponse[data](http.StatusNoContent, strings.NewReader(""), DecodeOptions{})
	if err != nil {
		t.Fatalf("decode empty response: %v", err)
	}
	if empty.Data.ID != "" {
		t.Fatalf("empty response data: %#v", empty.Data)
	}
}

func TestDecodeResponse_returnsTypedErrorsForMalformedAndOversizedJSON(t *testing.T) {
	_, err := decodeResponse[struct{}](200, strings.NewReader("{"), DecodeOptions{})
	var decodeErr *JSONDecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("malformed error: got %T %v, want *JSONDecodeError", err, err)
	}

	_, err = decodeResponse[struct{}](200, bytes.NewReader([]byte(`{"data":"too large"}`)), DecodeOptions{Limits: BodyLimits{Response: 4}})
	var limitErr *BodyLimitError
	if !errors.As(err, &limitErr) {
		t.Fatalf("oversized error: got %T %v, want *BodyLimitError", err, err)
	}
}

func TestDecodeResponse_usesLocalCodec(t *testing.T) {
	called := false
	codec := JSONCodec{
		Marshal: json.Marshal,
		Unmarshal: func(data []byte, value any) error {
			called = true
			return json.Unmarshal(data, value)
		},
	}
	_, err := decodeResponse[struct{}](200, strings.NewReader(`{"data":{}}`), DecodeOptions{Codec: codec})
	if err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !called {
		t.Fatal("local JSON codec was not used")
	}
}

func TestReadErrorExcerpt_isBounded(t *testing.T) {
	excerpt, truncated, err := readErrorExcerpt(strings.NewReader(strings.Repeat("x", 10)), BodyLimits{ErrorExcerpt: 4})
	if err != nil {
		t.Fatalf("read excerpt: %v", err)
	}
	if string(excerpt) != "xxxx" || !truncated {
		t.Fatalf("excerpt: got %q, truncated=%t", excerpt, truncated)
	}
	if _, err := io.ReadAll(strings.NewReader("")); err != nil {
		t.Fatalf("unrelated reader check: %v", err)
	}
}
