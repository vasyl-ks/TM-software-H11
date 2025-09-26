package internal

import (
	"encoding/json"
	"fmt"
	"net"
	"github.com/vasyl-ks/TM-software-H11/config"
	"github.com/vasyl-ks/TM-software-H11/model"
)

/*
Sender receives ResultData from a channel, marshals it to JSON and sends it via UDP to a localhost client.
*/
func Sender(resultChan <-chan model.ResultData) {

	// Client address
	clientAddr := net.UDPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: config.Sender.ClientPort,
	}

	// Create connetion
	conn, err := net.DialUDP("udp", nil, &clientAddr)
	if err != nil {
		fmt.Println("Error connecting", err)
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Server sending data to", clientAddr.String())

	for resultData := range resultChan { 
		// Marshal ResultData to JSON
		data, err := json.Marshal(resultData)
		if err != nil {
			fmt.Println("Error marshalling", err)
			continue
		}

		// Send JSON via UDP
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error sending:", err)
			continue
		}
	}
}