package manifest

// Stability marks how trustworthy an endpoint's contract is.
type Stability string

const (
	StabilityStable Stability = "stable"
	StabilityNew    Stability = "NEW"
	StabilityBeta   Stability = "BETA"
)

// TokenClass enumerates the Twitch token kinds an operation accepts.
type TokenClass string

const (
	TokenClassApp       TokenClass = "app"
	TokenClassUser      TokenClass = "user"
	TokenClassExtension TokenClass = "extension"
	TokenClassUnknown   TokenClass = "unknown"
)

// Operation describes one Twitch Helix endpoint: how the client reaches it,
// which credentials it requires, and how the transport may replay it.
//
// Operations are defined as Go literals in operations_<group>.go files via
// defineOperation, which fills every mechanically derived field (OperationID,
// Anchor, Replay.Replayable, Implementation.Anchor/Stability/TestIDs).
type Operation struct {
	OperationID    string
	Anchor         string
	Group          string
	Name           string
	Method         string
	Path           string
	Stability      Stability
	TokenClasses   []TokenClass
	Scopes         []string
	SubjectBinding string
	Request        RequestSpec
	Response       ResponseSpec
	Pagination     PaginationSpec
	Replay         ReplaySpec
	Implementation ImplementationSpec
	// Source is the Twitch documentation URL for hand-editing provenance.
	Source string
}

type RequestSpec struct {
	Locations           map[string][]RequestField
	BodyReconstructible bool
}

type RequestField struct {
	Name     string
	Type     string
	Required bool
}

type ResponseSpec struct {
	Format        string
	Status        []int
	StatusUnknown bool
}

type PaginationSpec struct {
	Shape           string
	CursorParameter string
}

type ReplaySpec struct {
	Replayable     bool
	BucketWaitable bool
}

type ImplementationSpec struct {
	Anchor         string
	Selector       string
	ServiceType    string
	Method         string
	Signature      string
	PagerSignature *string
	RequestType    string
	DataType       string
	Stability      Stability
	TestIDs        []string
}
