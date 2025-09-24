# TM Software H11: Task 1

This Go project models a lightweight sensor analytics pipeline: a sensor goroutine synthesizes pressure/temperature samples, a processor goroutine batches and summarises them, and a logger goroutine persists each summary with simple file rotation logic.

## Table of Contents
- [Features](#features)
- [Repository Layout](#repository-layout)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Telemetry Flow](#telemetry-flow)
- [Development Notes](#development-notes)

## Features
- Sensor that generates random pressure and temperature readings at a configurable interval.
- Processor that buffers readings, computes average/min/max statistics concurrently, and emits structured `Result` values.
- Logger that writes each result to timestamped log files, enforcing a per-file line limit and creating new files as needed.
- Centralised configuration loader that maps `config.json` values into typed packages, exposing ready-to-use `time.Duration` intervals.

## Repository Layout
```
TM-software-H11
│   .gitignore
│   config.json
│   go.mod
│   main.go
│   README.md
├───config
│       config.go
└───modules
        logger.go
        processor.go
        sensor.go
```

## Getting Started
1. Install Go 1.21 or newer.
2. Review `config.json` to tune sensor ranges, batch cadence, and max log lines.
3. Run the application from the repository root: 
```bash
go run main.go
```
4. Observe console output for generated samples. Telemetry logs are written under the directory specified by `client.fileDir` (default: `logs`).

## Configuration
`config.json` is read once at startup by `config.LoadConfig()` and drives runtime behaviour:
- `sensor.intervalSeconds`, `sensor.minPressure`, `sensor.maxPressure`, `sensor.minTemp`, `sensor.maxTemp` define the random value ranges and sampling cadence.
- `processor.intervalSeconds` controls how long the processor buffers readings before producing a batch summary.
- `logger.maxLines` and `logger.fileDir` determine when log rotation occurs and where summary files are stored.

Configuration is loaded at startup by `config.LoadConfig()`; updates require restarting the program.

## Telemetry Flow
1. `modules.Sensor` ticks based on `config.Sensor.Interval`, synthesises a `SensorData`, and sends it to the shared `dataChan`.
2. `modules.Processor` collects readings until the processor interval elapses, fans calculations out to goroutines, then sends a consolidated `Result` into `resultChan`.
3. `modules.Logger` reads each `Result`, appends a formatted line to the active log file, and rotates when the configured maximum line count is reached.

## Development Notes
- `main.go` wires the pipeline by creating channels and starting the three goroutines, using `select {}` to keep the program alive.
- The processor's fan-out/fan-in pattern is educational; a single-pass reducer would be more efficient but less illustrative of concurrency.