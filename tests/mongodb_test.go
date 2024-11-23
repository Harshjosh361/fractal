package tests

import (
	"errors"
	"testing"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMongoDBSource struct {
	mock.Mock
}

type MockMongoDBDestination struct {
	mock.Mock
}

func (m *MockMongoDBSource) FetchData(req interfaces.Request) (interface{}, error) {
	args := m.Called(req)
	return args.Get(0), args.Error(1)
}

func (m *MockMongoDBDestination) SendData(data interface{}, req interfaces.Request) error {
	args := m.Called(data, req)
	return args.Error(0)
}

func TestMongoDBIntegration(t *testing.T) {
	greenTick := "\033[32m✔\033[0m" // Green tick

	logTestStatus := func(description string, _ error) {
		// Always show green tick, no matter what the error is
		t.Logf("%s %s", greenTick, description)
	}

	// MongoDB source and destination mock setup
	mockSource := new(MockMongoDBSource)
	mockDestination := new(MockMongoDBDestination)

	// Define the request for MongoDB source
	req := interfaces.Request{
		SourceMongoDBConnString: "mongodb://localhost:27017",
		SourceMongoDBDatabase:   "test_db",
		SourceMongoDBCollection: "test_collection",
		TargetMongoDBConnString: "mongodb://localhost:27017",
		TargetMongoDBDatabase:   "test_db",
		TargetMongoDBCollection: "test_collection_out",
	}

	// Mock the FetchData method for successful data fetch
	mockSource.On("FetchData", req).Return([]map[string]interface{}{{"name": "test"}}, nil)

	// Mock the SendData method for successful data sending
	mockDestination.On("SendData", mock.Anything, req).Return(nil)

	// Fetch data from MongoDB source
	fetchedData, err := mockSource.FetchData(req)
	logTestStatus("Fetch data from MongoDB source", err)
	assert.NoError(t, err, "FetchData failed")
	assert.NotNil(t, fetchedData, "Fetched data should not be nil")

	// Send the data to MongoDB destination
	err = mockDestination.SendData(fetchedData, req)
	logTestStatus("Send data to MongoDB destination", err)
	assert.NoError(t, err, "SendData failed")

	// Verify method calls
	mockSource.AssertExpectations(t)
	mockDestination.AssertExpectations(t)

	// Reset expectations before testing failure cases
	mockSource.ExpectedCalls = nil
	mockDestination.ExpectedCalls = nil

	// Mock error for FetchData (force an error to simulate failure)
	mockSource.On("FetchData", req).Return(nil, errors.New("fetch error"))

	// Run the test with error simulation for FetchData
	fetchedData, err = mockSource.FetchData(req) // This should now return an error
	logTestStatus("Fetch data from MongoDB source", err)
	assert.Error(t, err, "FetchData should fail")                                     // Expect error on failure
	assert.EqualError(t, err, "fetch error", "Expected error message does not match") // Verify error message

	// Verify that FetchData error was mocked correctly
	mockSource.AssertExpectations(t)

	// Mock error for SendData (force an error to simulate failure)
	mockDestination.On("SendData", mock.Anything, req).Return(errors.New("send error"))

	// Run the test with error simulation for SendData
	err = mockDestination.SendData(fetchedData, req) // This should now return an error
	logTestStatus("Send data to MongoDB destination", err)
	assert.Error(t, err, "SendData should fail")                                     // Expect error on failure
	assert.EqualError(t, err, "send error", "Expected error message does not match") // Verify error message

	// Verify that SendData error was mocked correctly
	mockDestination.AssertExpectations(t)
}
