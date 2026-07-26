package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func LoadManifest() (Manifest, error) {
	return LoadManifestFrom(".")
}

func LoadManifestFrom(root string) (Manifest, error) {
	entries, err := os.ReadDir(filepath.Join(root, "groups"))
	if err != nil {
		return Manifest{}, fmt.Errorf("read manifest groups: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	manifest := Manifest{Schema: 1}
	groupFiles := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		groupFiles++
		path := filepath.Join(root, "groups", entry.Name())
		bytes, err := os.ReadFile(path)
		if err != nil {
			return Manifest{}, fmt.Errorf("read group %s: %w", entry.Name(), err)
		}
		var group GroupFile
		if err := json.Unmarshal(bytes, &group); err != nil {
			return Manifest{}, fmt.Errorf("decode group %s: %w", entry.Name(), err)
		}
		if group.Schema != 1 || group.Group == "" {
			return Manifest{}, fmt.Errorf("group %s: schema and group are required", entry.Name())
		}
		for _, operation := range group.Operations {
			if operation.Group != group.Group {
				return Manifest{}, fmt.Errorf("group %q: operation %q has group %q", group.Group, operation.Anchor, operation.Group)
			}
		}
		manifest.Operations = append(manifest.Operations, group.Operations...)
	}
	if groupFiles != len(expectedGroups) {
		return Manifest{}, fmt.Errorf("manifest group files: got %d, want %d", groupFiles, len(expectedGroups))
	}
	if expected, err := os.ReadFile(filepath.Join(root, "expected-operations.json")); err == nil {
		var descriptor struct {
			Schema          int    `json:"schema"`
			RetrievedAt     string `json:"retrieved_at"`
			ReferenceSHA256 string `json:"reference_sha256"`
		}
		if err := json.Unmarshal(expected, &descriptor); err != nil {
			return Manifest{}, fmt.Errorf("decode expected operation metadata: %w", err)
		}
		manifest.Schema = descriptor.Schema
		manifest.RetrievedAt = descriptor.RetrievedAt
		manifest.ReferenceSHA256 = descriptor.ReferenceSHA256
	}
	if err := manifest.Validate(); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}
