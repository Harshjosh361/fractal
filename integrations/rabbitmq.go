package integrations

import (
	"errors"

	"github.com/SkySingh04/fractal/logger"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
)

type RabbitMQSource struct {
	URL       string `json:"rabbitmq_input_url"`
	QueueName string `json:"rabbitmq_input_queue_name"`
}

type RabbitMQDestination struct {
	URL       string `json:"rabbitmq_output_url"`
	QueueName string `json:"rabbitmq_output_queue_name"`
}

func (r RabbitMQSource) FetchData(req interfaces.Request) (interface{}, error) {
	if req.RabbitMQInputURL == "" || req.RabbitMQInputQueueName == "" {
		return nil, errors.New("missing RabbitMQ source details")
	}
	logger.Infof("Fetching data from RabbitMQ...")
	// Add RabbitMQ fetch logic here
	return "RabbitMQData", nil
}

func (r RabbitMQDestination) SendData(data interface{}, req interfaces.Request) error {
	if req.RabbitMQOutputURL == "" || req.RabbitMQOutputQueueName == "" {
		return errors.New("missing RabbitMQ target details")
	}
	logger.Infof("Sending data to RabbitMQ...")
	// Add RabbitMQ send logic here
	return nil
}

func init() {
	registry.RegisterSource("RabbitMQ", RabbitMQSource{})
	registry.RegisterDestination("RabbitMQ", RabbitMQDestination{})
}
