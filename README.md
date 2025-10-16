# TM Software H11: Final Stage

This Go project simulates a complete vehicle telemetry system.
It connects multiple components — a **Generator**, **Hub**, **Consumer**, and a **Frontend (WebSocket)** — that exchange telemetry, control commands in real time and visualizes metrics.

> This project was developed over the course of **four weeks** as part of the **Training Month (TM)** as a **Backend Engineer** for **Hyperloop UPV**.  
> The different *branches* reflect the project’s progress throughout the four weeks.  
> **Backend** components (contained in this repository) were implemented by me, while the **frontend** —included as a Git *submodule*— was developed by **[maximka76667](https://github.com/maximka76667)**.  
> In the final stage, both parts were integrated into a fully functional end-to-end system.

## Table of Contents
* [Features](#features)
* [Repository Layout](#repository-layout)
* [Getting Started](#getting-started)
* [Configuration](#configuration)
* [System Flow](#system-flow)
* [Development Notes](#development-notes)

## Features
* Synthetic **Generator** models speed, pressure, and temperature, reacting to `start`, `stop`, `accelerate`, and `mode` commands with mode-aware speed caps.
* **Hub** that routes data and commands between Generator, Consumer, and Frontend:
  * fans telemetry to the frontend (WebSocket) and consumer (UDP) while forwarding commands from the frontend to the generator (channels) and consumer (TCP).
  * `ResultData` is sent to both the Frontend (WS) and Consumer (UDP).
  * `Command` messages flow from the Frontend (WS) to the Generator and Consumer (TCP).
* **Consumer** listens on UDP/TCP, autodetects `ResultData` vs `Command` payloads, and rotates structured `.jsonl` logs across `logs/`, `logs/data/`, and `logs/commands/`.
* **React frontend** (Vite + Tailwind) offers connect/disconnect controls, command groups, toast feedback, and metric tiles that track the latest batch stats in real time.
* Central **config package** exposes runtime tuning parameters — settings that define how the system behaves when running, such as sensor cadence, aggregation windows, port bindings, log rotation, and vehicle identity.
* End-to-end **integration test** (`cmd/app/main_test.go`) spins up the stack, drives scripted WebSocket commands, and records the telemetry stream under `test/`.
* `start.sh` boots the Go backend and frontend dev server together for a single command developer experience.

## Repository Layout
```
TM-software-H11/
│   .gitignore
│   .gitmodules
│   config.json
│   go.mod
│   go.sum
│   README.md
│   start.sh
│   
├───cmd
│   └───app
│           main.go
│           main_test.go
│
├───config
│       config.go
│
├───frontend
│   │	package.json
│   │	package-lock.json
│   └───src
│
├───internal
│   ├───consumer
│   │       consumer.go
│   │       listener.go
│   │       logger.go
│   │       parser.go
│   │
│   ├───generator
│   │       generator.go
│   │       processor.go
│   │       sensor.go
│   │
│   ├───hub
│   │       hub.go
│   │       tcphandler.go
│   │       updhandler.go
│   │       wshandler.go
│   │
│   └───model
│           command.go
│           resultData.go
│           sensorData.go
│
├───logs
│   ├───commands
│   └───data
└───test
        test_logs.jsonl
```

## Getting Started
1. Install Go 1.25.1 or newer and Node.js 20+ (with npm).
2. Clone the repository and change into it:

   ```bash
   git clone https://github.com/vasyl-ks/TM-software-H11.git
   cd TM-software-H11
   ```
3. Inspect and adjust `config.json` for your desired intervals, ports, range and mode ratios.
4. Install frontend dependencies:

   ```bash
   cd frontend
   npm install
   ```
5. Run the stack:

   ```bash
   cd ..
   ./start.sh
   ```

   The script launches `go run ./cmd/app/main.go` and `npm run dev` (Vite). Stop with `Ctrl+C`.
6. To run only the backend:

   ```bash
   go run ./cmd/app/main.go
   ```

   Then start the frontend separately with `npm run dev` inside `frontend/` (use `-- --host` if you need LAN access).
7. Execute the integration test (writes logs under `test/`):

   ```bash
   go test ./cmd/app -run TestFrontendSimulation -v
   ```
8. Inspect telemetry and command logs in `logs/` after running. Files rotate automatically when `maxLines` is reached.

## Configuration
`config.json` governs how the system behaves:
* **vehicle**
  * `vehicleID`: identifier stamped on telemetry batches.
* **sensor**
  * `intervalMilliSeconds`: cadence for raw SensorData generation.
  * `minSpeed`, `maxSpeed`, `minPressure`, `maxPressure`, `minTemp`, `maxTemp`: randomization bounds.
  * `ecoMode`, `normalMode`, `speedMode`: relative limits applied when each driving mode is active.
* **processor**
  * `intervalMilliSeconds`: aggregation window for computing averages/min/max.
* **logger**
  * `maxLines`: number of log entries before a new file is created.
  * `fileDir`: root folder for combined, data-only, and command-only `.jsonl` logs.
* **hub**
  * `udpPort`, `tcpPort`, `wsPort`: loopback endpoints used by consumer and frontend.
  * `bufferSize`: byte buffer used by UDP/TCP readers.

Configuration loads once on startup via `config.LoadConfig()`. Update the file and restart to apply changes.

## System Flow
1. **Generator**
   * `Sensor` emits random-but-mode-aware speed, pressure, and temperature readings and reacts to incoming commands.
   * `Process` batches readings for the configured interval, fan-outs calculations across goroutines, and forwards summarized `ResultData`.
2. **Hub**
   * Registers `/api/stream` and upgrades HTTP requests to WebSocket connections.
   * Streams each `ResultData` batch to connected frontend and the consumer (UDP) while duplicating commands to generator (channels) and consumer (TCP).
3. **Consumer**
   * Opens UDP and TCP listeners (signalling readiness through `consumer.Ready`).
   * Differentiates telemetry vs command payloads, then logs each to rotating files with timestamps.
4. **Frontend**
   * Uses a WebSocket hook to connect on demand, show connection status, render the latest metrics, and send predefined commands or custom acceleration values.
   * Provides toast notifications for connect/disconnect, command results, and validation feedback.
5. **Tests**
   * `TestFrontendSimulation` spins up the services, drives a scripted command sequence, captures the WebSocket stream, and persists the interaction under `test/test_logs.jsonl`.

## Development Notes
* The system is fully concurrent, using goroutines and channels for communication.
* Goroutines and channels orchestrate concurrency; `consumer.Ready` ensures network listeners are up before the hub dials TCP/UDP.
* The system is fully concurrent, using goroutines and channels for communication.
* Each transport layer (UDP, TCP, WS) runs independently but shares data via the Hub.
* WebSocket handlers handle graceful close frames and distinguish expected vs unexpected disconnects for cleaner logs.
* Generator speed adjusts based on commands in real time.
* Temperature and pressure are randomly generated based on speed, and they increase or decrease at different rates depending on the mode.
* Frontend tests provide an end-to-end check of the communication pipeline.
* Logs in `.jsonl` format are machine- and human-readable, suitable for further analysis.
* Once the program starts, the vehicle begins sending telemetry automatically, but it must be started and accelerated through commands to simulate motion.