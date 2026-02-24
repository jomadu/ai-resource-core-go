package airesource

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadResource_ValidPrompt(t *testing.T) {
	yaml := `apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: summarize
  name: Summarize Text
spec:
  body: "Summarize the following text in 3-5 sentences."`

	path := createTempFile(t, "prompt.yaml", yaml)
	defer os.Remove(path)

	resource, err := LoadResource(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resource.APIVersion != "ai-resource/draft" {
		t.Errorf("expected apiVersion 'ai-resource/draft', got %s", resource.APIVersion)
	}
	if resource.Kind != KindPrompt {
		t.Errorf("expected kind Prompt, got %s", resource.Kind)
	}
	if resource.Metadata.ID != "summarize" {
		t.Errorf("expected id 'summarize', got %s", resource.Metadata.ID)
	}
}

func TestLoadResource_ValidJSON(t *testing.T) {
	json := `{
  "apiVersion": "ai-resource/draft",
  "kind": "Prompt",
  "metadata": {
    "id": "test"
  },
  "spec": {
    "body": "Test prompt"
  }
}`

	path := createTempFile(t, "prompt.json", json)
	defer os.Remove(path)

	resource, err := LoadResource(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resource.Kind != KindPrompt {
		t.Errorf("expected kind Prompt, got %s", resource.Kind)
	}
}

func TestLoadResource_FileNotFound(t *testing.T) {
	_, err := LoadResource("nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}

	if !strings.Contains(err.Error(), "file not found") {
		t.Errorf("expected 'file not found' error, got %v", err)
	}
}

func TestLoadResource_EmptyFile(t *testing.T) {
	path := createTempFile(t, "empty.yaml", "")
	defer os.Remove(path)

	_, err := LoadResource(path)
	if err == nil {
		t.Fatal("expected error for empty file")
	}

	if !strings.Contains(err.Error(), "empty file") {
		t.Errorf("expected 'empty file' error, got %v", err)
	}
}

func TestLoadResource_UnsupportedVersion(t *testing.T) {
	yaml := `apiVersion: ai-resource/v2
kind: Prompt
metadata:
  id: test
spec:
  body: "test"`

	path := createTempFile(t, "future.yaml", yaml)
	defer os.Remove(path)

	_, err := LoadResource(path)
	if err == nil {
		t.Fatal("expected error for unsupported version")
	}

	if !strings.Contains(err.Error(), "unsupported apiVersion: ai-resource/v2") {
		t.Errorf("expected unsupported version error, got %v", err)
	}
	if !strings.Contains(err.Error(), "ai-resource/draft") {
		t.Errorf("expected error to list supported versions, got %v", err)
	}
}

func TestLoadResource_MissingAPIVersion(t *testing.T) {
	yaml := `kind: Prompt
metadata:
  id: test
spec:
  body: "test"`

	path := createTempFile(t, "no-version.yaml", yaml)
	defer os.Remove(path)

	_, err := LoadResource(path)
	if err == nil {
		t.Fatal("expected error for missing apiVersion")
	}

	if !strings.Contains(err.Error(), "missing required field: apiVersion") {
		t.Errorf("expected missing apiVersion error, got %v", err)
	}
}

func TestLoadResource_MissingKind(t *testing.T) {
	yaml := `apiVersion: ai-resource/draft
metadata:
  id: test
spec:
  body: "test"`

	path := createTempFile(t, "no-kind.yaml", yaml)
	defer os.Remove(path)

	_, err := LoadResource(path)
	if err == nil {
		t.Fatal("expected error for missing kind")
	}

	if !strings.Contains(err.Error(), "missing required field: kind") {
		t.Errorf("expected missing kind error, got %v", err)
	}
}

func TestLoadResource_FileSizeLimit(t *testing.T) {
	yaml := `apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: test
spec:
  body: "test"`

	path := createTempFile(t, "large.yaml", yaml)
	defer os.Remove(path)

	_, err := LoadResource(path, WithMaxFileSize(10))
	if err == nil {
		t.Fatal("expected error for file size limit")
	}

	if !strings.Contains(err.Error(), "exceeds limit") {
		t.Errorf("expected size limit error, got %v", err)
	}
}

func TestLoadResources_MultiDocument(t *testing.T) {
	yaml := `apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: prompt1
spec:
  body: "First prompt"
---
apiVersion: ai-resource/draft
kind: Rule
metadata:
  id: rule1
spec:
  enforcement: must
  body: "First rule"`

	path := createTempFile(t, "resources.yaml", yaml)
	defer os.Remove(path)

	resources, err := LoadResources(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(resources))
	}

	if resources[0].Kind != KindPrompt {
		t.Errorf("expected first resource to be Prompt, got %s", resources[0].Kind)
	}
	if resources[1].Kind != KindRule {
		t.Errorf("expected second resource to be Rule, got %s", resources[1].Kind)
	}
}

func TestLoadResources_InvalidDocument(t *testing.T) {
	yaml := `apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: prompt1
spec:
  body: "First prompt"
---
apiVersion: ai-resource/v2
kind: Rule
metadata:
  id: rule1
spec:
  enforcement: must
  body: "Invalid version"`

	path := createTempFile(t, "invalid-multi.yaml", yaml)
	defer os.Remove(path)

	_, err := LoadResources(path)
	if err == nil {
		t.Fatal("expected error for invalid document")
	}

	if !strings.Contains(err.Error(), "document 2") {
		t.Errorf("expected error to indicate document 2, got %v", err)
	}
}

func TestLoadPrompt(t *testing.T) {
	yaml := `apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: test
spec:
  body: "Test prompt"`

	path := createTempFile(t, "prompt.yaml", yaml)
	defer os.Remove(path)

	prompt, err := LoadPrompt(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if prompt.Metadata.ID != "test" {
		t.Errorf("expected id 'test', got %s", prompt.Metadata.ID)
	}
}

func TestLoadPrompt_WrongKind(t *testing.T) {
	yaml := `apiVersion: ai-resource/draft
kind: Rule
metadata:
  id: test
spec:
  enforcement: must
  body: "Test rule"`

	path := createTempFile(t, "rule.yaml", yaml)
	defer os.Remove(path)

	_, err := LoadPrompt(path)
	if err == nil {
		t.Fatal("expected error for wrong kind")
	}

	if !strings.Contains(err.Error(), "expected kind Prompt") {
		t.Errorf("expected kind mismatch error, got %v", err)
	}
}

func createTempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	return path
}
