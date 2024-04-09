package filewriter

import (
	"fmt"
	"godocument/internal/config"
	"godocument/internal/stypes"
	"os"
)

// GenerateDynamicNavbar generates the dynamic navbar based on ./godocument.config.json
func GenerateDynamicNavbar(cnf stypes.DocConfig) {
	html := "<nav id='sitenav'><ul>"
	for i := 0; i < len(cnf); i++ {
		html = workOnNavbar(cnf[i], html)
	}
	html += "</ul></nav>"
	writeNavbarHTML(html)
}

func workOnNavbar(node stypes.DocNode, html string) string {
	switch n := node.(type) {
	case *stypes.ObjectNode:
		html += "<li class='sitenav-dropdown " + fmt.Sprintf("depth-%d", n.BaseNodeData.Depth) + "'>"
		html += "<button class='sitenav-dropdown-button sitenav-item'><summary>" + n.BaseNodeData.Name + "</summary><div class='dropdown-caret'>></div></button>"
		html += "<ul class='sitenav-dropdown-children hidden'>"
		for i := 0; i < len(n.Children); i++ {
			html = workOnNavbar(n.Children[i], html)
		}
		html += "</ul>"
		html += "</li>"
	case *stypes.MarkdownNode:
		html += "<li class='" + fmt.Sprintf("depth-%d", n.BaseNodeData.Depth) + "'><a class='sitenav-item' href='" + n.RouterPath + "'>" + n.BaseNodeData.Name + "</a></li>"
	}
	return html
}

// writeNavbarHTML writes the generated navbar html to ./template/generated-nav.html
func writeNavbarHTML(html string) {
	f, err := os.Create(config.GeneratedNavPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString("<!-- This file is auto-generated. Do not modify. -->\n")
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(html)
	if err != nil {
		panic(err)
	}
}
