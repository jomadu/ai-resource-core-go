package airesource

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testCase struct {
	name       string
	path       string
	shouldPass bool
}

func TestConformance(t *testing.T) {
	validCases := []testCase{
		{"Valid Prompt", "../../testdata/valid/prompt.yml", true},
		{"Valid Promptset", "../../testdata/valid/promptset.yml", true},
		{"Valid Rule", "../../testdata/valid/rule.yml", true},
		{"Valid Ruleset", "../../testdata/valid/ruleset.yml", true},
		{"Valid Fragment", "../../testdata/valid/fragment.yml", true},
		{"Valid Multi-Document", "../../testdata/valid/multi-doc.yml", true},
	}

	invalidCases := []testCase{
		{"Invalid Version", "../../testdata/invalid/invalid-version.yml", false},
		{"Missing ID", "../../testdata/invalid/missing-id.yml", false},
		{"Invalid ID Pattern", "../../testdata/invalid/invalid-id.yml", false},
		{"Missing Body", "../../testdata/invalid/missing-body.yml", false},
		{"Undefined Fragment", "../../testdata/invalid/undefined-fragment.yml", false},
	}

	t.Run("ValidCases", func(t *testing.T) {
		for _, tc := range validCases {
			t.Run(tc.name, func(t *testing.T) {
				if strings.Contains(tc.path, "multi-doc") {
					resources, err := LoadResources(tc.path)
					if err != nil {
						t.Fatalf("expected valid, got error: %v", err)
					}
					if len(resources) != 2 {
						t.Fatalf("expected 2 resources, got %d", len(resources))
					}
				} else {
					_, err := LoadResource(tc.path)
					if err != nil {
						t.Fatalf("expected valid, got error: %v", err)
					}
				}
			})
		}
	})

	t.Run("InvalidCases", func(t *testing.T) {
		for _, tc := range invalidCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := LoadResource(tc.path)
				if err == nil {
					t.Fatalf("expected error, got success")
				}
			})
		}
	})
}

func TestFragmentResolution(t *testing.T) {
	res, err := LoadPrompt("../../testdata/valid/fragment.yml")
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	resolved, err := ResolveBody(res.Spec.Body, res.Spec.Fragments)
	if err != nil {
		t.Fatalf("failed to resolve: %v", err)
	}

	expected := "Hello, World!"
	if resolved != expected {
		t.Errorf("expected %q, got %q", expected, resolved)
	}
}

func TestConformanceDiscovery(t *testing.T) {
	// Check for spec repository fixtures (preferred)
	specValidDir := "../../testdata/spec/examples/valid"
	specInvalidDir := "../../testdata/spec/examples/invalid"

	// Fallback to local fixtures
	localValidDir := "../../testdata/valid"
	localInvalidDir := "../../testdata/invalid"

	checkDir := func(dir string, shouldExist bool) {
		_, err := os.Stat(dir)
		if shouldExist && os.IsNotExist(err) {
			t.Errorf("directory %s should exist", dir)
		}
	}

	countFiles := func(dir string) int {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return 0
		}
		count := 0
		for _, e := range entries {
			if !e.IsDir() && (strings.HasSuffix(e.Name(), ".yml") || strings.HasSuffix(e.Name(), ".yaml")) {
				count++
			}
		}
		return count
	}

	// Try spec fixtures first
	specValidCount := countFiles(specValidDir)
	specInvalidCount := countFiles(specInvalidDir)

	if specValidCount > 0 || specInvalidCount > 0 {
		t.Logf("Using spec repository fixtures: %d valid, %d invalid", specValidCount, specInvalidCount)
		checkDir(specValidDir, true)
		checkDir(specInvalidDir, true)
		if specValidCount == 0 {
			t.Error("no valid test fixtures found in spec repository")
		}
		if specInvalidCount == 0 {
			t.Error("no invalid test fixtures found in spec repository")
		}
		return
	}

	// Fall back to local fixtures
	t.Log("Spec repository not found, using local fixtures")
	checkDir(localValidDir, true)
	checkDir(localInvalidDir, true)

	validCount := countFiles(localValidDir)
	invalidCount := countFiles(localInvalidDir)

	if validCount == 0 {
		t.Error("no valid test fixtures found")
	}
	if invalidCount == 0 {
		t.Error("no invalid test fixtures found")
	}

	t.Logf("Found %d valid and %d invalid test fixtures", validCount, invalidCount)
}

func TestAllResourceKinds(t *testing.T) {
	tests := []struct {
		name string
		path string
		kind Kind
	}{
		{"Prompt", "../../testdata/valid/prompt.yml", KindPrompt},
		{"Promptset", "../../testdata/valid/promptset.yml", KindPromptset},
		{"Rule", "../../testdata/valid/rule.yml", KindRule},
		{"Ruleset", "../../testdata/valid/ruleset.yml", KindRuleset},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := LoadResource(tt.path)
			if err != nil {
				t.Fatalf("failed to load: %v", err)
			}
			if res.Kind != tt.kind {
				t.Errorf("expected kind %s, got %s", tt.kind, res.Kind)
			}
		})
	}
}

func TestMultiDocumentLoading(t *testing.T) {
	resources, err := LoadResources("../../testdata/valid/multi-doc.yml")
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if len(resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(resources))
	}

	if resources[0].Kind != KindPrompt {
		t.Errorf("expected first resource to be Prompt, got %s", resources[0].Kind)
	}

	if resources[1].Kind != KindRule {
		t.Errorf("expected second resource to be Rule, got %s", resources[1].Kind)
	}
}

func TestConformanceReport(t *testing.T) {
	var total, passed, failed int

	runTest := func(path string, shouldPass bool) {
		total++
		_, err := LoadResource(path)
		if shouldPass && err == nil {
			passed++
		} else if !shouldPass && err != nil {
			passed++
		} else {
			failed++
		}
	}

	validFiles, _ := filepath.Glob("../../testdata/valid/*.yml")
	for _, f := range validFiles {
		if !strings.Contains(f, "multi-doc") {
			runTest(f, true)
		}
	}

	invalidFiles, _ := filepath.Glob("../../testdata/invalid/*.yml")
	for _, f := range invalidFiles {
		runTest(f, false)
	}

	t.Logf("Conformance Results: Total=%d Passed=%d Failed=%d", total, passed, failed)

	if failed > 0 {
		t.Errorf("%d conformance tests failed", failed)
	}
}
