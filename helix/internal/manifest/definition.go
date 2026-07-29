package manifest

import "net/http"

// operationsForGroup stamps the Twitch documentation group onto every
// operation defined in a per-group file.
func operationsForGroup(group string, operations ...Operation) []Operation {
	for index := range operations {
		operations[index].Group = group
	}
	return operations
}

// defineOperation fills every field that is mechanically derived from the
// anchor or from other explicit fields, so per-group definitions only carry
// facts that can vary independently.
func defineOperation(anchor string, operation Operation) Operation {
	operation.OperationID = anchor
	operation.Anchor = anchor
	operation.Replay.Replayable = (operation.Method == http.MethodGet || operation.Method == http.MethodHead) &&
		operation.Request.BodyReconstructible
	operation.Implementation.Anchor = anchor
	operation.Implementation.Stability = operation.Stability
	operation.Implementation.TestIDs = []string{
		"TestManifestConformance/happy/" + anchor,
		"TestManifestConformance/negative/" + anchor,
	}
	return operation
}

func pagerSignature(signature string) *string {
	return &signature
}
