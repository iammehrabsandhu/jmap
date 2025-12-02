# jmap - JSON Transformation Library

Go library for transforming JSON structures using declarative specifications. Supports nested objects, arrays, and intelligent spec generation.

## Features

- **Spec generation** - Automatically suggest transformation specs from examples, rigorously tested too !
- **Dynamic transformation** - No hardcoded mappings (mostly atleast), everything driven by specs
- **Nested JSON support** - Handle deeply nested objects and arrays
- **Default values** - Fallback values for missing fields
- **CLI tool** - Command-line interface for quick transformations
- **Industry-standard structure** - Follows Go project layout best practices

## Project Structure

```
jmap/                           # Root repository
├── cmd/
│   └── main.go                 # CLI application
├── pkg/
│   ├── types.go                # Public types (types.TransformSpec, FieldMapping, etc.)
│   ├── api.go                  # Public API (Transform, SuggestSpec)
│   └── api_test.go             # API tests
├── internal/
│   ├── pathutil/
│   │   └── parser.go           # Path parsing utilities
│   ├── transform/
│   │   ├── engine.go           # Transformation engine
│   │   └── accessor.go         # JSON value get/set operations
│   └── spec/
│       ├── analyzer.go         # Spec generation logic
│       └── matcher/
│           └── matcher.go      # Field name matching algorithms
├── examples/
│   └── basic.go                # Usage examples
├── go.mod
└── README.md
```

## Installation

```bash
go get github.com/iammehrabsandhu/jmap/pkg
```

## Quick Start

### 1. As a Library

```go
package main

import (
    "fmt"
    jmap "github.com/iammehrabsandhu/jmap/pkg"
    "github.com/iammehrabsandhu/jmap/types"
)

func main() {
    input := `{"user": {"name": "John", "age": 30}}`
    
    spec := &types.TransformSpec{
        Operations: []types.Operation{
            {
                Type: "shift",
                Spec: map[string]interface{}{
                    "user": map[string]interface{}{
                        "name": "fullName",
                        "age":  "userAge",
                    },
                },
            },
        },
    }
    
    result, _ := jmap.Transform(input, spec)
    fmt.Println(result)
}
```

### 2. Generate a Spec (Suggest)

```go
inputJSON := `{
    "566": {
        "tables": {
            "role_v2": {
                "data": {
                    "id": "f2e4bf21",
                    "name": "testRole2"
                }
            }
        }
    }
}`

outputJSON := `{
    "id": "",
    "name": ""
}`

spec, _ := jmap.SuggestSpec(inputJSON, outputJSON)
// spec now contains suggested mappings
```

### 3. Using the CLI

```bash
# Generate a spec
jmap suggest -input input.json -output template.json -spec generated_spec.json

# Transform JSON
jmap transform -input data.json -spec spec.json -output result.json
```

## Operation Types

| Type | Description | Example |
|------|-------------|---------|
| `shift` | Move/map data from input to output | `"source.field": "target.field"` |
| `default` | Provide fallback values for missing fields | `{"status": "ACTIVE"}` |

## Spec Format

```json
{
  "version": "1.0",
  "mappings": [
    {
      "source_path": "566.tables.role_v2.data.name",
      "target_path": "name",
      "transform": "direct"
    },
    {
      "source_path": "566.tables.role_v2.data.id",
      "target_path": "id",
      "transform": "direct"
    },
    {
      "source_path": "",
      "target_path": "status",
      "transform": "constant",
      "default_value": "ACTIVE"
    }
  ]
}
```

## Path Syntax

- **Nested objects**: `user.profile.firstName`
- **Array access**: `items[0].value`
- **Deep nesting**: `data.level1.level2.level3.field`

## Advanced Usage

### Handling Arrays

```go
spec := &types.TransformSpec{
    Operations: []types.Operation{
        {
            Type: "shift",
            Spec: map[string]interface{}{
                "permissions": map[string]interface{}{
                    "*": map[string]interface{}{
                        "scope": "permissions[&].accessType",
                    },
                },
            },
        },
    },
}
```

### Default Values

```go
{
    Type: "default",
    Spec: map[string]interface{}{
        "requiredField": "fallback_value",
        "status":        "ACTIVE",
    },
}
```

### Constant Fields

Use the `default` operation to set constant values:

```go
{
    Type: "default",
    Spec: map[string]interface{}{
        "createdAt": "2025-01-01T00:00:00Z",
        "version":   "1.0",
    },
}
```

## Real-World Example

Transform your complex role data structure:

```go
inputJSON := `{
    "566": {
        "tables": {
            "role_v2": {
                "data": {
                    "id": "f2e4bf21-d6be-4627-b971-6d2ba2ed4858",
                    "name": "testRole2",
                    "created_by": "Sharma, Amisha (NonEmp)",
                    "lookup": {
                        "geo": ["OH"]
                    }
                }
            },
            "permission_v2": {
                "data": [
                    {
                        "scope": "14",
                        "application_id": "214800a5"
                    }
                ]
            }
        }
    }
}`

spec := &types.TransformSpec{
    Operations: []types.Operation{
        {
            Type: "shift",
            Spec: map[string]interface{}{
                "566": map[string]interface{}{
                    "tables": map[string]interface{}{
                        "role_v2": map[string]interface{}{
                            "data": map[string]interface{}{
                                "id":         "id",
                                "name":       "name",
                                "created_by": "createdBy",
                            },
                        },
                        "permission_v2": map[string]interface{}{
                            "data": map[string]interface{}{
                                "*": map[string]interface{}{
                                    "scope": "permissions[&].accessType",
                                },
                            },
                        },
                    },
                },
            },
        },
    },
}

result, _ := jmap.Transform(inputJSON, spec)
```

## Testing

```bash
go test ./pkg/...
```

## Contributing

Contributions welcome! Please ensure:
- Tests pass
- Code follows Go conventions
- Documentation is updated

## License

MIT License