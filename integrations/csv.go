package integrations

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

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

// FetchData connects to CSV, retrieves data, and processes it concurrently.
func (r CSVSource) FetchData(req interfaces.Request) (interface{}, error) {
	logger.Infof("Reading data from CSV Source: %s", req.CSVSourceFileName)

	if req.CSVSourceFileName == "" {
		return nil, errors.New("missing CSV source file name")
	}

	// Create channels for processing pipeline
	dataChan := make(chan string, bufferSize)
	validChan := make(chan string, bufferSize)
	transformedChan := make(chan string, bufferSize)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup

	// Start concurrent CSV reading
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := readCSVConcurrently(req.CSVSourceFileName, dataChan, errChan); err != nil {
			errChan <- err
		}
		close(dataChan)
	}()

	// Start concurrent validation
	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range dataChan {
			if validData, err := validateCSVData(data); err != nil {
				errChan <- err
			} else {
				validChan <- validData
			}
		}
		close(validChan)
	}()

	// Start concurrent transformation
	wg.Add(1)
	go func() {
		defer wg.Done()
		for validData := range validChan {
			transformedChan <- transformCSVData(validData)
		}
		close(transformedChan)
	}()

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect results or errors
	var results []string
	for transformedData := range transformedChan {
		results = append(results, transformedData)
	}

	// Check for errors
	if err, ok := <-errChan; ok {
		return nil, err
	}

	return strings.Join(results, "\n"), nil
}

// SendData writes data to a CSV file concurrently.
func (r CSVDestination) SendData(data interface{}, req interfaces.Request) error {
	logger.Infof("Writing data to CSV Destination: %s", req.CSVDestinationFileName)

	if req.CSVDestinationFileName == "" {
		return errors.New("missing CSV destination file name")
	}

	// Convert data to a slice of strings for writing
	lines, ok := data.(string)
	if !ok {
		return errors.New("invalid data format for CSV destination")
	}
	records := strings.Split(lines, "\n")

	// Write concurrently
	errChan := make(chan error, 1)
	go func() {
		errChan <- writeCSVConcurrently(req.CSVDestinationFileName, records)
	}()

	// Check for errors
	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

// readCSVConcurrently reads the content of a CSV file and sends records to a channel.
func readCSVConcurrently(fileName string, out chan<- string, errChan chan<- error) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, os.ErrClosed) || errors.Is(err, io.EOF) {
				break
			}
			errChan <- err
			return err
		}
		out <- strings.Join(record, ",")
	}
	return nil
}

// writeCSVConcurrently writes data records to a CSV file concurrently.
func writeCSVConcurrently(fileName string, records []string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for _, record := range records {
		if err := writer.Write(strings.Split(record, ",")); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

// validateCSVData ensures the input data meets the required criteria.
func validateCSVData(data string) (string, error) {
	logger.Infof("Validating data: %s", data)

	// Example: Check if data is non-empty
	if strings.TrimSpace(data) == "" {
		return "", errors.New("data is empty")
	}

	// Add custom validation logic here
	return data, nil
}

// transformCSVData modifies the input data as per business logic.
func transformCSVData(data string) string {
	logger.Infof("Transforming data: %s", data)

	// Example: Convert data to uppercase (modify as needed)
	return strings.ToUpper(data)
}

// Initialize the CSV integrations by registering them with the registry.
func init() {
	registry.RegisterSource("CSV", CSVSource{})
	registry.RegisterDestination("CSV", CSVDestination{})
}
