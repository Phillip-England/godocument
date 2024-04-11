package filewriter

import (
	"fmt"
	"godocument/internal/config"
	"godocument/internal/stypes"
	"os"
)

// GenerateDynamicNavbar generates the dynamic navbar based on ./godocument.config.json
func GenerateDynamicNavbar(cnf stypes.DocConfig) {
	html := `
		<nav id='sitenav'>
			<div class='sitenav-mobile-header'>
				<div class='sitenav-mobile-header-logo-wrapper'>
					<div class="sitenav-mobile-header-logo">
						<img src="/static/img/logo.svg" alt="logo" id="logo">
					</div>
					<svg class="sitenav-mobile-header-sun-icon sun-icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
					  <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5V3m0 18v-2M7.05 7.05 5.636 5.636m12.728 12.728L16.95 16.95M5 12H3m18 0h-2M7.05 16.95l-1.414 1.414M18.364 5.636 16.95 7.05M16 12a4 4 0 1 1-8 0 4 4 0 0 1 8 0Z"/>
					</svg>
				</div>
				<svg class='sitenav-mobile-header-close-icon' aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
					<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18 17.94 6M18 18 6.06 6"/>
			  	</svg>
			</div>
			<ul class='sitenav-list'>
	`
	for i := 0; i < len(cnf); i++ {
		html = workOnNavbar(cnf[i], html)
	}
	html += "</ul></nav>"
	writeNavbarHTML(html)
}

func workOnNavbar(node stypes.DocNode, html string) string {
	switch n := node.(type) {
	case *stypes.ObjectNode:
		dropdownCarat := `
			<svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="none" viewBox="0 0 24 24">
  				<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m9 5 7 7-7 7"/>
			</svg>
		`
		html += "<li class='sitenav-dropdown " + fmt.Sprintf("depth-%d", n.BaseNodeData.Depth) + "'>"
		html += "<button class='sitenav-dropdown-button sitenav-item'><summary>" + n.BaseNodeData.Name + "</summary><div class='dropdown-caret'>" + dropdownCarat + "</div></button>"
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
