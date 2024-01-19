#!/bin/bash

# Define the number of servers to start
NUM_SERVERS=16

# Define the range of ports for the servers
START_PORT=5000
END_PORT=$((START_PORT+NUM_SERVERS-1))
rm -rf logs/*.log
for (( port=START_PORT; port<=END_PORT; port++ ))
do
  kill $(lsof -t -i:$port)
done
sleep 2
# Start the servers
for (( port=START_PORT; port<=END_PORT; port++ ))
do
  echo "Starting server on port $port"
  #kill $(lsof -t -i:$port)
  go run cmd/server/main.go $port data/block_info.txt &
done
sleep 5
# Start the client
echo "Starting client"
go run cmd/client/main.go "14-1500"
