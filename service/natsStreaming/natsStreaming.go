package natsStreaming

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"

	"service/internal/cache"
	"service/internal/http-server/model"
	"service/internal/storage/postgres"
)

func RunNatsStreaming(storage *postgres.Storage, cache *cache.Cache) error {
	const op = "RunNatsStreaming"

	sc, err := stan.Connect("test-cluster", "subscriber-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Printf("%s: %v", op, err)
		return err
	}
	defer func() {
		if err := sc.Close(); err != nil {
			log.Printf("%s: %v", op, err)
		}
	}()

	subject := "test"

	_, err = sc.Subscribe(subject, func(m *stan.Msg) {
		handleMessage(storage, cache, m)
	}, stan.StartWithLastReceived())
	if err != nil {
		log.Printf("%s: %v", op, err)
		return err
	}

	log.Println("Subscriber is listening...")

	select {}
	return nil
}

func handleMessage(storage *postgres.Storage, cache *cache.Cache, m *stan.Msg) {
	log.Printf("Received a message: %s\n", string(m.Data))

	var model model.Model
	err := json.Unmarshal(m.Data, &model)
	if err != nil {
		log.Println(err)
	}

	if err = storage.AddOrder(model.OrderUID, string(m.Data)); err != nil {
		log.Println(err)
	}

	cache.SetDefault(model.OrderUID, model)
}
