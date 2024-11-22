package integrations

import (
	"errors"
	"log"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
)

// KafkaSource implements the DataSource interface
type KafkaSource struct{}

// FetchData fetches data from Kafka
func (k KafkaSource) FetchData(req interfaces.Request) (interface{}, error) {
	if err := validateKafkaRequest(req, true); err != nil {
		return nil, err
	}
	log.Println("Fetching data from Kafka...")
	// Your Kafka fetch logic here
	return "KafkaData", nil
}

// KafkaDestination implements the DataDestination interface
type KafkaDestination struct{}

// SendData sends data to Kafka
func (k KafkaDestination) SendData(data interface{}, req interfaces.Request) error {
	if err := validateKafkaRequest(req, false); err != nil {
		return err
	}
	log.Println("Sending data to Kafka...")
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
