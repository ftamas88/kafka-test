// Code generated by mockery v2.6.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/ftamas88/kafka-test/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// Producer is an autogenerated mock type for the Producer type
type Producer struct {
	mock.Mock
}

// Ingest provides a mock function with given fields: msg
func (_m *Producer) Ingest(msg domain.Payload) {
	_m.Called(msg)
}