package contentrouter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	DocRoot            = "Root"
	DocJSONKey         = "docs"
	JSONConfigPath     = "./godocument.config.json"
	IntroductionString = "Introduction"
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

// GetBaseData returns a string representation of the BaseDocData
func (b *BaseDocData) GetBaseData() string {
	return fmt.Sprintf("%d | %s | %s", b.Depth, b.Parent, b.Name)
}

// DocNodeString represents a leaf node in the structured data
type DocNodeString struct {
	BaseDocData  *BaseDocData
	MarkdownFile string
	RouterPath   string
	Next         *DocNodeString
	Prev         *DocNodeString
}

// Print prints the DocNodeString data
func (b *DocNodeString) Print() {
	baseData := b.BaseDocData.GetBaseData()
	fmt.Printf("%s | %s | %s\n", baseData, b.MarkdownFile, b.RouterPath)
}

// DocNodeObject represents a non-leaf node in the structured data
type DocNodeObject struct {
	BaseDocData *BaseDocData
	Children    DocNodes
}

// Print prints the DocNodeObject data
func (b *DocNodeObject) Print() {
	baseData := b.BaseDocData.GetBaseData()
	fmt.Println(baseData)
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
			routerPath := ""
			if key == IntroductionString && depth == 0 {
				routerPath = "/"
			} else {
				routerPath = strings.TrimPrefix(strings.TrimSuffix(value, ".md"), "./docs")
			}
			docNode := &DocNodeString{
				BaseDocData: &BaseDocData{
					Depth:  depth,
					Parent: parent,
					Name:   key,
				},
				MarkdownFile: value,
				RouterPath:   routerPath,
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

// GenerateRoutes generates code for application routes based on the ./godocument.config.json file "docs" section
// this function populates ./internal/generated/generated.go
func GenerateRoutes() {
	uDocs := GetUnstructuredDocs()
	docs := GetStructuredDocs(uDocs, DocRoot, DocNodes{}, 0)
	for _, doc := range docs {
		doc.Print()
	}
}
