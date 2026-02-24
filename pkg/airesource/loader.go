package airesource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/ai-resource-core-go/internal/schema"

	"gopkg.in/yaml.v3"
)

func LoadResource(path string, opts ...LoadOption) (*Resource, error) {
	options := DefaultLoadOptions()
	for _, opt := range opts {
		opt(&options)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &LoadError{Path: path, Message: "file not found", Cause: err}
	}

	if int64(len(data)) > options.MaxFileSize {
		return nil, &LoadError{
			Path:    path,
			Message: fmt.Sprintf("file size %d exceeds limit %d", len(data), options.MaxFileSize),
		}
	}

	if len(data) == 0 {
		return nil, &LoadError{Path: path, Message: "empty file"}
	}

	ext := strings.ToLower(filepath.Ext(path))
	var resource Resource

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &resource); err != nil {
			return nil, &LoadError{Path: path, Message: "invalid JSON format", Cause: err}
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &resource); err != nil {
			return nil, &LoadError{Path: path, Message: "invalid YAML format", Cause: err}
		}
	default:
		return nil, &LoadError{Path: path, Message: fmt.Sprintf("unsupported file extension: %s", ext)}
	}

	if resource.APIVersion == "" {
		return nil, &LoadError{Path: path, Message: "missing required field: apiVersion"}
	}

	if !IsSupportedVersion(resource.APIVersion) {
		return nil, &LoadError{Path: path, Message: UnsupportedVersionError(resource.APIVersion).Error()}
	}

	if resource.Kind == "" {
		return nil, &LoadError{Path: path, Message: "missing required field: kind"}
	}

	if err := schema.ValidateSchema(&resource); err != nil {
		return nil, err
	}

	return &resource, nil
}

func LoadResources(path string, opts ...LoadOption) ([]*Resource, error) {
	options := DefaultLoadOptions()
	for _, opt := range opts {
		opt(&options)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &LoadError{Path: path, Message: "file not found", Cause: err}
	}

	if int64(len(data)) > options.MaxFileSize {
		return nil, &LoadError{
			Path:    path,
			Message: fmt.Sprintf("file size %d exceeds limit %d", len(data), options.MaxFileSize),
		}
	}

	if len(data) == 0 {
		return nil, &LoadError{Path: path, Message: "empty file"}
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".json" {
		return nil, &LoadError{Path: path, Message: "multi-document loading not supported for JSON files"}
	}

	if ext != ".yaml" && ext != ".yml" {
		return nil, &LoadError{Path: path, Message: fmt.Sprintf("unsupported file extension: %s", ext)}
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var resources []*Resource

	for {
		var resource Resource
		if err := decoder.Decode(&resource); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, &LoadError{
				Path:    path,
				Message: fmt.Sprintf("invalid YAML in document %d", len(resources)+1),
				Cause:   err,
			}
		}

		if resource.APIVersion == "" {
			return nil, &LoadError{
				Path:    path,
				Message: fmt.Sprintf("document %d: missing required field: apiVersion", len(resources)+1),
			}
		}

		if !IsSupportedVersion(resource.APIVersion) {
			return nil, &LoadError{
				Path:    path,
				Message: fmt.Sprintf("document %d: %s", len(resources)+1, UnsupportedVersionError(resource.APIVersion).Error()),
			}
		}

		if resource.Kind == "" {
			return nil, &LoadError{
				Path:    path,
				Message: fmt.Sprintf("document %d: missing required field: kind", len(resources)+1),
			}
		}

		if err := schema.ValidateSchema(&resource); err != nil {
			return nil, err
		}

		resources = append(resources, &resource)
	}

	if len(resources) == 0 {
		return nil, &LoadError{Path: path, Message: "no documents found"}
	}

	return resources, nil
}

func LoadPrompt(path string, opts ...LoadOption) (*Prompt, error) {
	resource, err := LoadResource(path, opts...)
	if err != nil {
		return nil, err
	}
	return resource.AsPrompt()
}

func LoadPromptset(path string, opts ...LoadOption) (*Promptset, error) {
	resource, err := LoadResource(path, opts...)
	if err != nil {
		return nil, err
	}
	return resource.AsPromptset()
}

func LoadRule(path string, opts ...LoadOption) (*Rule, error) {
	resource, err := LoadResource(path, opts...)
	if err != nil {
		return nil, err
	}
	return resource.AsRule()
}

func LoadRuleset(path string, opts ...LoadOption) (*Ruleset, error) {
	resource, err := LoadResource(path, opts...)
	if err != nil {
		return nil, err
	}
	return resource.AsRuleset()
}
