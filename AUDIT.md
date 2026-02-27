# Specification to Implementation Gap Analysis

**Date:** 2026-02-26  
**Iteration:** 1 of 5  
**Status:** Bootstrap phase - Implementation partially complete

## Executive Summary

The ai-resource-core-go implementation has made significant progress with core functionality in place: resource loading, schema validation, semantic validation, fragment resolution, and conformance testing infrastructure. However, **6 specification gaps** remain that affect validation completeness and semantic correctness.

**Critical Finding:** Fragment input validation exists but is not integrated into the semantic validation pipeline, allowing invalid inputs to pass validation and fail later during resolution.

## Gaps Found

### Spec-to-Implementation Gaps (Specified but Not Implemented)

#### 1. Fragment Input Validation Not Integrated [HIGH PRIORITY]

**Specification:** `fragment-input-validation.md`, `semantic-validation.md`

**Requirement:**
- Fragment input validation (type checking, required fields, defaults) must run during semantic validation phase
- ValidateInputs should be called for each fragment reference in validateFragmentRef

**Current State:**
- `ValidateInputs()` function exists in `pkg/airesource/fragment.go` (line 193)
- Function is NOT called in `pkg/airesource/validator.go`
- Fragment references are checked for existence but inputs are not validated

**Impact:**
- Invalid fragment inputs pass validation
- Errors occur later during resolution instead of validation
- Poor user experience - late-stage error detection
- Type mismatches, missing required inputs, undefined inputs not caught

**Example:**
```yaml
# This should fail validation but currently passes
spec:
  fragments:
    read:
      inputs:
        path: {type: string, required: true}
      body: "{{path}}"
  body:
    - fragment: read
      inputs: {count: 123}  # Wrong input name, wrong type
```

**Recommended Action:**
Modify `validateFragmentRef()` in `pkg/airesource/validator.go` to call `ValidateInputs()` for each fragment reference.

---

#### 2. Collection Key Pattern Validation Missing [HIGH PRIORITY]

**Specification:** `semantic-validation.md`

**Requirement:**
- Fragment map keys must match pattern `^[a-zA-Z0-9_-]+$`
- Promptset prompt keys must match pattern
- Ruleset rule keys must match pattern

**Current State:**
- Only `metadata.id` is validated against pattern (validator.go line 16)
- Collection keys (fragments, prompts, rules) are not validated

**Impact:**
- Invalid keys like "my fragment!" or "rule#1" are accepted
- May cause issues in downstream systems expecting valid identifiers
- Inconsistent with spec requirements

**Example:**
```yaml
# This should fail but currently passes
spec:
  fragments:
    "invalid key!": {...}  # Should be rejected
    "another bad key": {...}  # Should be rejected
```

**Recommended Action:**
Add key validation loops in:
- `validatePromptSpec()` - validate fragment keys
- `validatePromptsetSpec()` - validate fragment keys and prompt keys
- `validateRuleSpec()` - validate fragment keys
- `validateRulesetSpec()` - validate fragment keys and rule keys

---

#### 3. Empty Collection Validation Missing [MEDIUM PRIORITY]

**Specification:** `semantic-validation.md`

**Requirement:**
- Promptset must have at least one prompt
- Ruleset must have at least one rule

**Current State:**
- `validatePromptsetSpec()` does not check `len(spec.Prompts) > 0`
- `validateRulesetSpec()` does not check `len(spec.Rules) > 0`

**Impact:**
- Empty promptsets/rulesets are accepted
- Resources that define no behavior pass validation
- Violates spec semantic rules

**Example:**
```yaml
# This should fail but currently passes
apiVersion: ai-resource/draft
kind: Promptset
metadata:
  id: empty-set
spec:
  prompts: {}  # Should require at least one
```

**Recommended Action:**
Add length checks:
```go
if len(spec.Prompts) == 0 {
    errors = append(errors, &ValidationError{
        Field: "spec.prompts",
        Message: "promptset must have at least one prompt",
    })
}
```

