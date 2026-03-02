package airesource

import (
	"encoding/json"
	"strings"

	"github.com/aws/ai-resource-core-go/internal/schema"
	"github.com/aws/ai-resource-core-go/internal/template"
)

// Fragment represents a reusable template with typed inputs.
type Fragment struct {
	Inputs map[string]InputDefinition `yaml:"inputs,omitempty" json:"inputs,omitempty"`
	Body   string                     `yaml:"body" json:"body"`
}

// InputDefinition defines the type and constraints for a fragment input.
type InputDefinition struct {
	Type       InputType                  `yaml:"type" json:"type"`
	Required   bool                       `yaml:"required,omitempty" json:"required,omitempty"`
	Default    interface{}                `yaml:"default,omitempty" json:"default,omitempty"`
	Items      *InputDefinition           `yaml:"items,omitempty" json:"items,omitempty"`
	Properties map[string]InputDefinition `yaml:"properties,omitempty" json:"properties,omitempty"`
}

// InputType represents the type of a fragment input.
type InputType string

const (
	InputTypeString  InputType = "string"
	InputTypeNumber  InputType = "number"
	InputTypeBoolean InputType = "boolean"
	InputTypeArray   InputType = "array"
	InputTypeObject  InputType = "object"
)

// Body represents the content of a prompt or rule.
// Exactly one of String or Array must be non-nil.
type Body struct {
	String *string
	Array  []BodyItem
}

// UnmarshalYAML implements custom YAML unmarshaling for Body.
func (b *Body) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err == nil {
		b.String = &s
		return nil
	}

	var arr []BodyItem
	if err := unmarshal(&arr); err == nil {
		b.Array = arr
		return nil
	}

	return &ValidationError{Field: "body", Message: "must be string or array"}
}

// MarshalJSON implements custom JSON marshaling for Body.
func (b Body) MarshalJSON() ([]byte, error) {
	if b.String != nil {
		return json.Marshal(*b.String)
	}
	if b.Array != nil {
		return json.Marshal(b.Array)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements custom JSON unmarshaling for Body.
func (b *Body) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		b.String = &s
		return nil
	}

	var arr []BodyItem
	if err := json.Unmarshal(data, &arr); err == nil {
		b.Array = arr
		return nil
	}

	return &ValidationError{Field: "body", Message: "must be string or array"}
}

// BodyItem represents an item in a body array.
// Exactly one of String or FragmentRef must be non-nil.
type BodyItem struct {
	String      *string
	FragmentRef *FragmentRef
}

// UnmarshalYAML implements custom YAML unmarshaling for BodyItem.
func (bi *BodyItem) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err == nil {
		bi.String = &s
		return nil
	}

	var ref FragmentRef
	if err := unmarshal(&ref); err == nil {
		bi.FragmentRef = &ref
		return nil
	}

	return &ValidationError{Field: "body item", Message: "must be string or fragment reference"}
}

// MarshalJSON implements custom JSON marshaling for BodyItem.
func (bi BodyItem) MarshalJSON() ([]byte, error) {
	if bi.String != nil {
		return json.Marshal(*bi.String)
	}
	if bi.FragmentRef != nil {
		return json.Marshal(bi.FragmentRef)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements custom JSON unmarshaling for BodyItem.
func (bi *BodyItem) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		bi.String = &s
		return nil
	}

	var ref FragmentRef
	if err := json.Unmarshal(data, &ref); err == nil {
		bi.FragmentRef = &ref
		return nil
	}

	return &ValidationError{Field: "body item", Message: "must be string or fragment reference"}
}

