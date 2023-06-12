package suite

import (
	"context"
	"errors"
	"fmt"
	"github.com/yusufcanb/kubot/pkg/cluster"
	"github.com/yusufcanb/kubot/pkg/workspace"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os/exec"
	"time"
)

type Volume struct {
	cluster *cluster.Cluster

	volume  *corev1.Volume // volume to extract workspace into
	initPod *Pod
}

func (it *Volume) Exists() bool {
	if it.volume != nil {
		return true
	}
	return false
}

func (it *Volume) Create() error {
	size := "1Gi"
	accessModes := []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
	storageClassName := "azurefile-premium"

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "pvc-kubot-",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      accessModes,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(size),
				},
			},
		},
	}

	pvc, err := it.cluster.Client().CoreV1().PersistentVolumeClaims(it.cluster.DefaultNamespace()).Create(context.Background(), pvc, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	it.volume = &corev1.Volume{
		Name: pvc.ObjectMeta.Name,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvc.ObjectMeta.Name,
			},
		},
	}

	// Wait until the volume is Bound
	for {
		pvc, err = it.cluster.Client().CoreV1().PersistentVolumeClaims(it.cluster.DefaultNamespace()).Get(context.Background(), pvc.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if pvc.Status.Phase == corev1.ClaimBound {
			break
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func (it *Volume) Destroy() error {
	if it.volume == nil {
		return nil
	}

	err := it.cluster.Client().CoreV1().PersistentVolumeClaims(it.cluster.DefaultNamespace()).Delete(context.Background(), it.volume.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	it.volume = nil

	err = it.initPod.destroy()
	if err != nil {
		return err
	}

	return nil
}

func (it *Volume) InitDirectories(w *workspace.Workspace) error {
	suitePod, err := NewSuitePod(it, "docker.io/ubuntu:bionic")
	if err != nil {
		return errors.New(fmt.Sprintf("init directories: ", err))
	}

	err = suitePod.exec([]string{"mkdir", "/data/workspace", "/data/output", "/data/console"})
	if err != nil {
		return errors.New(fmt.Sprintf("init directories: ", err))
	}

	err = suitePod.copy(w.Root().Path, "/data/workspace/")
	if err != nil {
		return errors.New(fmt.Sprintf("copy workspace: %s", err))
	}

	it.initPod = suitePod

	return nil
}

func (it *Volume) DownloadOutput() error {

	cmd := exec.Command("kubectl", "cp", fmt.Sprintf("%s:%s", it.initPod.pod.Name, "/data/output/"), ".kubot/", "-n", it.cluster.DefaultNamespace())

	// Run the command and capture the output and error streams
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to copy tar file to pod: %s. Error: %v", string(output), err)
	}

	return nil
}

func NewVolume(c *cluster.Cluster) (*Volume, error) {
	v := Volume{
		cluster: c,
	}

	err := v.Create()
	if err != nil {
		return nil, err
	}

	return &v, nil
}
