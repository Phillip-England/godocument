package contentrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"godocument/internal/middleware"
	"godocument/internal/util"

	"github.com/iancoleman/orderedmap"
	"github.com/yuin/goldmark"
)

const (
	DocRoot             = "Root"
	DocJSONKey          = "docs"
	JSONConfigPath      = "./godocument.config.json"
	IntroductionString  = "Introduction"
	GeneratedRoutesFile = "./internal/generated/generated.go"
)

// each line in the godocument.config.json under the "docs" section is a DocNode
type DocNode interface {
	Print()
}

// a slice of DocConfig representing the structured data
type DocConfig []DocNode

// all DocConfig should implement this type in their struct
type BaseNodeData struct {
	Depth   int
	Parent  string
	Name    string
	NavHTML string
}

// GetBaseData returns a string representation of the BaseNodeData
func (b *BaseNodeData) GetBaseData() string {
	return fmt.Sprintf("%d | %s | %s", b.Depth, b.Parent, b.Name)
}

// MarkdownNode represents a leaf node in the structured data
type MarkdownNode struct {
	BaseNodeData        *BaseNodeData
	MarkdownFile        string
	RouterPath          string
	Sequence            int
	Next                *MarkdownNode
	Prev                *MarkdownNode
	HandlerName         string
	HandlerUniqueString string
	HandlerFunc         middleware.CustomHandler
}

// Print prints the MarkdownNode data
func (b *MarkdownNode) Print() {
	baseData := b.BaseNodeData.GetBaseData()
	if b.Prev == nil {
		fmt.Printf("%s | %s | %s | %d | %s | %d | %s\n", baseData, b.MarkdownFile, b.RouterPath, b.Sequence, "nil", b.Next.Sequence, b.HandlerName)
		return

	}
	if b.Next == nil {
		fmt.Printf("%s | %s | %s | %d | %d | %s | %s\n", baseData, b.MarkdownFile, b.RouterPath, b.Sequence, b.Prev.Sequence, "nil", b.HandlerName)
		return

	}
	fmt.Printf("%s | %s | %s | %d | %d | %d | %s\n", baseData, b.MarkdownFile, b.RouterPath, b.Sequence, b.Prev.Sequence, b.Next.Sequence, b.HandlerName)
}

func (b *MarkdownNode) AssignHandlerName() {
	// generate 8 random characters to prevent name collisions
	b.HandlerUniqueString = util.RandomString(8)
	nameWithoutSpaces := strings.ReplaceAll(b.BaseNodeData.Name, " ", "")
	b.HandlerName = fmt.Sprintf("%s%s", nameWithoutSpaces, b.HandlerUniqueString)
}

// ObjectNode represents a non-leaf node in the structured data
type ObjectNode struct {
	BaseNodeData *BaseNodeData
	Children     DocConfig
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
func getLinearDocs(om interface{}, parent string, docConfig DocConfig, depth int) DocConfig {
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
				docConfig = append(docConfig, docNode)
			case orderedmap.OrderedMap:
				docNode := &ObjectNode{
					BaseNodeData: &BaseNodeData{
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
func sequenceMarkdownNodes(docConfig DocConfig) {
	sequence := 0
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *ObjectNode:
			continue
		case *MarkdownNode:
			docConfig[i].(*MarkdownNode).Sequence = sequence
			sequence++
		}
	}
}

// links each markdown node to the next markdown node based on their sequence number
func linkMarkdownNodes(docConfig DocConfig) {
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *ObjectNode:
			continue
		case *MarkdownNode:
			for j := 0; j < len(docConfig); j++ {
				switch docConfig[j].(type) {
				case *ObjectNode:
					continue
				case *MarkdownNode:
					if docConfig[j].(*MarkdownNode).Sequence == docConfig[i].(*MarkdownNode).Sequence+1 {
						docConfig[i].(*MarkdownNode).Next = docConfig[j].(*MarkdownNode)
					}
					if docConfig[j].(*MarkdownNode).Sequence == docConfig[i].(*MarkdownNode).Sequence-1 {
						docConfig[i].(*MarkdownNode).Prev = docConfig[j].(*MarkdownNode)
					}
				}
			}
		}

	}
}

// assignChildNodes assigns markdownNodes to each ObjectNode
func assignMarkdownNodes(docConfig DocConfig) {
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *ObjectNode:
			for j := 0; j < len(docConfig); j++ {
				switch docConfig[j].(type) {
				case *MarkdownNode:
					if docConfig[j].(*MarkdownNode).BaseNodeData.Parent == docConfig[i].(*ObjectNode).BaseNodeData.Name {
						docConfig[i].(*ObjectNode).Children = append(docConfig[i].(*ObjectNode).Children, docConfig[j])
					}
				}
			}
		}
	}
}

