package registry

import "github.com/SkySingh04/fractal/interfaces"

var (
	dataSources      = make(map[string]interfaces.DataSource)
	dataDestinations = make(map[string]interfaces.DataDestination)
)

func RegisterSource(name string, source interfaces.DataSource) {
	dataSources[name] = source
}

func RegisterDestination(name string, destination interfaces.DataDestination) {
	dataDestinations[name] = destination
}

func GetSource(name string) (interfaces.DataSource, bool) {
	source, exists := dataSources[name]
	return source, exists
}

func GetDestination(name string) (interfaces.DataDestination, bool) {
	destination, exists := dataDestinations[name]
	return destination, exists
}

// GetSources returns all registered data sources
func GetSources() map[string]interfaces.DataSource {
	return dataSources
}

// GetDestinations returns all registered data destinations
func GetDestinations() map[string]interfaces.DataDestination {
	return dataDestinations
}
