package main

import (
	"fmt"

	"github.com/SkySingh04/fractal/config"
	"github.com/SkySingh04/fractal/controller"
	_ "github.com/SkySingh04/fractal/integrations"
	"github.com/SkySingh04/fractal/logger"
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
			configuration = make(map[string]string)
			for key, value := range configMap {
				if strValue, ok := value.(string); ok {
					configuration[key] = strValue
				} else {
					logger.Fatalf("Invalid configuration value for key %s: %v", key, value)
				}
			}
			if err != nil {
				logger.Fatalf("Failed to set up configuration: %v", err)
			}
		}

		logger.Infof("Configuration loaded successfully: %+v", configuration)
		// Here you can add further CLI-based data processing logic
	}
}
