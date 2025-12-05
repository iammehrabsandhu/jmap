package jmap_test

import (
	"encoding/json"
	"testing"

	jmap "github.com/iammehrabsandhu/jmap/pkg"
	"github.com/iammehrabsandhu/jmap/types"
)

func TestAncestorLookupOutOfBounds(t *testing.T) {
	input := `{
		"a": 1
	}`

	// Spec with out-of-bounds ancestor lookup (&5).
	// Since we are at root -> "a", depth is 1 (keyStack: ["a"]).
	// &0 = "a" (index 0)
	// &5 would be index -5, which is out of bounds.
	// We expect the token "&5" to be preserved in the output path.
	specJSON := `{
		"operations": [
			{
				"type": "shift",
				"spec": {
					"a": "result.&5.val"
				}
			}
		]
	}`

	var spec types.TransformSpec
	if err := json.Unmarshal([]byte(specJSON), &spec); err != nil {
		t.Fatalf("Failed to parse spec: %v", err)
	}

	result, err := jmap.Transform(input, &spec)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	var output map[string]interface{}
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	resMap, ok := output["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result map, got %T", output["result"])
	}

	// Check if "&5" key exists
	valMap, ok := resMap["&5"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected '&5' key to be preserved, got map: %v", resMap)
		return
	}

	if valMap["val"] != 1.0 {
		t.Errorf("Expected val=1, got %v", valMap["val"])
	}
}
