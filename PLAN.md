# Implementation Plan: Spec Drift Correction

**Date:** 2026-02-24  
**Status:** Published to TODO.json (2026-02-24T17:00:02-08:00)  
**Audit Reference:** AUDIT.md (2026-02-24)

## Overview

This plan addresses 6 gaps identified between specification and implementation. Gaps are prioritized by user impact and system stability.

## Priority 1: Critical (Blocks Core Functionality)

### Gap 1: Mustache Template Rendering

**Problem:** Fragment resolution cannot render templates with variables, conditionals, or iteration.

**Spec:** `specs/fragment-resolution.md`

**Tasks:**
1. Add Mustache dependency: `go get github.com/cbroglie/mustache`
2. Create `internal/template/renderer.go`:
   - `Render(template string, context map[string]interface{}) (string, error)`
   - Wrap mustache library with error context
3. Update `ResolveBody()` in `pkg/airesource/fragment.go`:
   - Replace placeholder with actual Mustache rendering
   - Pass fragment inputs as context
   - Wrap template errors with fragment ID
4. Add tests in `pkg/airesource/fragment_test.go`:
   - Variable substitution: `{{name}}`
   - Conditionals: `{{#show}}...{{/show}}`
   - Array iteration: `{{#items}}...{{/items}}`
   - Missing variables (should render empty)
   - Template syntax errors

**Acceptance:**
- All fragment resolution tests pass
- `go test ./...` passes
- Examples from `specs/fragment-resolution.md` work

**Estimated Effort:** 2-3 hours

---

## Priority 2: High (Spec Compliance)

### Gap 2: Fragment Input Validation via JSON Schema

**Problem:** Implementation uses manual type checking instead of JSON Schema validation as specified.

**Spec:** `specs/fragment-input-validation.md`

**Tasks:**
1. Refactor `ValidateInputs()` in `pkg/airesource/fragment.go`:
   - Use `BuildSchemaFromInputs()` to generate schema
   - Call `internal/schema/validator.go` for validation
   - Remove manual `validateType()` function
   - Keep `applyDefaults()` before validation
2. Update error handling:
   - Convert JSON Schema errors to InputError format
   - Preserve error messages for undefined inputs
3. Update tests in `pkg/airesource/fragment_test.go`:
   - Verify JSON Schema validation behavior
   - Test nested structures (arrays, objects)
   - Test additionalProperties: false enforcement

**Acceptance:**
- Validation uses JSON Schema library
- All existing tests pass
- Error messages match spec requirements
- No manual type checking code remains

**Estimated Effort:** 2-3 hours

---

### Gap 3: Conformance Test Suite Structure

**Problem:** No official test fixtures from spec repository.

**Spec:** `specs/conformance-testing.md`

**Tasks:**
1. Add git submodule (if spec repo exists):
   ```bash
   git submodule add <spec-repo-url> testdata/spec
   git submodule update --init --recursive
   ```
2. Update `conformance_test.go`:
   - Discover fixtures from `testdata/spec/examples/`
   - Run valid examples (expect pass)
   - Run invalid examples (expect fail)
   - Generate conformance report
3. Document in README.md:
   - How to initialize submodules
   - How to update test fixtures
   - How to run conformance tests

**Acceptance:**
- Submodule added and initialized
- Tests discover and run fixtures automatically
- Conformance report generated
- CI/CD documentation updated

**Estimated Effort:** 1-2 hours

**Note:** If spec repository doesn't exist yet, defer this task and document in AUDIT.md.

---

## Priority 3: Medium (Safety & Robustness)

### Gap 4: LoadOptions Completeness

**Problem:** Missing safety limits for arrays, nesting depth, and timeouts.

**Spec:** `specs/resource-loading.md`

**Tasks:**
1. Audit `pkg/airesource/options.go`:
   - Verify `MaxFileSize` exists
   - Add `MaxArraySize` (default: 10000)
   - Add `MaxNestingDepth` (default: 100)
   - Add `Timeout` (default: 30s)
