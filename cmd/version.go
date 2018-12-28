package cmd

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type versionInfo struct {
	Tag       string `json:"tag"`
	Rev       string `json:"rev"`
	BuildDate string `json:"buildDate"`
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:               "version",
	Short:             "Prints the version of ruruku",
	PersistentPreRunE: remoteCmdValuesPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		client := versionInfo{
			Tag:       version.Tag,
			Rev:       version.Rev,
			BuildDate: version.BuildDate,
		}
		server := getVersionInfoFromServer()

		versions := make(map[string]versionInfo)
		versions["client"] = client
		if server != nil {
			versions["server"] = *server
		}

		tpl := `{{- range $k, $v := . -}}
{{ $k }}: {{ $v.Tag }} ({{ $v.Rev }}) built on {{ $v.BuildDate }}
{{ end -}}
`
		ctnt := remoteCmdValues.GetOutputFormat(versions, tpl, `{.client.tag}`)
		if err := ctnt.Print(); err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func getVersionInfoFromServer() *versionInfo {
	cfg, err := GetConfigFromViper()
	if err != nil {
		return nil
	}

	conn, err := cfg.Connect()
	if err != nil {
		return nil
	}
	defer conn.Close()
	client := api.NewVersionServiceClient(conn)

	ctx, cancel := cfg.GetContext(true)
	defer cancel()

	res, err := client.Get(ctx, &api.GetVersionRequest{})
	if err != nil {
		return nil
	}

	return &versionInfo{
		Tag:       res.Tag,
		Rev:       res.Rev,
		BuildDate: res.BuildDate,
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	registerRemoteCmdValueFlags(versionCmd)
}
