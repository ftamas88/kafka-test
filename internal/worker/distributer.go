package worker

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ftamas88/kafka-test/internal/domain"
	"golang.org/x/sync/errgroup"
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
	eg, egCtx := errgroup.WithContext(ctx)

	for k, worker := range d.workers {
		fmt.Printf("Running Worker [%d][%s]\n", k, reflect.TypeOf(worker).String())

		go func(w Worker) {
			eg.Go(func() error {
				return w.Run(egCtx)
			})
		}(worker)

	}

	return eg.Wait()
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
