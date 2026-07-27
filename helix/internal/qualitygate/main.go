package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const defaultSummary = "149 operations, 30 groups, 127 stable, 10 NEW, 12 BETA, 0 missing, 0 extra, 0 unclassified, 0 duplicate mappings"

type testEvent struct {
	Action string `json:"Action"`
	Output string `json:"Output"`
}

func main() {
	mode := flag.String("mode", "", "quality gate mode: test-json or evidence")
	directory := flag.String("dir", ".", "evidence directory")
	logPath := flag.String("log", "", "log file for evidence mode")
	requiredReceipts := flag.Int("required-receipts", 0, "number of task receipts required")
	canonicalSummary := flag.String("canonical-summary", defaultSummary, "exact canonical summary line")
	taskReceipt := flag.String("task13-receipt", "", "task 13 receipt to validate")
	finalHead := flag.String("final-head", "", "expected final HEAD recorded by receipts")
	finalReceipts := flag.String("final-receipts", "", "comma-separated final receipt names")
	flag.Parse()

	var err error
	switch *mode {
	case "test-json":
		err = validateTestJSON(os.Stdin)
	case "evidence":
		err = runEvidence(*directory, *logPath, *requiredReceipts, *canonicalSummary, *taskReceipt, *finalHead, *finalReceipts)
	default:
		err = errors.New("qualitygate: -mode must be test-json or evidence")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func validateTestJSON(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 4096), 4<<20)
	seen := 0
	for scanner.Scan() {
		line := scanner.Text()
		var event testEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return fmt.Errorf("qualitygate: invalid go test JSON: %w", err)
		}
		seen++
		if event.Action == "skip" || event.Action == "fail" {
			return fmt.Errorf("qualitygate: go test event action %q is not allowed", event.Action)
		}
		lower := strings.ToLower(line)
		if strings.Contains(lower, "external-dial") || strings.Contains(lower, "external dial") || strings.Contains(lower, "external_dial") {
			return errors.New("qualitygate: external-dial marker detected")
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("qualitygate: read test JSON: %w", err)
	}
	if seen == 0 {
		return errors.New("qualitygate: empty go test JSON input")
	}
	return nil
}

func validateEvidence(log io.Reader, requiredReceipts int, summary string) error {
	if requiredReceipts < 0 {
		return errors.New("qualitygate: required receipt count cannot be negative")
	}
	if summary == "" {
		return errors.New("qualitygate: canonical summary is required")
	}
	scanner := bufio.NewScanner(log)
	receipts := 0
	summaryFound := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "receipt") {
			receipts++
		}
		if line == summary {
			summaryFound = true
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("qualitygate: read evidence log: %w", err)
	}
	if receipts < requiredReceipts {
		return fmt.Errorf("qualitygate: found %d receipt markers, want at least %d", receipts, requiredReceipts)
	}
	if !summaryFound {
		return errors.New("qualitygate: canonical summary line is missing")
	}
	return nil
}

func runEvidence(directory, logPath string, requiredReceipts int, summary, taskReceipt, finalHead, finalReceipts string) error {
	if logPath == "" {
		return errors.New("qualitygate: -log is required in evidence mode")
	}
	log, err := openRegular(logPath)
	if err != nil {
		return err
	}
	defer log.Close()
	if err := validateEvidence(log, requiredReceipts, summary); err != nil {
		return err
	}
	if requiredReceipts > 0 {
		if err := validateTaskReceipts(directory, requiredReceipts); err != nil {
			return err
		}
	}
	if taskReceipt != "" {
		if _, err := testedCommit(taskReceipt); err != nil {
			return err
		}
	}
	if finalHead != "" || finalReceipts != "" {
		if err := validateFinalReceipts(directory, finalHead, finalReceipts); err != nil {
			return err
		}
	}
	return nil
}

func validateTaskReceipts(directory string, required int) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("qualitygate: read evidence directory: %w", err)
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "task-") || !strings.HasSuffix(entry.Name(), "-receipt.json") {
			continue
		}
		if _, err := readJSON(filepath.Join(directory, entry.Name())); err != nil {
			return err
		}
		count++
	}
	if count < required {
		return fmt.Errorf("qualitygate: found %d task receipts, want at least %d", count, required)
	}
	return nil
}

func validateFinalReceipts(directory, finalHead, names string) error {
	for _, name := range strings.Split(names, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if !strings.HasPrefix(name, "final-") {
			name = "final-" + name + "-receipt.json"
		}
		value, err := readJSON(filepath.Join(directory, name))
		if err != nil {
			return err
		}
		if finalHead != "" && !strings.Contains(string(value), finalHead) {
			return fmt.Errorf("qualitygate: final receipt %q does not contain final HEAD", name)
		}
	}
	return nil
}

func testedCommit(path string) (string, error) {
	value, err := readJSON(path)
	if err != nil {
		return "", err
	}
	var receipt struct {
		TestedCommit string `json:"testedCommit"`
	}
	if err := json.Unmarshal(value, &receipt); err != nil || receipt.TestedCommit == "" {
		return "", fmt.Errorf("qualitygate: receipt %q has no testedCommit", path)
	}
	return receipt.TestedCommit, nil
}

func readJSON(path string) ([]byte, error) {
	file, err := openRegular(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	value, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("qualitygate: read %q: %w", path, err)
	}
	if !json.Valid(value) {
		return nil, fmt.Errorf("qualitygate: invalid JSON receipt %q", path)
	}
	return value, nil
}

func openRegular(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return nil, fmt.Errorf("qualitygate: open %q: %w", path, err)
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("qualitygate: stat %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		file.Close()
		return nil, fmt.Errorf("qualitygate: %q is not a regular file", path)
	}
	return file, nil
}
