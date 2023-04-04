package executor

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type K8sExecutor struct {
	Executor

	client *kubernetes.Clientset
	config *rest.Config
	pod    *corev1.Pod

	Namespace  string
	JobName    string
	JobCommand string
	JobImage   string
}

func (it *K8sExecutor) createPod() error {
	// Create a new PodSpec with the job container
	podSpec := corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:    "job-container",
				Image:   it.JobImage,
				Command: []string{"sleep", "infinity"},
			},
		},
		RestartPolicy: corev1.RestartPolicyNever,
	}

	// Create a new Pod object with the PodSpec
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: it.JobName + "-",
			Namespace:    it.Namespace,
		},
		Spec: podSpec,
	}

	// Create the Pod in the cluster
	podInterface, err := it.client.CoreV1().Pods(it.Namespace).Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	it.pod = podInterface

	return nil
}

func (it *K8sExecutor) copyTarFileToPod(pod *corev1.Pod, tarFilePath string, destinationPath string) error {
	// Construct the kubectl cp command
	cmd := exec.Command("kubectl", "cp", tarFilePath, fmt.Sprintf("%s:%s", pod.Name, destinationPath), "-n", pod.Namespace)

	// Run the command and capture the output and error streams
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to copy tar file to pod: %s. Error: %v", string(output), err)
	}

	return nil
}

func (it *K8sExecutor) createArchiveFromWorkspace(path *string) (string, error) {
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

	return tempFile.Name(), nil
}

//func (it *K8sExecutor) createArchiveFromWorkspace(path *string) (string, error) {
//	// create a new zip file
//	zipFilePath := filepath.Join(os.TempDir(), "workspace.zip")
//	zipFile, err := os.Create(zipFilePath)
//	if err != nil {
//		return "", err
//	}
//	defer zipFile.Close()
//
//	// create a new zip writer
//	zipWriter := zip.NewWriter(zipFile)
//	defer zipWriter.Close()
//
//	// walk the directory tree and add all files to the zip file
//	err = filepath.Walk(*path, func(filePath string, fileInfo os.FileInfo, err error) error {
//		if err != nil {
//			return err
//		}
//
//		// skip directories and the zip file itself
//		if fileInfo.IsDir() || filePath == zipFilePath {
//			return nil
//		}
//
//		// add the file to the zip archive
//		fileRelPath, err := filepath.Rel(*path, filePath)
//		if err != nil {
//			return err
//		}
//		zipFile, err := zipWriter.Create(fileRelPath)
//		if err != nil {
//			return err
//		}
//		file, err := os.Open(filePath)
//		if err != nil {
//			return err
//		}
//		defer file.Close()
//		_, err = io.Copy(zipFile, file)
//		if err != nil {
//			return err
//		}
//
//		return nil
//	})
//	if err != nil {
//		return "", err
//	}
//
//	return zipFilePath, nil
//}

func (it *K8sExecutor) copyWorkspaceToPod(pod *corev1.Pod, srcDir string, destDir string) error {
	filePath, err := it.createArchiveFromWorkspace(&srcDir)
	if err != nil {
		return err
	}

	// replace Windows drive letter prefix with a backslash
	if runtime.GOOS == "windows" && (strings.HasPrefix(filePath, "C:") || strings.HasPrefix(filePath, "D:")) {
		if len(filePath) < 3 {
			return errors.New("invalid path")
		}
		filePath = "\\" + strings.ToLower(filePath[0:1]) + filePath[2:]
	}

	err = it.copyTarFileToPod(pod, filePath, destDir)
	if err != nil {
		return err
	}

	return nil
}

func (it *K8sExecutor) deletePod() error {
	// Delete the Pod with the specified name in the specified namespace
	err := it.client.CoreV1().Pods(it.Namespace).Delete(context.Background(), it.pod.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (it *K8sExecutor) Configure() {
	kubeConfigPath := ""
	if home, _ := os.UserHomeDir(); home != "" {
		kubeConfigPath = fmt.Sprintf("%s/.kube/config", home)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		panic(err.Error())
	}

	it.config = config

	// create the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	it.client = client
}

func (it *K8sExecutor) Execute() error {
	it.Configure()

	err := it.createPod()
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	err = it.copyWorkspaceToPod(it.pod, ".", "/opt/robotframework/reports")
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	defer func(it *K8sExecutor) {
		err := it.deletePod()
		if err != nil {
			panic(err)
		}
	}(it)

	return nil
}
