package airesource

import "encoding/json"

// Fragment represents a reusable template with typed inputs.
type Fragment struct {
	Inputs map[string]InputDefinition `yaml:"inputs,omitempty"`
	Body   string                     `yaml:"body"`
}

// InputDefinition defines the type and constraints for a fragment input.
type InputDefinition struct {
	Type       InputType                  `yaml:"type"`
	Required   bool                       `yaml:"required,omitempty"`
	Default    interface{}                `yaml:"default,omitempty"`
	Items      *InputDefinition           `yaml:"items,omitempty"`
	Properties map[string]InputDefinition `yaml:"properties,omitempty"`
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
	Fragment string                 `yaml:"fragment"`
	Inputs   map[string]interface{} `yaml:"inputs,omitempty"`
}

// ValidateInputs validates provided inputs against fragment input definitions.
// It applies default values and returns the validated inputs map.
func ValidateInputs(fragmentID string, fragment Fragment, providedInputs map[string]interface{}) (map[string]interface{}, error) {
	inputsWithDefaults := applyDefaults(providedInputs, fragment.Inputs)
	
	if err := validateInputTypes(fragmentID, fragment.Inputs, inputsWithDefaults); err != nil {
		return nil, err
	}
	
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
	
	return inputsWithDefaults, nil
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

func validateInputTypes(fragmentID string, definitions map[string]InputDefinition, inputs map[string]interface{}) error {
	for name, def := range definitions {
		value, exists := inputs[name]
		
		if !exists {
			if def.Required {
				return &InputError{
					FragmentID: fragmentID,
					InputName:  name,
					Expected:   "required input",
					Got:        "missing",
				}
			}
			continue
		}
		
		if err := validateType(fragmentID, name, def, value); err != nil {
			return err
		}
	}
	
	return nil
}

func validateType(fragmentID, inputName string, def InputDefinition, value interface{}) error {
	switch def.Type {
	case InputTypeString:
		if _, ok := value.(string); !ok {
			return &InputError{
				FragmentID: fragmentID,
				InputName:  inputName,
				Expected:   "string",
				Got:        getTypeName(value),
			}
		}
	case InputTypeNumber:
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			// Valid number types
		default:
			return &InputError{
				FragmentID: fragmentID,
				InputName:  inputName,
				Expected:   "number",
				Got:        getTypeName(value),
			}
		}
	case InputTypeBoolean:
		if _, ok := value.(bool); !ok {
			return &InputError{
				FragmentID: fragmentID,
				InputName:  inputName,
				Expected:   "boolean",
				Got:        getTypeName(value),
			}
		}
	case InputTypeArray:
		arr, ok := value.([]interface{})
		if !ok {
			return &InputError{
				FragmentID: fragmentID,
				InputName:  inputName,
				Expected:   "array",
				Got:        getTypeName(value),
			}
		}
		if def.Items != nil {
			for i, item := range arr {
				if err := validateType(fragmentID, inputName+"["+string(rune(i))+ "]", *def.Items, item); err != nil {
					return err
				}
			}
		}
	case InputTypeObject:
		obj, ok := value.(map[string]interface{})
		if !ok {
			return &InputError{
				FragmentID: fragmentID,
				InputName:  inputName,
				Expected:   "object",
				Got:        getTypeName(value),
			}
		}
		if def.Properties != nil {
			for propName, propDef := range def.Properties {
				if propValue, exists := obj[propName]; exists {
					if err := validateType(fragmentID, inputName+"."+propName, propDef, propValue); err != nil {
						return err
					}
				}
			}
		}
	}
	
	return nil
}

func getTypeName(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return "number"
	case bool:
		return "boolean"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}
