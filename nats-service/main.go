package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		slog.Error("unable to connect to NATS server", "err", err)
		os.Exit(1)
	}
	defer nc.Close()

	greeterConfig := micro.Config{
		Name:        "greeter",
		Version:     "0.0.1",
		Description: "A simple greeter service",
	}

	greeterService, err := micro.AddService(nc, greeterConfig)
	if err != nil {
		slog.Error("unable to create service", "err", err)
		os.Exit(1)
	}

	err = greeterService.AddEndpoint("hello", micro.HandlerFunc(helloHandler))
	if err != nil {
		slog.Error("unable to create endpoint", "err", err)
		os.Exit(1)
	}

	// keep the program running forever
	select {}
}

func helloHandler(r micro.Request) {
	data := string(r.Data())
	if data == "" {
		data = "world"
	}

	msg := fmt.Sprintf("Hello, %s!", data)
	r.Respond([]byte(msg))
}
