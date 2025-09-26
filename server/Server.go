package server

import (
	"github.com/vasyl-ks/TM-software-H11/model"
	"github.com/vasyl-ks/TM-software-H11/server/internal"
)

/*
Server initializes the dataChan and resultChan channels, and calls the Sensor, Processos and Logger goroutines.
- Sensor runs independently and produces random values for SensorData.
- Processor consumes SensorData from dataChan and produces ResultData.
- Sender consumes ResultData from resultChan, marshals it to JSON and sends it via UDP.
*/
func Server() {
	// Create unbuffered channels.
	dataChan := make(chan model.SensorData)
	resultChan := make(chan model.ResultData)

	// Launch concurrent goroutines.
	go internal.Sensor(dataChan)
	go internal.Processor(dataChan, resultChan)
	//go internal.Logger(resultChan)
	go internal.Sender(resultChan)
}
