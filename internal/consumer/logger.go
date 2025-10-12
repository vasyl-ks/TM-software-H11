package consumer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/vasyl-ks/TM-software-H11/config"
	"github.com/vasyl-ks/TM-software-H11/internal/model"
)

// Helper type for grouped loggers
type Loggers struct {
	Main    *log.Logger
	Data    *log.Logger
	Command *log.Logger
}

// Helper function to create a logger for a given subdirectory and prefix
func createLogger(baseDir, subDir, prefix string) (*log.Logger, *os.File, error) {
	dirPath := filepath.Join(baseDir, subDir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, nil, fmt.Errorf("[ERROR][Consumer][Log] Error creating directory %s: %w", dirPath, err)
	}

	filename := fmt.Sprintf("%s_%s.jsonl", prefix, time.Now().Format("20060102_150405"))
	fullPath := filepath.Join(dirPath, filename)

	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("[ERROR][Consumer][Log] Error opening log file %s: %w", fullPath, err)
	}

	logger := log.New(file, "", 0)
	return logger, file, nil
}

// Helper function to write a ResultData
func writeResult(loggers Loggers, r model.ResultData) {
	msg := fmt.Sprintf(
		"[DATA] Created at %s, Processed at %s, Logged at %s | "+
			"AvgSpeed: %5.2f, MinSpeed: %5.2f, MaxSpeed: %5.2f | "+
			"AvgTemp: %5.2f, MinTemp: %5.2f, MaxTemp: %5.2f | "+
			"AvgPressure: %4.2f, MinPressure: %4.2f, MaxPressure: %4.2f",
		r.CreatedAt.Format("15:04:05.000000"),
		r.ProcessedAt.Format("15:04:05.000000"),
		time.Now().Local().Format("15:04:05.000000"),
		r.AverageSpeed, r.MinimumSpeed, r.MaximumSpeed,
		r.AverageTemperature, r.MinimumTemperature, r.MaximumTemperature,
		r.AveragePressure, r.MinimumPressure, r.MaximumPressure,
	)

	loggers.Main.Println(msg)
	loggers.Data.Println(msg)
}

// Helper function to write a Command
func writeCommand(loggers Loggers, cmd model.Command) {
	msg := fmt.Sprintf(
		"[COMMAND] Received at %s | Action: %-12s | Params: %-8v",
		time.Now().Local().Format("15:04:05.000000"),
		cmd.Action,
		cmd.Params,
	)

	loggers.Main.Println(msg)
	loggers.Command.Println(msg)
}

/*
Log receives ResultData and Command messages from their respective channels
and logs them to rotating log files.
- Each file contains up to maxLines entries.
- Once the limit is reached, the current file is closed and a new file is created.
- Files are named using the creation timestamp in the format "YYYYMMDD_hhmmss".
- If terminated early, the current file may have fewer than maxLines; a new file is created on the next run.
*/
func Log(inResultChan <-chan model.ResultData, inCommandChan <-chan model.Command) {
	lineCount := 0
	fileDir := config.Logger.FileDir   // defines directory where the log is saved.
	maxLines := config.Logger.MaxLines // defines the maximum number of ResultData to log in a single file.

	// Create base directory
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		log.Println("[ERROR][Consumer][Log] Error creating directory:", err)
		return
	}

	// Create all loggers
	mainLogger, mainFile, err := createLogger(fileDir, "", "log")
	if err != nil {
		log.Println(err)
		return
	}
	defer mainFile.Close()

	dataLogger, dataFile, err := createLogger(fileDir, "data", "data")
	if err != nil {
		log.Println(err)
		return
	}
	defer dataFile.Close()

	commandLogger, commandFile, err := createLogger(fileDir, "commands", "command")
	if err != nil {
		log.Println(err)
		return
	}
	defer commandFile.Close()

	// Group them for easier access
	loggers := Loggers{
		Main:    mainLogger,
		Data:    dataLogger,
		Command: commandLogger,
	}

	log.Println("[INFO][Consumer][Log] Running.")

	for {
		select {
		// Receive ResultData
		case resultData, ok := <-inResultChan:
			if !ok {
				inResultChan = nil // channel closed
				continue
			}
			// Log in the file
			writeResult(loggers, resultData)

		// Receive Command
		case cmd, ok := <-inCommandChan:
			if !ok {
				inCommandChan = nil // channel closed
				continue
			}
			// Log in the file
			writeCommand(loggers, cmd)
		}

		// Exit if both channels are closed
		if inResultChan == nil && inCommandChan == nil {
			break
		}

		lineCount++
		if lineCount >= maxLines {
			mainFile.Close()
			dataFile.Close()
			commandFile.Close()

			// Create a new file
			mainLogger, mainFile, err = createLogger(fileDir, "", "log")
			if err != nil {
				fmt.Println(err)
				return
			}

			dataLogger, dataFile, err = createLogger(fileDir, "data", "data")
			if err != nil {
				fmt.Println(err)
				return
			}

			commandLogger, commandFile, err = createLogger(fileDir, "commands", "command")
			if err != nil {
				fmt.Println(err)
				return
			}

			// Group them for easier access
			loggers = Loggers{
				Main:    mainLogger,
				Data:    dataLogger,
				Command: commandLogger,
			}

			lineCount = 0
		}
	}
}
