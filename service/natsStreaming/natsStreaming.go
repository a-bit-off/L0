package natsStreaming

import (
	"context"
	"log"

	"github.com/nats-io/stan.go"
)

func RunNatsStreaming(ctx context.Context) error {
	const op = "RunNatsStreaming"

	sc, err := stan.Connect("test-cluster", "subscriber-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		return err
	}
	defer sc.Close()

	subject := "test"

	_, err = sc.Subscribe(subject, func(m *stan.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
	})
	if err != nil {
		log.Printf("Error subscribing: %s: %v", op, err)
		return err
	}

	log.Println("Subscriber is listening...")

	<-ctx.Done() // Ожидание сигнала завершения
	log.Println("Stopping subscriber...")
	return nil
}
