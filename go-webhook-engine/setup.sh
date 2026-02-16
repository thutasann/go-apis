#!/bin/bash

# Project root
PROJECT_NAME="webhook-system"

# Create root directory
mkdir -p $PROJECT_NAME
cd $PROJECT_NAME || exit

# Create folders
mkdir -p cmd/server

mkdir -p internal/app
mkdir -p internal/config
mkdir -p internal/delivery/http
mkdir -p internal/domain
mkdir -p internal/repository
mkdir -p internal/queue
mkdir -p internal/worker
mkdir -p internal/service

mkdir -p pkg/logger

# Create empty files
touch cmd/server/main.go

touch internal/app/app.go
touch internal/config/config.go

touch internal/delivery/http/handler.go
touch internal/delivery/http/middleware.go

touch internal/domain/event.go
touch internal/domain/status.go

touch internal/repository/interface.go
touch internal/repository/mongo_event_repository.go

touch internal/queue/interface.go
touch internal/queue/redis_queue.go

touch internal/worker/pool.go
touch internal/worker/worker.go

touch internal/service/processor.go

touch pkg/logger/logger.go

touch go.mod
touch README.md

echo "Webhook system folder structure created successfully!"
