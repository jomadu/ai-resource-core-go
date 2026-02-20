# AI Resource Core (Go) - Implementation Task

Reference implementation of the [AI Resource Specification](https://github.com/maxdunn/ai-resource-spec).

## Specification Summary

The AI Resource Specification defines a versioned, declarative format for portable AI artifacts:

- **4 resource kinds**: Prompt, Promptset, Rule, Ruleset
- **Kubernetes-style envelope**: apiVersion, kind, metadata, spec
- **Inline fragment composition**: Reusable Mustache templates scoped to resources
- **Rule enforcement levels**: may, should, must
- **Rule scoping**: Optional glob-based file targeting
- **JSON Schema validation**: Structural + semantic validation

Current version: `ai-resource/draft`

## Project Structure

```
ai-resource-core-go/
├── go.mod
├── go.sum
├── README.md
├── LICENSE
├── TASK.md                  # This file
│
├── pkg/
│   └── airesource/          # Public API
│       ├── resource.go      # Core types (Resource, Metadata, envelope)
│       ├── version.go       # Version constants and validation
│       ├── prompt.go        # Prompt and Promptset types
│       ├── rule.go          # Rule and Ruleset types
│       ├── fragment.go      # Fragment types and resolution
│       ├── loader.go        # File loading (YAML/JSON)
│       ├── validator.go     # Validation logic
│       ├── errors.go        # Error types
│       └── options.go       # Configuration options
│
├── internal/
│   ├── schema/              # Schema validation (private)
│   │   ├── validator.go
│   │   └── schemas.go       # Embedded JSON schemas
│   ├── template/            # Mustache rendering (private)
│   │   └── renderer.go
│   └── types/               # Internal type validation
│       └── validator.go
│
├── testdata/                # Test fixtures
│   ├── valid/
│   │   ├── prompt.yml
│   │   ├── promptset.yml
│   │   ├── rule.yml
│   │   └── ruleset.yml
│   └── invalid/
│       ├── missing-body.yml
│       └── invalid-enforcement.yml
│
├── examples/                # Usage examples
│   ├── basic/
│   │   └── main.go
│   └── advanced/
│       └── main.go
│
└── cmd/                     # Optional CLI tools
    └── airesource/
        └── main.go          # Validation CLI
```

## Version Support

**Current implementation: Single-version (draft only)**

The spec currently has only one version: `ai-resource/draft`. The implementation should:

1. Support only `ai-resource/draft` initially
2. Return clear errors for unsupported versions
3. Be structured to easily add version support when v1 is released

```go
// pkg/airesource/version.go
const (
    APIVersionDraft = "ai-resource/draft"
)

var SupportedVersions = []string{
    APIVersionDraft,
}

func IsSupportedVersion(version string) bool {
    for _, v := range SupportedVersions {
        if v == version {
            return true
        }
    }
    return false
}
```

**Rationale:**
- Spec is still in draft, no stable versions exist
- YAGNI principle - don't build multi-version support prematurely
- Easy to refactor when v1 is released
- Clear error messages guide users to supported versions

**Future:** When `ai-resource/v1` is released, add it to `SupportedVersions` and handle any breaking changes in validation/resolution logic.

## Core Types

### Resource Envelope

```go
type Resource struct {
    APIVersion string
    Kind       Kind
    Metadata   Metadata
    Spec       interface{} // Use type assertions or accessors
}

type Kind string

const (
    KindPrompt     Kind = "Prompt"
    KindPromptset  Kind = "Promptset"
    KindRule       Kind = "Rule"
    KindRuleset    Kind = "Ruleset"
)

type Metadata struct {
    ID          string  // Required, pattern: ^[a-zA-Z0-9_-]+$
    Name        string  // Optional
    Description string  // Optional
}

// Type-safe accessors
func (r *Resource) AsPrompt() (*Prompt, error)
func (r *Resource) AsPromptset() (*Promptset, error)
func (r *Resource) AsRule() (*Rule, error)
func (r *Resource) AsRuleset() (*Ruleset, error)
```

### Prompt Types

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

### Rule Types

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
    Priority    int         // Default: 100
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
    Files []string // Glob patterns
}
```

### Fragment Types

```go
type Fragment struct {
    Inputs map[string]InputDefinition
    Body   string // Mustache template
}

