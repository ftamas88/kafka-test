package app

import (
	"fmt"
	"github.com/ftamas88/kafka-test/internal/config"
	"github.com/ftamas88/kafka-test/internal/consumer"
	"github.com/ftamas88/kafka-test/internal/controller"
	"github.com/ftamas88/kafka-test/internal/domain"
	"github.com/ftamas88/kafka-test/internal/producer"
	"github.com/ftamas88/kafka-test/internal/routing"
	"github.com/ftamas88/kafka-test/internal/watcher"
	"os"
)

// New returns a new instance of the App
func New() (*App, error) {
	conf, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("unable to initialise config: %s", err.Error())
	}

	schemaData, err := os.ReadFile(conf.SchemaPath)
	if err != nil {
		return nil, fmt.Errorf("could not read avro schema file: %s", err)
	}

	incomingCh := make(chan domain.Payload, 1)
	cloudIncomingCh := make(chan domain.Payload, 1)
	outgoingCh := make(chan domain.Payload, 1)

	return &App{
		version:    version,
		commitHash: commitHash,
		consumer:   consumer.NewKafkaConsumer(conf),
		producer:   producer.NewKafkaProducer(conf, string(schemaData), incomingCh, outgoingCh),
		cloudProducer: producer.NewCloudKafkaProducer(conf, string(schemaData), cloudIncomingCh),
		watcher: watcher.NewFileWatcher(incomingCh, cloudIncomingCh, outgoingCh),
		router: routing.NewRouter(
			conf.HTTPPort,
			&routing.RouterConfig{
				WellKnown: &controller.WellKnownController{
					AppVersion:    Version(),
					AppCommitHash: CommitHash(),
				},
			},
		),
	}, nil
}

