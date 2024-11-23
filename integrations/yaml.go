package integrations

import (
	"errors"
	"io/ioutil"

	"github.com/SkySingh04/fractal/logger"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
	"gopkg.in/yaml.v3"
)

// YAMLSource implements the DataSource interface
type YAMLSource struct {
	FilePath string `json:"file_path"`
}

// YAMLDestination implements the DataDestination interface
type YAMLDestination struct {
	FilePath string `json:"file_path"`
}

// FetchData reads data from a YAML file
func (y YAMLSource) FetchData(req interfaces.Request) (interface{}, error) {
	if err := validateYAMLRequest(req, true); err != nil {
		return nil, err
	}
	logger.Infof("Fetching data from YAML...")
	data, err := ioutil.ReadFile(req.YAMLSourceFilePath)
	if err != nil {
		return nil, err
	}
	var result interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SendData writes data to a YAML file
func (y YAMLDestination) SendData(data interface{}, req interfaces.Request) error {
	if err := validateYAMLRequest(req, false); err != nil {
		return err
	}
	logger.Infof("Sending data to YAML...")
	outputData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(req.YAMLDestinationFilePath, outputData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// validateYAMLRequest validates the request fields for YAML
func validateYAMLRequest(req interfaces.Request, isSource bool) error {
	if isSource && req.YAMLSourceFilePath == "" {
		return errors.New("missing YAML source file path")
	}
	if !isSource && req.YAMLDestinationFilePath == "" {
		return errors.New("missing YAML destination file path")
	}
	return nil
}

func init() {
	registry.RegisterSource("YAML", YAMLSource{})
	registry.RegisterDestination("YAML", YAMLDestination{})
}