// BuildSchemaFromInputs converts a map of InputDefinitions to a JSON Schema.
func BuildSchemaFromInputs(inputs map[string]InputDefinition) map[string]interface{} {
	properties := make(map[string]interface{})
	required := []string{}

	for name, def := range inputs {
		properties[name] = convertInputDefToSchema(def)
		if def.Required {
			required = append(required, name)
		}
	}

	schema := map[string]interface{}{
		"type":                 "object",
		"properties":           properties,
		"additionalProperties": false,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

func convertInputDefToSchema(def InputDefinition) map[string]interface{} {
	schema := map[string]interface{}{
		"type": string(def.Type),
	}

	if def.Type == InputTypeArray && def.Items != nil {
		schema["items"] = convertInputDefToSchema(*def.Items)
	}

	if def.Type == InputTypeObject && def.Properties != nil {
		properties := make(map[string]interface{})
		for propName, propDef := range def.Properties {
			properties[propName] = convertInputDefToSchema(propDef)
		}
		schema["properties"] = properties
	}

	return schema
}

type FragmentRef struct {
	Fragment string                 `yaml:"fragment" json:"fragment"`
	Inputs   map[string]interface{} `yaml:"inputs,omitempty" json:"inputs,omitempty"`
}

// ValidateInputs validates provided inputs against fragment input definitions.
// It applies default values and returns the validated inputs map.
func ValidateInputs(fragmentID string, fragment Fragment, providedInputs map[string]interface{}) (map[string]interface{}, error) {
	inputsWithDefaults := applyDefaults(providedInputs, fragment.Inputs)
	
	// Check for undefined inputs
	for name := range inputsWithDefaults {
		if _, exists := fragment.Inputs[name]; !exists {
			return nil, &InputError{
				FragmentID: fragmentID,
				InputName:  name,
				Expected:   "defined input",
				Got:        "undefined input",
			}
		}
	}
	
	// Build JSON Schema from input definitions
	jsonSchema := BuildSchemaFromInputs(fragment.Inputs)
	
	// Validate using JSON Schema
	if err := schema.ValidateAgainstSchema(inputsWithDefaults, jsonSchema); err != nil {
		return nil, convertSchemaErrorToInputError(fragmentID, err)
	}
	
	return inputsWithDefaults, nil
}

func convertSchemaErrorToInputError(fragmentID string, err error) error {
	if schemaErr, ok := err.(*schema.SchemaError); ok {
		inputName := strings.TrimPrefix(schemaErr.Path, "(root).")
		if inputName == "(root)" {
			inputName = ""
		}
		
		expected := "valid value"
		got := "invalid value"
		
		// Parse "Invalid type. Expected: number, given: string"
		if strings.Contains(schemaErr.Message, "Invalid type") {
			parts := strings.Split(schemaErr.Message, ",")
			if len(parts) >= 2 {
				// Extract expected type
				expectedPart := strings.TrimSpace(parts[0])
				if idx := strings.Index(expectedPart, "Expected:"); idx != -1 {
					expected = strings.TrimSpace(expectedPart[idx+9:])
				}
				// Extract given type
				givenPart := strings.TrimSpace(parts[1])
				if idx := strings.Index(givenPart, "given:"); idx != -1 {
					got = strings.TrimSpace(givenPart[idx+6:])
				}
			}
		} else if strings.Contains(schemaErr.Message, "required") {
			expected = "required input"
			got = "missing"
		}
		
		return &InputError{
			FragmentID: fragmentID,
			InputName:  inputName,
			Expected:   expected,
			Got:        got,
		}
	}
	
	if multiErr, ok := err.(*schema.MultiError); ok {
		if len(multiErr.Errors) > 0 {
			return convertSchemaErrorToInputError(fragmentID, multiErr.Errors[0])
		}
	}
	
	return &InputError{
		FragmentID: fragmentID,
		InputName:  "",
		Expected:   "valid input",
		Got:        err.Error(),
	}
}

func applyDefaults(provided map[string]interface{}, definitions map[string]InputDefinition) map[string]interface{} {
	result := make(map[string]interface{})
	
	for k, v := range provided {
		result[k] = v
	}
	
	for name, def := range definitions {
		if _, exists := result[name]; !exists && def.Default != nil {
			result[name] = def.Default
		}
	}
	
	return result
}

// ResolveBody resolves a body by processing fragment references and rendering templates.
// Returns the final rendered string.
func ResolveBody(body Body, fragments map[string]Fragment) (string, error) {
	if body.String != nil {
		return *body.String, nil
	}

	if body.Array == nil {
		return "", &FragmentError{
			Message: "body must have String or Array set",
		}
	}

	var parts []string
	for _, item := range body.Array {
		if item.String != nil {
			parts = append(parts, *item.String)
		} else if item.FragmentRef != nil {
			fragment, exists := fragments[item.FragmentRef.Fragment]
			if !exists {
				return "", &FragmentError{
					FragmentID: item.FragmentRef.Fragment,
					Message:    "fragment not found",
				}
			}

			rendered, err := template.Render(fragment.Body, item.FragmentRef.Inputs)
			if err != nil {
				return "", &FragmentError{
					FragmentID: item.FragmentRef.Fragment,
					Message:    "template rendering failed",
					Cause:      err,
				}
			}
			parts = append(parts, rendered)
		} else {
			return "", &FragmentError{
				Message: "BodyItem must have String or FragmentRef set",
			}
		}
	}

	return strings.Join(parts, "\n\n"), nil
}
