package assets

import (
	"io/fs"
	"testing"
)

func TestGetSchema(t *testing.T) {
	tests := []struct {
		version string
		kind    string
		wantErr bool
	}{
		{"draft", "Prompt", false},
		{"draft", "Promptset", false},
		{"draft", "Rule", false},
		{"draft", "Ruleset", false},
		{"draft", "Unknown", true},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			data, err := GetSchema(tt.version, tt.kind)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(data) == 0 {
				t.Errorf("GetSchema() returned empty data")
			}
		})
	}
}

func TestValidFixtures(t *testing.T) {
	fixtures := ValidFixtures("draft")
	entries, err := fs.ReadDir(fixtures, ".")
	if err != nil {
		t.Fatalf("ValidFixtures() failed to read: %v", err)
	}
	if len(entries) == 0 {
		t.Error("ValidFixtures() returned no files")
	}
}

func TestInvalidFixtures(t *testing.T) {
	fixtures := InvalidFixtures("draft")
	entries, err := fs.ReadDir(fixtures, ".")
	if err != nil {
		t.Fatalf("InvalidFixtures() failed to read: %v", err)
	}
	if len(entries) == 0 {
		t.Error("InvalidFixtures() returned no files")
	}
}
