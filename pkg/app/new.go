package app

import (
	"kubot/pkg/cluster"
	"kubot/pkg/suite"
	"kubot/pkg/workspace"
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
