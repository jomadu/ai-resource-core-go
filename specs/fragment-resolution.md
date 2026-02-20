# Fragment Resolution

## Job to be Done
Transform resource bodies containing fragment references into final rendered strings by resolving references and rendering Mustache templates.

## Activities
- Parse body structure (string or array)
- Identify fragment references in body
- Lookup fragments from spec.fragments
- Render Mustache templates with provided inputs
- Join resolved parts into final string
- Handle nested template structures

## Acceptance Criteria
- [ ] Simple string bodies are returned unchanged
- [ ] Array bodies are processed item by item
- [ ] Fragment references are replaced with rendered templates
- [ ] Mustache variables are substituted correctly
- [ ] Mustache conditionals are evaluated correctly
- [ ] Mustache array iteration works for primitives and objects
- [ ] Resolved parts are joined with double newline separator
- [ ] Fragment not found returns clear error
- [ ] Template rendering errors include fragment context
- [ ] Empty body array returns empty string

## Data Structures

### ResolveOptions
```go
type ResolveOptions struct {
    Fragments map[string]Fragment
}
```

**Fields:**
- `Fragments` - Map of fragment definitions available for resolution

## Algorithm

1. Check body type (string or array)
2. If string: return as-is
3. If array: process each item
   - If item is string: add to parts
   - If item is fragment reference:
     - Lookup fragment
     - Render template with inputs
     - Add rendered result to parts
4. Join parts with "\n\n"
5. Return final string

**Pseudocode:**
```
function ResolveBody(body, fragments):
    if is_string(body):
        return body
    
    if not is_array(body):
        return error("body must be string or array")
    
    parts = []
    
    for item in body:
        if is_string(item):
            parts.append(item)
        else if is_fragment_ref(item):
            fragment = fragments[item.Fragment]
            if not fragment:
                return error("fragment not found: {item.Fragment}")
            
            rendered = render_mustache(fragment.Body, item.Inputs)
            parts.append(rendered)
        else:
            return error("invalid body item type")
    
    return join(parts, "\n\n")
```

### Mustache Rendering

```
function render_mustache(template, inputs):
    context = prepare_context(inputs)
    result = mustache.Render(template, context)
    return result
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Empty string body | Return empty string |
| Empty array body | Return empty string |
| Fragment not found | Error: "fragment 'x' not found" |
| Template syntax error | Error with template context |
| Missing variable in template | Render as empty string (Mustache default) |
| Conditional with false value | Omit conditional section |
| Array iteration with empty array | Omit iteration section |
| Single item array body | Return item without extra newlines |
| Mixed string and fragment array | Join with "\n\n" separator |
| Fragment with no inputs | Render template with empty context |

## Dependencies

- `core-types.md` - Body and Fragment types
- `fragment-input-validation.md` - Input validation before resolution
- Mustache library (github.com/cbroglie/mustache)

## Implementation Mapping

**Source files:**
- `pkg/airesource/fragment.go` - ResolveBody function
- `internal/template/renderer.go` - Mustache rendering logic

**Related specs:**
- `core-types.md` - Fragment and Body types
- `fragment-input-validation.md` - Input validation
- `semantic-validation.md` - Fragment reference validation
- `error-handling.md` - Error types

## Examples

### Example 1: Simple String Body

**Input:**
```go
body := "Simple text body"
result, err := ResolveBody(body, fragments)
```

**Expected Output:**
```go
result == "Simple text body"
err == nil
```

**Verification:**
- String returned unchanged
- No errors

### Example 2: Fragment Reference

**Input:**
```go
fragments := map[string]Fragment{
    "greet": {
        Body: "Hello, {{name}}!",
    },
}

body := []interface{}{
    map[string]interface{}{
        "fragment": "greet",
        "inputs": map[string]interface{}{
            "name": "World",
        },
    },
}

result, err := ResolveBody(body, fragments)
```

**Expected Output:**
```go
result == "Hello, World!"
err == nil
```

**Verification:**
- Template rendered with input
- Variable substituted correctly

### Example 3: Mixed Array Body

**Input:**
```go
fragments := map[string]Fragment{
    "read": {
        Body: "Read file: {{path}}",
    },
}

body := []interface{}{
    "Introduction text",
    map[string]interface{}{
        "fragment": "read",
        "inputs": map[string]interface{}{
            "path": "data.txt",
        },
    },
    "Conclusion text",
}

result, err := ResolveBody(body, fragments)
```

**Expected Output:**
```go
result == "Introduction text\n\nRead file: data.txt\n\nConclusion text"
err == nil
```

**Verification:**
- Three parts joined with double newlines
- Fragment rendered in middle

### Example 4: Mustache Conditional

**Input:**
```go
fragments := map[string]Fragment{
    "conditional": {
        Body: "{{#show}}This is shown{{/show}}{{^show}}This is hidden{{/show}}",
    },
}

body := []interface{}{
    map[string]interface{}{
        "fragment": "conditional",
        "inputs": map[string]interface{}{
            "show": true,
        },
    },
}

result, err := ResolveBody(body, fragments)
```

**Expected Output:**
```go
result == "This is shown"
err == nil
```

**Verification:**
- Conditional evaluated correctly
- True branch rendered

### Example 5: Mustache Array Iteration

**Input:**
```go
fragments := map[string]Fragment{
    "list": {
        Body: "Files:\n{{#files}}- {{.}}\n{{/files}}",
    },
}

body := []interface{}{
    map[string]interface{}{
        "fragment": "list",
        "inputs": map[string]interface{}{
            "files": []string{"a.txt", "b.txt", "c.txt"},
        },
    },
}

result, err := ResolveBody(body, fragments)
```

**Expected Output:**
```go
result == "Files:\n- a.txt\n- b.txt\n- c.txt\n"
err == nil
```

**Verification:**
- Array iterated correctly
- Each item rendered

### Example 6: Fragment Not Found

**Input:**
```go
body := []interface{}{
    map[string]interface{}{
        "fragment": "missing",
        "inputs": map[string]interface{}{},
    },
}

result, err := ResolveBody(body, fragments)
```

**Expected Output:**
```go
err.Message == "fragment 'missing' not found"
```

**Verification:**
- Error indicates missing fragment
- Error includes fragment ID

## Notes

- Fragment resolution happens after validation passes
- Fragments are resource-scoped (cannot reference fragments from other resources)
- Fragments cannot reference other fragments (no nesting)
- Mustache is chosen for its simplicity and logic-less nature
- The double newline separator ("\n\n") creates clear visual separation between parts
- Missing Mustache variables render as empty strings (Mustache default behavior)
- Template syntax errors should include the fragment ID for debugging
- Resolution is deterministic - same inputs always produce same output

## Known Issues

None.

## Areas for Improvement

- Could add caching for rendered templates
- Could support custom Mustache delimiters
- Could add template linting/validation
- Could provide better error messages for common Mustache mistakes
- Could support partial templates in future versions
