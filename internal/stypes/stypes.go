package stypes

import "godocument/internal/middleware"

// represents our base template
type BaseTemplate struct {
	Title           string
	Content         string
	Prev            *MarkdownNode
	Next            *MarkdownNode
	MarkdownHeaders []MarkdownHeader
	MetaTags        []MarkdownMetaTag
}

// a slice of DocConfig representing the structured data
type DocConfig []DocNode

// each line in the godocument.config.json under the "docs" section is a DocNode
type DocNode interface{}

// all DocConfig should implement this type in their struct
type BaseNodeData struct {
	Depth   int
	Parent  string
	Name    string
	NavHTML string
}

// MarkdownNode represents a leaf node in the structured data
type MarkdownNode struct {
	BaseNodeData        *BaseNodeData
	MarkdownFile        string
	RouterPath          string
	StaticAssetPath     string
	Sequence            int
	Next                *MarkdownNode
	Prev                *MarkdownNode
	HandlerName         string
	HandlerUniqueString string
	HandlerFunc         middleware.CustomHandler
}

// MarkdownHeader represents a header in the markdown file
type MarkdownHeader struct {
	Line       string
	Link       string
	DepthClass string
}

// MarkdownMetaTag represets a meta tag found within a markdown file
type MarkdownMetaTag struct {
	Tag string
}

// ObjectNode represents a non-leaf node in the structured data
type ObjectNode struct {
	BaseNodeData *BaseNodeData
	Children     DocConfig
}
