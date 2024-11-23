package integrations

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
)

type JSONSource struct {
	Data string `json:"json_source_data"`
}

type JSONDestination struct {
	Filename string `json:"json_output_filename"`
}

// FetchData retrieves and processes JSON source data
func (j JSONSource) FetchData(req interfaces.Request) (interface{}, error) {
	if req.JSONSourceData == "" {
		return nil, errors.New("missing JSON source data")
	}

	// Validate and sanitize JSON data
	validatedData, err := ValidateJSONData(req.JSONSourceData)
	if err != nil {
		logger.Fatalf("Validation error: %v", err)
		return nil, err
	}

	// Transform JSON data
	transformedData, err := transformJSONData(validatedData)
	if err != nil {
		logger.Fatalf("Transformation error: %v", err)
		return nil, err
	}

	return transformedData, nil
}

// SendData writes JSON data to a destination file
func (j JSONDestination) SendData(data interface{}, req interfaces.Request) error {
	if req.JSONOutputFilename == "" {
		return errors.New("missing JSON destination filename")
	}

	logger.Infof("Sending data to JSON destination...")
	logger.Infof("Data: %v", data)

	// Write data to a JSON file
	err := writeJSONFile(req.JSONOutputFilename, data)
	if err != nil {
		logger.Fatalf("Error writing data to JSON file: %v", err)
		return err
	}

	logger.Infof("Data successfully written to %s", req.JSONOutputFilename)
	return nil
}

func init() {
	registry.RegisterSource("JSON", JSONSource{})
	registry.RegisterDestination("JSON", JSONDestination{})
}

// ValidateJSONData validates, sanitizes, and unmarshals JSON data
func ValidateJSONData(data string) (interface{}, error) {
	var jsonData interface{}
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		return nil, errors.New("invalid JSON format")
	}

	// Sanitize the JSON data
	sanitizedData := sanitizeJSONData(jsonData)
	logger.Infof("Validation and sanitization successful for JSON data")
	return sanitizedData, nil
}

// sanitizeJSONData recursively sanitizes the JSON data to ensure consistency
func sanitizeJSONData(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		// Recursively sanitize each key-value pair in the map
		for key, value := range v {
			v[key] = sanitizeJSONData(value)
		}
		return v
	case []interface{}:
		// Recursively sanitize each element in the array
		for i, value := range v {
			v[i] = sanitizeJSONData(value)
		}
		return v
	case string:
		// Optionally trim strings or apply further sanitization
		return v
	case float64, bool, nil:
		// Leave primitive types as-is
		return v
	default:
		// Convert unsupported types to their string representations
		logger.Warnf("Unsupported data type %T sanitized to string: %v", v, v)
		return reflect.TypeOf(v).String()
	}
}

// writeJSONFile writes the provided data to a JSON file with proper formatting
func writeJSONFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

// transformJSONData applies transformations to the JSON data
func transformJSONData(data interface{}) (interface{}, error) {
	// Example transformation: Add a key-value pair if the data is a map
	if jsonMap, ok := data.(map[string]interface{}); ok {
		jsonMap["transformed"] = true
		return jsonMap, nil
	}

	// If no transformation is required, return data as is
	logger.Infof("No transformation applied to JSON data")
	return data, nil
}
