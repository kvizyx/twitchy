package manifest

import (
	"fmt"
	"net/http"
)

type Stability string

const (
	StabilityStable Stability = "stable"
	StabilityNew    Stability = "NEW"
	StabilityBeta   Stability = "BETA"
)

type TokenClass string

const (
	TokenClassApp       TokenClass = "app"
	TokenClassUser      TokenClass = "user"
	TokenClassExtension TokenClass = "extension"
	TokenClassUnknown   TokenClass = "unknown"
)

type Manifest struct {
	Schema          int         `json:"schema"`
	RetrievedAt     string      `json:"retrieved_at"`
	ReferenceSHA256 string      `json:"reference_sha256"`
	Operations      []Operation `json:"operations"`
}

type GroupFile struct {
	Schema     int         `json:"schema"`
	Group      string      `json:"group"`
	Operations []Operation `json:"operations"`
}

type Operation struct {
	OperationID    string             `json:"operation_id"`
	Anchor         string             `json:"anchor"`
	Group          string             `json:"group"`
	Name           string             `json:"name"`
	Method         string             `json:"method"`
	Path           string             `json:"path"`
	Stability      Stability          `json:"stability"`
	TokenClasses   []TokenClass       `json:"token_classes"`
	Scopes         []string           `json:"scopes"`
	SubjectBinding string             `json:"subject_binding"`
	Request        RequestSpec        `json:"request"`
	Response       ResponseSpec       `json:"response"`
	Pagination     PaginationSpec     `json:"pagination"`
	Replay         ReplaySpec         `json:"replay"`
	Implementation ImplementationSpec `json:"implementation"`
	Fixtures       []string           `json:"fixtures"`
	Source         string             `json:"source"`
}

type RequestSpec struct {
	Locations           map[string][]RequestField `json:"locations"`
	BodyReconstructible bool                      `json:"body_reconstructible"`
}

type RequestField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type ResponseSpec struct {
	Format        string `json:"format"`
	Status        []int  `json:"status"`
	StatusUnknown bool   `json:"status_unknown"`
}

type PaginationSpec struct {
	Shape           string `json:"shape"`
	CursorParameter string `json:"cursor_parameter"`
}

type ReplaySpec struct {
	Replayable     bool `json:"replayable"`
	BucketWaitable bool `json:"bucket_waitable"`
}

type ImplementationSpec struct {
	Anchor         string    `json:"anchor"`
	Selector       string    `json:"selector"`
	ServiceType    string    `json:"service_type"`
	Method         string    `json:"method"`
	Signature      string    `json:"signature"`
	PagerSignature *string   `json:"pager_signature"`
	RequestType    string    `json:"request_type"`
	DataType       string    `json:"data_type"`
	Stability      Stability `json:"stability"`
}

var expectedGroups = map[string]int{
	"Ads": 3, "Analytics": 2, "Bits": 4, "Channels": 5, "Channel Points": 6,
	"Charity": 2, "Chat": 19, "Clips": 4, "Conduits": 6, "CCLs": 1,
	"Entitlements": 2, "Extensions": 12, "EventSub": 3, "Games": 2,
	"Goals": 1, "Guest Star": 12, "Hype Train": 1, "Moderation": 25,
	"Polls": 3, "Predictions": 3, "Raids": 2, "Schedule": 6,
	"Search": 2, "Streams": 5, "Subscriptions": 2, "Tags": 2, "Teams": 2,
	"Users": 9, "Videos": 2, "Whispers": 1,
}

func (m Manifest) Validate() error {
	if len(m.Operations) != 149 {
		return fmt.Errorf("manifest operations: got %d operations, want 149", len(m.Operations))
	}

	groups := make(map[string]int, len(expectedGroups))
	labels := make(map[Stability]int, 3)
	anchors := make(map[string]struct{}, len(m.Operations))
	identities := make(map[string]struct{}, len(m.Operations))
	for _, operation := range m.Operations {
		if _, exists := anchors[operation.Anchor]; exists {
			return fmt.Errorf("operation %q: duplicate anchor", operation.Anchor)
		}
		anchors[operation.Anchor] = struct{}{}
		identity := operation.Method + " " + operation.Path
		if _, exists := identities[identity]; exists {
			return fmt.Errorf("operation %q: duplicate method+path %q", operation.Anchor, identity)
		}
		identities[identity] = struct{}{}
		if err := operation.validate(); err != nil {
			return err
		}
		groups[operation.Group]++
		labels[operation.Stability]++
	}
	if len(groups) != len(expectedGroups) {
		return fmt.Errorf("manifest groups: got %d unique groups, want %d", len(groups), len(expectedGroups))
	}
	for group, want := range expectedGroups {
		if got := groups[group]; got != want {
			return fmt.Errorf("group %q: got %d operations, want %d", group, got, want)
		}
	}
	for label, want := range map[Stability]int{StabilityStable: 127, StabilityNew: 10, StabilityBeta: 12} {
		if got := labels[label]; got != want {
			return fmt.Errorf("stability %q: got %d operations, want %d", label, got, want)
		}
	}
	return nil
}

