package airesource

import "fmt"

// Resource represents the generic envelope for all AI resources.
type Resource struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       Kind        `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       interface{} `yaml:"spec"`
}

// Kind represents the type of AI resource.
type Kind string

const (
	KindPrompt    Kind = "Prompt"
	KindPromptset Kind = "Promptset"
	KindRule      Kind = "Rule"
	KindRuleset   Kind = "Ruleset"
)

// Metadata contains identifying information for a resource.
type Metadata struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
}

// AsPrompt returns the resource as a Prompt if the kind matches.
func (r *Resource) AsPrompt() (*Prompt, error) {
	if r.Kind != KindPrompt {
		return nil, fmt.Errorf("expected kind Prompt, got %s", r.Kind)
	}
	return &Prompt{
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
		Metadata:   r.Metadata,
		Spec:       r.Spec.(PromptSpec),
	}, nil
}

// AsPromptset returns the resource as a Promptset if the kind matches.
func (r *Resource) AsPromptset() (*Promptset, error) {
	if r.Kind != KindPromptset {
		return nil, fmt.Errorf("expected kind Promptset, got %s", r.Kind)
	}
	return &Promptset{
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
		Metadata:   r.Metadata,
		Spec:       r.Spec.(PromptsetSpec),
	}, nil
}

// AsRule returns the resource as a Rule if the kind matches.
func (r *Resource) AsRule() (*Rule, error) {
	if r.Kind != KindRule {
		return nil, fmt.Errorf("expected kind Rule, got %s", r.Kind)
	}
	return &Rule{
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
		Metadata:   r.Metadata,
		Spec:       r.Spec.(RuleSpec),
	}, nil
}

// AsRuleset returns the resource as a Ruleset if the kind matches.
func (r *Resource) AsRuleset() (*Ruleset, error) {
	if r.Kind != KindRuleset {
		return nil, fmt.Errorf("expected kind Ruleset, got %s", r.Kind)
	}
	return &Ruleset{
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
		Metadata:   r.Metadata,
		Spec:       r.Spec.(RulesetSpec),
	}, nil
}
