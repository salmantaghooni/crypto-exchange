# Dockerfile

# Use the official Golang image as the build stage
FROM golang:1.20-alpine AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to the workspace
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o crypto-exchange main.go

# Start a new stage from scratch
FROM alpine:latest

# Set environment variables
ENV PORT=8080

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/crypto-exchange .

# Copy the config.yaml
COPY --from=builder /app/config.yaml .

# Create logs directory
RUN mkdir -p logs

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./crypto-exchange"]