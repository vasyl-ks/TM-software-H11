package main

import (
	"github.com/vasyl-ks/TM-software-H11/config"
	generator "github.com/vasyl-ks/TM-software-H11/internal/generator"
	hub		  "github.com/vasyl-ks/TM-software-H11/internal/hub"
	consumer  "github.com/vasyl-ks/TM-software-H11/internal/consumer"
	modelPkg  "github.com/vasyl-ks/TM-software-H11/internal/model"
)

/*
main loads configuration values, creates an internal channel, and then calls the internal goroutines.
- Generator produces SensorData, process it into ResultData and then sends it through internalChan.
- Hub receives ResultData from internalChan, and sends it via UDP to TelemetryLogger.
- Consumer listens for raw JSON datagrams, parses them to ResultData and logs them.
The final "select {}" keep the program running indefinitely.
*/
func main() {
	// Load configuration (const variables)
	config.LoadConfig()

	// Creates internal channel of ResultData and Command between Generator and Hub.
	resultChan := make(chan modelPkg.ResultData)
	commandChan := make(chan modelPkg.Command)

	// Run Generator, Hub and Consumer.
	go generator.Run(commandChan, resultChan)
	go hub.Run(resultChan, commandChan)
	go consumer.Run()

	select {}
}