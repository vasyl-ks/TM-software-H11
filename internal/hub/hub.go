package hub

import "github.com/vasyl-ks/TM-software-H11/internal/model"

/*
Generator calls the UDPSend goroutine.
- SendUDP receives ResultData from a channel, marshals it to JSON-encoded []byte and sends it via UDP to a localhost client.
*/
func Run(resultChan <-chan model.ResultData) {
	// Launch concurrent goroutines.
	go SendUDP(resultChan)
}
