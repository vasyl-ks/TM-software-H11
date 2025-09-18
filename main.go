package main

import(
	"github.com/vasyl-ks/TM-software-H11/modules"
)

/*
main initializes the dataChan and resultChan channels, and calls the Sensor, Processos and Logger goroutines.
- Sensor runs independently and produces random values for SensorData.
- Processor consumes SensorData from dataChan and produces Result.
- Logger consumes Results from resultChan and logs them in a file.
The final "select {}" keep the program running indefinitely.
*/

func main() {
	// Create unbuffered channels.
	dataChan := make(chan modules.SensorData)
	resultChan := make(chan modules.Result)

	// Launch concurrent goroutines.
	go modules.Sensor(dataChan)
	go modules.Processor(dataChan, resultChan)
	go modules.Logger(resultChan)

	// Block forever to keep the main goroutine alive.
	select {}
}