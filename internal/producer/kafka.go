package producer

import (
	"context"
	"fmt"
	"github.com/dangkaka/go-kafka-avro"
	"github.com/ftamas88/kafka-test/internal/config"
	"github.com/ftamas88/kafka-test/internal/domain"
	"log"
	"time"
)

type KafkaProducer struct {
	cfg *config.Config
	schema string
	incomingCh <- chan domain.Payload
	notifyCh chan<- domain.Payload
}

func NewKafkaProducer(cfg *config.Config, s string, ch <- chan domain.Payload, oCh chan<- domain.Payload) *KafkaProducer {
	return &KafkaProducer{
		cfg: cfg,
		schema: s,
		incomingCh: ch,
		notifyCh: oCh,
	}
}

func (k KafkaProducer) Run(ctx context.Context) error {
	log.Println("Producer running..")

	producer, err := kafka.NewAvroProducer(k.cfg.Servers, k.cfg.SchemaRegistryServers)
	if err != nil {
		return fmt.Errorf("could not create avro producer: %s", err)
	}

	if err := createTopicFromAdmin(k.cfg, "test"); err != nil {
		log.Printf("could not create cloud topic: %s", err.Error())
		return fmt.Errorf("could not create cloud topic: %s", err)
	}

	go func() {
		for {
			select {
			case msg := <-k.incomingCh:
				k.addMsg(producer, msg)
			}
		}
	}()

	return nil
}

func (k KafkaProducer) addMsg(producer *kafka.AvroProducer, msg domain.Payload) {
	defer func() {
		k.notifyCh <- msg
	}()
	key := time.Now().String()

	err := producer.Add("test", k.schema, []byte(key), msg.Data)
	if err != nil {
		msg.Status = false
		log.Printf("Could not add a msg: [%s] msg: %s", err, msg.Filename)

		return
	}

	msg.Status = true
}