type InputDefinition struct {
    Type       InputType
    Required   bool
    Default    interface{}
    Items      *InputDefinition            // For array type
    Properties map[string]InputDefinition  // For object type
}

type InputType string

const (
    InputTypeString  InputType = "string"
    InputTypeNumber  InputType = "number"
    InputTypeBoolean InputType = "boolean"
    InputTypeArray   InputType = "array"
    InputTypeObject  InputType = "object"
)

type Body interface{} // string | []BodyItem

type BodyItem interface{} // string | FragmentRef

type FragmentRef struct {
    Fragment string
    Inputs   map[string]interface{}
}
```

## Public API

### Core Functions

```go
// Load and validate from file
func LoadResource(path string, opts ...LoadOption) (*Resource, error)
func LoadResources(path string, opts ...LoadOption) ([]*Resource, error) // Multi-doc YAML

// Parse from bytes
func ParseResource(data []byte, opts ...LoadOption) (*Resource, error)
func ParseResources(data []byte, opts ...LoadOption) ([]*Resource, error)

// Validate
func Validate(resource *Resource) error

// Fragment resolution
func ResolveBody(body Body, fragments map[string]Fragment) (string, error)

// Type-safe loading
func LoadPrompt(path string, opts ...LoadOption) (*Prompt, error)
func LoadPromptset(path string, opts ...LoadOption) (*Promptset, error)
func LoadRule(path string, opts ...LoadOption) (*Rule, error)
func LoadRuleset(path string, opts ...LoadOption) (*Ruleset, error)
```

### Configuration Options

```go
type LoadOptions struct {
    MaxFileSize     int64
    MaxArraySize    int
    MaxNestingDepth int
    Timeout         time.Duration
}

type LoadOption func(*LoadOptions)

func WithMaxFileSize(size int64) LoadOption
func WithMaxArraySize(size int) LoadOption
func WithMaxNestingDepth(depth int) LoadOption
func WithTimeout(timeout time.Duration) LoadOption

// Defaults
const (
    DefaultMaxFileSize     = 1 * 1024 * 1024 // 1MB
    DefaultMaxArraySize    = 1000
    DefaultMaxNestingDepth = 10
    DefaultTimeout         = 5 * time.Second
)
```

### Error Types

```go
type ValidationError struct {
    Field   string
    Message string
    Cause   error
}

type FragmentError struct {
    FragmentID string
    Message    string
    Cause      error
}

type InputError struct {
    FragmentID string
    InputName  string
    Expected   string
    Got        string
}

type SchemaError struct {
    Path    string
    Message string
}
```

## Implementation Requirements

### 1. File Loading

- Parse YAML and JSON
- Support multi-document YAML (separated by `---`)
- Validate file size limits
- Check apiVersion is supported
- Return clear parse errors

**Version validation:**
```go
func LoadResource(path string, opts ...LoadOption) (*Resource, error) {
    // Parse YAML/JSON
    var raw map[string]interface{}
    // ...
    
    apiVersion, ok := raw["apiVersion"].(string)
    if !ok {
        return nil, errors.New("missing or invalid apiVersion")
    }
    
    if !IsSupportedVersion(apiVersion) {
        return nil, fmt.Errorf("unsupported apiVersion: %s (supported: %v)", 
            apiVersion, SupportedVersions)
    }
    
    // Continue with validation...
}
```

### 2. Schema Validation

- Embed JSON schemas from spec repo
- Validate against appropriate schema based on `kind`
- Use `github.com/xeipuuv/gojsonschema` or similar
- Return structured validation errors

### 3. Semantic Validation

Beyond schema validation:
- Fragment references must resolve to defined fragments
- Fragment input types must match definitions
- Required fragment inputs must be provided
- metadata.id must match pattern `^[a-zA-Z0-9_-]+$`
- Collections must have minimum 1 item

### 4. Fragment Resolution

Algorithm:
```
1. Parse body (string | array)
2. For each fragment reference:
   a. Lookup fragment in spec.fragments
   b. Validate inputs (required, types)
   c. Apply defaults
   d. Render Mustache template
