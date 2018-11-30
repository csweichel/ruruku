package cmd

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/server/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

// sessionDescribeCmd represents the sessionDescribe command
var sessionDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Prints session details and its testcases",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(sessionFlagValues.server, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := api.NewSessionServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sessionFlagValues.timeout)*time.Second)
		defer cancel()

		status, err := client.Status(ctx, &api.SessionStatusRequest{Id: args[0]})
		if err != nil {
			log.WithError(err).Fatal()
		}

		log.WithField("status", status).Info()
	},
}

func init() {
	sessionCmd.AddCommand(sessionDescribeCmd)
}
