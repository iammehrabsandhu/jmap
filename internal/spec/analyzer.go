package spec

import (
	"fmt"
	"reflect"
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
	// We need to handle array indices like "orders[0].items[0].qty"
	// and convert them to "orders" -> "*" -> "items" -> "*" -> "qty"

	// Split by dot, but we need to handle the array parts within segments
	parts := strings.Split(keyPath, ".")
	current := spec

	for i, part := range parts {
		// Check for array index like "orders[0]"
		arrayIdx := -1
		key := part

		if idxStart := strings.Index(part, "["); idxStart != -1 {
			key = part[:idxStart]
			arrayIdx = 0 // Just a flag that it's an array
		}

		// If it's an array, we need to add the key, then the wildcard
		if arrayIdx != -1 {
			// Add the key (e.g., "orders")
			if _, exists := current[key]; !exists {
				current[key] = make(map[string]interface{})
			}

			// Move into "orders"
			if nextMap, ok := current[key].(map[string]interface{}); ok {
				current = nextMap
			} else {
				return // Conflict
			}

			// Handle the wildcard "*"
			// If this is the last part, we set the value here
			if i == len(parts)-1 {
				current["*"] = value
				return
			}

			// Otherwise, ensure "*" is a map and move into it
			if _, exists := current["*"]; !exists {
				current["*"] = make(map[string]interface{})
			}

			// Move into "*"
			if nextMap, ok := current["*"].(map[string]interface{}); ok {
				current = nextMap
			} else {
				return // Conflict
			}
		} else {
			// Not an array part (or at least not this segment)
			// If it's the last part and NOT an array, set the value
			if i == len(parts)-1 {
				current[key] = value
				return
			}

			// Otherwise create map and traverse
			if _, exists := current[key]; !exists {
				current[key] = make(map[string]interface{})
			}
			if nextMap, ok := current[key].(map[string]interface{}); ok {
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

	// Extract field name from path
	targetField := extractFieldName(targetPath)
	targetParent := extractParentFieldName(targetPath)

	var bestPath string
	var bestScore float64

	// Search for matching fields
	for sourcePath, sourceValue := range inputPaths {
		sourceField := extractFieldName(sourcePath)
		sourceParent := extractParentFieldName(sourcePath)

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

// extractFieldName gets the last segment of a path
func extractFieldName(path string) string {
	// Replace array indices [0], [1], etc. with empty string to get the "schema" path
	var cleanPath strings.Builder
	inBracket := false
	for _, char := range path {
		if char == '[' {
			inBracket = true
			continue
		}
		if char == ']' {
			inBracket = false
			continue
		}
		if !inBracket {
			cleanPath.WriteRune(char)
		}
	}

	path = cleanPath.String()

	// Get last segment
	segments := strings.Split(path, ".")
	if len(segments) > 0 {
		return segments[len(segments)-1]
	}

	return path
}

// extractParentFieldName gets the second to last segment of a path
func extractParentFieldName(path string) string {
	// Clean array indices first
	var cleanPath strings.Builder
	inBracket := false
	for _, char := range path {
		if char == '[' {
			inBracket = true
			continue
		}
		if char == ']' {
			inBracket = false
			continue
		}
		if !inBracket {
			cleanPath.WriteRune(char)
		}
	}

	path = cleanPath.String()

	segments := strings.Split(path, ".")
	if len(segments) > 1 {
		return segments[len(segments)-2]
	}

	return ""
}
