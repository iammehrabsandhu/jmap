package types

// TransformType defines how a value should be transformed
type TransformType string

const (
	// TransformDirect performs direct 1:1 mapping
	TransformDirect TransformType = "direct"

	// TransformArray wraps value in an array
	TransformArray TransformType = "array"

	// TransformObject maps to nested object with optional nested spec
	TransformObject TransformType = "object"

	// TransformFirstElem extracts first element from array
	TransformFirstElem TransformType = "first_elem"

	// TransformConstant uses a constant value
	TransformConstant TransformType = "constant"

	// TransformConcat concatenates multiple source values
	TransformConcat TransformType = "concat"
)

// FieldMapping defines how to map a single field from source to target
type FieldMapping struct {
	// SourcePath is the JSONPath-like source location (e.g., "566.tables.role_v2.data.name")
	SourcePath string `json:"source_path"`

	// TargetPath is the JSONPath-like target location (e.g., "name")
	TargetPath string `json:"target_path"`

	// Transform specifies the transformation type to apply
	Transform TransformType `json:"transform,omitempty"`

	// DefaultValue is used if source path doesn't exist or is null
	DefaultValue interface{} `json:"default_value,omitempty"`

	// ArrayIndex specifies a specific array index to use (optional)
	ArrayIndex *int `json:"array_index,omitempty"`

	// NestedMapping allows recursive transformation for nested objects
	NestedMapping *TransformSpec `json:"nested_mapping,omitempty"`
}

// types.TransformSpec defines the complete transformation specification
type TransformSpec struct {
	// Version of the spec format
	Version string `json:"version"`

	// Mappings is the list of field transformations to apply
	Mappings []FieldMapping `json:"mappings"`
}
