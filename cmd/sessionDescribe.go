package cmd

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const sessionDescribeTpl = `ID:	{{ .Id }}
Name:	{{ .Name }}
Plan:	{{ .PlanID }}
Result:	{{ .State }}
{{- if .Annotations }}
Annotations:	{{ range $k, $v := .Annotations }}{{ $k }}:	{{ $v }}
	{{ end -}}
{{ end }}
Tests:
{{- range .Status }}
  {{ .Case.Id }}:
    Name:	{{ .Case.Name }}
    Result:	{{ .State }}
    Claims:	{{ len .Claim -}}
{{ end }}
`

// sessionDescribeCmd represents the sessionDescribe command
var sessionDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Prints session details and its testcases",
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

		resp, err := client.Status(ctx, &api.SessionStatusRequest{Id: args[0]})
		if err != nil {
			log.WithError(err).Fatal()
		}

		ctnt := remoteCmdValues.GetOutputFormat(resp.Status, sessionDescribeTpl)
		if err := ctnt.Print(); err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	sessionCmd.AddCommand(sessionDescribeCmd)
}
