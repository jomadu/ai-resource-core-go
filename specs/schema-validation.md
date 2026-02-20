# Schema Validation

## Job to be Done
Ensure loaded AI Resources conform to the structural requirements defined in the JSON Schema specification.

## Activities
- Embed JSON schemas for each resource kind
- Select appropriate schema based on resource kind
- Validate resource structure against schema
- Return structured validation errors with field paths
- Support draft API version schemas

## Acceptance Criteria
- [ ] All valid resources pass schema validation
- [ ] Invalid resources fail with specific field errors
- [ ] Error messages include JSON path to invalid field
- [ ] Missing required fields are detected
- [ ] Type mismatches are detected (string vs number, etc.)
- [ ] Pattern constraints are validated (e.g., metadata.id)
- [ ] Enum constraints are validated (e.g., enforcement levels)
- [ ] Array and object structures are validated
- [ ] Schemas are embedded in binary (no external file dependencies)

## Data Structures

### SchemaError
```go
type SchemaError struct {
    Path    string
    Message string
}
```

**Fields:**
- `Path` - JSON path to invalid field (e.g., "spec.body", "metadata.id")
- `Message` - Human-readable error description

## Algorithm

1. Determine resource kind from parsed resource
2. Select corresponding JSON schema
3. Convert resource to JSON for validation
4. Run JSON schema validator
5. Collect validation errors
6. Convert errors to structured format with paths
7. Return validation result

**Pseudocode:**
```
function ValidateSchema(resource):
    schema = get_schema_for_kind(resource.Kind)
    
    json_data = convert_to_json(resource)
    
    result = schema.Validate(json_data)
    
    if not result.Valid():
        errors = []
        for err in result.Errors():
            errors.append(SchemaError{
                Path: err.Field(),
                Message: err.Description()
            })
        return errors
    
    return nil
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Unknown kind | Return error: "no schema for kind: {kind}" |
| Null required field | Return error with field path |
| Wrong type for field | Return error: "expected {type}, got {actual}" |
| Invalid enum value | Return error listing valid values |
| Pattern mismatch | Return error with pattern requirement |
| Multiple validation errors | Return all errors, not just first |
| Nested object errors | Return full path (e.g., "spec.fragments.read.inputs.path") |
| Array item errors | Return path with index (e.g., "spec.scope[0].files") |

## Dependencies

- `core-types.md` - Resource types being validated
- `resource-loading.md` - Resources to validate
- JSON Schema library (github.com/xeipuuv/gojsonschema)
- Embedded JSON schemas from ai-resource-spec repo

## Implementation Mapping

**Source files:**
- `internal/schema/validator.go` - Schema validation logic
- `internal/schema/schemas.go` - Embedded JSON schemas
- `pkg/airesource/errors.go` - SchemaError type

**Related specs:**
- `core-types.md` - Types being validated
- `semantic-validation.md` - Additional validation beyond schema
- `error-handling.md` - Error type definitions

## Examples

### Example 1: Valid Resource

**Input:**
```go
resource := &Resource{
    APIVersion: "ai-resource/draft",
    Kind: KindPrompt,
    Metadata: Metadata{
        ID: "test-prompt",
    },
    Spec: PromptSpec{
        Body: "Test body",
    },
}

err := ValidateSchema(resource)
```

**Expected Output:**
```go
err == nil
```

**Verification:**
- No validation errors returned

### Example 2: Missing Required Field

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  name: "Missing ID"
spec:
  body: "Test"
```

**Expected Output:**
```go
err.Path == "metadata.id"
err.Message contains "required"
```

**Verification:**
- Error indicates missing required field
- Path points to metadata.id

### Example 3: Invalid Pattern

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: "invalid id with spaces"
spec:
  body: "Test"
```

**Expected Output:**
```go
err.Path == "metadata.id"
err.Message contains "pattern" or "^[a-zA-Z0-9_-]+$"
```

**Verification:**
- Error indicates pattern violation
- Message shows expected pattern

### Example 4: Invalid Enum Value

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Rule
metadata:
  id: "test-rule"
spec:
  enforcement: "always"
  body: "Test"
```

**Expected Output:**
```go
err.Path == "spec.enforcement"
err.Message contains "may, should, must"
```

**Verification:**
- Error indicates invalid enum value
- Message lists valid values

### Example 5: Multiple Errors

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Rule
metadata:
  id: "invalid id!"
spec:
  enforcement: "invalid"
```

**Expected Output:**
```go
len(errors) >= 2
// One error for metadata.id pattern
// One error for spec.enforcement enum
// One error for missing spec.body
```

**Verification:**
- All validation errors are returned
- Each error has distinct path

## Notes

- Schema validation is structural only - it checks types, required fields, patterns, and enums
- Semantic validation (e.g., fragment references exist) is handled separately
- Schemas are embedded using `go:embed` to avoid runtime file dependencies
- Schema files should be copied from the ai-resource-spec repository
- Validation happens after loading but before semantic validation
- The validator should be version-aware (different schemas for different API versions)

## Known Issues

None.

## Areas for Improvement

- Could cache compiled schemas for performance
- Could provide schema validation for fragments separately
- Could add custom error messages for common mistakes
- Could support schema validation without full resource loading
