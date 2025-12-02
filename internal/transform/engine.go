package transform

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/iammehrabsandhu/jmap/types"
)

// Engine runs the show.
type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

// Transform applies the spec to input.
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
	if err := e.processShift(input, spec, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (e *Engine) processShift(input interface{}, spec interface{}, output map[string]interface{}) error {
	specMap, ok := spec.(map[string]interface{})
	if !ok {
		// Spec needs to be a map at the top level.
		return fmt.Errorf("invalid shift spec: expected map, got %T", spec)
	}

	inputMap, ok := input.(map[string]interface{})
	if !ok {
		// Input not a map? Can't traverse.
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
				if err := e.processField(v, k, specVal, output); err != nil {
					return err
				}
			}
			continue
		}

		// Exact match.
		if val, exists := inputMap[key]; exists {
			if err := e.processField(val, key, specVal, output); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Engine) processField(val interface{}, key string, specVal interface{}, output map[string]interface{}) error {
	switch s := specVal.(type) {
	case string:
		// Direct mapping.
		e.placeValue(output, s, val, key)
	case []interface{}:
		// Multiple mappings.
		for _, item := range s {
			if str, ok := item.(string); ok {
				e.placeValue(output, str, val, key)
			}
		}
	case map[string]interface{}:
		if nestedMap, ok := val.(map[string]interface{}); ok {
			if err := e.processShift(nestedMap, s, output); err != nil {
				return err
			}
		} else if nestedArr, ok := val.([]interface{}); ok {
			// Array input? Iterate spec for "*" or indices.
			for k, v := range s {
				if k == "*" {
					// Wildcard: all items.
					for i, item := range nestedArr {
						if err := e.processField(item, fmt.Sprintf("%d", i), v, output); err != nil {
							return err
						}
					}
				} else {
					// Specific index.
					if idx, err := strconv.Atoi(k); err == nil {
						if idx >= 0 && idx < len(nestedArr) {
							if err := e.processField(nestedArr[idx], k, v, output); err != nil {
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

func (e *Engine) placeValue(output map[string]interface{}, path string, val interface{}, key string) {
	// Handle "&" lookup.
	if strings.Contains(path, "&") {
		path = strings.ReplaceAll(path, "&", key)
	}

	// Fast path traversal.
	current := output
	start := 0
	pathLen := len(path)

	for start < pathLen {
		end := strings.IndexByte(path[start:], '.')
		var part string
		var isLast bool

		if end == -1 {
			part = path[start:]
			start = pathLen
			isLast = true
		} else {
			part = path[start : start+end]
			start += end + 1
			isLast = false
		}

		if isLast {
			// Leaf node.
			// Key exists? Make it a list.
			if existing, exists := current[part]; exists {
				if arr, ok := existing.([]interface{}); ok {
					current[part] = append(arr, val)
				} else {
					current[part] = []interface{}{existing, val}
				}
			} else {
				current[part] = val
			}
		} else {
			// Intermediate node.
			if _, exists := current[part]; !exists {
				current[part] = make(map[string]interface{})
			}
			if nextMap, ok := current[part].(map[string]interface{}); ok {
				current = nextMap
			} else {
				// Conflict? Overwrite with map.
				newMap := make(map[string]interface{})
				current[part] = newMap
				current = newMap
			}
		}
	}
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
