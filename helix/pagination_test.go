package helix

import "testing"

type paginationRequestFixture struct {
	After  *string  `query:"after,omitempty"`
	Before *string  `query:"before,omitempty"`
	IDs    []string `query:"id"`
}

func TestPagination_rejectsBothDirections(t *testing.T) {
	// Given a request containing both cursor directions.
	after := "after"
	before := "before"
	request := paginationRequestFixture{After: &after, Before: &before}

	// When its endpoint pagination plan is built.
	_, err := newPaginationPlan(request, "after")

	// Then the existing typed exclusivity error is returned before any request.
	if err == nil {
		t.Fatal("newPaginationPlan() error = nil")
	}
	if _, ok := err.(*ExclusiveParametersError); !ok {
		t.Fatalf("newPaginationPlan() error = %T, want *ExclusiveParametersError", err)
	}
}

func TestPagination_preservesBackwardDirectionAndDuplicateValues(t *testing.T) {
	// Given a backward request and duplicate IDs.
	before := "old"
	request := paginationRequestFixture{Before: &before, IDs: []string{"same", "same"}}

	// When its endpoint plan selects the request cursor.
	plan, err := newPaginationPlan(request, "after")
	if err != nil {
		t.Fatalf("newPaginationPlan() error = %v", err)
	}
	cloned, err := plan.withCursor(request, "older")
	if err != nil {
		t.Fatalf("withCursor() error = %v", err)
	}
	got := cloned.(paginationRequestFixture)

	// Then backward paging updates only its cursor and preserves duplicate items.
	if got.Before == nil || *got.Before != "older" || got.After != nil {
		t.Fatalf("cursor fields = after %v, before %v", got.After, got.Before)
	}
	if len(got.IDs) != 2 || got.IDs[0] != "same" || got.IDs[1] != "same" {
		t.Fatalf("IDs = %v, want duplicate values preserved", got.IDs)
	}
}

func TestPagination_supportsStringCursorException(t *testing.T) {
	// Given an endpoint whose cursor field is a string pointer.
	request := paginationRequestFixture{}

	// When the endpoint plan prepares its first cursor.
	plan, err := newPaginationPlan(request, "after")
	if err != nil {
		t.Fatalf("newPaginationPlan() error = %v", err)
	}
	cloned, err := plan.withCursor(request, "extension-cursor")
	if err != nil {
		t.Fatalf("withCursor() error = %v", err)
	}

	// Then the opaque string cursor is preserved exactly.
	got := cloned.(paginationRequestFixture)
	if got.After == nil || *got.After != "extension-cursor" {
		t.Fatalf("after cursor = %v, want extension-cursor", got.After)
	}
}
