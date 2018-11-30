package cmd

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/server/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
	"time"
)

// sessionListCmd represents the sessionList command
var sessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "Prints a table of the available sessions and their status",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(sessionFlagValues.server, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := api.NewSessionServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sessionFlagValues.timeout)*time.Second)
		defer cancel()

		stream, err := client.List(ctx, &api.ListSessionsRequest{})
		if err != nil {
			log.WithError(err).Fatal()
		}

		for {
			session, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.WithError(err).Fatal()
			}

			log.WithField("ID", session.Id).WithField("Name", session.Name).WithField("Open", session.IsOpen).Info()
		}
	},
}

func init() {
	sessionCmd.AddCommand(sessionListCmd)
}
