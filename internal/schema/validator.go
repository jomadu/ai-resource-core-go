package schema

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type Resource struct {
	APIVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Metadata   interface{} `json:"metadata"`
	Spec       interface{} `json:"spec"`
}

type SchemaError struct {
	Path    string
	Message string
}

func (e *SchemaError) Error() string {
	return fmt.Sprintf("schema validation failed at %s: %s", e.Path, e.Message)
}

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return "no errors"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	msg := fmt.Sprintf("%d errors occurred:", len(e.Errors))
	for _, err := range e.Errors {
		msg += "\n  - " + err.Error()
	}
	return msg
}

func (e *MultiError) Unwrap() []error {
	return e.Errors
}

func ValidateSchema(resource interface{}) error {
	resourceJSON, err := json.Marshal(resource)
	if err != nil {
		return fmt.Errorf("failed to marshal resource: %w", err)
	}

	var resData map[string]interface{}
	if err := json.Unmarshal(resourceJSON, &resData); err != nil {
		return fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	kind, ok := resData["kind"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid kind field")
	}

	schema, err := getSchemaForKind(kind)
	if err != nil {
		return err
	}

	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(resourceJSON)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	if !result.Valid() {
		errors := make([]error, 0, len(result.Errors()))
		for _, err := range result.Errors() {
			errors = append(errors, &SchemaError{
				Path:    err.Field(),
				Message: err.Description(),
			})
		}

		if len(errors) == 1 {
			return errors[0]
		}

		return &MultiError{Errors: errors}
	}

	return nil
}

func getSchemaForKind(kind string) (string, error) {
	switch kind {
	case "Prompt":
		return promptSchema, nil
	case "Promptset":
		return promptsetSchema, nil
	case "Rule":
		return ruleSchema, nil
	case "Ruleset":
		return rulesetSchema, nil
	default:
		return "", fmt.Errorf("no schema for kind: %s", kind)
	}
}
