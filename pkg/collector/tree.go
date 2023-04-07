package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DirectoryNode struct {
	Path     string
	Files    []string
	Children []DirectoryNode
}

func BuildDirectoryTree(rootPath string) (DirectoryNode, error) {
	node := DirectoryNode{
		Path: rootPath,
	}
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == rootPath {
			return nil
		}
		relPath, err := filepath.Rel(rootPath, path)
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

func PrintDirectoryTree(node DirectoryNode, indent string) {
	fmt.Printf("%s%s/\n", indent, node.Path)
	for _, child := range node.Children {
		PrintDirectoryTree(child, indent+"  ")
	}

	for _, file := range node.Files {
		fmt.Printf("%s  %s\n", indent, file)
	}
}
