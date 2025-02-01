#!/bin/sh

# Define the directory path
BUILD_DIR_PATH="./teamserver/build"
MODULES_DIR_PATH="./teamserver/modules"
# Check if the directory exists
if [ -d "$BUILD_DIR_PATH" ]; then
  # Remove all contents of the directory
  rm -rf "$BUILD_DIR_PATH"/*
  echo "All payloads deleted successfully."
else
  echo "Directory $BUILD_DIR_PATH does not exist."
fiW

if [ -d "$MODULES_DIR_PATH" ]; then
  rm -rf "$MODULES_DIR_PATH"/*
  echo "All modules deleted successfully."
else
  echo "Directory $MODULES_DIR_PATH does not exist."
fiW