package producer

import "github.com/ftamas88/kafka-test/internal/domain"

//go:generate mockery --name=Producer
type Producer interface {
	Ingest(msg domain.Payload)
}
