package producer

import (
	"context"
	"fmt"
	"github.com/ftamas88/kafka-test/internal/config"
	"github.com/ftamas88/kafka-test/internal/domain"
	ka "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"time"
)

type CloudKafkaProducer struct {
	cfg *config.Config
	schema string
	incomingCh <- chan domain.Payload
}

func NewCloudKafkaProducer(cfg *config.Config, s string, ch <- chan domain.Payload) *CloudKafkaProducer {
	return &CloudKafkaProducer{
		cfg: cfg,
		schema: s,
		incomingCh: ch,
	}
}

func (k CloudKafkaProducer) Run(ctx context.Context) error {
	if k.cfg.CloudServer == "" {
		return nil
	}

	log.Printf("Cloud Producer running.. %+v", k.cfg.CloudServer)

	producer, err := ka.NewProducer(&ka.ConfigMap{
		"bootstrap.servers":       k.cfg.CloudServer,
		"security.protocol":       "SASL_SSL",
		"sasl.mechanisms":         "PLAIN",
		"sasl.username":           k.cfg.KafkaCloudKey,
		"sasl.password":           k.cfg.KafkaCloudSecret})

	if err != nil {
		log.Printf("cloud not create cloud producer: %s", err.Error())
		return fmt.Errorf("could not create cloud producer: %s", err)
	}

	if err := createTopic(producer, "test"); err != nil {
		log.Printf("could not create cloud topic: %s", err.Error())
		return fmt.Errorf("could not create cloud topic: %s", err)
	}

	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			select {
			case msg := <-k.incomingCh:
				log.Printf("received message to publish for cloud: %s", msg.Filename)
				k.addMsg(ctx, producer, msg)
			}
		}
	}()

	return nil
}

func (k CloudKafkaProducer) addMsg(ctx context.Context, producer *ka.Producer, msg domain.Payload) {
	key := time.Now().String()

	if err := producer.BeginTransaction(); err != nil {
		return
	}

	topic := "test"
	deliveryChan := make(chan ka.Event)

	if err := producer.Produce(&ka.Message{
		TopicPartition: ka.TopicPartition{Topic: &topic, Partition: ka.PartitionAny},
		Value:          msg.Data,
		Headers:        []ka.Header{{Key: key, Value: []byte("header values are binary")}},
	}, deliveryChan); err != nil {
		log.Printf("Failed to produce message: %s\n", err.Error())
		return
	}

	e := <-deliveryChan
	m := e.(*ka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		log.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	if err := producer.CommitTransaction(ctx); err != nil {
		log.Printf("failed to commit trx: %s", err.Error())
		return
	}

	log.Printf("message published tot he cloud :%s\n", m.Key)
}