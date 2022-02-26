package cmd

import (
	server "github.com/flutterbar/chess-explorer-go/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverPort int
var startBrowser bool

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a web server to access data via a web browser",
	Long:  `Start a web server to access data via a web browser`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVar(&serverPort, "server-port", 52825, "server http port")
	serverCmd.Flags().BoolVar(&startBrowser, "start-browser", false, "automatically start a browser (default false)")

	// To be able to support the config file, we need to bind with viper (and read with viper.GetString())
	viper.BindPFlag("server-port", serverCmd.Flags().Lookup("server-port"))
	viper.BindPFlag("start-browser", serverCmd.Flags().Lookup("start-browser"))
}
