package jmap_test

import (
	"encoding/json"
	"testing"

	jmap "github.com/iammehrabsandhu/jmap/pkg"
	"github.com/iammehrabsandhu/jmap/types"
)

func TestShiftTransform(t *testing.T) {
	input := `{
		"rating": {
			"primary": {
				"value": 3
			},
			"quality": {
				"value": 3
			}
		}
	}`

	specJSON := `{
		"operations": [
			{
				"type": "shift",
				"spec": {
					"rating": {
						"primary": {
							"value": "Rating"
						},
						"quality": {
							"value": "SecondaryRating"
						}
					}
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

	if output["Rating"] != 3.0 {
		t.Errorf("Expected Rating 3, got %v", output["Rating"])
	}
	if output["SecondaryRating"] != 3.0 {
		t.Errorf("Expected SecondaryRating 3, got %v", output["SecondaryRating"])
	}
}

func TestShiftWildcard(t *testing.T) {
	input := `{
		"data": {
			"key1": "val1",
			"key2": "val2"
		}
	}`

	specJSON := `{
		"operations": [
			{
				"type": "shift",
				"spec": {
					"data": {
						"*": "values.&"
					}
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

	values, ok := output["values"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected values map, got %T", output["values"])
	}

	if values["key1"] != "val1" {
		t.Errorf("Expected key1=val1, got %v", values["key1"])
	}
}

func TestSuggestSpec(t *testing.T) {
	input := `{
		"user": {
			"name": "John",
			"age": 30
		}
	}`

	output := `{
		"fullName": "John",
		"years": 30
	}`

	spec, err := jmap.SuggestSpec(input, output)
	if err != nil {
		t.Fatalf("SuggestSpec failed: %v", err)
	}

	if len(spec.Operations) == 0 {
		t.Fatal("Expected operations in suggested spec")
	}

	op := spec.Operations[0]
	if op.Type != "shift" {
		t.Errorf("Expected shift operation, got %s", op.Type)
	}

	// Verify the spec works
	result, err := jmap.Transform(input, spec)
	if err != nil {
		t.Fatalf("Transform with suggested spec failed: %v", err)
	}

	var resMap map[string]interface{}
	json.Unmarshal([]byte(result), &resMap)

	if resMap["fullName"] != "John" {
		t.Errorf("Expected fullName=John, got %v", resMap["fullName"])
	}
}
