package main

import (
	"github.com/vasyl-ks/TM-software-H11/config"
	clientpkg "github.com/vasyl-ks/TM-software-H11/internal/client"
    serverpkg "github.com/vasyl-ks/TM-software-H11/internal/server"
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
	go serverpkg.Server()
	go clientpkg.Client()

	select {}
}