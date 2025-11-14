package transform

import (
	"fmt"
	"strings"

	"github.com/iammehrabsandhu/jmap/internal/pathutil"
)

// GetValue retrieves a value from nested JSON using dot notation path
func GetValue(data map[string]interface{}, path string) (interface{}, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	segments, err := pathutil.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path '%s': %w", path, err)
	}

	current := interface{}(data)

	for i, segment := range segments {
		if segment.IsArray {
			arr, ok := current.([]interface{})
			if !ok {
				return nil, fmt.Errorf("expected array at segment %d in path '%s', got %T",
					i, path, current)
			}

			if segment.Index < 0 || segment.Index >= len(arr) {
				return nil, fmt.Errorf("array index %d out of bounds (length %d) at segment %d in path '%s'",
					segment.Index, len(arr), i, path)
			}

			current = arr[segment.Index]
		} else {
			obj, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("expected object at segment %d in path '%s', got %T",
					i, path, current)
			}

			val, exists := obj[segment.Key]
			if !exists {
				return nil, fmt.Errorf("key '%s' not found at segment %d in path '%s'",
					segment.Key, i, path)
			}

			current = val
		}
	}

	return current, nil
}

// SetValue sets a value in nested JSON using dot notation path
// Creates intermediate objects as needed
func SetValue(data map[string]interface{}, path string, value interface{}) error {
	if path == "" {
		return fmt.Errorf("empty path")
	}

	segments, err := pathutil.Parse(path)
	if err != nil {
		return fmt.Errorf("invalid path '%s': %w", path, err)
	}

	if len(segments) == 0 {
		return fmt.Errorf("no segments in path '%s'", path)
	}

	// Navigate to the parent of the final segment
	current := data
	for i := 0; i < len(segments)-1; i++ {
		segment := segments[i]

		if segment.IsArray {
			return fmt.Errorf("arrays in middle of target path are not supported at segment %d in '%s'",
				i, path)
		}

		nextSegment := segments[i+1]

		// Ensure current segment exists
		if _, exists := current[segment.Key]; !exists {
			if nextSegment.IsArray {
				// Next level is array
				current[segment.Key] = make([]interface{}, 0)
			} else {
				// Next level is object
				current[segment.Key] = make(map[string]interface{})
			}
		}

		// Handle next level being an array
		if nextSegment.IsArray {
			arr, ok := current[segment.Key].([]interface{})
			if !ok {
				return fmt.Errorf("expected array at '%s'", segment.Key)
			}

			// Extend array if needed
			for len(arr) <= nextSegment.Index {
				arr = append(arr, make(map[string]interface{}))
			}
			current[segment.Key] = arr

			// Move to the array element
			elem, ok := arr[nextSegment.Index].(map[string]interface{})
			if !ok {
				elem = make(map[string]interface{})
				arr[nextSegment.Index] = elem
			}
			current = elem
			i++ // Skip the array segment
		} else {
			// Move to nested object
			next, ok := current[segment.Key].(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected object at '%s', got %T", segment.Key, current[segment.Key])
			}
			current = next
		}
	}

	// Set the final value
	lastSegment := segments[len(segments)-1]
	if lastSegment.IsArray {
		return fmt.Errorf("cannot set array element directly in path '%s'", path)
	}

	current[lastSegment.Key] = value
	return nil
}

// pathContext builds a human-readable context string for error messages
func pathContext(segments []pathutil.Segment, index int) string {
	if index < 0 || index >= len(segments) {
		return ""
	}

	var parts []string
	for i := 0; i <= index && i < len(segments); i++ {
		seg := segments[i]
		if seg.IsArray {
			if len(parts) > 0 {
				parts[len(parts)-1] += fmt.Sprintf("[%d]", seg.Index)
			}
		} else {
			parts = append(parts, seg.Key)
		}
	}

	return strings.Join(parts, ".")
}
