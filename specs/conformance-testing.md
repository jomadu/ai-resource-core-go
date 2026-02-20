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

## Implementation Mapping

**Source files:**
- `conformance_test.go` - Main conformance test suite
- `testdata/spec/` - Git submodule to ai-resource-spec repository

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