3. Join resolved parts with "\n\n"
4. Return final string
```

### 5. Mustache Rendering

Use `github.com/cbroglie/mustache`

Support:
- Variable substitution: `{{path}}`
- Conditionals: `{{#focus}}...{{/focus}}`
- Array iteration (primitives): `{{#files}}{{.}}{{/files}}`
- Array iteration (objects): `{{#tasks}}{{name}}: {{description}}{{/tasks}}`

### 6. Type Validation

Validate fragment inputs against definitions:
- string, number, boolean primitives
- Arrays with item type validation
- Objects with property type validation
- Recursive validation for nested structures

### 7. Error Messages

Must be clear and actionable:
```
Error: Fragment 'read-file' not found in spec.fragments
Error: Required input 'path' not provided for fragment 'read-file'
Error: Input 'count' expects type 'number', got 'string'
Error: metadata.id 'my prompt' does not match pattern ^[a-zA-Z0-9_-]+$
```

## Testing Requirements

### Unit Tests

```go
// resource_test.go
func TestLoadPrompt(t *testing.T)
func TestLoadPromptset(t *testing.T)
func TestLoadRule(t *testing.T)
func TestLoadRuleset(t *testing.T)
func TestLoadMultiDocument(t *testing.T)

// version_test.go
func TestSupportedVersion(t *testing.T)
func TestUnsupportedVersion(t *testing.T)

// fragment_test.go
func TestResolveBody(t *testing.T)
func TestFragmentInputValidation(t *testing.T)
func TestMustacheRendering(t *testing.T)
func TestFragmentDefaults(t *testing.T)

// validator_test.go
func TestValidatePrompt(t *testing.T)
func TestValidateInvalidResources(t *testing.T)
func TestSchemaValidation(t *testing.T)
func TestSemanticValidation(t *testing.T)
```

### Conformance Tests

Test against spec test suite from ai-resource-spec repo:
- All valid examples in `/schema/draft/tests/valid/` must pass
- All invalid examples in `/schema/draft/tests/invalid/` must fail
- All examples in `/examples/draft/` must pass

```go
// conformance_test.go
func TestSpecConformance(t *testing.T)
```

## Dependencies

```go
// go.mod
module github.com/yourusername/ai-resource-core-go

go 1.21

require (
    github.com/cbroglie/mustache v1.4.0      // Mustache templates
    github.com/xeipuuv/gojsonschema v1.2.0   // JSON Schema validation
    gopkg.in/yaml.v3 v3.0.1                  // YAML parsing
)
```

## Example Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yourusername/ai-resource-core-go/pkg/airesource"
)

func main() {
    // Load and validate
    prompt, err := airesource.LoadPrompt("prompts/summarize.yml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Resolve body with fragments
    resolved, err := airesource.ResolveBody(
        prompt.Spec.Body,
        prompt.Spec.Fragments,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(resolved)
}
```

## Design Principles

1. **Strict spec adherence** - Implement exactly what the spec defines
2. **Deterministic behavior** - Same input always produces same output
3. **Minimal public API** - Small, focused surface area
4. **Type safety** - Leverage Go's type system
5. **Clear errors** - Actionable error messages
6. **No runtime dependencies** - Embed schemas, no file I/O for validation
7. **Separation of concerns** - Core validates, doesn't execute
8. **Testability** - Comprehensive test coverage

## Key Specification Details

### API Version
- Current: `ai-resource/draft`
- Implementation supports only draft version initially
- Clear error for unsupported versions
- Easy to extend when v1 is released

### Pattern Constraints
- IDs: `^[a-zA-Z0-9_-]+$` (alphanumeric, underscore, hyphen)
- Fragment IDs: same pattern
- Map keys in collections: same pattern

### Required Fields
- All resources: `apiVersion`, `kind`, `metadata`, `spec`
- Metadata: `id`
- Prompt/Promptset: `body` (or `prompts` for set)
- Rule/Ruleset: `enforcement`, `body` (or `rules` for set)
- Fragment definition: `body`
- Fragment input definition: `type`

### Defaults
- Input `required`: false
- Rule `priority`: 100

### Body Formats
1. Simple string: `body: "text"`
2. Array of strings: `body: ["text1", "text2"]`
3. Array with fragment refs: `body: [{fragment: "id", inputs: {...}}, "text"]`

### Fragment Resolution
- Fragments are resource-scoped (not cross-resource)
- No fragment nesting (fragments can't reference other fragments)
- Join resolved parts with `\n\n`

### Rule Enforcement
- `may` - optional suggestion
- `should` - recommended practice
- `must` - required constraint

### Rule Scope
- Optional array of scope entries
- Each entry has `files` array with glob patterns
- Examples: `**/*.py`, `src/**/*.ts`, `*.json`

## Implementation Phases

### Phase 1: Core Types and Loading
- [ ] Define all types in `pkg/airesource/`
- [ ] Implement version constants and validation
- [ ] Implement YAML/JSON parsing
- [ ] Implement multi-document support
- [ ] Basic error types

### Phase 2: Schema Validation
- [ ] Embed JSON schemas
- [ ] Implement schema validator
- [ ] Structured error reporting
- [ ] Test against invalid examples

### Phase 3: Fragment System
- [ ] Fragment resolution algorithm
- [ ] Mustache template rendering
- [ ] Input validation
- [ ] Default value application

### Phase 4: Semantic Validation
- [ ] Fragment reference validation
- [ ] Input type validation
- [ ] Pattern validation
- [ ] Collection validation

### Phase 5: Testing
- [ ] Unit tests for all components
- [ ] Conformance tests against spec
- [ ] Example programs
- [ ] Documentation

### Phase 6: Polish
- [ ] Optimize performance
- [ ] Add caching where appropriate
- [ ] CLI tool for validation
- [ ] Comprehensive README

## References

- [AI Resource Specification](https://github.com/maxdunn/ai-resource-spec)
- [Spec: Envelope](https://github.com/maxdunn/ai-resource-spec/blob/main/spec/draft/envelope.md)
- [Spec: Resource Kinds](https://github.com/maxdunn/ai-resource-spec/blob/main/spec/draft/resource-kinds.md)
- [Spec: Fragments](https://github.com/maxdunn/ai-resource-spec/blob/main/spec/draft/fragments.md)
- [Spec: Rules](https://github.com/maxdunn/ai-resource-spec/blob/main/spec/draft/rules.md)
- [Spec: Validation](https://github.com/maxdunn/ai-resource-spec/blob/main/spec/draft/validation.md)
- [Spec: Implementation Guide](https://github.com/maxdunn/ai-resource-spec/blob/main/spec/draft/implementation.md)
- [JSON Schemas](https://github.com/maxdunn/ai-resource-spec/tree/main/schema/draft)
- [Examples](https://github.com/maxdunn/ai-resource-spec/tree/main/examples/draft)
- [Test Suite](https://github.com/maxdunn/ai-resource-spec/tree/main/schema/draft/tests)

## Notes

- This is the reference implementation - it sets the standard for other implementations
- Focus on correctness over performance initially
- All behavior must match the specification exactly
- When in doubt, refer to the spec and test suite
