package integrations

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
)

// ReadCSV reads the content of a CSV file and returns it as a byte slice.
func ReadCSV(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []byte
	for _, record := range records {
		data = append(data, []byte(strings.Join(record, ","))...)
		data = append(data, '\n')
	}

	return data, nil
}

func WriteCSV(fileName string, data []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	records := strings.Split(strings.TrimSpace(string(data)), "\n") // Trim trailing newlines
	for _, record := range records {
		fields := strings.Split(record, ",")
		err := writer.Write(fields)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error() // Ensure to check for flush errors
}

// CSVSource struct represents the configuration for consuming messages from CSV.
type CSVSource struct {
	CSVSourceFileName string `json:"csv_source_file_name"`
}

// CSVDestination struct represents the configuration for publishing messages to CSV.
type CSVDestination struct {
	CSVDestinationFileName string `json:"csv_destination_file_name"`
}

// FetchData connects to CSV, retrieves data, and passes it through validation and transformation pipelines.
func (r CSVSource) FetchData(req interfaces.Request) (interface{}, error) {
	logger.Infof("Reading data from CSV Source: %s", req.CSVSourceFileName)

	if req.CSVSourceFileName == "" {
		return nil, errors.New("missing CSV source file name")
	}

	// Read data from CSV
	data, err := ReadCSV(req.CSVSourceFileName)
	if err != nil {
		return nil, err
	}

	// Validate the data
	validatedData, err := validateCSVData(data)
	if err != nil {
		return nil, err
	}

	// Transform the data
	transformedData, err := transformCSVData(validatedData)
	if err != nil {
		return nil, err
	}
	return transformedData, nil

}

// SendData connects to CSV and publishes data to the specified queue.
func (r CSVDestination) SendData(data interface{}, req interfaces.Request) error {
	logger.Infof("Writing data to CSV Destination: %s", req.CSVDestinationFileName)

	if req.CSVDestinationFileName == "" {
		return errors.New("missing CSV destination file name")
	}

	// Write data to CSV
	err := WriteCSV(req.CSVDestinationFileName, data.([]byte))
	if err != nil {
		return err
	}
	return nil
}

// Initialize the CSV integrations by registering them with the registry.
func init() {
	registry.RegisterSource("CSV", CSVSource{})
	registry.RegisterDestination("CSV", CSVDestination{})
}

// validateCSVData ensures the input data meets the required criteria.
func validateCSVData(data []byte) ([]byte, error) {
	logger.Infof("Validating data: %s", data)

	// Example: Check if data is non-empty
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	// Add custom validation logic here
	return data, nil
}

// transformCSVData modifies the input data as per business logic.
func transformCSVData(data []byte) ([]byte, error) {
	logger.Infof("Transforming data: %s", data)

	// Example: Convert data to uppercase (modify as needed)
	transformed := []byte(strings.ToUpper(string(data)))
	return transformed, nil
}
