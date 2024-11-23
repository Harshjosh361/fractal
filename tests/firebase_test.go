package tests

import (
	"context"
	"fmt"
	"github.com/SkySingh04/fractal/integrations"
	"github.com/SkySingh04/fractal/interfaces"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time" // Import the time package
)

type MockFirestoreClient struct {
	mock.Mock
}

func (m *MockFirestoreClient) Collection(collection string) *MockFirestoreCollection {
	args := m.Called(collection)
	return args.Get(0).(*MockFirestoreCollection)
}

type MockFirestoreCollection struct {
	mock.Mock
}

func (m *MockFirestoreCollection) Doc(document string) *MockFirestoreDocument {
	args := m.Called(document)
	return args.Get(0).(*MockFirestoreDocument)
}

type MockFirestoreDocument struct {
	mock.Mock
}

func (m *MockFirestoreDocument) Get(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestFirebaseIntegration(t *testing.T) {

	const (
		GreenTick = "\033[32m✔\033[0m" // Green tick
		RedCross  = "\033[31m✘\033[0m" // Red cross
	)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFirestoreClient := new(MockFirestoreClient)
	mockFirestoreCollection := new(MockFirestoreCollection)
	mockFirestoreDocument := new(MockFirestoreDocument)

	// Setting up mocks
	mockFirestoreClient.On("Collection", "testCollection").Return(mockFirestoreCollection)
	mockFirestoreCollection.On("Doc", "testDocument").Return(mockFirestoreDocument)
	mockFirestoreDocument.On("Get", context.Background()).Return(map[string]interface{}{
		"data": "test data", // Mock the data for "testDocument"
	}, nil) // Ensure no error

	firebaseSource := integrations.FirebaseSource{
		CredentialFileAddr: "firebaseConfig.json",
		Collection:         "testCollection",
		Document:           "testDocument",
	}

	req := interfaces.Request{
		Collection:         "testCollection",
		Document:           "testDocument",
		CredentialFileAddr: "firebaseConfig.json",
	}

	// Fetch data from Firebase
	data, err := firebaseSource.FetchData(req)

	// Assertions
	if assert.NoError(t, err) {
		fmt.Printf("%s FetchData passed\n", GreenTick)
	} else {
		fmt.Printf("%s FetchData failed\n", RedCross)
	}

	// Assert that the data is of type map[string]interface{}
	resultData, ok := data.(map[string]interface{})

	if assert.True(t, ok) {
		fmt.Printf("%s Data type validation passed\n", GreenTick)
	} else {
		fmt.Printf("%s Data type validation failed\n", RedCross)
	}

	// Check the transformed data (data should be "test data")
	if assert.Equal(t, "TEST DATA", resultData["data"]) {
		fmt.Printf("%s Data validation passed\n", GreenTick)
	} else {
		fmt.Printf("%s Data validation failed\n", RedCross)
	}// Parse the "processed" time string
parsedTime, err := time.Parse(time.RFC3339, resultData["processed"].(string))
if assert.NoError(t, err) {
    // Check the processed field (it should match the current time within a small tolerance)
    if assert.WithinDuration(t, time.Now(), parsedTime, time.Second) {
        fmt.Printf("%s Processed field validation passed\n", GreenTick)
    } else {
        fmt.Printf("%s Processed field validation failed\n", RedCross)
    }
} else {
    fmt.Printf("%s Error parsing processed time: %v\n", RedCross, err)
}

}
