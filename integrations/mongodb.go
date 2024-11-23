package integrations

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const bufferSize = 10 // Buffer size for channels

// MongoDBSource struct represents the configuration for consuming messages from MongoDB.
type MongoDBSource struct {
	ConnString string `json:"source_mongodb_conn_string"`
	Database   string `json:"source_mongodb_database"`
	Collection string `json:"source_mongodb_collection"`
}

// MongoDBDestination struct represents the configuration for publishing messages to MongoDB.
type MongoDBDestination struct {
	ConnString string `json:"target_mongodb_conn_string"`
	Database   string `json:"target_mongodb_database"`
	Collection string `json:"target_mongodb_collection"`
}

// FetchData connects to MongoDB, retrieves data, and returns it.
func (m MongoDBSource) FetchData(req interfaces.Request) (interface{}, error) {
	if req.SourceMongoDBConnString == "" || req.SourceMongoDBDatabase == "" || req.SourceMongoDBCollection == "" {
		return nil, errors.New("missing MongoDB source connection details")
	}
	logger.Infof("Connecting to MongoDB source...")

	clientOptions := options.Client().ApplyURI(req.SourceMongoDBConnString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database(req.SourceMongoDBDatabase).Collection(req.SourceMongoDBCollection)

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Use a buffered channel to handle data
	dataChannel := make(chan bson.M, bufferSize)
	var wg sync.WaitGroup
	var allResults []bson.M
	var fetchError error
	mutex := &sync.Mutex{}

	// Goroutine to collect data from the channel
	go func() {
		for doc := range dataChannel {
			mutex.Lock()
			allResults = append(allResults, doc)
			mutex.Unlock()
		}
	}()

	// Process each document
	for cursor.Next(context.TODO()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			fetchError = err
			break
		}

		wg.Add(1)
		go func(d bson.M) {
			defer wg.Done()
			dataChannel <- d
		}(doc)
	}

	wg.Wait()
	close(dataChannel)

	if fetchError != nil {
		return nil, fetchError
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	logger.Infof("Data fetched from MongoDB: %d documents", len(allResults))
	return allResults, nil
}

// SendData connects to MongoDB and publishes data to the specified collection.
func (m MongoDBDestination) SendData(data interface{}, req interfaces.Request) error {
	if req.TargetMongoDBConnString == "" || req.TargetMongoDBDatabase == "" || req.TargetMongoDBCollection == "" {
		return errors.New("missing MongoDB target connection details")
	}
	logger.Infof("Connecting to MongoDB destination...")

	clientOptions := options.Client().ApplyURI(req.TargetMongoDBConnString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			logger.Errorf("Error disconnecting MongoDB client: %v", err)
		}
	}()

	collection := client.Database(req.TargetMongoDBDatabase).Collection(req.TargetMongoDBCollection)

	// Assert that data is a slice of bson.M
	dataSlice, ok := data.([]bson.M)
	if !ok {
		logger.Errorf("data must be a slice of bson.M representing documents")
		return errors.New("invalid data format: expected []bson.M")
	}

	// Buffered channel for sending documents
	dataChannel := make(chan bson.M, bufferSize)
	errorChannel := make(chan error, bufferSize)
	var wg sync.WaitGroup

	// Goroutines for worker pool
	for i := 0; i < bufferSize; i++ { // Worker pool
		wg.Add(1)
		go func() {
			defer wg.Done()
			for doc := range dataChannel {
				if _, err := collection.InsertOne(context.TODO(), doc); err != nil {
					logger.Errorf("Error inserting into collection %s: %v", req.TargetMongoDBCollection, err)
					errorChannel <- err
				} else {
					logger.Infof("Data sent to MongoDB target collection %s: %v", req.TargetMongoDBCollection, doc)
				}
			}
		}()
	}

	// Feed data into the channel
	go func() {
		for _, doc := range dataSlice {
			dataChannel <- doc
		}
		close(dataChannel)
	}()

	wg.Wait()
	close(errorChannel)

	// Check for errors in the error channel
	if len(errorChannel) > 0 {
		return errors.New("one or more errors occurred while inserting data")
	}

	return nil
}

// Initialize the MongoDB integrations by registering them with the registry.
func init() {
	registry.RegisterSource("MongoDB", MongoDBSource{})
	registry.RegisterDestination("MongoDB", MongoDBDestination{})
}
