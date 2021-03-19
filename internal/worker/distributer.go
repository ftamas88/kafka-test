package worker

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/ftamas88/kafka-test/internal/domain"
)

type distributor struct {
	ingest chan domain.Payload
	quit   chan struct{}

	workers []Worker
}

func NewDistributor(ingest chan domain.Payload, q chan struct{}) *distributor {
	return &distributor{
		ingest: ingest,
		quit:   q,
	}
}

func (d *distributor) AddWorker(w Worker) {
	d.workers = append(d.workers, w)
}

func (d *distributor) Start(ctx context.Context) error {
	for k, worker := range d.workers {
		fmt.Printf("Running Worker [%d][%s]\n", k, reflect.TypeOf(worker).String())

		go func(w Worker) {
			if err := w.Run(ctx); err != nil {
				log.Fatalf("unable to run worker: %s", err.Error())
			}
		}(worker)

	}

	return nil
}

func (d *distributor) Listen() {
	go func() {
		for {
			select {
			case p := <-d.ingest:
				for _, worker := range d.workers {
					go worker.Ingest(p)
				}
			case <-d.quit:
				return
			}
		}
	}()
}
