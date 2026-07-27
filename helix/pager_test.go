package helix

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
)

type pageFixtureTransport struct {
	fixtures []string
	calls    atomic.Int32
}

func (t *pageFixtureTransport) RoundTrip(*http.Request) (*http.Response, error) {
	index := int(t.calls.Add(1)) - 1
	if index >= len(t.fixtures) {
		return nil, io.ErrUnexpectedEOF
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(t.fixtures[index])),
		Header:     make(http.Header),
	}, nil
}

func fixtureFetcher(t *testing.T, fixtures ...string) (pageFetcher[[]int], *pageFixtureTransport) {
	t.Helper()
	transport := &pageFixtureTransport{fixtures: fixtures}
	client := &http.Client{Transport: transport}
	return func(ctx context.Context, cursor string) (*Response[[]int], error) {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://fixture.test/helix?after="+cursor, nil)
		if err != nil {
			return nil, err
		}
		response, err := client.Do(request)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()
		return decodeResponse[[]int](response.StatusCode, response.Body, DecodeOptions{})
	}, transport
}

func TestPager_doesNotRequestBeforeFirstNext(t *testing.T) {
	// Given a fixture-backed page fetcher and a new pager.
	fetch, transport := fixtureFetcher(t, `{"data":[1],"pagination":{}}`)
	pager, err := newPager(fetch)
	if err != nil {
		t.Fatalf("newPager() error = %v", err)
	}

	// When no page has been requested yet.
	// Then the pager is empty and the transport is untouched.
	if pager.Page() != nil || pager.Err() != nil {
		t.Fatalf("initial pager state = page %v, err %v", pager.Page(), pager.Err())
	}
	if got := transport.calls.Load(); got != 0 {
		t.Fatalf("initial request count = %d, want 0", got)
	}
}

func TestPager_iteratesLazilyAndExhaustsOnAbsentCursor(t *testing.T) {
	// Given two fixture pages with a cursor only on the first page.
	fetch, transport := fixtureFetcher(t,
		`{"data":[1,1],"pagination":{"cursor":"next"}}`,
		`{"data":[2],"pagination":{}}`,
	)
	pager, err := newPager(fetch)
	if err != nil {
		t.Fatalf("newPager() error = %v", err)
	}

	// When each page is requested explicitly.
	if !pager.Next(context.Background()) {
		t.Fatalf("first Next() = false, err = %v", pager.Err())
	}
	if got := pager.Page().Data; len(got) != 2 || got[0] != 1 || got[1] != 1 {
		t.Fatalf("first page data = %v, want duplicate values preserved", got)
	}
	if pager.Err() != nil || transport.calls.Load() != 1 {
		t.Fatalf("after first page = err %v, calls %d", pager.Err(), transport.calls.Load())
	}
	if !pager.Next(context.Background()) {
		t.Fatalf("second Next() = false, err = %v", pager.Err())
	}
	if got := pager.Page().Data; len(got) != 1 || got[0] != 2 {
		t.Fatalf("second page data = %v, want [2]", got)
	}
	if pager.Next(context.Background()) {
		t.Fatal("exhausted Next() = true, want false")
	}

	// Then no request occurs after clean exhaustion and the last page remains visible.
	if pager.Err() != nil || transport.calls.Load() != 2 {
		t.Fatalf("exhausted state = err %v, calls %d", pager.Err(), transport.calls.Load())
	}
	if got := pager.Page().Data[0]; got != 2 {
		t.Fatalf("retained page value = %d, want 2", got)
	}
}

func TestPager_emptyPageWithCursorContinues(t *testing.T) {
	// Given an empty page that still supplies a cursor.
	fetch, transport := fixtureFetcher(t,
		`{"data":[],"pagination":{"cursor":"next"}}`,
		`{"data":[3],"pagination":{}}`,
	)
	pager, err := newPager(fetch)
	if err != nil {
		t.Fatalf("newPager() error = %v", err)
	}

	// When the empty page is consumed.
	if !pager.Next(context.Background()) {
		t.Fatalf("empty page Next() = false, err = %v", pager.Err())
	}
	if !pager.Next(context.Background()) {
		t.Fatalf("empty-page flow stopped early, err = %v", pager.Err())
	}

	// Then the cursor, rather than the item count, controls continuation.
	if got := pager.Page().Data; len(got) != 1 || got[0] != 3 || transport.calls.Load() != 2 {
		t.Fatalf("final page = %v, calls = %d", got, transport.calls.Load())
	}
}

