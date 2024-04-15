package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"godocument/internal/stypes"

	"github.com/iancoleman/orderedmap"
)

const (
	DocRoot              = "Root"
	DocJSONKey           = "docs"
	JSONConfigPath       = "./godocument.config.json"
	IntroductionString   = "Introduction"
	GeneratedNavPath     = "./html/components/sitenav.html"
	StaticMarkdownPrefix = "./docs"
	StaticAssetsDir      = "./out"
)

// taks the "docs" section of godocument.json.config and generates a workable data structure from it
func GetDocConfig() stypes.DocConfig {

	// here is how we take the json found in ./godocument.config.json and generate the data in a structured format
	// each "line" in the json file is a DocNode (an interface that represents all lines in the "docs" section of the json file)
	// first, we get each line using orderedmap.OrderedMap in getUnstructuredDocs()
	// the order of the lines is important because it will dictate the arrangement of html components
	// then we generate a slice of each line in the json file using getLinearDocs()
	// we sequence each markdown node so it is easy for us to link them together
	// then, each markdown node is linked to eachother based on their sequence number
	// nodes not found at the root level have a parent node assigned to them
	// in assignMarkdownNodes, we assign each markdown node to their respective parent object node
	// in assignSubObjectNodes, we assign each object node to their respective parent object node
	// after doing this, we purge all markdown nodes that are not at the root level (because they now exist in objectnode.Children)
	// we also purge all object nodes that are not at the root level (because they will be assigned to another object node)

	u := getUnstructuredDocs()
	c := stypes.DocConfig{}
	for i := 0; i < len(u.Keys()); i++ {
		key := u.Keys()[i]
		if key == DocJSONKey {
			value, _ := u.Get(key)
			c = getLinearDocs(value, DocRoot, c, 0)
		}
	}
	sequenceMarkdownNodes(c)
	linkMarkdownNodes(c)
	assignMarkdownNodes(c)
	assignSubObjectNodes(c)
	c = purgeMarkdownNodes(c)
	c = purgeNonRootObjectNodes(c)
	ensureMarkdownsFileExists(c)
	return c
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
func getLinearDocs(om interface{}, parent string, docConfig stypes.DocConfig, depth int) stypes.DocConfig {
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
				markdownFile := StaticMarkdownPrefix + value
				staticAssetPath := markdownFile
				if key == IntroductionString && depth == 0 {
					routerPath = "/"
					staticAssetPath = StaticAssetsDir + "/index.html"
				} else {
					routerPath = strings.TrimSuffix(value, ".md")
					staticAssetPath = strings.TrimPrefix(staticAssetPath, StaticMarkdownPrefix)
					staticAssetPath = StaticAssetsDir + staticAssetPath
					staticAssetPath = strings.Replace(staticAssetPath, ".md", ".html", 1)
				}
				docNode := &stypes.MarkdownNode{
					BaseNodeData: &stypes.BaseNodeData{
						Depth:  depth,
						Parent: parent,
						Name:   key,
					},
					MarkdownFile:    markdownFile,
					RouterPath:      routerPath,
					StaticAssetPath: staticAssetPath,
				}
				docConfig = append(docConfig, docNode)
			case orderedmap.OrderedMap:
				docNode := &stypes.ObjectNode{
					BaseNodeData: &stypes.BaseNodeData{
						Depth:  depth,
						Parent: parent,
						Name:   key,
					},
					Children: nil,
				}
				docConfig = append(docConfig, docNode)
				docConfig = getLinearDocs(value, key, docConfig, depth+1)
			}
		}
	case string:
		return docConfig
	case nil:
		return docConfig
	default:
		panic("Invalid type")
	}
	return docConfig
}

// sequenceMarkdownNodes assigns a sequence number to each MarkdownNode
func sequenceMarkdownNodes(docConfig stypes.DocConfig) {
	sequence := 0
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *stypes.ObjectNode:
			continue
		case *stypes.MarkdownNode:
			docConfig[i].(*stypes.MarkdownNode).Sequence = sequence
			sequence++
		}
	}
}

