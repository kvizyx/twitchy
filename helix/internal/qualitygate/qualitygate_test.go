package main

import (
	"strings"
	"testing"
)

func TestValidateTestJSONRejectsSkipFailAndExternalDial(t *testing.T) {
	input := strings.NewReader(`{"Action":"pass","Package":"example"}
{"Action":"skip","Test":"TestSkipped"}
{"Action":"output","Output":"external-dial: attempted"}
`)
	if err := validateTestJSON(input); err == nil {
		t.Fatal("validateTestJSON() unexpectedly accepted skip/external-dial output")
	}
}

func TestValidateEvidenceRequiresReceiptsAndCanonicalSummary(t *testing.T) {
	log := strings.NewReader("task-12 receipt.json\n149 operations, 30 groups, 127 stable, 10 NEW, 12 BETA, 0 missing, 0 extra, 0 unclassified, 0 duplicate mappings\n")
	if err := validateEvidence(log, 1, "149 operations, 30 groups, 127 stable, 10 NEW, 12 BETA, 0 missing, 0 extra, 0 unclassified, 0 duplicate mappings"); err != nil {
		t.Fatalf("validateEvidence() error = %v", err)
	}
}
