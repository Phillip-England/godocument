package filewriter

import (
	"bytes"
	"fmt"
	"godocument/internal/config"
	"godocument/internal/stypes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

// GenerateDynamicNavbar generates the dynamic navbar based on ./godocument.config.json
func GenerateDynamicNavbar(cnf stypes.DocConfig) {
	html := `
		<nav id='sitenav' class='flex-col fixed z-max h-screen'>
			<div class='sitenav-mobile-header flex flex-row justify-between items-center text-md'>
				<div class='sitenav-mobile-header-logo-wrapper flex flex-row items-center justify-between'>
					<div class="sitenav-mobile-header-logo">
						<img class='logo' src="/static/img/logo.svg" alt="logo" id="logo">
					</div>
				</div>
				<div class='flex items-center'>
					<svg class="sun-icon cursor-pointer" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
						<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5V3m0 18v-2M7.05 7.05 5.636 5.636m12.728 12.728L16.95 16.95M5 12H3m18 0h-2M7.05 16.95l-1.414 1.414M18.364 5.636 16.95 7.05M16 12a4 4 0 1 1-8 0 4 4 0 0 1 8 0Z"/>
					</svg>
					<svg class="moon-icon cursor-pointer" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
  						<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 21a9 9 0 0 1-.5-17.986V3c-.354.966-.5 1.911-.5 3a9 9 0 0 0 9 9c.239 0 .254.018.488 0A9.004 9.004 0 0 1 12 21Z"/>
					</svg>
				</div>
			</div>
			<ul class='sitenav-list flex flex-col'>
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
		// to ensure tailwind classes are caught during build, we need to include this string
		// tailwind will find this sting and include all these classes in our output.css
		// this is needed because we are using integer values for padding in the navbar
		// and tailwind does not include these classes by default
		// for example, in the fmt.Sprintf() below we have included pl-%d which will be replaced by the integer value
		_ = "pl-0 pl-1 pl-2 pl-3 pl-4 pl-5 pl-6 pl-7"
		innerHTML := ""
		for i := 0; i < len(n.Children); i++ {
			innerHTML = workOnNavbar(n.Children[i], innerHTML)
		}
		html += fmt.Sprintf(`
		<li class='sitenav-dropdown flex flex-col pl-%d'>
			<button class='sitenav-dropdown-button sitenav-item flex flex-row justify-between items-center rounded-sm font-bold'>
				<summary>%s</summary>
				<div class='dropdown-caret'>
					<svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="none" viewBox="0 0 24 24">
						<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m9 5 7 7-7 7"/>
					</svg>
				</div>
			</button>
			<ul class='sitenav-dropdown-children hidden flex flex-col'>
				%s
			</ul>
		</li>
		`, n.BaseNodeData.Depth, n.BaseNodeData.Name, innerHTML)
	case *stypes.MarkdownNode:
		html += fmt.Sprintf(`
			<li class='pl-%d'>
				<a class='sitenav-item flex flex-row justify-between items-center rounded-sm font-bold' href='%s'>%s</a>
			</li>
		`, n.BaseNodeData.Depth, n.RouterPath, n.BaseNodeData.Name)
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

func GenerateStaticAssets(cnf stypes.DocConfig) {
	m := prepareMinify()
	ResetOutDir()
	copyDir(config.DevStaticPrefix, config.ProdStaticPrefix)
	generateStaticHTML(cnf)
	minifyStaticFiles(m, config.StaticAssetsDir)
}

func prepareMinify() *minify.M {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	return m
}

func ResetDocsDir() {
	err := os.RemoveAll(config.StaticMarkdownPrefix)
	if err != nil {
		fmt.Printf("Error removing directory: %s\n", err)
		return
	}
	if _, err := os.Stat(config.StaticMarkdownPrefix); os.IsNotExist(err) {
		err := os.Mkdir(config.StaticMarkdownPrefix, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %s\n", err)
			return
		}
	}
	filePath := fmt.Sprintf("%s/introduction.md", config.StaticMarkdownPrefix)
	content := "# Introduction"
	err = ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %s\n", err)
	}
}

func ResetOutDir() {
	err := os.RemoveAll(config.StaticAssetsDir)
	if err != nil {
		fmt.Printf("Error removing directory: %s\n", err)
		return
	}
	if _, err := os.Stat(config.StaticAssetsDir); os.IsNotExist(err) {
		err := os.Mkdir(config.StaticAssetsDir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %s\n", err)
			return
		}
	}
}

func ResetGodocumentConfig() {
	path := config.JSONConfigPath
	jsonData := "{\n\t\"docs\": {\n\t\t\"Introduction\": \"/introduction.md\"\n\t}\n}"

	// Write the JSON data to the file, creating it if it doesn't exist
	err := ioutil.WriteFile(path, []byte(jsonData), 0644)
	if err != nil {
		panic(err)
	}
}

func makeRequiredStaticDirs(staticAssetPath string) {
	dirPath := filepath.Dir(staticAssetPath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Printf("Failed to create directories: %v", err)
		return
	}
}

func createStaticAssetFile(staticAssetPath string) *os.File {
	f, err := os.Create(staticAssetPath)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		return nil
	}
	return f
}

func generateStaticHTML(cnf stypes.DocConfig) {
	config.WorkOnMarkdownNodes(cnf, func(n *stypes.MarkdownNode) {
		client := &http.Client{}
		res, err := client.Get("http://localhost:8080" + n.RouterPath)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading body: %s\n", err)
			return
		}
		makeRequiredStaticDirs(n.StaticAssetPath)
		f := createStaticAssetFile(n.StaticAssetPath)
		if f == nil {
			return
		}
		defer f.Close()
		doc, err := getQueryDoc(body)
		if err != nil {
			fmt.Printf("Error *goquery.Document from res.Body(): %s\n", err)
			return
		}
		modifyAnchorTagsForStatic(doc)
		setOtherAbsolutePaths(doc)
		htmlString, err := doc.Html()
		if err != nil {
			fmt.Printf("Error converting doc to html: %s\n", err)
			return
		}
		body = []byte(htmlString)
		_, err = f.Write(body)
		if err != nil {
			fmt.Printf("Error writing body to file: %s\n", err)
			return
		}
	})
}

func getQueryDoc(body []byte) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// prepares all anchor tags in the static html files to point to the correct .html file
// also converts relative paths to absolute paths
func modifyAnchorTagsForStatic(doc *goquery.Document) {
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		fmt.Println("SERVER_URL not set. Defaulting to http://localhost:8080")
		serverURL = "http://localhost:8080"
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			if len(href) > 3 && href[0:4] == "http" {
				return
			}
			if href[0] == '#' {
				return
			}
			if href == "/" {
				s.SetAttr("href", serverURL+"/")
				return
			}

			s.SetAttr("href", serverURL+href+".html")
		}
	})
}

// finds all local paths to static assets (other than anchor links) and converts them to absolute paths
func setOtherAbsolutePaths(doc *goquery.Document) {
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		fmt.Println("SERVER_URL not set. Defaulting to http://localhost:8080")
		serverURL = "http://localhost:8080"
	}
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			if len(href) > 3 && href[0:4] == "http" {
				return
			}
			s.SetAttr("href", serverURL+href)
		}
	})
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			if len(src) > 3 && src[0:4] == "http" {
				return
			}
			s.SetAttr("src", serverURL+src)
		}
	})
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			if len(src) > 3 && src[0:4] == "http" {
				return
			}
			s.SetAttr("src", serverURL+src)
		}
	})
}

func minifyStaticFiles(m *minify.M, dirPath string) {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		var mimetype string
		switch ext {
		case ".css":
			mimetype = "text/css"
		case ".html":
			mimetype = "text/html"
		case ".js":
			mimetype = "application/javascript"
		case ".json":
			mimetype = "application/json"
		case ".svg":
			mimetype = "image/svg+xml"
		case ".xml":
			mimetype = "text/xml"
		default:
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file: %s\n", err)
			return err
		}
		defer f.Close()
		fileBytes, err := io.ReadAll(f)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err)
			return err
		}
		minifiedBytes, err := m.Bytes(mimetype, fileBytes)
		if err != nil {
			fmt.Printf("Error minifying file: %s\n", err)
			return err
		}
		err = os.WriteFile(path, minifiedBytes, info.Mode()) // Preserving original file permissions
		if err != nil {
			fmt.Printf("Error writing minified file: %s\n", err)
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the directory: %s\n", err)
	}
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	return os.Chmod(dstFile.Name(), srcInfo.Mode())
}

// copyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func copyDir(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}
	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}
	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
