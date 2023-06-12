package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yusufcanb/kubot/pkg/app"
	"os"
	"path/filepath"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a job",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil || name == "" {
			log.Fatalf("Error getting name flag: %s", err)
		}

		image, err := cmd.Flags().GetString("image")
		if err != nil || image == "" {
			log.Fatalf("Error getting image flag: %s", err)
		}

		batchSize, err := cmd.Flags().GetInt("batchsize")
		if err != nil || batchSize == 0 {
			log.Fatalf("Error getting batchsize flag: %s", err)
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
			selector = workspace
		}

		k, err := app.New(app.RuntimeArgs{
			TopLevelSuiteName: name,
			Namespace:         namespace,
			Image:             image,
			WorkspacePath:     workspace,
			Selector:          selector,
			BatchSize:         batchSize,
		})

		if err != nil {
			log.Fatal(err)
		}

		err = k.Run()
		if err != nil {
			log.Fatal(err)
		}

		defer k.Clean()
	},
}

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	execCmd.Flags().StringP("workspace", "w", filepath.Dir(ex), "workspace path")
	execCmd.Flags().StringP("name", "n", "Kubot Results", "top level suite name for logs and reports")
	execCmd.Flags().StringP("namespace", "", "", "kubernetes namespace to create workloads in it")
	execCmd.Flags().StringP("image", "i", "", "docker image for execution for pods and jobs")
	execCmd.Flags().IntP("batchsize", "b", 25, "execution batch size")
	execCmd.Flags().StringP("selector", "s", "", "script selector. e.g. tasks/*")

	rootCmd.AddCommand(execCmd)
}
