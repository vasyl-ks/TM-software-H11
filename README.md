# TM Software H11: Task 1

This Go project simulates a basic vehicle telemetry system: the server generates sensor readings, groups them into summaries, and sends those over UDP; the client listens for these summaries, turns them into structured data, and saves them to log files.

## Table of Contents
- [Features](#features)
- [Repository Layout](#repository-layout)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Telemetry Flow](#telemetry-flow)
- [Development Notes](#development-notes)

## Features
- Synthetic sensor that produces speed, pressure, and temperature samples at a configurable cadence.
- Statistics processor that calculates min/avg/max values per batch using fan-out/fan-in goroutines before emitting `ResultData` summaries.
- UDP sender that serializes each summary to JSON and delivers it to a localhost client on the configured port.
- Client-side listener, parser, and rotating logger that persist the received summaries with consistent formatting.
- Central configuration package that exposes strongly typed settings (`time.Duration` intervals, ports, limits) to every component.
- Shared `model` package defining sensor and result payload schemas for server and client interoperability.

## Repository Layout
```
TM-software-H11/
│   .gitignore
│   config.json
│   go.mod
│   main.go
│   README.md
├───client
│   │   client.go
│   └───internal
│           listener.go
│           logger.go
│           parser.go
├───config
│       config.go
├───logs
│       sensor_log_20250928_193501.log
├───model
│       resultData.go
│       sensorData.go
└───server
    │   Server.go
    └───internal
            processor.go
            sender.go
            sensor.go
```

## Getting Started
1. Install Go 1.21 or newer.
2. Clone the repository and enter the project directory:
   ```bash
   git clone https://github.com/vasyl-ks/TM-software-H11.git
   cd TM-software-H11
   ```
3. Inspect `config.json` to set the vehicle identifier, sensor ranges, batch interval, logger settings, and UDP port.
4. Ensure the selected UDP port (default 10000) is available on localhost.
5. From the repository root, run:
   ```bash
   go run main.go
   ```
6. Watch the terminal for server/client status messages; log files are created under the directory specified by `logger.fileDir` (default `logs/`).

## Configuration
`config.json` is loaded at startup by `config.LoadConfig()` and controls runtime behaviour:
- `vehicle.vehicleID` tags every reading and result with the configured identifier.
- `sensor.intervalMilliSeconds`, `sensor.minSpeed`, `sensor.maxSpeed`, `sensor.minPressure`, `sensor.maxPressure`, `sensor.minTemp`, `sensor.maxTemp` define the random value ranges and sampling cadence.
- `processor.intervalMilliSeconds` sets how often batches close and statistics are emitted.
- `logger.maxLines` and `logger.fileDir` govern log rotation and storage location for client output files.
- `senderANDlistener.udpPort` specifies the UDP port used by both the server sender and client listener on localhost.

Configuration is read once at process start; update the file and restart the application to apply changes.

## Telemetry Flow
1. `server/internal.Sensor` wakes on `config.Sensor.Interval`, generates `model.SensorData`, and publishes it to `dataChan`.
2. `server/internal.Process` buffers readings until `config.Processor.Interval` elapses, calculates statistics via concurrent helpers, and produces a `model.ResultData` summary.
3. `server/internal.Send` serializes the summary to JSON and sends the datagram to `127.0.0.1:udpPort`.
4. `client/internal.Listen` binds the same UDP port, clones each datagram, and forwards the raw bytes to `byteChan`.
5. `client/internal.Parse` unmarshals the payload into `model.ResultData`, validating the JSON before passing it along.
6. `client/internal.Log` writes formatted log entries, rotating files after `config.Logger.MaxLines` entries and timestamping every stage (created, processed, logged).

## Development Notes
- `main.go` loads configuration once, launches the server and client goroutines, and blocks forever with `select {}` to keep both halves running in a single binary.
- The processor intentionally uses a fan-out/fan-in pattern for educational concurrency, even though a single-pass reducer would be more efficient.
- UDP transport keeps server and client decoupled; both rely on the shared `model` definitions to guarantee compatible payloads.
- The client logger sets `log` output per rotation to ensure each file receives consistent formatting and timestamps for troubleshooting end-to-end latency.
