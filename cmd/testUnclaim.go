package cmd

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

// testUnclaimCmd represents the testUnclaim command
var testUnclaimCmd = &cobra.Command{
	Use:   "unclaim <testcase-id>",
	Short: "Unclaim a test for a participant",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := remoteCmdValues.Connect()
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := api.NewSessionServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(remoteCmdValues.timeout)*time.Second)
		defer cancel()

		_, err = client.Claim(ctx, &api.ClaimRequest{
			ParticipantToken: testCmdToken,
			TestcaseID:       args[0],
			Claim:            false,
		})
		if err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	testCmd.AddCommand(testUnclaimCmd)
}
