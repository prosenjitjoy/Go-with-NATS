package main

import (
	"fmt"
	"log"
	"main/util"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	cfg, err := util.LoadConfig(".env")
	if err != nil {
		log.Fatal("failed to load configuration file:", err)
	}

	nc, err := nats.Connect(cfg.NatsUrl)
	if err != nil {
		log.Fatal("failed to connect to nats server:", err)
	}
	defer nc.Close()

	for i := 0; i < 10; i++ {
		msg, err := nc.Request("time.of.day", []byte(""), cfg.Timeout)
		if err != nil {
			log.Fatal("failed to send payload request:", err)
		}

		fmt.Printf("INFO - Got reply - %s\n", string(msg.Data))
		time.Sleep(time.Second)
	}
}
