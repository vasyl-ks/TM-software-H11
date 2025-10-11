package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/vasyl-ks/TM-software-H11/config"
	"github.com/vasyl-ks/TM-software-H11/internal/model"
)

// CreateConnUDP establishes a UDP connection to the configured address and port.
func CreateConnUDP() *net.UDPConn {
	// Client address
	address := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: config.Hub.UDPPort,
	}
	conn, err := net.DialUDP("udp", nil, &address)
	if err != nil {
		log.Println("[ERROR][Hub][UDP] Error connecting via UDP", err)
		panic(err)
	}
	log.Printf("[INFO][Hub][UDP] Established UDP connection from Hub to Consumer, on %s", fmt.Sprintf("%s:%d", address.IP, address.Port))
	return conn
}

/*
SendResultToConsumer receives ResultData from a channel,
marshals it to JSON-encoded []byte
and sends it via UDP to a localhost client.
*/
func SendResultToConsumer(conn *net.UDPConn, inChan <-chan model.ResultData) {
	defer conn.Close()
	
	// Receive ResultData from channel
	for resultData := range inChan {
		// Marshal ResultData to JSON-encoded []byte
		data, err := json.Marshal(resultData)
		if err != nil {
			log.Println("[ERROR][Hub][UDP] Error marshalling WS result JSON", err)
			continue
		}

		// Send JSON via UDP
		_, err = conn.Write(data)
		if err != nil {
			log.Println("[ERROR][Hub][UDP] Error sending via UDP:", err)
			continue
		}
	}
}