---

#### 4. Scope Glob Pattern Validation Missing [MEDIUM PRIORITY]

**Specification:** `semantic-validation.md`

**Requirement:**
- Validate scope.files glob patterns are well-formed
- Catch syntax errors in glob patterns during validation

**Current State:**
- `validateRuleSpec()` does not validate scope patterns
- Malformed globs like `[invalid` are accepted

**Impact:**
- Invalid glob patterns accepted during validation
- Errors occur when patterns are evaluated at runtime
- Poor error messages for users

**Example:**
```yaml
# This should fail but currently passes
spec:
  scope:
    - files: ["[invalid-glob"]  # Malformed pattern
```

**Recommended Action:**
Add glob validation using `filepath.Match()` or similar:
```go
for _, entry := range spec.Scope {
    for _, pattern := range entry.Files {
        if _, err := filepath.Match(pattern, ""); err != nil {
            errors = append(errors, &ValidationError{
                Field: "spec.scope.files",
                Message: fmt.Sprintf("invalid glob pattern: %s", pattern),
            })
        }
    }
}
```

---

#### 5. Rule Priority Default Not Applied [LOW PRIORITY]

**Specification:** `core-types.md`

**Requirement:**
- "Rule priority defaults to 100 when not specified"
- Default should be applied during parsing

**Current State:**
- No default application visible in loader or unmarshaling
- Priority field is `int` which defaults to 0

**Impact:**
- Rules without explicit priority get 0 instead of 100
- Affects rule ordering/sorting behavior
- Violates spec default behavior

**Recommended Action:**
Apply default in one of:
1. Custom UnmarshalYAML for RuleItem
2. Post-load processing in LoadRule/LoadRuleset
3. Change Priority to `*int` and apply default when nil

---

#### 6. Validation Pipeline Fail-Fast Behavior [LOW PRIORITY]

**Specification:** `error-handling.md`

**Requirement:**
- Schema validation → Semantic validation (fail-fast between phases)
- Within each phase, collect all errors
- Between phases, stop on first failure

**Current State:**
- Need to verify `loadResourceSync()` implements this correctly
- Code structure suggests it does, but not explicitly documented

**Impact:**
- If not implemented: cascading errors, confusing error messages
- If implemented: documentation/clarity issue only

**Recommended Action:**
1. Verify loader.go implements fail-fast between phases
2. Add code comments documenting pipeline behavior
3. Add test cases verifying fail-fast behavior

---

### Implementation-to-Spec Gaps (Implemented but Not Specified)

**None identified.** All implemented functionality appears to be specified.

---

## Impact Assessment

### By Priority

| Priority | Count | Description |
|----------|-------|-------------|
| HIGH     | 2     | Fragment input validation, collection key validation |
| MEDIUM   | 2     | Empty collections, glob pattern validation |
| LOW      | 2     | Priority default, pipeline documentation |

### By Category

| Category | Count | Description |
|----------|-------|-------------|
| Validation | 4 | Missing validation checks |
| Defaults | 1 | Missing default value application |
| Documentation | 1 | Unclear implementation behavior |

### User-Facing Impact

