package helper

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

	// Extract headers from the first line
	headers := strings.Split(records[0], ",")
	for i := range headers {
		headers[i] = strings.TrimSpace(headers[i]) // Trim whitespace from headers
	}
	//remove the header
	records = records[1:]
	logger.Infof("Validating %d records", len(records))
	for _, record := range records {
		logger.Infof("Validating record: %s", record)
		logger.Infof("Number of validation rules: %d", len(rulesAST.Children))
		for _, ruleNode := range rulesAST.Children {
			err := applyValidationRule(record, headers, *ruleNode)
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
func applyValidationRule(record string, headers []string, ruleNode language.Node) error {
	// Split the record into fields (assuming CSV format)
	logger.Infof("Inside of applyValidationRule , record: %s , headers: %v , ruleNode: %v", record, headers, ruleNode)
	fields := strings.Split(record, ",")

	// Ensure the number of fields matches the number of headers
	if len(fields) != len(headers) {
		return fmt.Errorf("mismatch between headers and fields: headers=%v, fields=%v", headers, fields)
	}
	logger.Infof("Headers: %v", headers)
	logger.Infof("Fields: %v", fields)

	// Map field names to their values based on headers
	fieldMap := map[string]string{}
	for i, header := range headers {
		fieldMap[header] = strings.TrimSpace(fields[i])
	}

	// Log the constructed field map
	fmt.Printf("Constructed FieldMap: %v", fieldMap)

	// Evaluate the ruleNode recursively
	logger.Infof("Evaluating rule: %s", ruleNode.Value)
	if err := evaluateNode(&ruleNode, fieldMap); err != nil {
		return err
	}

	return nil
}

// Helper function to resolve field names
func resolveField(field string) string {
	// Check if the field is in the format FIELD("field_name")
	if strings.HasPrefix(field, `FIELD("`) && strings.HasSuffix(field, `")`) {
		return strings.TrimSuffix(strings.TrimPrefix(field, `FIELD("`), `")`)
	}
	return field // Return as-is if not wrapped
}

// Recursive function to evaluate nodes
func evaluateNode(node *language.Node, fieldMap map[string]string) error {
	fmt.Printf("Fieldmap in evaluateNode: %v", fieldMap)
	//print the fields in the fieldMap
	for key, value := range fieldMap {
		fmt.Printf("Key: %s, Value: %s", key, value)
	}
	switch node.Type {
	case language.TokenField:
		return nil // This case is handled within expressions

	case "EXPRESSION":
		fieldNode := node.Children[0]
		conditionNode := node.Children[1]
		valueNode := node.Children[2]

		resolvedField := resolveField(fieldNode.Value) // Resolve FIELD("...") to actual field name
		logger.Infof("Evaluating expression: %s %s %s", resolvedField, conditionNode.Value, valueNode.Value)

		// Display the FieldMap with fmt.Printf for debugging

		fieldMap := map[string]string{
			"name": "Alice",
			"age":  "30",
			"city": "New York",
		}

		fmt.Printf("FieldMap: %v\n", fieldMap)
		// Check if the field exists in FieldMap
		fieldValue, exists := fieldMap[resolvedField]
		if !exists {
			return fmt.Errorf("field %s not found", resolvedField)
		}

		// Remove quotes from fieldValue
		fieldValue = strings.Trim(fieldValue, "'\"")

		switch conditionNode.Value {
		case "TYPE":
			return evaluateTypeCondition(fieldValue, valueNode.Value)
		case "RANGE":
			// fieldValue = strings.ReplaceAll(fieldValue, "'", "")
			logger.Infof("Field Value: %s", fieldValue)
			// logger.Infof("Value Node Children: %v", valueNode.Children)
			// logger.Infof("Value Node Value: %s", valueNode.Value)
			return evaluateRangeCondition(fieldValue, valueNode.Value)
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
func evaluateRangeCondition(value string, rangeNodes string) error {
	//here rangeNodes is a string of min and max values separated by comma like (30 , 50)
	//split the rangeNodes string into min and max values
	rangeValues := strings.Split(rangeNodes, ",")
	//remove the ( from the min value
	rangeValues[0] = strings.ReplaceAll(rangeValues[0], "(", "")
	//remove the ) from the max value
	rangeValues[1] = strings.ReplaceAll(rangeValues[1], ")", "")
	logger.Infof("Range Values: %v", rangeValues)
	if len(rangeValues) != 2 {
		return errors.New("range condition should have two values")
	}
	minValue, err1 := strconv.ParseInt(rangeValues[0], 10, 64)
	maxValue, err2 := strconv.ParseInt(rangeValues[1], 10, 64)
	if err1 != nil || err2 != nil {
		return errors.New("range values should be numbers")
	}
	fieldValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return errors.New("field value should be a number")
	}
	if fieldValue < minValue || fieldValue > maxValue {
		return fmt.Errorf("value '%s' not in range", value)
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
