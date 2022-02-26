package cmd

import (
	pgntodb "github.com/flutterbar/chess-explorer-go/internal/pgntodb"
	"github.com/spf13/cobra"
)

var username string

var pgnToDbCmd = &cobra.Command{
	Use:   "pgntodb [pgn file]",
	Short: "Parse a pgn file and feed mongo database",
	Long:  `Parse a pgn file and feed mongo database. Designed for chess.com and lichess.org`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lastGame := pgntodb.LastGame{Username: username}
		pgntodb.Process(args[0], &lastGame)
	},
}

func init() {
	rootCmd.AddCommand(pgnToDbCmd)

	pgnToDbCmd.Flags().StringVar(&username, "username", "", "username for whom you are downloading games")

}
