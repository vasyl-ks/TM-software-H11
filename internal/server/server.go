package server

import (
	"github.com/vasyl-ks/TM-software-H11/internal/model"
)

/*
Server initializes the dataChan and resultChan channels, and calls the Sensor, Process and Log goroutines.
- Sensor runs independently, generates random values, SensorData, and sends it through dataChan.
- Process receives SensorData, calculates statistics, builds a Result, and sends it through resultChan.
- Send receives ResultData from resultChan, marshals it to JSON and sends it via UDP.
*/
func Server() {
	// Create unbuffered channels.
	dataChan := make(chan model.SensorData)
	resultChan := make(chan model.ResultData)

	// Launch concurrent goroutines.
	go Sensor(dataChan)
	go Process(dataChan, resultChan)
	go Send(resultChan)
}
