# Conformance Testing

## Job to be Done
Verify that the Go implementation correctly interprets AI Resources according to the specification by testing against the official test suite.

## Activities
- Load test fixtures from ai-resource-spec repository
- Run valid examples and verify they pass
- Run invalid examples and verify they fail with appropriate errors
- Test all resource kinds (Prompt, Promptset, Rule, Ruleset)
- Test fragment resolution with various templates
- Test edge cases and boundary conditions
- Report conformance results

## Acceptance Criteria
- [ ] All valid examples from spec repo pass loading and validation
- [ ] All invalid examples from spec repo fail with errors
- [ ] Fragment resolution produces expected output for spec examples
- [ ] Multi-document YAML examples are handled correctly
- [ ] Error messages match expected failure reasons
- [ ] Test suite can be run automatically (CI/CD)
- [ ] Test results clearly indicate pass/fail for each case
- [ ] New spec versions can be tested by updating fixtures

## Data Structures

### TestCase
```go
type TestCase struct {
    Name        string
    Path        string
    ShouldPass  bool
    ExpectedErr string
}
```

**Fields:**
- `Name` - Human-readable test case name
- `Path` - Path to test fixture file
- `ShouldPass` - Whether resource should be valid
- `ExpectedErr` - Expected error message pattern (for invalid cases)

### ConformanceResult
```go
type ConformanceResult struct {
    Total   int
    Passed  int
    Failed  int
    Skipped int
    Cases   []CaseResult
}

type CaseResult struct {
    Name    string
    Status  string
    Message string
}
```

## Algorithm

### Test Suite Execution

1. Discover test fixtures in testdata/
2. Categorize as valid or invalid
3. For each test case:
   - Load resource
   - Validate (schema + semantic)
   - Check result matches expectation
   - Record pass/fail
4. Generate conformance report

**Pseudocode:**
```
function RunConformanceTests():
    results = ConformanceResult{}
    
    valid_cases = discover_fixtures("testdata/valid/")
    invalid_cases = discover_fixtures("testdata/invalid/")
    
    for case in valid_cases:
        result = test_valid_case(case)
        results.add(result)
    
    for case in invalid_cases:
        result = test_invalid_case(case)
        results.add(result)
    
    return results
```

### Valid Case Testing

```
function test_valid_case(case):
    resource, err = LoadResource(case.Path)
    
    if err != nil:
        return CaseResult{
            Name: case.Name,
            Status: "FAIL",
            Message: "expected valid, got error: {err}"
        }
    
    err = Validate(resource)
    
    if err != nil:
        return CaseResult{
            Name: case.Name,
            Status: "FAIL",
            Message: "validation failed: {err}"
        }
    
    return CaseResult{
        Name: case.Name,
        Status: "PASS"
    }
```

### Invalid Case Testing

```
function test_invalid_case(case):
    resource, err = LoadResource(case.Path)
    
    if err == nil:
        err = Validate(resource)
    
    if err == nil:
        return CaseResult{
            Name: case.Name,
            Status: "FAIL",
            Message: "expected error, got success"
        }
    
    if case.ExpectedErr != "" and not matches(err, case.ExpectedErr):
        return CaseResult{
            Name: case.Name,
            Status: "FAIL",
            Message: "wrong error: expected '{case.ExpectedErr}', got '{err}'"
        }
    
    return CaseResult{
        Name: case.Name,
        Status: "PASS"
    }
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Missing test fixture | Skip test with warning |
| Malformed test fixture | Fail test with parse error |
| Test fixture with wrong extension | Fail with unsupported format error |
| Empty test directory | Report 0 tests run |
| Test timeout | Fail test with timeout error |
| Unexpected panic | Catch and fail test gracefully |

## Dependencies

- `resource-loading.md` - Load test fixtures
- `schema-validation.md` - Validate test cases
- `semantic-validation.md` - Validate test cases
- `fragment-resolution.md` - Test fragment resolution
- Test fixtures from ai-resource-spec repository (via git submodule)

## Test Fixture Management

Test fixtures are maintained in the official ai-resource-spec repository and referenced via git submodule:

```bash
# Initial setup
git submodule add https://github.com/org/ai-resource-spec testdata/spec
git submodule update --init --recursive

# Update to latest spec version
cd testdata/spec
git pull origin main
cd ../..
git add testdata/spec
git commit -m "Update spec test fixtures to v1.2.3"
```

**Directory structure:**
```
testdata/
  spec/              # Git submodule → ai-resource-spec repo
    examples/
      valid/
      invalid/
