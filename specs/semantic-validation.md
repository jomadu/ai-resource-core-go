# Semantic Validation

## Job to be Done
Enforce semantic rules and constraints that go beyond structural JSON Schema validation to ensure resources are logically valid and internally consistent.

## Activities
- Validate fragment references exist in spec.fragments
- Validate collection keys match ID pattern
- Validate collections have minimum 1 item
- Validate fragment input names match definitions
- Validate no circular dependencies (future-proofing)
- Validate scope glob patterns are well-formed
- Validate body format (string or array)

## Acceptance Criteria
- [ ] Fragment references in body must exist in spec.fragments
- [ ] Collection keys (prompts, rules, fragments) match pattern ^[a-zA-Z0-9_-]+$
- [ ] Promptset has at least one prompt
- [ ] Ruleset has at least one rule
- [ ] Fragment references include all required inputs
- [ ] Fragment references don't include undefined inputs
- [ ] Body is either string or array (not other types)
- [ ] Array body items are strings or fragment references
- [ ] Scope file patterns are valid glob syntax
- [ ] Clear error messages indicate which rule was violated

## Data Structures

### ValidationError
```go
type ValidationError struct {
    Field   string
    Message string
    Cause   error
}
```

**Fields:**
- `Field` - Path to invalid field (e.g., "spec.body[0].fragment")
- `Message` - Human-readable error description
- `Cause` - Optional underlying error

### FragmentError
```go
type FragmentError struct {
    FragmentID string
    Message    string
    Cause      error
}
```

**Fields:**
- `FragmentID` - ID of fragment with error
- `Message` - Error description
- `Cause` - Optional underlying error

## Algorithm

1. Check resource kind
2. Validate metadata constraints
3. Validate spec.fragments map keys match pattern
4. For each resource kind:
   - Validate kind-specific requirements
   - Validate body structure
   - Validate fragment references
   - Validate collections
5. Return all validation errors

**Pseudocode:**
```
function ValidateSemantic(resource):
    errors = []
    
    // Validate metadata
    if not matches_pattern(resource.Metadata.ID, "^[a-zA-Z0-9_-]+$"):
        errors.append("metadata.id does not match pattern")
    
    // Validate fragment keys
    for key in resource.Spec.Fragments:
        if not matches_pattern(key, "^[a-zA-Z0-9_-]+$"):
            errors.append("fragment key '{key}' invalid")
    
    // Kind-specific validation
    switch resource.Kind:
        case Prompt:
            validate_body(resource.Spec.Body, resource.Spec.Fragments, errors)
        case Promptset:
            validate_promptset(resource.Spec, errors)
        case Rule:
            validate_rule(resource.Spec, errors)
        case Ruleset:
            validate_ruleset(resource.Spec, errors)
    
    return errors
```

### Body Validation

```
function validate_body(body, fragments, errors):
    if is_string(body):
        return // Simple string body is valid
    
    if not is_array(body):
        errors.append("body must be string or array")
        return
    
    for item in body:
        if is_string(item):
            continue // String items are valid
        
        if is_fragment_ref(item):
            validate_fragment_ref(item, fragments, errors)
        else:
            errors.append("body item must be string or fragment reference")
```

### Fragment Reference Validation

```
function validate_fragment_ref(ref, fragments, errors):
    fragment = fragments[ref.Fragment]
    
    if not fragment:
        errors.append("fragment '{ref.Fragment}' not found")
        return
    
    // Check required inputs provided
    for input_name, input_def in fragment.Inputs:
        if input_def.Required and not ref.Inputs[input_name]:
            errors.append("required input '{input_name}' not provided")
    
    // Check no undefined inputs
    for input_name in ref.Inputs:
        if not fragment.Inputs[input_name]:
            errors.append("undefined input '{input_name}' for fragment '{ref.Fragment}'")
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| metadata.id with spaces | Error: "does not match pattern ^[a-zA-Z0-9_-]+$" |
| metadata.id with special chars | Error: "does not match pattern ^[a-zA-Z0-9_-]+$" |
| Fragment reference not found | Error: "fragment 'x' not found in spec.fragments" |
| Missing required fragment input | Error: "required input 'x' not provided for fragment 'y'" |
| Extra undefined fragment input | Error: "undefined input 'x' for fragment 'y'" |
| Empty promptset.prompts | Error: "promptset must have at least one prompt" |
| Empty ruleset.rules | Error: "ruleset must have at least one rule" |
| Body is number | Error: "body must be string or array" |
| Body array with object | Error: "body item must be string or fragment reference" |
| Invalid glob pattern | Error: "invalid glob pattern in scope.files" |
| Collection key with spaces | Error: "collection key 'x y' does not match pattern" |

## Dependencies

- `core-types.md` - Resource types being validated
- `schema-validation.md` - Must pass schema validation first
- Glob pattern validation library

## Implementation Mapping

**Source files:**
- `pkg/airesource/validator.go` - Main semantic validation
- `internal/types/validator.go` - Type-specific validation helpers
- `pkg/airesource/errors.go` - ValidationError and FragmentError types

**Related specs:**
- `core-types.md` - Types being validated
- `schema-validation.md` - Structural validation (runs first)
- `fragment-resolution.md` - Uses validated fragment references
- `error-handling.md` - Error type definitions

## Examples

### Example 1: Valid Resource

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: test-prompt
spec:
  fragments:
    read-file:
      inputs:
        path:
          type: string
          required: true
      body: "Read file: {{path}}"
  body:
    - fragment: read-file
      inputs:
        path: "test.txt"
```

**Expected Output:**
```go
err == nil
```

**Verification:**
- No validation errors

### Example 2: Fragment Not Found

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: test
spec:
  body:
    - fragment: missing-fragment
      inputs: {}
```

**Expected Output:**
```go
err.Message == "fragment 'missing-fragment' not found in spec.fragments"
```

**Verification:**
- Error indicates missing fragment
- Error includes fragment ID

### Example 3: Missing Required Input

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: test
spec:
  fragments:
    read-file:
      inputs:
        path:
          type: string
          required: true
      body: "Read {{path}}"
  body:
    - fragment: read-file
      inputs: {}
```

**Expected Output:**
```go
err.Message contains "required input 'path' not provided"
```

**Verification:**
- Error indicates missing required input
- Error includes input name and fragment ID

### Example 4: Invalid Collection Key

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Promptset
metadata:
  id: test
spec:
  prompts:
    "invalid key!":
      body: "Test"
```

**Expected Output:**
```go
err.Message contains "does not match pattern"
```

**Verification:**
- Error indicates invalid key pattern
- Error includes the invalid key

### Example 5: Empty Collection

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Promptset
metadata:
  id: test
spec:
  prompts: {}
```

**Expected Output:**
```go
err.Message == "promptset must have at least one prompt"
```

**Verification:**
- Error indicates empty collection
- Error is specific to resource kind

## Notes

- Semantic validation runs after schema validation passes
- Fragment references are validated but not resolved (resolution is separate)
- Pattern validation for IDs uses regex: `^[a-zA-Z0-9_-]+$`
- Collections (prompts, rules, fragments) must have at least one entry
- Fragment input validation checks existence and required status, not types (types validated during resolution)
- Glob pattern validation ensures patterns are syntactically valid, not that they match files
- All errors should be collected and returned together, not fail-fast

## Known Issues

None.

## Areas for Improvement

- Could add warnings for unused fragments
- Could validate fragment input default values match declared types
- Could detect potential circular dependencies in future versions
- Could provide suggestions for common typos in fragment references
