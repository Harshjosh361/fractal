package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/SkySingh04/fractal/integrations"
	"github.com/SkySingh04/fractal/interfaces"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func createTempYAMLFile(content string) (string, error) {
	tmpFile, err := ioutil.TempFile("", "*.yaml")
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
	logTestStatus := func(description string, err error) {
		if err == nil {
			fmt.Printf("✅ %s\n", description)
		} else {
			fmt.Printf("❌ %s: %v\n", description, err)
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
	writtenData, err := ioutil.ReadFile(destinationFilePath)
	logTestStatus("Read data from YAML destination file", err)
	assert.NoError(t, err, "Failed to read destination file")
	defer os.Remove(destinationFilePath)

	var result map[string]interface{}
	err = yaml.Unmarshal(writtenData, &result)
	logTestStatus("Unmarshal YAML data from destination file", err)
	assert.NoError(t, err, "Unmarshalling written YAML failed")

	// Validate content
	assert.Equal(t, "TestUser", result["name"], "Name should match")
	logTestStatus("Validate 'name' field", nil)

	assert.Equal(t, 30, result["age"], "Age should match")
	logTestStatus("Validate 'age' field", nil)

	assert.Equal(t, []interface{}{"Go", "Kubernetes"}, result["skills"], "Skills should match")
	logTestStatus("Validate 'skills' field", nil)

	assert.Equal(t, true, result["transformed"], "Expected 'transformed' key in output")
	logTestStatus("Validate 'transformed' key in output", nil)
}
