package workspace

type DirectoryNode struct {
	Path     string
	Files    []string
	Children []DirectoryNode
}

type Workspace struct {
	root DirectoryNode
}

func (it *Workspace) Root() DirectoryNode {
	return it.root
}

func New(basePath string) (*Workspace, error) {

	var w Workspace
	var err error

	w = Workspace{}
	w.root, err = buildDirectoryTree(basePath)
	if err != nil {
		return nil, err
	}

	return &w, nil
}
