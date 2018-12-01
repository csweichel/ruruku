package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var testCmdToken string

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Claim or contribute to tests",
	Args:  cobra.ExactArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := remoteCmdValuesPreRun(cmd, args); err != nil {
			return err
		}
		if testCmdToken == "" {
			return fmt.Errorf("--token is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatalf("Unkown command: %s", args[0])
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	registerRemoteCmdValueFlags(testCmd)
	testCmd.PersistentFlags().StringVarP(&testCmdToken, "token", "t", "", "Session participant token (use session join to get one)")
}
