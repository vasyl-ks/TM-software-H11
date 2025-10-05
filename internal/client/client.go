package client

import (
	"github.com/vasyl-ks/TM-software-H11/internal/model"
)

/*
Client initializes the byteChan and jsonChan channels, and calls the Listen, Parse and Log goroutines.
- Listen runs independently, listens for UDP datagrams and sends it through byteChan.
- Parse receives a JSON from byteChan, parses it to ResultData and sends it through resultChan. 
- Log  receives a ResultData from resultChan and logs it.
*/
func Client() {
	// Create unbuffered channels.
	byteChan := make(chan []byte)
	resultChan := make(chan model.ResultData)

	// Launch concurrent goroutines.
	go Listen(byteChan)
	go Parse(byteChan, resultChan)
	go Log(resultChan)
}
