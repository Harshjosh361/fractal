package integrations

import (
	"errors"
	"strings"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
	"github.com/gorilla/websocket"
)

// WebSocketSource struct represents the configuration for consuming messages from WebSocket.
type WebSocketSource struct {
	URL string `json:"websocket_source_url"`
}

// WebSocketDestination struct represents the configuration for publishing messages to WebSocket.
type WebSocketDestination struct {
	URL string `json:"websocket_dest_url"`
}

// FetchData connects to WebSocket, retrieves data, and passes it through validation and transformation pipelines.
func (ws WebSocketSource) FetchData(req interfaces.Request) (interface{}, error) {
	logger.Infof("Connecting to WebSocket Source: URL=%s", req.WebSocketSourceURL)

	if req.WebSocketSourceURL == "" {
		return nil, errors.New("missing WebSocket source details")
	}

	// Connect to WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(req.WebSocketSourceURL, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Read message from WebSocket
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	logger.Infof("Message received from WebSocket: %s", msg)

	// Validation
	validatedData, err := validateWebSocketData(msg)
	if err != nil {
		logger.Fatalf("Validation failed for message: %s, Error: %s", msg, err)
		return nil, err
	}

	// Transformation
	transformedData := transformWebSocketData(validatedData)

	logger.Infof("Message successfully processed and routed: %s", transformedData)
	return transformedData, nil
}

// SendData connects to WebSocket and publishes data to the specified WebSocket server.
func (ws WebSocketDestination) SendData(data interface{}, req interfaces.Request) error {
	logger.Infof("Connecting to WebSocket Destination: URL=%s", req.WebSocketDestURL)

	if req.WebSocketDestURL == "" {
		return errors.New("missing WebSocket destination details")
	}

	// Connect to WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(req.WebSocketDestURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Convert data to string if necessary
	var msg string
	switch v := data.(type) {
	case string:
		msg = v
	case []byte:
		msg = string(v)
	default:
		return errors.New("data should be a string or byte slice to send over WebSocket")
	}

	// Send the message to WebSocket
	err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return err
	}

	logger.Infof("Message sent to WebSocket server: %s", msg)
	return nil
}

// Initialize the WebSocket integrations by registering them with the registry.
func init() {
	registry.RegisterSource("WebSocket", WebSocketSource{})
	registry.RegisterDestination("WebSocket", WebSocketDestination{})
}

// validateWebSocketData ensures the input data meets the required criteria.
func validateWebSocketData(data []byte) ([]byte, error) {
	logger.Infof("Validating data: %s", data)

	// Example: Check if data is non-empty
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	// Add custom validation logic here
	return data, nil
}

// transformWebSocketData modifies the input data as per business logic.
func transformWebSocketData(data []byte) []byte {
	logger.Infof("Transforming data: %s", data)

	// Example: Convert data to uppercase (modify as needed)
	transformed := []byte(strings.ToUpper(string(data)))
	return transformed
}
