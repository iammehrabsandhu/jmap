package matcher

import (
	"reflect"
	"strings"
)

// FieldMatcher matches field names.
type FieldMatcher struct {
	caseSensitive bool
}

// NewFieldMatcher creates a matcher.
func NewFieldMatcher(caseSensitive bool) *FieldMatcher {
	return &FieldMatcher{
		caseSensitive: caseSensitive,
	}
}

// Match scores similarity (0.0 to 1.0).
func (m *FieldMatcher) Match(field1, field2 string) float64 {
	if field1 == "" || field2 == "" {
		return 0.0
	}

	original1, original2 := field1, field2

	if !m.caseSensitive {
		field1 = strings.ToLower(field1)
		field2 = strings.ToLower(field2)
	}

	// Exact match.
	if field1 == field2 {
		return 1.0
	}

	// Normalize.
	norm1 := normalize(field1)
	norm2 := normalize(field2)

	if norm1 == norm2 {
		return 0.95
	}

	// Substring check.
	if strings.Contains(norm1, norm2) || strings.Contains(norm2, norm1) {
		shorter := len(norm1)
		if len(norm2) < shorter {
			shorter = len(norm2)
		}
		longer := len(norm1)
		if len(norm2) > longer {
			longer = len(norm2)
		}

		// Base score.
		score := 0.6

		// Overlap bonus.
		ratio := float64(shorter) / float64(longer)
		score += 0.3 * ratio

		// Prefix/Suffix bonus.
		if strings.HasPrefix(norm1, norm2) || strings.HasPrefix(norm2, norm1) ||
			strings.HasSuffix(norm1, norm2) || strings.HasSuffix(norm2, norm1) {
			score += 0.1
		}

		return score
	}

	// Pattern check (camel vs snake).
	if camelToSnake(original1) == camelToSnake(original2) {
		return 0.9
	}

	// Fuzzy match (Levenshtein).
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
		return similarity * 0.8 // Scale down.
	}

	return 0.0
}

// normalize cleans field names.
func normalize(field string) string {
	field = strings.ToLower(field)
	field = strings.ReplaceAll(field, "_", "")
	field = strings.ReplaceAll(field, "-", "")
	field = strings.ReplaceAll(field, " ", "")
	return field
}

// camelToSnake converts to snake_case.
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

// levenshtein distance.
func levenshtein(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Init matrix.
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Calc distances.
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

// min of three.
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

// TypesCompatible checks type safety.
func TypesCompatible(source, target interface{}) bool {
	if source == nil || target == nil {
		return true
	}

	sourceType := reflect.TypeOf(source).Kind()
	targetType := reflect.TypeOf(target).Kind()

	// Exact match.
	if sourceType == targetType {
		return true
	}

	// Strings match anything.
	if sourceType == reflect.String || targetType == reflect.String {
		return true
	}

	// Numerics match.
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

	// Arrays/Slices match.
	if (sourceType == reflect.Array || sourceType == reflect.Slice) &&
		(targetType == reflect.Array || targetType == reflect.Slice) {
		return true
	}

	// Maps match.
	if sourceType == reflect.Map && targetType == reflect.Map {
		return true
	}

	return false
}
