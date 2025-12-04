package main

import (
	"encoding/json"
	"fmt"
	"log"

	jmap "github.com/iammehrabsandhu/jmap/pkg"
	"github.com/iammehrabsandhu/jmap/types"
)

func main() {
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

	// 	outputJSON := `{
	//   "id": "",
	//   "name": "",
	//   "createdAt": "2025-03-02T14:16:49.240985632Z",
	//   "createdBy": "",
	//   "permissions": [
	//     {
	//       "accessType": "",
	//       "dimensions": [
	//         {
	//           "facility": {
	//             "value": [
	//               ""
	//             ]
	//           },
	//           "geography": {
	//             "value": [
	//               ""
	//             ]
	//           },
	//           "family_tree": {
	//             "value": [
	//               ""
	//             ]
	//           }
	//         }
	//       ],
	//       "application": ""
	//     }
	//   ]
	// }
	// `
	// spec1, err := jmap.SuggestSpec(inputJSON, outputJSON)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	spec1 := &types.TransformSpec{
		Operations: []types.Operation{
			{
				Type: "shift",
				Spec: map[string]interface{}{
					"566": map[string]interface{}{
						"tables": map[string]interface{}{
							"role_v2": map[string]interface{}{
								"data": map[string]interface{}{
									"id":         "id",
									"name":       "name",
									"created_at": "createdAt",
									"created_by": "createdBy",
									"lookup": map[string]interface{}{
										"facility": "permissions[0].dimensions[0].facility.value",
										"geo":      "permissions[0].dimensions[0].geography.value",
										"ft":       "permissions[0].dimensions[0].family_tree.value",
									},
								},
							},
							"permission_v2": map[string]interface{}{
								"data": map[string]interface{}{
									"*": map[string]interface{}{
										"scope":          "permissions[&1].accessType",
										"application_id": "permissions[&1].application",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	specJSON, _ := json.MarshalIndent(spec1, "", "  ")

	fmt.Println("Generated Spec:")

	fmt.Println(string(specJSON))

	fmt.Println("\n=== Example transformation ===")

	defaultResult, err := jmap.Transform(inputJSON, spec1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(defaultResult)
}
