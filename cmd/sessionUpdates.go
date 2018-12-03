package cmd

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
)

var waitUntilCloses bool

// sessionUpdatesCmd represents the sessionUpdates command
var sessionUpdatesCmd = &cobra.Command{
	Use:   "updates <session-id>",
	Short: "Listens for changes in the session",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(remoteCmdValues.server, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := api.NewSessionServiceClient(conn)

		clnt, err := client.Updates(context.Background(), &api.SessionUpdatesRequest{Id: args[0]})
		if err != nil {
			log.WithError(err).Fatal()
		}

		for {
			resp, err := clnt.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.WithError(err).Fatal()
				break
			}

			ctnt := remoteCmdValues.GetOutputFormat(resp.Status, sessionDescribeTpl)
			if err := ctnt.Print(); err != nil {
				log.WithError(err).Fatal()
			}

			if waitUntilCloses && !resp.Status.Open {
				break
			}
		}
	},
}

func init() {
	sessionCmd.AddCommand(sessionUpdatesCmd)
	sessionUpdatesCmd.Flags().BoolVarP(&waitUntilCloses, "wait-until-closed", "w", true, "Stop listening once the session is closed")
}
