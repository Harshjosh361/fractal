package integrations

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// MockDynamoDB is a mock struct for simulating DynamoDB operations.
type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (m *MockDynamoDB) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	// Mocking data returned by Scan based on table name
	if *input.TableName == "input" {
		return &dynamodb.ScanOutput{
			Items: []map[string]*dynamodb.AttributeValue{
				{
					"KeyAttribute": {S: aws.String("sampleKey1")},
					"Data":         {S: aws.String("sampleData1")},
				},
				{
					"KeyAttribute": {S: aws.String("sampleKey2")},
					"Data":         {S: aws.String("sampleData2")},
				},
			},
		}, nil
	}
	return nil, errors.New("table not found")
}

func (m *MockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	// Simulate a successful PutItem operation
	return &dynamodb.PutItemOutput{}, nil
}

// DynamoDBSource represents the configuration for reading data from DynamoDB.
type DynamoDBSource struct {
	TableName string `json:"table_name"`
	Region    string `json:"region"`
}

// DynamoDBDestination represents the configuration for writing data to DynamoDB.
type DynamoDBDestination struct {
	TableName string `json:"table_name"`
	Region    string `json:"region"`
}

// FetchData retrieves data from the source DynamoDB table in the specified region.
func (d DynamoDBSource) FetchData(req interfaces.Request) (interface{}, error) {
	logger.Infof("Connecting to DynamoDB Source: Table=%s, Region=%s", req.DynamoDBSourceTable, req.DynamoDBSourceRegion)

	// Validate the request
	if err := validateDynamoDBRequest(req, true); err != nil {
		return nil, err
	}

	// Mock DynamoDB client
	mockDynamoDB := &MockDynamoDB{}

	// Scan the table
	input := &dynamodb.ScanInput{
		TableName: aws.String(req.DynamoDBSourceTable),
	}

	result, err := mockDynamoDB.Scan(input)
	if err != nil {
		return nil, err
	}

	// Handle empty result
	if len(result.Items) == 0 {
		logger.Logf("No data retrieved from DynamoDB table: %s", req.DynamoDBSourceTable)
		return nil, errors.New("no data retrieved from DynamoDB")
	}

	// Create channels for concurrency
	dataChannel := make(chan map[string]interface{}, len(result.Items))
	errorChannel := make(chan error, len(result.Items))
	var wg sync.WaitGroup

	// Process and transform items concurrently using goroutines
	for _, item := range result.Items {
		wg.Add(1)
		go func(item map[string]*dynamodb.AttributeValue) {
			defer wg.Done()

			// Validate data
			validatedData, err := validateDynamoDBData(item)
			if err != nil {
				errorChannel <- fmt.Errorf("validation failed for item: %v, Error: %s", item, err)
				return
			}

			// Transform data
			transformedData := transformDynamoDBData(validatedData)

			// Convert transformed data (map[string]*dynamodb.AttributeValue) to map[string]interface{}
			interfaceData := make(map[string]interface{})
			for key, value := range transformedData {
				if value.S != nil {
					interfaceData[key] = *value.S
				} else if value.N != nil {
					interfaceData[key] = *value.N
				} else if value.BOOL != nil {
					interfaceData[key] = *value.BOOL
				}
			}

			// Send processed data to the channel
			dataChannel <- interfaceData
		}(item)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close channels after processing
	close(dataChannel)
	close(errorChannel)

	// Check for errors
	if len(errorChannel) > 0 {
		return nil, <-errorChannel
	}

	// Collect and return the processed data
	var processedData []map[string]interface{}
	for data := range dataChannel {
		processedData = append(processedData, data)
	}

	if len(processedData) == 0 {
		return nil, errors.New("no valid data processed from DynamoDB")
	}

	return processedData, nil
}

// SendData writes data to the target DynamoDB table in the specified region.
func (d DynamoDBDestination) SendData(data interface{}, req interfaces.Request) error {
	logger.Infof("Connecting to DynamoDB Destination: Table=%s, Region=%s", req.DynamoDBTargetTable, req.DynamoDBTargetRegion)

	// Validate the request
	if err := validateDynamoDBRequest(req, false); err != nil {
		return err
	}

	// Ensure the data is of the correct type (map[string]interface{})
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		// Attempt to convert the data to map[string]interface{}
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data for conversion: %v", err)
		}

		err = json.Unmarshal(dataBytes, &dataMap)
		if err != nil {
			return fmt.Errorf("failed to unmarshal data for conversion: %v", err)
		}
	}

	// Mock DynamoDB client
	mockDynamoDB := &MockDynamoDB{}

	// Prepare the item
	item, err := prepareDynamoDBItem(dataMap)
	if err != nil {
		return err
	}

	// Put the item into the target table
	input := &dynamodb.PutItemInput{
		TableName: aws.String(req.DynamoDBTargetTable),
		Item:      item,
	}

	_, err = mockDynamoDB.PutItem(input)
	if err != nil {
		return err
	}

	logger.Infof("Data successfully written to DynamoDB table %s: %v", req.DynamoDBTargetTable, data)
	return nil
}

// prepareDynamoDBItem converts a map[string]interface{} to a map[string]*dynamodb.AttributeValue
func prepareDynamoDBItem(data map[string]interface{}) (map[string]*dynamodb.AttributeValue, error) {
	// Convert the map to a DynamoDB-compatible item
	item := make(map[string]*dynamodb.AttributeValue)
	for k, v := range data {
		switch v := v.(type) {
		case string:
			item[k] = &dynamodb.AttributeValue{S: aws.String(v)}
		case int, int64:
			item[k] = &dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%v", v))}
		case bool:
			item[k] = &dynamodb.AttributeValue{BOOL: aws.Bool(v)}
		default:
			return nil, fmt.Errorf("unsupported attribute type for key '%s'", k)
		}
	}

	return item, nil
}

// validateDynamoDBData ensures the input DynamoDB data meets required criteria.
func validateDynamoDBData(data map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
	logger.Infof("Validating DynamoDB data: %v", data)

	// Example: Ensure a specific attribute exists and is not empty
	if val, ok := data["KeyAttribute"]; !ok || val.S == nil || *val.S == "" {
		return nil, errors.New("missing or empty KeyAttribute")
	}

	return data, nil
}

// transformDynamoDBData modifies the input DynamoDB data as per business logic.
func transformDynamoDBData(data map[string]*dynamodb.AttributeValue) map[string]*dynamodb.AttributeValue {
	logger.Infof("Transforming DynamoDB data: %v", data)

	// Example: Convert a string attribute to uppercase
	if val, ok := data["KeyAttribute"]; ok && val.S != nil {
		val.S = aws.String(strings.ToUpper(*val.S))
	}

	return data
}

// validateDynamoDBRequest validates the request fields for DynamoDB operations.
func validateDynamoDBRequest(req interfaces.Request, isSource bool) error {
	if isSource {
		if req.DynamoDBSourceTable == "" || req.DynamoDBSourceRegion == "" {
			return errors.New("missing source DynamoDB table or region")
		}
	} else {
		if req.DynamoDBTargetTable == "" || req.DynamoDBTargetRegion == "" {
			return errors.New("missing target DynamoDB table or region")
		}
	}
	return nil
}

// Register DynamoDB source and destination
func init() {
	registry.RegisterSource("DynamoDB", DynamoDBSource{})
	registry.RegisterDestination("DynamoDB", DynamoDBDestination{})
}
