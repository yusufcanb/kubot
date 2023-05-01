package suite

import (
	"fmt"
	"time"
)

type Merger struct {
	topLevelSuiteName string
}

func (it *Merger) MergeResults(v *Volume, image string, startedAt *time.Time, completedAt *time.Time) error {

	suitePod, err := NewSuitePod(v, image)
	if err != nil {
		return err
	}

	cmd := []string{
		"rebot",
		"--name", it.topLevelSuiteName,
		"--starttime", startedAt.UTC().String(),
		"--endtime", completedAt.UTC().String(),
		"--outputdir", "/data/output", "/data/output/*/output.xml",
	}

	err = suitePod.exec(cmd)
	defer suitePod.destroy()
	if err != nil {
		fmt.Errorf("merger failed: %s", err)
		return err
	}

	return nil
}
