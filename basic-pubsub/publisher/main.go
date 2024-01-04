package main

import (
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

	for i := 0; ; i++ {
		log.Println("sent message", i)
		nc.Publish("intros", []byte("Hello World!"))
		time.Sleep(5 * time.Second)
	}
}
