package cmd

import (
	"context"
	"fmt"
	api "github.com/32leaves/ruruku/pkg/server/api/v1"
	"google.golang.org/grpc"
	"log"
	"github.com/spf13/cobra"
	"io"
	"time"
)

// grpcClientCmd represents the grpcClient command
var grpcClientCmd = &cobra.Command{
	Use:   "grpc-client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := api.NewSessionServiceClient(conn)

		ctx, cancelVersionReq := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelVersionReq()
		version, _ := client.Version(ctx, &api.VersionRequest{})
		fmt.Printf("%s\n", version)

		ctx, cancelUpdates := context.WithTimeout(context.Background(), 1*time.Second)
		updates, err := client.Updates(ctx, &api.SessionUpdatesRequest{})
		if err != nil {
			log.Fatalf("Error while subscribing: %v", err)
		}
		for {
			update, err := updates.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("Error while receving update: %v", err)
			} else {
				fmt.Printf("Update: %s\n", update)
			}
		}
		cancelUpdates()

		_, err = client.Register(context.Background(), &api.RegistrationRequest{
			Name:      "foobar",
			SessionID: "foo",
		})
		if err != nil {
			log.Fatalf("cannot register: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(grpcClientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grpcClientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grpcClientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
