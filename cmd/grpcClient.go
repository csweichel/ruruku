// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
    "time"
    "context"
    "io"
	"github.com/spf13/cobra"
    api "github.com/32leaves/ruruku/pkg/server/api"
    "google.golang.org/grpc"
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
            Name: "foobar",
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
