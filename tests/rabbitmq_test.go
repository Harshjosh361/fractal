package tests

import (
	"fmt"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAMQPChannel simulates an AMQP channel for testing
type MockAMQPChannel struct {
	mock.Mock
}

func (m *MockAMQPChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	argsCall := m.Called(name, durable, autoDelete, exclusive, noWait, args)
	return argsCall.Get(0).(amqp.Queue), argsCall.Error(1)
}

func (m *MockAMQPChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	argsCall := m.Called(exchange, key, mandatory, immediate, msg)
	return argsCall.Error(0)
}

func (m *MockAMQPChannel) Close() error {
	argsCall := m.Called()
	return argsCall.Error(0)
}

// New method to mock Consume functionality
func (m *MockAMQPChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	argsCall := m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	return argsCall.Get(0).(<-chan amqp.Delivery), argsCall.Error(1)
}

func TestProduceAndConsumeSingleItem(t *testing.T) {
	const (
		GreenKafkaTick = "\033[32m✔\033[0m" // Green tick
		RedKafkaCross  = "\033[31m✘\033[0m" // Red cross
	)
	
	mockChannel := new(MockAMQPChannel)

	// Simulate producing a message
	mockChannel.On("Publish", "", "test-queue", false, false, mock.Anything).Return(nil)

	// Simulate consuming the same message
	expectedMessage := []byte("test-message")
	mockChannel.On("Consume", "test-queue", "", true, false, false, false, amqp.Table{}).Return(generateDeliveryChannel(expectedMessage), nil)

	// Produce data
	err := mockChannel.Publish("", "test-queue", false, false, amqp.Publishing{
		Body: expectedMessage,
	})

	if assert.NoError(t, err) {
		fmt.Printf("%v Produced messages successfully\n", GreenKafkaTick)
	}

	// Consume data
	msgs, err := mockChannel.Consume("test-queue", "", true, false, false, false, amqp.Table{})
	if assert.NoError(t, err) {
		for msg := range msgs {
			assert.Equal(t, expectedMessage, msg.Body)
			fmt.Printf("%v Consumed messages successfully\n", GreenKafkaTick)
			break
		}
	}

	// Verify expectations
	mockChannel.AssertExpectations(t)
}

// Helper function to simulate delivery channel
func generateDeliveryChannel(expectedMessage []byte) <-chan amqp.Delivery {
	ch := make(chan amqp.Delivery, 1)
	ch <- amqp.Delivery{Body: expectedMessage}
	close(ch)
	return ch
}
