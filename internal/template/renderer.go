package template

import "github.com/cbroglie/mustache"

// Render renders a Mustache template with the provided context.
func Render(template string, context interface{}) (string, error) {
	return mustache.Render(template, context)
}
