package integrations

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	logger.Infof("request: %v", req)

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

	// logger.Infof("Validating data: %s", data)

	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	logger.Infof("validationRules: %s", validationRules)

	// Initialize lexer and tokenize the validation rules
	lexer := language.NewLexer(validationRules)
	tokens, err := lexer.Tokenize(validationRules)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize validation rules: %v", err)
	}
	logger.Infof("Tokens: %v", tokens)

	// Parse the tokens into an AST
	parser := language.NewParser()
	// var tokenValues []string
	// for _, token := range tokens {
	// 	tokenValues = append(tokenValues, token.Value)
	// }
	rulesAST, err := parser.ParseRules(tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to parse validation rules: %v", err)
	}

	// Apply validation rules to data
	records := strings.Split(strings.TrimSpace(string(data)), "\n")
	logger.Infof("Validating %d records", len(records))
	for _, record := range records {
		logger.Infof("Validating record: %s", record)
		logger.Infof("Number of validation rules: %d", len(rulesAST.Children))
		for _, ruleNode := range rulesAST.Children {
			logger.Infof("Applying rule: %s", ruleNode.Value)
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
	// logger.Infof("Transforming data: %s", data)

	// Initialize lexer and tokenize the transformation rules
	lexer := language.NewLexer(transformationRules)
	tokens, err := lexer.Tokenize(transformationRules)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize transformation rules: %v", err)
	}

	// Parse the tokens into an AST
	parser := language.NewParser()
	// var tokenValues []string
	// for _, token := range tokens {
	// 	tokenValues = append(tokenValues, token.Value)
	// }
	rulesAST, err := parser.ParseRules(tokens)
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
	// Split the record into fields (assuming CSV format)
	fields := strings.Split(record, ",")
	fieldMap := map[string]string{}

	// Map field names to their values based on index (assuming header order)
	// READ the first line of the CSV file to get the headers
	headers := strings.Split(fields[0], ",")
	for i, header := range headers {
		if i < len(fields) {
			fieldMap[header] = strings.TrimSpace(fields[i])
		}
	}

	// Evaluate the ruleNode recursively
	logger.Infof("Evaluating rule: %s", ruleNode.Value)
	if err := evaluateNode(ruleNode, fieldMap); err != nil {
		return err
	}

	return nil
}

// Recursive function to evaluate nodes
func evaluateNode(node *language.Node, fieldMap map[string]string) error {
	switch node.Type {
	case language.TokenField:
		return nil // This case is handled within expressions

	case "EXPRESSION":
		fieldNode := node.Children[0]
		conditionNode := node.Children[1]
		valueNode := node.Children[2]

		fieldValue, exists := fieldMap[fieldNode.Value]
		if !exists {
			return fmt.Errorf("field %s not found", fieldNode.Value)
		}

		switch conditionNode.Value {
		case "TYPE":
			return evaluateTypeCondition(fieldValue, valueNode.Value)
		case "RANGE":
			return evaluateRangeCondition(fieldValue, valueNode.Children)
		case "MATCHES":
			return evaluateRegexCondition(fieldValue, valueNode.Value)
		case "IN":
			return evaluateInCondition(fieldValue, valueNode.Children)
		case "REQUIRED":
			return evaluateRequiredCondition(fieldValue)
		default:
			return fmt.Errorf("unsupported condition: %s", conditionNode.Value)
		}

	case language.TokenLogical:
		return evaluateLogicalCondition(node, fieldMap)
	}

	return fmt.Errorf("unknown node type: %s", node.Type)
}
func evaluateTypeCondition(value, expectedType string) error {
	switch expectedType {
	case "STRING":
		return nil // All values are strings by default
	case "INT":
		if _, err := strconv.Atoi(value); err != nil {
			return fmt.Errorf("value '%s' is not an integer", value)
		}
	case "FLOAT":
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("value '%s' is not a float", value)
		}
	case "BOOL":
		if _, err := strconv.ParseBool(value); err != nil {
			return fmt.Errorf("value '%s' is not a boolean", value)
		}
	case "DATE":
		if _, err := time.Parse("2006-01-02", value); err != nil {
			return fmt.Errorf("value '%s' is not a valid date", value)
		}
	default:
		return fmt.Errorf("unknown type: %s", expectedType)
	}
	return nil
}
func evaluateRangeCondition(value string, rangeNodes []*language.Node) error {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("value '%s' is not numeric", value)
	}

	min, err := strconv.ParseFloat(rangeNodes[0].Value, 64)
	if err != nil {
		return fmt.Errorf("invalid minimum range value: %v", err)
	}
	max, err := strconv.ParseFloat(rangeNodes[1].Value, 64)
	if err != nil {
		return fmt.Errorf("invalid maximum range value: %v", err)
	}

	if val < min || val > max {
		return fmt.Errorf("value '%s' out of range (%f, %f)", value, min, max)
	}
	return nil
}
func evaluateRegexCondition(value, pattern string) error {
	matched, err := regexp.MatchString(pattern, value)
	if err != nil || !matched {
		return fmt.Errorf("value '%s' does not match pattern", value)
	}
	return nil
}
func evaluateInCondition(value string, allowedValues []*language.Node) error {
	for _, valNode := range allowedValues {
		if value == valNode.Value {
			return nil
		}
	}
	return fmt.Errorf("value '%s' not in allowed list", value)
}
func evaluateRequiredCondition(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("field is required and cannot be empty")
	}
	return nil
}

