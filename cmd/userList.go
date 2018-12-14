package cmd

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type userListOutput struct {
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Permissions []types.Permission `json:"permissions"`
}

// userListCmd represents the sessionClose command
var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all users in the system (requires user.list permission)",
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

		ctx, cancel := cfg.GetContext(true)
		defer cancel()

		resp, err := client.List(ctx, &api.ListUsersRequest{})
		if err != nil {
			log.WithError(err).Fatal()
		}

		out := make([]userListOutput, len(resp.User))
		for idx, u := range resp.User {
			perms := make([]types.Permission, len(u.Permission))
			for jdx, p := range u.Permission {
				perms[jdx] = p.Convert()
			}
			out[idx] = userListOutput{
				Name:        u.Name,
				Email:       u.Email,
				Permissions: perms,
			}
		}

		tpl := `NAME	EMAIL	PERMISSION
{{- range . }}
{{ .Name }}	{{ .Email }}	{{ .Permissions -}}
{{ end }}
`
		ctnt := remoteCmdValues.GetOutputFormat(out, tpl)
		if err := ctnt.Print(); err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	userCmd.AddCommand(userListCmd)
}
