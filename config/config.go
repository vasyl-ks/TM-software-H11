package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type sensor struct {
	Interval 	time.Duration
	I 			int		`json:"intervalSeconds"`
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

var Sensor sensor
var Processor processor
var Logger logger

func LoadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
		return
	}
	defer file.Close()

	temp := struct {
		S	sensor    `json:"sensor"`
		P	processor `json:"processor"`
		L   logger    `json:"logger"`
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