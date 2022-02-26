package cmd

import (
	"github.com/flutterbar/chess-explorer-go/internal/delete"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [user]",
	Short: "Delete user in database",
	Long: `Delete user in database ...
Username can have 3 forms:
- username
- lichess.org:username
- chess.com:username`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		delete.Games(args[0])
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
