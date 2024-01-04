package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

const putFileHelp = "put <store> <file>"
const getFileHelp = "get <store> <file> <destination>"

var (
	username string
	password string
	hostname = "localhost"
	port     = 4222
)

func init() {
	flag.StringVar(&username, "user", username, "username for NATS server")
	flag.StringVar(&password, "pass", password, "password for NATS server")
	flag.StringVar(&hostname, "host", hostname, "NATS server hostname")
	flag.IntVar(&port, "port", port, "NATS server port")
	flag.Parse()
}

func main() {
	if len(flag.Args()) == 0 {
		log.Fatalf("missing arguments, user either %s or %s\n", putFileHelp, getFileHelp)
	}

	natsUrl := fmt.Sprintf("nats://%v:%v", hostname, port)
	if username != "" {
		natsUrl = fmt.Sprintf("nats://%v:%v@%v:%v", username, password, hostname, port)
	}

	nc, err := nats.Connect(natsUrl)
	fatalOnErr(err)

	js, err := nc.JetStream()
	fatalOnErr(err)

	operation := flag.Arg(0)
	if operation == "" {
		log.Fatalf("missing operation, user either %s or %s\n", putFileHelp, getFileHelp)
	}

	objStoreName := flag.Arg(1)
	if objStoreName == "" {
		log.Fatalf("missing objStoreName, user either %s or %s\n", putFileHelp, getFileHelp)
	}

	objStore, err := js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket: objStoreName,
	})
	if err != nil {
		objStore, err = js.ObjectStore(objStoreName)
		fatalOnErr(err)
	}

	fileName := flag.Arg(2)
	if fileName == "" {
		log.Fatalf("missing fileName, user either %s or %s\n", putFileHelp, getFileHelp)
	}

	switch operation {
	case "put":
		_, err := objStore.PutFile(fileName)
		fatalOnErr(err)
	case "get":
		err := objStore.GetFile(fileName, flag.Arg(3))
		fatalOnErr(err)
	}
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatal("fatal error:", err)
	}
}
