package jmap

import (
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
	if err != nil {
		t.Fatalf("SuggestSpec failed: %v", err)
	}

	// We expect "facility" to be shifted, not defaulted.
	// Check if "facility" appears in the shift spec.

	// Helper to check if a key exists deeply in the map
	var findKey func(m map[string]interface{}, key string) bool
	findKey = func(m map[string]interface{}, key string) bool {
		for k, v := range m {
			if k == key {
				return true
			}
			if nested, ok := v.(map[string]interface{}); ok {
				if findKey(nested, key) {
					return true
				}
			}
		}
		return false
	}

	hasShift := false
	for _, op := range spec.Operations {
		if op.Type == "shift" {
			if specMap, ok := op.Spec.(map[string]interface{}); ok {
				// We want to see if "facility" from input is being mapped.
				// Input path has "facility".
				// The shift spec keys should reflect input structure.
				// So we look for "facility" key in the spec.
				if findKey(specMap, "facility") {
					hasShift = true
				}
			}
		}
	}

	if !hasShift {
		t.Errorf("Expected 'facility' to be in shift spec, but it likely fell back to default")
	}
}
