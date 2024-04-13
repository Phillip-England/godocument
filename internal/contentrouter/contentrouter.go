package contentrouter

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"godocument/internal/config"
	"godocument/internal/middleware"
	"godocument/internal/stypes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// GenerateRoutes generates code for application routes based on the ./godocument.config.json file "docs" section
// this function populates ./internal/generated/generated.go
func GenerateRoutes(mux *http.ServeMux, templates *template.Template) stypes.DocConfig {
	cnf := config.GetDocConfig()
	assignHandlers(cnf)
	hookDocRoutes(mux, templates, cnf)
	return cnf
}

// hookDocRoutes links our routes to the http.ServeMux
func hookDocRoutes(mux *http.ServeMux, templates *template.Template, cnf stypes.DocConfig) {
	config.WorkOnMarkdownNodes(cnf, func(m *stypes.MarkdownNode) {
		if m.BaseNodeData.Parent == config.DocRoot && m.BaseNodeData.Name == config.IntroductionString {
			mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/" {
					http.NotFound(w, r)
					return
				}
				middleware.Chain(w, r, templates, m.HandlerFunc)
			})
			return
		}
		mux.HandleFunc("GET "+m.RouterPath, func(w http.ResponseWriter, r *http.Request) {
			middleware.Chain(w, r, templates, m.HandlerFunc)
		})
	})
}

func assignHandlers(cnf stypes.DocConfig) {
	config.WorkOnMarkdownNodes(cnf, func(m *stypes.MarkdownNode) {
		m.HandlerFunc = func(cc *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
			md := goldmark.New(
				goldmark.WithParserOptions(
					parser.WithAutoHeadingID(),
				),
				goldmark.WithRendererOptions(
					html.WithUnsafe(),
				),
				// goldmark.WithExtensions(
				// 	highlighting.NewHighlighting(
				// 		highlighting.WithStyle("github"),
				// 	),
				// ),
			)
			mdContent, err := os.ReadFile(m.MarkdownFile)
			if err != nil {
				// Handle error (e.g., file not found)
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			var mdBuf bytes.Buffer
			if err := md.Convert(mdContent, &mdBuf); err != nil {
				http.Error(w, "Error converting markdown", http.StatusInternalServerError)
				return
			}

			// reading all headers from the markdown file content
			file, err := os.Open(m.MarkdownFile)
			if err != nil {
				http.Error(w, "Error reading markdown file", http.StatusInternalServerError)
				return
			}
			defer file.Close()
			// read the file line by line
			scanner := bufio.NewScanner(file)
			headers := []stypes.MarkdownHeader{}
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, "# ") {
					if strings.Count(line, "#") == 1 {
						continue
					}
					header := stypes.MarkdownHeader{
						Line:       strings.TrimLeft(line, "# "),
						DepthClass: "depth-" + fmt.Sprintf("%d", strings.Count(line, "#")-2), // ensures we start at a 0 index
						Link:       "#" + strings.ToLower(strings.ReplaceAll(strings.TrimLeft(line, "# "), " ", "-")),
					}
					headers = append(headers, header)
				}
			}

			// Create a new instance of tdata.Base with the title and markdown content as HTML
			baseData := &stypes.BaseTemplate{
				Title:           "Godocument - " + m.BaseNodeData.Name,
				Content:         mdBuf.String(),
				Prev:            m.Prev,
				Next:            m.Next,
				MarkdownHeaders: headers,
			}

			// Assuming you have already parsed your templates (including the base template) elsewhere
			tmpl := cc.Templates.Lookup("base.html")
			if tmpl == nil {
				http.Error(w, "Base template not found", http.StatusInternalServerError)
				return
			}

			// Execute the base template with the baseData instance
			var htmlBuf bytes.Buffer
			if err := tmpl.Execute(&htmlBuf, baseData); err != nil {
				fmt.Println(err)
				http.Error(w, "Error executing template", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			w.Write(htmlBuf.Bytes())
		}
	})
}
