package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"main/model"
	"main/util"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	var filePath string
	var streamName string
	flag.StringVar(&filePath, "path", "data.json", "path to the JSON playlist data file")
	flag.StringVar(&streamName, "stream", "songs", "stream name")
	flag.Parse()

	playlists, err := loadPlaylists(filePath)
	fatalOnErr(err)

	cfg, err := util.LoadConfig(".env")
	fatalOnErr(err)

	nc, err := nats.Connect(cfg.NatsUrl)
	fatalOnErr(err)
	defer nc.Close()

	js, err := nc.JetStream()
	fatalOnErr(err)

	// create JetStream kv store bucket call songs
	songBucket, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: streamName,
	})
	if err != nil {
		songBucket, err = js.KeyValue(streamName)
		if err != nil {
			log.Fatalf("error creating or accessing '%v' stream bucket: %v\n", streamName, err)
		}
	}

	for genre, tracks := range playlists {
		log.Printf("creating genre, stream-bucket: %v, genre: %v", streamName, genre)
		go streamPlaylist(songBucket, genre, tracks)
	}

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

func loadPlaylists(filepath string) (model.Playlists, error) {
	var playlists model.Playlists
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening playlist data file: %v", err)
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&playlists); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return playlists, nil
}

func streamPlaylist(jsKeyVal nats.KeyValue, playlist string, tracks []model.Track) {
	numTracks := len(tracks)

	for {
		// select a random track
		trackIndex := r.Intn(numTracks)
		track := tracks[trackIndex]

		log.Printf("playing track playlist: %v, artist: %v, title: %v\n", playlist, track.Artist, track.Title)

		// store track to the playlist key
		jsKeyVal.Put(playlist, track.ToJsonBytes())
		trackDuration := 5 + time.Duration(r.Intn(5))
		time.Sleep(trackDuration * time.Second)
	}
}
