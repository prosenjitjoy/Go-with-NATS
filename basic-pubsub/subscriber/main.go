package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("can't connect to nats:", err)
	}
	defer nc.Close()

	nc.Subscribe("intros", func(msg *nats.Msg) {
		fmt.Printf("I got a message: %s\n", string(msg.Data))
	})

	time.Sleep(1 * time.Hour)
}
