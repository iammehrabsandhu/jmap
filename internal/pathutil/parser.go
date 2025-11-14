package pathutil

import (
	"fmt"
	"strconv"
	"strings"
)

// Segment represents a single segment in a JSON path
type Segment struct {
	Key     string
	IsArray bool
	Index   int
}

// Parse converts a dot-notation path into segments
// Examples:
//   - "user.name" -> [{Key: "user"}, {Key: "name"}]
//   - "items[0].value" -> [{Key: "items"}, {IsArray: true, Index: 0}, {Key: "value"}]
//   - "data.list[2].nested.field" -> multiple segments with array access
func Parse(path string) ([]Segment, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	var segments []Segment
	parts := strings.Split(path, ".")

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Check for array notation: field[0] or just [0]
		if idx := strings.Index(part, "["); idx >= 0 {
			// Validate closing bracket
			if !strings.HasSuffix(part, "]") {
				return nil, fmt.Errorf("invalid array notation in '%s': missing closing bracket", part)
			}

			key := part[:idx]
			indexStr := part[idx+1 : len(part)-1]

			// Parse index
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil, fmt.Errorf("invalid array index '%s' in path: %w", indexStr, err)
			}

			if index < 0 {
				return nil, fmt.Errorf("array index cannot be negative: %d", index)
			}

			// Add the key segment if it exists (not just [0])
			if key != "" {
				segments = append(segments, Segment{Key: key, IsArray: false})
			}

			// Add the array access segment
			segments = append(segments, Segment{IsArray: true, Index: index})
		} else {
			// Regular key
			segments = append(segments, Segment{Key: part, IsArray: false})
		}
	}

	if len(segments) == 0 {
		return nil, fmt.Errorf("path resulted in zero segments")
	}

	return segments, nil
}

// Join converts segments back into a path string
func Join(segments []Segment) string {
	if len(segments) == 0 {
		return ""
	}

	var result strings.Builder

	for i, seg := range segments {
		if seg.IsArray {
			// Array segment - append [index] to previous part
			result.WriteString(fmt.Sprintf("[%d]", seg.Index))
		} else {
			// Key segment
			if i > 0 && !segments[i-1].IsArray {
				result.WriteString(".")
			}
			result.WriteString(seg.Key)
		}
	}

	return result.String()
}

// ExtractArrayIndex extracts array index from path segment like "data[0]"
// Returns -1 if no array notation found
func ExtractArrayIndex(pathSegment string) int {
	start := strings.LastIndex(pathSegment, "[")
	end := strings.LastIndex(pathSegment, "]")

	if start >= 0 && end > start {
		indexStr := pathSegment[start+1 : end]
		if index, err := strconv.Atoi(indexStr); err == nil {
			return index
		}
	}

	return -1
}

// IsValid checks if a path string has valid syntax
func IsValid(path string) bool {
	if path == "" {
		return false
	}

	// Check for obviously invalid patterns
	invalid := []string{"  ", "..", ".[", "].", "[]"}
	for _, pattern := range invalid {
		if strings.Contains(path, pattern) {
			return false
		}
	}

	// Try to parse it
	_, err := Parse(path)
	return err == nil
}
