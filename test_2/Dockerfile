# Use the official Golang image as the base image
FROM golang:1.25.5-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Generate swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go --output docs

# Build the Go app with minimal memory usage
RUN CGO_ENABLED=0 GOOS=linux GOGC=50 go build -o main ./cmd/api

# Start a new stage from scratch
FROM alpine:latest

# Install ca-certificates to enable HTTPS requests
RUN apk --no-cache add ca-certificates

RUN apk add --no-cache bash

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
# Copy the docs folder for swagger
COPY --from=builder /app/docs ./docs

# Expose port 8090 to allow communication to/from the blog-api container
EXPOSE 8090

# Command to run the executable
CMD ["./main"]