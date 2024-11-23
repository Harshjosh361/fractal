package factory

import (
	"fmt"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
)

func CreateSource(name string) (interfaces.DataSource, error) {
	source, exists := registry.GetSource(name)
	if !exists {
		return nil, fmt.Errorf("source %s not found", name)
	}
	return source, nil
}

func CreateDestination(name string) (interfaces.DataDestination, error) {
	destination, exists := registry.GetDestination(name)
	if !exists {
		return nil, fmt.Errorf("destination %s not found", name)
	}
	return destination, nil
}