**High Impact (Gaps #1, #2):**
- Invalid resources pass validation
- Errors occur at wrong stage (resolution instead of validation)
- Confusing error messages
- Poor developer experience

**Medium Impact (Gaps #3, #4):**
- Semantically invalid resources accepted
- Runtime errors instead of validation errors
- Reduced reliability

**Low Impact (Gaps #5, #6):**
- Incorrect default behavior
- Documentation clarity
- Edge case handling

---

## Recommended Actions

### Immediate (High Priority)

1. **Integrate fragment input validation**
   - File: `pkg/airesource/validator.go`
   - Function: `validateFragmentRef()`
   - Action: Call `ValidateInputs()` for each fragment reference
   - Estimated effort: 30 minutes

2. **Add collection key validation**
   - File: `pkg/airesource/validator.go`
   - Functions: All validate*Spec functions
   - Action: Validate all map keys match `^[a-zA-Z0-9_-]+$`
   - Estimated effort: 1 hour

### Short-term (Medium Priority)

3. **Add empty collection checks**
   - File: `pkg/airesource/validator.go`
   - Functions: `validatePromptsetSpec()`, `validateRulesetSpec()`
   - Action: Require len > 0 for prompts/rules maps
   - Estimated effort: 15 minutes

4. **Add glob pattern validation**
   - File: `pkg/airesource/validator.go`
   - Function: `validateRuleSpec()`
   - Action: Validate scope.files patterns
   - Estimated effort: 30 minutes

### Long-term (Low Priority)

5. **Apply rule priority default**
   - File: `pkg/airesource/rule.go` or `loader.go`
   - Action: Apply default value of 100 when priority not specified
   - Estimated effort: 30 minutes

6. **Document validation pipeline**
   - File: `pkg/airesource/loader.go`
   - Action: Add comments explaining fail-fast behavior
   - Estimated effort: 15 minutes

---

## Testing Recommendations

For each gap fix:

1. Add test case with invalid input that should be rejected
2. Verify error message is clear and actionable
3. Add test case with valid input that should pass
4. Update conformance tests if official fixtures exist

Example test structure:
```go
func TestFragmentInputValidation(t *testing.T) {
    // Should reject: wrong input type
    // Should reject: missing required input
    // Should reject: undefined input
    // Should pass: valid inputs with defaults
}
```

---

## Conformance Status

**Overall:** Partial conformance

**Passing:**
- Resource loading (YAML, JSON, multi-doc)
- Schema validation
- Fragment resolution
- Error types
- Core type definitions

**Failing:**
- Semantic validation completeness (gaps #1-4)
- Default value application (gap #5)

**Recommendation:** Address high-priority gaps before claiming full spec conformance.

---

## Operational Notes

**Last Verified:** 2026-02-26

**Bootstrap Status:**
- Go module initialized ✓
- Core implementation complete ✓
- Conformance test infrastructure ready ✓
- Semantic validation incomplete (gaps identified)

**Next Steps:**
1. Fix high-priority validation gaps (#1, #2)
2. Run conformance tests to verify fixes
3. Address medium-priority gaps (#3, #4)
4. Document validation pipeline behavior
5. Apply low-priority fixes (#5, #6)

---

## Appendix: Validation Coverage Matrix

| Validation Rule | Spec Location | Implemented | Notes |
|----------------|---------------|-------------|-------|
| metadata.id pattern | semantic-validation.md | ✓ | Working |
| Fragment exists | semantic-validation.md | ✓ | Working |
| Fragment input types | fragment-input-validation.md | ✗ | Gap #1 |
| Fragment input required | fragment-input-validation.md | ✗ | Gap #1 |
| Fragment input defaults | fragment-input-validation.md | ✗ | Gap #1 |
| Undefined inputs | fragment-input-validation.md | ✗ | Gap #1 |
| Fragment key pattern | semantic-validation.md | ✗ | Gap #2 |
| Prompt key pattern | semantic-validation.md | ✗ | Gap #2 |
| Rule key pattern | semantic-validation.md | ✗ | Gap #2 |
| Promptset min size | semantic-validation.md | ✗ | Gap #3 |
| Ruleset min size | semantic-validation.md | ✗ | Gap #3 |
| Scope glob syntax | semantic-validation.md | ✗ | Gap #4 |
| Body format | semantic-validation.md | ✓ | Working |
| Schema structure | schema-validation.md | ✓ | Working |
| API version | resource-loading.md | ✓ | Working |

**Coverage:** 7/15 validation rules fully implemented (47%)

---

## Conclusion

The implementation has strong foundations with core functionality working correctly. The identified gaps are primarily in semantic validation completeness. Addressing the 2 high-priority gaps will significantly improve validation robustness and user experience. All gaps are fixable with targeted, minimal code changes.

**Estimated Total Effort:** 3-4 hours to address all gaps.
