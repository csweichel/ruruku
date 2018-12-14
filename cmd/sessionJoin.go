package cmd

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sessionJoinCmd represents the sessionJoin command
var sessionJoinCmd = &cobra.Command{
	Use:   "join <session-id>",
	Short: "Joins a testing session",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := GetConfigFromViper()
		if err != nil {
			log.Fatalf("Error while loading the configuration: %v", err)
		}

		conn, err := cfg.Connect()
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := api.NewSessionServiceClient(conn)

		ctx, cancel := cfg.GetContext(true)
		defer cancel()

		_, err = client.Register(ctx, &api.RegistrationRequest{Session: args[0]})
		if err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	sessionCmd.AddCommand(sessionJoinCmd)
}
