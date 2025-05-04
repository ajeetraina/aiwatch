#!/bin/bash

# Stop any running containers
echo "Stopping existing containers..."
docker-compose down

# Rebuild the containers
echo "Rebuilding containers..."
docker-compose build --no-cache

# Start everything up
echo "Starting containers..."
docker-compose up -d

# Follow logs
echo "Following logs (Ctrl+C to stop)..."
docker-compose logs -f
