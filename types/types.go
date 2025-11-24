package types

// TransformSpec defines the complete transformation specification
// It is a list of operations to be applied in order
type TransformSpec struct {
	Operations []Operation `json:"operations"`
}

// Operation defines a single transformation step
type Operation struct {
	// Type of operation: "shift", "default", "remove", "sort", "cardinality", "modify-overwrite-beta"
	Type string `json:"type"`

	// Spec is the configuration for this operation
	// For "shift", it's a map defining the mapping
	Spec interface{} `json:"spec"`
}
