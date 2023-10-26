#!/usr/bin/bash

# Build the docker image
docker build -t p2p-ses-clocksync .

# Define the base IP address
base_ip="10.10.10"

# Create 10 containers with sequential IPs from .1 to .10
for i in {1..10}; do
    container_name="Node$i"
    container_ip="$base_ip.$i"
    docker run -d --net p2p --ip "$container_ip" --name "$container_name" p2p-ses-clocksync
done

