package spec

import (
	"fmt"
	"reflect"

	"github.com/iammehrabsandhu/jmap/internal/pathutil"
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

// Analyze examines input and output JSON to suggest a shift spec
func (a *Analyzer) Analyze(input, output map[string]interface{}) (*types.TransformSpec, error) {
	// Flatten both JSONs to get all paths
	inputPaths := a.flattenJSON(input, "")
	outputPaths := a.flattenJSON(output, "")

	shiftSpec := make(map[string]interface{})
	defaultSpec := make(map[string]interface{})

	// For each output path, find best matching input path
	for outPath, outValue := range outputPaths {
		sourcePath := a.findBestMapping(outPath, outValue, inputPaths)

		if sourcePath != "" {
			// Add to shift spec
			a.addToSpec(shiftSpec, sourcePath, outPath)
		} else {
			// Add to default spec
			a.addToSpec(defaultSpec, outPath, outValue)
		}
	}

	ops := []types.Operation{}

	if len(shiftSpec) > 0 {
		ops = append(ops, types.Operation{
			Type: "shift",
			Spec: shiftSpec,
		})
	}

	if len(defaultSpec) > 0 {
		ops = append(ops, types.Operation{
			Type: "default",
			Spec: defaultSpec,
		})
	}

	return &types.TransformSpec{
		Operations: ops,
	}, nil
}

// addToSpec adds a mapping to the nested spec map
// For shift: key=sourcePath, value=targetPath
// For default: key=targetPath, value=defaultValue
func (a *Analyzer) addToSpec(spec map[string]interface{}, keyPath string, value interface{}) {
	// Use pathutil to parse the path correctly
	segments, err := pathutil.Parse(keyPath)
	if err != nil {
		return // Should not happen with flattened paths
	}

	current := spec

	for i, seg := range segments {
		// 1. Handle the key part of the segment (if any)
		if seg.Key != "" {
			// If this is the last segment and NOT an array access, set value
			if i == len(segments)-1 && !seg.IsArray {
				current[seg.Key] = value
				return
			}

			// Otherwise create/traverse map
			if _, exists := current[seg.Key]; !exists {
				current[seg.Key] = make(map[string]interface{})
			}

			if nextMap, ok := current[seg.Key].(map[string]interface{}); ok {
				current = nextMap
			} else {
				return // Conflict
			}
		}

		// 2. Handle the array part of the segment (if any)
		if seg.IsArray {
			// In spec generation, we replace specific indices with wildcard "*"
			wildcard := "*"

			// If this is the very last part of the path (e.g. "items[0]"), set value
			if i == len(segments)-1 {
				current[wildcard] = value
				return
			}

			// Otherwise create/traverse map
			if _, exists := current[wildcard]; !exists {
				current[wildcard] = make(map[string]interface{})
			}

			if nextMap, ok := current[wildcard].(map[string]interface{}); ok {
				current = nextMap
			} else {
				return // Conflict
			}
		}
	}
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
	inputPaths map[string]interface{}) string {

	// Extract field name from path using pathutil
	targetField, targetParent := pathutil.GetSchemaNames(targetPath)

	var bestPath string
	var bestScore float64

	// Search for matching fields
	for sourcePath, sourceValue := range inputPaths {
		sourceField, sourceParent := pathutil.GetSchemaNames(sourcePath)

		// Calculate match score
		// 1. Leaf vs Leaf
		score := a.matcher.Match(sourceField, targetField)

		// 2. Leaf vs Parent (Target) - e.g. source="facility", target="...facility.value"
		// This handles cases where target structure is deeper/different
		if score < 0.9 && targetParent != "" {
			parentScore := a.matcher.Match(sourceField, targetParent)
			if parentScore > score {
				score = parentScore * 0.95 // Slight penalty but high enough to win
			}
		}

		// 3. Parent (Source) vs Leaf (Target)
		if score < 0.9 && sourceParent != "" {
			parentScore := a.matcher.Match(sourceParent, targetField)
			if parentScore > score {
				score = parentScore * 0.9
			}
		}

		// 4. Parent vs Parent
		if score < 0.8 && sourceParent != "" && targetParent != "" {
			parentScore := a.matcher.Match(sourceParent, targetParent)
			if parentScore > score {
				score = parentScore * 0.8
			}
		}

		// 5. Exact Value Match (Heuristic)
		// If values are identical (and not nil/empty), it's a strong indicator
		if sourceValue != nil && targetValue != nil {
			// Use DeepEqual to handle slices/maps safely
			if reflect.DeepEqual(sourceValue, targetValue) {
				// Avoid matching simple common values like boolean true/false or small numbers purely by value
				// unless we have no other choice.
				isSimple := false
				switch v := sourceValue.(type) {
				case bool:
					isSimple = true
				case float64:
					// If it's 0 or 1, maybe simple?
					if v == 0 || v == 1 {
						isSimple = true
					}
				}

				if !isSimple {
					// Boost score significantly
					if score < 0.9 {
						score = 0.9
					}
				}
			}
		}

		// Check type compatibility
		if !matcher.TypesCompatible(sourceValue, targetValue) {
			score *= 0.1 // Heavy penalty for type mismatch
		}

		// Track best match
		if score > bestScore {
			bestScore = score
			bestPath = sourcePath
		}
	}

	// If no good match found (score < 0.6), return empty string (will use default)
	if bestScore < 0.6 {
		return ""
	}

	return bestPath
}