func (operation Operation) validate() error {
	if operation.OperationID == "" {
		return fmt.Errorf("operation %q: operation_id is required", operation.Anchor)
	}
	if operation.Anchor == "" {
		return fmt.Errorf("operation %q: anchor is required", operation.OperationID)
	}
	if operation.OperationID != operation.Anchor {
		return fmt.Errorf("operation %q: operation_id does not match anchor", operation.Anchor)
	}
	for field, value := range map[string]string{
		"group": operation.Group, "name": operation.Name, "method": operation.Method,
		"path": operation.Path, "subject_binding": operation.SubjectBinding,
		"source": operation.Source,
	} {
		if value == "" {
			return fmt.Errorf("operation %q: %s is required", operation.Anchor, field)
		}
	}
	if operation.Stability != StabilityStable && operation.Stability != StabilityNew && operation.Stability != StabilityBeta {
		return fmt.Errorf("operation %q: invalid stability %q", operation.Anchor, operation.Stability)
	}
	if operation.Method != http.MethodGet && operation.Method != http.MethodHead && operation.Method != http.MethodPost && operation.Method != http.MethodPut && operation.Method != http.MethodPatch && operation.Method != http.MethodDelete {
		return fmt.Errorf("operation %q: invalid method %q", operation.Anchor, operation.Method)
	}
	if len(operation.TokenClasses) == 0 {
		return fmt.Errorf("operation %q: token_classes is required", operation.Anchor)
	}
	for _, tokenClass := range operation.TokenClasses {
		if tokenClass == "" {
			return fmt.Errorf("operation %q: token class is required", operation.Anchor)
		}
	}
	if len(operation.Scopes) == 0 {
		return fmt.Errorf("operation %q: scopes is required", operation.Anchor)
	}
	for _, scope := range operation.Scopes {
		if scope == "" {
			return fmt.Errorf("operation %q: scope is required", operation.Anchor)
		}
	}
	if len(operation.Request.Locations) == 0 {
		return fmt.Errorf("operation %q: request locations are required", operation.Anchor)
	}
	for location, fields := range operation.Request.Locations {
		if location == "" {
			return fmt.Errorf("operation %q: request location is required", operation.Anchor)
		}
		for _, field := range fields {
			if field.Name == "" || field.Type == "" {
				return fmt.Errorf("operation %q: request field is incomplete", operation.Anchor)
			}
		}
	}
	if operation.Response.Format == "" || (!operation.Response.StatusUnknown && len(operation.Response.Status) == 0) {
		return fmt.Errorf("operation %q: response format and status are required", operation.Anchor)
	}
	if operation.Response.StatusUnknown && len(operation.Response.Status) != 0 {
		return fmt.Errorf("operation %q: unknown response status cannot contain numeric statuses", operation.Anchor)
	}
	if operation.Pagination.Shape == "" || operation.Pagination.CursorParameter == "" {
		return fmt.Errorf("operation %q: pagination shape is required", operation.Anchor)
	}
	if len(operation.Fixtures) == 0 {
		return fmt.Errorf("operation %q: fixture references are required", operation.Anchor)
	}
	expectedReplayable := (operation.Method == http.MethodGet || operation.Method == http.MethodHead) && operation.Request.BodyReconstructible
	if operation.Replay.Replayable != expectedReplayable {
		return fmt.Errorf("operation %q: replayable must be %t for %s with body_reconstructible=%t", operation.Anchor, expectedReplayable, operation.Method, operation.Request.BodyReconstructible)
	}
	for _, status := range operation.Response.Status {
		if status == 0 {
			return fmt.Errorf("operation %q: response status is invalid", operation.Anchor)
		}
		if status == http.StatusTooManyRequests && operation.Replay.BucketWaitable {
			return fmt.Errorf("operation %q: endpoint-specific 429 cannot be bucket-waitable", operation.Anchor)
		}
	}
	if operation.Implementation.Anchor != operation.Anchor || operation.Implementation.Selector == "" || operation.Implementation.ServiceType == "" || operation.Implementation.Method == "" || operation.Implementation.Signature == "" || operation.Implementation.RequestType == "" || operation.Implementation.DataType == "" {
		return fmt.Errorf("operation %q: implementation symbol is incomplete", operation.Anchor)
	}
	return nil
}
