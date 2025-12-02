package pathutil

import (
	"fmt"
	"strconv"
	"strings"
)

// Segment is a path part.
type Segment struct {
	Key     string
	IsArray bool
	Index   int
}

// Parse converts dot-notation to segments.
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

		// Check array notation.
		if idx := strings.Index(part, "["); idx >= 0 {
			// Validate bracket.
			if !strings.HasSuffix(part, "]") {
				return nil, fmt.Errorf("invalid array notation in '%s': missing closing bracket", part)
			}

			key := part[:idx]
			indexStr := part[idx+1 : len(part)-1]

			// Parse index.
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil, fmt.Errorf("invalid array index '%s' in path: %w", indexStr, err)
			}

			if index < 0 {
				return nil, fmt.Errorf("array index cannot be negative: %d", index)
			}

			// Add key if exists.
			if key != "" {
				segments = append(segments, Segment{Key: key, IsArray: false})
			}

			// Add array access.
			segments = append(segments, Segment{IsArray: true, Index: index})
		} else {
			// Regular key.
			segments = append(segments, Segment{Key: part, IsArray: false})
		}
	}

	if len(segments) == 0 {
		return nil, fmt.Errorf("path resulted in zero segments")
	}

	return segments, nil
}

// Join segments to path.
func Join(segments []Segment) string {
	if len(segments) == 0 {
		return ""
	}

	var result strings.Builder

	for i, seg := range segments {
		if seg.IsArray {
			// Array segment.
			result.WriteString(fmt.Sprintf("[%d]", seg.Index))
		} else {
			// Key segment.
			if i > 0 && !segments[i-1].IsArray {
				result.WriteString(".")
			}
			result.WriteString(seg.Key)
		}
	}

	return result.String()
}

// ExtractArrayIndex gets index from "data[0]".
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

func IsValid(path string) bool {
	if path == "" {
		return false
	}

	// Check invalid patterns.
	invalid := []string{"  ", "..", ".[", "].", "[]"}
	for _, pattern := range invalid {
		if strings.Contains(path, pattern) {
			return false
		}
	}

	_, err := Parse(path)
	return err == nil
}

// GetSchemaNames gets field and parent.
func GetSchemaNames(path string) (string, string) {
	if path == "" {
		return "", ""
	}

	segments, err := Parse(path)
	if err != nil {
		// for invalid paths.
		parts := strings.Split(path, ".")
		if len(parts) == 0 {
			return "", ""
		}

		field := cleanIndex(parts[len(parts)-1])
		parent := ""
		if len(parts) > 1 {
			parent = cleanIndex(parts[len(parts)-2])
		}
		return field, parent
	}

	// Find last non-array segment for field.
	var field string
	var fieldIdx = -1
	for i := len(segments) - 1; i >= 0; i-- {
		if !segments[i].IsArray {
			field = segments[i].Key
			fieldIdx = i
			break
		}
	}

	// If no key return empty.
	if fieldIdx == -1 {
		return "", ""
	}

	parent := ""
	if fieldIdx > 0 {
		// Find last non-array before field.
		for i := fieldIdx - 1; i >= 0; i-- {
			if !segments[i].IsArray {
				parent = segments[i].Key
				break
			}
		}
	}

	return field, parent
}

func cleanIndex(s string) string {
	if idx := strings.Index(s, "["); idx != -1 {
		return s[:idx]
	}
	return s
}
