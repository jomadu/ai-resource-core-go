package schema

import (
	"strings"
	"testing"
)

func TestValidateSchema_ValidPrompt(t *testing.T) {
	resource := map[string]interface{}{
		"apiVersion": "ai-resource/draft",
		"kind":       "Prompt",
		"metadata": map[string]interface{}{
			"id": "test-prompt",
		},
		"spec": map[string]interface{}{
			"body": "Test body",
		},
	}

	err := ValidateSchema(resource)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateSchema_MissingRequiredField(t *testing.T) {
	resource := map[string]interface{}{
		"apiVersion": "ai-resource/draft",
		"kind":       "Prompt",
		"metadata": map[string]interface{}{
			"name": "Missing ID",
		},
		"spec": map[string]interface{}{
			"body": "Test",
		},
	}

	err := ValidateSchema(resource)
	if err == nil {
		t.Fatal("expected error for missing metadata.id, got nil")
	}

	schemaErr, ok := err.(*SchemaError)
	if !ok {
		t.Fatalf("expected SchemaError, got %T", err)
	}

	if !strings.Contains(schemaErr.Path, "metadata") {
		t.Errorf("expected path to contain 'metadata', got: %s", schemaErr.Path)
	}
}

func TestValidateSchema_InvalidPattern(t *testing.T) {
	resource := map[string]interface{}{
		"apiVersion": "ai-resource/draft",
		"kind":       "Prompt",
		"metadata": map[string]interface{}{
			"id": "invalid id with spaces",
		},
		"spec": map[string]interface{}{
			"body": "Test",
		},
	}

	err := ValidateSchema(resource)
	if err == nil {
		t.Fatal("expected error for invalid metadata.id pattern, got nil")
	}

	schemaErr, ok := err.(*SchemaError)
	if !ok {
		t.Fatalf("expected SchemaError, got %T", err)
	}

	if !strings.Contains(schemaErr.Path, "metadata.id") {
		t.Errorf("expected path to contain 'metadata.id', got: %s", schemaErr.Path)
	}
}

func TestValidateSchema_InvalidEnumValue(t *testing.T) {
	resource := map[string]interface{}{
		"apiVersion": "ai-resource/draft",
		"kind":       "Rule",
		"metadata": map[string]interface{}{
			"id": "test-rule",
		},
		"spec": map[string]interface{}{
			"enforcement": "always",
			"body":        "Test",
		},
	}

	err := ValidateSchema(resource)
	if err == nil {
		t.Fatal("expected error for invalid enforcement value, got nil")
	}

	schemaErr, ok := err.(*SchemaError)
	if !ok {
		t.Fatalf("expected SchemaError, got %T", err)
	}

	if !strings.Contains(schemaErr.Path, "spec.enforcement") {
		t.Errorf("expected path to contain 'spec.enforcement', got: %s", schemaErr.Path)
	}
}

func TestValidateSchema_MultipleErrors(t *testing.T) {
	resource := map[string]interface{}{
		"apiVersion": "ai-resource/draft",
		"kind":       "Rule",
		"metadata": map[string]interface{}{
			"id": "invalid id!",
		},
		"spec": map[string]interface{}{
			"enforcement": "invalid",
		},
	}

	err := ValidateSchema(resource)
	if err == nil {
		t.Fatal("expected errors, got nil")
	}

	multiErr, ok := err.(*MultiError)
	if !ok {
		t.Fatalf("expected MultiError, got %T", err)
	}

	if len(multiErr.Errors) < 2 {
		t.Errorf("expected at least 2 errors, got %d", len(multiErr.Errors))
	}
}

func TestValidateSchema_UnknownKind(t *testing.T) {
	resource := map[string]interface{}{
		"apiVersion": "ai-resource/draft",
		"kind":       "Unknown",
		"metadata": map[string]interface{}{
			"id": "test",
		},
		"spec": map[string]interface{}{},
	}

	err := ValidateSchema(resource)
	if err == nil {
		t.Fatal("expected error for unknown kind, got nil")
	}

	if !strings.Contains(err.Error(), "no schema for kind") {
		t.Errorf("expected 'no schema for kind' error, got: %v", err)
	}
}

func TestValidateSchema_ValidRule(t *testing.T) {
	resource := map[string]interface{}{
		"apiVersion": "ai-resource/draft",
		"kind":       "Rule",
		"metadata": map[string]interface{}{
			"id": "test-rule",
		},
		"spec": map[string]interface{}{
			"enforcement": "must",
			"body":        "Test rule",
		},
	}

	err := ValidateSchema(resource)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateSchema_ValidPromptset(t *testing.T) {
	resource := map[string]interface{}{
		"apiVersion": "ai-resource/draft",
		"kind":       "Promptset",
		"metadata": map[string]interface{}{
			"id": "test-promptset",
		},
		"spec": map[string]interface{}{
			"prompts": map[string]interface{}{
				"greeting": map[string]interface{}{
					"body": "Hello",
				},
			},
		},
	}

	err := ValidateSchema(resource)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateSchema_ValidRuleset(t *testing.T) {
	resource := map[string]interface{}{
		"apiVersion": "ai-resource/draft",
		"kind":       "Ruleset",
		"metadata": map[string]interface{}{
			"id": "test-ruleset",
		},
		"spec": map[string]interface{}{
			"rules": map[string]interface{}{
				"no-console": map[string]interface{}{
					"enforcement": "should",
					"body":        "No console.log",
				},
			},
		},
	}

	err := ValidateSchema(resource)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
