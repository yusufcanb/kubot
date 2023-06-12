package app

import (
	"github.com/yusufcanb/kubot/pkg/cluster"
	"github.com/yusufcanb/kubot/pkg/suite"
	"github.com/yusufcanb/kubot/pkg/workspace"
)

func New(args RuntimeArgs) (*App, error) {
	var err error
	var app = App{}

	app.topLevelSuiteName = args.TopLevelSuiteName
	app.batchSize = args.BatchSize

	app.cluster, err = cluster.NewCluster("", args.Namespace)
	if err != nil {
		return nil, err
	}

	app.workspace, err = workspace.New(args.WorkspacePath)
	if err != nil {
		return nil, err
	}

	app.suiteVolume, err = suite.NewVolume(app.cluster)
	if err != nil {
		return nil, err
	}

	err = app.suiteVolume.InitDirectories(app.workspace)
	if err != nil {
		return nil, err
	}

	app.suiteRunner = suite.NewRunner(app.cluster, args.Image, app.topLevelSuiteName)

	return &app, nil
}
