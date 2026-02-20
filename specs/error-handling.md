# Error Handling

## Job to be Done
Provide clear, actionable error messages that help developers quickly identify and fix issues with their AI Resources.

## Activities
- Define structured error types for different failure modes
- Include context (file paths, field paths, fragment IDs) in errors
- Provide expected vs actual information for validation errors
- Support error aggregation (multiple errors at once)
- Format errors for human readability
- Enable programmatic error inspection

## Acceptance Criteria
- [ ] Each error type has clear purpose and fields
- [ ] Errors include enough context to locate the problem
- [ ] Error messages are human-readable
- [ ] Errors can be inspected programmatically
- [ ] Multiple errors can be collected and returned together
- [ ] Error messages suggest fixes where possible
- [ ] File and line information included when available
- [ ] Errors implement standard Go error interface

## Data Structures

### ValidationError
```go
type ValidationError struct {
    Field   string
    Message string
    Cause   error
}

func (e *ValidationError) Error() string
```

**Fields:**
- `Field` - Path to invalid field (e.g., "spec.body", "metadata.id")
- `Message` - Human-readable error description
- `Cause` - Optional underlying error

**Usage:** Schema and semantic validation failures

### SchemaError
```go
type SchemaError struct {
    Path    string
    Message string
}

func (e *SchemaError) Error() string
```

**Fields:**
- `Path` - JSON path to invalid field
- `Message` - Schema violation description

**Usage:** JSON Schema validation failures

### FragmentError
```go
type FragmentError struct {
    FragmentID string
    Message    string
    Cause      error
}

func (e *FragmentError) Error() string
```

**Fields:**
- `FragmentID` - ID of fragment with error
- `Message` - Error description
- `Cause` - Optional underlying error

**Usage:** Fragment resolution and rendering failures

### InputError
```go
type InputError struct {
    FragmentID string
    InputName  string
    Expected   string
    Got        string
}

func (e *InputError) Error() string
```

**Fields:**
- `FragmentID` - ID of fragment being validated
- `InputName` - Name of input with error
- `Expected` - Expected type or constraint
- `Got` - Actual value or type received

**Usage:** Fragment input validation failures

### LoadError
```go
type LoadError struct {
    Path    string
    Message string
    Cause   error
}

func (e *LoadError) Error() string
```

**Fields:**
- `Path` - File path that failed to load
- `Message` - Error description
- `Cause` - Underlying error (parse error, file not found, etc.)

**Usage:** File loading and parsing failures

### MultiError
```go
type MultiError struct {
    Errors []error
}

func (e *MultiError) Error() string
func (e *MultiError) Unwrap() []error
```

**Fields:**
- `Errors` - Collection of errors

**Usage:** Aggregating multiple validation errors

## Algorithm

### Error Formatting

```
function Error() string:
    return format_error_message(error_type, fields)
```

**Examples:**
- ValidationError: "validation error at spec.body: body must be string or array"
- SchemaError: "schema validation failed at metadata.id: does not match pattern ^[a-zA-Z0-9_-]+$"
- FragmentError: "fragment 'read-file' error: template rendering failed"
- InputError: "fragment 'read-file' input 'path': expected string, got number"
- LoadError: "failed to load resource.yaml: invalid YAML syntax at line 5"

### Error Aggregation

```
function CollectErrors(operations):
    errors = []
    
    for op in operations:
        result, err = op()
        if err:
            errors.append(err)
    
    if len(errors) == 0:
        return nil
    
    if len(errors) == 1:
        return errors[0]
    
    return MultiError{Errors: errors}
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Single error | Return error directly, not wrapped in MultiError |
| Multiple errors | Return MultiError with all errors |
| Nested errors | Preserve cause chain with Unwrap |
| Empty error message | Use default message for error type |
| Very long field path | Truncate with ellipsis if needed |
| Binary data in error | Show hex representation or length |

## Dependencies

- `core-types.md` - Types referenced in errors
- All validation specs - Use these error types

## Implementation Mapping

**Source files:**
- `pkg/airesource/errors.go` - All error type definitions

**Related specs:**
- `resource-loading.md` - Uses LoadError
- `schema-validation.md` - Uses SchemaError
- `semantic-validation.md` - Uses ValidationError, FragmentError
- `fragment-input-validation.md` - Uses InputError
- `fragment-resolution.md` - Uses FragmentError

## Examples

### Example 1: ValidationError

**Input:**
```go
err := &ValidationError{
    Field:   "spec.body",
    Message: "body must be string or array",
}
```

**Expected Output:**
```go
err.Error() == "validation error at spec.body: body must be string or array"
```

**Verification:**
- Error message includes field path
- Error message includes description

### Example 2: InputError

**Input:**
```go
err := &InputError{
    FragmentID: "read-file",
    InputName:  "path",
    Expected:   "string",
    Got:        "number",
}
```

**Expected Output:**
```go
err.Error() == "fragment 'read-file' input 'path': expected string, got number"
```

**Verification:**
- Error identifies fragment
- Error identifies input
- Error shows type mismatch

### Example 3: MultiError

**Input:**
```go
err := &MultiError{
    Errors: []error{
        &ValidationError{Field: "metadata.id", Message: "invalid pattern"},
        &ValidationError{Field: "spec.body", Message: "missing required field"},
    },
}
```

**Expected Output:**
```go
err.Error() contains "2 errors"
err.Error() contains "metadata.id"
err.Error() contains "spec.body"
```

**Verification:**
- Error indicates multiple failures
- All errors included in message

### Example 4: LoadError with Cause

**Input:**
```go
parseErr := errors.New("invalid YAML syntax at line 5")
err := &LoadError{
    Path:    "resource.yaml",
    Message: "failed to parse YAML",
    Cause:   parseErr,
}
```

**Expected Output:**
```go
err.Error() contains "resource.yaml"
err.Error() contains "failed to parse YAML"
errors.Unwrap(err) == parseErr
```

**Verification:**
- Error includes file path
- Error includes description
- Cause is accessible via Unwrap

### Example 5: FragmentError

**Input:**
```go
err := &FragmentError{
    FragmentID: "read-file",
    Message:    "fragment not found in spec.fragments",
}
```

**Expected Output:**
```go
err.Error() == "fragment 'read-file': fragment not found in spec.fragments"
```

**Verification:**
- Error identifies fragment
- Error describes problem

## Notes

- All error types implement the standard Go `error` interface
- Errors support unwrapping for cause chains (Go 1.13+)
- MultiError implements `Unwrap() []error` for Go 1.20+ multi-error support
- Error messages should be lowercase (Go convention) unless starting with proper noun
- Field paths use dot notation: "spec.body", "spec.fragments.read-file.inputs.path"
- Array indices use bracket notation: "spec.body[0]", "spec.scope[1].files"
- Error messages should not end with punctuation
- Errors should be actionable - tell the user what's wrong and ideally how to fix it

## Known Issues

None.

## Areas for Improvement

- Could add error codes for programmatic handling
- Could add suggestions for common mistakes
- Could add "did you mean?" suggestions for typos
- Could add links to documentation for error types
- Could add structured logging support
