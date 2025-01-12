#!/bin/sh

# Define the directory path
DIR_PATH="./teamserver/build"

# Check if the directory exists
if [ -d "$DIR_PATH" ]; then
  # Remove all contents of the directory
  rm -rf "$DIR_PATH"/*
  echo "All payloads deleted successfully."
else
  echo "Directory $DIR_PATH does not exist."
fi