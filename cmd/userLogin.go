package cmd

import (
	"context"
	"time"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// userAuthCmd represents the sessionClose command
var userAuthCmd = &cobra.Command{
	Use:   "login <username>",
	Short: "Authenticates a user against an API server",
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
		client := api.NewUserServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.CLI.Timeout)*time.Second)
		defer cancel()

		password, err := cli.GetPassword(cmd)
		if err != nil {
			log.WithError(err).Fatal()
		}

		resp, err := client.AuthenticateCredentials(ctx, &api.AuthenticationRequest{Username: args[0], Password: password})
		if err != nil {
			log.WithError(err).Fatal()
		}

		tpl := `export RURUKU_TOKEN={{ .Token }}
export RURUKU_USER={{ .User }}
`
		result := struct {
			Token string
			User  string
		}{
			Token: resp.Token,
			User:  args[0],
		}
		ctnt := remoteCmdValues.GetOutputFormat(result, tpl)
		if err := ctnt.Print(); err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	userCmd.AddCommand(userAuthCmd)

	userAuthCmd.Flags().String("password", "", "The user's password. Passwords on the command line are unsafe: use the interactive mode or RURUKU_PASSWORD env var instead")
}
