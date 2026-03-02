package airesource

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation failure.
type ValidationError struct {
	Field   string
	Message string
	Cause   error
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return "validation error at " + e.Field + ": " + e.Message
	}
	return "validation error: " + e.Message
}

func (e *ValidationError) Unwrap() error {
	return e.Cause
}

// SchemaError represents a JSON Schema validation failure.
type SchemaError struct {
	Path    string
	Message string
}

func (e *SchemaError) Error() string {
	return "schema validation failed at " + e.Path + ": " + e.Message
}

// FragmentError represents a fragment resolution or rendering failure.
type FragmentError struct {
	FragmentID string
	Message    string
	Cause      error
}

func (e *FragmentError) Error() string {
	return "fragment '" + e.FragmentID + "': " + e.Message
}

func (e *FragmentError) Unwrap() error {
	return e.Cause
}

// InputError represents a fragment input validation failure.
type InputError struct {
	FragmentID string
	InputName  string
	Expected   string
	Got        string
}

func (e *InputError) Error() string {
	return "fragment '" + e.FragmentID + "' input '" + e.InputName + "': expected " + e.Expected + ", got " + e.Got
}

// LoadError represents a file loading or parsing failure.
type LoadError struct {
	Path    string
	Message string
	Cause   error
}

func (e *LoadError) Error() string {
	return "failed to load " + e.Path + ": " + e.Message
}

func (e *LoadError) Unwrap() error {
	return e.Cause
}

// MultiError represents multiple errors collected together.
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
	
	var b strings.Builder
	fmt.Fprintf(&b, "%d errors: ", len(e.Errors))
	for i, err := range e.Errors {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(err.Error())
	}
	return b.String()
}

func (e *MultiError) Unwrap() []error {
	return e.Errors
}
