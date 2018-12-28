package cmd

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// planFromSessionCmd represents the planFromSession command
var planFromSessionCmd = &cobra.Command{
	Use:               "from-session <sessionID>",
	Short:             "Extracts a testplan from a test session",
	Args:              cobra.ExactArgs(1),
	PersistentPreRunE: remoteCmdValuesPreRun,
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

		status, err := client.Status(ctx, &api.SessionStatusRequest{Id: args[0]})
		if err != nil {
			log.Fatalf("Unable to get session status: %v", err)
		}

		stat := status.Status.Convert()
		plan := stat.ToTestplan()
		initFlags.Plan = &plan
		if err := initFlags.Run(); err != nil {
			log.WithError(err).Error()
			return
		}
	},
}

func init() {
	initCmd.AddCommand(planFromSessionCmd)
	registerRemoteCmdValueFlags(planFromSessionCmd)
}
