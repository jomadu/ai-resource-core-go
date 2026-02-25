# Test Fixtures

This directory contains test fixtures for the ai-resource-core-go project.

## Directory Structure

```
testdata/
├── spec/              # Git submodule → official AI Resource Specification test suite
│   └── schema/
│       └── draft/
│           └── tests/
│               ├── valid/    # Valid resource examples from spec
│               └── invalid/  # Invalid resource examples from spec
└── README.md          # This file
```

## Test Types

### Conformance Tests

**Location:** `pkg/airesource/conformance_test.go`

**Purpose:** Verify that the Go implementation correctly interprets AI Resources according to the official specification.

**Fixtures:** Uses official test fixtures from `testdata/spec/` (git submodule)

**Requirements:**
- MUST use official fixtures from ai-resource-spec repository
- MUST fail explicitly if submodule is not initialized
- MUST NOT fall back to local fixtures

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

The official test suite is maintained in the [AI Resource Specification](https://github.com/jomadu/ai-resource-spec) repository and referenced via git submodule.

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

### Local Fixtures (Unit Tests)

Unit tests create fixtures dynamically using `t.TempDir()` and `os.WriteFile()`. This approach:
- Keeps test code self-contained
- Makes test intent clear
- Avoids fixture file management
- Allows testing edge cases not covered by spec

## Best Practices

1. **Conformance tests** - Use official fixtures only, no local fallbacks
2. **Unit tests** - Generate fixtures in test code, don't commit fixture files
3. **Submodule** - Always initialize before running conformance tests
4. **Updates** - Use `make update-spec` to update official fixtures
5. **CI/CD** - Use `make test` (handles submodule initialization automatically)

## Troubleshooting

**Error: "Conformance test fixtures not found"**
```bash
make test  # Auto-initializes submodule
```

**Error: "No test fixtures found"**
```bash
git submodule update --init --recursive
```

**Submodule out of date:**
```bash
make update-spec
git add testdata/spec
git commit -m "Update spec to vX.Y.Z"
```
