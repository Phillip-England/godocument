package contentrouter

import (
	"encoding/json"
	"fmt"
	"os"
)

func GenerateRoutes() {
	uDocs := GetUnstructuredDocs()
	TraverseUnstructuredDocs(uDocs, 0)
}

func GetUnstructuredDocs() map[string]interface{} {

	// Open the JSON file
	file, err := os.ReadFile("./godocument.config.json") // Ensure the file path and extension are correct
	if err != nil {
		panic(err)
	}

	// Unmarshal JSON to a map of interfaces
	var result map[string]interface{}
	err = json.Unmarshal(file, &result)
	if err != nil {
		panic(err)
	}

	// Get the "docs" key from the map
	docs := result["docs"]
	return IToStrMap(docs)

}

func IToStrMap(i interface{}) map[string]interface{} {
	return i.(map[string]interface{})
}

func TraverseUnstructuredDocs(docs map[string]interface{}, depth int) {
	for key, value := range docs {
		switch value := value.(type) {
		case string:
			fmt.Printf("%d | %s | %s\n", depth, key, value)
		case map[string]interface{}:
			if depth > 0 {
				fmt.Println("I am within another interface")
			}
			fmt.Printf("%d | %s\n", depth, key)
			TraverseUnstructuredDocs(value, depth+1)
		default:
			panic("Invalid type")
		}
	}
}

// func indent(depth int) string {
// 	indentation := ""
// 	for i := 0; i < depth; i++ {
// 		indentation += "  " // Using a tab for each level of depth
// 	}
// 	return indentation
// }
