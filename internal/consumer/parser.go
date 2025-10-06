package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/vasyl-ks/TM-software-H11/internal/model"
)

/*
Parse consumes raw JSON datagrams from the input channel,
attempts to decode each into either a ResultData or a Command object,
then send parsed messages to their respective output channels.
*/
func Parse(in <-chan []byte, outResult chan<- model.ResultData, outCommand chan<- model.Command) {
	for payload := range in {
		// First, try to unmarshal as ResultData
		var res model.ResultData
		if err := json.Unmarshal(payload, &res); err == nil && res.VehicleID != "" {
			outResult <- res
			continue
		}

		// Otherwise, try to unmarshal as Command
		var cmd model.Command
		if err := json.Unmarshal(payload, &cmd); err == nil && cmd.Action != "" {
			outCommand <- cmd
			continue
		}

		// If neither works, log error
		fmt.Printf("Parse: unrecognized JSON payload: %s\n", string(payload))
	}
}