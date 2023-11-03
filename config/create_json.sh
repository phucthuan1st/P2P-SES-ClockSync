#!/bin/bash

# Define the list of node IDs from A to J
node_ids=("A" "B" "C" "D" "E" "F" "G" "H" "I" "J")

# Loop through each node
for ((i = 0; i < 10; i++)); do
  node_id=${node_ids[$i]}
  node_ip="127.0.0.1"
  node_port=$((9000 + i))

  # Create an array for the nodes to connect to
  connect_to=()

  # Loop to establish connections as described
  for ((j = i + 1; j < 10; j++)); do
    connect_node_id=${node_ids[$j]}
    connect_node_ip="127.0.0.1"
    connect_node_port=$((9000 + j))
    connect_to+=("{\"nodeID\": \"$connect_node_id\", \"nodeIP\": \"$connect_node_ip\", \"nodePort\": $connect_node_port}")
  done

  # Convert the array to JSON format
  connect_to_json=$(IFS=,; echo "[${connect_to[*]}]")

  # Create the JSON configuration file
  config='{
    "nodeID": "'$node_id'",
    "nodeIP": "'$node_ip'",
    "nodePort": '$node_port',
    "connectTo": '$connect_to_json'
  }'

  echo $config > node${node_id}.json
  echo "Created node${node_id}.json"
done
