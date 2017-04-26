package main

import (
	"fmt"
	"log"

	"encoding/json"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"

	"golang.org/x/net/context"

	"github.com/pkg/errors"

	proto "github.com/chronojam/trubble/api/proto"
	"io"
)

// GameServer ...
// GRPC Server for game management.
type GameServer struct{}
type GameBinary struct {
	Version string
	Game    string
	Size    int64
	Data    []byte
}

// CreateGame ...
//
func (g *GameServer) CreateGame(ctx context.Context, req *proto.CreateGameRequest) (*proto.CreateGameResponse, error) {
	if req.GetName() == "" {
		return nil, errors.New("name is a required field")
	}

	if req.GetVersion() == "" {
		return nil, errors.New("version is a required field")
	}

	gameName := req.GetName()
	gameVersion := req.GetVersion()

	// "github.com/chronojam/trubble/api/main.go:79"
	accountID := ctx.Value("_account_id").(float64)

	var games []Game

	query := datastore.NewQuery("Game").
		Filter("Name =", gameName).
		Filter("Version =", gameVersion).
		Limit(1)

	_, err := dsClient.GetAll(ctx, query, &games)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current games from datastore")
	}

	if len(games) > 0 {
		return nil, errors.New("a game with that version already exists")
	}

	key := datastore.NameKey("Game", "", nil)
	game := Game{
		Name:    gameName,
		Version: gameVersion,
		Owner:   accountID,
	}

	_, err = dsClient.Put(ctx, key, &game)
	if err != nil {
		return nil, errors.Wrap(err, "could not write to datastore")
	}

	resp := &proto.CreateGameResponse{}

	return resp, nil
}

// UploadBinary ...
func (g *GameServer) UploadBinary(stream proto.Game_UploadBinaryServer) error {
	gameBinary := new(GameBinary)
	ctx := stream.Context()
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			if int64(len(gameBinary.Data)) < gameBinary.Size {
				log.Printf("gameBinary.Data: %v, gameBinary.Size: %v", int64(len(gameBinary.Data)), gameBinary.Size)
				return errors.Wrap(err, "Binary/Size mismatch")
			}
			gameName := gameBinary.Game
			gameVersion := gameBinary.Version

			// "github.com/chronojam/trubble/api/main.go:79"
			accountID := ctx.Value("_account_id").(float64)

			var games []Game

			query := datastore.NewQuery("Game").
				Filter("Name =", gameName).
				Filter("Version =", gameVersion).
				Limit(1)

			_, err := dsClient.GetAll(ctx, query, &games)
			if err != nil {
				return errors.Wrap(err, "failed to get current games from datastore")
			}

			if len(games) == 0 {
				return errors.New("no game with that version exists")
			}

			storClient, err := storage.NewClient(ctx)
			if err != nil {
				return errors.Wrap(err, "unable to initialize storage client")
			}

			bkt := storClient.Bucket("trubble-data")
			obj := bkt.Object(fmt.Sprintf("%v/games/%s/versions/%s/binaries/server", int64(accountID), gameName, gameVersion))

			w := obj.NewWriter(ctx)
			w.Write(gameBinary.Data)

			if err := w.Close(); err != nil {
				return errors.Wrap(err, "could not close connection?")
			}

			// Put a message onto pubsub queue for build server.
			go PutPubsubMessage("builder", BuilderMessage{
				Game:      gameName,
				Version:   gameVersion,
				AccountID: accountID,
			})

			return stream.SendAndClose(&proto.UploadBinaryResponse{})

		}
		if err != nil {
			return err
		}
		switch x := message.Value.(type) {
		case *proto.UploadBinaryRequest_Key_:
			gameBinary.Size = x.Key.Size
			gameBinary.Version = x.Key.Version
			gameBinary.Game = x.Key.Game
		case *proto.UploadBinaryRequest_Chunk_:
			gameBinary.Data = append(gameBinary.Data, x.Chunk.Data...)

		default:
			return fmt.Errorf("GameBinary has unexpected type %T", x)

		}

	}
}

// BuilderMessage ...
// Message to instruct the build service to start building a container
type BuilderMessage struct {
	Game      string  `json:"game"`
	Version   string  `json:"version"`
	AccountID float64 `json:"account_id"`
}

// PutPubsubMessage ...
// Retry on failures.
func PutPubsubMessage(topicName string, message interface{}) {
	ctx := context.Background()
	for {
		// Messages are parsed as json blobs
		b, err := json.Marshal(message)
		if err != nil {
			log.Printf("failed to publish message %s", err.Error())
			continue
		}

		top := psClient.Topic(topicName)
		res := top.Publish(ctx, &pubsub.Message{Data: b})
		id, err := res.Get(ctx)
		if err != nil {
			// Failed to put onto queue for whatever reason
			log.Printf("failed to publish message %s", err.Error())
			continue
		}

		log.Printf("Published message with ID %s\n", id)
		return
	}
}
