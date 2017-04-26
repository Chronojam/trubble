package main

import (
	"log"
	"fmt"

	"golang.org/x/net/context"
	"encoding/json"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"os"
	"os/exec"
	"io/ioutil"
)

type BuilderMessage struct {
	Game	string 	`json:"game"`
	Version	string	`json:"version"`
	AccountID float64 `json:"account_id"`
}

func main(){
	ctx := context.Background()
	// Setup pubsub Client
	p, err := pubsub.NewClient(ctx, "chronojam-trubble")
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Storage Client
	storClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	bkt := storClient.Bucket("trubble-data")

	topic := p.Topic("builder")
	sub, err := p.CreateSubscription(context.Background(), "builder", topic, 0, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			log.Print("Got Message")
			mess := BuilderMessage{}
			err := json.Unmarshal(m.Data, &mess)
			if err != nil {
				log.Printf(err.Error())
				m.Nack()
				return
			}

			// We've got the message, and it is intended for us.
			// Lets build our container.
			// Download GCS object.
			obj := bkt.Object(fmt.Sprintf("%v/games/%s/versions/%s/binaries/server", int64(mess.AccountID), mess.Game, mess.Version))
			r, err := obj.NewReader(ctx)
			if err != nil {
				log.Printf(err.Error())
				m.Nack()
				return
			}

			defer r.Close()
			path := fmt.Sprintf("/tmp/%v/%v/%v", int64(mess.AccountID), mess.Game, mess.Version)
			err = os.MkdirAll(path, 0777)
			if err != nil {
				log.Printf(err.Error())
				m.Nack()
				return
			}

			b, err := ioutil.ReadAll(r)
			if err != nil {
				log.Printf(err.Error())
				m.Nack()
				return
			}

			err = ioutil.WriteFile(path + "/server", b ,0755)
			if err != nil {
				log.Printf(err.Error())
				m.Nack()	
				return
			}

			// Run some docker commands to build and push the image.
			// ...
			d, err := ioutil.ReadFile("/Dockerfile.Template")
			if err != nil {
				log.Printf(err.Error())
				m.Nack()
				return
			}

			err = ioutil.WriteFile(path + "/Dockerfile", d, 0755)
			if err != nil {
				log.Printf(err.Error())
				m.Nack()
				return
			}

			// Docker commands

			_, err = exec.Command("docker", "login", "-e", "trubble-worker@google.com", "-u", "_json_key", "-p", "\"$(cat /keyfile.json)\"", "https://gcr.io").CombinedOutput()
			if err != nil {
				log.Printf(err.Error())
				m.Nack()
				return
			}

			_, err = exec.Command("docker", "build", "-t", fmt.Sprintf("trubble/%v-%v:%v", int64(mess.AccountID), mess.Game, mess.Version), path).CombinedOutput()
			if err != nil {
				log.Printf(err.Error())
				m.Nack()
				return
			}

			_, err = exec.Command("docker", "tag", fmt.Sprintf("trubble/%v-%v:%v", int64(mess.AccountID), mess.Game, mess.Version), 
			fmt.Sprintf("gcr.io/chronojam-trubble/%v-%v:%v", int64(mess.AccountID), mess.Game, mess.Version)).CombinedOutput()
			if err != nil {
				log.Printf(err.Error())
				m.Nack()
				return
			}

			_, err = exec.Command("docker", "push", fmt.Sprintf("gcr.io/chronojam-trubble/%v-%v:%v", int64(mess.AccountID), mess.Game, mess.Version)).CombinedOutput()
			if err != nil {
				log.Printf(err.Error())
				m.Nack()
				return
			}

			m.Ack()
		})
	}
}