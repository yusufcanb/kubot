package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BuildDirectoryTree builds the workspace tree
func buildDirectoryTree(basePath string) (DirectoryNode, error) {
	node := DirectoryNode{
		Path: basePath,
	}
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == basePath {
			return nil
		}
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			parentNode := &node
			for _, dir := range strings.Split(relPath, string(filepath.Separator)) {
				found := false
				for i, child := range parentNode.Children {
					if child.Path == dir {
						parentNode = &parentNode.Children[i]
						found = true
						break
					}
				}
				if !found {
					newNode := DirectoryNode{
						Path:     dir,
						Children: make([]DirectoryNode, 0),
						Files:    make([]string, 0),
					}
					parentNode.Children = append(parentNode.Children, newNode)
					parentNode = &newNode
				}
			}
		} else {
			parentNode := &node
			dirs := strings.Split(relPath, string(filepath.Separator))
			for _, dir := range dirs[:len(dirs)-1] {
				found := false
				for i, child := range parentNode.Children {
					if child.Path == dir {
						parentNode = &parentNode.Children[i]
						found = true
						break
					}
				}
				if !found {
					newNode := DirectoryNode{
						Path:     dir,
						Children: make([]DirectoryNode, 0),
						Files:    make([]string, 0),
					}
					parentNode.Children = append(parentNode.Children, newNode)
					parentNode = &newNode
				}
			}
			parentNode.Files = append(parentNode.Files, filepath.Base(path))
		}
		return nil
	})
	return node, err
}

func printDirectoryTree(node DirectoryNode, indent string) {
	fmt.Printf("%s%s/\n", indent, node.Path)
	for _, child := range node.Children {
		printDirectoryTree(child, indent+"  ")
	}

	for _, file := range node.Files {
		fmt.Printf("%s  %s\n", indent, file)
	}
}
