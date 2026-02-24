package airesource

import (
	"fmt"
	"regexp"
)

var idPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateSemantic performs semantic validation on a resource.
// This includes fragment reference checking, collection validation,
// and other rules that go beyond structural JSON Schema validation.
func ValidateSemantic(resource *Resource) error {
	var errors []error

	// Validate metadata.id
	if !idPattern.MatchString(resource.Metadata.ID) {
		errors = append(errors, &ValidationError{
			Field:   "metadata.id",
			Message: fmt.Sprintf("does not match pattern ^[a-zA-Z0-9_-]+$: %q", resource.Metadata.ID),
		})
	}

	// Kind-specific validation
	switch resource.Kind {
	case KindPrompt:
		prompt, err := resource.AsPrompt()
		if err != nil {
			errors = append(errors, err)
		} else {
			errors = append(errors, validatePromptSpec(&prompt.Spec)...)
		}
	case KindPromptset:
		promptset, err := resource.AsPromptset()
		if err != nil {
			errors = append(errors, err)
		} else {
			errors = append(errors, validatePromptsetSpec(&promptset.Spec)...)
		}
	case KindRule:
		rule, err := resource.AsRule()
		if err != nil {
			errors = append(errors, err)
		} else {
			errors = append(errors, validateRuleSpec(&rule.Spec)...)
		}
	case KindRuleset:
		ruleset, err := resource.AsRuleset()
		if err != nil {
			errors = append(errors, err)
		} else {
			errors = append(errors, validateRulesetSpec(&ruleset.Spec)...)
		}
	}

	if len(errors) > 0 {
		return &MultiError{Errors: errors}
	}
	return nil
}

func validatePromptSpec(spec *PromptSpec) []error {
	var errors []error

	// Validate fragment keys
	for key := range spec.Fragments {
		if !idPattern.MatchString(key) {
			errors = append(errors, &ValidationError{
				Field:   fmt.Sprintf("spec.fragments[%s]", key),
				Message: fmt.Sprintf("fragment key does not match pattern ^[a-zA-Z0-9_-]+$: %q", key),
			})
		}
	}

	// Validate body
	errors = append(errors, validateBody(spec.Body, spec.Fragments)...)

	return errors
}

func validateRuleSpec(spec *RuleSpec) []error {
	var errors []error

	// Validate fragment keys
	for key := range spec.Fragments {
		if !idPattern.MatchString(key) {
			errors = append(errors, &ValidationError{
				Field:   fmt.Sprintf("spec.fragments[%s]", key),
				Message: fmt.Sprintf("fragment key does not match pattern ^[a-zA-Z0-9_-]+$: %q", key),
			})
		}
	}

	// Validate body
	errors = append(errors, validateBody(spec.Body, spec.Fragments)...)

	return errors
}

func validateBody(body Body, fragments map[string]Fragment) []error {
	var errors []error

	if body.String != nil {
		return nil
	}

	if body.Array == nil {
		errors = append(errors, &ValidationError{
			Field:   "spec.body",
			Message: "body must be string or array",
		})
		return errors
	}

	for i, item := range body.Array {
		if item.String != nil {
			continue
		}

		if item.FragmentRef != nil {
			errs := validateFragmentRef(i, item.FragmentRef, fragments)
			errors = append(errors, errs...)
		} else {
			errors = append(errors, &ValidationError{
				Field:   fmt.Sprintf("spec.body[%d]", i),
				Message: "body item must be string or fragment reference",
			})
		}
	}

	return errors
}

func validateFragmentRef(index int, ref *FragmentRef, fragments map[string]Fragment) []error {
	var errors []error

	fragment, exists := fragments[ref.Fragment]
	if !exists {
		errors = append(errors, &FragmentError{
			FragmentID: ref.Fragment,
			Message:    fmt.Sprintf("fragment %q not found in spec.fragments", ref.Fragment),
		})
		return errors
	}

	// Validate fragment inputs
	_, err := ValidateInputs(ref.Fragment, fragment, ref.Inputs)
	if err != nil {
		errors = append(errors, &ValidationError{
			Field:   fmt.Sprintf("spec.body[%d].inputs", index),
			Message: fmt.Sprintf("fragment %q input validation failed", ref.Fragment),
			Cause:   err,
		})
	}

	return errors
}

