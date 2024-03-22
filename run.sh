#!/bin/bash

while true; do
    # Check if the process does not exist
    if ! ps aux | grep -v grep | grep soroban-indexer > /dev/null; then
        echo "Process not running. Starting..."
        timestamp=$(date +"%Y-%m-%d_%H-%M-%S") # Get current timestamp
        cd /home/puncsky/soroban-indexer && sudo nohup make dev-ubuntu-mainnet < /dev/null > "output_${timestamp}.log" 2>&1 &
    else
        echo "Process running."
    fi
    sleep 10
done
