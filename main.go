package main

import (
	"fmt"

	"github.com/SkySingh04/fractal/config"
	"github.com/SkySingh04/fractal/controller"
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
		app.Logger().Fatalf("Failed to select application mode: %v", err)
	}

	if mode == "Start HTTP Server" {
		app.Logger().Infof("Starting HTTP Server... Welcome to the Fractal API!")

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
			app.Logger().Logf("Config file not found. Let's set up the input and output methods.")
			configuration, err = config.SetupConfigInteractively()
			if err != nil {
				app.Logger().Fatalf("Failed to set up configuration: %v", err)
			}
		}

		app.Logger().Infof("Configuration loaded successfully: %+v", configuration)
		// Here you can add further CLI-based data processing logic
	}
}
