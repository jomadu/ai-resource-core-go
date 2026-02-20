# Fragment Input Validation

## Job to be Done
Validate that inputs provided to fragment references match the type definitions and constraints specified in the fragment's input schema.

## Activities
- Validate input types match definitions (string, number, boolean, array, object)
- Apply default values for missing optional inputs
- Validate required inputs are provided
- Recursively validate array item types
- Recursively validate object property types
- Return clear type mismatch errors

## Acceptance Criteria
- [ ] String inputs are validated as strings
- [ ] Number inputs are validated as numbers
- [ ] Boolean inputs are validated as booleans
- [ ] Array inputs are validated with item type checking
- [ ] Object inputs are validated with property type checking
- [ ] Missing optional inputs receive default values
- [ ] Missing required inputs return error
- [ ] Type mismatches return error with expected vs actual
- [ ] Nested structures are validated recursively
- [ ] Extra inputs not in definition return error

## Data Structures

### InputError
```go
type InputError struct {
    FragmentID string
    InputName  string
    Expected   string
    Got        string
}
```

**Fields:**
- `FragmentID` - ID of fragment being validated
- `InputName` - Name of input with error
- `Expected` - Expected type or constraint
- `Got` - Actual value or type received

## Algorithm

1. Get fragment definition
2. For each defined input:
   - Check if provided in inputs
   - If missing and required: error
   - If missing and has default: apply default
   - If provided: validate type
3. For each provided input:
   - Check if defined in fragment
   - If not defined: error
4. Return validated inputs with defaults applied

**Pseudocode:**
```
function ValidateInputs(fragmentID, fragment, providedInputs):
    validated = {}
    errors = []
    
    // Check defined inputs
    for name, definition in fragment.Inputs:
        value = providedInputs[name]
        
        if value == nil:
            if definition.Required:
                errors.append("required input '{name}' not provided")
                continue
            if definition.Default != nil:
                validated[name] = definition.Default
                continue
        else:
            if not validate_type(value, definition):
                errors.append("input '{name}' type mismatch")
            else:
                validated[name] = value
    
    // Check for undefined inputs
    for name in providedInputs:
        if not fragment.Inputs[name]:
            errors.append("undefined input '{name}'")
    
    if errors:
        return nil, errors
    
    return validated, nil
```

### Type Validation

```
function validate_type(value, definition):
    switch definition.Type:
        case "string":
            return is_string(value)
        case "number":
            return is_number(value)
        case "boolean":
            return is_boolean(value)
        case "array":
            if not is_array(value):
                return false
            for item in value:
                if not validate_type(item, definition.Items):
                    return false
            return true
        case "object":
            if not is_object(value):
                return false
            for prop_name, prop_def in definition.Properties:
                if not validate_type(value[prop_name], prop_def):
                    return false
            return true
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Missing required input | Error: "required input 'x' not provided" |
| Missing optional input with default | Apply default value |
| Missing optional input without default | Omit from validated inputs |
| Type mismatch | Error: "expected string, got number" |
| Array with wrong item type | Error: "array item expected string, got number" |
| Object with wrong property type | Error: "property 'x' expected number, got string" |
| Undefined input provided | Error: "undefined input 'x'" |
| Empty array | Valid if array type |
| Empty object | Valid if object type |
| Null value | Treated as missing |
| Nested array validation | Recursively validate items |
| Nested object validation | Recursively validate properties |

## Dependencies

- `core-types.md` - Fragment and InputDefinition types
- `semantic-validation.md` - Calls this for fragment reference validation

## Implementation Mapping

**Source files:**
- `pkg/airesource/fragment.go` - ValidateInputs function
- `internal/types/validator.go` - Type validation helpers

**Related specs:**
- `core-types.md` - InputDefinition types
- `fragment-resolution.md` - Uses validated inputs
- `semantic-validation.md` - Validates fragment references exist
- `error-handling.md` - InputError type

## Examples

### Example 1: Valid String Input

**Input:**
```go
fragment := Fragment{
    Inputs: map[string]InputDefinition{
        "path": {Type: InputTypeString, Required: true},
    },
}

