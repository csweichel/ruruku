package cmd

import (
	"github.com/32leaves/ruruku/pkg/cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initTestcaseCmd represents the initTestcase command
var initTestcaseCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a testcase to a plan",
	Run: func(cmd *cobra.Command, args []string) {
		initTestcaseFlags.Filename = initFlags.Filename
		initTestcaseFlags.NonInteractive = initFlags.NonInteractive
		initTestcaseFlags.MinTesterCountSet = cmd.Flag("min-tester-count").Changed

		if err := initTestcaseFlags.Run(); err != nil {
			log.WithError(err).Error()
		}
	},
}

var initTestcaseFlags = cli.InitTestcase{}

func init() {
	initCmd.AddCommand(initTestcaseCmd)

	initTestcaseCmd.Flags().StringVarP(&initTestcaseFlags.ID, "id", "", "", "testcase ID (must be unique within the group)")
	initTestcaseCmd.Flags().StringVarP(&initTestcaseFlags.Group, "group", "g", "", "testcase group")
	initTestcaseCmd.Flags().StringVarP(&initTestcaseFlags.Name, "name", "n", "", "name/short description of the testcase")
	initTestcaseCmd.Flags().StringVarP(&initTestcaseFlags.Description, "description", "d", "", "long description - can be markdown")
	initTestcaseCmd.Flags().StringVarP(&initTestcaseFlags.Steps, "steps", "s", "", "steps to perform during the test - can be markdown")
	initTestcaseCmd.Flags().Uint32VarP(&initTestcaseFlags.MinTesterCount, "min-tester-count", "", 0, "number of testers for this testcase")
	initTestcaseCmd.Flags().BoolVarP(&initTestcaseFlags.MustPass, "must-pass", "", false, "test must pass for the testsuite to pass")
}
