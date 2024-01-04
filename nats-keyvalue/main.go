package main

import (
	"fmt"
	"log"
	"main/util"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

func main() {
	cfg, err := util.LoadConfig(".env")
	fatalOnErr(err)

	nc, err := nats.Connect(cfg.NatsUrl)
	fatalOnErr(err)
	defer nc.Close()

	js, err := nc.JetStream()
	fatalOnErr(err)

	// list all buckets and their keys
	for bucketName := range js.KeyValueStoreNames() {
		// trim first 3 characters from bucketName
		bucketName = bucketName[3:]

		fmt.Printf("Bucket - %s\n", bucketName)

		kvBucket, err := js.KeyValue(bucketName)
		if err != nil {
			log.Println(err)
			continue
		}

		keys, err := kvBucket.Keys()
		if err != nil {
			log.Println(err)
			continue
		}

		for _, key := range keys {
			fmt.Printf("\t%s\n", key)
		}
	}

	// create a new bucket call "sensors"
	bucketName := "sensor"
	sensors, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: bucketName,
	})
	fatalOnErr(err)

	sensors.PutString("temperature", "48deg")
	sensors.PutString("humidity", "50%")
	sensors.PutString("pressure", "10bars")

	monitor, err := sensors.Watch("*")
	fatalOnErr(err)

	go handleReadings(monitor.Updates())

	// cleanly exit application if signal is caught
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Println("INFO - Exiting on signal")
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatal("fatal error:", err)
	}
}

func handleReadings(ch <-chan nats.KeyValueEntry) {
	for entry := range ch {
		if entry != nil {
			log.Printf("new reading sensor %v, value %s", entry.Key(), string(entry.Value()))
		}
	}
}
