package collector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildDirectoryTree(t *testing.T) {
	// create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "example")
	if err != nil {
		t.Fatalf("error creating temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	// create some subdirectories and files within the temporary directory
	err = os.Mkdir(filepath.Join(tempDir, "dir1"), 0777)
	if err != nil {
		t.Fatalf("error creating directory: %s", err)
	}
	err = os.Mkdir(filepath.Join(tempDir, "dir2"), 0777)
	if err != nil {
		t.Fatalf("error creating directory: %s", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("file1"), 0644)
	if err != nil {
		t.Fatalf("error creating file: %s", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "dir1", "file2.txt"), []byte("file2"), 0644)
	if err != nil {
		t.Fatalf("error creating file: %s", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "dir2", "file3.txt"), []byte("file3"), 0644)
	if err != nil {
		t.Fatalf("error creating file: %s", err)
	}

	// test the BuildDirectoryTree function
	tree, err := BuildDirectoryTree(tempDir)
	if err != nil {
		t.Fatalf("error building directory tree: %s", err)
	}

	// check that the directory tree includes the expected nodes and files
	expectedTree := DirectoryNode{
		Path: tempDir,
		Children: []DirectoryNode{
			{Path: "dir1", Files: []string{"file2.txt"}},
			{Path: "dir2", Files: []string{"file3.txt"}},
		},
		Files: []string{"file1.txt"},
	}

	if !compareDirectoryNodes(tree, expectedTree) {
		t.Errorf("unexpected directory tree:\n\nexpected:\n%+v\n\nactual:\n%+v", expectedTree, tree)
	}
}

// helper function for comparing DirectoryNodes
func compareDirectoryNodes(node1, node2 DirectoryNode) bool {
	if node1.Path != node2.Path {
		return false
	}
	if len(node1.Children) != len(node2.Children) {
		return false
	}
	if len(node1.Files) != len(node2.Files) {
		return false
	}
	for i := range node1.Children {
		if !compareDirectoryNodes(node1.Children[i], node2.Children[i]) {
			return false
		}
	}
	for i := range node1.Files {
		if node1.Files[i] != node2.Files[i] {
			return false
		}
	}
	return true
}
