package integrations

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
	"github.com/segmentio/kafka-go"
)

// KafkaSource struct represents the configuration for consuming messages from Kafka.
type KafkaSource struct {
	URL   string `json:"consumer_url"`
	Topic string `json:"consumer_topic"`
}

// KafkaDestination struct represents the configuration for publishing messages to Kafka.
type KafkaDestination struct {
	URL   string `json:"producer_url"`
	Topic string `json:"producer_topic"`
}

// FetchData connects to Kafka, retrieves data, and processes it.
func (k KafkaSource) FetchData(req interfaces.Request) (interface{}, error) {
	logger.Infof("Connecting to Kafka Source: URL=%s, Topic=%s", req.ConsumerURL, req.ConsumerTopic)

	if req.ConsumerURL == "" || req.ConsumerTopic == "" {
		return nil, errors.New("missing Kafka source details")
	}

	// Create Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(req.ConsumerURL, ","),
		Topic:    req.ConsumerTopic,
		GroupID:  "fractal-group", // Example: change as needed
		MinBytes: 10e3,            // 10KB
		MaxBytes: 10e6,            // 10MB
	})
	defer reader.Close()

	// Process messages
	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			return nil, err
		}

		logger.Infof("Message received from Kafka: %s", message.Value)

		// Validation
		validatedData, err := validateKafkaData(message.Value)
		if err != nil {
			logger.Fatalf("Validation failed for message: %s, Error: %s", message.Value, err)
			continue // Skip invalid message
		}

		// Transformation
		transformedData := transformKafkaData(validatedData)
		logger.Infof("Message successfully processed", transformedData)
		return transformedData, nil
	}

}

// SendData connects to Kafka and publishes data to the specified topic.
func (k KafkaDestination) SendData(data interface{}, req interfaces.Request) error {
	logger.Infof("Connecting to Kafka Destination: URL=%s, Topic=%s", req.ProducerURL, req.ProducerTopic)

	if req.ProducerURL == "" || req.ProducerTopic == "" {
		return errors.New("missing Kafka target details")
	}

	// Create Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: strings.Split(req.ProducerURL, ","),
		Topic:   req.ProducerTopic,
	})
	defer writer.Close()

	// Convert data to a string
	var message string
	switch v := data.(type) {
	case string:
		message = v
	case []byte:
		message = string(v) // Convert bytes to a string
	default:
		return fmt.Errorf("unsupported data type: %T", v)
	}

	// Publish message
	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte(message),
		},
	)
	if err != nil {
		return err
	}

	logger.Infof("Message sent to Kafka topic %s: %s", req.ProducerTopic, message)
	return nil
}

// Initialize the Kafka integrations by registering them with the registry.
func init() {
	registry.RegisterSource("Kafka", KafkaSource{})
	registry.RegisterDestination("Kafka", KafkaDestination{})
}

// validateData ensures the input data meets the required criteria.
func validateKafkaData(data []byte) ([]byte, error) {
	logger.Infof("Validating data: %s", data)

	// Example: Check if data is non-empty
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	// Add custom validation logic here
	return data, nil
}

// transformData modifies the input data as per business logic.
func transformKafkaData(data []byte) []byte {
	logger.Infof("Transforming data: %s", data)

	// Example: Convert data to uppercase (modify as needed)
	transformed := []byte(strings.ToUpper(string(data)))
	return transformed
}
