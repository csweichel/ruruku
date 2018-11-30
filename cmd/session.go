package cmd

import (
	"github.com/spf13/cobra"
)

type sessionFlags struct {
	server  string
	timeout uint32
}

var sessionFlagValues sessionFlags

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Controls a test session/run",
	Args:  cobra.ExactArgs(1),
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(sessionCmd)

	sessionCmd.PersistentFlags().StringVarP(&sessionFlagValues.server, "server", "", "localhost:1234", "Host address of the ruruku API server")
	sessionCmd.PersistentFlags().Uint32VarP(&sessionFlagValues.timeout, "timeout", "", 10, "Request timeout in seconds")
}
