package contentrouter

import (
	"encoding/json"
	"fmt"
	"os"
)

type DocNode struct {
	Depth        int
	Parent       string
	Name         string
	MarkdownFile *string
}

func GenerateRoutes() {
	uDocs := GetUnstructuredDocs()
	docs := GetStructuredDocs(uDocs, "Root", []*DocNode{}, 0)
	for _, doc := range docs {
		fmt.Println(doc)
	}
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

func GetStructuredDocs(docs map[string]interface{}, parent string, docSlice []*DocNode, depth int) []*DocNode {
	for key, value := range docs {
		switch value := value.(type) {
		case string:
			if depth == 0 {
				parent = "Root"
			}
			docNode := &DocNode{
				Depth:        depth,
				Parent:       parent,
				Name:         key,
				MarkdownFile: &value,
			}
			docSlice = append(docSlice, docNode)
		case map[string]interface{}:
			docNode := &DocNode{
				Depth:        depth,
				Parent:       parent,
				Name:         key,
				MarkdownFile: nil,
			}
			docSlice = append(docSlice, docNode)
			docSlice = GetStructuredDocs(value, key, docSlice, depth+1)
		default:
			panic("Invalid type")
		}
	}
	return docSlice
}
