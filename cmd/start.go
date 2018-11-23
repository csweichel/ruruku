package cmd

import (
	"fmt"
	"time"

	"github.com/32leaves/ruruku/pkg/rurukulib"
	"github.com/32leaves/ruruku/pkg/rurukulib/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sessionName string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <testsuite.yaml>",
	Short: "Starts a test session run",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := GetConfigFromViper()
		if err != nil {
			log.Fatalf("Error while loading the configuration: %s", err)
		}

		suite, err := rurukulib.LoadSuite(args[0])
		if err != nil {
			log.Fatalf("Error while loading test suite: %s", err)
		}

		srvcfg := cfg.Server
		if err := server.Start(&srvcfg, suite, sessionName); err != nil {
			log.Fatalf("Error while starting the ruruku server", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	sn := time.Now().Format("20060102-150405")
	startCmd.Flags().StringVar(&sessionName, "session-name", fmt.Sprintf("ruruku-%s.yaml", sn), "The name of the session to run")

	startCmd.Flags().Int("port", 8080, "Port to run the server on")
	viper.BindPFlag("server.port", startCmd.Flags().Lookup("port"))

	startCmd.Flags().String("token", "", "The authentication token users need to know (default is auto-generated)")
	viper.BindPFlag("server.token", startCmd.Flags().Lookup("token"))
}
