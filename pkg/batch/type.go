package batch

import "kubot/pkg/workspace"

type Batch struct {
	size   int
	cursor int

	workspace *workspace.Workspace
}

func (b *Batch) Next() []string {
	allFiles := b.workspace.Root().Files
	if b.cursor >= len(allFiles) {
		return nil
	}

	var batchSize int
	if b.cursor+b.size > len(allFiles) {
		batchSize = len(allFiles) - b.cursor
	} else {
		batchSize = b.size
	}

	batch := allFiles[b.cursor : b.cursor+batchSize]
	b.cursor += batchSize

	return batch
}

func NewBatch(size int, w *workspace.Workspace) *Batch {
	return &Batch{
		size:      size,
		cursor:    0,
		workspace: w,
	}
}
