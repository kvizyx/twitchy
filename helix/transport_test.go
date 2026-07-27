package helix

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type scriptedTransport struct {
	responses []int
	calls     atomic.Int32
}

func (t *scriptedTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if request.Body != nil {
		_, _ = io.ReadAll(request.Body)
	}
	index := int(t.calls.Add(1)) - 1
	if index >= len(t.responses) {
		return nil, io.ErrUnexpectedEOF
	}
	status := t.responses[index]
	header := make(http.Header)
	header.Set("X-Request-ID", "request-")
	if status == http.StatusTooManyRequests {
		header.Set("Ratelimit-Limit", "10")
		header.Set("Ratelimit-Remaining", "0")
		header.Set("Ratelimit-Reset", "1704067290")
	}
	if status == http.StatusServiceUnavailable {
		header.Set("Retry-After", "300")
	}
	body := `{"data":{}}`
	if status >= 400 {
		body = `{"error":"failure","status":` + string(rune(status/10+48)) + `,"message":"failure"}`
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     header,
		Request:    request,
	}, nil
}

type refreshTestSource struct {
	snapshot  CredentialSnapshot
	refreshes atomic.Int32
}

func (s *refreshTestSource) Token(context.Context) (CredentialSnapshot, error) {
	return s.snapshot, nil
}

func (s *refreshTestSource) Refresh(context.Context, CredentialSnapshot, RefreshReason) (CredentialSnapshot, error) {
	s.refreshes.Add(1)
	s.snapshot = NewCredentialSnapshot(Credential{
		AccessToken: "new-token",
		TokenClass:  TokenClassUser,
		Refreshable: true,
		Generation:  2,
	})
	return s.snapshot, nil
}

type testClock struct{ now time.Time }

func (c testClock) Now() time.Time { return c.now }

type recordingSleeper struct {
	durations chan time.Duration
}

func (s *recordingSleeper) Sleep(ctx context.Context, duration time.Duration) error {
	s.durations <- duration
	return ctx.Err()
}

func testOperation() manifest.Operation {
	return manifest.Operation{
		OperationID: "test-operation",
		Method:      http.MethodGet,
		Request:     manifest.RequestSpec{BodyReconstructible: true},
		Replay:      manifest.ReplaySpec{Replayable: true, BucketWaitable: true},
	}
}

func testRequest(t *testing.T) *http.Request {
	request, err := http.NewRequest(http.MethodGet, "https://api.example.test/helix", strings.NewReader(""))
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	return request
}

func TestAttemptStateMachine(t *testing.T) {
	tests := []struct {
		name      string
		responses []int
		wantCalls int
		refresh   bool
		wantError bool
	}{
		{name: "401 then success", responses: []int{401, 200}, wantCalls: 2, refresh: true},
		{name: "503 then success", responses: []int{503, 200}, wantCalls: 2},
		{name: "repeated 401", responses: []int{401, 401}, wantCalls: 2, refresh: true, wantError: true},
		{name: "401 503 success", responses: []int{401, 503, 200}, wantCalls: 3, refresh: true},
		{name: "429 401 success", responses: []int{429, 401, 200}, wantCalls: 3, refresh: true},
		{name: "401 503 429 terminal", responses: []int{401, 503, 429}, wantCalls: 3, refresh: true, wantError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := &scriptedTransport{responses: test.responses}
			source := &refreshTestSource{snapshot: NewCredentialSnapshot(Credential{
				AccessToken: "old-token",
				TokenClass:  TokenClassUser,
				Refreshable: test.refresh,
				Generation:  1,
			})}
			sleeper := &recordingSleeper{durations: make(chan time.Duration, 1)}
			executor := newTransportExecutor(&http.Client{Transport: transport}, source, RateLimitPolicy{Wait: true}, testClock{now: time.Unix(1704067200, 0)}, sleeper)

			_, _, err := executor.execute(context.Background(), testRequest(t), testOperation(), source.snapshot)

			if got := int(transport.calls.Load()); got != test.wantCalls {
				t.Fatalf("attempts = %d, want %d", got, test.wantCalls)
			}
			if test.wantError != (err != nil) {
				t.Fatalf("error = %v, want error=%t", err, test.wantError)
			}
			if test.responses[0] == http.StatusTooManyRequests {
				<-sleeper.durations
			}
		})
	}
}

func TestAttemptStateMachine_doesNotReplayMutationOrOneShotBody(t *testing.T) {
	for _, replayable := range []bool{false, true} {
		transport := &scriptedTransport{responses: []int{503, 200}}
		executor := newTransportExecutor(&http.Client{Transport: transport}, nil, RateLimitPolicy{}, testClock{now: time.Unix(1704067200, 0)}, nil)
		operation := testOperation()
		operation.Replay.Replayable = replayable
		operation.Request.BodyReconstructible = replayable
		body := &oneShotReader{}
		request, err := http.NewRequest(http.MethodPost, "https://api.example.test/helix", body)
		if err != nil {
			t.Fatalf("NewRequest() error = %v", err)
		}

		_, _, _ = executor.execute(context.Background(), request, operation, CredentialSnapshot{})
		if got := int(transport.calls.Load()); got != 1 {
			t.Fatalf("attempts = %d, want 1", got)
		}
		if got := body.reads.Load(); got != 1 {
			t.Fatalf("one-shot body reads = %d, want 1", got)
		}
	}
}

