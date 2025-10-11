package main

import (
	"github.com/vasyl-ks/TM-software-H11/config"
	consumer "github.com/vasyl-ks/TM-software-H11/internal/consumer"
	generator "github.com/vasyl-ks/TM-software-H11/internal/generator"
	hub "github.com/vasyl-ks/TM-software-H11/internal/hub"
	modelPkg "github.com/vasyl-ks/TM-software-H11/internal/model"
)

/*
Start loads configuration values, creates an internal channel, and then calls the internal goroutines.
- Generator produces SensorData, process it into ResultData and then sends it through internalChan.
- Hub receives ResultData from internalChan, and sends it via UDP to TelemetryLogger.
- Consumer listens for raw JSON datagrams, parses them to ResultData and logs them.
The final "select {}" keep the program running indefinitely.
*/
func Start() {
	// Load configuration (const variables)
	config.LoadConfig()

	// Wait for config to finish
	<-config.Done

	// Creates internal channel of ResultData and Command between Generator and Hub.
	resultChan := make(chan modelPkg.ResultData)
	commandChan := make(chan modelPkg.Command)

	// Run Generator, Hub and Consumer.
	go generator.Run(commandChan, resultChan)
	go consumer.Run()
	<-consumer.Ready // Wait for consumer to initialize UDP&TCP listeners, before Hub tries to connect.
	go hub.Run(resultChan, commandChan)

	select {}
}

func main() {
	Start()
}