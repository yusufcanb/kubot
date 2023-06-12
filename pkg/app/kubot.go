package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/yusufcanb/kubot/pkg/cluster"
	"github.com/yusufcanb/kubot/pkg/suite"
	"github.com/yusufcanb/kubot/pkg/workspace"
)

type App struct {
	cluster *cluster.Cluster

	workspace *workspace.Workspace

	suiteVolume *suite.Volume // volume to extract workspace into
	suiteRunner *suite.Runner

	topLevelSuiteName string
	batchSize         int
}

func (it *App) Run() error {
	err := it.suiteRunner.Run(it.workspace, it.suiteVolume, it.batchSize)
	if err != nil {
		return err
	}

	return nil
}

func (it *App) Clean() {
	err := it.suiteVolume.Destroy()
	if err != nil {
		log.Fatal(err)
	}
}
