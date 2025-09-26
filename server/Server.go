package server

import (
	"github.com/vasyl-ks/TM-software-H11/server/internal"
	"github.com/vasyl-ks/TM-software-H11/model"
)

/*
Server initializes the dataChan and resultChan channels, and calls the Sensor, Processos and Logger goroutines.
- Sensor runs independently and produces random values for SensorData.
- Processor consumes SensorData from dataChan and produces Result.
- Logger consumes Results from resultChan and logs them in a file.
*/
func Server() {
	// Create unbuffered channels.
	dataChan := make(chan model.SensorData)
	resultChan := make(chan model.Result)

	// Launch concurrent goroutines.
	go internal.Sensor(dataChan)
	go internal.Processor(dataChan, resultChan)
	go internal.Logger(resultChan)
}