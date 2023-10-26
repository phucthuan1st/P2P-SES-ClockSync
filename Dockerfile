# Use a base image with Golang preinstalled
FROM golang:latest

# Set the working directory
WORKDIR /go/src/app

# Install Git
RUN apt-get update && apt-get install -y git

# Clone your project from GitHub
RUN git clone https://github.com/phucthuan1st/p2p-ses-clocksync

# Set the entry point for your application
ENTRYPOINT ["go", "run", "/go/src/app/github.com/phucthuan1st/p2p-ses-clocksync/main.go"]
 
