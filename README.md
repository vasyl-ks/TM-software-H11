# TM Software H11: Task 2

This project simulates a simple telemetry pipeline using UDP sockets. A server produces random vehicle readings, sends them over UDP, and a client listens, decodes the JSON payloads, and writes them to timestamped log files.

## Table of Contents
- [Features](#features)
- [Repository Layout](#repository-layout)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Telemetry Flow](#telemetry-flow)
- [Development Notes](#development-notes)

## Features
- UDP server that generates vehicle telemetry (speed, RPM, temperature, pressure) with configurable ranges.
- UDP client that receives telemetry, parses the JSON payload, and logs the data with creation and receipt timestamps.
- Shared `modules/model` package defining the telemetry schema consumed by both server and client.
- Centralised configuration via `config.json` for ports, send intervals, telemetry range limits, and client log directory.

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
└───modules
    ├───client
    │       UDPClient.go
    ├───model
    │       telemetry.go
    └───server
            UDPServer.go
```

## Getting Started
1. Install Go 1.21 or newer.
2. Review and adjust `config.json` to match your port, interval, and telemetry range requirements.
3. Run the application from the repository root:
   ```bash
   go run main.go
   ```
4. Observe console output for server/client startup messages. Telemetry logs are written under the directory specified by `client.fileDir` (default: `logs`).

## Configuration
`config.json` controls runtime behaviour:
- `server.clientPort`: UDP port used by both server and client.
- `server.intervalMiliSeconds`: delay between telemetry messages.
- `server.vehicleID`: identifier embedded in each telemetry record.
- `server.*Min/*Max`: numeric ranges for randomly generated readings.
- `client.fileDir`: directory where the client writes log files.

Configuration is loaded at startup by `config.LoadConfig()`; updates require restarting the program.

## Telemetry Flow
1. The server builds a `model.Telemetry` instance, stamps it with `CreatedAt`, and marshals it to JSON.
2. The server sends the JSON datagram to the configured UDP port.
3. The client reads the datagram, unmarshals it into `model.Telemetry`, and writes a formatted log entry:
   - Log entries include both the `CreatedAt` timestamp and the local receipt time with microsecond precision.

## Development Notes
- The server and client run in separate goroutines launched from `main.go`; the main goroutine blocks with `select {}`.
- Telemetry values rely on `math/rand`; seeding or alternative distributions can be added in `modules/server/UDPServer.go`.
- Logs are appended via the standard `log` package with formatting customisations in `modules/client/UDPClient.go`.
