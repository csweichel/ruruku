package cmd

import (
    "os"
    "os/signal"
	"github.com/32leaves/ruruku/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sessionName string

// serveCmd represents the start command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts a ruruku API server",

	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := GetConfigFromViper()
		if err != nil {
			log.Fatalf("Error while loading the configuration: %s", err)
		}

		srvcfg := cfg.Server
		if err := server.Start(&srvcfg, server.NewMemoryBackedSession()); err != nil {
			log.Fatalf("Error while starting the ruruku server", err)
		}

        signal_channel := make(chan os.Signal, 1)
        signal.Notify(signal_channel, os.Interrupt)
        <-signal_channel
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Int("ui-port", 8080, "Port to run UI the server on")
	viper.BindPFlag("server.ui.port", serveCmd.Flags().Lookup("ui-port"))
}
