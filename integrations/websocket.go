package integrations

import (
	"errors"

	"github.com/SkySingh04/fractal/logger"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
	"github.com/gorilla/websocket"
)

// WebSocketSource implements the DataSource interface
type WebSocketSource struct {
	URL string `json:"url"`
}

// WebSocketDestination implements the DataDestination interface
type WebSocketDestination struct {
	URL string `json:"url"`
}

// FetchData fetches data from a WebSocket server
func (w WebSocketSource) FetchData(req interfaces.Request) (interface{}, error) {
	if err := validateWebSocketRequest(req, true); err != nil {
		return nil, err
	}
	logger.Infof("Fetching data from WebSocket...")
	conn, _, err := websocket.DefaultDialer.Dial(req.WebSocketSourceURL, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	// Receive data from WebSocket
	_, message, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return string(message), nil
}

// SendData sends data to a WebSocket server
func (w WebSocketDestination) SendData(data interface{}, req interfaces.Request) error {
	if err := validateWebSocketRequest(req, false); err != nil {
		return err
	}
	logger.Infof("Sending data to WebSocket...")
	conn, _, err := websocket.DefaultDialer.Dial(req.WebSocketDestURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	// Send data to WebSocket
	err = conn.WriteMessage(websocket.TextMessage, []byte(data.(string)))
	if err != nil {
		return err
	}
	return nil
}

// validateWebSocketRequest validates the request fields for WebSocket
func validateWebSocketRequest(req interfaces.Request, isSource bool) error {
	if isSource && req.WebSocketSourceURL == "" {
		return errors.New("missing WebSocket source URL")
	}
	if !isSource && req.WebSocketDestURL == "" {
		return errors.New("missing WebSocket destination URL")
	}
	return nil
}

func init() {
	registry.RegisterSource("WebSocket", WebSocketSource{})
	registry.RegisterDestination("WebSocket", WebSocketDestination{})
}
