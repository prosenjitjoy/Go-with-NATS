package main

import (
	"encoding/json"
	"fmt"
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

	sub, err := nc.QueueSubscribe("intros", "zip1", processMsg)
	if err != nil {
		log.Fatal("can't subscribe to nats queue 'zigp1':", err)
	}
	defer sub.Unsubscribe()
	time.Sleep(1 * time.Hour)
}

func processMsg(msg *nats.Msg) {
	pl := &model.Payload{}
	json.Unmarshal(msg.Data, pl)
	reply := fmt.Sprintf("ack message # %v", pl.Count)
	msg.Respond([]byte(reply))

	fmt.Printf("I got a message: %s, count: %v\n", pl.Data, pl.Count)
}
