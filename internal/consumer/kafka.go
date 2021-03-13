package consumer

import (
	"context"
	"fmt"
	"github.com/bsm/sarama-cluster"
	"github.com/dangkaka/go-kafka-avro"
	"github.com/ftamas88/kafka-test/internal/config"
	"log"
)

type KafkaConsumer struct {
	cfg *config.Config
}

func NewKafkaConsumer(cfg *config.Config) *KafkaConsumer {
	return &KafkaConsumer{
		cfg: cfg,
	}
}

func (k KafkaConsumer) Run(ctx context.Context) error {
	log.Printf("servers: %s", k.cfg.Servers)

	consumerCallbacks := kafka.ConsumerCallbacks{
		OnDataReceived: func(msg kafka.Message) {
			fmt.Println("Message: " + msg.Value)
		},
		OnError: func(err error) {
			fmt.Println("Consumer error", err)
		},
		OnNotification: func(notification *cluster.Notification) {
			fmt.Println(notification)
		},
	}

	c, err := kafka.NewAvroConsumer(
		k.cfg.Servers,
		k.cfg.SchemaRegistryServers,
		"test",
		"consumer-group",
		consumerCallbacks,
	)
	if err != nil {
		return fmt.Errorf("error during consuming data: %s", err.Error())
	}

	c.Consume()

	return nil
}
