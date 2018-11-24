package cmd

import (
	"fmt"
	"github.com/32leaves/ruruku/pkg/rurukulib/cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a new testsuite",
	Long: `Init creates a new testsuite. It goes into interactive mode if called without flags, e.g.:
    ruruku init

Most fields are available as flags as well, e.g.:
    ruruku init -f testsuite.yaml --name Demo
    ruruku init -f testsuite.yaml add --name "My testcase" --id tc1 --group grp
    `,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initFlags.Run(); err != nil {
			log.WithError(err).Error()
			return
		}
		log.Info(fmt.Sprintf("Use ruruku init testcase -f %s to add testcases", initFlags.Filename))
	},
}

var initFlags = cli.InitSuite{}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().BoolVarP(&initFlags.NonInteractive, "non-interactive", "y", false, "do not use an interactive prompt. Excepts all fields to be provided as flags.")
	initCmd.PersistentFlags().StringVarP(&initFlags.Filename, "filename", "f", "", "the output filename")

	initCmd.Flags().StringVarP(&initFlags.Name, "name", "n", "", "name of the testsuite")
}
