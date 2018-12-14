package cmd

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// userAuthCmd represents the sessionClose command
var userDeleteCmd = &cobra.Command{
	Use:   "delete <username>",
	Short: "Deletes a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := GetConfigFromViper()
		if err != nil {
			log.Fatalf("Error while loading the configuration: %v", err)
		}

		username := args[0]

		conn, err := cfg.Connect()
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := api.NewUserServiceClient(conn)

		ctx, cancel := cfg.GetContext(true)
		defer cancel()

		if _, err := client.Delete(ctx, &api.DeleteUserRequest{Username: username}); err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	userCmd.AddCommand(userDeleteCmd)
}
