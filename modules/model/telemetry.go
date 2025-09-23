package model

import "time"

/*
Telemetry represents a single vehicle reading,
containing speed, RPM, temperature, and pressure
*/
type Telemetry struct {
	VehicleID   string  	`json:"vehicle_id"`
	Speed       float32 	`json:"speed"`
	RPM         float32		`json:"rpm"`
	Temperature float32 	`json:"temperature"`
	Pressure    float32 	`json:"pressure"`
	CreatedAt	time.Time	`json:"created_at"`
}