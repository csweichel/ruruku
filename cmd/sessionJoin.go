package cmd

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/server/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

// sessionJoinCmd represents the sessionJoin command
var sessionJoinCmd = &cobra.Command{
	Use:   "join <session-id> <participant-name>",
	Short: "Joins a testing session",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(remoteCmdValues.server, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := api.NewSessionServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(remoteCmdValues.timeout)*time.Second)
		defer cancel()

		resp, err := client.Register(ctx, &api.RegistrationRequest{SessionID: args[0], Name: args[1]})
		if err != nil {
			log.WithError(err).Fatal()
		}

		tpl := `{{ .Token }}`
		ctnt := remoteCmdValues.GetOutputFormat(resp, tpl)
		if err := ctnt.Print(); err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	sessionCmd.AddCommand(sessionJoinCmd)
}
