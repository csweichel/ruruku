package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:               "session",
	Short:             "Controls a test session/run",
	Args:              cobra.ExactArgs(1),
	PersistentPreRunE: remoteCmdValuesPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatalf("Unkown command: %s", args[0])
	},
}

func init() {
	rootCmd.AddCommand(sessionCmd)
	registerRemoteCmdValueFlags(sessionCmd)
}