func validatePromptset(promptset *Promptset) []error {
	var errors []error

	// Validate fragment keys
	for key := range promptset.Spec.Fragments {
		if !idPattern.MatchString(key) {
			errors = append(errors, &ValidationError{
				Field:   fmt.Sprintf("spec.fragments[%s]", key),
				Message: fmt.Sprintf("fragment key does not match pattern ^[a-zA-Z0-9_-]+$: %q", key),
			})
		}
	}

	// Validate collection keys
	for key := range promptset.Spec.Prompts {
		if !idPattern.MatchString(key) {
			errors = append(errors, &ValidationError{
				Field:   fmt.Sprintf("spec.prompts[%s]", key),
				Message: fmt.Sprintf("prompt key does not match pattern ^[a-zA-Z0-9_-]+$: %q", key),
			})
		}
	}

	// Validate minimum one prompt
	if len(promptset.Spec.Prompts) == 0 {
		errors = append(errors, &ValidationError{
			Field:   "spec.prompts",
			Message: "promptset must have at least one prompt",
		})
	}

	// Validate each prompt body
	for key, prompt := range promptset.Spec.Prompts {
		bodyErrors := validateBody(prompt.Body, promptset.Spec.Fragments)
		for _, err := range bodyErrors {
			if ve, ok := err.(*ValidationError); ok {
				ve.Field = fmt.Sprintf("spec.prompts[%s].%s", key, ve.Field)
			}
			errors = append(errors, err)
		}
	}

	return errors
}

func validatePromptsetSpec(spec *PromptsetSpec) []error {
	var errors []error

	// Validate fragment keys
	for key := range spec.Fragments {
		if !idPattern.MatchString(key) {
			errors = append(errors, &ValidationError{
				Field:   fmt.Sprintf("spec.fragments[%s]", key),
				Message: fmt.Sprintf("fragment key does not match pattern ^[a-zA-Z0-9_-]+$: %q", key),
			})
		}
	}

	// Validate collection keys
	for key := range spec.Prompts {
		if !idPattern.MatchString(key) {
			errors = append(errors, &ValidationError{
				Field:   fmt.Sprintf("spec.prompts[%s]", key),
				Message: fmt.Sprintf("prompt key does not match pattern ^[a-zA-Z0-9_-]+$: %q", key),
			})
		}
	}

	// Validate minimum one prompt
	if len(spec.Prompts) == 0 {
		errors = append(errors, &ValidationError{
			Field:   "spec.prompts",
			Message: "promptset must have at least one prompt",
		})
	}

	// Validate each prompt body
	for key, prompt := range spec.Prompts {
		bodyErrors := validateBody(prompt.Body, spec.Fragments)
		for _, err := range bodyErrors {
			if ve, ok := err.(*ValidationError); ok {
				ve.Field = fmt.Sprintf("spec.prompts[%s].%s", key, ve.Field)
			}
			errors = append(errors, err)
		}
	}

	return errors
}

func validateRulesetSpec(spec *RulesetSpec) []error {
	var errors []error

	// Validate fragment keys
	for key := range spec.Fragments {
		if !idPattern.MatchString(key) {
			errors = append(errors, &ValidationError{
				Field:   fmt.Sprintf("spec.fragments[%s]", key),
				Message: fmt.Sprintf("fragment key does not match pattern ^[a-zA-Z0-9_-]+$: %q", key),
			})
		}
	}

	// Validate collection keys
	for key := range spec.Rules {
		if !idPattern.MatchString(key) {
			errors = append(errors, &ValidationError{
				Field:   fmt.Sprintf("spec.rules[%s]", key),
				Message: fmt.Sprintf("rule key does not match pattern ^[a-zA-Z0-9_-]+$: %q", key),
			})
		}
	}

	// Validate minimum one rule
	if len(spec.Rules) == 0 {
		errors = append(errors, &ValidationError{
			Field:   "spec.rules",
			Message: "ruleset must have at least one rule",
		})
	}

	// Validate each rule body
	for key, rule := range spec.Rules {
		bodyErrors := validateBody(rule.Body, spec.Fragments)
		for _, err := range bodyErrors {
			if ve, ok := err.(*ValidationError); ok {
				ve.Field = fmt.Sprintf("spec.rules[%s].%s", key, ve.Field)
			}
			errors = append(errors, err)
		}
	}

	return errors
}