func TestPager_repeatedCursorRetainsPreviousPage(t *testing.T) {
	// Given a response that repeats an already exposed cursor.
	fetch, transport := fixtureFetcher(t,
		`{"data":[1],"pagination":{"cursor":"same"}}`,
		`{"data":[2],"pagination":{"cursor":"same"}}`,
	)
	pager, err := newPager(fetch)
	if err != nil {
		t.Fatalf("newPager() error = %v", err)
	}
	if !pager.Next(context.Background()) {
		t.Fatalf("first Next() = false, err = %v", pager.Err())
	}

	// When the repeated cursor is received.
	if pager.Next(context.Background()) {
		t.Fatal("cycle Next() = true, want false")
	}

	// Then the returned page is not exposed and the terminal error is stable.
	if pager.Err() != ErrCursorCycle {
		t.Fatalf("cycle error = %v, want %v", pager.Err(), ErrCursorCycle)
	}
	if got := pager.Page().Data[0]; got != 1 {
		t.Fatalf("retained page value = %d, want 1", got)
	}
	if pager.Next(context.Background()) || pager.Err() != ErrCursorCycle || transport.calls.Load() != 2 {
		t.Fatal("terminal cycle state changed after another Next()")
	}
}

func TestPager_pageLimitDefersErrorUntilNextWithoutIO(t *testing.T) {
	// Given a one-page cap and a page that advertises another cursor.
	fetch, transport := fixtureFetcher(t, `{"data":[1],"pagination":{"cursor":"next"}}`)
	pager, err := newPager(fetch, WithPageLimit(1))
	if err != nil {
		t.Fatalf("newPager() error = %v", err)
	}
	if !pager.Next(context.Background()) {
		t.Fatalf("first Next() = false, err = %v", pager.Err())
	}

	// When the caller cancels before attempting the capped continuation.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if pager.Next(ctx) {
		t.Fatal("cap-pending Next() = true, want false")
	}

	// Then the cap wins over context cancellation and performs no I/O.
	if pager.Err() != ErrPageLimit || transport.calls.Load() != 1 {
		t.Fatalf("cap state = err %v, calls %d", pager.Err(), transport.calls.Load())
	}
}

func TestPager_pageLimitWithoutCursorExhaustsCleanly(t *testing.T) {
	// Given a one-page cap and a page without a continuation cursor.
	fetch, transport := fixtureFetcher(t, `{"data":[1],"pagination":{}}`)
	pager, err := newPager(fetch, WithPageLimit(1))
	if err != nil {
		t.Fatalf("newPager() error = %v", err)
	}

	// When the capped page is consumed and the caller asks for another page.
	if !pager.Next(context.Background()) || pager.Next(context.Background()) {
		t.Fatal("cleanly exhausted capped pager did not retain a successful page")
	}

	// Then the cap does not manufacture an error when the endpoint is already exhausted.
	if pager.Err() != nil || transport.calls.Load() != 1 || pager.Page().Data[0] != 1 {
		t.Fatalf("clean cap state = err %v, calls %d, page %v", pager.Err(), transport.calls.Load(), pager.Page())
	}
}

func TestPager_cancellationRetainsLastSuccess(t *testing.T) {
	// Given one successful page and a canceled context for the next request.
	fetch, transport := fixtureFetcher(t, `{"data":[1],"pagination":{"cursor":"next"}}`)
	pager, err := newPager(fetch)
	if err != nil {
		t.Fatalf("newPager() error = %v", err)
	}
	if !pager.Next(context.Background()) {
		t.Fatalf("first Next() = false, err = %v", pager.Err())
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// When the canceled continuation is attempted.
	if pager.Next(ctx) {
		t.Fatal("canceled Next() = true, want false")
	}

	// Then context error is terminal, the previous page remains, and no request is made.
	if pager.Err() != context.Canceled || pager.Page().Data[0] != 1 || transport.calls.Load() != 1 {
		t.Fatalf("canceled state = err %v, page %v, calls %d", pager.Err(), pager.Page(), transport.calls.Load())
	}
}

func TestPager_rejectsInvalidPageLimit(t *testing.T) {
	// Given limits outside the documented range.
	fetch, _ := fixtureFetcher(t, `{"data":[],"pagination":{}}`)
	for _, limit := range []int{0, 10001} {
		// When the pager is constructed.
		_, err := newPager(fetch, WithPageLimit(limit))

		// Then construction rejects the invalid option.
		if err == nil {
			t.Fatalf("limit %d: newPager() error = nil", limit)
		}
	}
	for _, limit := range []int{1, 10000} {
		if _, err := newPager(fetch, WithPageLimit(limit)); err != nil {
			t.Fatalf("limit %d: newPager() error = %v", limit, err)
		}
	}
}

func TestPager_fixturePayloadsRemainJSONCompatible(t *testing.T) {
	// Given a decoded page fixture.
	var response Response[[]int]
	if err := json.Unmarshal([]byte(`{"data":[1],"pagination":{"cursor":"next"}}`), &response); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Then the typed response retains its cursor for the pager contract.
	if response.Pagination.Cursor() != "next" {
		t.Fatalf("cursor = %q, want next", response.Pagination.Cursor())
	}
}
