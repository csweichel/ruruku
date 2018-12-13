package cmd

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var userAddCmdFile string

// userAuthCmd represents the sessionClose command
var userAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new user",
	Long: `Add a new user to a Ruruku installation. Users can be added from a file or using command line flags.

To add users from a file, use the -f flag pointing to a YAML file with username, password and email properties.
The optional permission property can be used to add permissions at the same time, e.g.

    username: foo
    password: ThisIsMySecretPassword
    email: foo@bar.com
    permission:
    - session.start
    - session.contribute

You can add multiple users in one go by using the YAML file separater ---.`,
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

		ctx, cancel := cfg.GetContext()
		defer cancel()

		if userAddCmdFile != "" {
			if err := cli.AddUserFromFile(client, ctx, userAddCmdFile); err != nil {
				log.WithError(err).Fatal()
			}
		} else {
			username := viper.GetString("username")
			email := viper.GetString("email")
			password, err := cli.GetPassword(cmd)
			if err != nil {
				log.WithError(err).Fatal()
			}

			if _, err := client.Add(ctx, &api.AddUserRequest{Username: username, Password: password, Email: email}); err != nil {
				log.WithError(err).Fatal()
			}
		}
	},
}

func init() {
	userCmd.AddCommand(userAddCmd)

	userAddCmd.Flags().StringVarP(&userAddCmdFile, "file", "f", "", "Add all users from a YAML file")

	userAddCmd.Flags().String("name", "", "The username")
	viper.BindPFlag("username", userAddCmd.Flags().Lookup("name"))
	userAddCmd.Flags().String("email", "", "The user's email address")
	viper.BindPFlag("email", userAddCmd.Flags().Lookup("email"))
	userAddCmd.Flags().String("password", "", "The user's password. Passwords on the command line are unsafe: use the interactive mode or RURUKU_PASSWORD env var instead")
}
