package airesource

import (
	"encoding/json"
	"fmt"
)

// Resource represents the generic envelope for all AI resources.
type Resource struct {
	APIVersion string      `yaml:"apiVersion" json:"apiVersion"`
	Kind       Kind        `yaml:"kind" json:"kind"`
	Metadata   Metadata    `yaml:"metadata" json:"metadata"`
	Spec       interface{} `yaml:"spec" json:"spec"`
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
	ID          string `yaml:"id" json:"id"`
	Name        string `yaml:"name,omitempty" json:"name,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

// AsPrompt returns the resource as a Prompt if the kind matches.
func (r *Resource) AsPrompt() (*Prompt, error) {
	if r.Kind != KindPrompt {
		return nil, fmt.Errorf("expected kind Prompt, got %s", r.Kind)
	}
	
	var spec PromptSpec
	if err := remarshal(r.Spec, &spec); err != nil {
		return nil, fmt.Errorf("invalid Prompt spec: %w", err)
	}
	
	return &Prompt{
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
		Metadata:   r.Metadata,
		Spec:       spec,
	}, nil
}

// AsPromptset returns the resource as a Promptset if the kind matches.
func (r *Resource) AsPromptset() (*Promptset, error) {
	if r.Kind != KindPromptset {
		return nil, fmt.Errorf("expected kind Promptset, got %s", r.Kind)
	}
	
	var spec PromptsetSpec
	if err := remarshal(r.Spec, &spec); err != nil {
		return nil, fmt.Errorf("invalid Promptset spec: %w", err)
	}
	
	return &Promptset{
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
		Metadata:   r.Metadata,
		Spec:       spec,
	}, nil
}

// AsRule returns the resource as a Rule if the kind matches.
func (r *Resource) AsRule() (*Rule, error) {
	if r.Kind != KindRule {
		return nil, fmt.Errorf("expected kind Rule, got %s", r.Kind)
	}
	
	var spec RuleSpec
	if err := remarshal(r.Spec, &spec); err != nil {
		return nil, fmt.Errorf("invalid Rule spec: %w", err)
	}
	
	return &Rule{
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
		Metadata:   r.Metadata,
		Spec:       spec,
	}, nil
}

// AsRuleset returns the resource as a Ruleset if the kind matches.
func (r *Resource) AsRuleset() (*Ruleset, error) {
	if r.Kind != KindRuleset {
		return nil, fmt.Errorf("expected kind Ruleset, got %s", r.Kind)
	}
	
	var spec RulesetSpec
	if err := remarshal(r.Spec, &spec); err != nil {
		return nil, fmt.Errorf("invalid Ruleset spec: %w", err)
	}
	
	return &Ruleset{
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
		Metadata:   r.Metadata,
		Spec:       spec,
	}, nil
}

func remarshal(from, to interface{}) error {
	data, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, to)
}
