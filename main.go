package main

import(
	"github.com/vasyl-ks/TM-software-H11/modules/server"
	"github.com/vasyl-ks/TM-software-H11/config"
)

func main() {
	config.LoadConfig()
	server.UDPServer()
}