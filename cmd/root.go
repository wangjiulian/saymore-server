package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	Version = "v0.1.0"
	Date    = "2022-04-17"
)

var rootCmd = &cobra.Command{
	Use:   "say-more-server",
	Short: "Say More Server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
