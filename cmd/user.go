package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:               "user",
	Short:             "Interacts with ruruku users, e.g. authenticates, adds or deletes users",
	Args:              cobra.ExactArgs(1),
	PersistentPreRunE: remoteCmdValuesPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatalf("Unkown command: %s", args[0])
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	registerRemoteCmdValueFlags(userCmd)
}
