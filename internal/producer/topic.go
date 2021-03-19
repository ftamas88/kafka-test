package producer

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/ftamas88/kafka-test/internal/config"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// createTopicFromAdmin creates a topic using the Admin Client API
func createTopicFromAdmin(cfg *config.Config, topic string) error {
	a, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": strings.Join(cfg.Servers, ",")})
	if err != nil {
		log.Printf("Failed to create new admin client: %s", err)
		return err
	}
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create topics on cluster.
	// Set Admin options to wait up to 60s for the operation to finish on the remote cluster
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		log.Printf("ParseDuration(60s): %s", err)
		return err
	}
	results, err := a.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 2}},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		log.Printf("Admin Client request error: %v\n", err)
		return err
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			log.Printf("Failed to create topic: %v\n", result.Error)
			return err
		}
		log.Printf("%v\n", result)
	}
	a.Close()

	return nil
}

// createTopic creates a topic using the Admin Client API
func createTopic(p *kafka.Producer, topic string) error {

	a, err := kafka.NewAdminClientFromProducer(p)
	if err != nil {
		log.Printf("Failed to create new admin client from producer: %s", err)
		return err
	}
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create topics on cluster.
	// Set Admin options to wait up to 60s for the operation to finish on the remote cluster
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		log.Printf("ParseDuration(60s): %s", err)
		return err
	}
	results, err := a.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 2}},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		log.Printf("Admin Client request error: %v\n", err)
		return err
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			log.Printf("Failed to create topic: %v\n", result.Error)
			return err
		}
		log.Printf("%v\n", result)
	}
	a.Close()

	return nil
}
