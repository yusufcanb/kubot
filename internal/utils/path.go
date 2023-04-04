package utils

import (
	"errors"
	"runtime"
	"strings"
)

func GetAbsolutePath(filePath string) (string, error) {
	// replace Windows drive letter prefix with a backslash
	if runtime.GOOS == "windows" && (strings.HasPrefix(filePath, "C:") || strings.HasPrefix(filePath, "D:")) {
		if len(filePath) < 3 {
			return "", errors.New("invalid path")
		}
		filePath = strings.Replace(filePath, "C:", "", -1)
	}
	return filePath, nil
}
