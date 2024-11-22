package integrations

import (
	"errors"
	"fmt"
	"log"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDBSource implements the DataSource interface
type DynamoDBSource struct {
	TableName string `json:"table_name"`
	Region    string `json:"region"`
}

// DynamoDBDestination implements the DataDestination interface
type DynamoDBDestination struct {
	TableName string `json:"table_name"`
	Region    string `json:"region"`
}

// FetchData fetches data from DynamoDB
func (d DynamoDBSource) FetchData(req interfaces.Request) (interface{}, error) {
	if err := validateDynamoDBRequest(req, true); err != nil {
		return nil, err
	}
	log.Println("Fetching data from DynamoDB...")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(req.DynamoDBSourceRegion),
	})
	if err != nil {
		return nil, err
	}

	svc := dynamodb.New(sess)
	params := &dynamodb.ScanInput{
		TableName: aws.String(req.DynamoDBSourceTable),
	}
	result, err := svc.Scan(params)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// SendData sends data to DynamoDB
func (d DynamoDBDestination) SendData(data interface{}, req interfaces.Request) error {
	if err := validateDynamoDBRequest(req, false); err != nil {
		return err
	}
	log.Println("Sending data to DynamoDB...")
	// Example logic for sending data, adapt as needed
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(req.DynamoDBTargetRegion),
	})
	if err != nil {
		return err
	}

	svc := dynamodb.New(sess)
	fmt.Println(svc)
	// Implement logic to send data to DynamoDB
	// For instance, using PutItem or BatchWriteItem APIs

	return nil
}

// validateDynamoDBRequest validates the request fields for DynamoDB
func validateDynamoDBRequest(req interfaces.Request, isSource bool) error {
	if isSource {
		if req.DynamoDBSourceTable == "" {
			return errors.New("missing source DynamoDB table")
		}
		if req.DynamoDBSourceRegion == "" {
			return errors.New("missing source DynamoDB region")
		}
	} else {
		if req.DynamoDBTargetTable == "" {
			return errors.New("missing target DynamoDB table")
		}
		if req.DynamoDBTargetRegion == "" {
			return errors.New("missing target DynamoDB region")
		}
	}
	return nil
}

func init() {
	registry.RegisterSource("DynamoDB", DynamoDBSource{})
	registry.RegisterDestination("DynamoDB", DynamoDBDestination{})
}
