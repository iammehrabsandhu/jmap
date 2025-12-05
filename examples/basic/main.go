package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	jmap "github.com/iammehrabsandhu/jmap/pkg"
)

func main() {
	fmt.Println("=== JMap Example: SuggestSpec + Transform ===")

	// Input JSON: E-commerce order data
	inputJSON := `{
		"orderId": "ORD-2024-001",
		"customer": {
			"profile": {
				"firstName": "Alice",
				"lastName": "Johnson",
				"contactEmail": "alice.johnson@example.com"
			},
			"address": {
				"street": "123 Main St",
				"city": "Seattle",
				"zipCode": "98101"
			}
		},
		"items": [
			{
				"productId": "PROD-A",
				"productName": "Wireless Headphones",
				"qty": 2,
				"unitPrice": 79.99
			},
			{
				"productId": "PROD-B",
				"productName": "Phone Case",
				"qty": 1,
				"unitPrice": 19.99
			}
		],
		"payment": {
			"method": "credit_card",
			"status": "completed"
		}
	}`

	// Expected output JSON: Nested structure for order management
	outputJSON := `{
		"order": {
			"reference": "ORD-2024-001",
			"payment": {
				"type": "credit_card",
				"state": "completed"
			}
		},
		"buyer": {
			"name": "Alice",
			"email": "alice.johnson@example.com"
		},
		"shipping": {
			"location": {
				"city": "Seattle",
				"postal": "98101"
			}
		}
	}`

	fmt.Println("Step 1: Generate spec from input/output examples")
	fmt.Println("------------------------------------------------")

	spec, err := jmap.SuggestSpec(inputJSON, outputJSON)
	if err != nil {
		log.Fatal("SuggestSpec failed:", err)
	}

	specJSON, _ := json.MarshalIndent(spec, "", "  ")
	fmt.Println("Generated Spec:")
	fmt.Println(string(specJSON))

	fmt.Println("\nStep 2: Transform input using generated spec")
	fmt.Println("---------------------------------------------")

	result, err := jmap.Transform(inputJSON, spec)
	if err != nil {
		log.Fatal("Transform failed:", err)
	}

	// Pretty print result
	var resultMap map[string]interface{}
	json.Unmarshal([]byte(result), &resultMap)
	prettyResult, _ := json.MarshalIndent(resultMap, "", "  ")
	fmt.Println("Transform Result:")
	fmt.Println(string(prettyResult))

	fmt.Println("\nStep 3: Verify result matches expected output")
	fmt.Println("----------------------------------------------")

	var expected, actual map[string]interface{}
	json.Unmarshal([]byte(outputJSON), &expected)
	json.Unmarshal([]byte(result), &actual)

	if reflect.DeepEqual(expected, actual) {
		fmt.Println("✓ SUCCESS: Transform result matches expected output!")
	} else {
		fmt.Println("✗ MISMATCH: Transform result differs from expected output")
		fmt.Println("\nExpected:")
		expJSON, _ := json.MarshalIndent(expected, "", "  ")
		fmt.Println(string(expJSON))
		fmt.Println("\nActual:")
		fmt.Println(string(prettyResult))
	}
}
