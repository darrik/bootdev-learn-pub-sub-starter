#!/bin/bash

case "$1" in
    start)
        echo "Starting RabbitMQ container..."
        podman run -d --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.13-management
        ;;
    stop)
        echo "Stopping RabbitMQ container..."
        podman stop rabbitmq
        ;;
    logs)
        echo "Fetching logs for RabbitMQ container..."
        podman logs -f rabbitmq
        ;;
    *)
        echo "Usage: $0 {start|stop|logs}"
        exit 1
esac
