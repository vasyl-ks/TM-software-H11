package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"
)

/*
Telemetry represents a single vehicle reading, 
containing speed, RPM, temperature, and pressure
*/
type Telemetry struct {
	VehicleID   string  `json:"vehicle_id"`
	Speed       float32 `json:"speed"`
	RPM         float32 `json:"rpm"`
	Temperature float32 `json:"temperature"`
	Pressure    float32 `json:"pressure"`
}

// telemetryInterval defines how often Telemetry is generated and sent.
const telemetryInterval = 100 * time.Millisecond

// clientPort defines the UDP port to which telemetry messages are sent.
const clientPort = 100000

/*
UDPServer simulates a vehicle telemetry server.
It generates random Telemetry data at a fixed interval, 
marshals it to JSON and sends it via UDP to localhost client address 
*/
func UDPServer() {
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
		t := Telemetry{
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