# jmap - JSON Transformation Library

A powerful, flexible Go library for transforming JSON structures using declarative specifications. Supports nested objects, arrays, and intelligent spec generation.

## Features

- **Dynamic transformation** - No hardcoded mappings, everything driven by specs
- **Nested JSON support** - Handle deeply nested objects and arrays
- **Spec generation** - Automatically suggest transformation specs from examples
- **Multiple transform types** - Direct mapping, array handling, constants, and more
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
    "github.com/iammehrabsandhu/jmap/pkg"
)

func main() {
    input := `{"user": {"name": "John", "age": 30}}`
    
    spec := &jmap.types.TransformSpec{
        Version: "1.0",
        Mappings: []jmap.FieldMapping{
            {
                SourcePath: "user.name",
                TargetPath: "fullName",
                Transform:  jmap.TransformDirect,
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

## Transform Types

| Type | Description | Example |
|------|-------------|---------|
| `TransformDirect` | Direct 1:1 mapping | `source.field → target.field` |
| `TransformFirstElem` | Extract first array element | `array[0] → field` |
| `TransformArray` | Wrap value in array | `field → [field]` |
| `TransformConstant` | Use constant value | `→ "PENDING"` |
| `TransformObject` | Apply nested transformation | Complex nested mapping |

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
spec := &jmap.types.TransformSpec{
    Version: "1.0",
    Mappings: []jmap.FieldMapping{
        {
            SourcePath: "permissions[0].scope",
            TargetPath: "firstPermission",
            Transform:  jmap.TransformDirect,
        },
    },
}
```

### Default Values

```go
{
    SourcePath:   "optional.field",
    TargetPath:   "requiredField",
    Transform:    jmap.TransformDirect,
    DefaultValue: "fallback_value",
}
```

### Constant Fields

```go
{
    SourcePath:   "",
    TargetPath:   "createdAt",
    Transform:    jmap.TransformConstant,
    DefaultValue: "2025-01-01T00:00:00Z",
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

spec := &jmap.types.TransformSpec{
    Version: "1.0",
    Mappings: []jmap.FieldMapping{
        {
            SourcePath: "566.tables.role_v2.data.id",
            TargetPath: "id",
            Transform:  jmap.TransformDirect,
        },
        {
            SourcePath: "566.tables.role_v2.data.name",
            TargetPath: "name",
            Transform:  jmap.TransformDirect,
        },
        {
            SourcePath: "566.tables.permission_v2.data[0].scope",
            TargetPath: "permissions[0].accessType",
            Transform:  jmap.TransformDirect,
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