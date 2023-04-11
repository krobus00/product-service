package cmd

import (
	"github.com/krobus00/product-service/internal/bootstrap"
	"github.com/spf13/cobra"
)

// workerCmd represents the worker command.
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "product service worker",
	Long:  `product service worker`,
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.StartWorker()
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
