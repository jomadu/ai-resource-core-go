# AI Resource Core Specifications Index

## Jobs to be Done (JTBDs)

1. **Provide type-safe representations of AI Resources** - Type-safe Go structures for Prompts, Promptsets, Rules, and Rulesets that developers can use in their applications
2. **Load AI Resources from files** - Parse YAML and JSON files (single and multi-document) into validated Go structures with safety limits
3. **Validate resource structure** - Ensure resources conform to JSON Schema requirements for types, patterns, and enums
4. **Validate resource semantics** - Enforce logical consistency rules beyond schema (fragment references, collections, patterns)
5. **Resolve fragment references** - Transform bodies with fragment references into rendered strings using Mustache templates
6. **Validate fragment inputs** - Type-check fragment inputs (string, number, boolean, array, object) with required validation and defaults
7. **Provide actionable errors** - Clear error messages with context, field paths, and expected vs actual information
8. **Verify spec compliance** - Test implementation against official AI Resource Specification test suite

## Topics of Concern

### Type System
- **Core Types** - Resource envelope, Kind enum, Metadata structure
- **Resource Types** - Prompt, Promptset, Rule, Ruleset with kind-specific specs
- **Fragment Types** - Fragment, InputDefinition, Body, FragmentRef
- **Type Accessors** - Safe conversion from generic Resource to specific types

### Loading & Parsing
- **File Loading** - YAML and JSON parsing from file extensions
- **Multi-Document Support** - Handle YAML documents separated by `---`
- **Safety Limits** - File size, array size, nesting depth, timeout constraints
- **Version Validation** - Check apiVersion is supported (currently draft only)

### Validation
- **Schema Validation** - JSON Schema validation with embedded schemas
- **Semantic Validation** - Fragment reference existence, collection constraints, pattern matching
- **Input Validation** - Fragment input type checking, required fields, default values
- **Error Aggregation** - Collect and return multiple validation errors

### Fragment System
- **Fragment Resolution** - Parse body structure, lookup fragments, render templates
- **Mustache Rendering** - Variable substitution, conditionals, array iteration
- **Input Processing** - Type validation, default application, recursive validation
- **Template Context** - Prepare inputs for Mustache rendering

### Error Handling
- **Error Types** - ValidationError, SchemaError, FragmentError, InputError, LoadError, MultiError
- **Error Context** - Field paths, fragment IDs, expected vs actual values
- **Error Formatting** - Human-readable messages with actionable information
- **Error Unwrapping** - Support Go 1.13+ error chains

### Testing
- **Conformance Testing** - Test against official spec test suite
- **Valid Cases** - Ensure valid resources pass all validation
- **Invalid Cases** - Ensure invalid resources fail with appropriate errors
- **Test Reporting** - Clear pass/fail results with context

## Conformance Philosophy

This implementation uses **submodule-based conformance testing** as the industry standard approach for ensuring interoperability across AI Resource implementations.

### Why Submodules

The [AI Resource Specification](https://github.com/jomadu/ai-resource-spec) repository is the **source of truth** for what constitutes a valid AI Resource. This implementation is an **interpreter** of that specification. By using the official test suite as a git submodule, we ensure:

- **No drift** - Tests come directly from the spec repository, not local copies
- **Interoperability** - All implementations test against the same fixtures
- **Traceability** - Each build pins to a specific spec version
- **Early detection** - Breaking changes in the spec are caught immediately

### Version Pinning Strategy

Git submodules pin to a specific commit by default. This provides:

- **Build stability** - Tests don't change unexpectedly
- **Controlled updates** - Spec updates are explicit, reviewed changes
- **Easier debugging** - Reproducible test failures across environments
- **CI/CD reliability** - Builds are deterministic

To update the spec version, use `make update-spec`, review changes, run tests, and commit the new pin if tests pass.

### Makefile Interface

The Makefile is the **documented interface** for all development tasks:

- `make test` - Runs all tests (auto-initializes submodule)
- `make test-conformance` - Runs conformance tests only
- `make build` - Builds all packages
- `make lint` - Runs linters
- `make update-spec` - Updates spec to latest version

This automation ensures developers don't need to manually manage submodules or remember git commands. The Makefile handles initialization, updates, and cleanup.

### Failure Policy

Conformance tests MUST fail if the submodule is missing or not initialized. There are no fallbacks to local fixtures. This hard requirement ensures:

- **Explicit dependencies** - Developers know the spec is required
- **No silent drift** - Missing fixtures cause immediate, visible failures
- **Clear error messages** - Failures direct users to `make test`

See [conformance-testing.md](conformance-testing.md) for detailed requirements and implementation guidance.

## Specification Documents

### Foundation
- [core-types.md](core-types.md) - Type system for AI Resources
- [error-handling.md](error-handling.md) - Structured error types and messages

### Loading & Validation
- [resource-loading.md](resource-loading.md) - File loading and parsing
- [schema-validation.md](schema-validation.md) - JSON Schema validation
- [semantic-validation.md](semantic-validation.md) - Logical consistency validation

### Fragment System
- [fragment-resolution.md](fragment-resolution.md) - Fragment reference resolution and Mustache rendering
- [fragment-input-validation.md](fragment-input-validation.md) - Fragment input type checking

### Testing
- [conformance-testing.md](conformance-testing.md) - Spec compliance verification
