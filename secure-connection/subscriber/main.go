package main

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("tls://127.0.0.1:4222", nats.UserInfo("a", "a"))
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	// Subscribe to the "foo" subject and print all message
	nc.Subscribe("foo", func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
	})

	// Keep the application running
	select {}
}
