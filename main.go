package main

import (
	"fmt"

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
 ______ _____            _____ _______       _      
 |  ____|  __ \     /\   / ____|__   __|/\   | |     
 | |__  | |__) |   /  \ | |       | |  /  \  | |     
 |  __| |  _  /   / /\ \| |       | | / /\ \ | |     
 | |    | | \ \  / ____ \ |____   | |/ ____ \| |____ 
 |_|    |_|  \_\/_/    \_\_____|  |_/_/    \_\______|
`
)

func main() {
	app := gofr.New()
	fmt.Print(logo)

	// Ask if the user wants to start HTTP Server or use CLI
	mode, err := config.AskForMode()
	if err != nil {
		logger.Fatalf("Failed to select application mode: %v", err)
	}

	if mode == "Start HTTP Server" {
		logger.Infof("Starting HTTP Server... Welcome to the Fractal API!")

		// Register route greet
		app.GET("/greet", func(ctx *gofr.Context) (interface{}, error) {
			return "Hello World!", nil
		})

		// Register other routes as necessary
		app.POST("/api/migration", controller.MigrationHandler)

		// Default port 8000
		app.Run()
	} else if mode == "Use CLI" {
		// Load or set up the configuration interactively for CLI mode
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
					configuration[key] = v // Optionally handle other types here
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

		// logger.Infof("Configuration loaded successfully: %+v", configuration)

		// Get the input and output methods from the configuration
		inputMethod, inputconfig := configuration["inputMethod"], configuration["inputconfig"].(map[string]interface{})
		outputMethod, outputconfig := configuration["outputMethod"], configuration["outputconfig"].(map[string]interface{})

		// Fetch data from the input method
		inputIntegration, found := registry.GetSource(inputMethod.(string))
		if !found {
			logger.Fatalf("Input method %s not registered", inputMethod)
		}

		// logger.Infof("Fetching data from %s...", inputMethod)
		// logger.Infof("Input configuration: %+v", inputconfig)

		inputRequest := mapConfigToRequest(inputconfig)
		data, err := inputIntegration.FetchData(inputRequest)

		if err != nil {
			logger.Fatalf("Failed to fetch data from %s: %v", inputMethod, err)
		}

		// logger.Infof("Data fetched from %s: %v", inputMethod, data)

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
		JSONOutputData:          getStringField(config, "data", ""),
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
	}
}
