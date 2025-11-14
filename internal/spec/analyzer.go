package spec

import (
	"fmt"
	"strings"

	"github.com/iammehrabsandhu/jmap/internal/spec/matcher"
	"github.com/iammehrabsandhu/jmap/types"
)

// Analyzer analyzes JSON structures to suggest transformation specs
type Analyzer struct {
	matcher *matcher.FieldMatcher
}

// NewAnalyzer creates a new spec analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		matcher: matcher.NewFieldMatcher(false), // case-insensitive by default
	}
}

// Analyze examines input and output JSON to suggest field mappings
func (a *Analyzer) Analyze(input, output map[string]interface{}) ([]types.FieldMapping, error) {
	// Flatten both JSONs to get all paths
	inputPaths := a.flattenJSON(input, "")
	outputPaths := a.flattenJSON(output, "")

	var mappings []types.FieldMapping

	// For each output path, find best matching input path
	for outPath, outValue := range outputPaths {
		mapping := a.findBestMapping(outPath, outValue, inputPaths)
		if mapping != nil {
			mappings = append(mappings, *mapping)
		}
	}

	return mappings, nil
}

// flattenJSON converts nested JSON to dot-notation paths with their values
func (a *Analyzer) flattenJSON(data interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})

	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			newPrefix := key
			if prefix != "" {
				newPrefix = prefix + "." + key
			}

			// Recursively flatten nested structures
			switch nested := val.(type) {
			case map[string]interface{}:
				for nk, nv := range a.flattenJSON(nested, newPrefix) {
					result[nk] = nv
				}
			case []interface{}:
				// Store array itself
				result[newPrefix] = nested
				// Also flatten array elements
				for i, item := range nested {
					indexPrefix := fmt.Sprintf("%s[%d]", newPrefix, i)
					if nestedMap, ok := item.(map[string]interface{}); ok {
						for nk, nv := range a.flattenJSON(nestedMap, indexPrefix) {
							result[nk] = nv
						}
					} else {
						result[indexPrefix] = item
					}
				}
			default:
				result[newPrefix] = val
			}
		}

	case []interface{}:
		result[prefix] = v
		for i, item := range v {
			indexPrefix := fmt.Sprintf("%s[%d]", prefix, i)
			if nestedMap, ok := item.(map[string]interface{}); ok {
				for nk, nv := range a.flattenJSON(nestedMap, indexPrefix) {
					result[nk] = nv
				}
			} else {
				result[indexPrefix] = item
			}
		}

	default:
		if prefix != "" {
			result[prefix] = v
		}
	}

	return result
}

// findBestMapping finds the best source path for a target path
func (a *Analyzer) findBestMapping(targetPath string, targetValue interface{},
	inputPaths map[string]interface{}) *types.FieldMapping {

	// Extract field name from path
	targetField := extractFieldName(targetPath)

	var bestMatch *types.FieldMapping
	var bestScore float64

	// Search for matching fields
	for sourcePath, sourceValue := range inputPaths {
		sourceField := extractFieldName(sourcePath)

		// Calculate match score
		score := a.matcher.Match(sourceField, targetField)

		// Check type compatibility
		if !matcher.TypesCompatible(sourceValue, targetValue) {
			score *= 0.5 // Penalize type mismatches
		}

		// Track best match
		if score > bestScore {
			bestScore = score
			transform := determineTransform(sourcePath, targetPath, sourceValue, targetValue)

			bestMatch = &types.FieldMapping{
				SourcePath: sourcePath,
				TargetPath: targetPath,
				Transform:  transform,
			}
		}
	}

	// If no good match found (score < 0.6), suggest constant mapping
	if bestScore < 0.6 {
		return &types.FieldMapping{
			SourcePath:   "",
			TargetPath:   targetPath,
			Transform:    types.TransformConstant,
			DefaultValue: targetValue,
		}
	}

	return bestMatch
}

// extractFieldName gets the last segment of a path
func extractFieldName(path string) string {
	// Remove array indices
	path = strings.Split(path, "[")[0]

	// Get last segment
	segments := strings.Split(path, ".")
	if len(segments) > 0 {
		return segments[len(segments)-1]
	}

	return path
}

// determineTransform determines the appropriate transform type
func determineTransform(sourcePath, targetPath string, sourceValue, targetValue interface{}) types.TransformType {
	sourceIsArray := strings.Contains(sourcePath, "[")
	targetIsArray := strings.Contains(targetPath, "[")

	// Source is array element, target is not
	if sourceIsArray && !targetIsArray {
		return types.TransformFirstElem
	}

	// Source is not array, target expects array
	if !sourceIsArray && targetIsArray {
		return types.TransformArray
	}

	// Check if source is array type but target is not
	if _, ok := sourceValue.([]interface{}); ok {
		if _, ok := targetValue.([]interface{}); !ok {
			return types.TransformFirstElem
		}
	}

	return types.TransformDirect
}
