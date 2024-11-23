package tests

import (
	"os"
	"testing"

	"github.com/SkySingh04/fractal/integrations"
	"github.com/SkySingh04/fractal/interfaces"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func createTempYAMLFile(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "*.yaml")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write([]byte(content))
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func TestYAMLIntegration(t *testing.T) {
	greenTick := "\033[32m✔\033[0m" // Green tick
	redCross := "\033[31m✘\033[0m"  // Red cross

	logTestStatus := func(description string, err error) {
		if err == nil {
			t.Logf("%s %s", greenTick, description)
		} else {
			t.Logf("%s %s: %v", redCross, description, err)
		}
	}

	// Create temporary source YAML file
	sourceContent := `
name: TestUser
age: 30
skills:
  - Go
  - Kubernetes
`
	sourceFilePath, err := createTempYAMLFile(sourceContent)
	logTestStatus("Create temporary source YAML file", err)
	assert.NoError(t, err)
	defer os.Remove(sourceFilePath)

	// Destination file path
	destinationFilePath := sourceFilePath + "_out.yaml"

	// Initialize YAMLSource and YAMLDestination
	yamlSource := integrations.YAMLSource{}
	yamlDestination := integrations.YAMLDestination{}

	// Define the request
	req := interfaces.Request{
		YAMLSourceFilePath:      sourceFilePath,
		YAMLDestinationFilePath: destinationFilePath,
	}

	// Fetch data from source
	fetchedData, err := yamlSource.FetchData(req)
	logTestStatus("Fetch data from YAML source", err)
	assert.NoError(t, err, "FetchData failed")
	assert.NotNil(t, fetchedData, "Fetched data should not be nil")

	// Write data to destination
	err = yamlDestination.SendData(fetchedData, req)
	logTestStatus("Write data to YAML destination", err)
	assert.NoError(t, err, "SendData failed")

	// Verify written data
	writtenData, err := os.ReadFile(destinationFilePath)
	logTestStatus("Read data from YAML destination file", err)
	assert.NoError(t, err, "Failed to read destination file")
	defer os.Remove(destinationFilePath)

	var result map[string]interface{}
	err = yaml.Unmarshal(writtenData, &result)
	logTestStatus("Unmarshal YAML data from destination file", err)
	assert.NoError(t, err, "Unmarshalling written YAML failed")

	// Validate content
	if assert.Equal(t, "TestUser", result["name"], "Name should match") {
		logTestStatus("Validate 'name' field", nil)
	} else {
		logTestStatus("Validate 'name' field", assert.AnError)
	}

	if assert.Equal(t, 30, result["age"], "Age should match") {
		logTestStatus("Validate 'age' field", nil)
	} else {
		logTestStatus("Validate 'age' field", assert.AnError)
	}

	if assert.Equal(t, []interface{}{"Go", "Kubernetes"}, result["skills"], "Skills should match") {
		logTestStatus("Validate 'skills' field", nil)
	} else {
		logTestStatus("Validate 'skills' field", assert.AnError)
	}

	if assert.Equal(t, true, result["transformed"], "Expected 'transformed' key in output") {
		logTestStatus("Validate 'transformed' key in output", nil)
	} else {
		logTestStatus("Validate 'transformed' key in output", assert.AnError)
	}
}
