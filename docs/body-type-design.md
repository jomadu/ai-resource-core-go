# Body Type Design

## Current Design (interface{})

```go
type Body interface{}

// Usage - requires runtime type checking
func processBody(body Body) error {
    switch v := body.(type) {
    case string:
        return processString(v)
    case []interface{}:
        return processArray(v)
    default:
        return fmt.Errorf("body must be string or array")
    }
}
```

**Unmarshaling:**
```go
// YAML automatically unmarshals to interface{}
var spec PromptSpec
yaml.Unmarshal(data, &spec)
// spec.Body is now interface{} containing string or []interface{}
```

## Proposed Design (Explicit Union)

```go
type Body struct {
    // Exactly one of these will be non-nil
    String *string
    Array  []BodyItem
}

type BodyItem struct {
    // Exactly one of these will be non-nil
    String        *string
    FragmentRef   *FragmentRef
}

// Usage - type-safe at compile time
func processBody(body Body) error {
    if body.String != nil {
        return processString(*body.String)
    }
    if body.Array != nil {
        return processArray(body.Array)
    }
    return fmt.Errorf("body must have either String or Array set")
}
```

**Unmarshaling (custom):**
```go
func (b *Body) UnmarshalYAML(node *yaml.Node) error {
    // Try string first
    var s string
    if err := node.Decode(&s); err == nil {
        b.String = &s
        return nil
    }
    
    // Try array
    var arr []BodyItem
    if err := node.Decode(&arr); err == nil {
        b.Array = arr
        return nil
    }
    
    return fmt.Errorf("body must be string or array")
}

func (bi *BodyItem) UnmarshalYAML(node *yaml.Node) error {
    // Try string first
    var s string
    if err := node.Decode(&s); err == nil {
        bi.String = &s
        return nil
    }
    
    // Try fragment reference (has "fragment" key)
    var ref FragmentRef
    if err := node.Decode(&ref); err == nil {
        bi.FragmentRef = &ref
        return nil
    }
    
    return fmt.Errorf("body item must be string or fragment reference")
}
```

## Comparison

### interface{} Approach
```go
// Consumer code - error-prone
body := prompt.Spec.Body
if str, ok := body.(string); ok {
    fmt.Println(str)
} else if arr, ok := body.([]interface{}); ok {
    for _, item := range arr {
        // More type assertions needed...
    }
}
```

### Explicit Union Approach
```go
// Consumer code - type-safe
body := prompt.Spec.Body
if body.String != nil {
    fmt.Println(*body.String)
}
if body.Array != nil {
    for _, item := range body.Array {
        if item.String != nil {
            fmt.Println(*item.String)
        }
        if item.FragmentRef != nil {
            fmt.Println(item.FragmentRef.Fragment)
        }
    }
}
```

## Tradeoffs

| Aspect | interface{} | Explicit Union |
|--------|-------------|----------------|
| Unmarshaling complexity | Simple (automatic) | Custom (20 lines) |
| Consumer code safety | Runtime errors | Compile-time safety |
| API clarity | Unclear from types | Self-documenting |
| Maintenance | Easy to break | Hard to misuse |
| Performance | Reflection overhead | Direct field access |

## Recommendation

Use explicit union. The custom unmarshaling is ~40 lines total but provides:
- Compile-time safety for all consumers
- Self-documenting API
- Better IDE autocomplete
- Catches errors at parse time, not usage time

The unmarshaling complexity is hidden from users and only written once.
