package suite

import (
	"bytes"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"kubot/pkg/cluster"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Pod struct {
	cluster *cluster.Cluster
	pod     *corev1.Pod

	deleted bool
}

// waitUntilPodHasStarted to ensure the executor's pod in Running state
func (it *Pod) waitUntilPodHasStarted() error {
	// define timeout and polling interval
	timeoutSeconds := 300
	pollIntervalSeconds := 5

	// create a context with timeout and cancel functions
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// poll the pod until it is in the `Running` state
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for pod %s/%s to start", it.pod.Namespace, it.pod.Name)
		default:
			pod, err := it.cluster.Client().CoreV1().Pods(it.pod.Namespace).Get(context.Background(), it.pod.Name, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("error getting pod %s/%s: %v", pod.Namespace, pod.Name, err)
			}
			if pod.Status.Phase == corev1.PodRunning {
				return nil
			}
			time.Sleep(time.Duration(pollIntervalSeconds) * time.Second)
		}
	}
}

// collectEnvironmentVariablesFromOs
func (it *Pod) collectEnvironmentVariablesFromOs() []corev1.EnvVar {
	var envVars []corev1.EnvVar

	envKeyRegex := regexp.MustCompile(`^[-._a-zA-Z][-._a-zA-Z0-9]*$`)

	for _, env := range os.Environ() {
		keyValue := strings.SplitN(env, "=", 2)
		if len(keyValue) == 2 {
			key := keyValue[0]
			value := keyValue[1]

			if envKeyRegex.MatchString(key) && envKeyRegex.MatchString(value) {
				envVar := corev1.EnvVar{
					Name:  key,
					Value: value,
				}
				envVars = append(envVars, envVar)
			}
		}
	}

	return envVars
}

// copy given file into the pod
func (it *Pod) copy(srcPath string, destinationPath string) error {
	fmt.Printf("%s >>> %s\n", it.pod.Name, []string{"kubectl", "cp", srcPath, fmt.Sprintf("%s:%s", it.pod.Name, destinationPath), "-n", it.pod.Namespace})

	// Construct the kubectl cp command
	cmd := exec.Command("kubectl", "cp", srcPath, fmt.Sprintf("%s:%s", it.pod.Name, destinationPath), "-n", it.pod.Namespace)

	// Run the command and capture the output and error streams
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to copy tar file to pod: %s. Error: %v", string(output), err)
	}

	return nil
}

// exec executes given command inside the pod
func (it *Pod) exec(cmd []string) error {
	fmt.Printf("%s >>> %s\n", it.pod.Name, cmd)

	buf := &bytes.Buffer{}
	request := it.cluster.Client().CoreV1().RESTClient().
		Post().
		Namespace(it.pod.Namespace).
		Resource("pods").
		Name(it.pod.Name).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: cmd,
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	spdyExec, err := remotecommand.NewSPDYExecutor(it.cluster.Config(), "POST", request.URL())
	err = spdyExec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: buf,
		Stderr: buf,
		Tty:    true,
	})

	if err != nil {
		fmt.Println(buf.String())
		return fmt.Errorf("%w Failed executing command %s on %v/%v", err, cmd, it.pod.Namespace, it.pod.Name)
	}
	return nil
}

func (it *Pod) destroy() error {
	// Delete the Pod with the specified name in the specified namespace
	err := it.cluster.Client().CoreV1().Pods(it.cluster.DefaultNamespace()).Delete(context.Background(), it.pod.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	it.pod = nil
	it.deleted = true

	return nil
}

func NewSuitePod(suiteVolume *Volume, image string) (*Pod, error) {
	suitePod := Pod{}
	suitePod.cluster = suiteVolume.cluster

	// Create a new PodSpec with the job container
	podSpec := corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:    "job-container",
				Image:   image,
				Command: []string{"sleep", "infinity"},
				Env:     suitePod.collectEnvironmentVariablesFromOs(),
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      suiteVolume.volume.Name,
						MountPath: "/data",
					},
				},
			},
		},
		Volumes: []corev1.Volume{
			{
				Name: suiteVolume.volume.Name,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: suiteVolume.volume.Name,
					},
				},
			},
		},
		RestartPolicy: corev1.RestartPolicyNever,
	}

	// Create a new Pod object with the PodSpec
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "kubot-",
			Namespace:    suitePod.cluster.DefaultNamespace(),
		},
		Spec: podSpec,
	}

	// Create the Pod in the cluster
	podInterface, err := suitePod.cluster.Client().CoreV1().Pods(suitePod.cluster.DefaultNamespace()).Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	suitePod.pod = podInterface
	err = suitePod.waitUntilPodHasStarted()
	if err != nil {
		return nil, err
	}

	return &suitePod, nil
}
