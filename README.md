# jmap - JSON Transformation Library

Go library for transforming JSON structures using declarative specifications. Supports nested objects, arrays, and intelligent spec generation.

## Features

- **Spec generation** - Automatically suggest transformation specs from examples, rigorously tested too !
- **Dynamic transformation** - No hardcoded mappings (mostly atleast), everything driven by specs
- **Nested JSON support** - Handle deeply nested objects and arrays
- **Default values** - Fallback values for missing fields
- **CLI tool** - Command-line interface for quick transformations

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
		"rating": {
			"primary": {
				"value": 3
			},
			"quality": {
				"value": 3
			}
		}
	}`

outputJSON := `{
		"operations": [
			{
				"type": "shift",
				"spec": {
					"rating": {
						"primary": {
							"value": "Rating"
						},
						"quality": {
							"value": "SecondaryRating"
						}
					}
				}
			}
		]
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
      "source_path": "organization.users.profile.name",
      "target_path": "userName",
      "transform": "direct"
    },
    {
      "source_path": "organization.users.profile.id",
      "target_path": "userId",
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

### Handling Arrays and Dynamic Keys

#### Wildcard Matching

Use `*` to iterate over all array elements or object keys:

```go
spec := &types.TransformSpec{
    Operations: []types.Operation{
        {
            Type: "shift",
            Spec: map[string]interface{}{
                "users": map[string]interface{}{
                    "*": "people.&.profile",  // & references the current key/index
                },
            },
        },
    },
}
```

#### Ancestor Lookup (`&`)

Reference keys from parent levels using `&N`:
- `&` or `&0`: Current key/index
- `&1`: Parent key
- `&2`: Grandparent key

Example:
```json
{
  "operations": [{
    "type": "shift",
    "spec": {
      "departments": {
        "*": {
          "employees": {
            "*": "result.&2.staff.&1"
          }
        }
      }
    }
  }]
}
```

#### Concatenation Function (`@concat`)

Build dynamic keys by concatenating field values:

```json
{
  "operations": [{
    "type": "shift",
    "spec": {
      "users": {
        "*": "byName.@concat(first, '_', last)"
      }
    }
  }]
}
```

Input: `{"users": [{"first": "John", "last": "Doe", "id": 1}]}`  
Output: `{"byName": {"John_Doe": {"first": "John", "last": "Doe", "id": 1}}}`

#### Lookup Function (`@lookup`)

Map values conditionally based on field values:

```json
{
  "operations": [{
    "type": "shift",
    "spec": {
      "items": {
        "*": "categorized.@lookup(type, 'premium', 'PremiumItems')"
      }
    }
  }]
}
```

If `type == 'premium'`, item is placed under `categorized.PremiumItems`.  
Otherwise, item is placed under `categorized.{original_type_value}`.

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

Transform a complex organization data structure:

```go
inputJSON := `{
    "organization": {
        "users": {
            "profile": {
                "id": "user-12345",
                "name": "John Doe",
                "email": "john@example.com",
                "metadata": {
                    "department": ["Engineering"]
                }
            }
        },
        "roles": {
            "assignments": [
                {
                    "scope": "full",
                    "role_type": "admin"
                }
            ]
        }
    }
}`

spec := &types.TransformSpec{
    Operations: []types.Operation{
        {
            Type: "shift",
            Spec: map[string]interface{}{
                "organization": map[string]interface{}{
                    "users": map[string]interface{}{
                        "profile": map[string]interface{}{
                            "id":    "userId",
                            "name":  "userName",
                            "email": "userEmail",
                        },
                    },
                    "roles": map[string]interface{}{
                        "assignments": map[string]interface{}{
                            "*": map[string]interface{}{
                                "scope": "roleAssignments[&1].accessScope",
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
