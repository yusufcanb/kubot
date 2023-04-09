package utils

import "io/ioutil"

func SaveBufferToFile(buffer []byte, filepath string) error {
	// Write the buffer to a file
	err := ioutil.WriteFile(filepath, buffer, 0644)
	if err != nil {
		return err
	}
	return nil
}
