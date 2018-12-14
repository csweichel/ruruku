package cmd

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"

	"github.com/32leaves/ruruku/pkg/server"
	"github.com/32leaves/ruruku/pkg/server/kvsession"
	"github.com/32leaves/ruruku/pkg/server/kvuser"
	bolt "github.com/etcd-io/bbolt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <testplan.yaml>",
	Short: "Creates a ruruku API server and starts a session",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := GetConfigFromViper()
		if err != nil {
			log.Fatalf("Error while loading the configuration: %v", err)
		}

		srvcfg := cfg.Server
		srvcfg.UI.Enabled = true

		db, err := bolt.Open(srvcfg.DB.Filename, 0666, nil)
		if err != nil {
			log.Fatalf("Error while opening database: %v", err)
		}
		log.WithField("filename", srvcfg.DB.Filename).Info("Opened database")

		userStore, err := kvuser.NewUserStore(db)
		if err != nil {
			log.Fatalf("Error while creating the user store: %v", err)
		}

		store, err := kvsession.NewSession(db, userStore)
		if err != nil {
			log.Fatalf("Error while creating the session store: %v", err)
		}
		log.WithField("filename", srvcfg.DB.Filename).Info("Opened database")

		if err := server.Start(&srvcfg, store, nil); err != nil {
			log.Fatalf("Error while starting the ruruku server: %v", err)
		}

		sessionStart := sessionStartFlags{planfn: args[0], quiet: true}
		if err := sessionStart.Run(); err != nil {
			log.WithError(err).Fatal("Cannot start session")
		}

		log.Infof("Started server at %s", serverUrl(srvcfg.UI.Port))

		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt)
		<-signalChannel
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func serverUrl(port uint32) string {
	protocol := "https"
	host := "localhost"
	wsURL := os.Getenv("GITPOD_WORKSPACE_URL")
	if wsURL != "" {
		parsedWsURL, err := url.Parse(wsURL)
		if err == nil {
			host = fmt.Sprintf("%d-%s", port, parsedWsURL.Host)
			protocol = parsedWsURL.Scheme
		}
	}
	return fmt.Sprintf("%s://%s", protocol, host)
}
