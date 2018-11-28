package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/32leaves/ruruku/pkg/storage"
    "github.com/32leaves/ruruku/pkg/cli"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status <testsuite.yaml> <session.yaml>",
	Short: "Displays the status of a test session",
    Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		suite, err := storage.LoadSuite(args[0])
        if err != nil {
            log.WithError(err).Fatalf("Unable to load testsuite")
        }

        session, err := storage.LoadFileStorage(args[1], suite)
        if err != nil {
            log.WithError(err).Fatalf("Unable to load session")
        }

        status := cli.ComputeStatus(session)
        status.Print(true)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
