# Test Fixtures

This directory contains test fixtures for the ai-resource-core-go project.

## Directory Structure

```
testdata/
└── README.md          # This file (documentation only)

internal/
  assets/
    spec/              # Git submodule → official AI Resource Specification
      schema/
        draft/
          *.schema.json       # JSON schemas
          tests/
            valid/            # Valid resource examples from spec
            invalid/          # Invalid resource examples from spec
    assets.go          # Embeds schemas and fixtures via go:embed
```

## Test Types

### Conformance Tests

**Location:** `pkg/airesource/conformance_test.go`

**Purpose:** Verify that the Go implementation correctly interprets AI Resources according to the official specification.

**Fixtures:** Uses official test fixtures embedded from `internal/assets/spec/` (git submodule)

**Requirements:**
- Fixtures are embedded in the binary via `go:embed` directives
- Submodule MUST be initialized before building
- Tests access fixtures via `assets.ValidFixtures()` and `assets.InvalidFixtures()`

**Running:**
```bash
make test-conformance
```

### Unit Tests

**Location:** `pkg/airesource/*_test.go`, `internal/**/*_test.go`

**Purpose:** Test implementation-specific behavior, edge cases, and error handling.

**Fixtures:** Uses dynamically generated fixtures (created in test code using `t.TempDir()`)

**Examples:**
- `loader_test.go` - Tests file size limits, timeouts, array size limits
- `validator_test.go` - Tests semantic validation rules
- `fragment_test.go` - Tests fragment resolution and Mustache rendering

**Running:**
```bash
make test
```

## Fixture Management

### Official Fixtures (Conformance)

The official test suite is maintained in the [AI Resource Specification](https://github.com/jomadu/ai-resource-spec) repository and embedded via git submodule at `internal/assets/spec/`.

**Initial setup:**
```bash
git submodule update --init --recursive
```

**Update to latest spec version:**
```bash
make update-spec
```

**Version pinning:**
- Submodule references a specific commit (not a branch)
- Updates are intentional and controlled
- Each update is reviewed, tested, and committed explicitly
- Fixtures are embedded at build time via `go:embed`

### Local Fixtures (Unit Tests)

Unit tests create fixtures dynamically using `t.TempDir()` and `os.WriteFile()`. This approach:
- Keeps test code self-contained
- Makes test intent clear
- Avoids fixture file management
- Allows testing edge cases not covered by spec

## Best Practices

1. **Conformance tests** - Use embedded fixtures via `internal/assets` package
2. **Unit tests** - Generate fixtures in test code, don't commit fixture files
3. **Submodule** - Always initialize before building (for go:embed to work)
4. **Updates** - Use `make update-spec` to update official fixtures
5. **CI/CD** - Use `make test` (handles submodule initialization automatically)

## Troubleshooting

**Error: "No valid/invalid test fixtures found"**
```bash
# Submodule not initialized before build
git submodule update --init --recursive
make build
make test
```

**Submodule out of date:**
```bash
make update-spec
git add internal/assets/spec
git commit -m "Update spec to vX.Y.Z"
```
