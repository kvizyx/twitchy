package manifest

import (
	"fmt"
	"slices"
)

// groupOperations lists every per-group definition in a fixed, explicit
// order. New groups are appended here by hand along with their
// operations_<group>.go file.
var groupOperations = [][]Operation{
	adsOperations,
	analyticsOperations,
	bitsOperations,
	cclsOperations,
	channelPointsOperations,
	channelsOperations,
	charityOperations,
	chatOperations,
	clipsOperations,
	conduitsOperations,
	entitlementsOperations,
	eventsubOperations,
	extensionsOperations,
	gamesOperations,
	goalsOperations,
	guestStarOperations,
	hypeTrainOperations,
	moderationOperations,
	pollsOperations,
	predictionsOperations,
	raidsOperations,
	scheduleOperations,
	searchOperations,
	streamsOperations,
	subscriptionsOperations,
	tagsOperations,
	teamsOperations,
	usersOperations,
	videosOperations,
	whispersOperations,
}

var orderedOperations = concatGroupOperations()

var (
	operationIndex    map[string]Operation
	operationIndexErr error
)

func init() {
	operationIndex, operationIndexErr = buildOperationIndex(orderedOperations)
}

func concatGroupOperations() []Operation {
	total := 0
	for _, group := range groupOperations {
		total += len(group)
	}
	operations := make([]Operation, 0, total)
	for _, group := range groupOperations {
		operations = append(operations, group...)
	}
	return operations
}

func buildOperationIndex(operations []Operation) (map[string]Operation, error) {
	if len(operations) == 0 {
		return nil, fmt.Errorf("manifest: no operations defined")
	}
	index := make(map[string]Operation, len(operations))
	for _, operation := range operations {
		if operation.Anchor == "" {
			return nil, fmt.Errorf("manifest: operation %q has an empty anchor", operation.OperationID)
		}
		if _, exists := index[operation.Anchor]; exists {
			return nil, fmt.Errorf("manifest: duplicate anchor %q", operation.Anchor)
		}
		index[operation.Anchor] = operation
	}
	return index, nil
}

// OperationByAnchor resolves the registered operation for a Twitch
// documentation anchor.
func OperationByAnchor(anchor string) (Operation, error) {
	if operationIndexErr != nil {
		return Operation{}, operationIndexErr
	}
	operation, ok := operationIndex[anchor]
	if !ok {
		return Operation{}, fmt.Errorf("manifest operation %q not found", anchor)
	}
	return operation, nil
}

// Operations returns every registered operation in registry order.
func Operations() []Operation {
	return slices.Clone(orderedOperations)
}
