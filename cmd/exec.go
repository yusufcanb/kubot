package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"kubot/pkg/app"
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
		if err != nil || image == "" {
			log.Fatalf("Error getting image flag: %s", err)
		}

		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil || namespace == "" {
			log.Fatalf("Error getting namespace flag: %s", err)
		}

		workspace, err := cmd.Flags().GetString("workspace")
		if err != nil || workspace == "" {
			log.Fatalf("Error getting workspace flag: %s", err)
		}

		selector, err := cmd.Flags().GetString("selector")
		if err != nil || selector == "" {
			log.Fatalf("Error getting selector flag: %s", err)
		}

		k, err := app.New(app.RuntimeArgs{
			Namespace:     namespace,
			Image:         image,
			WorkspacePath: workspace,
			Selector:      selector,
		})

		if err != nil {
			log.Fatal(err)
		}

		k.Run()

		defer k.Clean()
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
