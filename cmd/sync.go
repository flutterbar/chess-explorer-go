package cmd

import (
	"github.com/flutterbar/chess-explorer-go/internal/sync"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Download recent games for all users in database",
	Long:  `Download recent games for all users in database`,
	Run: func(cmd *cobra.Command, args []string) {
		sync.All()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