inputs := map[string]interface{}{
    "path": "test.txt",
}

validated, err := ValidateInputs("read-file", fragment, inputs)
```

**Expected Output:**
```go
validated["path"] == "test.txt"
err == nil
```

**Verification:**
- Input validated successfully
- Value preserved

### Example 2: Apply Default Value

**Input:**
```go
fragment := Fragment{
    Inputs: map[string]InputDefinition{
        "count": {
            Type:     InputTypeNumber,
            Required: false,
            Default:  10,
        },
    },
}

inputs := map[string]interface{}{}

validated, err := ValidateInputs("list", fragment, inputs)
```

**Expected Output:**
```go
validated["count"] == 10
err == nil
```

**Verification:**
- Default value applied
- No error

### Example 3: Type Mismatch

**Input:**
```go
fragment := Fragment{
    Inputs: map[string]InputDefinition{
        "count": {Type: InputTypeNumber, Required: true},
    },
}

inputs := map[string]interface{}{
    "count": "not a number",
}

validated, err := ValidateInputs("list", fragment, inputs)
```

**Expected Output:**
```go
err.Expected == "number"
err.Got == "string"
err.InputName == "count"
```

**Verification:**
- Error indicates type mismatch
- Error includes expected and actual types

### Example 4: Array Validation

**Input:**
```go
fragment := Fragment{
    Inputs: map[string]InputDefinition{
        "files": {
            Type: InputTypeArray,
            Items: &InputDefinition{
                Type: InputTypeString,
            },
        },
    },
}

inputs := map[string]interface{}{
    "files": []string{"a.txt", "b.txt"},
}

validated, err := ValidateInputs("process", fragment, inputs)
```

**Expected Output:**
```go
validated["files"] == []string{"a.txt", "b.txt"}
err == nil
```

**Verification:**
- Array validated successfully
- Item types checked

### Example 5: Object Validation

**Input:**
```go
fragment := Fragment{
    Inputs: map[string]InputDefinition{
        "config": {
            Type: InputTypeObject,
            Properties: map[string]InputDefinition{
                "host": {Type: InputTypeString},
                "port": {Type: InputTypeNumber},
            },
        },
    },
}

inputs := map[string]interface{}{
    "config": map[string]interface{}{
        "host": "localhost",
        "port": 8080,
    },
}

validated, err := ValidateInputs("connect", fragment, inputs)
```

**Expected Output:**
```go
validated["config"].(map[string]interface{})["host"] == "localhost"
validated["config"].(map[string]interface{})["port"] == 8080
err == nil
```

**Verification:**
- Object validated successfully
- Property types checked

### Example 6: Undefined Input

**Input:**
```go
fragment := Fragment{
    Inputs: map[string]InputDefinition{
        "path": {Type: InputTypeString},
    },
}

inputs := map[string]interface{}{
    "path": "test.txt",
    "extra": "not defined",
}

validated, err := ValidateInputs("read", fragment, inputs)
```

**Expected Output:**
```go
err.Message contains "undefined input 'extra'"
```

**Verification:**
- Error indicates undefined input
- Error includes input name

## Notes

- Input validation happens before fragment resolution
- Type validation is strict - no automatic coercion
- Default values are applied only for missing optional inputs
- Nested structures (arrays, objects) are validated recursively
- Validation errors should be collected and returned together
- The validated inputs map includes both provided and default values
- Extra inputs not in the definition are rejected to catch typos

## Known Issues

None.

## Areas for Improvement

- Could add type coercion (e.g., "123" to 123) with opt-in flag
- Could add validation for additional constraints (min, max, pattern)
- Could provide better error messages for nested validation failures
- Could add support for union types in future versions
