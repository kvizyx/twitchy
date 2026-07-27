package testkit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type Fixture struct {
	Name string
	Data []byte
}

func LoadFixture(root, name string) (Fixture, error) {
	clean := filepath.Clean(name)
	if name == "" || filepath.IsAbs(name) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return Fixture{}, fmt.Errorf("testkit: invalid fixture path %q", name)
	}
	path := filepath.Join(root, name)
	file, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return Fixture{}, fmt.Errorf("testkit: open fixture %q: %w", name, err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return Fixture{}, fmt.Errorf("testkit: stat fixture %q: %w", name, err)
	}
	if !info.Mode().IsRegular() {
		return Fixture{}, fmt.Errorf("testkit: fixture %q is not a regular file", name)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return Fixture{}, fmt.Errorf("testkit: read fixture %q: %w", name, err)
	}
	return Fixture{Name: name, Data: data}, nil
}

func LoadText(root, name string) (string, error) {
	fixture, err := LoadFixture(root, name)
	if err != nil {
		return "", err
	}
	return string(fixture.Data), nil
}

func LoadJSON[T any](root, name string) (T, error) {
	var value T
	fixture, err := LoadFixture(root, name)
	if err != nil {
		return value, err
	}
	if err := json.Unmarshal(fixture.Data, &value); err != nil {
		return value, fmt.Errorf("testkit: decode JSON fixture %q: %w", name, err)
	}
	return value, nil
}
