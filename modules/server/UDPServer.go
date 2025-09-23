package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"
	"github.com/vasyl-ks/TM-software-H11/config"
	"github.com/vasyl-ks/TM-software-H11/modules/model"
)



/*
UDPServer simulates a vehicle telemetry server.
It generates random Telemetry data at a fixed interval, 
marshals it to JSON and sends it via UDP to localhost client address 
*/
func UDPServer() {
	telemetryInterval := config.Server.Interval // defines how often Telemetry is generated and sent.
	clientPort := config.Server.ClientPort // defines the UDP port to which telemetry messages are sent.

	// Client address
	clientAddr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: clientPort,
	}

	// Create connection
	conn, err := net.DialUDP("udp", nil, &clientAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		panic(err)
	}
	defer conn.Close()
	fmt.Println("UDPServer sending telemetry to", clientAddr.String())

	// Create ticker
	ticker := time.NewTicker(telemetryInterval)
	defer ticker.Stop()

	for range ticker.C {
		// Create a random telemetry
		t := model.Telemetry{
			VehicleID:   "V123",
			Speed:       rand.Float32() * 200,   // 0–200 km/h
			RPM:         rand.Float32() * 8000,  // 0–8000 rpm
			Temperature: rand.Float32() * 100,   // 0–100 °C
			Pressure:    rand.Float32() * 10,    // 0–10 bar
		}

		// Marshal telemetry to JSON
		data, err := json.Marshal(t)
		if err != nil {
    		fmt.Println("Error marshalling:", err)
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