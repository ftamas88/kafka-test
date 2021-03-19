package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	cluster "github.com/bsm/sarama-cluster"
	"github.com/dangkaka/go-kafka-avro"
	"github.com/ftamas88/kafka-test/internal/config"
)

type KafkaConsumer struct {
	cfg *config.Config
}

// NewKafkaConsumer creates a new instance of the Kafka consumer
func NewKafkaConsumer(cfg *config.Config) *KafkaConsumer {
	return &KafkaConsumer{
		cfg: cfg,
	}
}

// Consume starts consuming from a Kafka topic
func (k KafkaConsumer) Consume() error {
	log.Println("Consumer running..")

	consumerCallbacks := kafka.ConsumerCallbacks{
		OnDataReceived: func(msg kafka.Message) {
			// Just for better readability
			var prettyJSON bytes.Buffer
			_ = json.Indent(&prettyJSON, []byte(msg.Value), "", "\t")
			fmt.Printf("Message: [%s]\n", prettyJSON.String())
		},
		OnError: func(err error) {
			fmt.Println("Consumer error", err)
		},
		OnNotification: func(notification *cluster.Notification) {
			fmt.Println(notification)
		},
	}

	cons, err := kafka.NewAvroConsumer(
		k.cfg.Servers,
		k.cfg.SchemaRegistryServers,
		k.cfg.Topic,
		"consumer-group",
		consumerCallbacks,
	)
	if err != nil {
		return fmt.Errorf("error during consuming data: %s", err.Error())
	}

	cons.Consume()

	return nil
}

// Run runs the Kafka consumer component
func (k KafkaConsumer) Run(ctx context.Context) error {
	return k.Consume()
}
