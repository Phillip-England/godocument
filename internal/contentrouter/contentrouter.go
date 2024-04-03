package contentrouter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/iancoleman/orderedmap"
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
func getUnstructuredDocs() *orderedmap.OrderedMap {
	file, err := os.ReadFile(JSONConfigPath) // Ensure the file path and extension are correct
	if err != nil {
		panic(err)
	}
	result := orderedmap.New()
	err = json.Unmarshal(file, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func getLinearDocs(om interface{}, parent string, docNodes DocNodes, depth int) DocNodes {
	switch om := om.(type) {
	case orderedmap.OrderedMap:
		for _, key := range om.Keys() {
			value, _ := om.Get(key)
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
			case orderedmap.OrderedMap:
				docNode := &DocNodeObject{
					BaseDocData: &BaseDocData{
						Depth:  depth,
						Parent: parent,
						Name:   key,
					},
					Children: nil,
				}
				docNodes = append(docNodes, docNode)
				docNodes = getLinearDocs(value, key, docNodes, depth+1)
			}
		}
	case string:
		return docNodes
	case nil:
		return docNodes
	default:
		panic("Invalid type")
	}
	return docNodes
}

// GenerateRoutes generates code for application routes based on the ./godocument.config.json file "docs" section
// this function populates ./internal/generated/generated.go
func GenerateRoutes() {
	uDocs := getUnstructuredDocs()

	docNodes := DocNodes{}
	for i := 0; i < len(uDocs.Keys()); i++ {
		key := uDocs.Keys()[i]
		value, _ := uDocs.Get(key)
		docNodes = getLinearDocs(value, DocRoot, docNodes, 0)
	}
	for _, docNode := range docNodes {
		docNode.Print()
	}
	// docs := GetStructuredDocs(uDocs, DocRoot, DocNodes{}, 0)

}