type oneShotReader struct{ reads atomic.Int32 }

func (r *oneShotReader) Read([]byte) (int, error) {
	if r.reads.Add(1) > 1 {
		return 0, errors.New("one-shot body read twice")
	}
	return 0, io.EOF
}

func TestRatePolicy(t *testing.T) {
	clock := testClock{now: time.Unix(1704067200, 0)}
	sleeper := &recordingSleeper{durations: make(chan time.Duration, 1)}
	transport := &scriptedTransport{responses: []int{429, 200}}
	executor := newTransportExecutor(&http.Client{Transport: transport}, nil, RateLimitPolicy{Wait: true}, clock, sleeper)

	_, _, err := executor.execute(context.Background(), testRequest(t), testOperation(), CredentialSnapshot{})
	if err != nil {
		t.Fatalf("execute() error = %v", err)
	}
	if got := <-sleeper.durations; got != time.Minute {
		t.Fatalf("wait = %s, want 1m", got)
	}

	transport = &scriptedTransport{responses: []int{429, 200}}
	sleeper = &recordingSleeper{durations: make(chan time.Duration, 1)}
	operation := testOperation()
	operation.Replay.BucketWaitable = false
	executor = newTransportExecutor(&http.Client{Transport: transport}, nil, RateLimitPolicy{Wait: true}, clock, sleeper)
	_, _, err = executor.execute(context.Background(), testRequest(t), operation, CredentialSnapshot{})
	if err == nil || transport.calls.Load() != 1 {
		t.Fatalf("cooldown result = %v, calls = %d, want immediate typed error and one call", err, transport.calls.Load())
	}

	transport = &scriptedTransport{responses: []int{429, 200}}
	executor = newTransportExecutor(&http.Client{Transport: transport}, nil, RateLimitPolicy{}, clock, sleeper)
	_, _, err = executor.execute(context.Background(), testRequest(t), testOperation(), CredentialSnapshot{})
	if err == nil || transport.calls.Load() != 1 {
		t.Fatalf("default result = %v, calls = %d, want immediate typed error and one call", err, transport.calls.Load())
	}

	transport = &scriptedTransport{responses: []int{429, 200}}
	sleeper = &recordingSleeper{durations: make(chan time.Duration, 1)}
	executor = newTransportExecutor(&http.Client{Transport: transport}, nil, RateLimitPolicy{Wait: true, MaxWait: 2 * time.Second}, clock, sleeper)
	_, _, err = executor.execute(context.Background(), testRequest(t), testOperation(), CredentialSnapshot{})
	if err != nil || <-sleeper.durations != 2*time.Second {
		t.Fatalf("configured wait result = %v, want a 2s bounded wait", err)
	}

	transport = &scriptedTransport{responses: []int{429, 200}}
	sleeper = &recordingSleeper{durations: make(chan time.Duration, 1)}
	executor = newTransportExecutor(&http.Client{Transport: transport}, nil, RateLimitPolicy{Wait: true}, clock, sleeper)
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, err = executor.execute(canceled, testRequest(t), testOperation(), CredentialSnapshot{})
	if !errors.Is(err, context.Canceled) || transport.calls.Load() != 0 {
		t.Fatalf("canceled result = %v, calls = %d, want context cancellation before I/O", err, transport.calls.Load())
	}
}

func TestAttemptStateMachine_coalescesConcurrentRefresh(t *testing.T) {
	transport := &alternatingTransport{barrier: make(chan struct{})}
	source := &coalescingSource{snapshot: NewCredentialSnapshot(Credential{
		AccessToken: "old-token",
		TokenClass:  TokenClassUser,
		Refreshable: true,
		Generation:  1,
	})}
	executor := newTransportExecutor(&http.Client{Transport: transport}, source, RateLimitPolicy{}, testClock{now: time.Unix(1704067200, 0)}, nil)

	results := make(chan error, 2)
	for range 2 {
		go func() {
			_, _, err := executor.execute(context.Background(), testRequest(t), testOperation(), source.snapshot)
			results <- err
		}()
	}
	for range 2 {
		if err := <-results; err != nil {
			t.Fatalf("concurrent execute() error = %v", err)
		}
	}
	if got := source.refreshes.Load(); got != 1 {
		t.Fatalf("refresh calls = %d, want 1", got)
	}
}

type alternatingTransport struct {
	calls   atomic.Int32
	barrier chan struct{}
	once    sync.Once
}

func (t *alternatingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	count := t.calls.Add(1)
	status := http.StatusOK
	if count <= 2 {
		if count == 2 {
			t.once.Do(func() { close(t.barrier) })
		}
		<-t.barrier
		status = http.StatusUnauthorized
	}
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(`{"data":{}}`)), Header: make(http.Header), Request: request}, nil
}

type coalescingSource struct {
	snapshot  CredentialSnapshot
	refreshes atomic.Int32
}

func (s *coalescingSource) Token(context.Context) (CredentialSnapshot, error) { return s.snapshot, nil }

func (s *coalescingSource) Refresh(context.Context, CredentialSnapshot, RefreshReason) (CredentialSnapshot, error) {
	s.refreshes.Add(1)
	return NewCredentialSnapshot(Credential{AccessToken: "new-token", TokenClass: TokenClassUser, Refreshable: true, Generation: 2}), nil
}