// purgeMarkdownNodes removes all MarkdownNodes (except Root-level markdown nodes) from the structured data and returns only ObjectNodes
func purgeMarkdownNodes(docConfig DocConfig) DocConfig {
	objectNodes := DocConfig{}
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *ObjectNode:
			objectNodes = append(objectNodes, docConfig[i])
		case *MarkdownNode:
			if docConfig[i].(*MarkdownNode).BaseNodeData.Parent == DocRoot {
				objectNodes = append(objectNodes, docConfig[i])
			}
		}
	}
	return objectNodes
}

// printDocConfig prints the structured data
func printDocConfig(docConfig DocConfig) {
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *ObjectNode:
			docConfig[i].(*ObjectNode).Print()
			printDocConfig(docConfig[i].(*ObjectNode).Children)
		case *MarkdownNode:
			docConfig[i].(*MarkdownNode).Print()
		}
	}
}

// workOnMarkdownNodes applies the action function to each MarkdownNode in the structured data
func workOnMarkdownNodes(docConfig DocConfig, action func(*MarkdownNode)) {
	for i := 0; i < len(docConfig); i++ {
		switch docConfig[i].(type) {
		case *ObjectNode:
			workOnMarkdownNodes(docConfig[i].(*ObjectNode).Children, action)
		case *MarkdownNode:
			action(docConfig[i].(*MarkdownNode))
		}
	}
}

// hookDocRoutes links our routes to the http.ServeMux
func hookDocRoutes(mux *http.ServeMux, docConfig DocConfig) {
	workOnMarkdownNodes(docConfig, func(m *MarkdownNode) {
		if m.BaseNodeData.Parent == DocRoot && m.BaseNodeData.Name == IntroductionString {
			mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/" {
					http.NotFound(w, r)
					return
				}
				middleware.Chain(w, r, docConfig[0].(*MarkdownNode).HandlerFunc)
			})
			return
		}
		mux.HandleFunc("GET "+m.RouterPath, func(w http.ResponseWriter, r *http.Request) {
			middleware.Chain(w, r, m.HandlerFunc)
		})
	})
}

func assignHandlers(docConfig DocConfig) {
	workOnMarkdownNodes(docConfig, func(m *MarkdownNode) {
		m.HandlerFunc = func(cc *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
			mdContent, err := os.ReadFile(m.MarkdownFile)
			if err != nil {
				// Handle error (e.g., file not found)
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			var buf bytes.Buffer
			if err := goldmark.Convert(mdContent, &buf); err != nil {
				http.Error(w, "Error converting markdown", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(buf.Bytes())
		}
	})
}

func workOnNavbar(node DocNode, html string) string {
	switch n := node.(type) {
	case *ObjectNode:
		html += "<ul class='dropdown'>"
		html += "<button class='dropbtn'>" + n.BaseNodeData.Name + "</button>"
		html += "<div class='dropdown-content'>"
		for i := 0; i < len(n.Children); i++ {
			html = workOnNavbar(n.Children[i], html)
		}
		html += "</ul>"
		html += "</div>"
	case *MarkdownNode:
		html += "<li><a href='" + n.RouterPath + "'>" + n.BaseNodeData.Name + "</a></li>"
	}
	return html
}

func generateDynamicNavbar(docConfig DocConfig) string {
	html := "<nav><ul>"
	for i := 0; i < len(docConfig); i++ {
		html = workOnNavbar(docConfig[i], html)
	}
	html += "</ul></nav>"
	return html
}

// GenerateRoutes generates code for application routes based on the ./godocument.config.json file "docs" section
// this function populates ./internal/generated/generated.go
func GenerateRoutes(mux *http.ServeMux) {
	uDocs := getUnstructuredDocs()
	docConfig := DocConfig{}
	for i := 0; i < len(uDocs.Keys()); i++ {
		key := uDocs.Keys()[i]
		value, _ := uDocs.Get(key)
		docConfig = getLinearDocs(value, DocRoot, docConfig, 0)
	}
	sequenceMarkdownNodes(docConfig)
	linkMarkdownNodes(docConfig)
	assignMarkdownNodes(docConfig)
	docConfig = purgeMarkdownNodes(docConfig)
	workOnMarkdownNodes(docConfig, func(m *MarkdownNode) {
		m.AssignHandlerName()
	})
	assignHandlers(docConfig)
	hookDocRoutes(mux, docConfig)
	navbarHTML := generateDynamicNavbar(docConfig)
	// write html to ./test.html
	f, err := os.Create("./test.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(navbarHTML)
	if err != nil {
		panic(err)
	}
	fmt.Println(navbarHTML)
}
