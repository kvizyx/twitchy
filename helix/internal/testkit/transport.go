package testkit

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

var ErrNoResponse = errors.New("testkit: recording transport has no response")

type RequestRecord struct {
	Method string
	Path   string
	Header http.Header
	Body   []byte
}

type RoundTripResponse struct {
	StatusCode int
	Header     http.Header
	Body       string
}

func ResponseFromFixture(response ContractResponse) RoundTripResponse {
	return RoundTripResponse{StatusCode: response.Status, Header: response.Headers.Clone(), Body: response.Body}
}

type RecordingRoundTripper struct {
	mu        sync.Mutex
	responses []RoundTripResponse
	requests  []RequestRecord
}

func NewRecordingRoundTripper(responses ...RoundTripResponse) *RecordingRoundTripper {
	return &RecordingRoundTripper{responses: append([]RoundTripResponse(nil), responses...)}
}

func (transport *RecordingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	if request == nil || request.URL == nil {
		return nil, errors.New("testkit: request and URL are required")
	}
	var body []byte
	var err error
	if request.Body != nil {
		body, err = io.ReadAll(request.Body)
		if err != nil {
			return nil, fmt.Errorf("testkit: read request body: %w", err)
		}
	}
	transport.mu.Lock()
	transport.requests = append(transport.requests, RequestRecord{Method: request.Method, Path: request.URL.RequestURI(), Header: request.Header.Clone(), Body: append([]byte(nil), body...)})
	index := len(transport.requests) - 1
	transport.mu.Unlock()
	if index >= len(transport.responses) {
		return nil, ErrNoResponse
	}
	response := transport.responses[index]
	status := response.StatusCode
	if status == 0 {
		status = http.StatusOK
	}
	return &http.Response{StatusCode: status, Header: response.Header.Clone(), Body: io.NopCloser(strings.NewReader(response.Body)), Request: request}, nil
}

func (transport *RecordingRoundTripper) Requests() []RequestRecord {
	transport.mu.Lock()
	defer transport.mu.Unlock()
	requests := make([]RequestRecord, len(transport.requests))
	for index, request := range transport.requests {
		requests[index] = RequestRecord{Method: request.Method, Path: request.Path, Header: request.Header.Clone(), Body: append([]byte(nil), request.Body...)}
	}
	return requests
}

type ExternalDialError struct{ Host string }

func (err *ExternalDialError) Error() string {
	return fmt.Sprintf("testkit: external dial blocked for host %q", err.Host)
}

type FailingRoundTripper struct {
	Next         http.RoundTripper
	AllowedHosts map[string]struct{}
}

func NewFailingRoundTripper(next http.RoundTripper, allowedHosts ...string) *FailingRoundTripper {
	allow := make(map[string]struct{}, len(allowedHosts))
	for _, host := range allowedHosts {
		allow[strings.ToLower(strings.TrimSpace(host))] = struct{}{}
	}
	return &FailingRoundTripper{Next: next, AllowedHosts: allow}
}

func (transport *FailingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	if request == nil || request.URL == nil {
		return nil, errors.New("testkit: request and URL are required")
	}
	host := strings.ToLower(request.URL.Hostname())
	if !isLoopbackHost(host) {
		if _, allowed := transport.AllowedHosts[host]; !allowed {
			return nil, &ExternalDialError{Host: host}
		}
	}
	if transport.Next == nil {
		return nil, ErrNoResponse
	}
	return transport.Next.RoundTrip(request)
}

func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

type OfflineRoundTripper = FailingRoundTripper

func NewOfflineRoundTripper(next http.RoundTripper, allowedHosts ...string) *OfflineRoundTripper {
	return NewFailingRoundTripper(next, allowedHosts...)
}

var _ http.RoundTripper = (*RecordingRoundTripper)(nil)
var _ http.RoundTripper = (*FailingRoundTripper)(nil)
