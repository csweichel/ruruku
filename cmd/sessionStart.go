package cmd

import (
	"context"
	"fmt"
	"github.com/32leaves/ruruku/pkg/cli"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/technosophos/moniker"
	"google.golang.org/grpc"
	"time"
)

type sessionStartFlags struct {
	name   string
	planfn string
}

var sessionStartFlagValues sessionStartFlags

func (s *sessionStartFlags) Run() error {
	req := &api.StartSessionRequest{
		Name: moniker.New().Name(),
	}

	if s.name == "" {
		log.WithField("name", req.Name).Info("Using an auto-generated session name")
	} else {
		req.Name = s.name
	}

	plan, err := cli.LoadTestplan(s.planfn)
	if err != nil {
		return err
	}
	req.Plan = api.ConvertTestPlan(plan)

	conn, err := grpc.Dial(remoteCmdValues.server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := api.NewSessionServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(remoteCmdValues.timeout)*time.Second)
	defer cancel()

	resp, err := client.Start(ctx, req)
	if err != nil {
		return err
	}

	tpl := `{{ .Id }}`
	ctnt := remoteCmdValues.GetOutputFormat(resp, tpl)
	if err := ctnt.Print(); err != nil {
		log.WithError(err).Fatal()
	}

	log.WithField("id", resp.Id).Info("Session started")

	return nil
}

// sessionStartCmd represents the sessionStart command
var sessionStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a new test session based on a test plan",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if sessionStartFlagValues.planfn == "" {
			return fmt.Errorf("Cannot start a session without a plan (use --plan)")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := sessionStartFlagValues.Run(); err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	sessionCmd.AddCommand(sessionStartCmd)

	sessionStartCmd.Flags().StringVarP(&sessionStartFlagValues.name, "name", "n", "", "Name of the session")
	sessionStartCmd.Flags().StringVarP(&sessionStartFlagValues.planfn, "plan", "p", "", "Path to the test plan of this session")
}