func evaluateLogicalCondition(node *language.Node, fieldMap map[string]string) error {
	leftErr := evaluateNode(node.Children[0], fieldMap)
	rightErr := evaluateNode(node.Children[1], fieldMap)

	switch node.Value {
	case "AND":
		if leftErr != nil || rightErr != nil {
			return errors.New("AND condition failed")
		}
	case "OR":
		if leftErr != nil && rightErr != nil {
			return errors.New("OR condition failed")
		}
	case "NOT":
		if leftErr == nil {
			return errors.New("NOT condition failed")
		}
	default:
		return fmt.Errorf("unknown logical operator: %s", node.Value)
	}
	return nil
}

// Evaluate individual condition (fieldValue compared to ruleValue)
func evaluateCondition(fieldValue, operator, ruleValue string) error {
	switch operator {
	case "==":
		if fieldValue != ruleValue {
			return fmt.Errorf("expected %s, got %s", ruleValue, fieldValue)
		}
	case "!=":
		if fieldValue == ruleValue {
			return fmt.Errorf("field value should not be %s", ruleValue)
		}
	case ">":
		fieldNum, err1 := strconv.ParseFloat(fieldValue, 64)
		ruleNum, err2 := strconv.ParseFloat(ruleValue, 64)
		if err1 != nil || err2 != nil || fieldNum <= ruleNum {
			return fmt.Errorf("expected greater than %s, got %s", ruleValue, fieldValue)
		}
	case "<":
		fieldNum, err1 := strconv.ParseFloat(fieldValue, 64)
		ruleNum, err2 := strconv.ParseFloat(ruleValue, 64)
		if err1 != nil || err2 != nil || fieldNum >= ruleNum {
			return fmt.Errorf("expected less than %s, got %s", ruleValue, fieldValue)
		}
	case ">=":
		fieldNum, err1 := strconv.ParseFloat(fieldValue, 64)
		ruleNum, err2 := strconv.ParseFloat(ruleValue, 64)
		if err1 != nil || err2 != nil || fieldNum < ruleNum {
			return fmt.Errorf("expected greater than or equal to %s, got %s", ruleValue, fieldValue)
		}
	case "<=":
		fieldNum, err1 := strconv.ParseFloat(fieldValue, 64)
		ruleNum, err2 := strconv.ParseFloat(ruleValue, 64)
		if err1 != nil || err2 != nil || fieldNum > ruleNum {
			return fmt.Errorf("expected less than or equal to %s, got %s", ruleValue, fieldValue)
		}

	default:
		return fmt.Errorf("unsupported operator %s", operator)
	}
	return nil
}

// applyTransformationRule processes a single record against a transformation rule AST node.
func applyTransformationRule(record string, ruleNode *language.Node) (string, error) {
	// Implementation details based on your business logic
	// Transform the record using the information from ruleNode
	return record, nil // Replace with actual transformation logic
}
