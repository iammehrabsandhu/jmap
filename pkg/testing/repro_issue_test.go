package jmap_test

import (
	"encoding/json"
	"testing"

	jmap "github.com/iammehrabsandhu/jmap/pkg"
)

func TestSuggestSpec_Transform_Mismatch(t *testing.T) {
	inputJSON := `{
		"data": {
			"items": ["a", "b"]
		}
	}`

	// We want to map input.data.items[0] to output.list[0]
	outputJSON := `{
		"list": ["a"]
	}`

	// 1. Suggest Spec
	spec, err := jmap.SuggestSpec(inputJSON, outputJSON)
	if err != nil {
		t.Fatalf("SuggestSpec failed: %v", err)
	}

	// 2. Transform using the suggested spec
	result, err := jmap.Transform(inputJSON, spec)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// 3. Verify Output
	var output map[string]interface{}
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	// Check if "list" is an array
	if _, ok := output["list"].([]interface{}); !ok {
		// If it's not an array, it might be a map with key "list[0]"
		t.Logf("Output: %s", result)
		t.Errorf("Expected 'list' to be an array, but got %T", output["list"])

		if val, exists := output["list[0]"]; exists {
			t.Errorf("Found literal key 'list[0]' with value: %v. This indicates Engine is not parsing array indices in target paths.", val)
		}
	}
}
