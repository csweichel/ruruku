package cmd

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var userGrantPermissions []string

// userAuthCmd represents the sessionClose command
var userGrantCmd = &cobra.Command{
	Use:   "grant <username>",
	Short: "Grants permissions to a user (requires user.grant permission)",
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

		ctx, cancel := cfg.GetContext()
		defer cancel()

		permissions := make([]api.Permission, len(userGrantPermissions))
		for idx, perm := range userGrantPermissions {
			permissions[idx] = api.ConvertPermission(types.Permission(perm))
		}

		req := api.GrantPermissionsRequest{
			Username:   args[0],
			Permission: permissions,
		}
		if _, err := client.Grant(ctx, &req); err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	userCmd.AddCommand(userGrantCmd)

	userGrantCmd.Flags().StringArrayVarP(&userGrantPermissions, "permission", "p", []string{}, "Permission to add to the user")
}
