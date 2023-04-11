package suite

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"kubot/pkg/cluster"
	"kubot/pkg/workspace"
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
	if err != nil {
		log.Warningf("robot script failed: %s", err)
		_ = suitePod.destroy()
		return err
	}

	err = suitePod.destroy()
	if err != nil {
		return err
	}

	return nil
}

func (it *Runner) Run(w *workspace.Workspace, v *Volume) error {
	it.startedAt = time.Now()

	var wg sync.WaitGroup

	for _, file := range w.Root().Files {
		wg.Add(1)
		go func(filename string) {
			err := it.executeSuite(v, filename)

			if err != nil {
				log.Warning(err)
			}
			defer wg.Done()
		}(file)
	}
	wg.Wait()

	it.completedAt = time.Now()

	err := it.merger.MergeResults(v, it.image, &it.startedAt, &it.completedAt)
	defer v.DownloadOutput()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func NewRunner(c *cluster.Cluster, image string) *Runner {
	return &Runner{
		cluster: c,
		image:   image,
		merger:  &Merger{},
	}
}
