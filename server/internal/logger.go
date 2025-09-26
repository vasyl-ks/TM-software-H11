package internal

import (
	"fmt"
	"os"
	"time"
	"log"
	"github.com/vasyl-ks/TM-software-H11/config"
	"github.com/vasyl-ks/TM-software-H11/model"
)

/*
Logger receives Results from a channel and logs them into file.
- Each file contains up to maxLines entries.
- Once the limit is reached, the current file is closed and a new file is created.
- Files are named using the creation timestamp in the format "YYYYMMDD_hhmmss".
- If terminated early, the current file may have fewer than maxLines; a new file is created on the next run.
*/
func Logger(resultChan <-chan model.Result) {
	lineCount := 0
	fileDir := config.Logger.FileDir // defines directory where the log is saved.
	maxLines := config.Logger.MaxLines // defines the maximum number of Results to log in a single file.

	// Check directory
	err := os.MkdirAll(fileDir, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Create a new file
	filename := fmt.Sprintf("%s/sensor_log_%s.log", fileDir, time.Now().Format("20060102_150405"))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()
	log.SetOutput(file)

	for result := range resultChan {
		// Write in the file
		log.Printf(
    		"|| AvgTemp: %5.2f, MinTemp: %5.2f, MaxTemp: %5.2f | AvgPressure: %4.2f, MinPressure: %4.2f, MaxPressure: %4.2f\n",
			result.AverageTemp, result.MinTemp, result.MaxTemp,
			result.AveragePressure, result.MinPressure, result.MaxPressure,
		)

		lineCount++
		if lineCount >= maxLines {
			file.Close()
			
			// Create a new file
			filename = fmt.Sprintf("logs/sensor_log_%s.txt", time.Now().Format("20060102_150405"))
			file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				fmt.Println("Error opening new log file:", err)
				return
			}
			defer file.Close()
			log.SetOutput(file)

			lineCount = 0
		}
	}
}