# Core Types

## Job to be Done
Provide type-safe representations of AI Resources that developers can use to work with Prompts, Promptsets, Rules, and Rulesets in their Go applications.

## Activities
- Define resource envelope structure (apiVersion, kind, metadata, spec)
- Define resource kinds (Prompt, Promptset, Rule, Ruleset)
- Define metadata structure with ID constraints
- Define spec structures for each resource kind
- Define fragment and body types
- Provide type-safe accessors for resource specs

## Acceptance Criteria
- [ ] All resource types match the AI Resource Specification structure
- [ ] Type system enforces required fields at compile time where possible
- [ ] Metadata ID pattern constraint is documented
- [ ] Fragment input types support string, number, boolean, array, object
- [ ] Body type supports both string and array formats
- [ ] Type-safe accessors return appropriate errors for kind mismatches
- [ ] Enforcement levels (may, should, must) are type-safe constants
- [ ] Rule priority defaults to 100 when not specified

## Data Structures

### Resource Envelope
```go
type Resource struct {
    APIVersion string
    Kind       Kind
    Metadata   Metadata
    Spec       interface{}
}
```

**Fields:**
- `APIVersion` - Version identifier (e.g., "ai-resource/draft")
- `Kind` - Resource type (Prompt, Promptset, Rule, Ruleset)
- `Metadata` - Resource metadata including required ID
- `Spec` - Kind-specific specification (use type accessors)

### Kind
```go
type Kind string

const (
    KindPrompt     Kind = "Prompt"
    KindPromptset  Kind = "Promptset"
    KindRule       Kind = "Rule"
    KindRuleset    Kind = "Ruleset"
)
```

### Metadata
```go
type Metadata struct {
    ID          string
    Name        string
    Description string
}
```

**Fields:**
- `ID` - Required unique identifier, must match pattern `^[a-zA-Z0-9_-]+$`
- `Name` - Optional human-readable name
- `Description` - Optional description

### Prompt
```go
type Prompt struct {
    APIVersion string
    Kind       Kind
    Metadata   Metadata
    Spec       PromptSpec
}

type PromptSpec struct {
    Fragments map[string]Fragment
    Body      Body
}
```

### Promptset
```go
type Promptset struct {
    APIVersion string
    Kind       Kind
    Metadata   Metadata
    Spec       PromptsetSpec
}

type PromptsetSpec struct {
    Fragments map[string]Fragment
    Prompts   map[string]PromptItem
}

type PromptItem struct {
    Name        string
    Description string
    Body        Body
}
```

### Rule
```go
type Rule struct {
    APIVersion string
    Kind       Kind
    Metadata   Metadata
    Spec       RuleSpec
}

type RuleSpec struct {
    Fragments   map[string]Fragment
    Enforcement Enforcement
    Scope       []ScopeEntry
    Body        Body
}

type Enforcement string

const (
    EnforcementMay    Enforcement = "may"
    EnforcementShould Enforcement = "should"
    EnforcementMust   Enforcement = "must"
)

type ScopeEntry struct {
    Files []string
}
```

### Ruleset
```go
type Ruleset struct {
    APIVersion string
    Kind       Kind
    Metadata   Metadata
    Spec       RulesetSpec
}

type RulesetSpec struct {
    Fragments map[string]Fragment
    Rules     map[string]RuleItem
}

type RuleItem struct {
    Name        string
    Description string
    Priority    int
    Enforcement Enforcement
    Scope       []ScopeEntry
    Body        Body
}
```

### Fragment
```go
type Fragment struct {
    Inputs map[string]InputDefinition
    Body   string
}

type InputDefinition struct {
    Type       InputType
    Required   bool
    Default    interface{}
    Items      *InputDefinition
    Properties map[string]InputDefinition
}

type InputType string

const (
    InputTypeString  InputType = "string"
    InputTypeNumber  InputType = "number"
    InputTypeBoolean InputType = "boolean"
    InputTypeArray   InputType = "array"
    InputTypeObject  InputType = "object"
)
```

### Body
```go
type Body interface{}

type BodyItem interface{}

type FragmentRef struct {
    Fragment string
    Inputs   map[string]interface{}
}
```

**Body formats:**
- Simple string: `"text content"`
- Array of strings: `["text1", "text2"]`
- Array with fragment refs: `[{fragment: "id", inputs: {...}}, "text"]`

## Algorithm

Type-safe accessors:

```go
func (r *Resource) AsPrompt() (*Prompt, error)
func (r *Resource) AsPromptset() (*Promptset, error)
func (r *Resource) AsRule() (*Rule, error)
func (r *Resource) AsRuleset() (*Ruleset, error)
```

**Pseudocode:**
```
function AsPrompt(resource):
    if resource.Kind != KindPrompt:
        return error("expected Prompt, got {Kind}")
    
    prompt = convert resource to Prompt type
    return prompt
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Wrong kind accessor | Return error with expected vs actual kind |
| Nil spec | Return error indicating invalid resource |
| Empty metadata.ID | Validation catches this (not type system) |
| Invalid enforcement value | Validation catches this (not type system) |

## Dependencies

- None (foundational types)

## Implementation Mapping

**Source files:**
- `pkg/airesource/resource.go` - Resource envelope and Kind
- `pkg/airesource/prompt.go` - Prompt and Promptset types
- `pkg/airesource/rule.go` - Rule and Ruleset types
- `pkg/airesource/fragment.go` - Fragment and Body types

**Related specs:**
- `resource-loading.md` - Uses these types for parsing
- `schema-validation.md` - Validates these types
- `semantic-validation.md` - Enforces constraints on these types
- `fragment-resolution.md` - Resolves Body and Fragment types

## Examples

### Example 1: Type-safe Resource Access

**Input:**
```go
resource, _ := LoadResource("prompt.yaml")
prompt, err := resource.AsPrompt()
```

**Expected Output:**
```go
// If resource.Kind == KindPrompt:
prompt != nil, err == nil

// If resource.Kind != KindPrompt:
prompt == nil, err == "expected kind Prompt, got Rule"
```

**Verification:**
- Check error is nil for matching kind
- Check error message for mismatched kind

### Example 2: Fragment Input Definition

**Input:**
```go
fragment := Fragment{
    Inputs: map[string]InputDefinition{
        "path": {
            Type:     InputTypeString,
            Required: true,
        },
        "count": {
            Type:     InputTypeNumber,
            Required: false,
            Default:  10,
        },
    },
    Body: "Read {{path}} with limit {{count}}",
}
```

**Expected Output:**
- Fragment structure is valid
- Input types are constrained to valid values
- Default value is stored as interface{}

**Verification:**
- Compile-time type safety for InputType constants
- Runtime validation of default value types (separate spec)

## Notes

- The type system prioritizes clarity over clever abstractions
- `interface{}` is used for `Spec`, `Body`, and `Default` because these have dynamic structure determined by kind/type
- Type accessors provide safe conversion from generic Resource to specific types
- Pattern validation for metadata.ID is enforced during validation, not in the type system
- Rule priority default (100) is applied during parsing, not in the type definition

## Known Issues

None.

## Areas for Improvement

- Consider using generics for type-safe resource loading when Go 1.18+ is baseline
- Could use sum types if Go adds discriminated unions in future