```

**Benefits:**
- Version-pinned test fixtures (reproducible builds)
- Intentional updates when ready to test new spec versions
- Works offline after initial clone
- Standard Git tooling

## Source of Truth

The official AI Resource Specification repository is the canonical source for conformance testing:

**Repository:** `https://github.com/jomadu/ai-resource-spec.git`

**Purpose:** Provides the authoritative test fixtures that all implementations MUST use to verify spec compliance.

**Why Required:**
- **Prevents drift** - Ensures all implementations test against identical fixtures
- **Ensures interoperability** - Resources that pass conformance tests work across all compliant implementations
- **Single source of truth** - Spec repository owns the definition of correctness
- **Version alignment** - Implementation can pin to specific spec versions for stability

**Structure:**
- Valid fixtures: `schema/draft/tests/valid/`
- Invalid fixtures: `schema/draft/tests/invalid/`

Implementations MUST NOT use local test fixtures as a substitute for the official spec repository. Local fixtures may be used for implementation-specific tests, but conformance testing MUST use the official fixtures.

## Makefile-Driven Workflow

The Makefile provides a standard interface for all development tasks, abstracting away submodule management and ensuring consistent workflows:

**Standard Commands:**
- `make test` - Run all tests (automatically initializes submodule if needed)
- `make test-conformance` - Run conformance tests only
- `make build` - Build all packages
- `make lint` - Run linters
- `make update-spec` - Update spec submodule to latest version
- `make clean` - Clean build artifacts
- `make help` - Show all available commands

**Rationale:**
- **Automation** - Developers don't need to manually manage git submodules
- **Consistency** - Same commands work for all developers and CI/CD
- **Discoverability** - `make help` shows available commands
- **Error prevention** - Makefile ensures submodule is initialized before running tests

**CI/CD Integration:**
```bash
# CI/CD should use Makefile commands, not direct go commands
make test
make lint
```

The Makefile handles submodule initialization automatically, so CI/CD configurations don't need explicit `git submodule update --init` commands.

## Submodule Requirements

Conformance testing MUST use the official spec repository via git submodule. This is a hard requirement, not optional.

**Requirements:**
- Submodule MUST be initialized at `testdata/spec/` before running conformance tests
- Conformance tests MUST fail with a clear error if `testdata/spec/` is missing or not initialized
- Tests MUST NOT fall back to local fixtures or skip conformance tests silently
- CI/CD MUST initialize submodules (handled automatically by `make test`)

**Error Handling:**
When `testdata/spec/` is missing or not initialized, tests MUST fail with an error message like:

```
Conformance test fixtures not found at testdata/spec/

The official AI Resource Specification test suite is required for conformance testing.

To initialize: make test

Or manually: git submodule update --init --recursive
```

**Rationale:**
- **Prevents silent failures** - Missing fixtures are caught immediately
- **Clear guidance** - Error message tells developers exactly how to fix the issue
- **No drift** - Impossible to accidentally test against stale or incorrect fixtures
- **Enforces best practice** - Makes the correct workflow the only workflow

## Version Pinning Strategy

Git submodules pin to a specific commit by default. This is the recommended approach for build stability.

**Strategy:**
- Submodule references a specific commit (not a branch)
- Updates are intentional and controlled via `make update-spec`
- Each update is reviewed, tested, and committed explicitly

**Rationale:**
- **Build stability** - Tests are reproducible across time and environments
- **Controlled updates** - New spec versions are adopted intentionally, not automatically
- **Easier debugging** - Known-good commit can be referenced when investigating failures
- **Bisection support** - Can identify which spec change caused a test failure

**Update Workflow:**
1. Run `make update-spec` to fetch latest spec version
2. Review changes: `git diff testdata/spec`
3. Run `make test` to verify compatibility
4. If tests pass, commit the new pin: `git add testdata/spec && git commit -m "Update spec to vX.Y.Z"`
5. If tests fail, investigate and fix implementation or document known issues

**Optional: Testing Against Latest Spec**

For early warning of upcoming spec changes, CI/CD can optionally run a separate job that tests against the latest spec version:

```yaml
# Stable build (required to pass)
- name: Test against pinned spec
  run: make test

# Early warning (allowed to fail)
- name: Test against latest spec
  run: |
    cd testdata/spec
    git pull origin main
    cd ../..
    make test
  continue-on-error: true
```

This provides advance notice of spec changes without blocking builds.

## Failure Policy

Conformance tests MUST fail explicitly when requirements are not met. Silent fallbacks or warnings are not acceptable.

