package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/nats-io/stan.go"
)

func main() {
	path := "../models"
	channel := "test"

	sc, err := stan.Connect("test-cluster", "publisher-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range files {
		msg, err := readData(path + "/" + v.Name())
		if err != nil {
			return
		}
		if err = sc.Publish(channel, msg); err != nil {
			log.Fatalf("Error publishing: %v", err)
		} else {
			log.Printf("Published from: %s", v.Name())
		}
	}

	log.Println("All messages published")
}

func readData(path string) ([]byte, error) {
	file, err := os.Open(path)
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
