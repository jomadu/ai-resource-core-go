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
