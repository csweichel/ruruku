package cmd

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// userAuthCmd represents the sessionClose command
var userChpwdCmd = &cobra.Command{
	Use:   "chpwd [username]",
	Short: "Changes the password of a user",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := GetConfigFromViper()
		if err != nil {
			log.Fatalf("Error while loading the configuration: %v", err)
		}

		username := os.Getenv("RURUKU_USER")
		if len(args) > 0 {
			username = args[0]
		}
		if username == "" {
			log.Fatal("no username")
		}

		password, err := cli.GetPassword(cmd)
		if err != nil {
			log.WithError(err).Fatal()
		}

		conn, err := cfg.Connect()
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := api.NewUserServiceClient(conn)

		ctx, cancel := cfg.GetContext()
		defer cancel()

		if _, err := client.ChangePassword(ctx, &api.ChangePasswordRequest{Username: username, NewPassword: password}); err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	userCmd.AddCommand(userChpwdCmd)

	userChpwdCmd.Flags().String("password", "", "The user's password. Passwords on the command line are unsafe: use the interactive mode or RURUKU_PASSWORD env var instead")
}
