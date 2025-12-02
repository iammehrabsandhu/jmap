package spec

import (
	"fmt"
	"reflect"

	"github.com/iammehrabsandhu/jmap/internal/pathutil"
	"github.com/iammehrabsandhu/jmap/internal/spec/matcher"
	"github.com/iammehrabsandhu/jmap/types"
)

type Analyzer struct {
	matcher *matcher.FieldMatcher
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		matcher: matcher.NewFieldMatcher(false), // case-insensitive by default
	}
}

func (a *Analyzer) Analyze(input, output map[string]interface{}) (*types.TransformSpec, error) {
	inputPaths := a.flattenJSON(input, "")
	outputPaths := a.flattenJSON(output, "")

	shiftSpec := make(map[string]interface{})
	defaultSpec := make(map[string]interface{})

	for outPath, outValue := range outputPaths {
		sourcePath := a.findBestMapping(outPath, outValue, inputPaths)

		if sourcePath != "" {
			a.addToSpec(shiftSpec, sourcePath, outPath)
		} else {
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

func (a *Analyzer) addToSpec(spec map[string]interface{}, keyPath string, value interface{}) {
	segments, err := pathutil.Parse(keyPath)
	if err != nil {
		return
	}

	current := spec

	for i, seg := range segments {
		// 1. Handle key.
		if seg.Key != "" {
			// Last segment? Set value.
			if i == len(segments)-1 && !seg.IsArray {
				current[seg.Key] = value
				return
			}

			// Traverse.
			if _, exists := current[seg.Key]; !exists {
				current[seg.Key] = make(map[string]interface{})
			}

			if nextMap, ok := current[seg.Key].(map[string]interface{}); ok {
				current = nextMap
			} else {
				return // Conflict.
			}
		}

		// 2. Handle array.
		if seg.IsArray {
			// Use wildcard "*" for indices.
			wildcard := "*"

			// Last part? Set value.
			if i == len(segments)-1 {
				current[wildcard] = value
				return
			}

			// Traverse.
			if _, exists := current[wildcard]; !exists {
				current[wildcard] = make(map[string]interface{})
			}

			if nextMap, ok := current[wildcard].(map[string]interface{}); ok {
				current = nextMap
			} else {
				return // Conflict.
			}
		}
	}
}

// flattenJSON flattens JSON to dot-notation.
func (a *Analyzer) flattenJSON(data interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})

	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			newPrefix := key
			if prefix != "" {
				newPrefix = prefix + "." + key
			}

			// Recurse.
			switch nested := val.(type) {
			case map[string]interface{}:
				for nk, nv := range a.flattenJSON(nested, newPrefix) {
					result[nk] = nv
				}
			case []interface{}:
				// Store array.
				result[newPrefix] = nested
				// Flatten elements.
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

// findBestMapping matches target to input.
func (a *Analyzer) findBestMapping(targetPath string, targetValue interface{},
	inputPaths map[string]interface{}) string {

	// Extract field name.
	targetField, targetParent := pathutil.GetSchemaNames(targetPath)

	var bestPath string
	var bestScore float64

	// Find matches.
	for sourcePath, sourceValue := range inputPaths {
		sourceField, sourceParent := pathutil.GetSchemaNames(sourcePath)

		// 1. Leaf match.
		score := a.matcher.Match(sourceField, targetField)

		// 2. Leaf vs Parent (Target).
		if score < 0.9 && targetParent != "" {
			parentScore := a.matcher.Match(sourceField, targetParent)
			if parentScore > score {
				score = parentScore * 0.95 // Slight penalty.
			}
		}

		// 3. Parent (Source) vs Leaf.
		if score < 0.9 && sourceParent != "" {
			parentScore := a.matcher.Match(sourceParent, targetField)
			if parentScore > score {
				score = parentScore * 0.9
			}
		}

		// 4. Parent vs Parent.
		if score < 0.8 && sourceParent != "" && targetParent != "" {
			parentScore := a.matcher.Match(sourceParent, targetParent)
			if parentScore > score {
				score = parentScore * 0.8
			}
		}

		// 5. Exact Value Match.
		if sourceValue != nil && targetValue != nil {
			// DeepEqual for safety.
			if reflect.DeepEqual(sourceValue, targetValue) {
				// Skip simple values.
				isSimple := false
				switch v := sourceValue.(type) {
				case bool:
					isSimple = true
				case float64:
					if v == 0 || v == 1 {
						isSimple = true
					}
				}

				if !isSimple {
					if score < 0.9 {
						score = 0.9
					}
				}
			}
		}

		// Check types.
		if !matcher.TypesCompatible(sourceValue, targetValue) {
			score *= 0.1 // Penalty for mismatch.
		}

		// Track best match.
		if score > bestScore {
			bestScore = score
			bestPath = sourcePath
		}
	}

	// No good match? Return empty.
	if bestScore < 0.6 {
		return ""
	}

	return bestPath
}
