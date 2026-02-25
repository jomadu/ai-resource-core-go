package airesource

import (
	"os"
	"path/filepath"
	"testing"
)

// Conformance tests verify that the Go implementation correctly interprets
// AI Resources according to the official specification.
//
// These tests MUST use official fixtures from the ai-resource-spec repository
// (via git submodule at testdata/spec/). They MUST NOT fall back to local
// fixtures or skip tests silently.
//
// Run with: make test-conformance

const (
	specValidDir   = "../../testdata/spec/schema/draft/tests/valid"
	specInvalidDir = "../../testdata/spec/schema/draft/tests/invalid"
)

func checkSpecFixtures(t *testing.T) {
	t.Helper()
	
	if _, err := os.Stat(specValidDir); os.IsNotExist(err) {
		t.Fatal("Conformance test fixtures not found at testdata/spec/\n\n" +
			"The official AI Resource Specification test suite is required for conformance testing.\n\n" +
			"To initialize: make test\n\n" +
			"Or manually: git submodule update --init --recursive")
	}
	
	if _, err := os.Stat(specInvalidDir); os.IsNotExist(err) {
		t.Fatal("Expected structure: schema/draft/tests/valid/ and schema/draft/tests/invalid/")
	}
}

func TestConformance(t *testing.T) {
	checkSpecFixtures(t)
	
	t.Run("ValidCases", func(t *testing.T) {
		files, err := filepath.Glob(filepath.Join(specValidDir, "*.yml"))
		if err != nil {
			t.Fatalf("failed to glob valid fixtures: %v", err)
		}
		if len(files) == 0 {
			t.Fatal("No test fixtures found in testdata/spec/schema/draft/tests/valid/")
		}
		
		for _, f := range files {
			name := filepath.Base(f)
			t.Run(name, func(t *testing.T) {
				_, err := LoadResource(f)
				if err != nil {
					t.Fatalf("expected valid, got error: %v", err)
				}
			})
		}
	})
	
	t.Run("InvalidCases", func(t *testing.T) {
		files, err := filepath.Glob(filepath.Join(specInvalidDir, "*.yml"))
		if err != nil {
			t.Fatalf("failed to glob invalid fixtures: %v", err)
		}
		if len(files) == 0 {
			t.Fatal("No test fixtures found in testdata/spec/schema/draft/tests/invalid/")
		}
		
		for _, f := range files {
			name := filepath.Base(f)
			t.Run(name, func(t *testing.T) {
				_, err := LoadResource(f)
				if err == nil {
					t.Fatalf("expected error, got success")
				}
			})
		}
	})
}
