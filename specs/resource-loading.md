# Resource Loading

## Job to be Done
Enable developers to load AI Resources from YAML or JSON files into validated, type-safe Go structures.

## Activities
- Parse YAML and JSON file formats
- Support single-document files
- Support multi-document YAML files (separated by `---`)
- Enforce file size limits for safety
- Validate API version is supported
- Map parsed data to appropriate resource types
- Return clear parse errors with file context

## Acceptance Criteria
- [ ] YAML files are parsed correctly
- [ ] JSON files are parsed correctly
- [ ] Multi-document YAML files return array of resources
- [ ] Files exceeding size limit are rejected with clear error
- [ ] Unsupported API versions return error listing supported versions
- [ ] Missing required envelope fields return parse errors
- [ ] Invalid YAML/JSON syntax returns actionable error messages
- [ ] File not found returns clear error
- [ ] Loaded resources have correct Kind and APIVersion set

## Data Structures

### LoadOptions
```go
type LoadOptions struct {
    MaxFileSize     int64
    MaxArraySize    int
    MaxNestingDepth int
    Timeout         time.Duration
}
```

**Fields:**
- `MaxFileSize` - Maximum file size in bytes (default: 10MB, for extreme cases only)
- `MaxArraySize` - Maximum array length (default: 10000, for extreme cases only)
- `MaxNestingDepth` - Maximum nesting depth (default: 100, for extreme cases only)
- `Timeout` - Operation timeout (default: 30s)

### LoadOption
```go
type LoadOption func(*LoadOptions)
```

**Functions:**
- `WithMaxFileSize(size int64)` - Override max file size
- `WithMaxArraySize(size int)` - Override max array size
- `WithMaxNestingDepth(depth int)` - Override max nesting depth
- `WithTimeout(timeout time.Duration)` - Override timeout

## Algorithm

### Single Resource Loading

1. Read file from path
2. Check file size against limit
3. Detect format (YAML vs JSON) from file extension
4. Parse into raw map structure
5. Extract and validate apiVersion field
6. Check apiVersion is in supported versions list
7. Extract and validate kind field
8. Map to appropriate resource type based on kind
9. Return typed resource or error

**Pseudocode:**
```
function LoadResource(path, options):
    data = read_file(path)
    
    if len(data) > options.MaxFileSize:
        return error("file exceeds size limit")
    
    raw = parse_yaml_or_json(data)
    
    apiVersion = raw["apiVersion"]
    if not apiVersion:
        return error("missing apiVersion")
    
    if not is_supported_version(apiVersion):
        return error("unsupported apiVersion: {apiVersion}")
    
    kind = raw["kind"]
    if not kind:
        return error("missing kind")
    
    resource = map_to_resource_type(raw, kind)
    return resource
```

### Multi-Document Loading

1. Read file from path
2. Check file size against limit
3. Split YAML by document separator `---`
4. For each document:
   - Parse as single resource
   - Collect in array
5. Return array of resources or error

**Pseudocode:**
```
function LoadResources(path, options):
    data = read_file(path)
    
    if len(data) > options.MaxFileSize:
        return error("file exceeds size limit")
    
    documents = split_yaml_documents(data)
    resources = []
    
    for doc in documents:
        resource = parse_single_resource(doc, options)
        resources.append(resource)
    
    return resources
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| File does not exist | Return error: "file not found: {path}" |
| File exceeds size limit | Return error: "file size {size} exceeds limit {limit}" |
| Empty file | Return error: "empty file" |
| Invalid YAML syntax | Return error with line/column information |
| Invalid JSON syntax | Return error with position information |
| Unsupported apiVersion | Return error: "unsupported apiVersion: {version} (supported: {list})" |
| Missing apiVersion | Return error: "missing required field: apiVersion" |
| Missing kind | Return error: "missing required field: kind" |
| Invalid kind value | Return error: "invalid kind: {value}" |
| Multi-doc with one invalid | Return error indicating which document failed |
| Binary file content | Return error: "invalid YAML/JSON format" |

## Dependencies

- `core-types.md` - Defines Resource, Kind, and type structures
- YAML parsing library (gopkg.in/yaml.v3)
- JSON parsing (encoding/json)

## Implementation Mapping

**Source files:**
- `pkg/airesource/loader.go` - Main loading functions
- `pkg/airesource/version.go` - Version validation
- `pkg/airesource/options.go` - LoadOptions and functional options

**Related specs:**
- `core-types.md` - Types being loaded
- `schema-validation.md` - Validation after loading
- `error-handling.md` - Error types returned

## Examples

### Example 1: Load Single Prompt

**Input:**
```yaml
# prompt.yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: summarize
  name: Summarize Text
spec:
  body: "Summarize the following text in 3-5 sentences."
```

```go
prompt, err := LoadPrompt("prompt.yaml")
```

**Expected Output:**
```go
prompt.APIVersion == "ai-resource/draft"
prompt.Kind == KindPrompt
prompt.Metadata.ID == "summarize"
err == nil
```

**Verification:**
- Check all fields are populated correctly
- Check error is nil

### Example 2: Load Multi-Document YAML

**Input:**
```yaml
# resources.yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: prompt1
spec:
  body: "First prompt"
---
apiVersion: ai-resource/draft
kind: Rule
metadata:
  id: rule1
spec:
  enforcement: must
  body: "First rule"
```

```go
resources, err := LoadResources("resources.yaml")
```

**Expected Output:**
```go
len(resources) == 2
resources[0].Kind == KindPrompt
resources[1].Kind == KindRule
err == nil
```

**Verification:**
- Check array length
- Check each resource kind
- Check error is nil

### Example 3: Unsupported Version Error

**Input:**
```yaml
# future.yaml
apiVersion: ai-resource/v2
kind: Prompt
metadata:
  id: test
spec:
  body: "test"
```

```go
_, err := LoadResource("future.yaml")
```

**Expected Output:**
```go
err.Error() == "unsupported apiVersion: ai-resource/v2 (supported: [ai-resource/draft])"
```

**Verification:**
- Check error message contains unsupported version
- Check error message lists supported versions

### Example 4: File Size Limit

**Input:**
```go
_, err := LoadResource("huge.yaml", WithMaxFileSize(1024))
```

**Expected Output:**
```go
// If file is 2048 bytes:
err.Error() contains "exceeds size limit"
```

**Verification:**
- Check error indicates size limit exceeded
- Check file is not fully loaded into memory

## Notes

- File format is determined by extension (.yaml, .yml, .json)
- Multi-document support is YAML-only (JSON has no standard multi-document format)
- Size limits prevent DoS attacks from maliciously large files
- Timeout prevents hanging on slow filesystem operations
- Version validation happens early to fail fast on incompatible resources
- The loader does NOT perform schema or semantic validation - it only ensures structural parsing and version compatibility

## Known Issues

None.

## Areas for Improvement

- Could add content sniffing as fallback for files without extensions
- Could add streaming support for very large multi-document files
- Could add caching for repeatedly loaded files
- Could support loading from io.Reader for testing
- Could add progress callbacks for large file operations
