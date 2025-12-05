package jmap_test

import (
	"encoding/json"
	"testing"

	jmap "github.com/iammehrabsandhu/jmap/pkg"
	"github.com/iammehrabsandhu/jmap/types"
)

func TestConcatFunction(t *testing.T) {
	input := `{
		"users": [
			{
				"first": "John",
				"last": "Doe",
				"id": 1
			},
			{
				"first": "Jane",
				"last": "Smith",
				"id": 2
			}
		]
	}`

	// We want to rekey the list by a generated key.
	// "users": { "*": "ids.@concat(first, '_', last)" }
	specJSON := `{
		"operations": [
			{
				"type": "shift",
				"spec": {
					"users": {
						"*": "ids.@concat(first, '_', last)"
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

	ids, ok := output["ids"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected ids map, got %T", output["ids"])
	}

	// Check John_Doe
	john, ok := ids["John_Doe"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected John_Doe to be a map")
	} else {
		if john["id"] != 1.0 {
			t.Errorf("Expected John_Doe.id=1, got %v", john["id"])
		}
	}

	// Check Jane_Smith
	jane, ok := ids["Jane_Smith"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected Jane_Smith to be a map")
	} else {
		if jane["id"] != 2.0 {
			t.Errorf("Expected Jane_Smith.id=2, got %v", jane["id"])
		}
	}
}

func TestLookupFunction(t *testing.T) {
	input := `{
		"items": [
			{
				"type": "A",
				"value": 10
			},
			{
				"type": "B",
				"value": 20
			}
		]
	}`

	// Map item to a key based on type.
	specJSON := `{
		"operations": [
			{
				"type": "shift",
				"spec": {
					"items": {
						"*": "values.@lookup(type, 'A', 'TypeA')"
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

	// Check TypeA
	typeA, ok := values["TypeA"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected TypeA to be a map")
	} else {
		if typeA["value"] != 10.0 {
			t.Errorf("Expected TypeA.value=10, got %v", typeA["value"])
		}
	}

	// Check B (default fallback of lookup)
	typeB, ok := values["B"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected B to be a map")
	} else {
		if typeB["value"] != 20.0 {
			t.Errorf("Expected B.value=20, got %v", typeB["value"])
		}
	}
}
