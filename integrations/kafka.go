package integrations

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

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

// FetchData connects to Kafka, retrieves data, and processes it concurrently.
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

	var wg sync.WaitGroup
	msgChannel := make(chan interface{}, 100) // Buffered channel to collect results

	// Process messages concurrently
	go func() {
		for {
			message, err := reader.ReadMessage(context.Background())
			if err != nil {
				logger.Errorf("Error reading message from Kafka: %v", err)
				continue
			}

			logger.Infof("Message received from Kafka: %s", message.Value)

			// Validation
			validatedData, err := validateKafkaData(message.Value)
			if err != nil {
				logger.Errorf("Validation failed for message: %s, Error: %s", message.Value, err)
				continue // Skip invalid message
			}

			// Transformation
			transformedData := transformKafkaData(validatedData)

			// Send processed data to channel for further handling
			wg.Add(1)
			go func(data interface{}) {
				defer wg.Done()
				msgChannel <- data
			}(transformedData)
		}
	}()

	// Wait for all goroutines to finish processing
	go func() {
		wg.Wait()
		close(msgChannel)
	}()

	// Collect final data from the channel
	var result interface{}
	for data := range msgChannel {
		result = data
	}

	return result, nil
}

// SendData connects to Kafka and publishes data to the specified topic concurrently.
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

	// Convert data to string
	var message string
	switch v := data.(type) {
	case string:
		message = v
	case []byte:
		message = string(v) // Convert bytes to string
	default:
		return fmt.Errorf("unsupported data type: %T", v)
	}

	// Batch send messages concurrently
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Publish message
		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Value: []byte(message),
			},
		)
		if err != nil {
			errCh <- err
		}
	}()

	// Wait for all goroutines to finish and handle errors
	wg.Wait()
	close(errCh)

	// If any error occurred during message sending, return it
	if err := <-errCh; err != nil {
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

// validateKafkaData ensures the input data meets the required criteria.
func validateKafkaData(data []byte) ([]byte, error) {
	logger.Infof("Validating data: %s", data)

	// Example: Check if data is non-empty
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	// Add custom validation logic here
	return data, nil
}

// transformKafkaData modifies the input data as per business logic.
func transformKafkaData(data []byte) []byte {
	logger.Infof("Transforming data: %s", data)

	// Example: Convert data to uppercase (modify as needed)
	transformed := []byte(strings.ToUpper(string(data)))
	return transformed
}
