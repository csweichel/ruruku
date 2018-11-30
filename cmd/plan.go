package cmd

import (
	"fmt"
	"github.com/32leaves/ruruku/pkg/cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "plan",
	Short: "Creates a new testplan",
	Long: `Creates a new testplan. It goes into interactive mode if called without flags, e.g.:
    ruruku plan

Most fields are available as flags as well, e.g.:
    ruruku plan -f testplan.yaml --name Demo
    ruruku plan -f testplan.yaml add --name "My testcase" --id tc1 --group grp
    `,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initFlags.Run(); err != nil {
			log.WithError(err).Error()
			return
		}
		log.Info(fmt.Sprintf("Use ruruku plan add -f %s to add testcases", initFlags.Filename))
	},
}

var initFlags = cli.InitPlan{}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().BoolVarP(&initFlags.NonInteractive, "non-interactive", "y", false, "do not use an interactive prompt. Excepts all fields to be provided as flags.")
	initCmd.PersistentFlags().StringVarP(&initFlags.Filename, "filename", "f", "", "the output filename")

	initCmd.Flags().StringVarP(&initFlags.ID, "id", "i", "", "ID of the testplan")
	initCmd.Flags().StringVarP(&initFlags.Name, "name", "n", "", "name of the testplan")
}
