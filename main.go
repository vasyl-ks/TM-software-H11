package main

import (
	"github.com/vasyl-ks/TM-software-H11/config"
	"github.com/vasyl-ks/TM-software-H11/server"
	"github.com/vasyl-ks/TM-software-H11/client"
)

/*
main loads configuration values, and then calls the Server goroutine.
- Server generates SensorData, process it into Result and then sends it via UDP.
- Client listens for raw JSON datagrams, parses them to ResultData and logs them.
The final "select {}" keep the program running indefinitely.
*/
func main() {
	// Load configuration (const variables)
	config.LoadConfig()

	// Run Server and Client.
	go server.Server()
	go client.Client()

	select {}
}