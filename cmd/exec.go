package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"kubot/pkg/executor"
	"os"
	"path/filepath"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a job",
	Run: func(cmd *cobra.Command, args []string) {
		executor := executor.K8sExecutor{}
		err := executor.Configure(nil)
		if err != nil {
			log.Fatal(err)
		}

		workspace, err := cmd.Flags().GetString("workspace")
		if err != nil {
			log.Fatalf("Error getting file flag: %s", err)
		}
		executor.Namespace = "kubot"
		executor.JobName = "kubot-4"
		executor.JobImage = "yceiotc.azurecr.io/roc/roc-runner:k8s.2"
		executor.JobCommand = "robot /opt/robotframework/*.robot"

		err = executor.Execute(workspace)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execCmd.Flags().StringP("workspace", "w", filepath.Dir(ex), "workspace path")

	rootCmd.AddCommand(execCmd)
}
