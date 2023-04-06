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

		image, err := cmd.Flags().GetString("image")
		executor.JobImage = image
		if err != nil {
			log.Fatalf("Error getting image flag: %s", err)
		}

		namespace, err := cmd.Flags().GetString("namespace")
		executor.Namespace = namespace
		if err != nil {
			log.Fatalf("Error getting file flag: %s", err)
		}

		workspace, err := cmd.Flags().GetString("workspace")
		if err != nil {
			log.Fatalf("Error getting file flag: %s", err)
		}

		executor.JobName = "kubot-4"
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
	execCmd.Flags().StringP("namespace", "n", "", "kubernetes namespace")
	execCmd.Flags().StringP("image", "i", "", "docker image for execution")
	execCmd.Flags().StringP("selector", "s", "", "script selector. e.g. tasks/*")

	rootCmd.AddCommand(execCmd)
}
