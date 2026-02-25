package airesource

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMaxFileSizeLimit(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "large.yml")

	content := `apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: test
  name: Test
spec:
  body: "` + strings.Repeat("x", 2000) + `"`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadResource(path, WithMaxFileSize(1024))
	if err == nil {
		t.Fatal("expected error for file size limit, got nil")
	}

	if !strings.Contains(err.Error(), "exceeds limit") {
		t.Errorf("expected 'exceeds limit' error, got: %v", err)
	}
}

func TestMaxArraySizeLimit(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "large-array.yml")

	items := make([]string, 150)
	for i := range items {
		items[i] = `  - "item"`
	}

	content := `apiVersion: ai-resource/draft
kind: Promptset
metadata:
  id: test
  name: Test
spec:
  prompts:
` + strings.Join(items, "\n")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadResource(path, WithMaxArraySize(100))
	if err == nil {
		t.Fatal("expected error for array size limit, got nil")
	}

	if !strings.Contains(err.Error(), "array size") && !strings.Contains(err.Error(), "exceeds limit") {
		t.Errorf("expected array size error, got: %v", err)
	}
}

func TestMaxNestingDepthLimit(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "deep-nesting.yml")

	content := `apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: test
  name: Test
spec:
  fragments:
    test:
      body: "test"
      inputs:
        a:
          type: object
          properties:
            b:
              type: object
              properties:
                c:
                  type: object
                  properties:
                    d:
                      type: object
                      properties:
                        e:
                          type: object
                          properties:
                            f:
                              type: string
  body: "test"`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadResource(path, WithMaxNestingDepth(5))
	if err == nil {
		t.Fatal("expected error for nesting depth limit, got nil")
	}

	if !strings.Contains(err.Error(), "nesting depth") && !strings.Contains(err.Error(), "exceeds limit") {
		t.Errorf("expected nesting depth error, got: %v", err)
	}
}

func TestTimeoutLimit(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "normal.yml")

	content := `apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: test
  name: Test
spec:
  body: "test"`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadResource(path, WithTimeout(1*time.Nanosecond))
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("expected timeout error, got: %v", err)
	}
}

func TestDefaultLimitsPass(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "normal.yml")

	content := `apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: test
  name: Test
spec:
  body: "test"`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadResource(path)
	if err != nil {
		t.Fatalf("expected success with default limits, got: %v", err)
	}
}

func TestMultiDocumentWithLimits(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "multi.yml")

	items := make([]string, 150)
	for i := range items {
		items[i] = `  - "item"`
	}

	content := `apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: test1
  name: Test1
spec:
  body: "test"
---
apiVersion: ai-resource/draft
kind: Promptset
metadata:
  id: test2
  name: Test2
spec:
  prompts:
` + strings.Join(items, "\n")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadResources(path, WithMaxArraySize(100))
	if err == nil {
		t.Fatal("expected error for array size limit in multi-doc, got nil")
	}

	if !strings.Contains(err.Error(), "document 2") {
		t.Errorf("expected document 2 error, got: %v", err)
	}
}
