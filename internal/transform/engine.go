package transform

import (
	"encoding/json"
	"fmt"

	"github.com/iammehrabsandhu/jmap/types"
)

// Engine handles JSON transformation operations
type Engine struct {
	// Could add caching or other state here in future
}

// NewEngine creates a new transformation engine
func NewEngine() *Engine {
	return &Engine{}
}

// ApplyMapping applies a single field mapping from input to output
func (e *Engine) ApplyMapping(input, output map[string]interface{}, mapping types.FieldMapping) error {
	var value interface{}
	var err error

	// Get source value
	if mapping.SourcePath == "" {
		// Use default/constant value
		if mapping.DefaultValue == nil {
			return fmt.Errorf("no source path and no default value provided for target '%s'",
				mapping.TargetPath)
		}
		value = mapping.DefaultValue
	} else {
		value, err = GetValue(input, mapping.SourcePath)
		if err != nil {
			// Use default if available
			if mapping.DefaultValue != nil {
				value = mapping.DefaultValue
			} else {
				return fmt.Errorf("source path '%s' not found: %w", mapping.SourcePath, err)
			}
		}
	}

	// Apply transformation
	value, err = e.applyTransform(value, mapping)
	if err != nil {
		return fmt.Errorf("transform failed: %w", err)
	}

	// Set target value
	if err := SetValue(output, mapping.TargetPath, value); err != nil {
		return fmt.Errorf("failed to set target '%s': %w", mapping.TargetPath, err)
	}

	return nil
}

// applyTransform applies the specified transformation to a value
func (e *Engine) applyTransform(value interface{}, mapping types.FieldMapping) (interface{}, error) {
	// Handle nil values
	if value == nil && mapping.DefaultValue != nil {
		value = mapping.DefaultValue
	}

	switch mapping.Transform {
	case types.TransformDirect:
		return value, nil

	case types.TransformFirstElem:
		arr, ok := value.([]interface{})
		if !ok {
			return nil, fmt.Errorf("expected array for first_elem transform, got %T", value)
		}
		if len(arr) == 0 {
			if mapping.DefaultValue != nil {
				return mapping.DefaultValue, nil
			}
			return nil, fmt.Errorf("cannot get first element of empty array")
		}
		return arr[0], nil

	case types.TransformArray:
		// Wrap value in array if not already an array
		if arr, ok := value.([]interface{}); ok {
			return arr, nil
		}
		return []interface{}{value}, nil

	case types.TransformConstant:
		return mapping.DefaultValue, nil

	case types.TransformObject:
		// Apply nested mapping if provided
		if mapping.NestedMapping == nil {
			return value, nil
		}

		objMap, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected object for nested transform, got %T", value)
		}

		// Convert to JSON and back through nested transform
		_, err := json.Marshal(objMap)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal for nested transform: %w", err)
		}

		nestedOutput := make(map[string]interface{})
		for _, nestedMapping := range mapping.NestedMapping.Mappings {
			if err := e.ApplyMapping(objMap, nestedOutput, nestedMapping); err != nil {
				return nil, fmt.Errorf("nested mapping failed: %w", err)
			}
		}

		return nestedOutput, nil

	case types.TransformConcat:
		// For future implementation - concatenate multiple values
		return value, nil

	case "": // No transform specified, treat as direct
		return value, nil

	default:
		return nil, fmt.Errorf("unknown transform type: %s", mapping.Transform)
	}
}
