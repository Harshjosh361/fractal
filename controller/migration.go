package controller

import (
	"fmt"
	"log"

	"github.com/SkySingh04/fractal/factory"
	"github.com/SkySingh04/fractal/interfaces"
	"gofr.dev/pkg/gofr"
)

func RegisterRoutes(app *gofr.App) {
	app.POST("/migrate", MigrationHandler)
}

func MigrationHandler(ctx *gofr.Context) (interface{}, error) {
	var req interfaces.Request
	if err := ctx.Bind(&req); err != nil {
		// Log detailed error to understand the bind issue
		return nil, fmt.Errorf("failed to bind request: %v", err)
	}
	return runMigration(req)
}

func runMigration(req interfaces.Request) (interface{}, error) {
	// Create source
	input, err := factory.CreateSource(req.Input)
	if err != nil {
		log.Printf("Error creating source for input method %s: %v", req.Input, err)
		return nil, fmt.Errorf("failed to create source for input method %s: %v", req.Input, err)
	}

	// Create destination
	output, err := factory.CreateDestination(req.Output)
	if err != nil {
		log.Printf("Error creating destination for output method %s: %v", req.Output, err)
		return nil, fmt.Errorf("failed to create destination for output method %s: %v", req.Output, err)
	}

	// Fetch data from the source
	data, err := input.FetchData(req)
	if err != nil {
		log.Printf("Error fetching data from source: %v", err)
		return nil, fmt.Errorf("failed to fetch data from source: %v", err)
	}

	// Send data to the destination
	if err := output.SendData(data, req); err != nil {
		log.Printf("Error sending data to destination: %v", err)
		return nil, fmt.Errorf("failed to send data to destination: %v", err)
	}

	log.Println("Migration successful!")
	return map[string]string{"status": "success"}, nil
}
