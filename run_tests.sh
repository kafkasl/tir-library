#!/bin/bash -e

function stop_server {
    echo "Stopping local server ${APP_PID}"
    kill -9 ${APP_PID}
}

# Run tir-library tests

# Build app binary 
go build -o bin/tir-library -v .

# Clean test cache
go clean -testcache

# Run go tests
# godotenv -f .env.test go test models/* -v
# godotenv -f .env.test go test controllers/* -v

# Run Postman tests

# start local server
./bin/tir-library &
APP_PID=$!

# Set trap to kill local server
trap 'stop_server ${APP_PID}' EXIT

# Run tests
newman run tests/tir-library.postman_collection.json

# online collection
# newman run https://www.getpostman.com/collections/af5407f3fc9af7f9d0e3
