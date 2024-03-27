#!/bin/bash

echo "Starting..."
timestamp=$(date +"%Y-%m-%d_%H-%M-%S") # Get current timestamp
cd /home/puncsky/soroban-rpc && sudo nohup make dequeue < /dev/null > "consumer_${timestamp}.log" 2>&1 &
