package main

import (
	"fmt"

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
	fmt.Println(logo)

	// Print greeting message when the server starts
	app.Logger().Infof("Server is starting... Welcome to the Fractal API!")

	// register route greet
	app.GET("/greet", func(ctx *gofr.Context) (interface{}, error) {
		return "Hello World!", nil
	})

	// Routes
	app.POST("/api/migration", controller.MigrationHandler)

	//  default port 8000.
	app.Run()
}
