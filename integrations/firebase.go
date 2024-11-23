package integrations

import (
	"context"
	"errors"

	"firebase.google.com/go"
	// "firebase.google.com/go/db"
	// "github.com/SkySingh04/fractal/firebaseSetup"
	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
	"google.golang.org/api/option"
)

// FirebaseSource struct represents the configuration for fetching data from Firebase.
type FirebaseSource struct {
	FirebaseURL string `json:"firebase_database_url"`
	Credential  string `json:"firebase_credential_file"`
}

type FirebaseDestination struct {
	FirebaseURL string `json:"firebase_database_url"`
	Credential string  `json:"firebase_credential_file"`
}

// FetchData connects to Firebase Realtime Database, retrieves data, and passes it through validation and transformation pipelines.
func (f FirebaseSource) FetchData(req interfaces.Request) (interface{}, error) {
	logger.Infof("Connecting to Firebase Source: URL=%s, using service Account=%s", req.FirebaseURL, req.Credential)

	if req.FirebaseURL == "" || req.Credential == "" {
		return nil, errors.New("missing Firebase source details")
	}

	// Initialize Firebase app with service account
	opt := option.WithCredentialsFile(req.FirebaseURL)
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		DatabaseURL: req.FirebaseURL,
	}, opt)
	if err != nil {
		return nil, err
	}

	// Initialize the Database client
	client, err := app.Database(context.Background())
	if err != nil {
		return nil, err
	}

	// Specify the path to fetch data from (e.g., "messages")
	var fetchedData map[string]interface{}
	err = client.NewRef("messages").Get(context.Background(), &fetchedData)
	if err != nil {
		return nil, err
	}

	logger.Infof("Data fetched from Firebase: %v", fetchedData)

	// Validation
	validatedData, err := validateFirebaseData(fetchedData)
	if err != nil {
		logger.Fatalf("Validation failed for data: %v, Error: %s", fetchedData, err)
		return nil, err
	}

	// Transformation
	transformedData := transformFirebaseData(validatedData)

	logger.Infof("Data successfully processed: %v", transformedData)
	return transformedData, nil
}

// validateFirebaseData ensures the Firebase data meets the required criteria.
func validateFirebaseData(data map[string]interface{}) (map[string]interface{}, error) {
	logger.Infof("Validating Firebase data: %v", data)

	// Example: Ensure the "message" field exists and is not empty
	if message, ok := data["message"].(string); !ok || message == "" {
		return nil, errors.New("invalid or missing 'message' field")
	}

	return data, nil
}

// transformFirebaseData modifies the Firebase data as per business logic.
func transformFirebaseData(data map[string]interface{}) map[string]interface{} {
	logger.Infof("Transforming Firebase data: %v", data)

	// Example: Add a new field "processed" with a timestamp
	data["processed"] = true
	return data
}

// Initialize the Firebase integrations by registering them with the registry.
func init() {
	registry.RegisterSource("Firebase", FirebaseSource{})
}
