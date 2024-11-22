package integrations

import (
	"errors"
	"log"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
)

type MongoDBSource struct {
	ConnString string `json:"source_mongodb_conn_string"`
	Database   string `json:"source_mongodb_database"`
	Collection string `json:"source_mongodb_collection"`
}

type MongoDBDestination struct {
	ConnString string `json:"target_mongodb_conn_string"`
	Database   string `json:"target_mongodb_database"`
	Collection string `json:"target_mongodb_collection"`
}

func (m MongoDBSource) FetchData(req interfaces.Request) (interface{}, error) {
	if req.SourceMongoDBConnString == "" || req.SourceMongoDBDatabase == "" || req.SourceMongoDBCollection == "" {
		return nil, errors.New("missing MongoDB source connection details")
	}
	log.Println("Fetching data from MongoDB source...")
	// Add MongoDB fetch logic here
	return "MongoDBData", nil
}

func (m MongoDBDestination) SendData(data interface{}, req interfaces.Request) error {
	if req.TargetMongoDBConnString == "" || req.TargetMongoDBDatabase == "" || req.TargetMongoDBCollection == "" {
		return errors.New("missing MongoDB target connection details")
	}
	log.Println("Sending data to MongoDB target...")
	// Add MongoDB send logic here
	return nil
}

func init() {
	registry.RegisterSource("MongoDB", MongoDBSource{})
	registry.RegisterDestination("MongoDB", MongoDBDestination{})
}
