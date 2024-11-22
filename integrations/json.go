package integrations

import (
	"errors"
	"log"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
)

type JSONSource struct {
	Data string `json:"json_source_data"`
}

type JSONDestination struct {
	Data string `json:"json_output_data"`
}

func (j JSONSource) FetchData(req interfaces.Request) (interface{}, error) {
	if req.JSONSourceData == "" {
		return nil, errors.New("missing JSON source data")
	}
	log.Println("Fetching data from JSON source...")
	// Add JSON fetch logic here
	return req.JSONSourceData, nil
}

func (j JSONDestination) SendData(data interface{}, req interfaces.Request) error {
	if req.JSONOutputData == "" {
		return errors.New("missing JSON destination data")
	}
	log.Println("Sending data to JSON destination...")
	// Add JSON send logic here
	return nil
}

func init() {
	registry.RegisterSource("JSON", JSONSource{})
	registry.RegisterDestination("JSON", JSONDestination{})
}
