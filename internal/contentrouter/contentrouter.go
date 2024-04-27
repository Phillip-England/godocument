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
func GenerateRoutes(mux *http.ServeMux, templates *template.Template) {
	cnf := config.GetDocConfig()
	assignHandlers(cnf)
	hookDocRoutes(mux, templates, cnf)
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

			// reading all headers from the markdown file content and parsing out any <meta> tags
			file, err := os.Open(m.MarkdownFile)
			if err != nil {
				http.Error(w, "Error reading markdown file", http.StatusInternalServerError)
				return
			}
			defer file.Close()
			// read the file line by line
			scanner := bufio.NewScanner(file)
			headers := []stypes.MarkdownHeader{}
			metaTags := []stypes.MarkdownMetaTag{}
			backticksFound := 0
			for scanner.Scan() {
				line := scanner.Text()
				// skipping and backticks
				if strings.Contains(line, "```") {
					backticksFound++
					continue
				}
				// if we are in a code block, skip the line
				if backticksFound%2 == 1 {
					continue
				}
				// if we find a meta tag,
				if strings.Contains(line, "<meta") && strings.Contains(line, ">") {
					metaTag := stypes.MarkdownMetaTag{
						Tag: line,
					}
					metaTags = append(metaTags, metaTag)
					continue
				}
				if strings.Contains(line, "# ") {
					if strings.Count(line, "#") == 1 {
						continue
					}
					link := strings.ToLower(line)
					link = strings.TrimLeft(link, "# ")
					// when parsing headers using goldmark, the id of the header will have special chars stripped
					// for example a header of /docs will be converted to docs
					// also, spaces are replaced with dashes and text is converted to lowercase
					// for a header of Custom Components will be converted to custom-components
					// so when pointing to these headers, we need to ensure we are using the correct id
					// therefore, we need to also strip the hrefs of our pagenav links as well
					// here we are stripping those characters
					// if you find a case where a pagenav link is not working, this is most likely the issue
					// be sure to check their hrefs to see if they line up with the associated headers id
					unwantedLinkChars := []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+", "=", "{", "}", "[", "]", "|", "\\", ":", ";", "\"", "'", "<", ">", ",", ".", "/", "?"}
					for _, char := range unwantedLinkChars {
						link = strings.ReplaceAll(link, char, "")
					}
					link = strings.ReplaceAll(link, " ", "-")
					link = "#" + link
					depthClass := ""
					leftPadding := strings.Count(line, "#") - 2 // ensures we start at 0
					if leftPadding != 0 {
						depthClass = "pl-" + fmt.Sprintf("%d", leftPadding+2)
					}
					header := stypes.MarkdownHeader{
						Line:       strings.TrimLeft(line, "# "),
						DepthClass: depthClass, // ensures we start at a 0 index
						Link:       link,
					}
					headers = append(headers, header)
				}
			}

			// removing meta tags from the markdown content
			mdContentString := mdBuf.String()
			for _, metaTag := range metaTags {
				mdContentString = strings.Replace(mdContentString, metaTag.Tag, "", 1)
			}

			// getting title from godocument.config.json
			title := config.GetTitle()
			if title == "" {
				title = "Godocument"
			}

			// Create a new instance of tdata.Base with the title and markdown content as HTML
			baseData := &stypes.BaseTemplate{
				Title:           title + " - " + m.BaseNodeData.Name,
				Content:         mdContentString,
				Prev:            m.Prev,
				Next:            m.Next,
				MarkdownHeaders: headers,
				MetaTags:        metaTags,
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
