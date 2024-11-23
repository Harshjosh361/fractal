package integrations

import (
	"errors"

	"github.com/SkySingh04/fractal/logger"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
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
	logger.Infof("Fetching data from Kafka...")
	// Your Kafka fetch logic here
	return "KafkaData", nil
}

// SendData sends data to Kafka
func (k KafkaDestination) SendData(data interface{}, req interfaces.Request) error {
	if err := validateKafkaRequest(req, false); err != nil {
		return err
	}
	logger.Infof("Sending data to Kafka...")
	// Your Kafka send logic here
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
