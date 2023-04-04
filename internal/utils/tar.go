package utils

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//func CreateTarArchive(srcDir string) ([]byte, error) {
//	tarBuf := new(bytes.Buffer)
//	tarWriter := tar.NewWriter(tarBuf)
//	defer tarWriter.Close()
//
//	_, err := os.Stat(srcDir)
//	if err != nil {
//		return nil, err
//	}
//
//	err = filepath.Walk(srcDir, func(path string, fileInfo os.FileInfo, err error) error {
//		if err != nil {
//			return err
//		}
//		if fileInfo.Mode().IsRegular() {
//			relPath, err := filepath.Rel(srcDir, path)
//			if err != nil {
//				return err
//			}
//			header := &tar.Header{
//				Name:    relPath,
//				Mode:    int64(fileInfo.Mode()),
//				ModTime: fileInfo.ModTime(),
//				Size:    fileInfo.Size(),
//			}
//			if err := tarWriter.WriteHeader(header); err != nil {
//				return err
//			}
//			file, err := os.Open(path)
//			if err != nil {
//				return err
//			}
//			defer file.Close()
//			if _, err := io.Copy(tarWriter, file); err != nil {
//				return err
//			}
//		}
//		return nil
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	return tarBuf.Bytes(), nil
//}

func CreateTar(path string) ([]byte, error) {
	// Create a buffer to store the tar file
	buf := new(bytes.Buffer)

	// Create a new tar writer
	tw := tar.NewWriter(buf)

	// Open the source folder
	src, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open source folder: %v", err)
	}
	defer src.Close()

	// Get the list of files in the folder
	files, err := src.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of files in folder: %v", err)
	}

	// Add each file to the tar archive
	for _, fileInfo := range files {
		// Open the file
		file, err := os.Open(filepath.Join(path, fileInfo.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %v", fileInfo.Name(), err)
		}
		defer file.Close()

		// Create a new tar header for the file
		header := new(tar.Header)
		header.Name = fileInfo.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		// Write the header to the tar archive
		if err := tw.WriteHeader(header); err != nil {
			return nil, fmt.Errorf("failed to write tar header for file %s: %v", fileInfo.Name(), err)
		}

		// Write the file content to the tar archive
		if _, err := io.Copy(tw, file); err != nil {
			return nil, fmt.Errorf("failed to write content for file %s: %v", fileInfo.Name(), err)
		}
	}

	// Close the tar writer
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %v", err)
	}

	return buf.Bytes(), nil
}
