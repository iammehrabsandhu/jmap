package matcher

import (
	"reflect"
	"strings"
)

// FieldMatcher provides intelligent field name matching
type FieldMatcher struct {
	caseSensitive bool
}

// NewFieldMatcher creates a new field matcher
func NewFieldMatcher(caseSensitive bool) *FieldMatcher {
	return &FieldMatcher{
		caseSensitive: caseSensitive,
	}
}

// Match calculates similarity between two field names
// Returns a score from 0.0 (no match) to 1.0 (perfect match)
func (m *FieldMatcher) Match(field1, field2 string) float64 {
	if field1 == "" || field2 == "" {
		return 0.0
	}

	original1, original2 := field1, field2

	if !m.caseSensitive {
		field1 = strings.ToLower(field1)
		field2 = strings.ToLower(field2)
	}

	// Exact match
	if field1 == field2 {
		return 1.0
	}

	// Normalize and compare
	norm1 := normalize(field1)
	norm2 := normalize(field2)

	if norm1 == norm2 {
		return 0.95
	}

	// Check substring matches
	if strings.Contains(norm1, norm2) || strings.Contains(norm2, norm1) {
		shorter := len(norm1)
		if len(norm2) < shorter {
			shorter = len(norm2)
		}
		longer := len(norm1)
		if len(norm2) > longer {
			longer = len(norm2)
		}
		return 0.7 * float64(shorter) / float64(longer)
	}

	// Check for common patterns (camelCase vs snake_case)
	if camelToSnake(original1) == camelToSnake(original2) {
		return 0.9
	}

	// Levenshtein distance for fuzzy matching
	distance := levenshtein(norm1, norm2)
	maxLen := len(norm1)
	if len(norm2) > maxLen {
		maxLen = len(norm2)
	}

	if maxLen == 0 {
		return 0.0
	}

	similarity := 1.0 - float64(distance)/float64(maxLen)
	if similarity > 0.6 {
		return similarity * 0.8 // Scale down fuzzy matches
	}

	return 0.0
}

// normalize removes common separators and standardizes field names
func normalize(field string) string {
	field = strings.ToLower(field)
	field = strings.ReplaceAll(field, "_", "")
	field = strings.ReplaceAll(field, "-", "")
	field = strings.ReplaceAll(field, " ", "")
	return field
}

// camelToSnake converts camelCase to snake_case for comparison
func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// levenshtein calculates the Levenshtein distance between two strings
func levenshtein(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Calculate distances
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// TypesCompatible checks if source and target types are compatible for transformation
func TypesCompatible(source, target interface{}) bool {
	if source == nil || target == nil {
		return true
	}

	sourceType := reflect.TypeOf(source).Kind()
	targetType := reflect.TypeOf(target).Kind()

	// Exact type match
	if sourceType == targetType {
		return true
	}

	// String types are compatible with anything (can be converted)
	if sourceType == reflect.String || targetType == reflect.String {
		return true
	}

	// Numeric types are compatible with each other
	numericTypes := map[reflect.Kind]bool{
		reflect.Int:     true,
		reflect.Int8:    true,
		reflect.Int16:   true,
		reflect.Int32:   true,
		reflect.Int64:   true,
		reflect.Uint:    true,
		reflect.Uint8:   true,
		reflect.Uint16:  true,
		reflect.Uint32:  true,
		reflect.Uint64:  true,
		reflect.Float32: true,
		reflect.Float64: true,
	}

	if numericTypes[sourceType] && numericTypes[targetType] {
		return true
	}

	// Arrays and slices are compatible
	if (sourceType == reflect.Array || sourceType == reflect.Slice) &&
		(targetType == reflect.Array || targetType == reflect.Slice) {
		return true
	}

	// Maps are compatible
	if sourceType == reflect.Map && targetType == reflect.Map {
		return true
	}

	return false
}
