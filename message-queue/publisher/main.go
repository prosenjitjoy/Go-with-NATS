package main

import (
	"encoding/json"
	"log"
	"main/model"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("can't connect to nats:", err)
	}
	defer nc.Close()

	count := 0
	pl := &model.Payload{
		Data: "Hello World!",
	}

	for {
		pl.Count = count
		data, _ := json.Marshal(pl)
		reply, err := nc.Request("intros", data, 500*time.Millisecond)
		time.Sleep(5 * time.Second)
		if err != nil {
			log.Printf("error sending message count: %v, err: %v", count, err)
			continue
		}
		count++
		log.Printf("send message %v, got replay %v", count, string(reply.Data))
	}
}
