package internal

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vasyl-ks/TM-software-H11/config"
	"github.com/vasyl-ks/TM-software-H11/model"
)

/*
Sensor simulates a sensor by generating random speed, pressure and temperature
readings every sensorInterval and sending them to the provided channel.
*/
func Sensor(out chan<- model.SensorData) {
	sensorInterval := config.Sensor.Interval // defines how often a new sensor reading is generated.

	ticker := time.NewTicker(sensorInterval)
	defer ticker.Stop()

	minS, maxS := config.Sensor.MinSpeed, config.Sensor.MaxSpeed
	minP, maxP := config.Sensor.MinPressure, config.Sensor.MaxPressure
	minT, maxT := config.Sensor.MinTemp, config.Sensor.MaxTemp

	for range ticker.C {
		data := model.SensorData{
			VehicleID:   config.Vehicle.VehicleID,
			Speed:    	 rand.Float32()*(maxS-minS) + minS,	// 0-150 km/h
			Pressure:    rand.Float32()*(maxP-minP) + minP, // 0-10 bar
			Temperature: rand.Float32()*(maxT-minT) + minT, // 0-50 °C
			CreatedAt: time.Now().Local(),
		}
		out <- data
		fmt.Println("Sensor generated:", data)
	}
}
