#!/bin/bash

while true; do
    # Check if the process does not exist
    if ! ps aux | grep -v grep | grep queue_consumer > /dev/null; then
        echo "Process not running. Starting..."
        timestamp=$(date +"%Y-%m-%d_%H-%M-%S") # Get current timestamp
        cd /home/puncsky/soroban-rpc && sudo nohup make dequeue < /dev/null > "consumer_${timestamp}.log" 2>&1 &
    else
        echo "Process running."
    fi
    sleep 10
done
