package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"kubot/pkg/executor"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a job",
	Run: func(cmd *cobra.Command, args []string) {
		executor := executor.K8sExecutor{}
		executor.Namespace = "roc"
		executor.JobName = "kubot-4"
		executor.JobImage = "yceiotc.azurecr.io/roc/roc-runner:k8s.2"
		executor.JobCommand = "robot /opt/robotframework/*.robot"
		err := executor.Execute()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
