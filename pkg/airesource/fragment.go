package airesource

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

// FragmentRef represents a reference to a fragment with inputs.
type FragmentRef struct {
	Fragment string                 `yaml:"fragment"`
	Inputs   map[string]interface{} `yaml:"inputs,omitempty"`
}
