package cmd

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// testClaimCmd represents the testClaim command
var testClaimCmd = &cobra.Command{
	Use:   "claim <testcase-id>",
	Short: "Claim a test for a participant",
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

		_, err = client.Claim(ctx, &api.ClaimRequest{
			Session:    testCmdSession,
			TestcaseID: args[0],
			Claim:      true,
		})
		if err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	testCmd.AddCommand(testClaimCmd)
}
