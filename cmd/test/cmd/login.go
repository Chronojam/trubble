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

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if email == "" || password == "" {
			log.Printf("Pass email --email and --password")
			return
		}

		address := "127.0.0.1:3333"
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err.Error())
		}

		defer conn.Close()

		ctx := context.Background()
		auth := proto.NewAuthClient(conn)

		jwtReq := &proto.NewJwtRequest{
			Email:    email,
			Password: password,
		}

		jwtResp, err := auth.IssueNewJWT(ctx, jwtReq)
		if err != nil {
			log.Fatal(err.Error())
		}

		usr, err := user.Current()
		if err != nil {
			log.Fatal(err.Error())
		}

		tok := jwtResp.GetToken()
		err = ioutil.WriteFile(filepath.Join(usr.HomeDir, ".trubble-token"), []byte(tok), 0655)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("Wrote token to ~/.trubble-token")
	},
}

var email string
var password string

func init() {
	RootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	loginCmd.Flags().StringVarP(&email, "email", "e", "", "Email address.")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Password")

}
