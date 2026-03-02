package airesource

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/ai-resource-core-go/internal/assets"
)

// Conformance tests verify that the Go implementation correctly interprets
// AI Resources according to the official specification.
//
// These tests use official fixtures from the ai-resource-spec repository
// (embedded via internal/assets package).
//
// Run with: make test-conformance

func TestConformance(t *testing.T) {
	t.Run("ValidCases", func(t *testing.T) {
		validFS := assets.ValidFixtures("draft")
		entries, err := fs.ReadDir(validFS, ".")
		if err != nil {
			t.Fatalf("failed to read valid fixtures: %v", err)
		}
		if len(entries) == 0 {
			t.Fatal("No valid test fixtures found")
		}
		
		for _, entry := range entries {
			if filepath.Ext(entry.Name()) != ".yml" {
				continue
			}
			name := entry.Name()
			t.Run(name, func(t *testing.T) {
				data, err := fs.ReadFile(validFS, name)
				if err != nil {
					t.Fatalf("failed to read fixture: %v", err)
				}
				
				// Write to temp file for LoadResource
				tmpFile := filepath.Join(t.TempDir(), name)
				if err := os.WriteFile(tmpFile, data, 0644); err != nil {
					t.Fatalf("failed to write temp file: %v", err)
				}
				
				_, err = LoadResource(tmpFile)
				if err != nil {
					t.Fatalf("expected valid, got error: %v", err)
				}
			})
		}
	})
	
	t.Run("InvalidCases", func(t *testing.T) {
		invalidFS := assets.InvalidFixtures("draft")
		entries, err := fs.ReadDir(invalidFS, ".")
		if err != nil {
			t.Fatalf("failed to read invalid fixtures: %v", err)
		}
		if len(entries) == 0 {
			t.Fatal("No invalid test fixtures found")
		}
		
		for _, entry := range entries {
			if filepath.Ext(entry.Name()) != ".yml" {
				continue
			}
			name := entry.Name()
			t.Run(name, func(t *testing.T) {
				data, err := fs.ReadFile(invalidFS, name)
				if err != nil {
					t.Fatalf("failed to read fixture: %v", err)
				}
				
				// Write to temp file for LoadResource
				tmpFile := filepath.Join(t.TempDir(), name)
				if err := os.WriteFile(tmpFile, data, 0644); err != nil {
					t.Fatalf("failed to write temp file: %v", err)
				}
				
				_, err = LoadResource(tmpFile)
				if err == nil {
					t.Fatalf("expected error, got success")
				}
			})
		}
	})
}
