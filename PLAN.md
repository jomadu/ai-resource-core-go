# Refactoring Plan: Move Spec Submodule to internal/assets

## Problem Statement

The current architecture has a "smell" where:
- The spec submodule lives in `testdata/spec/`
- JSON schemas are duplicated in `internal/schema/schemas.go` (hardcoded strings)
- Test fixtures are duplicated in `testdata/valid/` and `testdata/invalid/`
- The submodule location doesn't reflect that it's an internal asset for the package

## Goals

1. Move spec submodule from `testdata/spec/` to `internal/assets/spec/`
2. Eliminate schema duplication by embedding directly from submodule
3. Eliminate test fixture duplication by using only submodule fixtures
4. Create clean `internal/assets` API for accessing schemas and test fixtures
5. Keep all spec-related assets internal (not exposed in public API)

## Architecture Changes

### Before
```
testdata/
  spec/                    # Git submodule
    schema/draft/*.schema.json
    schema/draft/tests/valid/*.yml
    schema/draft/tests/invalid/*.yml
  valid/*.yml              # Duplicated local fixtures
  invalid/*.yml            # Duplicated local fixtures

internal/
  schema/
    schemas.go             # Hardcoded schema strings (duplicated)
    validator.go
```

### After
```
internal/
  assets/
    spec/                  # Git submodule (moved)
      schema/draft/*.schema.json
      schema/draft/tests/valid/*.yml
      schema/draft/tests/invalid/*.yml
    assets.go              # Embed schemas and expose fixtures
  schema/
    validator.go           # Uses internal/assets for schemas
```

## Implementation Steps

### 1. Create internal/assets Package

**File:** `internal/assets/assets.go`

```go
package assets

import (
	"embed"
	"io/fs"
	"path/filepath"
)

//go:embed spec/schema/draft/*.schema.json
var schemas embed.FS

//go:embed spec/schema/draft/tests/valid/*.yml
var validFixtures embed.FS

//go:embed spec/schema/draft/tests/invalid/*.yml
var invalidFixtures embed.FS

// GetSchema returns the JSON schema for a given kind and version
func GetSchema(version, kind string) ([]byte, error) {
	path := filepath.Join("spec/schema", version, kind+".schema.json")
	return schemas.ReadFile(path)
}

// ValidFixtures returns an fs.FS for valid test fixtures
func ValidFixtures(version string) fs.FS {
	sub, _ := fs.Sub(validFixtures, filepath.Join("spec/schema", version, "tests/valid"))
	return sub
}

// InvalidFixtures returns an fs.FS for invalid test fixtures
func InvalidFixtures(version string) fs.FS {
	sub, _ := fs.Sub(invalidFixtures, filepath.Join("spec/schema", version, "tests/invalid"))
	return sub
}
```

**Rationale:**
- Version-aware API (ready for future versions beyond "draft")
- Returns `fs.FS` for test fixtures (standard library interface)
- Returns `[]byte` for schemas (ready for JSON parsing)
- All embedding happens in one place

### 2. Move Git Submodule

**Commands:**
```bash
# Remove old submodule
git submodule deinit -f testdata/spec
git rm -f testdata/spec
rm -rf .git/modules/testdata/spec

# Add submodule at new location
git submodule add https://github.com/jomadu/ai-resource-spec.git internal/assets/spec
git submodule update --init --recursive
```

**Update:** `.gitmodules`
```diff
-[submodule "testdata/spec"]
-	path = testdata/spec
+[submodule "internal/assets/spec"]
+	path = internal/assets/spec
 	url = https://github.com/jomadu/ai-resource-spec.git
```

### 3. Update internal/schema Package

**File:** `internal/schema/validator.go`

**Changes:**
- Import `internal/assets`
- Replace hardcoded schema strings with `assets.GetSchema()` calls
- Remove `internal/schema/schemas.go` entirely

**Before:**
```go
func ValidateSchema(resource interface{}) error {
    schema := getSchemaForKind(resource.Kind) // Uses hardcoded strings
    // ...
}
```

