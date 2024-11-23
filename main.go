package main

import (
	"fmt"
	"time"

	"github.com/SkySingh04/fractal/config"
	"github.com/SkySingh04/fractal/controller"
	_ "github.com/SkySingh04/fractal/integrations"
	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
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
	app := gofr.New()
	fmt.Print(logo)

	// Ask if the user wants to start HTTP Server or use CLI
	mode, err := config.AskForMode()
	if err != nil {
		logger.Fatalf("Failed to select application mode: %v", err)
	}

	// Add cron job to fetch and process data every 5 hours
	app.AddCronJob("* */5 * * *", "Data Fetch and Process", func(ctx *gofr.Context) {
		//  1: Fetch data
		data, err := fetchData()
		if err != nil {
			ctx.Logger.Errorf("Error fetching data: %v", err)
			return
		}
		ctx.Logger.Infof("Fetched %d records.", len(data))

		//  2: Process the data
		processedData, err := processData(data)
		if err != nil {
			ctx.Logger.Errorf("Error processing data: %v", err)
			return
		}

		ctx.Logger.Infof("Processed %d records.", len(processedData))
		ctx.Logger.Infof("Data fetching and processing completed successfully.")
	})

	if mode == "Start HTTP Server" {
		logger.Infof("Starting HTTP Server... Welcome to the Fractal API!")

		// Register routes
		app.GET("/greet", func(ctx *gofr.Context) (interface{}, error) {
			return "Hello World!", nil
		})

		// Register other routes as needed
		app.POST("/api/migration", controller.MigrationHandler)

		// Run HTTP server
		app.Run()

	} else if mode == "Use CLI" {
		// Load or set up the configuration for CLI mode
		configuration, err := config.LoadConfig("config.yaml")
		if err != nil {
			logger.Logf("Config file not found. Let's set up the input and output methods.")
			configMap, err := config.SetupConfigInteractively()
			if err != nil {
				logger.Fatalf("Failed to set up configuration: %v", err)
			}
			configuration := make(map[string]interface{})
			for key, value := range configMap {
				switch v := value.(type) {
				case string:
					configuration[key] = v
				case map[string]interface{}:
					logger.Logf("Key %s has a nested map value: %v", key, v)
					configuration[key] = v
				default:
					logger.Logf("Key %s has a value of unhandled type %T: %v", key, v, v)
					configuration[key] = v
				}
			}
		}
		logger.Infof("Configuration loaded successfully: %+v", configuration)
		if _, ok := configuration["inputconfig"]; !ok {
			logger.Fatalf("Missing 'inputconfig' in configuration")
		}

		if _, ok := configuration["outputconfig"]; !ok {
			logger.Fatalf("Missing 'outputconfig' in configuration")
		}

		// Get the input and output methods from the configuration
		inputMethod, inputconfig := configuration["inputMethod"], configuration["inputconfig"].(map[string]interface{})
		outputMethod, outputconfig := configuration["outputMethod"], configuration["outputconfig"].(map[string]interface{})

		// Fetch data from the input method
		inputIntegration, found := registry.GetSource(inputMethod.(string))
		if !found {
			logger.Fatalf("Input method %s not registered", inputMethod)
		}

		inputRequest := mapConfigToRequest(inputconfig)
		data, err := inputIntegration.FetchData(inputRequest)
		if err != nil {
			logger.Fatalf("Failed to fetch data from %s: %v", inputMethod, err)
		}

		// Send data to the output method
		outputIntegration, found := registry.GetDestination(outputMethod.(string))
		if !found {
			logger.Fatalf("Output method %s not registered", outputMethod)
		}

		logger.Infof("Sending data to %s...", outputMethod)
		logger.Infof("Output configuration: %+v", outputconfig)

		outputRequest := mapConfigToRequest(outputconfig)
		err = outputIntegration.SendData(data, outputRequest)
		if err != nil {
			logger.Fatalf("Failed to send data to %s: %v", outputMethod, err)
		}

		logger.Infof("Data sent to %s successfully", outputMethod)
	}
}

func getStringField(config map[string]interface{}, field string, defaultValue string) string {
	if value, ok := config[field]; ok && value != nil {
		return value.(string)
	}
	return defaultValue
}

	return interfaces.Request{
		Input:                   getStringField(config, "inputmethod", ""),
		Output:                  getStringField(config, "outputmethod", ""),
		RabbitMQInputURL:        getStringField(config, "url", ""),
		RabbitMQInputQueueName:  getStringField(config, "queuename", ""),
		RabbitMQOutputURL:       getStringField(config, "url", ""),
		RabbitMQOutputQueueName: getStringField(config, "queuename", ""),
		ConsumerURL:             getStringField(config, "url", ""),
		ConsumerTopic:           getStringField(config, "topic", ""),
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
		FTPFILEPATH:             getStringField(config, "ftpfilepath", ""),
		FTPURL:                  getStringField(config, "url", ""),
		FTPUser



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
		FTPFILEPATH:             getStringField(config, "ftpfilepath", ""),
		FTPURL:                  getStringField(config, "url", ""),
		FTPUser:                 getStringField(config, "user", ""),
		FTPPassword:             getStringField(config, "password", ""),
		SFTPFILEPATH:            getStringField(config, "sftpfilepath", ""),
		SFTPURL:                 getStringField(config, "url", ""),
		SFTPUser:                getStringField(config, "user", ""),
		SFTPPassword:            getStringField(config, "password", ""),
		WebSocketSourceURL:      getStringField(config, "url", ""),
		WebSocketDestURL:        getStringField(config, "url", ""),
		CredentialFileAddr:      getStringField(config, "credentialfileaddr", "firebaseConfig.json"),
		Document:                getStringField(config, "document", "sampledata"),
		Collection:              getStringField(config, "collection", "1"),
	}
}
