package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/SkySingh04/fractal/integrations"
	"github.com/SkySingh04/fractal/interfaces"
	"github.com/stretchr/testify/assert"
)


func TestCSVIntegration(t *testing.T) {
	const (
		greenTick = "\033[32m✔\033[0m" // Green tick
		redCross  = "\033[31m✘\033[0m" // Red cross
	)
	
	inputFileName := "test_input.csv"
	outputFileName := "test_output.csv"

	inputContent := `name,age,city
John,25,New York
Jane,30,San Francisco`

	// Create a temporary input file
	err := os.WriteFile(inputFileName, []byte(inputContent), 0644)
	if err != nil {
		t.Fatalf("%s Error creating test input file: %v\n", redCross, err)
	}
	defer os.Remove(inputFileName)
	defer os.Remove(outputFileName)

	req := interfaces.Request{
		CSVSourceFileName:      inputFileName,
		CSVDestinationFileName: outputFileName,
	}

	csvSource := integrations.CSVSource{}
	data, err := csvSource.FetchData(req)
	if assert.NoError(t, err, "Error fetching data from CSV source") {
		t.Logf("%s FetchData passed", greenTick)
	} else {
		t.Fatalf("%s FetchData failed", redCross)
	}

	dataStr, ok := data.(string)
	if !ok {
		t.Fatalf("%s Data type mismatch: expected string", redCross)
	}

	expectedTransformedData := "NAME,AGE,CITY\nJOHN,25,NEW YORK\nJANE,30,SAN FRANCISCO"
	dataStr = strings.TrimSpace(dataStr)
	expectedTransformedData = strings.TrimSpace(expectedTransformedData)

	if assert.Equal(t, expectedTransformedData, dataStr, "Transformed data mismatch") {
		t.Logf("%s Data validation passed", greenTick)
	} else {
		t.Fatalf("%s Data validation failed", redCross)
	}

	csvDestination := integrations.CSVDestination{}
	err = csvDestination.SendData(dataStr, req)
	if assert.NoError(t, err, "Error sending data to CSV destination") {
		t.Logf("%s SendData passed", greenTick)
	} else {
		t.Fatalf("%s SendData failed", redCross)
	}

	outputData, err := os.ReadFile(outputFileName)
	if assert.NoError(t, err, "Error reading test output file") {
		t.Logf("%s Output file reading passed", greenTick)
	} else {
		t.Fatalf("%s Output file reading failed", redCross)
	}

	outputDataStr := strings.TrimSpace(string(outputData))
	if assert.Equal(t, expectedTransformedData, outputDataStr, "Output file content mismatch") {
		t.Logf("%s Output file content validation passed", greenTick)
	} else {
		t.Fatalf("%s Output file content validation failed", redCross)
	}
}