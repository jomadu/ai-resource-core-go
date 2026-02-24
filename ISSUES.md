# Design Issues Report

**Generated:** 2026-02-23  
**Status:** Under Review

## Critical Design Issues

### 1. Fragment Nesting Limitation

**Issue:** Specs explicitly state "Fragments cannot reference other fragments (no nesting)" in `fragment-resolution.md`.

**Problem:** This severely limits composability. Real-world use cases often need:
- Shared sub-components across fragments
- Hierarchical prompt structures
- DRY principle for common patterns

**Options:**

**A. Keep flat (current design)**
- Pros: Simple implementation, no circular dependency concerns, predictable resolution
- Cons: Forces duplication, limits expressiveness, awkward for complex prompts
- Tradeoff: Simplicity vs. power

**B. Allow single-level nesting**
- Pros: Enables basic composition, manageable complexity
- Cons: Still limiting, "why not more?" questions
- Tradeoff: Moderate complexity for moderate power

**C. Allow unlimited nesting with cycle detection**
- Pros: Maximum flexibility, natural composition patterns
- Cons: Complex implementation, potential performance issues, harder debugging
- Tradeoff: Power vs. complexity

**Recommendation:** **Option B** - Allow fragments to reference other fragments, but only resolve one level deep. This gives 80% of the benefit with 20% of the complexity. Add cycle detection for safety.

---

### 2. Input Validation Timing Ambiguity

**Issue:** `semantic-validation.md` says "Fragment input validation checks existence and required status, not types (types validated during resolution)" but `fragment-input-validation.md` describes full type validation.

**Problem:** Unclear when type validation happens. This affects:
- Error reporting quality (fail early vs. late)
- Performance (validate once vs. multiple times)
- API design (separate validation step?)

**Options:**

**A. Validate during semantic validation (early)**
- Pros: Fail fast, all errors reported together, validate once
- Cons: Couples validation phases, harder to test independently
- Tradeoff: User experience vs. architectural purity

**B. Validate during resolution (late)**
- Pros: Clean separation, only validate what's used
- Cons: Errors appear during rendering, partial validation
- Tradeoff: Architecture vs. usability

**C. Validate in both places**
- Pros: Comprehensive checking, defense in depth
- Cons: Redundant work, maintenance burden
- Tradeoff: Safety vs. performance

**Recommendation:** **Option A** - Validate all fragment inputs during semantic validation. Users want to know about ALL errors upfront, not discover them during resolution. Make validation a distinct phase that runs after schema validation.

---

### 3. Body Type Ambiguity

**Issue:** `Body` is defined as `interface{}` with runtime type checking for string vs. array.

**Problem:** 
- No compile-time safety
- Requires type assertions everywhere
- Error-prone for API consumers
- Unclear from type signature what's valid

**Options:**

**A. Keep interface{} (current)**
- Pros: Matches YAML flexibility, simple unmarshaling
- Cons: No type safety, runtime errors, poor API ergonomics
- Tradeoff: Parsing simplicity vs. usage safety

**B. Use discriminated union type**
```go
type Body struct {
    String *string
    Array  []BodyItem
}
```
- Pros: Explicit, type-safe, self-documenting
- Cons: More verbose, requires custom unmarshaling
- Tradeoff: Safety vs. verbosity

**C. Use separate types per resource**
```go
type PromptBody interface { isPromptBody() }
type StringBody string
type ArrayBody []BodyItem
```
- Pros: Type-safe, idiomatic Go, extensible
- Cons: More complex, requires type switches
- Tradeoff: Type safety vs. simplicity

**Recommendation:** **Option B** - Use explicit struct with nullable fields. This makes the API self-documenting and catches errors at compile time. The verbosity is worth it for a public API.

---

### 4. Error Aggregation Strategy

**Issue:** `MultiError` returns all errors, but specs don't clarify how many errors to collect before stopping.

**Problem:**
- Validating everything can produce overwhelming error lists
- Some errors make subsequent validation meaningless
- No guidance on error limits or prioritization

**Options:**

**A. Collect all errors (current)**
- Pros: Complete picture, fix everything at once
- Cons: Can be overwhelming (100+ errors), some are cascading
- Tradeoff: Completeness vs. usability

**B. Fail fast on critical errors**
- Pros: Clear, focused feedback
- Cons: Multiple fix-test cycles, slower iteration
- Tradeoff: Clarity vs. efficiency

**C. Collect with limits and prioritization**
- Pros: Balanced, shows most important issues first
- Cons: More complex, requires error classification
- Tradeoff: Sophistication vs. implementation cost

**Recommendation:** **Option C** - Collect up to 10 errors per validation phase, prioritize by severity (schema > semantic > input). This prevents overwhelming users while still showing multiple issues. Add a flag for "show all errors" mode.

---

### 5. No Streaming/Incremental Loading

**Issue:** `LoadResources` loads entire multi-document file into memory.

**Problem:**
- 10MB limit is arbitrary
- Large resource bundles (1000+ resources) hit memory limits
- No way to process resources incrementally

**Options:**

