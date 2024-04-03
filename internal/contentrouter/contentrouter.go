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
type BaseNodeData struct {
	Depth  int
	Parent string
	Name   string
}

// GetBaseData returns a string representation of the BaseNodeData
func (b *BaseNodeData) GetBaseData() string {
	return fmt.Sprintf("%d | %s | %s", b.Depth, b.Parent, b.Name)
}

// MarkdownNode represents a leaf node in the structured data
type MarkdownNode struct {
	BaseNodeData *BaseNodeData
	MarkdownFile string
	RouterPath   string
	Sequence     int
	Next         *MarkdownNode
	Prev         *MarkdownNode
}

// Print prints the MarkdownNode data
func (b *MarkdownNode) Print() {
	baseData := b.BaseNodeData.GetBaseData()
	if b.Prev == nil {
		fmt.Printf("%s | %s | %s | %d | %s | %d\n", baseData, b.MarkdownFile, b.RouterPath, b.Sequence, "nil", b.Next.Sequence)
		return

	}
	if b.Next == nil {
		fmt.Printf("%s | %s | %s | %d | %d | %s\n", baseData, b.MarkdownFile, b.RouterPath, b.Sequence, b.Prev.Sequence, "nil")
		return

	}
	fmt.Printf("%s | %s | %s | %d | %d | %d\n", baseData, b.MarkdownFile, b.RouterPath, b.Sequence, b.Prev.Sequence, b.Next.Sequence)
}

// ObjectNode represents a non-leaf node in the structured data
type ObjectNode struct {
	BaseNodeData *BaseNodeData
	Children     DocNodes
}

// Print prints the ObjectNode data
func (b *ObjectNode) Print() {
	baseData := b.BaseNodeData.GetBaseData()
	fmt.Printf("%s\n", baseData)
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

// GetStructuredDocs recursively generates a structured representation of the unstructured doc config data
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
				docNode := &MarkdownNode{
					BaseNodeData: &BaseNodeData{
						Depth:  depth,
						Parent: parent,
						Name:   key,
					},
					MarkdownFile: value,
					RouterPath:   routerPath,
				}
				docNodes = append(docNodes, docNode)
			case orderedmap.OrderedMap:
				docNode := &ObjectNode{
					BaseNodeData: &BaseNodeData{
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

// sequenceMarkdownNodes assigns a sequence number to each MarkdownNode
func sequenceMarkdownNodes(docNodes DocNodes) {
	sequence := 0
	for i := 0; i < len(docNodes); i++ {
		switch docNodes[i].(type) {
		case *ObjectNode:
			continue
		case *MarkdownNode:
			docNodes[i].(*MarkdownNode).Sequence = sequence
			sequence++
		}
	}
}

// links each markdown node to the next markdown node based on their sequence number
func linkMarkdownNodes(docNodes DocNodes) {
	for i := 0; i < len(docNodes); i++ {
		switch docNodes[i].(type) {
		case *ObjectNode:
			continue
		case *MarkdownNode:
			for j := 0; j < len(docNodes); j++ {
				switch docNodes[j].(type) {
				case *ObjectNode:
					continue
				case *MarkdownNode:
					if docNodes[j].(*MarkdownNode).Sequence == docNodes[i].(*MarkdownNode).Sequence+1 {
						docNodes[i].(*MarkdownNode).Next = docNodes[j].(*MarkdownNode)
					}
					if docNodes[j].(*MarkdownNode).Sequence == docNodes[i].(*MarkdownNode).Sequence-1 {
						docNodes[i].(*MarkdownNode).Prev = docNodes[j].(*MarkdownNode)
					}
				}
			}
		}

	}
}

// assignChildNodes assigns markdownNodes to each ObjectNode
func assignMarkdownNodes(docNodes DocNodes) {
	for i := 0; i < len(docNodes); i++ {
		switch docNodes[i].(type) {
		case *ObjectNode:
			for j := 0; j < len(docNodes); j++ {
				switch docNodes[j].(type) {
				case *MarkdownNode:
					if docNodes[j].(*MarkdownNode).BaseNodeData.Parent == docNodes[i].(*ObjectNode).BaseNodeData.Name {
						docNodes[i].(*ObjectNode).Children = append(docNodes[i].(*ObjectNode).Children, docNodes[j])
					}
				}
			}
		}
	}
}

// purgeMarkdownNodes removes all MarkdownNodes (except Root-level markdown nodes) from the structured data and returns only ObjectNodes
func purgeMarkdownNodes(docNodes DocNodes) DocNodes {
	objectNodes := DocNodes{}
	for i := 0; i < len(docNodes); i++ {
		switch docNodes[i].(type) {
		case *ObjectNode:
			objectNodes = append(objectNodes, docNodes[i])
		case *MarkdownNode:
			if docNodes[i].(*MarkdownNode).BaseNodeData.Parent == DocRoot {
				objectNodes = append(objectNodes, docNodes[i])
			}
		}
	}
	return objectNodes
}

// printDocNodes prints the structured data
func printDocNodes(docNodes DocNodes) {
	for i := 0; i < len(docNodes); i++ {
		switch docNodes[i].(type) {
		case *ObjectNode:
			docNodes[i].(*ObjectNode).Print()
			printDocNodes(docNodes[i].(*ObjectNode).Children)
		case *MarkdownNode:
			docNodes[i].(*MarkdownNode).Print()
		}
	}

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
	sequenceMarkdownNodes(docNodes)
	linkMarkdownNodes(docNodes)
	assignMarkdownNodes(docNodes)
	sortedNodes := purgeMarkdownNodes(docNodes)
	printDocNodes(sortedNodes)
	// docs := GetStructuredDocs(uDocs, DocRoot, DocNodes{}, 0)

}
