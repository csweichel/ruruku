package cmd

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.CLI.Timeout)*time.Second)
		defer cancel()

		_, err = client.Claim(ctx, &api.ClaimRequest{
			ParticipantToken: testCmdToken,
			TestcaseID:       args[0],
			Claim:            true,
		})
		if err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	testCmd.AddCommand(testClaimCmd)
}
