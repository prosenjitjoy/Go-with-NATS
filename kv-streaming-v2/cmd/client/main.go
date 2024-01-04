package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"main/model"
	"main/util"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

func main() {
	var genreName string
	flag.StringVar(&genreName, "genre", "", "pass the genre name")
	flag.Parse()

	cfg, err := util.LoadConfig(".env")
	fatalOnErr(err)

	nc, err := nats.Connect(cfg.NatsUrl)
	fatalOnErr(err)
	defer nc.Close()

	js, err := nc.JetStream()
	fatalOnErr(err)

	// list all genres in the specified stream
	if genreName == "" {
		genres := getAllGenres(js)

		for i, s := range genres {
			fmt.Printf("%v - %v\n", i+1, s)
		}
		return
	}

	bucket, err := js.KeyValue(genreName)
	fatalOnErr(err)

	bucketWatcher, err := bucket.Watch("*")
	fatalOnErr(err)

	go playSongs(genreName, bucketWatcher.Updates())

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

func getAllGenres(js nats.JetStreamContext) []string {
	var genres []string

	for storeName := range js.KeyValueStoreNames() {
		// strip out the prefix "KV_"
		storeName = storeName[3:]
		genres = append(genres, storeName)
	}

	return genres
}

func playSongs(genre string, ch <-chan nats.KeyValueEntry) {
	log.Println("listening for songs, genre:", genre)

	var track model.Track

	fmt.Println("Songs:")
	for entry := range ch {
		// on first connect, entry is nil
		if entry == nil {
			continue
		}

		if err := json.Unmarshal(entry.Value(), &track); err != nil {
			continue
		}

		fmt.Printf("\t'%s' by %v\n", track.Title, track.Artist)
	}
}
