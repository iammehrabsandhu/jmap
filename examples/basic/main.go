package main

import (
	"encoding/json"
	"fmt"
	"log"

	jmap "github.com/iammehrabsandhu/jmap/pkg"
	"github.com/iammehrabsandhu/jmap/types"
)

func main() {
	// Example 1: Generate a spec
	fmt.Println("=== Example 1: Suggest Spec ===")

	inputJSON := `{
		"user": {
			"personal": {
				"firstName": "John",
				"lastName": "Doe"
			},
			"contact": {
				"email": "john@example.com"
			}
		}
	}`

	outputJSON := `{
		"name": "John Doe",
		"email": ""
	}`

	spec, err := jmap.SuggestSpec(inputJSON, outputJSON)
	if err != nil {
		log.Fatal(err)
	}

	specJSON, _ := json.MarshalIndent(spec, "", "  ")
	fmt.Println("Generated Spec:")
	fmt.Println(string(specJSON))

	// Example 2: Use a custom spec to transform
	fmt.Println("\n=== Example 2: Transform with Custom Spec ===")

	customSpec := &types.TransformSpec{
		Version: "1.0",
		Mappings: []types.FieldMapping{
			{
				SourcePath: "user.personal.firstName",
				TargetPath: "fullName",
				Transform:  types.TransformDirect,
			},
			{
				SourcePath: "user.contact.email",
				TargetPath: "contactEmail",
				Transform:  types.TransformDirect,
			},
		},
	}

	result, err := jmap.Transform(inputJSON, customSpec)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transformed Result:")
	fmt.Println(result)

	// Example 3: Handle arrays
	fmt.Println("\n=== Example 3: Array Handling ===")

	arrayInput := `{
		"users": [
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"}
		]
	}`

	arraySpec := &types.TransformSpec{
		Version: "1.0",
		Mappings: []types.FieldMapping{
			{
				SourcePath: "users[0].name",
				TargetPath: "firstUserName",
				Transform:  types.TransformDirect,
			},
			{
				SourcePath: "users[1].name",
				TargetPath: "secondUserName",
				Transform:  types.TransformDirect,
			},
		},
	}

	arrayResult, err := jmap.Transform(arrayInput, arraySpec)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Array Transform Result:")
	fmt.Println(arrayResult)

	// Example 4: Default values
	fmt.Println("\n=== Example 4: Default Values ===")

	incompleteInput := `{
		"user": {
			"name": "John"
		}
	}`

	defaultSpec := &types.TransformSpec{
		Version: "1.0",
		Mappings: []types.FieldMapping{
			{
				SourcePath: "user.name",
				TargetPath: "name",
				Transform:  types.TransformDirect,
			},
			{
				SourcePath:   "user.email",
				TargetPath:   "email",
				Transform:    types.TransformDirect,
				DefaultValue: "no-email@example.com",
			},
			{
				SourcePath:   "",
				TargetPath:   "status",
				Transform:    types.TransformConstant,
				DefaultValue: "ACTIVE",
			},
		},
	}

	defaultResult, err := jmap.Transform(incompleteInput, defaultSpec)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result with Defaults:")
	fmt.Println(defaultResult)
}
