#!/bin/bash

# Exit on error
set -e

echo "Setting up and starting AggieCode frontend..."

# Change to frontend directory
cd "$(dirname "$0")/../frontend"

# Install dependencies
echo "Installing frontend dependencies..."
npm install

# Return to root directory
cd ..

# Start the frontend using Bazel
echo "Starting frontend server..."
bazel run //frontend:dev
