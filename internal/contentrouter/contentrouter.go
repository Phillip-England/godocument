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

func (c ContentPath) Print() {
	fmt.Println("Printing ContentPath:")
	fmt.Printf("  ExactPath: %s\n", c.ExactPath)
	fmt.Printf("  Parts: %v\n", c.Parts)
}

func GenerateRoutes() {
	contentPaths := GetContentPaths()
	parentDirs := IdentifyParentDirs(contentPaths)
	childDirs := IdentifyChildDirs(contentPaths, parentDirs)
	for _, childDir := range childDirs {
		childDir.Print()
	}

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

func IdentifyParentDirs(contentPaths []ContentPath) []ContentPath {
	parentDirs := []ContentPath{}
	for _, contentPath := range contentPaths {
		lastPart := contentPath.Parts[len(contentPath.Parts)-1]
		// if last part contains a "."" it is not a dir and we can skip it
		if strings.Contains(lastPart, ".") {
			continue
		}
		// if the parts only contains one element, it is the root dir and therefore is a parent
		if len(contentPath.Parts) == 1 {
			parentDirs = append(parentDirs, contentPath)
			continue
		}
		dirEntries, err := os.ReadDir(contentPath.ExactPath)
		if err != nil {
			panic(err)
		}
		for _, dirEntry := range dirEntries {
			if dirEntry.IsDir() {
				parentDirs = append(parentDirs, contentPath)
			}
		}

	}
	return parentDirs
}

func IdentifyChildDirs(contentPaths []ContentPath, parentDirs []ContentPath) []ContentPath {
	childDirs := []ContentPath{}
	for _, contentPath := range contentPaths {
		lastPart := contentPath.Parts[len(contentPath.Parts)-1]
		// if last part contains a "."" it is not a dir and we can skip it
		if strings.Contains(lastPart, ".") {
			continue
		}
		// if the parts only contains one element, it is a root dir and therefore is not a child
		if len(contentPath.Parts) == 1 {
			continue
		}
	}
	return childDirs
}
