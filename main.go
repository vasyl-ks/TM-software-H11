package main

import (
	"github.com/vasyl-ks/TM-software-H11/config"
	"github.com/vasyl-ks/TM-software-H11/modules/server"
)

func main() {
	// Load runtime configuration
	config.LoadConfig()

	// Run Server
	server.UDPServer()
}