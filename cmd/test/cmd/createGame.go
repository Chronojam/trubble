// Copyright © 2017 Calum Gardner <calum@qubit.com>
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
	"context"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"

	proto "github.com/chronojam/trubble/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Pass game name, test game create {GameName} {Version}")
		}

		game := args[0]
		version := args[1]

		usr, err := user.Current()
		if err != nil {
			log.Fatal(err.Error())
		}

		b, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, ".trubble-token"))
		if err != nil {
			log.Fatal(err.Error())
		}

		address := "127.0.0.1:3334"

		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err.Error())
		}
		defer conn.Close()

		gameapi := proto.NewGameClient(conn)

		md := metadata.Pairs("jwt", string(b))
		ctx := metadata.NewContext(context.Background(), md)

		gameReq := &proto.CreateGameRequest{
			Name:    game,
			Version: version,
		}

		_, err = gameapi.CreateGame(ctx, gameReq)
		if err != nil {
			log.Printf(err.Error())
		}

		log.Printf("Created Game %s,\nVersion: %s", game, version)
	},
}

func init() {
	gameCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
