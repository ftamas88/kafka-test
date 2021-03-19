package consumer

//go:generate mockery --name=Consumer
type Consumer interface {
	Consume() error
}
