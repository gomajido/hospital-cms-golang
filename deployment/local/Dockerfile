# syntax=docker/dockerfile:1

FROM golang:1.22-alpine

# Install necessary packages
RUN apk add --no-cache curl

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Set environment variables for the build
ENV GOOS=linux GOARCH=arm64

# Expose port 8080
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["go", "run", "cmd/http/main.go"]
