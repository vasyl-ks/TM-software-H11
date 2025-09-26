package internal

import (
	"time"
	"github.com/vasyl-ks/TM-software-H11/config"
	"github.com/vasyl-ks/TM-software-H11/model"
)

// calculateAverage returns average temperature and pressure from a slice of SensorData.
func calculateAverage(data []model.SensorData) model.Result {
	var sumTemp, sumPressure float32
	n := float32(len(data))

	for _, d := range data {
		sumTemp += d.Temperature
		sumPressure += d.Pressure
	}

	return model.Result{
		AverageTemp:    sumTemp / n,
		AveragePressure: sumPressure / n,
	}
}

// calculateMin returns minimum temperature and pressure from a slice of SensorData.
func calculateMin(data []model.SensorData) model.Result {
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

	return model.Result{
		MinTemp:        minTemp,
		MinPressure:    minPressure,
	}
}

// calculateMax returns maximum temperature and pressure from a slice of SensorData.
func calculateMax(data []model.SensorData) model.Result {
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

	return model.Result{
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
func Processor(in <-chan model.SensorData, out chan<- model.Result) {
	batchInterval := config.Processor.Interval // defines how often results are calculated.

	var dataSlice []model.SensorData
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	for {
		select {
			case data := <- in:
				dataSlice = append(dataSlice, data)
			case <- ticker.C:
				// Channels for calculations
				avgChan := make(chan model.Result)
				minChan := make(chan model.Result)
				maxChan := make(chan model.Result)

				// Goroutines for calculations
				go func() { avgChan <- calculateAverage(dataSlice) }()
				go func() { minChan <- calculateMin(dataSlice) }()
				go func() { maxChan <- calculateMax(dataSlice) }()

				// Wait for results
				avg := <-avgChan
				min := <-minChan
				max := <-maxChan

				// Build Result
				result := model.Result{
					AverageTemp:    avg.AverageTemp,
					MinTemp:        min.MinTemp,
					MaxTemp:        max.MaxTemp,
					AveragePressure: avg.AveragePressure,
					MinPressure:     min.MinPressure,
					MaxPressure:     max.MaxPressure,
				}
				out <- result

				// Reset slice for next batch
				dataSlice = []model.SensorData{}
		}
	}
}