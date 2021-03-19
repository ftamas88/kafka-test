package producer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dangkaka/go-kafka-avro"
	"github.com/ftamas88/kafka-test/internal/config"
	"github.com/ftamas88/kafka-test/internal/domain"
)

type KafkaProducer struct {
	cfg        *config.Config
	schema     string
	incomingCh chan domain.Payload
	notifyCh   chan<- domain.Payload
}

func NewKafkaProducer(cfg *config.Config, s []byte, notify chan<- domain.Payload) *KafkaProducer {
	return &KafkaProducer{
		cfg:        cfg,
		schema:     string(s),
		incomingCh: make(chan domain.Payload, 1),
		notifyCh:   notify,
	}
}

func (k KafkaProducer) Run(ctx context.Context) error {
	log.Println("Producer running..")
	if err := createTopicFromAdmin(k.cfg, k.cfg.Topic); err != nil {
		log.Printf("could not create cloud topic: %s", err.Error())
		return fmt.Errorf("could not create cloud topic: %s", err)
	}

	producer, err := kafka.NewAvroProducer(k.cfg.Servers, k.cfg.SchemaRegistryServers)
	if err != nil {
		return fmt.Errorf("could not create avro producer: %s", err)
	}

	return k.Produce(producer)
}

func (k KafkaProducer) Produce(producer *kafka.AvroProducer) error {
	for {
		for msg := range k.incomingCh {
			k.addMsg(producer, msg)
		}
	}
}

func (k KafkaProducer) Ingest(msg domain.Payload) {
	k.incomingCh <- msg
}

func (k KafkaProducer) addMsg(producer *kafka.AvroProducer, msg domain.Payload) {
	key := time.Now().String()

	err := producer.Add(k.cfg.Topic, k.schema, []byte(key), msg.Data)
	if err != nil {
		log.Printf("Could not add a msg: [%s] msg: %s", err, msg.Filename)

		return
	}

	// instead of waiting for both cloud and regular producers to confirm
	// now only the internal cluster is enough
	defer func() {
		k.notifyCh <- msg
	}()
}
