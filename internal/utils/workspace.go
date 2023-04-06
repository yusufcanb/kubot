package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ArchiveWorkspace(path *string) (string, error) {
	// Create a temporary file to hold the archive
	tempFile, err := os.CreateTemp("", "kubot-*.tar.gz")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	//defer os.Remove(tempFile.Name()) // Remove the temp file when we're done

	// Create a gzip writer for the temporary file
	gzipWriter := gzip.NewWriter(tempFile)
	defer gzipWriter.Close()

	// Create a tar writer for the gzip writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Walk the directory tree starting at the given path and add each file to the tar archive
	err = filepath.Walk(*path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the header for the current file
		header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
		if err != nil {
			return err
		}

		// Update the header to use the relative path within the archive
		relPath, err := filepath.Rel(filepath.Dir(*path), filePath)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		// Write the header and file contents to the tar archive
		if err = tarWriter.WriteHeader(header); err != nil {
			return err
		}
		if fileInfo.Mode().IsRegular() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(tarWriter, file); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to create archive from workspace: %v", err)
	}

	filePath, err := GetAbsolutePath(tempFile.Name())
	if err != nil {
		return "", nil
	}

	return filePath, nil
}
