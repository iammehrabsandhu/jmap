package jmap

import (
	"encoding/json"
	"fmt"

	"github.com/iammehrabsandhu/jmap/internal/spec"
	"github.com/iammehrabsandhu/jmap/internal/transform"
	"github.com/iammehrabsandhu/jmap/types"
)

// Transform applies the transformation spec to input JSON and returns transformed JSON
func Transform(inputJSON string, types *types.TransformSpec) (string, error) {
	if types == nil {
		return "", fmt.Errorf("transform spec cannot be nil")
	}

	var input map[string]interface{}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", fmt.Errorf("invalid input JSON: %w", err)
	}

	output := make(map[string]interface{})
	engine := transform.NewEngine()

	for i, mapping := range types.Mappings {
		if err := engine.ApplyMapping(input, output, mapping); err != nil {
			return "", fmt.Errorf("failed to apply mapping #%d (%s -> %s): %w",
				i+1, mapping.SourcePath, mapping.TargetPath, err)
		}
	}

	result, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal output: %w", err)
	}

	return string(result), nil
}

// SuggestSpec analyzes input and output JSON to suggest a transformation spec
// This is useful for generating an initial spec that can be refined manually
func SuggestSpec(inputJSON, outputJSON string) (*types.TransformSpec, error) {
	var input, output map[string]interface{}

	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return nil, fmt.Errorf("invalid input JSON: %w", err)
	}

	if err := json.Unmarshal([]byte(outputJSON), &output); err != nil {
		return nil, fmt.Errorf("invalid output JSON: %w", err)
	}

	analyzer := spec.NewAnalyzer()
	mappings, err := analyzer.Analyze(input, output)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze JSONs: %w", err)
	}

	return &types.TransformSpec{
		Version:  "1.0",
		Mappings: mappings,
	}, nil
}
