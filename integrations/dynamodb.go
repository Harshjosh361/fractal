package integrations

import (
	"errors"
	"strings"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

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

	// Create a DynamoDB session for the specified region
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(req.DynamoDBSourceRegion),
		Endpoint: aws.String("http://localhost:8000"), // Specify DynamoDB Local endpoint here
	})
	if err != nil {
		return nil, err
	}

	svc := dynamodb.New(sess)

	// Scan the table
	input := &dynamodb.ScanInput{
		TableName: aws.String(req.DynamoDBSourceTable),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	// Handle empty result
	if len(result.Items) == 0 {
		logger.Warnf("No data retrieved from DynamoDB table: %s", req.DynamoDBSourceTable)
		return nil, errors.New("no data retrieved from DynamoDB")
	}

	// Process and transform items
	for _, item := range result.Items {
		// Validate data
		validatedData, err := validateDynamoDBData(item)
		if err != nil {
			logger.Fatalf("Validation failed for item: %v, Error: %s", item, err)
			continue
		}

		// Transform data
		transformedData := transformDynamoDBData(validatedData)

		logger.Infof("Item successfully processed: %v", transformedData)
		return transformedData, nil
	}

	return nil, errors.New("no valid data processed from DynamoDB")
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

// SendData writes data to the target DynamoDB table in the specified region.
func (d DynamoDBDestination) SendData(data interface{}, req interfaces.Request) error {
	logger.Infof("Connecting to DynamoDB Destination: Table=%s, Region=%s", req.DynamoDBTargetTable, req.DynamoDBTargetRegion)

	// Validate the request
	if err := validateDynamoDBRequest(req, false); err != nil {
		return err
	}

	// Create a DynamoDB session for the specified region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(req.DynamoDBTargetRegion),
	})
	if err != nil {
		return err
	}

	svc := dynamodb.New(sess)

	// Prepare the item
	item, err := prepareDynamoDBItem(data)
	if err != nil {
		return err
	}

	// Put the item into the target table
	input := &dynamodb.PutItemInput{
		TableName: aws.String(req.DynamoDBTargetTable),
		Item:      item,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	logger.Infof("Data successfully written to DynamoDB table %s: %v", req.DynamoDBTargetTable, data)
	return nil
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

// prepareDynamoDBItem converts generic data into a DynamoDB item.
func prepareDynamoDBItem(data interface{}) (map[string]*dynamodb.AttributeValue, error) {
	// Example: Assume data is a map[string]string
	dataMap, ok := data.(map[string]string)
	if !ok {
		return nil, errors.New("unsupported data type for DynamoDB item")
	}

	item := make(map[string]*dynamodb.AttributeValue)
	for k, v := range dataMap {
		item[k] = &dynamodb.AttributeValue{S: aws.String(v)}
	}

	return item, nil
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
