package contentrouter

import (
	"encoding/json"
	"fmt"
	"os"
)

type Page struct {
	Name         string
	RouterPath   string
	MarkdownPath string
}

func (p Page) Print(lineIndention string) {
	fmt.Println(lineIndention + "Page:")
	fmt.Println(lineIndention + "  Name: " + p.Name)
	fmt.Println(lineIndention + "  RouterPath: " + p.RouterPath)
	fmt.Println(lineIndention + "  MarkdownPath: " + p.MarkdownPath)
}

type Series struct {
	Name      string
	Pages     []Page
	Subseries []Series
}

func (s Series) Print(lineIndentation string) {
	fmt.Println(lineIndentation + "Series: " + s.Name)
	for _, page := range s.Pages {
		page.Print(lineIndentation + "  ")
	}
	for _, sub := range s.Subseries {
		sub.Print(lineIndentation + "  ")
	}
}

func GenerateRoutes() {

	docConfigLines := GetDocConfigLines()
	PrintDocConfigLines(docConfigLines, 0)
}

func GetDocConfigLines() []interface{} {

	// docConfigLines is a slice of interfaces that can be either a Page or a Series
	// we do this to ensure the order of the pages and series in the config file is maintained
	docConfigLines := []interface{}{}

	// reading godocument.config.json file
	data, err := os.ReadFile("./godocument.config.json")
	if err != nil {
		panic("Error reading godocument.config.json file")
	}

	// Unmarshal JSON data into an interface
	var result interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		panic("Error unmarshalling JSON data in godocument.config.json - check syntax")
	}

	// Assert result to a map to access its contents
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		panic("Error asserting JSON data to a map in godocument.config.json - check syntax")
	}

	// Extract "docs" value (which is a map) from the JSON object
	docs, ok := resultMap["docs"]
	if !ok {
		panic("docs key not found in godocument.config.json - check syntax")
	}

	for key, value := range docs.(map[string]interface{}) {
		// if value is a string, it is not a series
		if _, ok := value.(string); ok {
			markdownPath := value.(string)
			routerPath := markdownPath[6 : len(markdownPath)-3]
			docConfigLines = append(docConfigLines, Page{
				Name:         key,
				RouterPath:   routerPath,
				MarkdownPath: markdownPath,
			})
		}
		// if value is a map, it is a series
		if _, ok := value.(map[string]interface{}); ok {
			s := ParseSeries(key, value.(map[string]interface{}))
			docConfigLines = append(docConfigLines, s)
		}
	}

	return docConfigLines

}

func ParseSeries(seriesName string, series map[string]interface{}) Series {
	pages := []Page{}
	subseries := []Series{}
	for key, value := range series {
		// if value is a string, it is not a series
		if _, ok := value.(string); ok {
			markdownPath := value.(string)
			routerPath := markdownPath[6 : len(markdownPath)-3]
			pages = append(pages, Page{
				Name:         key,
				RouterPath:   routerPath,
				MarkdownPath: value.(string),
			})
		}
		// if the value is a map, it is a subseries
		if _, ok := value.(map[string]interface{}); ok {
			subseries = append(subseries, ParseSeries(key, value.(map[string]interface{})))
		}
	}
	return Series{
		Name:      seriesName,
		Pages:     pages,
		Subseries: subseries,
	}
}

func PrintDocConfigLines(docConfigLines []interface{}, baseIndentation int) {
	lineIndentation := ""
	for i := 0; i < baseIndentation; i++ {
		lineIndentation += " "
	}
	for _, line := range docConfigLines {
		switch v := line.(type) {
		case Page:
			v.Print(lineIndentation)
		case Series:
			v.Print(lineIndentation)
		}
	}
}
