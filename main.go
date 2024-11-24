package main

import (
	"context"
	"fmt"
	"time"

	"github.com/SkySingh04/fractal/config"
	"github.com/SkySingh04/fractal/controller"
	_ "github.com/SkySingh04/fractal/integrations"
	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/opentele"
	"github.com/SkySingh04/fractal/registry"
	"gofr.dev/pkg/gofr"
)

const (
	logo = `

	███████╗██████╗  █████╗  ██████╗████████╗ █████╗ ██╗     
	██╔════╝██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██║     
	█████╗  ██████╔╝███████║██║        ██║   ███████║██║     
	██╔══╝  ██╔══██╗██╔══██║██║        ██║   ██╔══██║██║     
	██║     ██║  ██║██║  ██║╚██████╗   ██║   ██║  ██║███████╗
	╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝
	`
)

func fetchData() ([]string, error) {
	// Simulate data fetching logic (e.g., from an API or database)
	fetchedData := []string{"Record 1", "Record 2", "Record 3"}
	if time.Now().Minute()%2 == 0 {
		return nil, fmt.Errorf("failed to fetch data")
	}
	return fetchedData, nil
}

func processData(data []string) ([]string, error) {
	// Simulate data processing (e.g., transforming, filtering)
	var processedData []string
	for _, record := range data {
		if record != "Record 2" {
			processedData = append(processedData, record)
		}
	}
	return processedData, nil
}

