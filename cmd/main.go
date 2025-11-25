package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"os"

	jmap "github.com/iammehrabsandhu/jmap/pkg"
	"github.com/iammehrabsandhu/jmap/types"
)

func main() {
	// Define command-line flags
	suggestCmd := flag.NewFlagSet("suggest", flag.ExitOnError)
	suggestInput := suggestCmd.String("input", "", "Input JSON file")
	suggestOutput := suggestCmd.String("output", "", "Output JSON file (template)")
	suggestSpecFile := suggestCmd.String("spec", "spec.json", "Output spec file")

	transformCmd := flag.NewFlagSet("transform", flag.ExitOnError)
	transformInput := transformCmd.String("input", "", "Input JSON file")
	transformSpec := transformCmd.String("spec", "", "Transformation spec file")
	transformOutput := transformCmd.String("output", "", "Output JSON file (optional)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "suggest":
		if err := suggestCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("Error parsing suggest flags: %v\n", err)
			os.Exit(1)
		}
		handleSuggest(*suggestInput, *suggestOutput, *suggestSpecFile)

	case "transform":
		if err := transformCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("Error parsing transform flags: %v\n", err)
			os.Exit(1)
		}
		handleTransform(*transformInput, *transformSpec, *transformOutput)

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("jmap - JSON Transformation Tool")
	fmt.Println("\nUsage:")
	fmt.Println("  jmap suggest -input <input.json> -output <output_template.json> [-spec <spec.json>]")
	fmt.Println("  jmap transform -input <input.json> -spec <spec.json> [-output <output.json>]")
	fmt.Println("\nCommands:")
	fmt.Println("  suggest    Generate a transformation spec by analyzing input and output JSONs")
	fmt.Println("  transform  Transform JSON using a specification")
}

func handleSuggest(inputFile, outputFile, specFile string) {
	if inputFile == "" || outputFile == "" {
		fmt.Println("Error: -input and -output flags are required")
		os.Exit(1)
	}

	// Read input JSON
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	// Read output template JSON
	outputData, err := os.ReadFile(outputFile)
	if err != nil {
		fmt.Printf("Error reading output file: %v\n", err)
		os.Exit(1)
	}

	// Generate spec
	spec, err := jmap.SuggestSpec(string(inputData), string(outputData))
	if err != nil {
		fmt.Printf("Error generating spec: %v\n", err)
		os.Exit(1)
	}

	// Write spec to file
	specJSON, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling spec: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(specFile, specJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing spec file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Spec generated successfully: %s\n", specFile)
	fmt.Println("\nGenerated Spec:")
	fmt.Println(string(specJSON))
}

func handleTransform(inputFile, specFile, outputFile string) {
	if inputFile == "" || specFile == "" {
		fmt.Println("Error: -input and -spec flags are required")
		os.Exit(1)
	}

	// Read input JSON
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	// Read spec
	specData, err := os.ReadFile(specFile)
	if err != nil {
		fmt.Printf("Error reading spec file: %v\n", err)
		os.Exit(1)
	}

	var spec types.TransformSpec
	err = json.Unmarshal(specData, &spec)
	if err != nil {
		fmt.Printf("Error parsing spec: %v\n", err)
		os.Exit(1)
	}

	// Transform
	result, err := jmap.Transform(string(inputData), &spec)
	if err != nil {
		fmt.Printf("Error transforming JSON: %v\n", err)
		os.Exit(1)
	}

	// Output result
	if outputFile != "" {
		err = os.WriteFile(outputFile, []byte(result), 0644)
		if err != nil {
			fmt.Printf("Error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Transformation complete: %s\n", outputFile)
	} else {
		fmt.Println(result)
	}
}
