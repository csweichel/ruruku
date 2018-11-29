package cmd

import (
	"fmt"
	"github.com/32leaves/ruruku/pkg/server"
	api "github.com/32leaves/ruruku/pkg/server/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net"
)

// grpcServerCmd represents the grpcServer command
var grpcServerCmd = &cobra.Command{
	Use:   "grpc-server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 1234))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		api.RegisterSessionServiceServer(grpcServer, server.NewMemoryBackedSession())
		grpcServer.Serve(lis)
	},
}

func init() {
	rootCmd.AddCommand(grpcServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grpcServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grpcServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
