package airesource

import (
	"strings"
	"testing"
)

func TestValidateSemantic_ValidPrompt(t *testing.T) {
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       KindPrompt,
		Metadata: Metadata{
			ID: "test-prompt",
		},
		Spec: PromptSpec{
			Body: Body{String: stringPtr("Simple body")},
		},
	}

	err := ValidateSemantic(resource)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateSemantic_InvalidMetadataID(t *testing.T) {
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       KindPrompt,
		Metadata: Metadata{
			ID: "invalid id!",
		},
		Spec: PromptSpec{
			Body: Body{String: stringPtr("Test")},
		},
	}

	err := ValidateSemantic(resource)
	if err == nil {
		t.Fatal("expected error for invalid metadata.id")
	}

	if !strings.Contains(err.Error(), "metadata.id") {
		t.Errorf("expected error to mention metadata.id, got: %v", err)
	}
}

func TestValidateSemantic_FragmentNotFound(t *testing.T) {
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       KindPrompt,
		Metadata: Metadata{
			ID: "test",
		},
		Spec: PromptSpec{
			Body: Body{
				Array: []BodyItem{
					{
						FragmentRef: &FragmentRef{
							Fragment: "missing",
							Inputs:   map[string]interface{}{},
						},
					},
				},
			},
		},
	}

	err := ValidateSemantic(resource)
	if err == nil {
		t.Fatal("expected error for missing fragment")
	}

	if !strings.Contains(err.Error(), "missing") {
		t.Errorf("expected error to mention missing fragment, got: %v", err)
	}
}

func TestValidateSemantic_MissingRequiredInput(t *testing.T) {
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       KindPrompt,
		Metadata: Metadata{
			ID: "test",
		},
		Spec: PromptSpec{
			Fragments: map[string]Fragment{
				"read": {
					Inputs: map[string]InputDefinition{
						"path": {
							Type:     "string",
							Required: true,
						},
					},
					Body: "Read {{path}}",
				},
			},
			Body: Body{
				Array: []BodyItem{
					{
						FragmentRef: &FragmentRef{
							Fragment: "read",
							Inputs:   map[string]interface{}{},
						},
					},
				},
			},
		},
	}

	err := ValidateSemantic(resource)
	if err == nil {
		t.Fatal("expected error for missing required input")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "read") || !strings.Contains(errMsg, "input") {
		t.Errorf("expected error to mention fragment and input validation, got: %v", err)
	}
}

func TestValidateSemantic_InvalidFragmentKey(t *testing.T) {
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       KindPrompt,
		Metadata: Metadata{
			ID: "test",
		},
		Spec: PromptSpec{
			Fragments: map[string]Fragment{
				"invalid key!": {
					Body: "Test",
				},
			},
			Body: Body{String: stringPtr("Test")},
		},
	}

	err := ValidateSemantic(resource)
	if err == nil {
		t.Fatal("expected error for invalid fragment key")
	}

	if !strings.Contains(err.Error(), "fragment key") {
		t.Errorf("expected error to mention fragment key, got: %v", err)
	}
}

func TestValidateSemantic_EmptyPromptset(t *testing.T) {
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       KindPromptset,
		Metadata: Metadata{
			ID: "test",
		},
		Spec: PromptsetSpec{
			Prompts: map[string]PromptItem{},
		},
	}

	err := ValidateSemantic(resource)
	if err == nil {
		t.Fatal("expected error for empty promptset")
	}

	if !strings.Contains(err.Error(), "at least one prompt") {
		t.Errorf("expected error about minimum prompts, got: %v", err)
	}
}

func TestValidateSemantic_InvalidPromptKey(t *testing.T) {
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       KindPromptset,
		Metadata: Metadata{
			ID: "test",
		},
		Spec: PromptsetSpec{
			Prompts: map[string]PromptItem{
				"invalid key!": {
					Body: Body{String: stringPtr("Test")},
				},
			},
		},
	}

	err := ValidateSemantic(resource)
	if err == nil {
		t.Fatal("expected error for invalid prompt key")
	}

	if !strings.Contains(err.Error(), "prompt key") {
		t.Errorf("expected error to mention prompt key, got: %v", err)
	}
}

func TestValidateSemantic_EmptyRuleset(t *testing.T) {
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       KindRuleset,
		Metadata: Metadata{
			ID: "test",
		},
		Spec: RulesetSpec{
			Rules: map[string]RuleItem{},
		},
	}

	err := ValidateSemantic(resource)
	if err == nil {
		t.Fatal("expected error for empty ruleset")
	}

	if !strings.Contains(err.Error(), "at least one rule") {
		t.Errorf("expected error about minimum rules, got: %v", err)
	}
}

func TestValidateSemantic_InvalidRuleKey(t *testing.T) {
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       KindRuleset,
		Metadata: Metadata{
			ID: "test",
		},
		Spec: RulesetSpec{
			Rules: map[string]RuleItem{
				"invalid key!": {
					Enforcement: "must",
					Body:        Body{String: stringPtr("Test")},
				},
			},
		},
	}

	err := ValidateSemantic(resource)
	if err == nil {
		t.Fatal("expected error for invalid rule key")
	}

	if !strings.Contains(err.Error(), "rule key") {
		t.Errorf("expected error to mention rule key, got: %v", err)
	}
}

func TestValidateSemantic_ValidPromptset(t *testing.T) {
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       KindPromptset,
		Metadata: Metadata{
			ID: "test",
		},
		Spec: PromptsetSpec{
			Prompts: map[string]PromptItem{
				"prompt1": {
					Body: Body{String: stringPtr("Test 1")},
				},
				"prompt2": {
					Body: Body{String: stringPtr("Test 2")},
				},
			},
		},
	}

	err := ValidateSemantic(resource)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateSemantic_ValidRuleset(t *testing.T) {
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       KindRuleset,
		Metadata: Metadata{
			ID: "test",
		},
		Spec: RulesetSpec{
			Rules: map[string]RuleItem{
				"rule1": {
					Enforcement: "must",
					Body:        Body{String: stringPtr("Test 1")},
				},
				"rule2": {
					Enforcement: "should",
					Body:        Body{String: stringPtr("Test 2")},
				},
			},
		},
	}

	err := ValidateSemantic(resource)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateSemantic_ValidFragmentReference(t *testing.T) {
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       KindPrompt,
		Metadata: Metadata{
			ID: "test",
		},
		Spec: PromptSpec{
			Fragments: map[string]Fragment{
				"greet": {
					Inputs: map[string]InputDefinition{
						"name": {
							Type:     "string",
							Required: true,
						},
					},
					Body: "Hello, {{name}}!",
				},
			},
			Body: Body{
				Array: []BodyItem{
					{String: stringPtr("Introduction")},
					{
						FragmentRef: &FragmentRef{
							Fragment: "greet",
							Inputs: map[string]interface{}{
								"name": "World",
							},
						},
					},
					{String: stringPtr("Conclusion")},
				},
			},
		},
	}

	err := ValidateSemantic(resource)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func stringPtr(s string) *string {
	return &s
}
