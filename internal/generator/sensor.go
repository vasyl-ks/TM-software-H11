package generator

import (
	"math/rand"
	"strings"
	"time"

	"github.com/vasyl-ks/TM-software-H11/config"
	"github.com/vasyl-ks/TM-software-H11/internal/model"
)

/*
Sensor simulates a sensor by generating random speed, pressure and temperature
readings every sensorInterval and sending them to the provided channel.

Speed behavior responds to control commands:
- "Start" → enables movement.
- "Stop" → sets speed to 0.
- "Accelerate n" → increases current speed by n.
- "Mode" → changes driving mode (eco|normal|speed).
*/
func Sensor(inCommandChan <-chan model.Command, outChan chan<- model.SensorData) {
	sensorInterval := config.Sensor.Interval // defines how often a new sensor reading is generated.

	ticker := time.NewTicker(sensorInterval)
	defer ticker.Stop()

	minS, maxS := config.Sensor.MinSpeed, config.Sensor.MaxSpeed
	minP, maxP := config.Sensor.MinPressure, config.Sensor.MaxPressure
	minT, maxT := config.Sensor.MinTemp, config.Sensor.MaxTemp

	// vehicle state
	currentSpeed := float32(0)
	mode := "normal"
	started := false

	for {
		select {
			case cmd := <- inCommandChan:
				switch strings.ToLower(cmd.Action) {
					case "start":
						started = true
					case "stop":
						started = false
						currentSpeed = 0
					case "accelerate":
						// Try to read numeric parameter
						if val, ok := cmd.Params.(float64); ok {
							currentSpeed += float32(val)
						}
					case "mode":
						if val, ok := cmd.Params.(string); ok {
							mode = strings.ToLower(val)
						}
				}
			
			case <- ticker.C:
				// adjust max speed depending on mode
				var maxAllowed float32
				switch mode {
					case "eco":
						maxAllowed = maxS * 0.5
					case "normal":
						maxAllowed = maxS * 0.8
					case "speed":
						maxAllowed = maxS
					default:
						maxAllowed = maxS * 0.8
				}

				// adjust growth factors based on mode
				var growthFactor float32
				switch mode {
				case "eco":
					growthFactor = 0.7   // slowest growth
				case "normal":
					growthFactor = 1.0   // standard growth
				case "speed":
					growthFactor = 1.3   // fastest growth
				default:
					growthFactor = 1.0
				}

				// simulate speed
				if started {
					if currentSpeed < minS {
						currentSpeed = minS
					}
					if currentSpeed > maxAllowed {
						currentSpeed = maxAllowed
					}
				} else {
					currentSpeed = 0
				}

				// Normalize the current speed into a [0,1] range
				// 0 means minimum speed, 1 means maximum speed
				speedRatio := (currentSpeed - minS) / (maxS - minS)
				if speedRatio < 0 {
					speedRatio = 0
				} else if speedRatio > 1 {
					speedRatio = 1
				}

				// Make pressure and temperature increase with speed
				// Both grow linearly from their minimum to maximum values
				// They grow faster on the sport mode and slower on eco.
				pressure := minP + (speedRatio * growthFactor) * (maxP - minP)
				temperature := minT + (speedRatio * growthFactor) * (maxT - minT)

				// Add a small random noise to simulate sensor variability
				pressure += rand.Float32()*0.1 - 0.05
				temperature += rand.Float32()*0.5 - 0.25

				// Adjust pressure to avoid negative values
				if pressure < 0 {
					pressure = 0
				}

				// Adjust pressure to avoid values above maximum.
				if pressure > maxP {
					pressure = maxP
				}
				if temperature > maxT {
					temperature = maxT
				}


				// Create SensorData
				data := model.SensorData{
					VehicleID:   config.Vehicle.VehicleID,
					Speed:    	 currentSpeed,
					Pressure:    pressure,
					Temperature: temperature,
					CreatedAt: time.Now().Local(),
				}
				outChan <- data
		}
	} 	
}