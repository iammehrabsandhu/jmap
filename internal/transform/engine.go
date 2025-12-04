package transform

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/iammehrabsandhu/jmap/types"
)

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Transform(input interface{}, spec *types.TransformSpec) (interface{}, error) {
	current := input
	var err error

	for _, op := range spec.Operations {
		switch op.Type {
		case "shift":
			current, err = e.applyShift(current, op.Spec)
		case "default":
			current, err = e.applyDefault(current, op.Spec)
		default:
			return nil, fmt.Errorf("unknown operation type: %s", op.Type)
		}

		if err != nil {
			return nil, fmt.Errorf("operation %s failed: %w", op.Type, err)
		}
	}

	return current, nil
}

func (e *Engine) applyShift(input interface{}, spec interface{}) (interface{}, error) {
	output := make(map[string]interface{})
	if err := e.processShift(input, spec, output, []string{}); err != nil {
		return nil, err
	}
	return output, nil
}

func (e *Engine) processShift(input interface{}, spec interface{}, output map[string]interface{}, keyStack []string) error {
	specMap, ok := spec.(map[string]interface{})
	if !ok {
		// Spec needs to be a map at the top level.
		return fmt.Errorf("invalid shift spec: expected map, got %T", spec)
	}

	inputMap, ok := input.(map[string]interface{})
	if !ok {
		// Input not a map, no traverse.
		return nil
	}

	// Sort keys for deterministic output.
	keys := make([]string, 0, len(specMap))
	for k := range specMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		specVal := specMap[key]

		if key == "*" {
			// Sort input keys too.
			inputKeys := make([]string, 0, len(inputMap))
			for k := range inputMap {
				inputKeys = append(inputKeys, k)
			}
			sort.Strings(inputKeys)

			for _, k := range inputKeys {
				v := inputMap[k]
				newStack := append(keyStack, k)
				if err := e.processField(v, k, specVal, output, newStack); err != nil {
					return err
				}
			}
			continue
		}

		// Exact match.
		if val, exists := inputMap[key]; exists {
			newStack := append(keyStack, key)
			if err := e.processField(val, key, specVal, output, newStack); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Engine) processField(val interface{}, key string, specVal interface{}, output map[string]interface{}, keyStack []string) error {
	switch s := specVal.(type) {
	case string:
		// Direct mapping.
		e.placeValue(output, s, val, keyStack)
	case []interface{}:
		// Multiple mappings.
		for _, item := range s {
			if str, ok := item.(string); ok {
				e.placeValue(output, str, val, keyStack)
			}
		}
	case map[string]interface{}:
		if nestedMap, ok := val.(map[string]interface{}); ok {
			if err := e.processShift(nestedMap, s, output, keyStack); err != nil {
				return err
			}
		} else if nestedArr, ok := val.([]interface{}); ok {
			// for array input iterate spec for "*" or indices.
			for k, v := range s {
				if k == "*" {
					// Wildcard: all items.
					for i, item := range nestedArr {
						newStack := append(keyStack, fmt.Sprintf("%d", i))
						if err := e.processField(item, fmt.Sprintf("%d", i), v, output, newStack); err != nil {
							return err
						}
					}
				} else {
					// Specific index.
					if idx, err := strconv.Atoi(k); err == nil {
						if idx >= 0 && idx < len(nestedArr) {
							newStack := append(keyStack, k)
							if err := e.processField(nestedArr[idx], k, v, output, newStack); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (e *Engine) placeValue(output map[string]interface{}, path string, val interface{}, keyStack []string) {
	// Handle "&" lookup.
	// & = &0 = current key (last in stack)
	// &1 = parent key (second to last)
	// etc.
	if strings.Contains(path, "&") {
		// We need to replace all occurrences of &N or &
		// Since we don't want to use regex if we can avoid it (to keep imports clean, though strings is imported),
		// let's iterate.
		// Actually, simple iteration is fine.

		var newPath strings.Builder
		n := len(path)
		for i := 0; i < n; i++ {
			if path[i] == '&' {
				// Check for number
				j := i + 1
				numStart := j
				for j < n && path[j] >= '0' && path[j] <= '9' {
					j++
				}

				level := 0
				if j > numStart {
					// We have a number
					val, err := strconv.Atoi(path[numStart:j])
					if err == nil {
						level = val
					}
					i = j - 1 // Advance i
				}

				// Get key from stack
				// stack: [root, ..., parent, current]
				// level 0 = current = len-1
				// level 1 = parent = len-2
				stackIdx := len(keyStack) - 1 - level
				if stackIdx >= 0 && stackIdx < len(keyStack) {
					newPath.WriteString(keyStack[stackIdx])
				} else {
					// Out of bounds, keep original or empty?
					// Jolt keeps empty or original? Let's write nothing or keep &?
					// Let's write the key if found, else nothing (empty string replacement).
				}
			} else {
				newPath.WriteByte(path[i])
			}
		}
		path = newPath.String()
	}

	// Fast path traversal.
	// We need to handle array notation: "permissions[0]"

	var current interface{} = output

	// Helper to set value in map or array
	// We need to parse the path into segments.
	// "a.b[0].c" -> "a", "b", "[0]", "c"

	segments := parsePath(path)

	for i, seg := range segments {
		isLast := i == len(segments)-1

		if isLast {
			e.setValue(current, seg, val)
		} else {
			current = e.ensureContainer(current, seg, segments[i+1])
			if current == nil {
				return
			}
		}
	}
}

func parsePath(path string) []string {
	var segments []string
	var sb strings.Builder

	for i := 0; i < len(path); i++ {
		c := path[i]
		if c == '.' {
			if sb.Len() > 0 {
				segments = append(segments, sb.String())
				sb.Reset()
			}
		} else if c == '[' {
			if sb.Len() > 0 {
				segments = append(segments, sb.String())
				sb.Reset()
			}
			// Read until ]
			start := i
			for i < len(path) && path[i] != ']' {
				i++
			}
			if i < len(path) {
				// Include brackets? No, let's just keep [N] as the segment identifier
				segments = append(segments, path[start:i+1])
			}
		} else {
			sb.WriteByte(c)
		}
	}
	if sb.Len() > 0 {
		segments = append(segments, sb.String())
	}
	return segments
}

func (e *Engine) setValue(container interface{}, key string, val interface{}) {
	if strings.HasPrefix(key, "[") && strings.HasSuffix(key, "]") {
		// Array index
		idxStr := key[1 : len(key)-1]
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			return // Ignore invalid index
		}

		if arr, ok := container.(*[]interface{}); ok {
			// Handle negative index
			if idx < 0 {
				idx = len(*arr) + idx
				if idx < 0 {
					return // Still negative, invalid
				}
			}

			// Grow if needed
			if idx >= len(*arr) {
				newArr := make([]interface{}, idx+1)
				copy(newArr, *arr)
				*arr = newArr
			}
			(*arr)[idx] = val
		}
	} else {
		// Map key
		if m, ok := container.(map[string]interface{}); ok {
			m[key] = val
		}
	}
}

func (e *Engine) ensureContainer(container interface{}, key string, nextKey string) interface{} {
	isNextArray := strings.HasPrefix(nextKey, "[")

	if strings.HasPrefix(key, "[") && strings.HasSuffix(key, "]") {
		// Current is array index
		idxStr := key[1 : len(key)-1]
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			return nil
		}

		if arr, ok := container.(*[]interface{}); ok {
			if idx >= len(*arr) {
				newArr := make([]interface{}, idx+1)
				copy(newArr, *arr)
				*arr = newArr
			}

			if (*arr)[idx] == nil {
				if isNextArray {
					newSlice := make([]interface{}, 0)
					(*arr)[idx] = &newSlice
				} else {
					(*arr)[idx] = make(map[string]interface{})
				}
			}
			return (*arr)[idx]
		}
	} else {
		// Current is map key
		if m, ok := container.(map[string]interface{}); ok {
			if m[key] == nil {
				if isNextArray {
					newSlice := make([]interface{}, 0)
					m[key] = &newSlice
				} else {
					m[key] = make(map[string]interface{})
				}
			}
			return m[key]
		}
	}
	return nil
}

// applyDefault fills missing fields.
func (e *Engine) applyDefault(input interface{}, spec interface{}) (interface{}, error) {
	specMap, ok := spec.(map[string]interface{})
	if !ok {
		return input, nil
	}

	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return input, nil
	}

	// Sort keys.
	keys := make([]string, 0, len(specMap))
	for k := range specMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		defaultVal := specMap[key]
		if _, exists := inputMap[key]; !exists {
			inputMap[key] = defaultVal
		} else if nestedSpec, ok := defaultVal.(map[string]interface{}); ok {
			// Recurse.
			if nestedInput, ok := inputMap[key].(map[string]interface{}); ok {
				res, _ := e.applyDefault(nestedInput, nestedSpec)
				inputMap[key] = res
			}
		}
	}

	return inputMap, nil
}

// Regex patterns for function parsing.
var (
	concatPattern = regexp.MustCompile(`^@concat\((.*)\)$`)
	lookupPattern = regexp.MustCompile(`^@lookup\(([^,]+),\s*'([^']*)',\s*'([^']*)'\)$`)
)

// evaluateFunction checks if a spec value is a function and evaluates it.
// Returns the original value if not a function.
func (e *Engine) evaluateFunction(spec string, input interface{}, keyStack []string) interface{} {
	// Check for @concat
	if match := concatPattern.FindStringSubmatch(spec); match != nil {
		return e.evaluateConcat(match[1], input, keyStack)
	}

	// Check for @lookup
	if match := lookupPattern.FindStringSubmatch(spec); match != nil {
		return e.evaluateLookup(match[1], match[2], match[3], input)
	}

	return nil // Not a function
}

// evaluateConcat joins multiple fields with separators.
// Args format: field1, ' ', field2, '-', field3
func (e *Engine) evaluateConcat(args string, input interface{}, keyStack []string) interface{} {
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return nil
	}

	var result strings.Builder
	parts := strings.Split(args, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		// Check if it's a literal string (quoted)
		if len(part) >= 2 && part[0] == '\'' && part[len(part)-1] == '\'' {
			result.WriteString(part[1 : len(part)-1])
		} else {
			// It's a field reference - look it up in input
			val := e.getNestedValue(inputMap, part)
			if val != nil {
				switch v := val.(type) {
				case string:
					result.WriteString(v)
				case []interface{}:
					// Join array elements with space
					for i, item := range v {
						if i > 0 {
							result.WriteString(" ")
						}
						if s, ok := item.(string); ok {
							result.WriteString(s)
						}
					}
				default:
					result.WriteString(fmt.Sprintf("%v", v))
				}
			}
		}
	}

	return result.String()
}

// evaluateLookup performs a key-value lookup.
// @lookup(field, 'key', 'value') -> if input[field] == 'key', return 'value'
func (e *Engine) evaluateLookup(field, key, value string, input interface{}) interface{} {
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return nil
	}

	field = strings.TrimSpace(field)
	fieldVal := e.getNestedValue(inputMap, field)
	if fieldVal == nil {
		return nil
	}

	// Check if field value matches the key
	fieldStr := fmt.Sprintf("%v", fieldVal)
	if fieldStr == key {
		return value
	}

	return fieldVal // Return original value if no match
}

// getNestedValue retrieves a value from a nested map using dot notation.
func (e *Engine) getNestedValue(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			current = m[part]
		} else if arr, ok := current.([]interface{}); ok {
			// Try to parse as array index
			if idx, err := strconv.Atoi(part); err == nil && idx >= 0 && idx < len(arr) {
				current = arr[idx]
			} else {
				return nil
			}
		} else {
			return nil
		}
	}

	return current
}
