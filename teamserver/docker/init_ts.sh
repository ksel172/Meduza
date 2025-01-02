#!/bin/bash

# Default mode is dev if TEAMSERVER_MODE is not set
MODE=${TEAMSERVER_MODE:-dev}

if [ "$MODE" = "debug" ]; then

  # Install https://github.com/go-delve/delve
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

  # Install https://github.com/air-verse/air for hot reloading
  echo "Installing air..."
  go install github.com/air-verse/air@latest
  echo "Finished installing air"

  echo "Starting server..."
  exec air ./build/server --port "$TEAMSERVER_PORT" -c .air.toml
fi
