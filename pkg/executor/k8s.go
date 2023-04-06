package executor

import (
	"bytes"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"kubot/internal/utils"
	"os"
	"os/exec"
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

// createPod creates new Pod workload using kubernetes client
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

// waitUntilPodHasStarted to ensure the executor's pod in Running state
func (it *K8sExecutor) waitUntilPodHasStarted(pod *corev1.Pod) error {
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
			return fmt.Errorf("timed out waiting for pod %s/%s to start", pod.Namespace, pod.Name)
		default:
			pod, err := it.client.CoreV1().Pods(pod.Namespace).Get(context.Background(), pod.Name, metav1.GetOptions{})
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

// copy given file into the pod
func (it *K8sExecutor) copy(pod *corev1.Pod, srcPath string, destinationPath string) error {
	// Construct the kubectl cp command
	cmd := exec.Command("kubectl", "cp", srcPath, fmt.Sprintf("%s:%s", pod.Name, destinationPath), "-n", pod.Namespace)

	// Run the command and capture the output and error streams
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to copy tar file to pod: %s. Error: %v", string(output), err)
	}

	return nil
}

// exec executes given command inside the pod
func (it *K8sExecutor) exec(pod *corev1.Pod, cmd []string) error {
	buf := &bytes.Buffer{}
	request := it.client.CoreV1().RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: cmd,
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)
	spdyExec, err := remotecommand.NewSPDYExecutor(it.config, "POST", request.URL())
	err = spdyExec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: buf,
		Stderr: buf,
		Tty:    true,
	})
	if err != nil {
		return fmt.Errorf("%w Failed executing command %s on %v/%v", err, cmd, pod.Namespace, pod.Name)
	}

	fmt.Println(buf.String())

	return nil
}

// deletePod deletes the executor's pod
func (it *K8sExecutor) deletePod() error {
	// Delete the Pod with the specified name in the specified namespace
	err := it.client.CoreV1().Pods(it.Namespace).Delete(context.Background(), it.pod.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Configure the executor
func (it *K8sExecutor) Configure(any) error {
	kubeConfigPath := ""
	if home, _ := os.UserHomeDir(); home != "" {
		kubeConfigPath = fmt.Sprintf("%s/.kube/config", home)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return err
	}

	it.config = config

	// create the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	it.client = client

	return nil
}

// Execute the executor
func (it *K8sExecutor) Execute(workspacePath string) error {
	err := it.createPod()
	if err != nil {
		return err
	}

	err = it.waitUntilPodHasStarted(it.pod)
	if err != nil {
		return err
	}

	log.Infof("copying workspace into pod %s", it.pod.Name)
	archivePath, err := utils.ArchiveWorkspace(&workspacePath)
	if err != nil {
		return err
	}

	err = it.copy(it.pod, archivePath, "/opt/robotframework/reports")
	if err != nil {
		return err
	}

	cmd := []string{"ls", "/", "-all"}
	log.Infof("executing command %s in pod %s", cmd, it.pod.Name)
	err = it.exec(it.pod, cmd)
	if err != nil {
		log.Fatal(err)
	}

	defer func(it *K8sExecutor) {
		log.Infof("removing pod %s", it.pod.Name)
		err := it.deletePod()
		if err != nil {
			panic(err)
		}
	}(it)

	return nil
}
