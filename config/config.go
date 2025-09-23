package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Server-side config
type server struct {
	ClientPort 	int	`json:"clientPort"`
	I 			int	`json:"intervalMiliSeconds"`
	Interval 	time.Duration
	VehicleID	string	`json:"vehicleID"`
	SpeedMax	float32	`json:"speedMax"`
	SpeedMin	float32	`json:"speedMin"`
	RPMMax		float32	`json:"rpmMax"`
	RPMMin		float32	`json:"rpmMix"`
	TempMax		float32	`json:"tempMax"`
	TempMin		float32	`json:"tempMin"`
	PressureMax	float32	`json:"pressureMax"`
	PressureMin	float32	`json:"pressureMin"`
}

// Client-side config
type client struct {
	FileDir 	string `json:"fileDir"`
}

// Global config instances
var Server server
var Client client

// LoadConfig reads config.json and configures Server and Client
func LoadConfig() {
	// Open the file
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
		return
	}
	defer file.Close()

	// Read the file
	decoder := json.NewDecoder(file)

	// Parse it into the temp struct
	temp := struct {
		S server `json:"server"`
		C client `json:"client"`
	}{}
	err = decoder.Decode(&temp)
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
		return
	}

	// Copy parsed values into globals
	Server = temp.S
	Client = temp.C

	// Derive time.Duration from milliseconds
	Server.Interval = time.Duration(Server.I) * time.Millisecond
}