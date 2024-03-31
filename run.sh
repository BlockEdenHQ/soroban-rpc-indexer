timestamp=$(date +"%Y-%m-%d_%H-%M-%S") # Get current timestamp
cd /home/puncsky/soroban-rpc && sudo nohup make dev-ubuntu-mainnet < /dev/null > "output_${timestamp}.log" 2>&1 &
