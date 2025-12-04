package jmap

// import (
// 	"fmt"
// 	"log"

// 	jmap "github.com/iammehrabsandhu/jmap/pkg"
// 	"github.com/iammehrabsandhu/jmap/types"
// )

// func main() {
// 	input := `{"items": [{"val": "a"}, {"val": "b"}]}`

// 	// Try to map to an array of objects
// 	spec := &types.TransformSpec{
// 		Operations: []types.Operation{
// 			{
// 				Type: "shift",
// 				Spec: map[string]interface{}{
// 					"items": map[string]interface{}{
// 						"*": map[string]interface{}{
// 							"val": "output[&1].value",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	result, err := jmap.Transform(input, spec)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Result 1 (Iterate):", result)

// 	// Try to map to a simple array
// 	spec2 := &types.TransformSpec{
// 		Operations: []types.Operation{
// 			{
// 				Type: "shift",
// 				Spec: map[string]interface{}{
// 					"items": map[string]interface{}{
// 						"*": map[string]interface{}{
// 							"val": "simpleList[0]", // This is still weird for list append, usually we want simpleList[&] or simpleList[]
// 							// But let's test explicit index
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	// For simple list append, jmap doesn't support "[]" syntax yet, but let's see if our fix supports [0]
// 	result2, err := jmap.Transform(input, spec2)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Result 2 (Simple List):", result2)
// }
