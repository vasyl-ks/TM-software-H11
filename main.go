package main
import (
	"github.com/vasyl-ks/TM-software-H11/config"
	"github.com/vasyl-ks/TM-software-H11/server"
)

/*
main loads configuration values, and then calls the Server goroutine.
- Server generates SensorData, proceess it into Result and then logs it.
The final "select {}" keep the program running indefinitely.
*/
func main() {
	// Load configuration (const variables)
	config.LoadConfig()

	go server.Server()

	select {}
}