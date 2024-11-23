package integrations

import (
	"context"
	"errors"
	"log"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	var allResults []bson.M
	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		allResults = append(allResults, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	logger.Infof("Data fetched from MongoDB: %v", allResults)
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

	for _, row := range dataSlice {
		if _, err := collection.InsertOne(context.TODO(), row); err != nil {
			logger.Errorf("Error inserting into collection %s: %v", req.TargetMongoDBCollection, err)
			return err
		}
		logger.Infof("Data sent to MongoDB target collection %s: %v", req.TargetMongoDBCollection, row)
	}

	return nil
}

// Initialize the MongoDB integrations by registering them with the registry.
func init() {
	registry.RegisterSource("MongoDB", MongoDBSource{})
	registry.RegisterDestination("MongoDB", MongoDBDestination{})
}
