package jmap

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestComplexRepro(t *testing.T) {
	inputJSON := `{
  "566": {
    "tables": {
      "role_v2": {
        "data": {
          "id": "f2e4bf21-d6be-4627-b971-6d2ba2ed4858",
          "name": "testRole2",
          "created_by": "Sharma, Amisha (NonEmp)",
          "created_at": 1740565554101839,
          "updated_at": null,
          "rank": 2004,
          "role_corpus": null,
          "status": "PENDING",
          "removed_apps": null,
          "updated_by": null,
          "lookup": {
            "app_id": [
              "214800a5-5291-4172-9c90-9b9c263a7849"
            ],
            "facility": [
              "ALL_FACILITIES"
            ],
            "ft": [
              "ALL_FT"
            ],
            "geo": [
              "OH"
            ]
          }
        }
      },
      "permission_v2": {
        "data": [
          {
            "id": "76a3d90d-6483-440c-b845-b2499f2b713d",
            "role_id": "f2e4bf21-d6be-4627-b971-6d2ba2ed4858",
            "application_id": "214800a5-5291-4172-9c90-9b9c263a7849",
            "scope": "14"
          }
        ]
      }
    }
  }
}`

	outputJSON := `{
  "id": "",
  "name": "",
  "createdAt": "2025-03-02T14:16:49.240985632Z",
  "createdBy": "",
  "permissions": [
    {
      "accessType": "",
      "dimensions": [
        {
          "facility": {
            "value": [
              ""
            ]
          },
          "geography": {
            "value": [
              ""
            ]
          },
          "family_tree": {
            "value": [
              ""
            ]
          }
        }
      ],
      "application": ""
    }
  ]
}`

	spec, err := SuggestSpec(inputJSON, outputJSON)
	fmt.Println()
	fmt.Println(spec)
	fmt.Println()
	if err != nil {
		t.Fatalf("SuggestSpec failed: %v", err)
	}

	// 1. Parse Output JSON to get all actual fields
	var output map[string]interface{}
	if err := json.Unmarshal([]byte(outputJSON), &output); err != nil {
		t.Fatalf("Invalid output JSON: %v", err)
	}
	outputFields := flatten(output, "")

	// 2. Collect Targets from Spec
	shiftTargets := make(map[string]bool)
	defaultKeys := make(map[string]bool)

	for _, op := range spec.Operations {
		switch op.Type {
		case "shift":
			collectShiftTargets(op.Spec, shiftTargets)
		case "default":
			collectDefaultKeys(op.Spec, "", defaultKeys)
		}
	}

	// 3. Count
	shiftCount := 0
	defaultCount := 0
	unaccountedCount := 0

	for _, field := range outputFields {
		schemaPath := normalizePath(field)

		if isCovered(schemaPath, shiftTargets) {
			shiftCount++
		} else if isCovered(schemaPath, defaultKeys) {
			defaultCount++
		} else {
			unaccountedCount++
		}
	}

	total := shiftCount + defaultCount + unaccountedCount
	t.Logf("Field Stats - Total: %d | Shift: %d | Default: %d | Unaccounted: %d",
		total, shiftCount, defaultCount, unaccountedCount)

	if total > 0 {
		ratio := float64(shiftCount) / float64(total)
		t.Logf("Mapping Success Rate (by field): %.1f%%", ratio*100)
	}

	// Original assertion logic (adapted)
	// We expect "permissions.dimensions.facility.value" to be shifted
	// Normalized path: "permissions.dimensions.facility.value"
	// Check if it's in shiftTargets
	targetPath := "permissions.dimensions.facility.value"
	if !isCovered(targetPath, shiftTargets) {
		t.Errorf("Expected '%s' to be in shift spec, but it was not found in targets", targetPath)
	}
}
