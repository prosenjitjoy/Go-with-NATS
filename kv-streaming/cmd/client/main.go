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
	var streamName string
	var genreName string
	flag.StringVar(&streamName, "stream", "songs", "pass the stream name")
	flag.StringVar(&genreName, "genre", "", "pass the genre name")
	flag.Parse()

	cfg, err := util.LoadConfig(".env")
	fatalOnErr(err)

	nc, err := nats.Connect(cfg.NatsUrl)
	fatalOnErr(err)
	defer nc.Close()

	js, err := nc.JetStream()
	fatalOnErr(err)

	// list all streams available if none was specified
	if streamName == "" {
		streams := getAllStreams(js)

		for i, s := range streams {
			fmt.Printf("%v - %v\n", i+1, s)
		}
		return
	}

	// list all genres in the specified stream
	if genreName == "" {
		genres := getAllGenres(js, streamName)

		for i, s := range genres {
			fmt.Printf("%v - %v\n", i+1, s)
		}
		return
	}

	bucket, err := js.KeyValue(streamName)
	fatalOnErr(err)

	genreKey, err := bucket.Watch(genreName)
	fatalOnErr(err)

	go playSongs(genreName, genreKey.Updates())

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

func getAllStreams(js nats.JetStreamContext) []string {
	var streams []string

	for storeName := range js.KeyValueStoreNames() {
		// strip out the prefix "KV_"
		storeName = storeName[3:]
		streams = append(streams, storeName)
	}

	return streams
}

func getAllGenres(js nats.JetStreamContext, stream string) []string {
	bucket, err := js.KeyValue(stream)
	if err != nil {
		return nil
	}

	genres, _ := bucket.Keys()
	return genres
}

func playSongs(genre string, ch <-chan nats.KeyValueEntry) {
	log.Println("listening for songs, genre:", genre)

	var track model.Track

	for entry := range ch {
		// on first connect, entry is nil
		if entry == nil {
			continue
		}

		if err := json.Unmarshal(entry.Value(), &track); err != nil {
			continue
		}

		fmt.Printf("Playing song '%s' by %v\n", track.Title, track.Artist)
	}
}
