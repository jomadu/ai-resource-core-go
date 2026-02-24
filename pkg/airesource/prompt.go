package airesource

// Prompt represents a single prompt resource.
type Prompt struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       Kind        `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       PromptSpec  `yaml:"spec"`
}

// PromptSpec defines the specification for a Prompt.
type PromptSpec struct {
	Fragments map[string]Fragment `yaml:"fragments,omitempty"`
	Body      Body                `yaml:"body"`
}

// Promptset represents a collection of prompts.
type Promptset struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       Kind           `yaml:"kind"`
	Metadata   Metadata       `yaml:"metadata"`
	Spec       PromptsetSpec  `yaml:"spec"`
}

// PromptsetSpec defines the specification for a Promptset.
type PromptsetSpec struct {
	Fragments map[string]Fragment    `yaml:"fragments,omitempty"`
	Prompts   map[string]PromptItem  `yaml:"prompts"`
}

// PromptItem represents a prompt within a promptset.
type PromptItem struct {
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
	Body        Body   `yaml:"body"`
}
