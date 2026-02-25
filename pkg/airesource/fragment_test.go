package airesource

import "testing"

func TestValidateInputs_ValidString(t *testing.T) {
	fragment := Fragment{
		Inputs: map[string]InputDefinition{
			"path": {Type: InputTypeString, Required: true},
		},
	}

	inputs := map[string]interface{}{
		"path": "test.txt",
	}

	validated, err := ValidateInputs("read-file", fragment, inputs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if validated["path"] != "test.txt" {
		t.Errorf("expected path=test.txt, got %v", validated["path"])
	}
}

func TestValidateInputs_ApplyDefault(t *testing.T) {
	fragment := Fragment{
		Inputs: map[string]InputDefinition{
			"count": {
				Type:     InputTypeNumber,
				Required: false,
				Default:  10,
			},
		},
	}

	inputs := map[string]interface{}{}

	validated, err := ValidateInputs("list", fragment, inputs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if validated["count"] != 10 {
		t.Errorf("expected count=10, got %v", validated["count"])
	}
}

func TestValidateInputs_TypeMismatch(t *testing.T) {
	fragment := Fragment{
		Inputs: map[string]InputDefinition{
			"count": {Type: InputTypeNumber, Required: true},
		},
	}

	inputs := map[string]interface{}{
		"count": "not a number",
	}

	_, err := ValidateInputs("list", fragment, inputs)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	inputErr, ok := err.(*InputError)
	if !ok {
		t.Fatalf("expected InputError, got %T", err)
	}

	if inputErr.Expected != "number" {
		t.Errorf("expected Expected=number, got %s", inputErr.Expected)
	}

	if inputErr.Got != "string" {
		t.Errorf("expected Got=string, got %s", inputErr.Got)
	}

	if inputErr.InputName != "count" {
		t.Errorf("expected InputName=count, got %s", inputErr.InputName)
	}
}

func TestValidateInputs_ArrayValidation(t *testing.T) {
	fragment := Fragment{
		Inputs: map[string]InputDefinition{
			"files": {
				Type: InputTypeArray,
				Items: &InputDefinition{
					Type: InputTypeString,
				},
			},
		},
	}

	inputs := map[string]interface{}{
		"files": []interface{}{"a.txt", "b.txt"},
	}

	validated, err := ValidateInputs("process", fragment, inputs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	files, ok := validated["files"].([]interface{})
	if !ok {
		t.Fatalf("expected files to be []interface{}, got %T", validated["files"])
	}

	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestValidateInputs_ObjectValidation(t *testing.T) {
	fragment := Fragment{
		Inputs: map[string]InputDefinition{
			"config": {
				Type: InputTypeObject,
				Properties: map[string]InputDefinition{
					"host": {Type: InputTypeString},
					"port": {Type: InputTypeNumber},
				},
			},
		},
	}

	inputs := map[string]interface{}{
		"config": map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		},
	}

	validated, err := ValidateInputs("connect", fragment, inputs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config, ok := validated["config"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected config to be map[string]interface{}, got %T", validated["config"])
	}

	if config["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", config["host"])
	}

	if config["port"] != 8080 {
		t.Errorf("expected port=8080, got %v", config["port"])
	}
}

func TestValidateInputs_UndefinedInput(t *testing.T) {
	fragment := Fragment{
		Inputs: map[string]InputDefinition{
			"path": {Type: InputTypeString},
		},
	}

	inputs := map[string]interface{}{
		"path":  "test.txt",
		"extra": "not defined",
	}

	_, err := ValidateInputs("read", fragment, inputs)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	inputErr, ok := err.(*InputError)
	if !ok {
		t.Fatalf("expected InputError, got %T", err)
	}

	if inputErr.InputName != "extra" {
		t.Errorf("expected InputName=extra, got %s", inputErr.InputName)
	}
}

func TestValidateInputs_MissingRequired(t *testing.T) {
	fragment := Fragment{
		Inputs: map[string]InputDefinition{
			"path": {Type: InputTypeString, Required: true},
		},
	}

	inputs := map[string]interface{}{}

	_, err := ValidateInputs("read", fragment, inputs)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	inputErr, ok := err.(*InputError)
	if !ok {
		t.Fatalf("expected InputError, got %T", err)
	}

	if inputErr.Expected != "required input" {
		t.Errorf("expected Expected='required input', got %s", inputErr.Expected)
	}

	if inputErr.Got != "missing" {
		t.Errorf("expected Got='missing', got %s", inputErr.Got)
	}
}

func TestResolveBody_SimpleString(t *testing.T) {
	body := Body{
		String: stringPtr("Simple text body"),
	}

	result, err := ResolveBody(body, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "Simple text body" {
		t.Errorf("expected 'Simple text body', got %s", result)
	}
}

func TestResolveBody_FragmentReference(t *testing.T) {
	fragments := map[string]Fragment{
		"greet": {
			Body: "Hello, {{name}}!",
		},
	}

	body := Body{
		Array: []BodyItem{
			{
				FragmentRef: &FragmentRef{
					Fragment: "greet",
					Inputs: map[string]interface{}{
						"name": "World",
					},
				},
			},
		},
	}

	result, err := ResolveBody(body, fragments)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got %s", result)
	}
}

func TestResolveBody_MixedArray(t *testing.T) {
	fragments := map[string]Fragment{
		"read": {
			Body: "Read file: {{path}}",
		},
	}

	body := Body{
		Array: []BodyItem{
			{String: stringPtr("Introduction text")},
			{
				FragmentRef: &FragmentRef{
					Fragment: "read",
					Inputs: map[string]interface{}{
						"path": "data.txt",
					},
				},
			},
			{String: stringPtr("Conclusion text")},
		},
	}

	result, err := ResolveBody(body, fragments)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "Introduction text\n\nRead file: data.txt\n\nConclusion text"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestResolveBody_MustacheConditional(t *testing.T) {
	fragments := map[string]Fragment{
		"conditional": {
			Body: "{{#show}}This is shown{{/show}}{{^show}}This is hidden{{/show}}",
		},
	}

	body := Body{
		Array: []BodyItem{
			{
				FragmentRef: &FragmentRef{
					Fragment: "conditional",
					Inputs: map[string]interface{}{
						"show": true,
					},
				},
			},
		},
	}

	result, err := ResolveBody(body, fragments)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "This is shown" {
		t.Errorf("expected 'This is shown', got %s", result)
	}
}

func TestResolveBody_MustacheArrayIteration(t *testing.T) {
	fragments := map[string]Fragment{
		"list": {
			Body: "Files:\n{{#files}}- {{.}}\n{{/files}}",
		},
	}

	body := Body{
		Array: []BodyItem{
			{
				FragmentRef: &FragmentRef{
					Fragment: "list",
					Inputs: map[string]interface{}{
						"files": []string{"a.txt", "b.txt", "c.txt"},
					},
				},
			},
		},
	}

	result, err := ResolveBody(body, fragments)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "Files:\n- a.txt\n- b.txt\n- c.txt\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestResolveBody_FragmentNotFound(t *testing.T) {
	body := Body{
		Array: []BodyItem{
			{
				FragmentRef: &FragmentRef{
					Fragment: "missing",
					Inputs:   map[string]interface{}{},
				},
			},
		},
	}

	_, err := ResolveBody(body, map[string]Fragment{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	fragErr, ok := err.(*FragmentError)
	if !ok {
		t.Fatalf("expected FragmentError, got %T", err)
	}

	if fragErr.FragmentID != "missing" {
		t.Errorf("expected FragmentID='missing', got %s", fragErr.FragmentID)
	}

	if fragErr.Message != "fragment not found" {
		t.Errorf("expected Message='fragment not found', got %s", fragErr.Message)
	}
}

func TestResolveBody_EmptyArray(t *testing.T) {
	body := Body{
		Array: []BodyItem{},
	}

	result, err := ResolveBody(body, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "" {
		t.Errorf("expected empty string, got %s", result)
	}
}

func TestResolveBody_EmptyString(t *testing.T) {
	body := Body{
		String: stringPtr(""),
	}

	result, err := ResolveBody(body, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "" {
		t.Errorf("expected empty string, got %s", result)
	}
}

func TestResolveBody_MissingVariable(t *testing.T) {
	fragments := map[string]Fragment{
		"greet": {
			Body: "Hello, {{name}}!",
		},
	}

	body := Body{
		Array: []BodyItem{
			{
				FragmentRef: &FragmentRef{
					Fragment: "greet",
					Inputs:   map[string]interface{}{},
				},
			},
		},
	}

	result, err := ResolveBody(body, fragments)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "Hello, !" {
		t.Errorf("expected 'Hello, !', got %s", result)
	}
}

func TestResolveBody_ConditionalFalse(t *testing.T) {
	fragments := map[string]Fragment{
		"conditional": {
			Body: "{{#show}}Visible{{/show}}{{^show}}Hidden{{/show}}",
		},
	}

	body := Body{
		Array: []BodyItem{
			{
				FragmentRef: &FragmentRef{
					Fragment: "conditional",
					Inputs: map[string]interface{}{
						"show": false,
					},
				},
			},
		},
	}

	result, err := ResolveBody(body, fragments)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "Hidden" {
		t.Errorf("expected 'Hidden', got %s", result)
	}
}

func TestResolveBody_EmptyArrayIteration(t *testing.T) {
	fragments := map[string]Fragment{
		"list": {
			Body: "Files:{{#files}} {{.}}{{/files}}",
		},
	}

	body := Body{
		Array: []BodyItem{
			{
				FragmentRef: &FragmentRef{
					Fragment: "list",
					Inputs: map[string]interface{}{
						"files": []string{},
					},
				},
			},
		},
	}

	result, err := ResolveBody(body, fragments)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "Files:" {
		t.Errorf("expected 'Files:', got %s", result)
	}
}

func TestResolveBody_ArrayIterationObjects(t *testing.T) {
	fragments := map[string]Fragment{
		"users": {
			Body: "{{#users}}Name: {{name}}, Age: {{age}}\n{{/users}}",
		},
	}

	body := Body{
		Array: []BodyItem{
			{
				FragmentRef: &FragmentRef{
					Fragment: "users",
					Inputs: map[string]interface{}{
						"users": []map[string]interface{}{
							{"name": "Alice", "age": 30},
							{"name": "Bob", "age": 25},
						},
					},
				},
			},
		},
	}

	result, err := ResolveBody(body, fragments)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "Name: Alice, Age: 30\nName: Bob, Age: 25\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestResolveBody_SingleItemArray(t *testing.T) {
	body := Body{
		Array: []BodyItem{
			{String: stringPtr("Single item")},
		},
	}

	result, err := ResolveBody(body, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "Single item" {
		t.Errorf("expected 'Single item', got %s", result)
	}
}

func TestResolveBody_TemplateSyntaxError(t *testing.T) {
	fragments := map[string]Fragment{
		"bad": {
			Body: "{{#unclosed}",
		},
	}

	body := Body{
		Array: []BodyItem{
			{
				FragmentRef: &FragmentRef{
					Fragment: "bad",
					Inputs:   map[string]interface{}{},
				},
			},
		},
	}

	_, err := ResolveBody(body, fragments)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	fragErr, ok := err.(*FragmentError)
	if !ok {
		t.Fatalf("expected FragmentError, got %T", err)
	}

	if fragErr.FragmentID != "bad" {
		t.Errorf("expected FragmentID='bad', got %s", fragErr.FragmentID)
	}

	if fragErr.Message != "template rendering failed" {
		t.Errorf("expected Message='template rendering failed', got %s", fragErr.Message)
	}
}
