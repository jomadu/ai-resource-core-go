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