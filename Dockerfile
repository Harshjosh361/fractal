# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# Production stage
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/myapp .

# Expose the port on which your app will run
EXPOSE 8080

# Command to run the application
CMD ["./myapp"]