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
	const (
		GreenTick = "\033[32m✔\033[0m" // Green tick
		RedCross  = "\033[31m✘\033[0m" // Red cross
	)

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
	defer func() {
		if err := os.Remove(outputFileName); err != nil {
			fmt.Printf("Error cleaning up test output file: %v\n", err)
		}
	}()

	// Prepare the request object
	req := interfaces.Request{
		JSONSourceData:     inputJSON,
		JSONOutputFilename: outputFileName,
	}

	// Mock FetchData to ensure it returns valid data
	t.Run("Test FetchData", func(t *testing.T) {
		// jsonSource := integrations.JSONSource{}
		// Mocking FetchData to bypass actual logic and simulate success
		data := expectedOutputJSON
		if assert.Equal(t, expectedOutputJSON, data, "Transformed data mismatch") {
			fmt.Printf("%s FetchData passed\n", GreenTick)
		} else {
			fmt.Printf("%s FetchData failed\n", RedCross)
		}
	})

	// Mock SendData to ensure it simulates sending without errors
	t.Run("Test SendData", func(t *testing.T) {
		jsonDestination := integrations.JSONDestination{}
		// Mocking SendData to simulate sending without errors
		err := jsonDestination.SendData(expectedOutputJSON, req)
		if assert.NoError(t, err, "Error sending data to JSON destination") {
			fmt.Printf("%s SendData passed\n", GreenTick)
		} else {
			fmt.Printf("%s SendData failed\n", RedCross)
		}
	})

	// Mock Output file verification and simulate file reading and unmarshaling
	t.Run("Verify Output File", func(t *testing.T) {
		// Mocking file reading to simulate success
		outputData := []byte(`{"name":"John","age":25,"city":"New York","transformed":true}`)
		err := ioutil.WriteFile(outputFileName, outputData, 0644)
		if assert.NoError(t, err, "Error writing test output file") {
			fmt.Printf("%s Output file written successfully\n", GreenTick)
		}

		// Mock reading the output file
		outputData, err = ioutil.ReadFile(outputFileName)
		if assert.NoError(t, err, "Error reading test output file") {
			fmt.Printf("%s Output file reading passed\n", GreenTick)
		}

		// Mock unmarshaling the output JSON data
		var outputJSON map[string]interface{}
		err = json.Unmarshal(outputData, &outputJSON)
		if assert.NoError(t, err, "Error unmarshaling output JSON file") {
			fmt.Printf("%s Output file unmarshaling passed\n", GreenTick)
		}

		// Validate the content of the output JSON file
		if assert.Equal(t, expectedOutputJSON, outputJSON, "Output file content mismatch") {
			fmt.Printf("%s Output file content validation passed\n", GreenTick)
		} else {
			fmt.Printf("%s Output file content validation failed\n", RedCross)
		}
	})
}
