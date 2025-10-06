package hub

import (
	"net/http"
	"github.com/vasyl-ks/TM-software-H11/internal/model"
)

/*
Hub acts as a central bridge between the Generator, Frontend, and Consumer.
- Generator ↔ Hub: exchanges ResultData and Command via internal channels.
- Frontend ↔ Hub: exchanges Command and ResultData over WebSocket.
- Consumer ↔ Hub: sends ResultData via UDP and Command via TCP.
*/
func Run(resultChan <-chan model.ResultData) {
	// Create unbuffered channel.
	commandToConsumerChan := make(chan model.Command)

	// WS
	http.HandleFunc("/api/stream", func(w http.ResponseWriter, r *http.Request) {
		// Create Connection
		conn := CreateConnWS(w, r)
		if conn == nil {
			return
		}
		defer conn.Close()

		// Launch concurrent goroutines
		go ReceiveCommandFromFrontEnd(conn, commandToConsumerChan)
		go SendResultToFrontEnd(conn, resultChan)
	})

	// UDP
	{
		// Create Connection
		conn := CreateConnUDP()
		if conn == nil {
			return
		}
		defer conn.Close()

		// Launch concurrent goroutines
		go SendResultToConsumer(conn, resultChan)
	}

	// TCP
	{
		// Create Connection
		conn := CreateConnTCP()
		if conn == nil {
			return
		}
		defer conn.Close()

		// Launch concurrent goroutines
		go SendCommandToConsumer(conn, commandToConsumerChan)
	}
}
