package main

import (
	"encoding/json"
	"fmt"
	"log"

	jmap "github.com/iammehrabsandhu/jmap/pkg"
)

func main() {
	// Example 1: Generate a spec
	fmt.Println("=== Example 1: Suggest Spec ===")

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
}
`

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
}
`

	// spec := &types.TransformSpec{
	// 	Operations: []types.Operation{
	// 		{
	// 			// ---- SHIFT -------------------------------------------------
	// 			Type: "shift",
	// 			Spec: map[string]interface{}{
	// 				// ── orders ───────────────────────────────────────
	// 				"orders": map[string]interface{}{
	// 					// “*” iterates over every element of the orders array.
	// 					"*": map[string]interface{}{
	// 						// orderId ← orders[i].id
	// 						"id": "orderId",

	// 						// ── customer ─────────────────────────────
	// 						"customer": map[string]interface{}{
	// 							// We expose the two name parts separately – they can be
	// 							// concatenated later in Go (jmap has no built‑in concat).
	// 							"firstName": "firstName",
	// 							"lastName":  "lastName",

	// 							// city ← orders[i].customer.address.city
	// 							"address": map[string]interface{}{
	// 								"city": "city",
	// 							},
	// 						},

	// 						// ── items (used only to collect qty/price) ─────
	// 						"items": map[string]interface{}{
	// 							"*": map[string]interface{}{
	// 								// Store each qty / price in a temporary array.
	// 								// These helpers will be removed (or used for a post‑step
	// 								// calculation) after the shift.
	// 								"qty":   "tmpQty[]",
	// 								"price": "tmpPrice[]",
	// 							},
	// 						},

	// 						// status ← orders[i].status
	// 						"status": "status",
	// 					},
	// 				},

	// 				// ── metadata ───────────────────────────────────────
	// 				"metadata": map[string]interface{}{
	// 					// generatedAt is copied unchanged.
	// 					"generatedAt": "generatedAt",
	// 				},
	// 			},
	// 		},

	// 		// -----------------------------------------------------------------
	// 		// NOTE: JMap currently does **not** have an expression engine, so the
	// 		//       `total` field (sum of qty*price) and the combined
	// 		//       (firstName + " " + lastName) must be derived in Go after the
	// 		//       shift step.  No `default` operation is required because every
	// 		//       target field is produced by the shift.
	// 		// -----------------------------------------------------------------
	// 	},
	// }

	spec1, err := jmap.SuggestSpec(inputJSON, outputJSON)
	if err != nil {
		log.Fatal(err)
	}

	specJSON, _ := json.MarshalIndent(spec1, "", "  ")
	fmt.Println("Generated Spec:")
	fmt.Println(string(specJSON))

	// Example 4: Default values
	fmt.Println("\n=== Example 4: Default Values ===")

	// incompleteInput := `{
	// 	"user": {
	// 		"name": "John"
	// 	}
	// }`

	// defaultSpec := &types.TransformSpec{
	// 	Version: "1.0",
	// 	Mappings: []types.FieldMapping{
	// 		{
	// 			SourcePath: "user.name",
	// 			TargetPath: "name",
	// 			Transform:  types.TransformDirect,
	// 		},
	// 		{
	// 			SourcePath:   "user.email",
	// 			TargetPath:   "email",
	// 			Transform:    types.TransformDirect,
	// 			DefaultValue: "no-email@example.com",
	// 		},
	// 		{
	// 			SourcePath:   "",
	// 			TargetPath:   "status",
	// 			Transform:    types.TransformConstant,
	// 			DefaultValue: "ACTIVE",
	// 		},
	// 	},
	// }

	defaultResult, err := jmap.Transform(inputJSON, spec1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result with Defaults:")
	fmt.Println(defaultResult)
}
