package controller

import (
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
		return nil, err
	}
	return runMigration(req)
}

func runMigration(req interfaces.Request) (interface{}, error) {
	input, err := factory.CreateSource(req.Input)
	if err != nil {
		return nil, err
	}
	output, err := factory.CreateDestination(req.Output)
	if err != nil {
		return nil, err
	}

	data, err := input.FetchData(req)
	if err != nil {
		return nil, err
	}

	if err := output.SendData(data, req); err != nil {
		return nil, err
	}

	return map[string]string{"status": "success"}, nil
}