**After:**
```go
import "github.com/yourorg/ai-resource-core-go/internal/assets"

func ValidateSchema(resource interface{}) error {
    schemaBytes, err := assets.GetSchema("draft", resource.Kind)
    if err != nil {
        return fmt.Errorf("schema not found for kind %s: %w", resource.Kind, err)
    }
    schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)
    // ...
}
```

**Delete:** `internal/schema/schemas.go`

### 4. Update Conformance Tests

**File:** `pkg/airesource/conformance_test.go`

**Changes:**
- Import `internal/assets`
- Replace hardcoded paths with `assets.ValidFixtures()` and `assets.InvalidFixtures()`
- Use `fs.ReadDir()` and `fs.ReadFile()` to iterate fixtures

**Before:**
```go
const (
	specValidDir   = "../../testdata/spec/schema/draft/tests/valid"
	specInvalidDir = "../../testdata/spec/schema/draft/tests/invalid"
)

func TestConformance(t *testing.T) {
	checkSpecFixtures(t) // Checks filesystem paths
	files, _ := filepath.Glob(filepath.Join(specValidDir, "*.yml"))
	// ...
}
```

**After:**
```go
import (
	"io/fs"
	"github.com/yourorg/ai-resource-core-go/internal/assets"
)

func TestConformance(t *testing.T) {
	t.Run("ValidCases", func(t *testing.T) {
		fixtures := assets.ValidFixtures("draft")
		entries, err := fs.ReadDir(fixtures, ".")
		if err != nil {
			t.Fatalf("failed to read valid fixtures: %v", err)
		}
		
		for _, entry := range entries {
			name := entry.Name()
			t.Run(name, func(t *testing.T) {
				data, err := fs.ReadFile(fixtures, name)
				if err != nil {
					t.Fatalf("failed to read fixture: %v", err)
				}
				// Parse and validate data
			})
		}
	})
	
	t.Run("InvalidCases", func(t *testing.T) {
		fixtures := assets.InvalidFixtures("draft")
		// Similar pattern
	})
}
```

### 5. Update Makefile

**File:** `Makefile`

**Changes:**
- Update submodule path from `testdata/spec` to `internal/assets/spec`

**Before:**
```makefile
test:
	@if [ ! -d testdata/spec/schema ]; then \
		echo "Initializing submodule..."; \
		git submodule update --init --recursive; \
	fi
	go test ./...

update-spec:
	git submodule update --remote testdata/spec
```

**After:**
```makefile
test:
	@if [ ! -d internal/assets/spec/schema ]; then \
		echo "Initializing submodule..."; \
		git submodule update --init --recursive; \
	fi
	go test ./...

update-spec:
	git submodule update --remote internal/assets/spec
	@echo "Spec updated. Review changes with: git diff internal/assets/spec"
```

### 6. Delete Duplicated Test Fixtures

**Commands:**
```bash
rm -rf testdata/valid
rm -rf testdata/invalid
rm -rf testdata/README.md
rmdir testdata  # Remove if empty
```

**Rationale:**
- Eliminates duplication
- Single source of truth (the submodule)
- Reduces maintenance burden

### 7. Update Specifications

**Files to update:**

#### `specs/conformance-testing.md`

**Changes:**
- Update "Test Fixture Management" section
- Update "Directory structure" example
- Update "Submodule Requirements" section
- Update "Makefile-Driven Workflow" section

**Key updates:**
```diff
-git submodule add https://github.com/org/ai-resource-spec testdata/spec
+git submodule add https://github.com/jomadu/ai-resource-spec.git internal/assets/spec

 **Directory structure:**
-testdata/
-  spec/              # Git submodule → ai-resource-spec repo
+internal/
+  assets/
+    spec/            # Git submodule → ai-resource-spec repo
+      schema/draft/*.schema.json
+      schema/draft/tests/valid/*.yml
+      schema/draft/tests/invalid/*.yml
+    assets.go        # Embed and expose schemas/fixtures
```

#### `specs/schema-validation.md`

**Changes:**
- Update "Implementation Mapping" section
- Update notes about embedding

**Key updates:**
```diff
 **Source files:**
 - `internal/schema/validator.go` - Schema validation logic
-- `internal/schema/schemas.go` - Embedded JSON schemas
+- `internal/assets/assets.go` - Embedded JSON schemas from submodule
 - `pkg/airesource/errors.go` - SchemaError type
```

