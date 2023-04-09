package suite

import "time"

type Merger struct {
}

func (it Merger) MergeResults(v *Volume, image string, startedAt *time.Time, completedAt *time.Time) error {
	suitePod, err := NewSuitePod(v, image)
	if err != nil {
		return err
	}

	cmd := []string{
		"rebot",
		"--starttime", startedAt.UTC().String(),
		"--endtime", completedAt.UTC().String(),
		"--outputdir", "/data/output", "/data/output/*/output.xml",
	}

	err = suitePod.exec(cmd)
	if err != nil {
		return err
	}

	err = suitePod.destroy()
	if err != nil {
		return err
	}

	return nil
}
