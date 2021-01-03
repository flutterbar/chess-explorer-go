package cmd

import (
	"github.com/spf13/cobra"
	server "github.com/yafred/chess-explorer/internal/server"
)

var serverPort int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a web server",
	Long:  `Start a web server`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Start(serverPort)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVar(&serverPort, "server-port", 52825, "Server http port")
}
