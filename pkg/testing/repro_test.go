package jmap

import (
	"encoding/json"
	"fmt"
	"testing"

	jmap "github.com/iammehrabsandhu/jmap/pkg"
)

func TestComplexRepro(t *testing.T) {
	inputJSON := `{
  "organization": {
    "users": {
      "profile": {
        "id": "user-12345",
        "name": "John Doe",
        "email": "john@example.com",
        "created_at": 1740565554101839,
        "metadata": {
          "department": ["Engineering"],
          "location": ["San Francisco"],
          "team": ["Platform"]
        }
      }
    },
    "roles": {
      "assignments": [
        {
          "id": "role-001",
          "user_id": "user-12345",
          "role_type": "admin",
          "scope": "full"
        }
      ]
    }
  }
}`

	outputJSON := `{
  "userId": "",
  "userName": "",
  "userEmail": "",
  "createdAt": "",
  "roleAssignments": [
    {
      "roleId": "",
      "roleType": "",
      "accessScope": ""
    }
  ],
  "metadata": {
    "department": [""],
    "location": [""],
    "team": [""]
  }
}`

	spec, err := jmap.SuggestSpec(inputJSON, outputJSON)
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
	// We expect "metadata.department" to be shifted
	targetPath := "metadata.department"
	if !isCovered(targetPath, shiftTargets) {
		t.Errorf("Expected '%s' to be in shift spec, but it was not found in targets", targetPath)
	}
}
