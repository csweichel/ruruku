package cmd

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

type testContributionFlags struct {
	comment string
}

var testContributionFlagValues testContributionFlags

// testContributeCmd represents the testContribute command
var testContributeCmd = &cobra.Command{
	Use:       "contribute <testcase-id> passed|undecided|failed",
	Short:     "Contribute results of a test execution",
	Args:      cobra.ExactArgs(2),
	ValidArgs: []string{"passed", "undecided", "failed"},
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := remoteCmdValues.Connect()
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := api.NewSessionServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(remoteCmdValues.timeout)*time.Second)
		defer cancel()

		result := api.TestRunState_FAILED
		if args[1] == "passed" {
			result = api.TestRunState_PASSED
		} else if args[1] == "undecided" {
			result = api.TestRunState_UNDECIDED
		} else if args[1] == "failed" {
			result = api.TestRunState_FAILED
		} else {
			log.Fatalf("Testcase result must be passed, undecided or failed. Not %s", args[1])
		}

		_, err = client.Contribute(ctx, &api.ContributionRequest{
			ParticipantToken: testCmdToken,
			TestcaseID:       args[0],
			Comment:          testContributionFlagValues.comment,
			Result:           result,
		})
		if err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	testCmd.AddCommand(testContributeCmd)
	testContributeCmd.Flags().StringVarP(&testContributionFlagValues.comment, "comment", "m", "", "Additional comment")
}
