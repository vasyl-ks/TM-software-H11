#!/bin/bash

# start.sh - runs backend (Go) and frontend (Vite) for local development

set -e

echo "Starting backend..."
go run ./cmd/app/main.go &
BACK_PID=$!

echo "Starting frontend..."
(
  cd frontend
  npm run dev
)

echo "Shutting down backend..."
kill $BACK_PID 2>/dev/null || true