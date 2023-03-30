package cmd

import (
	"fmt"
	"os"

	"github.com/krobus00/product-service/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "product-service",
	Short: "product service",
	Long:  `product service`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Init(name, version string) {
	log.Info(fmt.Sprintf("starting %s:%s ...", name, version))
	if err := config.LoadConfig(); err != nil {
		log.Fatalln(err.Error())
	}
}
