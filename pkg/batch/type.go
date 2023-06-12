package batch

import "github.com/yusufcanb/kubot/pkg/workspace"

// Batch struct represents a batch of files to be processed.
// It contains a size field which represents the size of each batch,
// and a cursor field which represents the current position in the file list.
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
