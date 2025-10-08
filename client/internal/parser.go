package internal

import (
	"encoding/json"
	"fmt"
	"github.com/vasyl-ks/TM-software-H11/model"
)

/*
Parse consumes a raw JSON datagram from in, unmarshals it to ResultData, and sends it to out.
*/
func Parse(in <-chan []byte, out chan<- model.ResultData) {
	for payload := range in { // Receive Datagram
		// Decode JSON
		var entry model.ResultData
		if err := json.Unmarshal(payload, &entry); err != nil {
			fmt.Println("Error decoding datagram:", err)
			continue
		}

		// Send the JSON through the channel
		out <- entry
	}
}