package hub

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/vasyl-ks/TM-software-H11/internal/model"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// CreateConnWS upgrades an HTTP connection to a WebSocket and returns the connection object.
func CreateConnWS(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR][Hub][WS] Error upgrading to WebSocket:", err)
		return nil
	}
	log.Printf("[INFO][Hub][WS] Established WS connection between Hub and Frontend, on %s", conn.LocalAddr())
	return conn
}

/*
ListenCommandWS listens for a command from the WebSocket
parses it to a Go struct
and forwards it to a channel.
*/
func ReceiveCommandFromFrontEnd(conn *websocket.Conn, outChan1 chan<- model.Command, outChan2 chan<- model.Command) {
	defer conn.Close()
	
	for {
		// Listen for WS Command JSON
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived) {
				log.Printf("[INFO][Hub][WS] Client disconnected normally: %v", err) // Expected error
			} else {
				log.Printf("[ERROR][Hub][WS] Error reading WS command: %v", err) // Unexpected error
			}
			break
		}

		// Parse it to Command struct
		var cmd model.Command
		if err := json.Unmarshal(msg, &cmd); err != nil {
			log.Println("[ERROR][Hub][WS] Error parsing WS command JSON:", err)
			continue
		}

		// Sends it to channel
		outChan1 <- cmd
		outChan2 <- cmd
	}
}

/*
SendResultToFrontEnd receives ResultData from a channel, 
marshals it to JSON-encoded []byte 
and sends it via WS to the WebSocket client.
*/
func SendResultToFrontEnd(conn *websocket.Conn, inChan <-chan model.ResultData) {
	defer func() {
		conn.Close()
		log.Printf("[INFO][Hub][WS] Writer closed connection: %s", conn.RemoteAddr())
	}()
	
	// Receive ResultData from channel
	for result := range inChan {
		// Marshal ResultData to JSON-encoded []byte
		data, err := json.Marshal(result)
		if err != nil {
			log.Println("[ERROR][Hub][WS] Error marshalling WS result JSON:", err)
			continue
		}

		// Send JSON via WS
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
        		websocket.CloseGoingAway,
        		websocket.CloseNoStatusReceived) ||
        		strings.Contains(err.Error(), "close sent") {
				log.Printf("[INFO][Hub][WS] Client disconnected during write: %v", err) // Expected error
			} else {
				log.Printf("[ERROR][Hub][WS] Error sending via WS: %v", err) // Unexpected error
			}
			break
		}
	}
}
