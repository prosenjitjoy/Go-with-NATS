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

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("failed to create nats jetstream:", err)
	}

	// js.Subscribe("orders.us", processMsg, nats.BindStream("ORDERS"))
	_, err = js.AddConsumer("ORDERS", &nats.ConsumerConfig{
		Durable:      "my-consumer-1",
		Description:  "this is my awesome consumer",
		ReplayPolicy: nats.ReplayInstantPolicy,
	})
	if err != nil {
		log.Fatal("failed to add stream consumer:", err)
	}

	sub, err := js.PullSubscribe("orders.us", "my")
	if err != nil {
		log.Fatal("failed to create pull subscription:", err)
	}

	go processMsg(sub)

	time.Sleep(10 * time.Second)
	sub.Unsubscribe()
	log.Println("shutting down application...")
}

func processMsg(sub *nats.Subscription) error {
	for sub.IsValid() {
		msgs, err := sub.Fetch(1)
		if err != nil {
			return err
		}
		fmt.Printf("INFO - Got reply: %s\n", string(msgs[0].Data))
	}
	return nil
}
