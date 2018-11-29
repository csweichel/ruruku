package cmd

import (
    api "github.com/32leaves/ruruku/pkg/server/api/v1"
	"google.golang.org/grpc"
	"github.com/spf13/cobra"
    "context"
    "time"
    "github.com/technosophos/moniker"
    log "github.com/sirupsen/logrus"
)

type sessionStartFlags struct {
    name string
    planfn string
}

var sessionStartFlagValues sessionStartFlags

func (s *sessionStartFlags) Run() error {
    name := moniker.New().Name()
    if s.name == "" {
        log.WithField("name", name).Info("Using an auto-generated session name")
    } else {
        name = s.name
    }

    if s.planfn == "" {
        log.Warn("Starting a session without a plan")
    }

    conn, err := grpc.Dial(sessionFlagValues.server, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("fail to dial: %v", err)
    }
    defer conn.Close()
    client := api.NewSessionServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sessionFlagValues.timeout)*time.Second)
    defer cancel()

    req := &api.StartSessionRequest{
        Name: name,
    }
    resp, err := client.Start(ctx, req)
    if err != nil {
        return err
    }

    log.WithField("id", resp.Id).Info("Session started")

    return nil
}

// sessionStartCmd represents the sessionStart command
var sessionStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a new test session based on a test plan",
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
