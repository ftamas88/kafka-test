package worker

import (
	"context"

	"github.com/ftamas88/kafka-test/internal/domain"
)

//go:generate mockery --name=Worker
type Worker interface {
	Run(ctx context.Context) error
	Ingest(payload domain.Payload)
}
