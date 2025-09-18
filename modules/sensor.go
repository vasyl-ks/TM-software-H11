package modules

import (
	"fmt"
	"math/rand"
	"time"
)

/*
SensorData represents a single sensor reading,
containing pressure and temperature values.
*/
type SensorData struct {
	Pressure float32
	Temperature float32
}

// sensorInterval defines how often a new sensor reading is generated.
const sensorInterval = 1 * time.Second

/*
Sensor simulates a sensor by generating random pressure and temperature
readings every sensorInterval and sending them to the provided channel.
*/ 
func Sensor(out chan<- SensorData) {
	rand.Seed(time.Now().UnixNano())

	for {
		data := SensorData{
			Pressure:    rand.Float32() * 10,   // 0-10 bar
			Temperature: rand.Float32() * 50,   // 0-50 °C
		}
		out <- data
		fmt.Println("Sensor generated:", data)
		time.Sleep(sensorInterval)
	}
}