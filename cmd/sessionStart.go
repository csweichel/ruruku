package cmd

import (
	"fmt"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/technosophos/moniker"
)

type sessionStartFlags struct {
	name        string
	planfn      string
	quiet       bool
	annotations map[string]string
	modifiable  bool
}

var sessionStartFlagValues sessionStartFlags

func (s *sessionStartFlags) Run() error {
	cfg, err := GetConfigFromViper()
	if err != nil {
		log.Fatalf("Error while loading the configuration: %v", err)
	}

	req := &api.StartSessionRequest{
		Name:        moniker.New().Name(),
		Annotations: s.annotations,
		Modifiable:  s.modifiable,
	}

	if s.name == "" {
		log.WithField("name", req.Name).Info("Using an auto-generated session name")
	} else {
		req.Name = s.name
	}

	if s.planfn != "" {
		plan, err := cli.LoadTestplan(s.planfn)
		if err != nil {
			return err
		}
		req.Plan = api.ConvertTestPlan(plan)
	} else {
		req.Plan = &api.TestPlan{
			Id:   req.Name,
			Name: req.Name,
			Case: []*api.Testcase{},
		}
	}

	conn, err := cfg.Connect()
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := api.NewSessionServiceClient(conn)

	ctx, cancel := cfg.GetContext(true)
	defer cancel()

	resp, err := client.Start(ctx, req)
	if err != nil {
		return err
	}

	if !s.quiet {
		tpl := `{{ .Id }}`
		ctnt := remoteCmdValues.GetOutputFormat(resp, tpl, `{.id}`)
		if err := ctnt.Print(); err != nil {
			log.WithError(err).Fatal()
		}
	}

	log.WithField("id", resp.Id).Info("Session started")

	return nil
}

// sessionStartCmd represents the sessionStart command
var sessionStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a new test session based on a test plan",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !sessionStartFlagValues.modifiable && sessionStartFlagValues.planfn == "" {
			return fmt.Errorf("Cannot start an unmodifiable session without a plan (use --plan or --modifiable)")
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
	sessionStartCmd.Flags().StringToStringVarP(&sessionStartFlagValues.annotations, "annotations", "a", map[string]string{}, "Metadata for this session")
	sessionStartCmd.Flags().BoolVar(&sessionStartFlagValues.modifiable, "modifiable", false, "Make this session modifiable")
}
