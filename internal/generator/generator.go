package generator

import (
	"github.com/vasyl-ks/TM-software-H11/internal/model"
)

/*
Generator initializes the dataChan and resultChan channels, then calls the Sensor and Process goroutines.
- Sensor runs independently, generates random values, SensorData, and sends it through dataChan.
- Process receives SensorData, calculates statistics, builds a Result, and sends it through resultChan.
*/
func Run(resultChan chan<- model.ResultData) {
	// Create unbuffered channels.
	dataChan := make(chan model.SensorData)

	// Launch concurrent goroutines.
	go Sensor(dataChan)
	go Process(dataChan, resultChan)
}
