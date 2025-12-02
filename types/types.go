package types

// TransformSpec is a list of ops.
type TransformSpec struct {
	Operations []Operation `json:"operations"`
}

// Operation is one step.
type Operation struct {
	// Type: "shift", "default"
	Type string `json:"type"`

	// Spec config.
	Spec interface{} `json:"spec"`
}
