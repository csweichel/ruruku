package cmd

import (
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/32leaves/ruruku/pkg/server"
	"github.com/32leaves/ruruku/pkg/server/kvsession"
	"github.com/32leaves/ruruku/pkg/server/kvuser"
	bolt "github.com/etcd-io/bbolt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootTokenFile string
var rootTokenStdout bool

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
		db, err := bolt.Open(cfg.Server.DB.Filename, 0666, nil)
		if err != nil {
			log.Fatalf("Error while opening database: %v", err)
		}
		log.WithField("filename", srvcfg.DB.Filename).Info("Opened database")

		userStore, err := kvuser.NewUserStore(db)
		if err != nil {
			log.Fatalf("Error while creating the user store: %v", err)
		}

		sessionStore, err := kvsession.NewSession(db, userStore)
		if err != nil {
			log.Fatalf("Error while creating the session store: %v", err)
		}

		rootTkn, err := userStore.GetUserToken("root")
		if err != nil {
			log.Fatalf("Cannot get root user token: %v", err)
		}
		if rootTokenFile != "" {
			if err := ioutil.WriteFile(rootTokenFile, []byte(rootTkn), 0600); err != nil {
				log.Fatalf("Unable to write root token file: %v", err)
			}
			log.WithField("filename", rootTokenFile).Info("Wrote root token to file")
		} else if rootTokenStdout {
			log.WithTime(time.Now()).WithField("token", rootTkn).WithField("lifetime", userStore.TokenLifetime).Infof("Root user token ... keep this token safe!")
		}

		if err := server.Start(&srvcfg, sessionStore, userStore); err != nil {
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

	serveCmd.Flags().StringVar(&rootTokenFile, "root-token-file", "", "Write a root token to a file")
	serveCmd.Flags().BoolVar(&rootTokenStdout, "root-token-stdout", false, "Print a root token to stdout")

	serveCmd.Flags().Int("ui-port", 8080, "Port to run UI the server on")
	viper.BindPFlag("server.ui.port", serveCmd.Flags().Lookup("ui-port"))
	serveCmd.Flags().String("db", "ruruku.db", "Path to the data storage location")
	viper.BindPFlag("server.DB.Filename", serveCmd.Flags().Lookup("db"))
	viper.BindEnv("server.DB.Filename", "RURUKU_DB")
}
