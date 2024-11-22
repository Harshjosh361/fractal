package integrations

import (
	"errors"
	"log"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
)

type SQLSource struct {
	ConnString string `json:"sql_source_conn_string"`
}

type SQLDestination struct {
	ConnString string `json:"sql_target_conn_string"`
}

func (s SQLSource) FetchData(req interfaces.Request) (interface{}, error) {
	if req.SQLSourceConnString == "" {
		return nil, errors.New("missing SQL source connection string")
	}
	log.Println("Fetching data from SQL source...")
	// Add SQL fetching logic here
	return "SQLData", nil
}

func (s SQLDestination) SendData(data interface{}, req interfaces.Request) error {
	if req.SQLTargetConnString == "" {
		return errors.New("missing SQL target connection string")
	}
	log.Println("Sending data to SQL target...")
	// Add SQL sending logic here
	return nil
}

func init() {
	registry.RegisterSource("SQL", SQLSource{})
	registry.RegisterDestination("SQL", SQLDestination{})
}
