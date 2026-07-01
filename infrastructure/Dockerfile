FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o notifications ./cmd/notifications/main.go

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/notifications .
# Copy the .env file
COPY .env .

# Expose port 9002 to the outside world
EXPOSE 9002

# Command to run the executable
CMD ["./notifications"]
