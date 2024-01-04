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
	flag.StringVar(&filePath, "path", "data.json", "path to the JSON playlist data file")
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

	for genre, tracks := range playlists {
		log.Printf("creating genre, stream-bucket: %v, genre: %v", genre, genre)
		go streamPlaylist(js, genre, tracks)
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

func streamPlaylist(js nats.JetStreamContext, genre string, tracks []model.Track) {
	genreBucket, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: genre,
		TTL:    1 * time.Minute,
	})
	if err != nil {
		genreBucket, err = js.KeyValue(genre)
		if err != nil {
			log.Fatalf("error creating or accessing stream bucket, genre: %v, error: %v", genre, err)
			return
		}
	}

	// post the first track on startup, then delay 5 seconds adding new tracks every 10 seconds
	track := getNextTrack(tracks)
	genreBucket.Put(genre+"-"+track.Id, track.ToJsonBytes())
	time.Sleep(5 * time.Second)

	for {
		track := getNextTrack(tracks)
		if _, err := genreBucket.Create(genre+"-"+track.Id, track.ToJsonBytes()); err != nil {
			continue
		}

		log.Printf("playing track, genre: %v, artist: %s, title: %s\n", genre, track.Artist, track.Title)
	}
}

func getNextTrack(tracks []model.Track) model.Track {
	numTracks := len(tracks)
	trackIndex := r.Intn(numTracks)
	track := tracks[trackIndex]
	return track
}
