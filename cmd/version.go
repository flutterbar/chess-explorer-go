package cmd

import (
	"github.com/flutterbar/chess-explorer-go/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display build information",
	Long:  `Display build information`,
	Run: func(cmd *cobra.Command, args []string) {
		version.DisplayVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
