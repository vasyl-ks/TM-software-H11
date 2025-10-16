package model

import "time"

/*
ResultData represents statistics for a batch of SensorData,
containing average, minimum, and maximum values for both speed, temperature and pressure
and indemnifications such as its ID and the time it was generated and processed.
*/
type ResultData struct {
	AverageSpeed       float32
	MinimumSpeed       float32
	MaximumSpeed       float32
	AverageTemperature float32
	MinimumTemperature float32
	MaximumTemperature float32
	AveragePressure    float32
	MinimumPressure    float32
	MaximumPressure    float32
	VehicleID          string
	CreatedAt          time.Time
	ProcessedAt        time.Time
}
