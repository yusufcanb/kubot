package suite

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/yusufcanb/kubot/pkg/batch"
	"github.com/yusufcanb/kubot/pkg/cluster"
	"github.com/yusufcanb/kubot/pkg/workspace"
	"sync"
	"time"
)

type Runner struct {
	cluster *cluster.Cluster

	merger *Merger
	image  string

	startedAt   time.Time
	completedAt time.Time
}

func (it *Runner) executeSuite(v *Volume, suiteName string) error {
	suitePod, err := NewSuitePod(v, it.image)
	if err != nil {
		return err
	}

	err = suitePod.exec([]string{
		"robot", "--log", "NONE", "--report", "NONE", "--outputdir", fmt.Sprintf("/data/output/%s", suiteName), fmt.Sprintf("/data/workspace/scripts/%s", suiteName),
	})
	defer suitePod.destroy()
	if err != nil {
		log.Errorf("robot script failed: %s", err)
		return err
	}

	return nil
}

func (it *Runner) Run(w *workspace.Workspace, v *Volume, batchSize int) error {

	scriptBatch := batch.NewBatch(batchSize, w)

	it.startedAt = time.Now()

	for {
		items := scriptBatch.Next()
		if items == nil {
			break
		}
		var wg sync.WaitGroup

		for _, file := range items {
			wg.Add(1)
			go func(filename string) {
				_ = it.executeSuite(v, filename)
				defer wg.Done()
			}(file)
		} // execute every script in the batch
		wg.Wait()
	} // process batches

	it.completedAt = time.Now()
	time.Sleep(5 * time.Second) // Wait for all the buffers to be completed.

	err := it.merger.MergeResults(v, it.image, &it.startedAt, &it.completedAt)
	defer v.DownloadOutput()
	if err != nil {
		log.Errorf("merging failed: %s", err)
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func NewRunner(c *cluster.Cluster, image string, topLevelSuiteName string) *Runner {
	return &Runner{
		cluster: c,
		image:   image,
		merger: &Merger{
			topLevelSuiteName: topLevelSuiteName,
		},
	}
}
