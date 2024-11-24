package integrations

import (
	"errors"
	"strings"
	"sync"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
	"github.com/streadway/amqp"
)

// RabbitMQSource struct represents the configuration for consuming messages from RabbitMQ.
type RabbitMQSource struct {
	URL       string `json:"rabbitmq_input_url"`
	QueueName string `json:"rabbitmq_input_queue_name"`
}

// RabbitMQDestination struct represents the configuration for publishing messages to RabbitMQ.
type RabbitMQDestination struct {
	URL       string `json:"rabbitmq_output_url"`
	QueueName string `json:"rabbitmq_output_queue_name"`
}

// FetchData connects to RabbitMQ, retrieves data, and processes it concurrently.
func (r RabbitMQSource) FetchData(req interfaces.Request) (interface{}, error) {
	logger.Infof("Connecting to RabbitMQ Source: URL=%s, Queue=%s", req.RabbitMQInputURL, req.RabbitMQInputQueueName)

	if req.RabbitMQInputURL == "" || req.RabbitMQInputQueueName == "" {
		return nil, errors.New("missing RabbitMQ source details")
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(req.RabbitMQInputURL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	// Consume messages
	msgs, err := ch.Consume(
		req.RabbitMQInputQueueName, // queue
		"",                         // consumer
		true,                       // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
	)
	if err != nil {
		return nil, err
	}

	// Use a buffered channel for processing messages
	messageChannel := make(chan []byte, 10)
	var wg sync.WaitGroup

	// Start multiple goroutines for concurrent processing
	for i := 0; i < 5; i++ { // Number of workers
		wg.Add(1)
		go func() {
			defer wg.Done()
			for message := range messageChannel {
				processRabbitMQMessage(message)
			}
		}()
	}

	// Read messages from RabbitMQ and send to the channel
	go func() {
		for msg := range msgs {
			messageChannel <- msg.Body
		}
		close(messageChannel)
	}()

	wg.Wait()
	return nil, nil // Return nil as we process messages asynchronously
}

// SendData connects to RabbitMQ and publishes data to the specified queue.
func (r RabbitMQDestination) SendData(data interface{}, req interfaces.Request) error {
	logger.Infof("Connecting to RabbitMQ Destination: URL=%s, Queue=%s", req.RabbitMQOutputURL, req.RabbitMQOutputQueueName)

	if req.RabbitMQOutputURL == "" || req.RabbitMQOutputQueueName == "" {
		return errors.New("missing RabbitMQ target details")
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(req.RabbitMQOutputURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the queue to ensure it exists
	_, err = ch.QueueDeclare(
		req.RabbitMQOutputQueueName, // queue name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)
	if err != nil {
		return err
	}

	// Convert the data to a byte slice
	messageBody, ok := data.([]byte)
	if !ok {
		return errors.New("unsupported data type for RabbitMQ message")
	}

	// Publish the message
	err = ch.Publish(
		"",                          // exchange
		req.RabbitMQOutputQueueName, // routing key
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        messageBody,
		},
	)
	if err != nil {
		return err
	}

	logger.Infof("Message sent to RabbitMQ queue %s: %s", req.RabbitMQOutputQueueName, string(messageBody))
	return nil
}

// processRabbitMQMessage handles individual RabbitMQ messages.
func processRabbitMQMessage(message []byte) {
	logger.Infof("Processing RabbitMQ message: %s", message)

	// Validation
	validatedData, err := validateRabbitMQData(message)
	if err != nil {
		logger.Errorf("Validation failed: %s", err)
		return
	}

	// Transformation
	transformedData := transformRabbitMQData(validatedData)

	logger.Infof("Message processed successfully: %s", transformedData)
}

// validateRabbitMQData ensures the input data meets the required criteria.
func validateRabbitMQData(data []byte) ([]byte, error) {
	logger.Infof("Validating data: %s", data)

	// Example: Check if data is non-empty
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	return data, nil
}

// transformRabbitMQData modifies the input data as per business logic.
func transformRabbitMQData(data []byte) []byte {
	logger.Infof("Transforming data: %s", data)

	// Example: Convert data to uppercase
	return []byte(strings.ToUpper(string(data)))
}

// Initialize the RabbitMQ integrations by registering them with the registry.
func init() {
	registry.RegisterSource("RabbitMQ", RabbitMQSource{})
	registry.RegisterDestination("RabbitMQ", RabbitMQDestination{})
}
