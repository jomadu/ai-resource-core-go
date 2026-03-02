# AI Resource Core (Go)

`ai-resource-core-go` is the reference Go implementation of the AI Resource Specification.

It is responsible for deterministic interpretation of AI Resources. Core guarantees that any loaded resource:

- Conforms structurally to the schema
- Passes semantic validation rules
- Is normalized into a consistent internal representation
- Matches a supported `apiVersion`

## Responsibilities

### Core provides:

- YAML parsing
- JSON Schema validation
- Envelope enforcement
- Semantic validation
- Multi-document bundle loading
- Structured validation errors

### Core does not:

- Execute resources
- Manage workspaces
- Persist state
- Integrate with runtimes
- Apply business logic

Those responsibilities belong to higher-level systems such as `ai-resource-manager`.

## Example

```go
res, err := core.LoadFile("resource.yaml")
if err != nil {
    log.Fatal(err)
}
```

If `LoadFile` succeeds, the resource is spec-compliant and normalized.

## Design Principles

- Strict adherence to the spec
- Deterministic behavior
- Minimal public API
- Clear separation from lifecycle and execution concerns

Core enforces the contract. Other systems build on top of it.

## Quick Start

```bash
make test   # Run all tests
make build  # Build packages
make lint   # Run linters
make help   # Show all commands
```

Conformance tests automatically use the official [AI Resource Specification](https://github.com/jomadu/ai-resource-spec) test suite embedded via git submodule. The Makefile handles initialization.

## Testing

### Conformance Testing

The project includes a conformance test suite to verify correct interpretation of AI Resources according to the specification.

**Running Tests:**
```bash
make test              # Run all tests (auto-initializes submodule)
make test-conformance  # Run conformance tests only
make update-spec       # Update spec to latest version
```

The conformance tests use the official [AI Resource Specification](https://github.com/jomadu/ai-resource-spec) test suite embedded via git submodule at `internal/assets/spec/`. Schemas and test fixtures are embedded in the binary using `go:embed`. The Makefile automatically initializes the submodule when running tests.

**Manual Workflow (Alternative):**
```bash
git submodule update --init --recursive
go test ./...
```

**Test Coverage:**
- Valid resource loading (Prompt, Promptset, Rule, Ruleset)
- Invalid resource rejection with appropriate errors
- Fragment resolution with Mustache templates
- Multi-document YAML loading
- Schema and semantic validation

See `AUDIT.md` for conformance testing status and recommendations.