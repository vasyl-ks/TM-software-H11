package modules

import (
	"time"
)

/*
Result represents statistics for a batch of SensorData,
containing average, minimum, and maximum values for both temperature and pressure.
*/
type Result struct {
	AverageTemp    float32
	MinTemp        float32
	MaxTemp        float32
	AveragePressure float32
	MinPressure     float32
	MaxPressure     float32
}

// batchInterval defines how often results are calculated.
const batchInterval = 10 * time.Second

// calculateAverage returns average temperature and pressure from a slice of SensorData.
func calculateAverage(data []SensorData) Result {
	var sumTemp, sumPressure float32
	n := float32(len(data))

	for _, d := range data {
		sumTemp += d.Temperature
		sumPressure += d.Pressure
	}

	return Result{
		AverageTemp:    sumTemp / n,
		AveragePressure: sumPressure / n,
	}
}

// calculateMin returns minimum temperature and pressure from a slice of SensorData.
func calculateMin(data []SensorData) Result {
	minTemp := data[0].Temperature
	minPressure := data[0].Pressure

	for _, d := range data[1:] {
		if d.Temperature < minTemp {
			minTemp = d.Temperature
		}
		if d.Pressure < minPressure {
			minPressure = d.Pressure
		}
	}

	return Result{
		MinTemp:        minTemp,
		MinPressure:    minPressure,
	}
}

// calculateMax returns maximum temperature and pressure from a slice of SensorData.
func calculateMax(data []SensorData) Result {
	maxTemp := data[0].Temperature
	maxPressure := data[0].Pressure

	for _, d := range data[1:] {
		if d.Temperature > maxTemp {
			maxTemp = d.Temperature
		}
		if d.Pressure > maxPressure {
			maxPressure = d.Pressure
		}
	}

	return Result{
		MaxTemp:        maxTemp,
		MaxPressure:    maxPressure,
	}
}

/*
Processor collects SensorData values from the input channel into a slice.
Every batchInterval, it calculates statistics (average, min, max) using
separate goroutines (fan-out/fan-in pattern), builds a Result, and sends it
to the output channel.

Note:
- The slice is cleared after each batch, so results are not cumulative.
- Calculations are split into separate functions/goroutines for concurrency practice,
  even though a single-pass calculation would be faster and use less computational overhead.
*/ 
func Processor(in <-chan SensorData, out chan<- Result) {
	var dataSlice []SensorData
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	for {
		select {
			case data := <- in:
				dataSlice = append(dataSlice, data)
			case <- ticker.C:
				// Channels for calculations
				avgChan := make(chan Result)
				minChan := make(chan Result)
				maxChan := make(chan Result)

				// Goroutines for calculations
				go func() { avgChan <- calculateAverage(dataSlice) }()
				go func() { minChan <- calculateMin(dataSlice) }()
				go func() { maxChan <- calculateMax(dataSlice) }()

				// Wait for results
				avg := <-avgChan
				min := <-minChan
				max := <-maxChan

				// Build Result
				result := Result{
					AverageTemp:    avg.AverageTemp,
					MinTemp:        min.MinTemp,
					MaxTemp:        max.MaxTemp,
					AveragePressure: avg.AveragePressure,
					MinPressure:     min.MinPressure,
					MaxPressure:     max.MaxPressure,
				}
				out <- result

				// Reset slice for next batch
				dataSlice = []SensorData{}
		}
	}
}