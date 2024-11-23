package integrations

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/language"
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

// WriteCSV writes data to a CSV file.
func WriteCSV(fileName string, data []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	records := strings.Split(strings.TrimSpace(string(data)), "\n")
	for _, record := range records {
		fields := strings.Split(record, ",")
		err := writer.Write(fields)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
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
	validatedData, err := validateCSVData(data, req.ValidationRules)
	if err != nil {
		return nil, err
	}

	// Transform the data
	transformedData, err := transformCSVData(validatedData, req.TransformationRules)
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

// validateCSVData ensures the input data meets the required criteria using validation rules.
func validateCSVData(data []byte, validationRules string) ([]byte, error) {
	
	logger.Infof("Validating data: %s", data)

	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	// Initialize lexer and tokenize the validation rules
	lexer := language.NewLexer(validationRules)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize validation rules: %v", err)
	}

	// Parse the tokens into an AST
	parser := language.NewParser(tokens)
	rulesAST, err := parser.ParseRules()
	if err != nil {
		return nil, fmt.Errorf("failed to parse validation rules: %v", err)
	}

	// Apply validation rules to data
	records := strings.Split(strings.TrimSpace(string(data)), "\n")
	for _, record := range records {
		for _, ruleNode := range rulesAST.Children {
			err := applyValidationRule(record, ruleNode)
			if err != nil {
				return nil, err // Return the first validation error encountered
			}
		}
	}

	return data, nil

}

// transformCSVData modifies the input data as per business logic using transformation rules.
func transformCSVData(data []byte, transformationRules string) ([]byte, error) {
	logger.Infof("Transforming data: %s", data)

	// Initialize lexer and tokenize the transformation rules
	lexer := language.NewLexer(transformationRules)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize transformation rules: %v", err)
	}

	// Parse the tokens into an AST
	parser := language.NewParser(tokens)
	rulesAST, err := parser.ParseRules()
	if err != nil {
		return nil, fmt.Errorf("failed to parse transformation rules: %v", err)
	}

	// Apply transformation rules to data
	var transformedRecords []string
	records := strings.Split(strings.TrimSpace(string(data)), "\n")
	for _, record := range records {
		for _, ruleNode := range rulesAST.Children {
			transformedRecord, err := applyTransformationRule(record, ruleNode)
			if err != nil {
				return nil, err
			}
			record = transformedRecord // Apply each rule sequentially
		}
		transformedRecords = append(transformedRecords, record)
	}

	return []byte(strings.Join(transformedRecords, "\n")), nil
}

// applyValidationRule processes a single record against a validation rule AST node.
func applyValidationRule(record string, ruleNode *language.Node) error {
	// Implementation details based on your business rules
	// Validate the record using the information from ruleNode
	// Example: Check if a specific field meets the condition
	return nil // Replace with actual validation logic
}

// applyTransformationRule processes a single record against a transformation rule AST node.
func applyTransformationRule(record string, ruleNode *language.Node) (string, error) {
	// Implementation details based on your business logic
	// Transform the record using the information from ruleNode
	return record, nil // Replace with actual transformation logic
}
