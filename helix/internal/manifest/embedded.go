package manifest

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sync"
)

//go:embed groups/*.json
var embeddedGroups embed.FS

var (
	embeddedOperationsOnce sync.Once
	embeddedOperationMap   map[string]Operation
	embeddedOperationsErr  error
)

func loadEmbeddedOperations() {
	paths, err := fs.Glob(embeddedGroups, "groups/*.json")
	if err != nil {
		embeddedOperationsErr = fmt.Errorf("find embedded manifest groups: %w", err)
		return
	}
	operations := make(map[string]Operation)
	for _, groupPath := range paths {
		data, err := fs.ReadFile(embeddedGroups, path.Clean(groupPath))
		if err != nil {
			embeddedOperationsErr = fmt.Errorf("read embedded manifest group %q: %w", groupPath, err)
			return
		}
		var group GroupFile
		if err := json.Unmarshal(data, &group); err != nil {
			embeddedOperationsErr = fmt.Errorf("decode embedded manifest group %q: %w", groupPath, err)
			return
		}
		for _, operation := range group.Operations {
			if _, exists := operations[operation.Anchor]; exists {
				embeddedOperationsErr = fmt.Errorf("embedded manifest operation %q: duplicate anchor", operation.Anchor)
				return
			}
			operations[operation.Anchor] = operation
		}
	}
	embeddedOperationMap = operations
}

func OperationByAnchor(anchor string) (Operation, error) {
	embeddedOperationsOnce.Do(loadEmbeddedOperations)
	if embeddedOperationsErr != nil {
		return Operation{}, embeddedOperationsErr
	}
	operation, ok := embeddedOperationMap[anchor]
	if !ok {
		return Operation{}, fmt.Errorf("manifest operation %q not found", anchor)
	}
	return operation, nil
}
