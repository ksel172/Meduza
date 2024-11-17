#!/bin/bash

# Default mode is dev if TEAMSERVER_MODE is not set
MODE=${TEAMSERVER_MODE:-dev}

if [ "$MODE" = "debug" ]; then

  # Install
  echo "Installing delve..."
  go install github.com/go-delve/delve/cmd/dlv@latest
  echo "Finished installing delve"

  # Build
  # -N == disable optimizations
  # -l == disable inlining
  echo "Building server binary..."
  go build -gcflags="-N -l" -o server ./cmd/main
  echo "Finished building server binary"

  echo "Starting in DEBUG mode with Delve..."
  # Default to port 2345 if DLV_PORT is not set
  DLV_PORT=${DLV_PORT:-2345}
  exec dlv --listen=:"$DLV_PORT" --headless=true --api-version=2 --accept-multiclient exec ./server
else
  echo "Building server binary..."
  go build -o server ./cmd/main
  echo "Finished building server binary"
  echo "Starting server..."
  exec ./server
fi
