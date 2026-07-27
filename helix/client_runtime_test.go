package helix

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestClientRejectsRedirects(t *testing.T) {
	// Given a server with a redirect followed by a target endpoint
	var targetHits atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/target" {
			targetHits.Add(1)
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		http.Redirect(writer, request, "/target", http.StatusFound)
	}))
	defer server.Close()

	// When the caller HTTP client is used to construct a Helix client
	client, err := New(WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	_, err = client.httpClient.Get(server.URL)

	// Then the redirect is rejected before the target receives a request.
	if err == nil {
		t.Fatal("redirect unexpectedly succeeded")
	}
	if targetHits.Load() != 0 {
		t.Fatal("redirect target received a request")
	}
}

func TestClientZeroValueIsInvalid(t *testing.T) {
	// Given a zero-value Client
	var client Client

	// When its construction invariant is checked
	err := client.validClient()

	// Then it is rejected as an invalid client.
	if !errors.Is(err, ErrInvalidClient) {
		t.Fatalf("validClient() error = %v, want ErrInvalidClient", err)
	}
}