**Failure Conditions:**

| Condition | Behavior | Error Message |
|-----------|----------|---------------|
| `testdata/spec/` missing | FAIL immediately | "Conformance test fixtures not found. Run: make test" |
| `testdata/spec/` empty | FAIL immediately | "Conformance test fixtures directory is empty. Run: git submodule update --init" |
| Invalid directory structure | FAIL immediately | "Expected structure: schema/draft/tests/valid/ and schema/draft/tests/invalid/" |
| No test files found | FAIL immediately | "No test fixtures found in testdata/spec/schema/draft/tests/" |

**Rationale:**
- **No silent failures** - Missing fixtures are always caught
- **Clear error messages** - Developers know exactly what's wrong and how to fix it
- **Fail fast** - Problems are detected immediately, not during test execution
- **Consistent behavior** - Same failure modes across all environments

**Not Allowed:**
- Falling back to local fixtures when submodule is missing
- Skipping conformance tests with a warning
- Continuing with partial test coverage
- Assuming fixtures are optional

## Implementation Mapping

**Source files:**
- `conformance_test.go` - Main conformance test suite
- `testdata/spec/` - Git submodule to ai-resource-spec repository
- `Makefile` - Standard development interface

**Related specs:**
- All specs - Conformance tests verify all functionality

## Examples

### Example 1: Valid Prompt Test

**Input:**
```yaml
# testdata/valid/prompt.yml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: test-prompt
  name: Test Prompt
spec:
  body: "This is a test prompt"
```

**Expected Output:**
```go
// Test passes
result.Status == "PASS"
```

**Verification:**
- Resource loads successfully
- Validation passes
- No errors

### Example 2: Invalid Metadata Test

**Input:**
```yaml
# testdata/invalid/missing-id.yml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  name: "No ID"
spec:
  body: "Test"
```

**Expected Output:**
```go
// Test passes (correctly rejects invalid resource)
result.Status == "PASS"
result.Message contains "metadata.id"
```

**Verification:**
- Resource fails validation
- Error indicates missing ID
- Test correctly expects failure

### Example 3: Fragment Resolution Test

**Input:**
```yaml
# testdata/valid/prompt-with-fragment.yml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: fragment-test
spec:
  fragments:
    greet:
      inputs:
        name:
          type: string
          required: true
      body: "Hello, {{name}}!"
  body:
    - fragment: greet
      inputs:
        name: "World"
```

**Expected Output:**
```go
// Test passes
result.Status == "PASS"

// Resolution test
resolved, _ := ResolveBody(prompt.Spec.Body, prompt.Spec.Fragments)
resolved == "Hello, World!"
```

**Verification:**
- Resource loads and validates
- Fragment resolves correctly
- Output matches expected

### Example 4: Multi-Document Test

**Input:**
```yaml
# testdata/valid/multi-doc.yml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: prompt1
spec:
  body: "First"
---
apiVersion: ai-resource/draft
kind: Rule
metadata:
  id: rule1
spec:
  enforcement: must
  body: "Second"
```

**Expected Output:**
```go
// Test passes
result.Status == "PASS"

// Load test
resources, _ := LoadResources("testdata/valid/multi-doc.yml")
len(resources) == 2
resources[0].Kind == KindPrompt
resources[1].Kind == KindRule
```

**Verification:**
- Both documents load
- Each validates independently
- Correct kinds assigned

### Example 5: Conformance Report

**Input:**
```go
results := RunConformanceTests()
```

**Expected Output:**
```
Conformance Test Results
========================
Total:   50
Passed:  50
Failed:  0
Skipped: 0

Valid Resources:   25/25 passed
Invalid Resources: 25/25 passed

All tests passed!
```

**Verification:**
- All valid cases pass
- All invalid cases correctly fail
- Clear summary provided

## Notes

- Test fixtures are managed via git submodule to ai-resource-spec repository
- Tests should be organized by resource kind and validation type
- Conformance tests are separate from unit tests
- Tests should run quickly (< 1 second for full suite)
- Failed tests should provide enough context to debug
- Tests should be deterministic (no flaky tests)
- Update submodule to test against new spec versions
- Consider using table-driven tests for clarity
- CI/CD should initialize submodules: `git submodule update --init --recursive`

## Known Issues

None.

## Areas for Improvement

- Could add performance benchmarks
- Could add fuzzing tests for robustness
- Could generate conformance report in multiple formats (JSON, HTML)
- Could compare against other implementations
- Could add regression tests for bug fixes
- Could test against large resources for performance
