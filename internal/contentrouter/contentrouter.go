package contentrouter

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	DocRoot        = "Root"
	DocJSONKey     = "docs"
	JSONConfigPath = "./godocument.config.json"
)

// each line in the godocument.config.json under the "docs" section is a DocNode
type DocNode interface {
	Print()
}

// a slice of DocNodes representing the structured data
type DocNodes []DocNode

// all DocNodes should implement this type in their struct
type BaseDocData struct {
	Depth  int
	Parent string
	Name   string
}

func (b *BaseDocData) GetBaseData() string {
	return fmt.Sprintf("%d | %s | %s | ", b.Depth, b.Parent, b.Name)
}

type DocNodeString struct {
	BaseDocData  *BaseDocData
	MarkdownFile string
}

func (b *DocNodeString) Print() {
	baseData := b.BaseDocData.GetBaseData()
	fmt.Println(baseData + b.MarkdownFile)
}

type DocNodeObject struct {
	BaseDocData *BaseDocData
	Children    DocNodes
}

func (b *DocNodeObject) Print() {
	baseData := b.BaseDocData.GetBaseData()
	fmt.Println(baseData)
}

func GenerateRoutes() {
	uDocs := GetUnstructuredDocs()
	docs := GetStructuredDocs(uDocs, DocRoot, DocNodes{}, 0)
	for _, doc := range docs {
		doc.Print()
	}
}

// GetUnstructuredDocs reads the godocument.config.json file and returns the unstructured data
func GetUnstructuredDocs() map[string]interface{} {
	file, err := os.ReadFile(JSONConfigPath) // Ensure the file path and extension are correct
	if err != nil {
		panic(err)
	}
	var result map[string]interface{}
	err = json.Unmarshal(file, &result)
	if err != nil {
		panic(err)
	}
	docs := result[DocJSONKey]
	return docs.(map[string]interface{})
}

// GetStructuredDocs takes the unstructured data and returns a structured data
func GetStructuredDocs(docs map[string]interface{}, parent string, docNodes DocNodes, depth int) DocNodes {
	for key, value := range docs {
		switch value := value.(type) {
		case string:
			if depth == 0 {
				parent = DocRoot
			}
			docNode := &DocNodeString{
				BaseDocData: &BaseDocData{
					Depth:  depth,
					Parent: parent,
					Name:   key,
				},
				MarkdownFile: value,
			}
			docNodes = append(docNodes, docNode)
		case map[string]interface{}:
			docNode := &DocNodeObject{
				BaseDocData: &BaseDocData{
					Depth:  depth,
					Parent: parent,
					Name:   key,
				},
				Children: nil,
			}
			docNodes = append(docNodes, docNode)
			docNodes = GetStructuredDocs(value, key, docNodes, depth+1)
		default:
			panic("Invalid type")
		}
	}
	return docNodes
}
