package contentrouter

import (
	"encoding/json"
	"fmt"
	"os"
)

type NodeWithoutDropdown struct {
	Parent       string
	Name         string
	MarkdownPath string
}

type NodeWithDropdown struct {
	Parent   string
	Name     string
	Children map[string]interface{}
}

func GenerateRoutes() {
	uDocs := GetUnstructuredDocs()
	TraverseUnstructuredDocs(uDocs, "Root", 0)
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

func TraverseUnstructuredDocs(docs map[string]interface{}, parent string, depth int) {
	docConfigLines := []interface{}{}
	for key, value := range docs {
		switch value := value.(type) {
		case string:
			if depth == 0 {
				parent = "Root"
			}
			fmt.Printf("%d | %s | %s | %s\n", depth, parent, key, value)
		case map[string]interface{}:
			fmt.Printf("%d | %s | %s\n", depth, parent, key)
			TraverseUnstructuredDocs(value, key, depth+1)
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
