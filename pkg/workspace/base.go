package workspace

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

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

func (it *Workspace) SubDirectoryNames() []string {
	subdirnames := make([]string, 0)
	for _, dir := range it.Root().Children {
		subdirnames = append(subdirnames, dir.Path)
	}

	return subdirnames
}

func New(basePath string) (*Workspace, error) {

	var w Workspace
	var err error

	w = Workspace{}
	w.root, err = buildDirectoryTree(basePath)
	if err != nil {
		return nil, err
	}

	if len(w.Root().Children) > 0 {
		log.Warnf("Sub-directories [%s] will be discarded. Kubot only supports root directories.", strings.Join(w.SubDirectoryNames()[:], ","))
	}

	return &w, nil
}