// links each markdown node to the next markdown node based on their sequence number
func linkMarkdownNodes(docConfig stypes.DocConfig) {
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *stypes.ObjectNode:
			continue
		case *stypes.MarkdownNode:
			for j := 0; j < len(docConfig); j++ {
				switch docConfig[j].(type) {
				case *stypes.ObjectNode:
					continue
				case *stypes.MarkdownNode:
					if docConfig[j].(*stypes.MarkdownNode).Sequence == docConfig[i].(*stypes.MarkdownNode).Sequence+1 {
						docConfig[i].(*stypes.MarkdownNode).Next = docConfig[j].(*stypes.MarkdownNode)
					}
					if docConfig[j].(*stypes.MarkdownNode).Sequence == docConfig[i].(*stypes.MarkdownNode).Sequence-1 {
						docConfig[i].(*stypes.MarkdownNode).Prev = docConfig[j].(*stypes.MarkdownNode)
					}
				}
			}
		}

	}
}

// assignChildNodes assigns markdownNodes to each ObjectNode
func assignMarkdownNodes(docConfig stypes.DocConfig) {
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *stypes.ObjectNode:
			for j := 0; j < len(docConfig); j++ {
				switch docConfig[j].(type) {
				case *stypes.MarkdownNode:
					if docConfig[j].(*stypes.MarkdownNode).BaseNodeData.Parent == docConfig[i].(*stypes.ObjectNode).BaseNodeData.Name {
						docConfig[i].(*stypes.ObjectNode).Children = append(docConfig[i].(*stypes.ObjectNode).Children, docConfig[j])
					}
				}
			}
		}
	}
}

// assignSubObjectNodes assigns ObjectNodes to their respective parent ObjectNode
func assignSubObjectNodes(docConfig stypes.DocConfig) {
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *stypes.ObjectNode:
			for j := 0; j < len(docConfig); j++ {
				switch docConfig[j].(type) {
				case *stypes.ObjectNode:
					if docConfig[j].(*stypes.ObjectNode).BaseNodeData.Parent == docConfig[i].(*stypes.ObjectNode).BaseNodeData.Name {
						docConfig[i].(*stypes.ObjectNode).Children = append(docConfig[i].(*stypes.ObjectNode).Children, docConfig[j])
					}
				}
			}
		}
	}

}

// purgeMarkdownNodes removes all MarkdownNodes (except Root-level markdown nodes) from the structured data and returns only ObjectNodes
func purgeMarkdownNodes(docConfig stypes.DocConfig) stypes.DocConfig {
	objectNodes := stypes.DocConfig{}
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *stypes.ObjectNode:
			objectNodes = append(objectNodes, docConfig[i])
		case *stypes.MarkdownNode:
			if docConfig[i].(*stypes.MarkdownNode).BaseNodeData.Parent == DocRoot {
				objectNodes = append(objectNodes, docConfig[i])
			}
		}
	}
	return objectNodes
}

// purgeNonRootObjectNodes removes all ObjectNodes that are not at the root level
func purgeNonRootObjectNodes(docConfig stypes.DocConfig) stypes.DocConfig {
	rootObjectNodes := stypes.DocConfig{}
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *stypes.ObjectNode:
			if docConfig[i].(*stypes.ObjectNode).BaseNodeData.Parent == DocRoot {
				rootObjectNodes = append(rootObjectNodes, docConfig[i])
			}
		case *stypes.MarkdownNode:
			rootObjectNodes = append(rootObjectNodes, docConfig[i])
		}
	}
	return rootObjectNodes

}

// workOnMarkdownNodes applies the action function to each MarkdownNode in the structured data
func WorkOnMarkdownNodes(docConfig stypes.DocConfig, action func(*stypes.MarkdownNode)) {
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *stypes.ObjectNode:
			WorkOnMarkdownNodes(docConfig[i].(*stypes.ObjectNode).Children, action)
		case *stypes.MarkdownNode:
			action(docConfig[i].(*stypes.MarkdownNode))
		}
	}
}

// ensures all markdown files in godocument.config.json exist
func ensureMarkdownsFileExists(docConfig stypes.DocConfig) {
	WorkOnMarkdownNodes(docConfig, func(m *stypes.MarkdownNode) {
		if _, err := os.Stat(m.MarkdownFile); os.IsNotExist(err) {
			panic(fmt.Sprintf("Markdown file %s does not exist", m.MarkdownFile))
		}
	})
}
