package cmd

import (
	"io"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sessionListCmd represents the sessionList command
var sessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "Prints a table of the available sessions and their status",
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

		stream, err := client.List(ctx, &api.ListSessionsRequest{})
		if err != nil {
			log.WithError(err).Fatal()
		}

		resp := make([]*api.ListSessionsResponse, 0)
		for {
			session, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.WithError(err).Fatal()
			}
			resp = append(resp, session)
		}

		tpl := `ID	IS OPEN	NAME
{{- range . }}
{{ .Id }}	{{ .IsOpen }}	{{ .Name -}}
{{ end }}
`
		ctnt := remoteCmdValues.GetOutputFormat(resp, tpl)
		if err := ctnt.Print(); err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	sessionCmd.AddCommand(sessionListCmd)
}
