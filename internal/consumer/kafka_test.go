package consumer_test

import (
	"fmt"
	"testing"

	"github.com/ftamas88/kafka-test/internal/consumer/mocks"
	"github.com/stretchr/testify/assert"
)

func TestKafkaConsumer_Run(t *testing.T) {
	type fields struct {
		service mocks.Consumer
	}
	tests := []struct {
		name   string
		setup  func(*fields)
		assert func(*testing.T, *fields, error)
	}{
		{
			name: "success - no error",
			setup: func(fields *fields) {
				fields.service.On("Consume").
					Once().
					Return(nil)
			},
			assert: func(t *testing.T, fields *fields, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "failure",
			setup: func(fields *fields) {
				fields.service.On("Consume").
					Once().
					Return(fmt.Errorf("test error"))
			},
			assert: func(t *testing.T, fields *fields, err error) {
				assert.Equal(t, "test error", err.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := &fields{
				service: mocks.Consumer{},
			}

			tt.setup(fields)
			err := fields.service.Consume()
			tt.assert(t, fields, err)
			fields.service.AssertExpectations(t)
		})
	}
}
