package integrations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	// "firebase.google.com/go/db"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
)

// FirebaseSource struct represents the configuration for fetching data from Firebase.
type FirebaseSource struct {
	CredentialFileAddr string `json:"firebase_credential_file"` // Path to service account JSON file
	Collection         string `json:"firebase_collection"`      // Collection name in Firebase
	Document           string `json:"firebase_document"`
}

// FirebaseDestination struct represents the configuration for writing data to Firebase.
type FirebaseDestination struct {
	CredentialFileAddr string `json:"firebase_credential_file"`
	Collection         string `json:"firebase_collection"`
	Document           string `json:"firebase_document"`
}

func (f FirebaseSource) FetchData(req interfaces.Request) (interface{}, error) {
	logger.Infof("Connecting to Firebase Source: Collection=%s, Document=%s, using Service Account=%s",
		req.Collection, req.Document, req.CredentialFileAddr)

	// Initialize Firebase app with service account
	opt := option.WithCredentialsFile(req.CredentialFileAddr)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	// Initialize Firestore client
	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firestore client: %w", err)
	}
	defer client.Close()

	dsnap, err := client.Collection(req.Collection).Doc(req.Document).Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch document from Firestore: %w", err)
	}

	// Check if document exists
	if !dsnap.Exists() {
		return nil, fmt.Errorf("document not found in Firestore: Collection=%s, Document=%s", req.Collection, req.Document)
	}

	data := dsnap.Data()
	logger.Infof("Data fetched from Firestore: %v", data)

	// Validate fetched data
	validatedData, err := validateFirebaseData(data)
	if err != nil {
		logger.Errorf("Validation failed for data: %v, Error: %s", data, err)
		return nil, err
	}

	// Transform data
	transformedData := transformFirebaseData(validatedData)
	logger.Infof("Data successfully processed: %v", transformedData)

	return transformedData, nil
}

// SendData connects to Firebase and writes data to the specified collection and document.
func (f FirebaseDestination) SendData(data interface{}, req interfaces.Request) error {
	logger.Infof("Writing data to Firebase database: Collection=%s, Document=%s", req.Collection, req.Document)

	// Initialize Firebase app with service account
	opt := option.WithCredentialsFile(req.CredentialFileAddr)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	// Initialize Firestore client
	client, err := app.Firestore(context.Background())
	if err != nil {
		return fmt.Errorf("failed to initialize Firestore client: %w", err)
	}
	defer client.Close()

	// Ensure data is in the expected format
	var post map[string]interface{}
	if err := convertToMap(data, &post); err != nil {
		logger.Errorf("Error converting data to Firestore format: %v", err)
		return err
	}

	// Write data to the Firestore collection/document
	_, err = client.Collection(req.Collection).NewDoc().Create(context.Background(), post)
	if err != nil {
		logger.Errorf("Error writing to Firestore: %v", err)
		return err
	}

	logger.Infof("Successfully written data to Firestore: Collection=%s, Document=%s", req.Collection, req.Document)
	return nil
}

// convertToMap converts an interface{} to a map[string]interface{} for Firestore compatibility.
func convertToMap(data interface{}, result *map[string]interface{}) error {
	temp, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	if err := json.Unmarshal(temp, result); err != nil {
		return fmt.Errorf("failed to unmarshal JSON to map: %v", err)
	}
	return nil
}

func validateFirebaseData(data map[string]interface{}) (map[string]interface{}, error) {
	logger.Infof("Validating Firebase data: %v", data)

	// Example: Ensure the "data" field exists and is a non-empty string
	message, ok := data["data"].(string)
	if !ok || strings.TrimSpace(message) == "" {
		return nil, errors.New("invalid or missing 'data' field")
	}

	return data, nil
}

// transformFirebaseData modifies the Firebase data as per business logic.
func transformFirebaseData(data map[string]interface{}) map[string]interface{} {
	logger.Infof("Transforming Firebase data: %v", data)

	// Example: Transform the "data" field to uppercase if it exists
	if message, ok := data["data"].(string); ok {
		data["data"] = strings.ToUpper(message)
	}

	// Add a new field "processed" with the current timestamp
	data["processed"] = time.Now().Format(time.RFC3339)

	return data
}

// Initialize the Firebase integrations by registering them with the registry.
func init() {
	registry.RegisterSource("Firebase", FirebaseSource{})
	registry.RegisterDestination("Firebase", FirebaseDestination{})
}
