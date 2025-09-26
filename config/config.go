package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type vehicle struct {
	VehicleID	string `json:"vehicleID"`
}

type sensor struct {
	Interval 	time.Duration
	I 			int		`json:"intervalSeconds"`
	MaxSpeed	float32 `json:"maxSpeed"`
	MinSpeed	float32	`json:"minSpeed"`
	MaxPressure float32	`json:"maxPressure"`
	MinPressure float32	`json:"minPressure"`
	MaxTemp     float32 `json:"maxTemp"`
	MinTemp     float32 `json:"minTemp"`
}

type processor struct {
	Interval 	time.Duration
	I	int		`json:"intervalSeconds"`
}

type logger struct {
	MaxLines 	int		`json:"maxLines"`
	FileDir 	string	`json:"fileDir"`
}

var Vehicle 	vehicle
var Sensor 		sensor
var Processor 	processor
var Logger 		logger

func LoadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
		return
	}
	defer file.Close()

	temp := struct {
		V	vehicle		`json:"vehicle"`
		S	sensor    	`json:"sensor"`
		P	processor 	`json:"processor"`
		L   logger    	`json:"logger"`
	}{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&temp)
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
		return
	}

	Sensor = temp.S
	Processor = temp.P
	Logger = temp.L

	Sensor.Interval = time.Duration(Sensor.I) * time.Second
	Processor.Interval = time.Duration(Processor.I) * time.Second
}