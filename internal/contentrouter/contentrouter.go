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

type DocConfig []*DocNode

func GenerateRoutes() {
	uDocs := GetUnstructuredDocs()
	docs := GetStructuredDocs(uDocs, "Root", DocConfig{}, 0)
	for _, doc := range docs {
		fmt.Println(doc)
	}
}

// GetUnstructuredDocs reads the godocument.config.json file and returns the unstructured data
func GetUnstructuredDocs() map[string]interface{} {
	file, err := os.ReadFile("./godocument.config.json") // Ensure the file path and extension are correct
	if err != nil {
		panic(err)
	}
	var result map[string]interface{}
	err = json.Unmarshal(file, &result)
	if err != nil {
		panic(err)
	}
	docs := result["docs"]
	return docs.(map[string]interface{})
}

// GetStructuredDocs takes the unstructured data and returns a structured data
func GetStructuredDocs(docs map[string]interface{}, parent string, docSlice DocConfig, depth int) DocConfig {
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
