# Use a base image with Golang preinstalled
FROM golang:latest

# Set the working directory
WORKDIR /go/src/app

# Copy the source code to the working directory
COPY ./node        /go/src/app//node
COPY ./vectorclock /go/src/app/vectorclock
COPY ./go.mod      /go/src/app/go.mod
COPY ./main.go     /go/src/app/main.go

# Build the source code
RUN go build -o node

ENV NODE_PORT = 9999

ENTRYPOINT [ "./node" ]