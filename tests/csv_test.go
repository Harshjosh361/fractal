package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/SkySingh04/fractal/integrations"
	"github.com/SkySingh04/fractal/interfaces"
	"github.com/stretchr/testify/assert"
)

// const (
// 	GreenTick = "\033[32m✔\033[0m" // Green tick
// 	RedCross  = "\033[31m✘\033[0m" // Red cross
// )

func TestCSVIntegration(t *testing.T) {

	const (
		GreenTick = "\033[32m✔\033[0m" // Green tick
		RedCross  = "\033[31m✘\033[0m" // Red cross
	)
	// Set up
	inputFileName := "test_input.csv"
	outputFileName := "test_output.csv"

	inputContent := `name,age,city
John,25,New York
Jane,30,San Francisco`

	// Create a temporary input file
	err := os.WriteFile(inputFileName, []byte(inputContent), 0644)
	if err != nil {
		fmt.Printf("%s Error creating test input file: %v\n", RedCross, err)
		t.FailNow()
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
		fmt.Printf("%s FetchData passed\n", GreenTick)
	} else {
		fmt.Printf("%s FetchData failed\n", RedCross)
	}

	// Validate transformed data
	expectedTransformedData := `NAME,AGE,CITY
JOHN,25,NEW YORK
JANE,30,SAN FRANCISCO
`
	if assert.Equal(t, expectedTransformedData, string(data.([]byte)), "Transformed data mismatch") {
		fmt.Printf("%s Data validation passed\n", GreenTick)
	} else {
		fmt.Printf("%s Data validation failed\n", RedCross)
	}

	csvDestination := integrations.CSVDestination{}
	err = csvDestination.SendData(dataStr, req)
	if assert.NoError(t, err, "Error sending data to CSV destination") {
		fmt.Printf("%s SendData passed\n", GreenTick)
	} else {
		fmt.Printf("%s SendData failed\n", RedCross)
	}

	outputData, err := os.ReadFile(outputFileName)
	if assert.NoError(t, err, "Error reading test output file") {
		fmt.Printf("%s Output file reading passed\n", GreenTick)
	} else {
		fmt.Printf("%s Output file reading failed\n", RedCross)
	}

	if assert.Equal(t, expectedTransformedData, string(outputData), "Output file content mismatch") {
		fmt.Printf("%s Output file content validation passed\n", GreenTick)
	} else {
		fmt.Printf("%s Output file content validation failed\n", RedCross)
	}
}
