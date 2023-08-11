package main

import (
	"log"

	"github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("test-cluster", "publisher-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	subject := "test-subject"
	msg := []byte("Hello, NATS Streaming!")

	if err := sc.Publish(subject, msg); err != nil {
		log.Fatalf("Error publishing: %v", err)
	}

	log.Printf("Published: %s", msg)
}
