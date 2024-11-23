package integrations

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	// "firebase.google.com/go/db"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	// "google.golang.org/grpc/internal/resolver/dns"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
)

// FirebaseSource struct represents the configuration for fetching data from Firebase.
type FirebaseSource struct {
	CredentialFileAddr string `json:"firebase_credential_file"` // Path to service account JSON file
	Collection         string `json:"firebase_collection"`     // Collection name in Firebase
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
		f.Collection, f.Document, f.CredentialFileAddr)

	// Validate configuration
	if f.CredentialFileAddr == "" || f.Collection == "" || f.Document == "" {
		return nil, errors.New("missing Firebase source configuration details")
	}

	// Initialize Firebase app with service account
	opt := option.WithCredentialsFile(f.CredentialFileAddr)
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

	// Fetch document
	logger.Infof("Fetching data from Firestore: Collection=%s, Document=%s", f.Collection, f.Document)
	dsnap, err := client.Collection(f.Collection).Doc(f.Document).Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch document from Firestore: %w", err)
	}

	// Check if document exists
	if !dsnap.Exists() {
		return nil, fmt.Errorf("document not found in Firestore: Collection=%s, Document=%s", f.Collection, f.Document)
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
}
