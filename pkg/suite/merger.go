package suite

type Merger struct {
}

func (it Merger) MergeResults(v *Volume, image string) error {
	suitePod, err := NewSuitePod(v, image)
	if err != nil {
		return err
	}

	cmd := []string{"rebot", "--outputdir", "/data/output", "/data/output/*/output.xml"}

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
