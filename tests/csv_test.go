package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/SkySingh04/fractal/integrations"
	"github.com/SkySingh04/fractal/interfaces"
	"github.com/stretchr/testify/assert"
)

const (
	greenTick = "\033[32m✔\033[0m" // Green tick
	redCross  = "\033[31m✘\033[0m" // Red cross
)

func TestCSVIntegration(t *testing.T) {
	// Set up
	inputFileName := "test_input.csv"
	outputFileName := "test_output.csv"

	inputContent := `name,age,city
John,25,New York
Jane,30,San Francisco`

	// Create a temporary input file
	err := ioutil.WriteFile(inputFileName, []byte(inputContent), 0644)
	if err != nil {
		fmt.Printf("%s Error creating test input file: %v\n", redCross, err)
		t.FailNow()
	}
	defer os.Remove(inputFileName) // Cleanup input file after test

	// Clean up output file after test
	defer os.Remove(outputFileName)

	// Prepare the request object for CSV source
	req := interfaces.Request{
		CSVSourceFileName:      inputFileName,
		CSVDestinationFileName: outputFileName,
	}

	// Test FetchData
	csvSource := integrations.CSVSource{}
	data, err := csvSource.FetchData(req)
	if assert.NoError(t, err, "Error fetching data from CSV source") {
		fmt.Printf("%s FetchData passed\n", greenTick)
	} else {
		fmt.Printf("%s FetchData failed\n", redCross)
	}

	// Validate transformed data
	expectedTransformedData := `NAME,AGE,CITY
JOHN,25,NEW YORK
JANE,30,SAN FRANCISCO
`
	if assert.Equal(t, expectedTransformedData, string(data.([]byte)), "Transformed data mismatch") {
		fmt.Printf("%s Data validation passed\n", greenTick)
	} else {
		fmt.Printf("%s Data validation failed\n", redCross)
	}

	// Test SendData
	csvDestination := integrations.CSVDestination{}
	err = csvDestination.SendData(data, req)
	if assert.NoError(t, err, "Error sending data to CSV destination") {
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

	if assert.Equal(t, expectedTransformedData, string(outputData), "Output file content mismatch") {
		fmt.Printf("%s Output file content validation passed\n", greenTick)
	} else {
		fmt.Printf("%s Output file content validation failed\n", redCross)
	}
}
