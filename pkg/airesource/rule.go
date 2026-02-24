package airesource

// Rule represents a single rule resource.
type Rule struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       Kind     `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       RuleSpec `yaml:"spec"`
}

// RuleSpec defines the specification for a Rule.
type RuleSpec struct {
	Fragments   map[string]Fragment `yaml:"fragments,omitempty" json:"fragments,omitempty"`
	Enforcement Enforcement         `yaml:"enforcement,omitempty" json:"enforcement,omitempty"`
	Scope       []ScopeEntry        `yaml:"scope,omitempty" json:"scope,omitempty"`
	Body        Body                `yaml:"body" json:"body"`
}

// Enforcement represents the enforcement level of a rule.
type Enforcement string

const (
	EnforcementMay    Enforcement = "may"
	EnforcementShould Enforcement = "should"
	EnforcementMust   Enforcement = "must"
)

// ScopeEntry defines file patterns for rule scope.
type ScopeEntry struct {
	Files []string `yaml:"files"`
}

// Ruleset represents a collection of rules.
type Ruleset struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       Kind        `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       RulesetSpec `yaml:"spec"`
}

// RulesetSpec defines the specification for a Ruleset.
type RulesetSpec struct {
	Fragments map[string]Fragment  `yaml:"fragments,omitempty" json:"fragments,omitempty"`
	Rules     map[string]RuleItem  `yaml:"rules" json:"rules"`
}

// RuleItem represents a rule within a ruleset.
type RuleItem struct {
	Name        string      `yaml:"name,omitempty" json:"name,omitempty"`
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
	Priority    int         `yaml:"priority,omitempty" json:"priority,omitempty"`
	Enforcement Enforcement `yaml:"enforcement,omitempty" json:"enforcement,omitempty"`
	Scope       []ScopeEntry `yaml:"scope,omitempty" json:"scope,omitempty"`
	Body        Body        `yaml:"body" json:"body"`
}
