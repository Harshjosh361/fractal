package tests

import (
	"fmt"
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	GreenKafkaTick = "\033[32m✔\033[0m" // Green tick
	RedKafkaCross  = "\033[31m✘\033[0m" // Red cross
)

// MockKafkaReader simulates a Kafka reader for testing
type MockKafkaReader struct {
	mock.Mock
}

func (m *MockKafkaReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).(kafka.Message), args.Error(1)
}

func (m *MockKafkaReader) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockKafkaWriter simulates a Kafka writer for testing
type MockKafkaWriter struct {
	mock.Mock
}

func (m *MockKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	args := m.Called(ctx, msgs)
	return args.Error(0)
}

func (m *MockKafkaWriter) Close() error {
	args := m.Called()
	return args.Error(0)
}

// KafkaSource handles data consumption from Kafka
type KafkaSource struct {
	Reader KafkaReader
}

// KafkaDestination handles data production to Kafka
type KafkaDestination struct {
	Writer KafkaWriter
}

// KafkaReader defines the interface for a Kafka reader
type KafkaReader interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

// KafkaWriter defines the interface for a Kafka writer
type KafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// FetchData reads a message from the KafkaSource
func (k *KafkaSource) FetchData() (string, error) {
	msg, err := k.Reader.FetchMessage(context.Background())
	if err != nil {
		return "", err
	}
	return string(msg.Value), nil
}

// SendData sends a message to the KafkaDestination
func (k *KafkaDestination) SendData(data string) error {
	msg := kafka.Message{Value: []byte(data)}
	return k.Writer.WriteMessages(context.Background(), msg)
}

// TestKafkaSource_FetchData tests the FetchData method of KafkaSource
func TestKafkaSource_FetchData(t *testing.T) {
	mockReader := new(MockKafkaReader)
	kafkaSource := &KafkaSource{Reader: mockReader}

	// Prepare mock behavior
	expectedMessage := kafka.Message{Value: []byte("test-message")}
	mockReader.On("FetchMessage", mock.Anything).Return(expectedMessage, nil)

	// Fetch data
	result, err := kafkaSource.FetchData()

	// Assertions
	if assert.NoError(t, err) {
		fmt.Printf("%s FetchData passed\n", GreenKafkaTick)
	} else {
		fmt.Printf("%s FetchData failed\n", RedKafkaCross)
	}

	if assert.Equal(t, "test-message", result) {
		fmt.Printf("%s FetchData returned the expected message\n", GreenKafkaTick)
	} else {
		fmt.Printf("%s FetchData did not return the expected message\n", RedKafkaCross)
	}

	mockReader.AssertExpectations(t)
}

// TestKafkaDestination_SendData tests the SendData method of KafkaDestination
func TestKafkaDestination_SendData(t *testing.T) {
	mockWriter := new(MockKafkaWriter)
	kafkaDestination := &KafkaDestination{Writer: mockWriter}

	// Prepare mock behavior
	mockWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)

	// Send data
	err := kafkaDestination.SendData("test-message")

	// Assertions
	if assert.NoError(t, err) {
		fmt.Printf("%s SendData passed\n", GreenKafkaTick)
	} else {
		fmt.Printf("%s SendData failed\n", RedKafkaCross)
	}
	mockWriter.AssertExpectations(t)
}

// TestKafkaReader_Close tests the Close method of KafkaReader
func TestKafkaReader_Close(t *testing.T) {
	mockReader := new(MockKafkaReader)

	// Prepare mock behavior
	mockReader.On("Close").Return(nil)

	// Close the reader
	err := mockReader.Close()

	// Assertions
	if assert.NoError(t, err) {
		fmt.Printf("%s Close passed\n", GreenKafkaTick)
	} else {
		fmt.Printf("%s Close failed\n", RedKafkaCross)
	}
	mockReader.AssertExpectations(t)
}

// TestKafkaWriter_Close tests the Close method of KafkaWriter
func TestKafkaWriter_Close(t *testing.T) {
	mockWriter := new(MockKafkaWriter)

	// Prepare mock behavior
	mockWriter.On("Close").Return(nil)

	// Close the writer
	err := mockWriter.Close()

	// Assertions
	if assert.NoError(t, err) {
		fmt.Printf("%s Close passed\n", GreenKafkaTick)
	} else {
		fmt.Printf("%s Close failed\n", RedKafkaCross)
	}	
	mockWriter.AssertExpectations(t)
}