**A. Keep in-memory loading (current)**
- Pros: Simple, works for 99% of cases
- Cons: Doesn't scale to large bundles
- Tradeoff: Simplicity vs. scalability

**B. Add streaming API**
```go
func StreamResources(path string) (<-chan Resource, <-chan error)
```
- Pros: Handles arbitrary sizes, memory-efficient
- Cons: More complex API, harder error handling
- Tradeoff: Scalability vs. API complexity

**C. Add both APIs**
- Pros: Simple for simple cases, powerful for complex
- Cons: Two APIs to maintain, documentation burden
- Tradeoff: Flexibility vs. maintenance

**Recommendation:** **Option A for MVP, plan for B**. The current design is fine for initial release. Add streaming in v2 when there's proven demand. Document the 10MB limit clearly.

---

### 6. Fragment Resolution Context

**Issue:** No way to pass external context into fragment resolution (only fragment inputs).

**Problem:**
- Can't inject runtime values (timestamps, user info, environment)
- Forces all data through fragment inputs
- Limits dynamic prompt generation

**Options:**

**A. Keep fragment-only inputs (current)**
- Pros: Pure, deterministic, testable
- Cons: Inflexible for runtime scenarios
- Tradeoff: Purity vs. practicality

**B. Add global context parameter**
```go
func ResolveBody(body Body, fragments map[string]Fragment, context map[string]interface{})
```
- Pros: Flexible, supports runtime injection
- Cons: Breaks determinism, harder to test, security concerns
- Tradeoff: Flexibility vs. predictability

**C. Add context as special fragment input**
- Pros: Explicit, traceable, type-safe
- Cons: Verbose, requires schema changes
- Tradeoff: Safety vs. ergonomics

**Recommendation:** **Option A for core, B for higher-level API**. Keep core resolution pure and deterministic. Let `ai-resource-manager` (the higher-level system) handle context injection by preprocessing inputs. This maintains clean separation of concerns.

---

## Summary of Recommendations

| Issue | Decision | Priority | Action Required |
|-------|----------|----------|-----------------|
| Fragment nesting | Do not implement | High | None - closed |
| Input validation timing | Validate during semantic phase | High | Clarify specs |
| Body type safety | Use explicit struct | Medium | Implement custom unmarshaling |
| Error aggregation | Fail-fast between phases | Low | Clarify specs |
| Streaming loading | Keep in-memory only | Low | None - closed |
| Resolution context | Keep pure, defer to future | Medium | None - closed |

## Next Steps

1. ✅ Clarify input validation timing in specs (validate during semantic phase)
2. ✅ Document error aggregation strategy (fail-fast between phases)
3. 🔨 Update specs to reflect validation pipeline: Schema → Semantic (w/ input validation) → Success
4. 🔨 Implement Body type with custom unmarshaling (see `docs/body-type-design.md`)
5. 📝 Update `semantic-validation.md` to include fragment input validation
6. 📝 Update `fragment-input-validation.md` to clarify it runs during semantic phase

## Decision Log

### 2026-02-23

**Issue #1: Fragment Nesting**
- **Decision:** Do not implement fragment nesting
- **Rationale:** Introduces complexity creep without sufficient benefit. Fragments provide composition within a resource; if more complex composition is needed, use multiple resources.
- **Status:** CLOSED - Will not implement

**Issue #2: Input Validation Timing**
- **Decision:** Validate fragment inputs during semantic validation phase (Option A)
- **Rationale:** Users should receive all errors upfront. Validation pipeline: Schema → Semantic (including input validation) → Success. Each phase returns immediately if errors found.
- **Status:** RESOLVED - Spec clarification needed

**Issue #3: Body Type Safety**
- **Decision:** Use explicit union type with custom unmarshaling (Option B)
- **Rationale:** Compile-time safety and self-documenting API worth the ~40 lines of custom unmarshaling code. See `docs/body-type-design.md` for implementation details.
- **Status:** RESOLVED - Implementation required

**Issue #4: Error Aggregation Strategy**
- **Decision:** Fail-fast between phases, collect within phases (Modified Option B)
- **Rationale:** Schema errors → return. Semantic errors (including input validation) → return. This prevents cascading errors while still showing all errors within a validation phase.
- **Status:** RESOLVED - Spec clarification needed

**Issue #5: Streaming/Incremental Loading**
- **Decision:** Keep in-memory loading only (Option A)
- **Rationale:** 10MB limit is sufficient for realistic use cases. Resource bundles are unlikely to exceed this. Streaming adds complexity without proven need.
- **Status:** CLOSED - Will not implement

**Issue #6: Fragment Resolution Context**
- **Decision:** Keep core resolution pure, no runtime context injection (Option A)
- **Rationale:** Core responsibility is deterministic interpretation, not dynamic generation. Runtime customization should happen via: (1) explicit fragment inputs, (2) resource overlays/patches (kustomize-style), or (3) higher-level systems like ai-resource-manager. Runtime injection may be reconsidered as future enhancement if clear use cases emerge.
- **Status:** CLOSED - Deferred to future consideration
