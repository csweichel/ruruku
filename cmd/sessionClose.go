package cmd

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

// sessionCloseCmd represents the sessionClose command
var sessionCloseCmd = &cobra.Command{
	Use:   "close <session-id>",
	Short: "Closes a testing session",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(remoteCmdValues.server, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := api.NewSessionServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(remoteCmdValues.timeout)*time.Second)
		defer cancel()

		_, err = client.Close(ctx, &api.CloseSessionRequest{Id: args[0]})
		if err != nil {
			log.WithError(err).Fatal()
		}
	},
}

func init() {
	sessionCmd.AddCommand(sessionCloseCmd)
}
