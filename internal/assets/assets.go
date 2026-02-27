package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
)

//go:embed spec/schema/draft/*.schema.json
var schemas embed.FS

//go:embed spec/schema/draft/tests/valid/*.yml
var validFixtures embed.FS

//go:embed spec/schema/draft/tests/invalid/*.yml
var invalidFixtures embed.FS

// GetSchema returns the JSON schema for a given kind and version.
func GetSchema(version, kind string) ([]byte, error) {
	path := filepath.Join("spec/schema", version, kind+".schema.json")
	data, err := schemas.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema not found for kind %s (version %s): %w", kind, version, err)
	}
	return data, nil
}

// ValidFixtures returns an fs.FS for valid test fixtures.
func ValidFixtures(version string) fs.FS {
	sub, _ := fs.Sub(validFixtures, filepath.Join("spec/schema", version, "tests/valid"))
	return sub
}

// InvalidFixtures returns an fs.FS for invalid test fixtures.
func InvalidFixtures(version string) fs.FS {
	sub, _ := fs.Sub(invalidFixtures, filepath.Join("spec/schema", version, "tests/invalid"))
	return sub
}