```diff
-- Schemas are embedded using `go:embed` to avoid runtime file dependencies
-- Schema files should be copied from the ai-resource-spec repository
+- Schemas are embedded directly from the ai-resource-spec submodule using `go:embed`
+- No duplication - schemas are embedded from `internal/assets/spec/schema/`
```

#### `specs/README.md`

**Changes:**
- Update "Conformance Philosophy" section if it mentions testdata paths

### 8. Update Documentation

**Files to update:**

#### `README.md`

**Changes:**
- Update conformance testing section if it mentions testdata paths
- Ensure examples use `make test` (which handles submodule automatically)

#### `AGENTS.md`

**Changes:**
- Update "Specification Definition" section if needed
- Update any references to testdata structure

## Testing Strategy

### Verification Steps

1. **Move submodule:**
   ```bash
   git submodule deinit -f testdata/spec
   git rm -f testdata/spec
   git submodule add https://github.com/jomadu/ai-resource-spec.git internal/assets/spec
   git submodule update --init --recursive
   ```

2. **Create internal/assets package:**
   - Implement `assets.go` with embed directives
   - Verify embeds work: `go build ./internal/assets`

3. **Update internal/schema:**
   - Modify `validator.go` to use `assets.GetSchema()`
   - Delete `schemas.go`
   - Verify builds: `go build ./internal/schema`

4. **Update conformance tests:**
   - Modify `conformance_test.go` to use `assets.ValidFixtures()` and `assets.InvalidFixtures()`
   - Run tests: `make test-conformance`

5. **Update Makefile:**
   - Change paths from `testdata/spec` to `internal/assets/spec`
   - Test: `make test`, `make update-spec`

6. **Delete old fixtures:**
   ```bash
   rm -rf testdata/valid testdata/invalid testdata/README.md
   ```

7. **Run full test suite:**
   ```bash
   make test
   make lint
   make build
   ```

8. **Update specs and docs:**
   - Update all references to old paths
   - Verify documentation accuracy

### Success Criteria

- [ ] All tests pass (`make test`)
- [ ] Conformance tests pass (`make test-conformance`)
- [ ] Linting passes (`make lint`)
- [ ] Build succeeds (`make build`)
- [ ] No hardcoded schema strings in codebase
- [ ] No duplicated test fixtures
- [ ] Submodule at `internal/assets/spec/`
- [ ] `internal/assets` package exposes clean API
- [ ] All specs updated with new paths
- [ ] Documentation reflects new structure

## Benefits

1. **Single Source of Truth:** Schemas and test fixtures come directly from submodule
2. **No Duplication:** Eliminates hardcoded schemas and duplicated fixtures
3. **Better Organization:** Submodule location reflects its purpose (internal asset)
4. **Clean API:** `internal/assets` provides version-aware access to schemas and fixtures
5. **Maintainability:** Updates to spec automatically update schemas and fixtures
6. **Type Safety:** Embedding ensures assets are available at compile time

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Breaking existing tests | Run full test suite after each step |
| Submodule move issues | Follow git submodule best practices, test in clean clone |
| Embed path issues | Verify embed paths match submodule structure |
| Import cycle | Keep `internal/assets` dependency-free (only stdlib) |
| CI/CD breakage | Update Makefile paths, test locally first |

## Rollback Plan

If issues arise:
1. Revert commits in reverse order
2. Re-add submodule at old location: `git submodule add https://github.com/jomadu/ai-resource-spec.git testdata/spec`
3. Restore `internal/schema/schemas.go` from git history
4. Restore old conformance test implementation

## Timeline

Estimated: 2-3 hours

1. **Phase 1:** Move submodule and create `internal/assets` (30 min)
2. **Phase 2:** Update `internal/schema` (20 min)
3. **Phase 3:** Update conformance tests (30 min)
4. **Phase 4:** Update Makefile and delete old fixtures (10 min)
5. **Phase 5:** Update specs and documentation (45 min)
6. **Phase 6:** Testing and verification (30 min)
