# Use the official Golang image
FROM golang:1.21.8-alpine3.19 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code to the working directory
COPY . .

# Build the Go application
RUN go build -o /go/bin/app

# Start a new stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the /app directory in the new stage
COPY --from=builder /go/bin/app /app

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./app"]
