package main

import (
	"github.com/example/controller"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	// register route greet
	app.GET("/greet", func(ctx *gofr.Context) (interface{}, error) {

		return "Hello World!", nil
	})

	// Routes
	app.POST("/api/migration", controller.MigrationHandler)

	//  default port 8000.
	app.Run()
}
