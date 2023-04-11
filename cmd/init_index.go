package cmd

import (
	"github.com/krobus00/product-service/internal/bootstrap"
	"github.com/spf13/cobra"
)

// initIndexCmd represents the initIndex command.
var initIndexCmd = &cobra.Command{
	Use:   "init-index",
	Short: "init opensearch index",
	Long:  `init opensearch index`,
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.StartInitIndex()
	},
}

func init() {
	rootCmd.AddCommand(initIndexCmd)
}
