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
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"

	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	proto "github.com/chronojam/trubble/api/proto"

	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("pass game name, version and binary\n test game upload {game_name} {version} {path_to_binary}")
		}

		game := args[0]
		version := args[1]
		bpath := args[2]

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

		stream, err := gameapi.UploadBinary(ctx)
		if err != nil {
			log.Fatalf("%v.UploadBinary(_) = _, %v", game, err)
		}

		chunkSize := 1048576
		f, err := ioutil.ReadFile(bpath)
		if err != nil {
			panic(err)
		}
		size := len(f)

		log.Printf("%v", size)
		numberOfChunks := size / chunkSize
		log.Printf("%v", numberOfChunks)

		request := &proto.UploadBinaryRequest{
			Value: &proto.UploadBinaryRequest_Key_{
				Key: &proto.UploadBinaryRequest_Key{
					Game:    game,
					Version: version,
					Size:    int64(size),
				},
			},
		}

		if err := stream.Send(request); err != nil {
			log.Fatalf("[KEY] %v.Send(%v) = %v", stream, request, err)
		}

		for chunkNum := 0; chunkNum <= numberOfChunks; chunkNum++ {
			var chunkEnd int
			chunkStart := chunkNum * chunkSize
			if chunkNum == numberOfChunks {
				chunkEnd = size
			} else {
				nextChunk := chunkNum + 1
				chunkEnd = nextChunk * chunkSize
			}

			log.Printf("ChunkStart: %v -> %v", chunkStart, chunkEnd)
			chunk := f[chunkStart:chunkEnd]
			request := &proto.UploadBinaryRequest{
				Value: &proto.UploadBinaryRequest_Chunk_{
					Chunk: &proto.UploadBinaryRequest_Chunk{
						Data: chunk,
					},
				},
			}

			if err := stream.Send(request); err != nil {
				log.Fatalf("[CHUNK] %v.Send(%v) = %v", stream, request, err)
			}
		}

		reply, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
		}

		log.Printf("UploadGameBinaryResponse: %v", reply)
	},
}

func init() {
	gameCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
