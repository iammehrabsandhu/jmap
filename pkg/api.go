package jmap

import (
	"encoding/json"
	"fmt"

	"github.com/iammehrabsandhu/jmap/internal/spec"
	"github.com/iammehrabsandhu/jmap/internal/transform"
	"github.com/iammehrabsandhu/jmap/types"
)

// Transform applies the transformation spec to input JSON and returns transformed JSON
func Transform(inputJSON string, spec *types.TransformSpec) (string, error) {
	if spec == nil {
		return "", fmt.Errorf("transform spec cannot be nil")
	}

	var input interface{}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", fmt.Errorf("invalid input JSON: %w", err)
	}

	engine := transform.NewEngine()
	output, err := engine.Transform(input, spec)
	if err != nil {
		return "", fmt.Errorf("transformation failed: %w", err)
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
	return analyzer.Analyze(input, output)
}
