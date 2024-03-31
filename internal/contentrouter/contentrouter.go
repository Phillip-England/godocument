package contentrouter

import (
	"fmt"
	"godocument/internal/config"
	"os"
	"path/filepath"
	"strings"
)

// ContentPath is a struct that represents all contents of ./content directory.
// It can be either a directory itself, or a file.
type ContentPath struct {
	ExactPath string
	Parts     []string
}

func NewContentPath(filePath string) ContentPath {
	exactPath := "./" + filePath
	trimmedPath := strings.TrimPrefix(filePath, config.ContentPrefix)
	pathParts := strings.Split(trimmedPath, "/")
	return ContentPath{
		ExactPath: exactPath,
		Parts:     pathParts,
	}
}

// ContentDirectory is a struct that represents a directory in ./content directory
type ContentDirectory struct {
	ContentPath      ContentPath
	ChildDirectories []ContentDirectory
	ChildFiles       []ContentMarkdownFile
}

func NewContentDirectory(contentPath ContentPath) ContentDirectory {
	childDirectories := []ContentDirectory{}
	dirEntries, _ := os.ReadDir(contentPath.ExactPath)
	for _, entry := range dirEntries {
		if entry.IsDir() {
			contentDir := NewContentDirectory(NewContentPath(entry.Name()))
			childDirectories = append(childDirectories, contentDir)
		}
	}
	return ContentDirectory{
		ContentPath:      contentPath,
		ChildDirectories: childDirectories,
		ChildFiles:       []ContentMarkdownFile{},
	}
}

// ContentMarkdownFile is a struct that represents a markdown file in ./content directory
type ContentMarkdownFile struct {
	ContentPath ContentPath
}

func GenerateRoutes() {
	contentPaths := GetContentPaths()
	_ = GetContentDirectories(contentPaths)

}

func GetContentPaths() []ContentPath {
	contentPaths := []ContentPath{}
	err := filepath.WalkDir(config.ContentDirRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == config.ContentDirRoot {
			return nil
		}
		contentPath := NewContentPath(path)
		contentPaths = append(contentPaths, contentPath)
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", config.ContentDirRoot, err)
	}
	return contentPaths
}

func GetContentDirectories(contentPaths []ContentPath) []ContentDirectory {
	contentDirectories := []ContentDirectory{}
	for _, contentPath := range contentPaths {
		finalPart := contentPath.Parts[len(contentPath.Parts)-1]
		if !strings.Contains(finalPart, ".") {
			contentDirectory := NewContentDirectory(contentPath)
			contentDirectories = append(contentDirectories, contentDirectory)
		}

	}
	return contentDirectories
}