func main() {
	// Initialize OpenTelemetry tracing
	cleanup, err := opentele.InitTracing()
	if err != nil {
		logger.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer cleanup() // Ensure resources are flushed on exit
	app := gofr.New()
	fmt.Print(logo)

	// Ask for mode selection
	mode, err := config.AskForMode()
	if err != nil {
		logger.Fatalf("Failed to select application mode: %v", err)
	}
	// Ask the user for the cron job repeat interval in seconds
	var intervalSec int
	fmt.Print("Enter the interval for cron job to repeat (in seconds): ")
	_, err = fmt.Scanf("%d", &intervalSec)
	if err != nil || intervalSec <= 0 {
		logger.Fatalf("Invalid interval input. Please enter a positive integer: %v", err)
	}

	if mode == "Start HTTP Server" {
		logger.Infof("Starting HTTP Server... Welcome to the Fractal API!")

		// Register route greet
		app.GET("/greet", func(ctx *gofr.Context) (interface{}, error) {
			// Start a span for this route
			_, span := opentele.CreateSpan(ctx.Context, "HTTP GET /greet")
			defer span.End()

			// Perform the route logic
			return "Hello Fractal!", nil
		})

		// Register other routes as necessary
		app.POST("/api/migration", controller.MigrationHandler)

		// Default port 8000
		app.Run()
	} else if mode == "Use CLI" {
		// CLI Mode Logic
		// Load configuration
		configuration, err := config.LoadConfig("config.yaml")
		if err != nil {
			
			logger.Infof("Config file not found. Setup interactively: %v", err)

			configMap, err := config.SetupConfigInteractively()

			if err != nil { logger.Fatalf(`Failed to setup configuration interactively:`, err) }

			configuration = make(map[string]interface{})

			for key, value := range configMap {
				if strValue, ok := value.(string); ok {
					configuration[key] = strValue
				} else {
					logger.Fatalf("Invalid value for key %s: %v", key, value)
				}
			}
			if err != nil { logger.Fatalf(`Failed to setup configuration interactively:`, err) }
			
			
		}
		logger.Infof("Configuration loaded: %+v", configuration)

		// Define the task to be executed
		task := func() {
			// Create a root span for the entire task
			ctx, span := opentele.CreateSpan(context.Background(), "cron-job")
			defer span.End()

			logger.Infof("Cron job triggered at: %s", time.Now().Format(time.RFC3339))

			// Your task logic (e.g., data fetch and send to CSV)
			inputMethod, inputconfig := configuration["inputMethod"], configuration["inputconfig"].(map[string]interface{})
			outputMethod, outputconfig := configuration["outputMethod"], configuration["outputconfig"].(map[string]interface{})

			// Fetch data from input integration
			_, fetchSpan := opentele.CreateSpan(ctx, "fetch-data")
			inputIntegration, found := registry.GetSource(inputMethod.(string))
			if !found {
				fetchSpan.RecordError(fmt.Errorf("input method %s not registered", inputMethod))
				fetchSpan.End()
				logger.Fatalf("Input method %s not registered", inputMethod)
			}
			inputRequest := mapConfigToRequest(inputconfig)
			data, err := inputIntegration.FetchData(inputRequest)
			if err != nil {
				fetchSpan.RecordError(err)
				fetchSpan.End()
				logger.Fatalf("Failed to fetch data from %s: %v", inputMethod, err)
			}
			fetchSpan.End()

			// Send data to output integration
			_, sendSpan := opentele.CreateSpan(ctx, "send-data")
			outputIntegration, found := registry.GetDestination(outputMethod.(string))
			if !found {
				sendSpan.RecordError(fmt.Errorf("output method %s not registered", outputMethod))
				sendSpan.End()
				logger.Fatalf("Output method %s not registered", outputMethod)
			}
			outputRequest := mapConfigToRequest(outputconfig)
			err = outputIntegration.SendData(data, outputRequest)
			if err != nil {
				sendSpan.RecordError(err)
				sendSpan.End()
				logger.Fatalf("Failed to send data to %s: %v", outputMethod, err)
			}
			sendSpan.End()

			logger.Infof("Data sent successfully")
		}

		// Run the task immediately
		task()

		// Repeat the task every interval
		ticker := time.NewTicker(time.Duration(intervalSec) * time.Second) // Adjust the interval as needed
		defer ticker.Stop()

		// Infinite loop to keep executing the task every interval
		for range ticker.C {
			// Execute the task on each tick
			task()
		}

	}
}

func getStringField(config map[string]interface{}, field string, defaultValue string) string {
	if value, ok := config[field]; ok && value != nil {
		return value.(string)
	}
	return defaultValue
}

func mapConfigToRequest(config map[string]interface{}) interfaces.Request {

	return interfaces.Request{
		Input:                   getStringField(config, "inputmethod", ""),
		Output:                  getStringField(config, "outputmethod", ""),
		RabbitMQInputURL:        getStringField(config, "url", ""),
		RabbitMQInputQueueName:  getStringField(config, "queuename", ""),
		RabbitMQOutputURL:       getStringField(config, "url", ""),
		RabbitMQOutputQueueName: getStringField(config, "queuename", ""),
		ConsumerURL:             getStringField(config, "url", ""),
		ConsumerTopic:           getStringField(config, "topic", ""), // Default is empty if "topic" is missing
		ProducerURL:             getStringField(config, "url", ""),
		ProducerTopic:           getStringField(config, "topic", ""),
		SQLSourceConnString:     getStringField(config, "connstring", ""),
		SQLTargetConnString:     getStringField(config, "connstring", ""),
		SourceMongoDBConnString: getStringField(config, "connstring", ""),
		SourceMongoDBDatabase:   getStringField(config, "database", ""),
		SourceMongoDBCollection: getStringField(config, "collection", ""),
		TargetMongoDBConnString: getStringField(config, "connstring", ""),
		TargetMongoDBDatabase:   getStringField(config, "database", ""),
		TargetMongoDBCollection: getStringField(config, "collection", ""),
		OutputFileName:          getStringField(config, "filename", ""),
		CSVSourceFileName:       getStringField(config, "csvsourcefilename", ""),
		CSVDestinationFileName:  getStringField(config, "csvdestinationfilename", ""),
		JSONSourceData:          getStringField(config, "data", ""),
		JSONOutputFilename:      getStringField(config, "filename", ""),
		YAMLSourceFilePath:      getStringField(config, "filepath", ""),
		YAMLDestinationFilePath: getStringField(config, "filepath", ""),
		DynamoDBSourceTable:     getStringField(config, "tablename", ""),
		DynamoDBTargetTable:     getStringField(config, "tablename", ""),
		DynamoDBSourceRegion:    getStringField(config, "region", ""),
		DynamoDBTargetRegion:    getStringField(config, "region", ""),
		FTPURL:                  getStringField(config, "url", ""),
		FTPUser:                 getStringField(config, "user", ""),
		FTPPassword:             getStringField(config, "password", ""),
		SFTPURL:                 getStringField(config, "url", ""),
		SFTPUser:                getStringField(config, "user", ""),
		SFTPPassword:            getStringField(config, "password", ""),
		WebSocketSourceURL:      getStringField(config, "url", ""),
		WebSocketDestURL:        getStringField(config, "url", ""),
		CredentialFileAddr: 	getStringField(config, "credentialfileaddr", ""),
		Collection: 		getStringField(config, "collection", ""),
		Document: 		getStringField(config, "document", ""),
	}
}
