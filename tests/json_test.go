package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/SkySingh04/fractal/integrations"
	"github.com/SkySingh04/fractal/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestJSONIntegration(t *testing.T) {
	// Setup
	inputJSON := `{"name": "John", "age": 25, "city": "New York"}`
	expectedOutputJSON := map[string]interface{}{
		"name":        "John",
		"age":         float64(25),
		"city":        "New York",
		"transformed": true,
	}
	outputFileName := "test_output.json"

	// Clean up output file after test
	defer os.Remove(outputFileName)

	// Prepare the request object
	req := interfaces.Request{
		JSONSourceData:     inputJSON,
		JSONOutputFilename: outputFileName,
	}

	// Test FetchData
	jsonSource := integrations.JSONSource{}
	data, err := jsonSource.FetchData(req)
	if assert.NoError(t, err, "Error fetching data from JSON source") {
		fmt.Printf("%s FetchData passed\n", greenTick)
	} else {
		fmt.Printf("%s FetchData failed\n", redCross)
	}

	// Validate fetched and transformed data
	if assert.Equal(t, expectedOutputJSON, data, "Transformed data mismatch") {
		fmt.Printf("%s Data validation passed\n", greenTick)
	} else {
		fmt.Printf("%s Data validation failed\n", redCross)
	}

	// Test SendData
	jsonDestination := integrations.JSONDestination{}
	err = jsonDestination.SendData(data, req)
	if assert.NoError(t, err, "Error sending data to JSON destination") {
		fmt.Printf("%s SendData passed\n", greenTick)
	} else {
		fmt.Printf("%s SendData failed\n", redCross)
	}

	// Verify the output file contents
	outputData, err := ioutil.ReadFile(outputFileName)
	if assert.NoError(t, err, "Error reading test output file") {
		fmt.Printf("%s Output file reading passed\n", greenTick)
	} else {
		fmt.Printf("%s Output file reading failed\n", redCross)
	}

	// Validate the content of the output JSON file
	var outputJSON map[string]interface{}
	err = json.Unmarshal(outputData, &outputJSON)
	if assert.NoError(t, err, "Error unmarshaling output JSON file") {
		fmt.Printf("%s Output file unmarshaling passed\n", greenTick)
	} else {
		fmt.Printf("%s Output file unmarshaling failed\n", redCross)
	}

	if assert.Equal(t, expectedOutputJSON, outputJSON, "Output file content mismatch") {
		fmt.Printf("%s Output file content validation passed\n", greenTick)
	} else {
		fmt.Printf("%s Output file content validation failed\n", redCross)
	}
}
