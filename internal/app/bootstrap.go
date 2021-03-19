package app

import (
	"context"
	"fmt"
	"os"

	"github.com/ftamas88/kafka-test/internal/config"
	"github.com/ftamas88/kafka-test/internal/consumer"
	"github.com/ftamas88/kafka-test/internal/controller"
	"github.com/ftamas88/kafka-test/internal/domain"
	"github.com/ftamas88/kafka-test/internal/producer"
	"github.com/ftamas88/kafka-test/internal/routing"
	"github.com/ftamas88/kafka-test/internal/watcher"
	"github.com/ftamas88/kafka-test/internal/worker"
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
	outgoingCh := make(chan domain.Payload, 1)

	quitCh := make(chan struct{})

	d := worker.NewDistributor(incomingCh, quitCh)
	d.AddWorker(producer.NewCloudKafkaProducer(conf, schemaData))
	d.AddWorker(producer.NewKafkaProducer(conf, schemaData, outgoingCh))

	err = d.Start(context.Background())
	d.Listen()

	return &App{
		version:    version,
		commitHash: commitHash,
		consumer:   consumer.NewKafkaConsumer(conf),
		watcher:    watcher.NewFileWatcher(incomingCh, outgoingCh),
		router: routing.NewRouter(
			conf.HTTPPort,
			&routing.RouterConfig{
				WellKnown: &controller.WellKnownController{
					AppVersion:    Version(),
					AppCommitHash: CommitHash(),
				},
			},
		),
	}, err
}