2. Add functional options:
   - `WithMaxArraySize(size int)`
   - `WithMaxNestingDepth(depth int)`
   - `WithTimeout(timeout time.Duration)`
3. Implement enforcement in `pkg/airesource/loader.go`:
   - Check array sizes during parsing
   - Check nesting depth during parsing
   - Apply timeout to file operations
4. Add tests in `pkg/airesource/loader_test.go`:
   - Test array size limit violation
   - Test nesting depth limit violation
   - Test timeout behavior

**Acceptance:**
- All LoadOptions fields present
- Limits enforced during loading
- Tests verify limit violations fail gracefully
- `go test ./...` passes

**Estimated Effort:** 2-3 hours

---

## Priority 4: Low (Code Organization)

### Gap 5: internal/types/ Implementation

**Problem:** Empty directory, validation logic in public API.

**Spec:** `AGENTS.md` - "internal/types/*.go - Type validation"

**Tasks:**
1. Create `internal/types/validation.go`:
   - Move type-specific validation helpers from `pkg/airesource/validator.go`
   - Keep public API minimal
2. Update imports in `pkg/airesource/validator.go`
3. Run tests to verify no regressions

**Acceptance:**
- Type validation helpers in `internal/types/`
- Public API remains minimal
- All tests pass

**Estimated Effort:** 1 hour

**Note:** Low priority - current implementation works, this is purely organizational.

---

### Gap 6: internal/template/ Implementation

**Problem:** Empty directory, template logic needed for Gap 1.

**Spec:** `AGENTS.md` - "internal/template/*.go - Mustache rendering"

**Tasks:**
- Addressed by Gap 1 implementation
- Create `internal/template/renderer.go` as part of Mustache integration

**Acceptance:**
- Covered by Gap 1 acceptance criteria

**Estimated Effort:** Included in Gap 1

---

## Execution Order

1. **Gap 1** (P1) - Mustache rendering - IMMEDIATE
2. **Gap 2** (P2) - JSON Schema validation - After Gap 1
3. **Gap 3** (P2) - Conformance tests - After Gap 1 & 2
4. **Gap 4** (P3) - LoadOptions - After Gap 1 & 2
5. **Gap 5** (P4) - Code organization - After all functional gaps

**Rationale:**
- Gap 1 blocks core functionality - must be first
- Gap 2 affects Gap 1 testing - should be second
- Gap 3 validates Gaps 1 & 2 - should be third
- Gap 4 is independent safety feature - can be parallel or after
- Gap 5 is pure refactoring - last to avoid churn

---

## Validation Checkpoints

After each gap is addressed:

1. Run `go test ./...` - All tests must pass
2. Run `go vet ./...` - No warnings
3. Run `go build ./...` - Clean build
4. Update AUDIT.md - Mark gap as resolved
5. Commit with reference to gap number

---

## Success Criteria

**Complete when:**
- [ ] All P1 gaps resolved
- [ ] All P2 gaps resolved
- [ ] All tests pass
- [ ] No spec/implementation divergence
- [ ] AUDIT.md updated with resolution status

**Partial success:**
- P1 + P2 resolved = Core functionality complete
- P3 resolved = Production-ready
- P4 resolved = Fully organized codebase

---

## Risk Mitigation

**Risk:** Mustache library doesn't match spec behavior  
**Mitigation:** Test against spec examples first, adjust if needed

**Risk:** JSON Schema validation changes error messages  
**Mitigation:** Update tests to match new error format

**Risk:** Spec repository doesn't exist yet  
**Mitigation:** Document in AUDIT.md, create local fixtures as interim

**Risk:** Breaking changes to public API  
**Mitigation:** Keep API changes minimal, add deprecation warnings if needed

---

## Notes

- This plan is based on AUDIT.md dated 2026-02-24
- Estimated efforts are for a single developer
- Tasks can be parallelized where noted
- Update this plan if new gaps are discovered
- Archive this plan when all gaps resolved
