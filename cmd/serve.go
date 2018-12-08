package cmd

import (
	"github.com/32leaves/ruruku/pkg/server"
	"github.com/32leaves/ruruku/pkg/server/kvsession"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
)

var sessionName string

// serveCmd represents the start command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts a ruruku API server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := GetConfigFromViper()
		if err != nil {
			log.Fatalf("Error while loading the configuration: %v", err)
		}

		srvcfg := cfg.Server
		store, err := kvsession.NewSession(srvcfg.DB.Filename)
		if err != nil {
			log.Fatalf("Error while creating the session store: %v", err)
		}
		log.WithField("filename", srvcfg.DB.Filename).Info("Opened database")

		if err := server.Start(&srvcfg, store); err != nil {
			log.Fatalf("Error while starting the ruruku server: %v", err)
		}

		if srvcfg.UI.Enabled {
			log.Infof("Started server at %s", serverUrl(srvcfg.UI.Port))
		}

		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt)
		<-signalChannel
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Int("ui-port", 8080, "Port to run UI the server on")
	viper.BindPFlag("server.ui.port", serveCmd.Flags().Lookup("ui-port"))
	serveCmd.Flags().String("db", "ruruku.db", "Path to the data storage location")
	viper.BindPFlag("server.DB.Filename", serveCmd.Flags().Lookup("db"))
	viper.BindEnv("server.DB.Filename", "RURUKU_DB")
}
