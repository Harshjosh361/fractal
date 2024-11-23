package integrations

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
	"github.com/segmentio/kafka-go"
)

// KafkaSource implements the DataSource interface
type KafkaSource struct {
	ConsumerURL   string `json:"consumer_url"`
	ConsumerTopic string `json:"consumer_topic"`
}

// KafkaDestination implements the DataDestination interface
type KafkaDestination struct {
	ProducerURL   string `json:"producer_url"`
	ProducerTopic string `json:"producer_topic"`
}

// FetchData fetches data from Kafka
func (k KafkaSource) FetchData(req interfaces.Request) (interface{}, error) {
	if err := validateKafkaRequest(req, true); err != nil {
		return nil, err
	}
	log.Println("Fetching data from Kafka...")
	log.Printf("Connecting to Kafka Consumer at %s for topic %s\n", req.ConsumerURL, req.ConsumerTopic)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{req.ConsumerURL},
		GroupID:  "default-group", // Replace with dynamic value if required
		Topic:    req.ConsumerTopic,
		MinBytes: 10e3,  // 10KB (adjust if needed)
		MaxBytes: 500e6, // Increased to 500MB
	})

	defer reader.Close()

	log.Println("Waiting for messages from Kafka...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Timeout for fetching
	defer cancel()

	message, err := reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Message received: Key=%s, Value=%s\n", string(message.Key), string(message.Value))
	return string(message.Value), nil
}

// SendData sends data to Kafka
func (k KafkaDestination) SendData(data interface{}, req interfaces.Request) error {
	if err := validateKafkaRequest(req, false); err != nil {
		return err
	}

	log.Printf("Connecting to Kafka Producer at %s for topic %s\n", req.ProducerURL, req.ProducerTopic)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:   []string{req.ProducerURL},
		Topic:     req.ProducerTopic,
		Balancer:  &kafka.LeastBytes{},
		BatchSize: 10485760, // Increased to 10MB
	})

	defer writer.Close()

	message := kafka.Message{
		Key:   []byte("Key"),
		Value: []byte(data.(string)), // Ensure the data can be converted to string
	}

	log.Println("Sending message to Kafka...")
	err := writer.WriteMessages(context.Background(), message)
	if err != nil {
		return err
	}

	log.Println("Message successfully sent to Kafka.")
	return nil
}

// validateKafkaRequest validates the request fields for Kafka
func validateKafkaRequest(req interfaces.Request, isConsumer bool) error {
	if isConsumer {
		if req.ConsumerURL == "" {
			return errors.New("missing consumer URL for Kafka")
		}
		if req.ConsumerTopic == "" {
			return errors.New("missing consumer topic for Kafka")
		}
	} else {
		if req.ProducerURL == "" {
			return errors.New("missing producer URL for Kafka")
		}
		if req.ProducerTopic == "" {
			return errors.New("missing producer topic for Kafka")
		}
	}
	return nil
}

func init() {
	registry.RegisterSource("Kafka", KafkaSource{})
	registry.RegisterDestination("Kafka", KafkaDestination{})
}
