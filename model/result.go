package model

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